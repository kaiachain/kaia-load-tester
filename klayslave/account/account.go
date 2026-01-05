package account

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	kaia "github.com/kaiachain/kaia"
	"github.com/kaiachain/kaia/accounts/abi"
	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/api"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/types/accountkey"
	"github.com/kaiachain/kaia/client"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax/auction"
	auctionImpl "github.com/kaiachain/kaia/kaiax/auction/impl"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
)

const Letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

var (
	gasPrice *big.Int
	chainID  *big.Int
	baseFee  *big.Int
)

type Account struct {
	id                 int
	privateKey         []*ecdsa.PrivateKey
	key                []string
	address            common.Address
	nonce              uint64
	balance            *big.Int
	mutex              sync.Mutex
	lastBlocknumSentTx uint64
}

type Client interface {
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
	SendRawTransaction(ctx context.Context, tx *types.Transaction) (common.Hash, error)
	CreateAccessList(ctx context.Context, callMsg kaia.CallMsg) (*types.AccessList, uint64, string, error)
	TransactionReceiptRpcOutput(ctx context.Context, txHash common.Hash) (r map[string]interface{}, err error)
}

// TxSendFunc is a function type that sends a transaction
type TxSendFunc func() (*types.Transaction, error)

// RetryConfig configures retry behavior
type RetryConfig struct {
	SendRetryInterval time.Duration
	WaitMinedTimeout  time.Duration
	// ShouldSkip returns true if the error should skip retry and return immediately
	ShouldSkip func(error) bool
}

// DefaultRetryConfig returns default retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		SendRetryInterval: 1 * time.Second,
		WaitMinedTimeout:  60 * time.Second,
		ShouldSkip:        nil,
	}
}

// RunWithRetry executes a transaction with retry logic for both send and WaitMined failures
func (a *Account) RunWithRetry(c *client.Client, config RetryConfig, sendTx TxSendFunc) *types.Transaction {
	for {
		var lastTx *types.Transaction
		var err error

		// Inner loop: retry sending transaction
		for {
			lastTx, err = sendTx()
			if err == nil {
				break
			}
			// Check if this error should skip retry
			if config.ShouldSkip != nil && config.ShouldSkip(err) {
				return lastTx
			}
			log.Printf("Failed to send tx: err=%s", err.Error())
			time.Sleep(config.SendRetryInterval)
		}

		// WaitMined with timeout
		ctx, cancelFn := context.WithTimeout(context.Background(), config.WaitMinedTimeout)
		receipt, err := bind.WaitMined(ctx, c, lastTx)
		cancelFn()

		if err == nil && receipt != nil && receipt.Status == 1 {
			return lastTx // success
		}

		// WaitMined failed or tx reverted - reset nonce and retry
		log.Printf("WaitMined failed or tx reverted, retrying: err=%v, txHash=%s", err, lastTx.Hash().String())
		a.nonce = 0 // Reset nonce to fetch fresh nonce on retry
		time.Sleep(1 * time.Second)
	}
}

// TxValues is a type alias for transaction value map
type TxValues = map[types.TxValueKeyType]interface{}

// TxOption configures transaction sending behavior
type TxOption struct {
	FeePayer *Account // If set, sign with fee payer
	NoLock   bool     // If true, skip mutex lock (caller handles locking)
}

// sendRawTx sends an already created and signed transaction
// It handles SendRawTransaction and nonce error handling
func (self *Account) sendRawTx(c Client, tx *types.Transaction, nonce uint64) error {
	_, err := c.SendRawTransaction(context.Background(), tx)
	if err != nil {
		fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			self.nonce++
		}
		return err
	}
	self.nonce++
	return nil
}

// sendTransaction is a helper that handles the common transaction sending pattern
func (self *Account) sendTransaction(c Client, txType types.TxType, values TxValues, opt *TxOption) (*types.Transaction, error) {
	ctx := context.Background()

	if opt == nil || !opt.NoLock {
		self.mutex.Lock()
		defer self.mutex.Unlock()
	}

	nonce := self.GetNonce(c)

	// Set common values
	values[types.TxValueKeyNonce] = nonce
	values[types.TxValueKeyFrom] = self.address
	if _, ok := values[types.TxValueKeyGasPrice]; !ok {
		values[types.TxValueKeyGasPrice] = gasPrice
	}

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(txType, values)
	if err != nil {
		return nil, fmt.Errorf("failed to create tx: %w", err)
	}

	if err = tx.SignWithKeys(signer, self.privateKey); err != nil {
		return nil, fmt.Errorf("failed to sign tx: %w", err)
	}

	// Fee payer signing if needed
	if opt != nil && opt.FeePayer != nil {
		if err = tx.SignFeePayerWithKeys(signer, opt.FeePayer.privateKey); err != nil {
			return nil, fmt.Errorf("failed to fee payer sign tx: %w", err)
		}
	}

	_, err = c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return tx, err
	}

	self.nonce++
	return tx, nil
}

func init() {
	gasPrice = big.NewInt(0)
	chainID = big.NewInt(2018)
	baseFee = big.NewInt(0)
}

func SetGasPrice(gp *big.Int) {
	gasPrice = gp
}

func SetBaseFee(bf *big.Int) {
	baseFee = bf
}

func SetChainID(id *big.Int) {
	chainID = id
}

func (acc *Account) Lock() {
	acc.mutex.Lock()
}

func (acc *Account) UnLock() {
	acc.mutex.Unlock()
}

func GetAccountFromKey(id int, key string) *Account {
	acc, err := crypto.HexToECDSA(key)
	if err != nil {
		log.Fatalf("Key(%v): Failed to HexToECDSA %v", key, err)
	}

	return &Account{
		id:         0,
		privateKey: []*ecdsa.PrivateKey{acc},
		key:        []string{key},
		address:    crypto.PubkeyToAddress(acc.PublicKey),
		nonce:      0,
		balance:    big.NewInt(0),
	}
}

func (account *Account) ImportUnLockAccount(endpoint string) {
	key := account.key[0]
	acc, err := crypto.HexToECDSA(key)
	if err != nil {
		log.Fatalf("Key(%v): Failed to HexToECDSA %v", err)
	}

	testAddr := crypto.PubkeyToAddress(acc.PublicKey)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c, err := client.Dial(endpoint)
	if err != nil {
		log.Fatalf("ImportUnLockAccount(): Create Client %v", err)
	}

	addr, err := c.ImportRawKey(ctx, key, "")
	if err != nil {
		log.Fatalf("Account(%v) : Failed to import => %v\n", account.address, err)
	} else {
		if testAddr != addr {
			log.Fatalf("origial:%v, imported: %v\n", testAddr.String(), addr.String())
		}
	}

	res, err := c.UnlockAccount(ctx, account.address, "", 0)
	if err != nil {
		log.Fatalf("Account(%v) : Failed to Unlock: %v\n", account.address.String(), err)
	} else {
		log.Printf("Wallet UnLock Result: %v", res)
	}
}

func NewAccount(id int) *Account {
	acc, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("crypto.GenerateKey() : Failed to generateKey %v", err)
	}

	testKey := hex.EncodeToString(crypto.FromECDSA(acc))

	return &Account{
		id:         0,
		privateKey: []*ecdsa.PrivateKey{acc},
		key:        []string{testKey},
		address:    crypto.PubkeyToAddress(acc.PublicKey),
		nonce:      0,
		balance:    big.NewInt(0),
	}
}

func NewAccountOnNode(id int, endpoint string) *Account {

	tAcc := NewAccount(id)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c, err := client.Dial(endpoint)
	if err != nil {
		log.Fatalf("NewAccountOnNode() : Failed to create client %v", err)
	}

	addr, err := c.ImportRawKey(ctx, tAcc.key[0], "")
	if err != nil {
		//log.Printf("Account(%v) : Failed to import\n", tAcc.address, err)
	} else {
		if tAcc.address != addr {
			log.Fatalf("origial:%v, imported: %v\n", tAcc.address, addr.String())
		}
		//log.Printf("origial:%v, imported:%v\n", tAcc.address, addr.String())
	}

	_, err = c.UnlockAccount(ctx, tAcc.GetAddress(), "", 0)
	if err != nil {
		log.Printf("Account(%v) : Failed to Unlock: %v\n", tAcc.GetAddress().String(), err)
	}

	//log.Printf("Wallet UnLock Result: %v", flag)

	return tAcc
}

func NewKaiaAccount(id int) *Account {
	acc, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("crypto.GenerateKey() : Failed to generateKey %v", err)
	}

	testKey := hex.EncodeToString(crypto.FromECDSA(acc))
	randomAddr := common.BytesToAddress(crypto.Keccak256([]byte(testKey))[12:])

	return &Account{
		id:         0,
		privateKey: []*ecdsa.PrivateKey{acc},
		key:        []string{testKey},
		address:    randomAddr,
		nonce:      0,
		balance:    big.NewInt(0),
	}
}

func NewKaiaAccountWithAddr(id int, addr common.Address) *Account {
	acc, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("crypto.GenerateKey() : Failed to generateKey %v", err)
	}

	testKey := hex.EncodeToString(crypto.FromECDSA(acc))

	return &Account{
		id:         0,
		privateKey: []*ecdsa.PrivateKey{acc},
		key:        []string{testKey},
		address:    addr,
		nonce:      0,
		balance:    big.NewInt(0),
	}
}

func NewKaiaMultisigAccount(id int) *Account {
	k1, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("crypto.GenerateKey() : Failed to generateKey %v", err)
	}
	k2, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("crypto.GenerateKey() : Failed to generateKey %v", err)
	}
	k3, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("crypto.GenerateKey() : Failed to generateKey %v", err)
	}

	testKey := hex.EncodeToString(crypto.FromECDSA(k1))
	randomAddr := common.BytesToAddress(crypto.Keccak256([]byte(testKey))[12:])

	return &Account{
		id:         0,
		privateKey: []*ecdsa.PrivateKey{k1, k2, k3},
		key:        []string{testKey},
		address:    randomAddr,
		nonce:      0,
		balance:    big.NewInt(0),
	}
}

func UnlockAccount(c *client.Client, addr common.Address, pwd string) {
	ctx := context.Background()
	defer ctx.Done()

	if _, err := c.UnlockAccount(ctx, addr, pwd, 0); err != nil {
		fmt.Println(err)
	}
}

func (acc *Account) GetKey() *ecdsa.PrivateKey {
	return acc.privateKey[0]
}

func (acc *Account) GetAddress() common.Address {
	return acc.address
}

func (acc *Account) GetPrivateKey() string {
	return acc.key[0]
}

func (acc *Account) GetNonce(c Client) uint64 {
	if acc.nonce != 0 {
		return acc.nonce
	}
	ctx := context.Background()
	nonce, err := c.NonceAt(ctx, acc.GetAddress(), nil)
	if err != nil {
		log.Printf("GetNonce(): Failed to NonceAt() %v\n", err)
		return acc.nonce
	}
	acc.nonce = nonce

	//fmt.Printf("account= %v  nonce = %v\n", acc.GetAddress().String(), nonce)
	return acc.nonce
}

func (acc *Account) GetNonceFromBlock(c Client) uint64 {
	ctx := context.Background()
	nonce, err := c.NonceAt(ctx, acc.GetAddress(), nil)
	if err != nil {
		log.Printf("GetNonce(): Failed to NonceAt() %v\n", err)
		return acc.nonce
	}

	acc.nonce = nonce

	fmt.Printf("%v: account= %v  nonce = %v\n", os.Getpid(), acc.GetAddress().String(), nonce)
	return acc.nonce
}

func (acc *Account) UpdateNonce() {
	acc.nonce++
}

func (a *Account) GetReceipt(c *client.Client, txHash common.Hash) (*types.Receipt, error) {
	ctx := context.Background()
	return c.TransactionReceipt(ctx, txHash)
}

func (a *Account) GetBalance(c *client.Client) (*big.Int, error) {
	ctx := context.Background()
	balance, err := c.BalanceAt(ctx, a.GetAddress(), nil)
	if err != nil {
		return nil, err
	}
	return balance, err
}

func (self *Account) TransferSignedTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	tx, gasPrice, err := self.TransferSignedTxReturnTx(true, c, to, value)
	return tx.Hash(), gasPrice, err
}

func (self *Account) TransferSignedTxWithGuaranteeRetry(c *client.Client, to *Account, value *big.Int) *types.Transaction {
	config := DefaultRetryConfig()
	config.WaitMinedTimeout = 30 * time.Second
	return self.RunWithRetry(c, config, func() (*types.Transaction, error) {
		tx, _, err := self.TransferSignedTxReturnTx(true, c, to, value)
		return tx, err
	})
}

func (self *Account) TransferSignedTxWithoutLock(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	tx, gasPrice, err := self.TransferSignedTxReturnTx(false, c, to, value)
	return tx.Hash(), gasPrice, err
}

func (self *Account) TransferSignedTxReturnTx(withLock bool, c *client.Client, to *Account, value *big.Int) (*types.Transaction, *big.Int, error) {
	if withLock {
		self.mutex.Lock()
		defer self.mutex.Unlock()
	}

	nonce := self.GetNonce(c)
	tx := types.NewTransaction(nonce, to.GetAddress(), value, 21000, gasPrice, nil)
	signTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), self.privateKey[0])
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = self.sendRawTx(c, signTx, nonce)
	return signTx, tx.GasPrice(), err
}

func (self *Account) TransferNewValueTransferWithCancelTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	var txList []*types.Transaction
	nonce := self.GetNonce(c)

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeValueTransfer, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyTo:       to.GetAddress(),
		types.TxValueKeyAmount:   value,
		types.TxValueKeyGasLimit: uint64(100000),
		types.TxValueKeyGasPrice: gasPrice,
		types.TxValueKeyFrom:     self.address,
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	txList = append(txList, tx)

	cancelTx, err := types.NewTransactionWithMap(types.TxTypeCancel, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyFrom:     self.address,
		types.TxValueKeyGasLimit: uint64(100000000),
		types.TxValueKeyGasPrice: gasPrice,
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = cancelTx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	txList = append(txList, cancelTx)

	var hash common.Hash
	for _, tx := range txList {
		hash, err := c.SendRawTransaction(ctx, tx)
		if err != nil {
			if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
				fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
				fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
				self.nonce++
			} else {
				fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			}
			return hash, gasPrice, err
		}
	}

	self.nonce++
	return hash, gasPrice, nil
}

func (self *Account) TransferNewValueTransferTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	tx, err := self.sendTransaction(c, types.TxTypeValueTransfer, TxValues{
		types.TxValueKeyTo:       to.GetAddress(),
		types.TxValueKeyAmount:   value,
		types.TxValueKeyGasLimit: uint64(100000),
	}, nil)
	if tx != nil {
		return tx.Hash(), gasPrice, err
	}
	return common.Hash{}, gasPrice, err
}

func (self *Account) TransferNewFeeDelegatedValueTransferTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	tx, err := self.sendTransaction(c, types.TxTypeFeeDelegatedValueTransfer, TxValues{
		types.TxValueKeyTo:       to.GetAddress(),
		types.TxValueKeyAmount:   value,
		types.TxValueKeyGasLimit: uint64(100000),
		types.TxValueKeyFeePayer: to.address,
	}, &TxOption{FeePayer: to})
	if tx != nil {
		return tx.Hash(), gasPrice, err
	}
	return common.Hash{}, gasPrice, err
}

func (self *Account) TransferNewFeeDelegatedValueTransferWithRatioTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	tx, err := self.sendTransaction(c, types.TxTypeFeeDelegatedValueTransferWithRatio, TxValues{
		types.TxValueKeyTo:                 to.GetAddress(),
		types.TxValueKeyAmount:             value,
		types.TxValueKeyGasLimit:           uint64(100000),
		types.TxValueKeyFeePayer:           to.address,
		types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(30),
	}, &TxOption{FeePayer: to})
	if tx != nil {
		return tx.Hash(), gasPrice, err
	}
	return common.Hash{}, gasPrice, err
}

// createMemoTransferTx is a helper function that creates and sends a memo transfer transaction
func (self *Account) createMemoTransferTx(c *client.Client, to *Account, value *big.Int, data []byte, gasLimit uint64) (common.Hash, *big.Int, error) {
	tx, err := self.sendTransaction(c, types.TxTypeValueTransferMemo, TxValues{
		types.TxValueKeyTo:       to.GetAddress(),
		types.TxValueKeyAmount:   value,
		types.TxValueKeyGasLimit: gasLimit,
		types.TxValueKeyData:     data,
	}, nil)
	if tx != nil {
		return tx.Hash(), gasPrice, err
	}
	return common.Hash{}, gasPrice, err
}

func (self *Account) TransferNewValueTransferMemoTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	data := []byte("hello")
	return self.createMemoTransferTx(c, to, value, data, 150000)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

// increase memo size from 5 bytes to between 50 bytes and 2,000 bytes

func (self *Account) TransferNewValueTransferBigRandomStringMemoTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	minBytes := 50
	maxBytes := 2000
	data := []byte(randomString(randInt(minBytes, maxBytes)))
	return self.createMemoTransferTx(c, to, value, data, 200000)
}

// create 200 strings of memo
func (self *Account) TransferNewValueTransferSmallMemoTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	length := 200
	data := []byte(randomString(length))
	return self.createMemoTransferTx(c, to, value, data, 150000)
}

// create 2000 strings of memo
func (self *Account) TransferNewValueTransferLargeMemoTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	length := 2000
	data := []byte(randomString(length))
	return self.createMemoTransferTx(c, to, value, data, 200000)
}

func (self *Account) TransferNewFeeDelegatedValueTransferMemoTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	tx, err := self.sendTransaction(c, types.TxTypeFeeDelegatedValueTransferMemo, TxValues{
		types.TxValueKeyTo:       to.GetAddress(),
		types.TxValueKeyAmount:   value,
		types.TxValueKeyGasLimit: uint64(100000),
		types.TxValueKeyData:     []byte("hello"),
		types.TxValueKeyFeePayer: to.address,
	}, &TxOption{FeePayer: to})
	if tx != nil {
		return tx.Hash(), gasPrice, err
	}
	return common.Hash{}, gasPrice, err
}

func (self *Account) TransferNewFeeDelegatedValueTransferMemoWithRatioTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	tx, err := self.sendTransaction(c, types.TxTypeFeeDelegatedValueTransferMemoWithRatio, TxValues{
		types.TxValueKeyTo:                 to.GetAddress(),
		types.TxValueKeyAmount:             value,
		types.TxValueKeyGasLimit:           uint64(100000),
		types.TxValueKeyData:               []byte("hello"),
		types.TxValueKeyFeePayer:           to.address,
		types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(30),
	}, &TxOption{FeePayer: to})
	if tx != nil {
		return tx.Hash(), gasPrice, err
	}
	return common.Hash{}, gasPrice, err
}

func (self *Account) TransferNewAccountCreationTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	tx, err := self.sendTransaction(c, types.TxTypeAccountCreation, TxValues{
		types.TxValueKeyTo:            to.GetAddress(),
		types.TxValueKeyAmount:        value,
		types.TxValueKeyGasLimit:      uint64(1000000),
		types.TxValueKeyHumanReadable: false,
		types.TxValueKeyAccountKey:    accountkey.NewAccountKeyPublicWithValue(&to.privateKey[0].PublicKey),
	}, nil)
	if tx != nil {
		return tx.Hash(), gasPrice, err
	}
	return common.Hash{}, gasPrice, err
}

func (self *Account) TransferNewAccountUpdateTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	tx, err := self.sendTransaction(c, types.TxTypeAccountUpdate, TxValues{
		types.TxValueKeyGasLimit:   uint64(100000),
		types.TxValueKeyAccountKey: accountkey.NewAccountKeyPublicWithValue(&self.privateKey[0].PublicKey),
	}, nil)
	if tx != nil {
		return tx.Hash(), gasPrice, err
	}
	return common.Hash{}, gasPrice, err
}

func (self *Account) TransferNewFeeDelegatedAccountUpdateTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	tx, err := self.sendTransaction(c, types.TxTypeFeeDelegatedAccountUpdate, TxValues{
		types.TxValueKeyGasLimit:   uint64(100000),
		types.TxValueKeyAccountKey: accountkey.NewAccountKeyPublicWithValue(&self.privateKey[0].PublicKey),
		types.TxValueKeyFeePayer:   to.address,
	}, &TxOption{FeePayer: to})
	if tx != nil {
		return tx.Hash(), gasPrice, err
	}
	return common.Hash{}, gasPrice, err
}

func (self *Account) TransferNewFeeDelegatedAccountUpdateWithRatioTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	tx, err := self.sendTransaction(c, types.TxTypeFeeDelegatedAccountUpdateWithRatio, TxValues{
		types.TxValueKeyGasLimit:           uint64(100000),
		types.TxValueKeyAccountKey:         accountkey.NewAccountKeyPublicWithValue(&self.privateKey[0].PublicKey),
		types.TxValueKeyFeePayer:           to.address,
		types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(30),
	}, &TxOption{FeePayer: to})
	if tx != nil {
		return tx.Hash(), gasPrice, err
	}
	return common.Hash{}, gasPrice, err
}

func (self *Account) TransferNewSmartContractDeployTx(c *client.Client, to *Account, value *big.Int, data []byte, shouldFixNonceZero bool) (common.Address, *types.Transaction, *big.Int, error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)
	if shouldFixNonceZero {
		nonce = 0
	}

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractDeploy, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:         nonce,
		types.TxValueKeyFrom:          self.address,
		types.TxValueKeyTo:            (*common.Address)(nil),
		types.TxValueKeyAmount:        value,
		types.TxValueKeyGasLimit:      uint64(10000000),
		types.TxValueKeyGasPrice:      gasPrice,
		types.TxValueKeyHumanReadable: false,
		types.TxValueKeyCodeFormat:    params.CodeFormatEVM,
		types.TxValueKeyData:          data,
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}
	if err = tx.SignWithKeys(signer, self.privateKey); err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	_, err = c.SendRawTransaction(context.Background(), tx)
	if err != nil {
		fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		if !shouldFixNonceZero && (err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error()) {
			self.nonce++
		}
		return common.Address{}, tx, gasPrice, err
	}

	contractAddr := crypto.CreateAddress(self.address, self.nonce)
	self.nonce++
	return contractAddr, tx, gasPrice, nil
}

func (self *Account) TransferNewFeeDelegatedSmartContractDeployTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	code := "0x608060405234801561001057600080fd5b506101de806100206000396000f3006080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416631a39d8ef81146100805780636353586b146100a757806370a08231146100ca578063fd6b7ef8146100f8575b3360009081526001602052604081208054349081019091558154019055005b34801561008c57600080fd5b5061009561010d565b60408051918252519081900360200190f35b6100c873ffffffffffffffffffffffffffffffffffffffff60043516610113565b005b3480156100d657600080fd5b5061009573ffffffffffffffffffffffffffffffffffffffff60043516610147565b34801561010457600080fd5b506100c8610159565b60005481565b73ffffffffffffffffffffffffffffffffffffffff1660009081526001602052604081208054349081019091558154019055565b60016020526000908152604090205481565b336000908152600160205260408120805490829055908111156101af57604051339082156108fc029083906000818181858888f193505050501561019c576101af565b3360009081526001602052604090208190555b505600a165627a7a72305820627ca46bb09478a015762806cc00c431230501118c7c26c30ac58c4e09e51c4f0029"

	tx, err := self.sendTransaction(c, types.TxTypeFeeDelegatedSmartContractDeploy, TxValues{
		types.TxValueKeyTo:            (*common.Address)(nil),
		types.TxValueKeyAmount:        common.Big0,
		types.TxValueKeyGasLimit:      uint64(10000000),
		types.TxValueKeyHumanReadable: false,
		types.TxValueKeyData:          common.FromHex(code),
		types.TxValueKeyCodeFormat:    params.CodeFormatEVM,
		types.TxValueKeyFeePayer:      self.address,
	}, &TxOption{FeePayer: self})
	if tx != nil {
		return tx.Hash(), gasPrice, err
	}
	return common.Hash{}, gasPrice, err
}

func (self *Account) TransferNewFeeDelegatedSmartContractDeployWithRatioTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	code := "0x608060405234801561001057600080fd5b506101de806100206000396000f3006080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416631a39d8ef81146100805780636353586b146100a757806370a08231146100ca578063fd6b7ef8146100f8575b3360009081526001602052604081208054349081019091558154019055005b34801561008c57600080fd5b5061009561010d565b60408051918252519081900360200190f35b6100c873ffffffffffffffffffffffffffffffffffffffff60043516610113565b005b3480156100d657600080fd5b5061009573ffffffffffffffffffffffffffffffffffffffff60043516610147565b34801561010457600080fd5b506100c8610159565b60005481565b73ffffffffffffffffffffffffffffffffffffffff1660009081526001602052604081208054349081019091558154019055565b60016020526000908152604090205481565b336000908152600160205260408120805490829055908111156101af57604051339082156108fc029083906000818181858888f193505050501561019c576101af565b3360009081526001602052604090208190555b505600a165627a7a72305820627ca46bb09478a015762806cc00c431230501118c7c26c30ac58c4e09e51c4f0029"

	tx, err := self.sendTransaction(c, types.TxTypeFeeDelegatedSmartContractDeployWithRatio, TxValues{
		types.TxValueKeyTo:                 (*common.Address)(nil),
		types.TxValueKeyAmount:             common.Big0,
		types.TxValueKeyGasLimit:           uint64(10000000),
		types.TxValueKeyHumanReadable:      false,
		types.TxValueKeyData:               common.FromHex(code),
		types.TxValueKeyFeePayer:           self.address,
		types.TxValueKeyCodeFormat:         params.CodeFormatEVM,
		types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(30),
	}, &TxOption{FeePayer: self})
	if tx != nil {
		return tx.Hash(), gasPrice, err
	}
	return common.Hash{}, gasPrice, err
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = Letters[r.Intn(len(Letters))]
	}
	return string(b)
}

func (self *Account) ExecuteStorageTrieStore(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))

	abiStr := `[{"constant":true,"inputs":[],"name":"rootCaCertificate","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"_serialNumber","type":"string"}],"name":"getIdentity","outputs":[{"name":"","type":"string"},{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_caKey","type":"string"}],"name":"deleteCaCertificate","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_caKey","type":"string"},{"name":"_caCert","type":"string"}],"name":"insertCaCertificate","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"owner","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_serialNumber","type":"string"},{"name":"_publicKey","type":"string"},{"name":"_hash","type":"string"}],"name":"insertIdentity","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_serialNumber","type":"string"}],"name":"deleteIdentity","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"_caKey","type":"string"}],"name":"getCaCertificate","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[],"payable":false,"stateMutability":"nonpayable","type":"constructor"}]`

	abii, err := abi.JSON(strings.NewReader(string(abiStr)))
	if err != nil {
		log.Fatalf("failed to abi.JSON: %v", err)
	}
	data, err := abii.Pack("insertIdentity", randomString(39), randomString(814), randomString(40))
	if err != nil {
		log.Fatalf("failed to abi.Pack: %v", err)
	}

	tx, err := self.sendTransaction(c, types.TxTypeSmartContractExecution, TxValues{
		types.TxValueKeyGasLimit: uint64(5000000),
		types.TxValueKeyAmount:   common.Big0,
		types.TxValueKeyTo:       to.address,
		types.TxValueKeyData:     data,
	}, nil)
	if tx != nil {
		return tx.Hash(), gasPrice, err
	}
	return common.Hash{}, gasPrice, err
}

func (self *Account) TransferNewSmartContractExecutionTx(c *client.Client, to *Account, value *big.Int, data []byte) (*types.Transaction, *big.Int, error) {
	if value == nil {
		value = big.NewInt(0)
	}
	tx, err := self.sendTransaction(c, types.TxTypeSmartContractExecution, TxValues{
		types.TxValueKeyGasLimit: uint64(5000000),
		types.TxValueKeyAmount:   value,
		types.TxValueKeyTo:       to.address,
		types.TxValueKeyData:     data,
	}, nil)
	return tx, gasPrice, err
}

func (self *Account) TransferNewFeeDelegatedSmartContractExecutionTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	abiStr := `[{"constant":true,"inputs":[],"name":"totalAmount","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"receiver","type":"address"}],"name":"reward","outputs":[],"payable":true,"stateMutability":"payable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"safeWithdrawal","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"inputs":[],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"payable":true,"stateMutability":"payable","type":"fallback"}]`

	abii, err := abi.JSON(strings.NewReader(string(abiStr)))
	if err != nil {
		log.Fatalf("failed to abi.JSON: %v", err)
	}

	data, err := abii.Pack("reward", self.address)
	if err != nil {
		log.Fatalf("failed to abi.Pack: %v", err)
	}

	tx, err := self.sendTransaction(c, types.TxTypeFeeDelegatedSmartContractExecution, TxValues{
		types.TxValueKeyGasLimit: uint64(5000000),
		types.TxValueKeyAmount:   value,
		types.TxValueKeyTo:       to.address,
		types.TxValueKeyData:     data,
		types.TxValueKeyFeePayer: self.address,
	}, &TxOption{FeePayer: self})
	if tx != nil {
		return tx.Hash(), gasPrice, err
	}
	return common.Hash{}, gasPrice, err
}

func (self *Account) TransferNewFeeDelegatedSmartContractExecutionWithRatioTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	abiStr := `[{"constant":true,"inputs":[],"name":"totalAmount","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"receiver","type":"address"}],"name":"reward","outputs":[],"payable":true,"stateMutability":"payable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"safeWithdrawal","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"inputs":[],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"payable":true,"stateMutability":"payable","type":"fallback"}]`

	abii, err := abi.JSON(strings.NewReader(string(abiStr)))
	if err != nil {
		log.Fatalf("failed to abi.JSON: %v", err)
	}

	data, err := abii.Pack("reward", self.address)
	if err != nil {
		log.Fatalf("failed to abi.Pack: %v", err)
	}

	tx, err := self.sendTransaction(c, types.TxTypeFeeDelegatedSmartContractExecutionWithRatio, TxValues{
		types.TxValueKeyGasLimit:           uint64(5000000),
		types.TxValueKeyAmount:             value,
		types.TxValueKeyTo:                 to.address,
		types.TxValueKeyData:               data,
		types.TxValueKeyFeePayer:           self.address,
		types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(30),
	}, &TxOption{FeePayer: self})
	if tx != nil {
		return tx.Hash(), gasPrice, err
	}
	return common.Hash{}, gasPrice, err
}

func (self *Account) TransferNewCancelTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	tx, err := self.sendTransaction(c, types.TxTypeCancel, TxValues{
		types.TxValueKeyGasLimit: uint64(100000000),
	}, nil)
	if tx != nil {
		return tx.Hash(), gasPrice, err
	}
	return common.Hash{}, gasPrice, err
}

func (self *Account) TransferNewFeeDelegatedCancelTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	tx, err := self.sendTransaction(c, types.TxTypeFeeDelegatedCancel, TxValues{
		types.TxValueKeyGasLimit: uint64(100000000),
		types.TxValueKeyFeePayer: to.address,
	}, &TxOption{FeePayer: to})
	if tx != nil {
		return tx.Hash(), gasPrice, err
	}
	return common.Hash{}, gasPrice, err
}

func (self *Account) TransferNewFeeDelegatedCancelWithRatioTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	tx, err := self.sendTransaction(c, types.TxTypeFeeDelegatedCancelWithRatio, TxValues{
		types.TxValueKeyGasLimit:           uint64(100000000),
		types.TxValueKeyFeePayer:           to.address,
		types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(30),
	}, &TxOption{FeePayer: to})
	if tx != nil {
		return tx.Hash(), gasPrice, err
	}
	return common.Hash{}, gasPrice, err
}

// EthTxType defines the type of Ethereum transaction
type EthTxType int

const (
	EthTxTypeLegacy EthTxType = iota
	EthTxTypeAccessList
	EthTxTypeDynamicFee
)

// transferNewEthereumTxWithAccessList is a helper for Ethereum tx types that use CreateAccessList
func (self *Account) transferNewEthereumTxWithAccessList(c Client, to *Account, value *big.Int, input []byte, txType EthTxType) (common.Hash, *big.Int, error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)
	gas := uint64(5000000)
	var toAddress *common.Address
	if to != nil {
		toAddress = &to.address
	}

	accessList, _, _, err := c.CreateAccessList(context.Background(), kaia.CallMsg{
		From: self.address, To: toAddress, Gas: gas, GasPrice: gasPrice, Value: value, Data: input,
	})
	if err != nil {
		log.Fatalf("Failed to get accessList: %v", err)
	}

	var tx *types.Transaction
	switch txType {
	case EthTxTypeAccessList:
		tx = types.NewTx(&types.TxInternalDataEthereumAccessList{
			ChainID: chainID, AccountNonce: nonce, Recipient: toAddress,
			GasLimit: gas, Price: gasPrice, Amount: value, AccessList: *accessList, Payload: input,
		})
	case EthTxTypeDynamicFee:
		tx = types.NewTx(&types.TxInternalDataEthereumDynamicFee{
			ChainID: chainID, AccountNonce: nonce, Recipient: toAddress,
			GasLimit: gas, GasFeeCap: gasPrice, GasTipCap: gasPrice, Amount: value, AccessList: *accessList, Payload: input,
		})
	}

	if err = tx.SignWithKeys(types.LatestSignerForChainID(chainID), self.privateKey); err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	err = self.sendRawTx(c, tx, nonce)
	return tx.Hash(), gasPrice, err
}

func (self *Account) TransferNewEthereumAccessListTx(c Client, to *Account, value *big.Int, input []byte) (common.Hash, *big.Int, error) {
	return self.transferNewEthereumTxWithAccessList(c, to, value, input, EthTxTypeAccessList)
}

func (self *Account) TransferNewEthereumDynamicFeeTx(c Client, to *Account, value *big.Int, input []byte) (common.Hash, *big.Int, error) {
	return self.transferNewEthereumTxWithAccessList(c, to, value, input, EthTxTypeDynamicFee)
}

// transferNewEthStyleTx is a helper for simple Ethereum-style transactions (without CreateAccessList)
func (self *Account) transferNewEthStyleTx(c Client, to *Account, value *big.Int, input []byte, txType EthTxType) (common.Hash, *big.Int, error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)
	gas := uint64(100000)
	var tx *types.Transaction

	if to == nil {
		tx = types.NewContractCreation(nonce, value, gas*2, gasPrice, input)
	} else {
		switch txType {
		case EthTxTypeLegacy:
			tx = types.NewTransaction(nonce, to.address, value, gas, gasPrice, input)
		case EthTxTypeAccessList:
			tx = types.NewTx(&types.TxInternalDataEthereumAccessList{
				ChainID: chainID, AccountNonce: nonce, Recipient: &to.address,
				GasLimit: gas, Price: gasPrice, Amount: value, AccessList: types.AccessList{}, Payload: input,
			})
		case EthTxTypeDynamicFee:
			tx = types.NewTx(&types.TxInternalDataEthereumDynamicFee{
				ChainID: chainID, AccountNonce: nonce, Recipient: &to.address,
				GasLimit: gas, GasFeeCap: gasPrice, GasTipCap: gasPrice, Amount: value, AccessList: types.AccessList{}, Payload: input,
			})
		}
	}

	if err := tx.SignWithKeys(types.LatestSignerForChainID(chainID), self.privateKey); err != nil {
		return tx.Hash(), gasPrice, err
	}

	err := self.sendRawTx(c, tx, nonce)
	return tx.Hash(), gasPrice, err
}

func (self *Account) TransferNewLegacyTxWithEth(c Client, to *Account, value *big.Int, input []byte) (common.Hash, *big.Int, error) {
	return self.transferNewEthStyleTx(c, to, value, input, EthTxTypeLegacy)
}

// This function is responsible for sending both Gasless Approve Transactions and Gasless Swap Transactions.
func (self *Account) TransferNewGaslessTx(c *client.Client, testToken, gsr *Account) (common.Hash, common.Hash, *big.Int, error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)

	ctx := context.Background()
	suggestedGasPrice, err := c.SuggestGasPrice(ctx)
	if err != nil {
		fmt.Printf("Failed to fetch suggest gas price: %v\n", err.Error())
		return common.Hash{0}, common.Hash{0}, gasPrice, err
	}

	approveTx := types.NewTransaction(
		nonce,
		testToken.address,
		common.Big0,
		100000,
		suggestedGasPrice,
		TestContractInfos[ContractGaslessToken].GenData(gsr.GetAddress(), abi.MaxUint256)) // Approve maximum amount
	signApproveTx, err := types.SignTx(approveTx, types.NewEIP155Signer(chainID), self.privateKey[0])
	if err != nil {
		log.Fatalf("Failed to encode approve tx: %v", err)
	}

	swapTx := types.NewTransaction(
		nonce+1,
		gsr.address,
		common.Big0,
		500000,
		suggestedGasPrice,
		TestContractInfos[ContractGaslessSwapRouter].GenData(testToken.GetAddress(), suggestedGasPrice))
	signSwapTx, err := types.SignTx(swapTx, types.NewEIP155Signer(chainID), self.privateKey[0])
	if err != nil {
		log.Fatalf("Failed to encode swap tx: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	_, err = c.SendRawTransaction(ctx, signApproveTx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return approveTx.Hash(), swapTx.Hash(), suggestedGasPrice, err
	}

	_, err = c.SendRawTransaction(ctx, signSwapTx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return approveTx.Hash(), swapTx.Hash(), suggestedGasPrice, err
	}

	self.nonce += 2

	return approveTx.Hash(), swapTx.Hash(), suggestedGasPrice, nil
}

// This function is responsible for sending only Gasless Approve Transactions.
// This function won't increase nonce.
func (self *Account) TransferNewGaslessApproveTx(c *client.Client, testToken, gsr *Account) (common.Hash, *big.Int, error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)

	ctx := context.Background()
	suggestedGasPrice, err := c.SuggestGasPrice(ctx)
	if err != nil {
		fmt.Printf("Failed to fetch suggest gas price: %v\n", err.Error())
		return common.Hash{0}, gasPrice, err
	}

	approveTx := types.NewTransaction(
		nonce,
		testToken.address,
		common.Big0,
		100000,
		suggestedGasPrice,
		TestContractInfos[ContractGaslessToken].GenData(gsr.GetAddress(), abi.MaxUint256)) // Approve maximum amount
	signApproveTx, err := types.SignTx(approveTx, types.NewEIP155Signer(chainID), self.privateKey[0])
	if err != nil {
		log.Fatalf("Failed to encode approve tx: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	_, err = c.SendRawTransaction(ctx, signApproveTx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return approveTx.Hash(), suggestedGasPrice, err
	}

	return approveTx.Hash(), suggestedGasPrice, nil
}

// BidResult represents the result of a bid operation
type BidResult struct {
	BidHash   common.Hash
	Error     error
	RpcOutput map[string]interface{}
}

func (self *Account) AuctionBid(c *client.Client, auctionEntryPoint, targetContract *Account, targetTxTypeKey string) (common.Hash, common.Hash, *big.Int, error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	// create tmpAccount
	tmpAccount := NewAccount(0)

	/* ---------------- Generate target tx ---------------- */
	nonce := self.GetNonce(c)
	ctx := context.Background()
	suggestedGasPrice, err := c.SuggestGasPrice(ctx)
	if err != nil {
		fmt.Printf("Failed to fetch suggest gas price: %v\n", err.Error())
		return common.Hash{0}, common.Hash{0}, gasPrice, err
	}

	targetTxType := TargetTxTypeList[targetTxTypeKey]
	targetTx := targetTxType.GenerateTx(c, self, tmpAccount, nonce, suggestedGasPrice)
	if targetTx == nil {
		return common.Hash{0}, common.Hash{0}, suggestedGasPrice, errors.New("failed to generate target tx")
	}

	/* ---------------- Send bid -------------------------- */
	err = targetTxType.PreSendBid(c, self, tmpAccount, nonce, suggestedGasPrice)
	if err != nil {
		// If PreSendBid fails, the remaining steps do not need to be performed.
		return common.Hash{0}, common.Hash{0}, suggestedGasPrice, err
	}

	gas := uint64(5000000)

	// Get entrypoint nonce
	appNonce := getEntrypointNonce(c, self.GetAddress())

	// Create contract call data (CounterForAuction.incForAuction())
	contractCallData := TestContractInfos[ContractCounterForTestAuction].GenData(common.Address{}, common.Big0) // 0 means calling incForAuction()

	// Get current block number
	blockNumber, err := c.BlockNumber(ctx)
	if err != nil {
		return common.Hash{0}, common.Hash{0}, suggestedGasPrice, err
	}
	if self.isLastBlocknumSentTx(blockNumber.Uint64()) {
		return common.Hash{0}, common.Hash{0}, suggestedGasPrice, errors.New("this account has already sent a tx for the block")
	}

	// Create bids for blockNumber +1, +2
	numOfMergines := 2
	var bidInputs []*auctionImpl.BidInput
	var bidHashes []common.Hash

	for i := 1; i <= numOfMergines; i++ {
		blockNum := new(big.Int).Add(blockNumber, big.NewInt(int64(i)))

		// Create the bid for this specific block
		bid := &auction.Bid{
			BidData: auction.BidData{
				TargetTxHash: targetTx.Hash(),
				BlockNumber:  blockNum.Uint64(),
				Sender:       self.address,
				To:           targetContract.address,
				Nonce:        appNonce.Uint64(),
				Bid:          big.NewInt(2),
				CallGasLimit: gas,
				Data:         contractCallData,
			},
		}
		bidHashes = append(bidHashes, bid.Hash())

		// searcher sign bid
		searcherSignedBid, err := self.signAuctionBidAsSearcher(bid, auctionEntryPoint)
		if err != nil {
			return common.Hash{0}, common.Hash{0}, suggestedGasPrice, err
		}

		// auctioneer sign bid
		bidInput, err := Auctioneer.signAuctionBidAsAuctioneer(searcherSignedBid, toRlp(targetTx))
		if err != nil {
			return common.Hash{0}, common.Hash{0}, suggestedGasPrice, err
		}
		bidInputs = append(bidInputs, bidInput)
	}

	// Create channel for results
	resultChan := make(chan BidResult, numOfMergines)

	// Launch goroutines for SendAuctionTx only
	for i := 0; i < numOfMergines; i++ {
		go func(index int) {
			ctx := context.Background()
			rpcOutput, err := c.SendAuctionTx(ctx, *bidInputs[index])
			resultChan <- BidResult{
				BidHash:   bidHashes[index],
				Error:     err,
				RpcOutput: rpcOutput,
			}
		}(i)
	}

	// Collect results from all goroutines
	var results []BidResult
	for i := 0; i < numOfMergines; i++ {
		result := <-resultChan
		results = append(results, result)
	}

	/* ---------------- Handle rpc output -------------------------- */
	var successResult *BidResult
	var numNonceTooLowErr int
	var submitErr string
	for _, result := range results {
		if result.Error == nil && result.RpcOutput[auctionImpl.RPC_AUCTION_ERROR_PROP] == nil {
			// If the bid is successful, set the success result and break the loop.
			successResult = &result
			break
		}

		if result.RpcOutput[auctionImpl.RPC_AUCTION_ERROR_PROP] != nil {
			submitErr = result.RpcOutput[auctionImpl.RPC_AUCTION_ERROR_PROP].(string)
		}
		if submitErr != "" {
			// If the bid is failed due to nonce too low, increment the number of nonce too low errors.
			if submitErr == blockchain.ErrNonceTooLow.Error() || submitErr == blockchain.ErrReplaceUnderpriced.Error() {
				numNonceTooLowErr++
			}
		}
	}

	if successResult == nil {
		// If all the bids are failed due to nonce too low, add nonce and return error.
		if numNonceTooLowErr == len(results) {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, submitErr)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		}
		return targetTx.Hash(), common.Hash{0}, suggestedGasPrice, fmt.Errorf("failed to send auction bid: %v", submitErr)
	}

	targetTxType.PostSendBid(c, self, tmpAccount, nonce, suggestedGasPrice, blockNumber)

	return targetTx.Hash(), successResult.BidHash, suggestedGasPrice, nil
}

// AuctionRevertedBid is responsible for sending reverted bid.
// Using an invalid nonce in a bid will cause the bid tx to be reverted.
func (self *Account) AuctionRevertedBid(c *client.Client, auctionEntryPoint, targetContract *Account, targetTxTypeKey string) (common.Hash, common.Hash, *big.Int, error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	// create tmpAccount
	tmpAccount := NewAccount(0)

	/* ---------------- Generate target tx ---------------- */
	nonce := self.GetNonce(c)
	ctx := context.Background()
	suggestedGasPrice, err := c.SuggestGasPrice(ctx)
	if err != nil {
		fmt.Printf("Failed to fetch suggest gas price: %v\n", err.Error())
		return common.Hash{0}, common.Hash{0}, gasPrice, err
	}

	targetTxType := TargetTxTypeList[targetTxTypeKey]
	targetTx := targetTxType.GenerateTx(c, self, tmpAccount, nonce, suggestedGasPrice)
	if targetTx == nil {
		return common.Hash{0}, common.Hash{0}, suggestedGasPrice, errors.New("failed to generate target tx")
	}

	/* ---------------- Send bid -------------------------- */
	err = targetTxType.PreSendBid(c, self, tmpAccount, nonce, suggestedGasPrice)
	if err != nil {
		// If PreSendBid fails, the remaining steps do not need to be performed.
		return common.Hash{0}, common.Hash{0}, suggestedGasPrice, err
	}

	gas := uint64(5000000)

	// Create contract call data (CounterForAuction.incForAuction())
	contractCallData := TestContractInfos[ContractCounterForTestAuction].GenData(common.Address{}, common.Big0) // 0 means calling incForAuction()

	// Get current block number
	blockNumber, err := c.BlockNumber(ctx)
	if err != nil {
		return common.Hash{0}, common.Hash{0}, suggestedGasPrice, err
	}
	if self.isLastBlocknumSentTx(blockNumber.Uint64()) {
		return common.Hash{0}, common.Hash{0}, suggestedGasPrice, errors.New("this account has already sent a tx for the block")
	}

	// Create the bid
	bid := &auction.Bid{
		BidData: auction.BidData{
			TargetTxHash: targetTx.Hash(),
			BlockNumber:  new(big.Int).Add(blockNumber, common.Big1).Uint64(),
			Sender:       self.address,
			To:           targetContract.address,
			Nonce:        math.MaxInt64, // This causes a revert.
			Bid:          big.NewInt(2),
			CallGasLimit: gas,
			Data:         contractCallData,
		},
	}

	// searcher sign bid
	searcherSignedBid, err := self.signAuctionBidAsSearcher(bid, auctionEntryPoint)
	if err != nil {
		return common.Hash{0}, common.Hash{0}, suggestedGasPrice, err
	}

	// auctioneer sign bid
	bidInput, err := Auctioneer.signAuctionBidAsAuctioneer(searcherSignedBid, toRlp(targetTx))
	if err != nil {
		return common.Hash{0}, common.Hash{0}, suggestedGasPrice, err
	}

	// send bid
	var submitErr string
	rpcOutput, err := c.SendAuctionTx(ctx, *bidInput)
	if err != nil {
		return targetTx.Hash(), bid.Hash(), suggestedGasPrice, err
	}

	/* ---------------- Handle rpc output -------------------------- */
	if rpcOutput[auctionImpl.RPC_AUCTION_ERROR_PROP] != nil {
		submitErr = rpcOutput[auctionImpl.RPC_AUCTION_ERROR_PROP].(string)
	}
	if submitErr != "" {
		if submitErr == blockchain.ErrNonceTooLow.Error() || submitErr == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, submitErr)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		}
		return targetTx.Hash(), bid.Hash(), suggestedGasPrice, errors.New(submitErr)
	}

	targetTxType.PostSendBid(c, self, tmpAccount, nonce, suggestedGasPrice, blockNumber)

	return targetTx.Hash(), bid.Hash(), suggestedGasPrice, nil
}

func (self *Account) TransferNewEthAccessListTxWithEth(c *client.EthClient, to *Account, value *big.Int, input []byte) (common.Hash, *big.Int, error) {
	return self.transferNewEthStyleTx(c, to, value, input, EthTxTypeAccessList)
}

func (self *Account) TransferNewEthDynamicFeeTxWithEth(c Client, to *Account, value *big.Int, input []byte) (common.Hash, *big.Int, error) {
	return self.transferNewEthStyleTx(c, to, value, input, EthTxTypeDynamicFee)
}

func (self *Account) TransferUnsignedTx(c *client.Client, to *Account, value *big.Int) (common.Hash, error) {
	ctx := context.Background()

	fromAddr := self.GetAddress()
	toAddr := to.GetAddress()
	gasLimit := hexutil.Uint64(21000)

	// Initialize empty data and payload for value transfer transaction
	emptyData := hexutil.Bytes{}
	emptyPayload := hexutil.Bytes{}

	var err error
	hash, err := c.SendUnsignedTransaction(ctx, api.SendTxArgs{
		From:      fromAddr,
		Recipient: &toAddr,
		GasLimit:  &gasLimit,
		Price:     (*hexutil.Big)(gasPrice),
		Amount:    (*hexutil.Big)(value),
		Data:      &emptyData,
		Payload:   &emptyPayload,
	})
	if err != nil {
		log.Printf("Account(%v) : Failed to sendTransaction: %v\n", self.address[:5], err)
		return common.Hash{}, err
	}
	//log.Printf("Account(%v) : Success to sendTransaction: %v\n", self.address[:5], hash.String())
	return hash, nil
}

// SmartContractDeployWithGuaranteeRetry deploys only one smart contract among the slaves.
// It the contract is already deployed by other slave, it just calculates the address of the contract.
func (self *Account) SmartContractDeployWithGuaranteeRetry(gCli *client.Client, byteCode []byte, contractName string, shouldFixNonceZero bool) *Account {
	log.Println(contractName, "deployer", self.address.String())

	nonce := self.GetNonce(gCli)
	if shouldFixNonceZero {
		nonce = 0
	}

	config := RetryConfig{
		SendRetryInterval: 5 * time.Second,
		WaitMinedTimeout:  60 * time.Second,
		ShouldSkip: func(err error) bool {
			// Treat "known transaction" and ErrNonceTooLow (when shouldFixNonceZero) as success
			return strings.HasPrefix(err.Error(), "known transaction") ||
				(shouldFixNonceZero && err.Error() == blockchain.ErrNonceTooLow.Error())
		},
	}

	self.RunWithRetry(gCli, config, func() (*types.Transaction, error) {
		_, tx, _, err := self.TransferNewSmartContractDeployTx(gCli, nil, common.Big0, byteCode, shouldFixNonceZero)
		return tx, err
	})

	addr := crypto.CreateAddress(self.GetAddress(), nonce)
	log.Printf("%s has been deployed to : %s\n", contractName, addr.String())
	return NewKaiaAccountWithAddr(1, addr)
}

func (a *Account) SmartContractExecutionWithGuaranteeRetry(gCli *client.Client, to *Account, value *big.Int, data []byte) {
	a.RunWithRetry(gCli, DefaultRetryConfig(), func() (*types.Transaction, error) {
		tx, _, err := a.TransferNewSmartContractExecutionTx(gCli, to, value, data)
		return tx, err
	})
}

func (a *Account) TryRunTxSendFunctionWithGuaranteeRetry(gCli *client.Client, allowedErrors []error, txSendFunc func(gCli *client.Client, sender *Account) (*types.Transaction, error)) {
	config := DefaultRetryConfig()
	config.ShouldSkip = func(err error) bool {
		for _, allowError := range allowedErrors {
			if err.Error() == allowError.Error() {
				log.Printf("Skipping the transaction: err=%s", err.Error())
				return true
			}
		}
		return false
	}
	a.RunWithRetry(gCli, config, func() (*types.Transaction, error) {
		return txSendFunc(gCli, a)
	})
}

func (a *Account) CheckBalance(expectedBalance *big.Int, cli *client.Client) error {
	balance, _ := a.GetBalance(cli)
	if balance.Cmp(expectedBalance) != 0 {
		fmt.Println(a.address.String() + " expected : " + expectedBalance.Text(10) + " actual : " + balance.Text(10))
		return errors.New("expected : " + expectedBalance.Text(10) + " actual : " + balance.Text(10))
	}

	return nil
}

func (a *Account) isLastBlocknumSentTx(blockNumber uint64) bool {
	return a.lastBlocknumSentTx == blockNumber
}

func (a *Account) updateLastBlocknumSentTx(blockNumber uint64) {
	a.lastBlocknumSentTx = blockNumber
}

func (a *Account) signAuctionBidAsSearcher(bid *auction.Bid, auctionEntryPoint *Account) (*auction.Bid, error) {
	bidHash := bid.GetHashTypedData(chainID, auctionEntryPoint.GetAddress())
	searcherSig, err := crypto.Sign(bidHash, a.privateKey[0])
	if err != nil {
		return nil, err
	}
	searcherSig[crypto.RecoveryIDOffset] += 27
	bid.SearcherSig = searcherSig
	return bid, nil
}

func (a *Account) signAuctionBidAsAuctioneer(auctionBid *auction.Bid, targetTxRaw []byte) (*auctionImpl.BidInput, error) {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(auctionBid.SearcherSig), auctionBid.SearcherSig)
	hash := crypto.Keccak256([]byte(msg))
	auctioneerSig, err := crypto.Sign(hash, a.privateKey[0])
	if err != nil {
		return nil, err
	}

	auctioneerSig[crypto.RecoveryIDOffset] += 27
	bidInput := &auctionImpl.BidInput{
		TargetTxRaw:   targetTxRaw,
		TargetTxHash:  auctionBid.TargetTxHash,
		BlockNumber:   auctionBid.BlockNumber,
		Sender:        auctionBid.Sender,
		To:            auctionBid.To,
		Nonce:         auctionBid.Nonce,
		Bid:           hexutil.Big(*auctionBid.Bid),
		CallGasLimit:  auctionBid.CallGasLimit,
		Data:          auctionBid.Data,
		SearcherSig:   auctionBid.SearcherSig,
		AuctioneerSig: auctioneerSig,
	}
	return bidInput, nil
}

func toRlp(tx *types.Transaction) []byte {
	rlp, _ := rlp.EncodeToBytes(tx)
	return rlp
}

func ConcurrentTransactionSend(accs []*Account, maxConcurrency int, transactionSend func(int, *Account)) {
	if maxConcurrency <= 0 {
		maxConcurrency = runtime.NumCPU() * 10 // default value
	}

	total := len(accs)
	if total == 0 {
		return
	}

	var completed int64
	ch := make(chan int, maxConcurrency)
	wg := sync.WaitGroup{}

	// Progress reporting interval (every 10% or minimum 1)
	reportInterval := total / 10
	if reportInterval == 0 {
		reportInterval = 1
	}

	for idx, acc := range accs {
		ch <- 1
		wg.Add(1)
		go func() {
			transactionSend(idx, acc)
			<-ch

			current := atomic.AddInt64(&completed, 1)
			if current%int64(reportInterval) == 0 || current == int64(total) {
				log.Printf("Progress: %d/%d (%.1f%%)", current, total, float64(current)/float64(total)*100)
			}

			wg.Done()
		}()
	}
	wg.Wait()
}
