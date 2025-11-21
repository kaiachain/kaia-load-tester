package account

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/orderbook"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
)

const Letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

var (
	gasPrice  *big.Int
	chainID   *big.Int
	baseFee   *big.Int
	batchSize int
)

type Account struct {
	id         int
	privateKey *ecdsa.PrivateKey
	address    common.Address
	sessionKey []*ecdsa.PrivateKey
	sessionCtx []*types.SessionContext // not nil once session is registered
	nonce      uint64
	timenonce  uint64
	balance    *big.Int
	mutex      sync.Mutex
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

func SetBatchSize(bs int) {
	batchSize = bs
}

func (acc *Account) Lock() {
	acc.mutex.Lock()
}

func (acc *Account) UnLock() {
	acc.mutex.Unlock()
}

func (acc *Account) updateTimeNonce() {
	currentTimeMillis := uint64(time.Now().UnixMilli())

	if acc.timenonce == 0 || acc.timenonce < currentTimeMillis {
		acc.timenonce = currentTimeMillis + 120000 // 2 minutes = 120000 milliseconds
	} else {
		acc.timenonce++
	}
}

func GetAccountFromKey(id int, key string) *Account {
	// Normalize the private key format
	key = strings.TrimPrefix(key, "0x")

	// Ensure the key is exactly 64 hex characters (32 bytes)
	if len(key) != 64 {
		log.Fatalf("Private key must be exactly 64 hex characters (32 bytes), got %d characters: %s\nExample of valid key: 2ef07640fd8d3f568c23185799ee92e0154bf08ccfe5c509466d1d40baca3430", len(key), key)
	}

	acc, err := crypto.HexToECDSA(key)
	if err != nil {
		log.Fatalf("Key(%v): Failed to HexToECDSA %v", key, err)
	}

	tAcc := Account{
		0,
		acc,
		crypto.PubkeyToAddress(acc.PublicKey),
		[]*ecdsa.PrivateKey{},
		[]*types.SessionContext{},
		0,
		uint64(time.Now().UnixMilli()),
		big.NewInt(0),
		sync.Mutex{},
		// make(TransactionMap),
	}

	return &tAcc
}

func (acc *Account) ImportUnLockAccount(endpoint string) {
}

func NewAccount(id int) *Account {
	acc, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("crypto.GenerateKey() : Failed to generateKey %v", err)
	}

	tAcc := Account{
		0,
		acc,
		crypto.PubkeyToAddress(acc.PublicKey),
		[]*ecdsa.PrivateKey{},
		[]*types.SessionContext{},
		0,
		uint64(time.Now().UnixMilli()),
		big.NewInt(0),
		sync.Mutex{},
		// make(TransactionMap),
	}

	return &tAcc
}

func (acc *Account) GetAddress() common.Address {
	return acc.address
}

func (acc *Account) GetNonce(c *ethclient.Client) uint64 {
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

	// fmt.Printf("account= %v  nonce = %v\n", acc.GetAddress().String(), nonce)
	return acc.nonce
}

func (acc *Account) GetNonceFromBlock(c *ethclient.Client) uint64 {
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

func (acc *Account) GetSessionCtx() []*types.SessionContext {
	return acc.sessionCtx
}

func (acc *Account) UpdateNonce() {
	acc.nonce++
}

// NewSessionCreateCtx creates a new session
func (acc *Account) NewSessionCreateCtx(expiresAt uint64, nonce uint64) (*types.SessionContext, *ecdsa.PrivateKey, error) {
	sessionKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, nil, err
	}
	sessionAddr := crypto.PubkeyToAddress(sessionKey.PublicKey)
	session := types.Session{
		PublicKey: sessionAddr,
		ExpiresAt: expiresAt,
		Nonce:     nonce,
		Metadata:  nil,
	}
	typedData := types.ToTypedData(&session)
	_, sigHash, _ := types.SignEip712(typedData)
	sig, err := crypto.Sign(sigHash, acc.privateKey)
	if err != nil {
		return nil, nil, err
	}

	return &types.SessionContext{
		Command:     types.SessionCreate,
		Session:     session,
		L1Owner:     acc.GetAddress(),
		L1Signature: sig,
	}, sessionKey, nil
}

// NewSessionDeleteCtx creates a new session
func (acc *Account) NewSessionDeleteCtx(i int, nonce uint64) (*types.SessionContext, error) {
	target := acc.sessionCtx[i]
	sessionAddr := target.Session.PublicKey
	session := types.Session{
		PublicKey: sessionAddr,
		ExpiresAt: target.Session.ExpiresAt,
		Nonce:     nonce,
		Metadata:  nil,
	}
	typedData := types.ToTypedData(&session)
	_, sigHash, _ := types.SignEip712(typedData)
	sig, err := crypto.Sign(sigHash, acc.privateKey)
	if err != nil {
		return nil, err
	}

	return &types.SessionContext{
		Command:     types.SessionDelete,
		Session:     session,
		L1Owner:     acc.GetAddress(),
		L1Signature: sig,
	}, nil
}

func (acc *Account) NewValueTransferCtx(to *Account, value *big.Int) *types.ValueTransferContext {
	vtCtx := types.ValueTransferContext{
		L1Owner: acc.GetAddress(),
		To:      to.GetAddress(),
		Value:   value,
	}

	return &vtCtx
}

func (acc *Account) NewTokenTransferCtx(to *Account, value *big.Int, token string) *types.TokenTransferContext {
	ctx := types.TokenTransferContext{
		L1Owner: acc.GetAddress(),
		To:      to.GetAddress(),
		Value:   value,
		Token:   token,
	}

	return &ctx
}

func (acc *Account) NewOrderCtx(baseToken string, quoteToken string, side orderbook.Side, price *big.Int, quantity *big.Int, orderType orderbook.OrderType) *types.OrderContext {
	ctx := types.OrderContext{
		L1Owner:    acc.GetAddress(),
		BaseToken:  baseToken,
		QuoteToken: quoteToken,
		Side:       uint8(side),
		Price:      price,
		Quantity:   quantity,
		OrderType:  uint8(orderType),
		OrderMode:  0,
		TPSL:       nil,
	}

	return &ctx
}

func (acc *Account) NewOrderCtxWithTpsl(baseToken string, quoteToken string, side orderbook.Side, price *big.Int, quantity *big.Int, orderType orderbook.OrderType, tpLimit, slTrigger, slLimit *big.Int) *types.OrderContext {
	tpsl := types.TPSLContext{
		TPLimit:   tpLimit,
		SLTrigger: slTrigger,
		SLLimit:   slLimit,
	}

	ctx := acc.NewOrderCtx(baseToken, quoteToken, side, price, quantity, orderType)
	ctx.TPSL = &tpsl

	return ctx
}

func (acc *Account) NewStopOrderCtx(baseToken string, quoteToken string, side orderbook.Side, stopPrice, price *big.Int, quantity *big.Int, orderType orderbook.OrderType) *types.StopOrderContext {
	ctx := types.StopOrderContext{
		L1Owner:    acc.GetAddress(),
		BaseToken:  baseToken,
		QuoteToken: quoteToken,
		StopPrice:  stopPrice,
		Price:      price,
		Quantity:   quantity,
		Side:       uint8(side),
		OrderType:  uint8(orderType),
		OrderMode:  0,
	}

	return &ctx
}

func (acc *Account) NewCancelAllCtx() *types.CancelAllContext {
	ctx := types.CancelAllContext{
		L1Owner: acc.GetAddress(),
	}

	return &ctx
}

func (acc *Account) GetReceipt(c *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	ctx := context.Background()
	return c.TransactionReceipt(ctx, txHash)
}

func (acc *Account) GetBalance(c *ethclient.Client) (*big.Int, error) {
	ctx := context.Background()
	balance, err := c.BalanceAt(ctx, acc.GetAddress(), nil)
	if err != nil {
		return nil, err
	}
	return balance, err
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func (acc *Account) TransferSignedTxReturnTx(withLock bool, c *ethclient.Client, to *Account, value *big.Int) (*types.Transaction, *big.Int, error) {
	if withLock {
		acc.mutex.Lock()
		defer acc.mutex.Unlock()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	nonce, err := c.NonceAt(ctx, acc.GetAddress(), nil)
	if err != nil {
		return nil, big.NewInt(0), err
	}

	acc.nonce = nonce
	tx := types.NewTransaction(
		acc.nonce,
		to.GetAddress(),
		value,
		21000,
		gasPrice,
		nil)
	gasPrice := tx.GasPrice()

	signer := types.LatestSignerForChainID(chainID)
	tx, err = types.SignTx(tx, signer, acc.privateKey)
	if err != nil {
		return nil, gasPrice, err
	}
	err = c.SendTransaction(ctx, tx)
	if err != nil {
		return tx, gasPrice, err
	}

	acc.nonce++

	// fmt.Printf("%v transferSignedTx %v klay to %v klay.\n", self.GetAddress().Hex(), to.GetAddress().Hex(), value)

	return tx, gasPrice, nil
}

func (acc *Account) TransferSignedTxWithGuaranteeRetry(c *ethclient.Client, to *Account, value *big.Int) *types.Transaction {
	var (
		err    error
		lastTx *types.Transaction
	)

	for {
		lastTx, _, err = acc.TransferSignedTxReturnTx(true, c, to, value)
		// TODO-kaia-load-tester: return error if the error isn't able to handle
		if err == nil {
			break // Succeed, let's break the loop
		}
		log.Printf("Failed to execute: err=%s", err.Error())
		if strings.Contains(err.Error(), "time nonce already exists") {
			acc.timenonce = uint64(time.Now().UnixMilli()) - uint64(rand.Intn(1000))
		}
		log.Printf("Retrying...")
		time.Sleep(1 * time.Second) // Mostly, the err is `txpool is full`, retry after a while.
		// numChargedAcc, lastFailedNum = estimateRemainingTime(accGrp, numChargedAcc, lastFailedNum)
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFn()

	receipt, err := bind.WaitMined(ctx, c, lastTx)
	cancelFn()
	if err != nil || (receipt != nil && receipt.Status == 0) {
		// shouldn't happen. must check if contract is correct.
		log.Fatalf("tx mined but failed, err=%s, txHash=%s", err, lastTx.Hash().String())
	}
	return lastTx
}

func (acc *Account) TransferTokenSignedTxWithGuaranteeRetry(c *ethclient.Client, to *Account, value *big.Int, token string) *types.Transaction {
	var (
		err error
		tx  *types.Transaction
	)

	for {
		time.Sleep(1 * time.Second)
		tx, err = acc.GenTokenTransferTx(to, value, token)
		if err != nil {
			log.Printf("Failed to generate token transfer: err=%v, from=%v, to=%v, timenonce=%v, value=%v, token=%v",
				err.Error(), acc.GetAddress().String(), to.GetAddress().String(), acc.timenonce, value.String(), token)
			continue
		}
		_, err = acc.SendTx(c, tx)
		if err != nil {
			log.Printf("Failed to send token transfer tx: err=%v, from=%v, to=%v, timenonce=%v, value=%v, token=%v",
				err.Error(), acc.GetAddress().String(), to.GetAddress().String(), acc.timenonce, value.String(), token)
			if strings.Contains(err.Error(), "insufficient") {
				os.Exit(1)
			}
			continue
		}
		receipt, err := c.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Printf("Failed to fetch receipt of tx %s: err=%v, from=%v, to=%v, value=%v, token=%v",
				tx.Hash().String(), err.Error(), acc.GetAddress().String(), to.GetAddress().String(), value.String(), token)
			continue
		}

		if receipt.Status != types.ReceiptStatusSuccessful {
			log.Printf("Token transfer tx %s is failed %d, from=%v, to=%v, value=%v, token=%v",
				tx.Hash().String(), receipt.Status, acc.GetAddress().String(), to.GetAddress().String(), value.String(), token)
			continue
		}

		break
	}

	return tx
}

func (acc *Account) RegisterNewSession(c *ethclient.Client) error {
	tx, err := acc.GenSessionCreateTx()
	if err != nil {
		return err
	}

	_, err = acc.SendTx(c, tx)
	if err != nil {
		return err
	}

	return nil
}

func (acc *Account) GenLegacyTx(to *Account, value *big.Int, input []byte) (*types.Transaction, error) {
	acc.mutex.Lock()
	defer acc.mutex.Unlock()

	tx := types.NewTransaction(
		acc.nonce,
		to.GetAddress(),
		value,
		21000,
		big.NewInt(25e9),
		input,
	)

	signer := types.LatestSignerForChainID(chainID)
	tx, err := types.SignTx(tx, signer, acc.privateKey)
	if err != nil {
		return nil, err
	}
	acc.nonce++

	return tx, nil
}

func (acc *Account) GenSessionCreateTx() (*types.Transaction, error) {
	acc.mutex.Lock()
	defer acc.mutex.Unlock()
	acc.updateTimeNonce()

	sessionCtx, sessionKey, err := acc.NewSessionCreateCtx(uint64(1000000), acc.timenonce)
	if err != nil {
		return nil, err
	}

	acc.sessionCtx = append(acc.sessionCtx, sessionCtx)
	acc.sessionKey = append(acc.sessionKey, sessionKey)

	input, err := types.WrapTxAsInput(sessionCtx)
	if err != nil {
		return nil, err
	}

	tx := types.NewTransaction(
		sessionCtx.Session.Nonce,
		types.DexAddress,
		common.Big0,
		0,
		common.Big0,
		input,
	)

	signer := types.LatestSignerForChainID(chainID)
	tx, err = types.SignTx(tx, signer, sessionKey)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (acc *Account) GenSessionDeleteTx(i int) (*types.Transaction, error) {
	acc.mutex.Lock()
	defer acc.mutex.Unlock()
	acc.updateTimeNonce()

	sessionCtx, err := acc.NewSessionDeleteCtx(i, acc.timenonce)
	if err != nil {
		return nil, err
	}

	signer := types.LatestSignerForChainID(chainID)
	input, err := types.WrapTxAsInput(sessionCtx)
	if err != nil {
		return nil, err
	}

	tx := types.NewTransaction(
		sessionCtx.Session.Nonce,
		types.DexAddress,
		common.Big0,
		0,
		common.Big0,
		input,
	)

	tx, err = types.SignTx(tx, signer, acc.sessionKey[i])
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (acc *Account) GenTransferTx(to *Account, value *big.Int) (*types.Transaction, error) {
	acc.mutex.Lock()
	defer acc.mutex.Unlock()
	acc.updateTimeNonce()

	ctx := acc.NewValueTransferCtx(to, value)

	signer := types.LatestSignerForChainID(chainID)
	input, err := types.WrapTxAsInput(ctx)
	if err != nil {
		return nil, err
	}

	tx := types.NewTransaction(
		acc.timenonce,
		types.DexAddress,
		common.Big0,
		0,
		common.Big0,
		input,
	)

	tx, err = types.SignTx(tx, signer, acc.privateKey)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (acc *Account) GenTokenTransferTx(to *Account, value *big.Int, token string) (*types.Transaction, error) {
	acc.mutex.Lock()
	defer acc.mutex.Unlock()
	acc.updateTimeNonce()

	ctx := acc.NewTokenTransferCtx(to, value, token)

	signer := types.LatestSignerForChainID(chainID)
	input, err := types.WrapTxAsInput(ctx)
	if err != nil {
		return nil, err
	}

	tx := types.NewTransaction(
		acc.timenonce,
		types.DexAddress,
		common.Big0,
		0,
		common.Big0,
		input,
	)

	tx, err = types.SignTx(tx, signer, acc.privateKey)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (acc *Account) GenNewOrderTx(baseToken string, quoteToken string, side orderbook.Side, price *big.Int, quantity *big.Int, orderType orderbook.OrderType) (*types.Transaction, error) {
	acc.mutex.Lock()
	defer acc.mutex.Unlock()
	acc.updateTimeNonce()

	ctx := acc.NewOrderCtx(baseToken, quoteToken, side, price, quantity, orderType)

	signer := types.LatestSignerForChainID(chainID)
	input, err := types.WrapTxAsInput(ctx)
	if err != nil {
		return nil, err
	}

	tx := types.NewTransaction(
		acc.timenonce,
		types.DexAddress,
		common.Big0,
		0,
		common.Big0,
		input,
	)

	tx, err = types.SignTx(tx, signer, acc.privateKey)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (acc *Account) GenNewOrderTxWithTpsl(baseToken string, quoteToken string, side orderbook.Side, price *big.Int, quantity *big.Int, orderType orderbook.OrderType, tpLimit, slTrigger, slLimit *big.Int) (*types.Transaction, error) {
	acc.mutex.Lock()
	defer acc.mutex.Unlock()
	acc.updateTimeNonce()

	ctx := acc.NewOrderCtxWithTpsl(baseToken, quoteToken, side, price, quantity, orderType, tpLimit, slTrigger, slLimit)

	signer := types.LatestSignerForChainID(chainID)
	input, err := types.WrapTxAsInput(ctx)
	if err != nil {
		return nil, err
	}

	tx := types.NewTransaction(
		acc.timenonce,
		types.DexAddress,
		common.Big0,
		0,
		common.Big0,
		input,
	)

	tx, err = types.SignTx(tx, signer, acc.privateKey)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (acc *Account) GenNewStopOrderTx(baseToken string, quoteToken string, side orderbook.Side, stopPrice, price *big.Int, quantity *big.Int, orderType orderbook.OrderType) (*types.Transaction, error) {
	acc.mutex.Lock()
	defer acc.mutex.Unlock()
	acc.updateTimeNonce()

	ctx := acc.NewStopOrderCtx(baseToken, quoteToken, side, stopPrice, price, quantity, orderType)

	signer := types.LatestSignerForChainID(chainID)
	input, err := types.WrapTxAsInput(ctx)
	if err != nil {
		return nil, err
	}

	tx := types.NewTransaction(
		acc.timenonce,
		types.DexAddress,
		common.Big0,
		0,
		common.Big0,
		input,
	)

	tx, err = types.SignTx(tx, signer, acc.privateKey)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (acc *Account) GenCancelAllTx() (*types.Transaction, error) {
	acc.mutex.Lock()
	defer acc.mutex.Unlock()
	acc.updateTimeNonce()

	ctx := acc.NewCancelAllCtx()

	signer := types.LatestSignerForChainID(chainID)
	input, err := types.WrapTxAsInput(ctx)
	if err != nil {
		return nil, err
	}

	tx := types.NewTransaction(
		acc.timenonce,
		types.DexAddress,
		common.Big0,
		0,
		common.Big0,
		input,
	)

	tx, err = types.SignTx(tx, signer, acc.privateKey)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (acc *Account) SendTx(c *ethclient.Client, tx *types.Transaction) (common.Hash, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	err := c.SendTransaction(ctx, tx)
	if err != nil {
		if strings.Contains(err.Error(), "time nonce too low") {
			acc.mutex.Lock()
			acc.updateTimeNonce()
			acc.mutex.Unlock()
		}
		return common.Hash{}, err
	}

	return tx.Hash(), nil
}

func (acc *Account) SendTxBatch(c *ethclient.Client, txs []*types.Transaction) ([]*hexutil.Bytes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	reqs := make([]rpc.BatchElem, len(txs))
	ret := make([]*hexutil.Bytes, len(txs))
	for i := range txs {
		rlpTx, err := rlp.EncodeToBytes(txs[i])
		if err != nil {
			return nil, err
		}
		reqs[i] = rpc.BatchElem{
			Method: "eth_sendRawTransaction",
			Args:   []interface{}{hexutil.Encode(rlpTx)},
			Result: new(hexutil.Bytes),
		}
	}

	err := c.Client().BatchCallContext(ctx, reqs)
	if err != nil {
		return nil, err
	}

	for i := range txs {
		ret[i] = reqs[i].Result.(*hexutil.Bytes)
	}

	return ret, nil
}

/*
func (self *Account) TransferNewLegacyTxWithEthBatch(c *ethclient.Client, endpoint string, to *Account, value *big.Int, input string) ([]common.Hash, *big.Int, error) {
	self.mutex.lock()
	defer self.mutex.unlock()

	var toAddress common.Address
	if to == nil {
		toAddress = common.Address{}
	} else {
		toAddress = to.GetAddress()
	}

	txs := make([]*types.Transaction, batchSize)
	for i := range batchSize {
		nonce := self.GetNonce(c)
		txs[i] = types.NewTransaction(
			nonce,
			toAddress,
			value,
			100000,
			gasPrice,
			common.FromHex(input),
		)
		self.nonce++
	}

	signer := types.LatestSignerForChainID(chainID)

	var wg sync.WaitGroup
	errChan := make(chan error, len(txs))

	for i := range txs {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			if err := txs[idx].SignWithKeys(signer, self.privateKey); err != nil {
				errChan <- fmt.Errorf("failed to sign tx %d: %v", idx, err)
			}
		}(i)
	}

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		log.Fatalf("Failed to sign transactions: %v", err)
	}

	ctx := context.Background()
	err := c.SendTransactionBatch(ctx, txs)
	if err != nil {
		if err.Error() == core.ErrNonceTooLow.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), self.nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), self.nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), self.nonce, err)
		}
		return nil, gasPrice, err
	}

	hashes := make([]common.Hash, len(txs))
	for i := range txs {
		hashes[i] = txs[i].Hash()
	}

	return hashes, gasPrice, nil
}

func (self *Account) CheckReceiptsBatch(c *ethclient.Client, txHashes []common.Hash) ([]*types.Receipt, error) {
	ctx := context.Background()
	defer ctx.Done()

	receipts, err := c.TransactionReceiptBatch(ctx, txHashes)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction receipts: %v", err)
	}

	for i, receipt := range receipts {
		if receipt == nil {
			return nil, fmt.Errorf("receipt not found for transaction %s", txHashes[i].Hex())
		}
		if receipt.Status != types.ReceiptStatusSuccessful {
			return nil, fmt.Errorf("transaction %s failed with status %d", txHashes[i].Hex(), receipt.Status)
		}
	}

	return receipts, nil
}
*/

func (acc *Account) CheckBalance(expectedBalance *big.Int, cli *ethclient.Client) error {
	balance, _ := acc.GetBalance(cli)
	if balance.Cmp(expectedBalance) != 0 {
		fmt.Println(acc.address.String() + " expected : " + expectedBalance.Text(10) + " actual : " + balance.Text(10))
		return errors.New("expected : " + expectedBalance.Text(10) + " actual : " + balance.Text(10))
	}

	return nil
}

func ConcurrentTransactionSend(accs []*Account, transactionSend func(*Account)) {
	ch := make(chan int, runtime.NumCPU()*10)
	wg := sync.WaitGroup{}
	for _, acc := range accs {
		ch <- 1
		wg.Add(1)
		go func() {
			transactionSend(acc)
			<-ch
			wg.Done()
		}()
	}
	wg.Wait()
}

const (
	numChunks = 4
)

// calcTotalTxCount calculates the total number of transactions needed
// for hierarchical distribution of n accounts with branching factor `numChunks`
func calcTotalTxCount(n int) int {
	if n <= numChunks {
		return n // Base case: direct distribution
	}

	chunkSize := (n + numChunks - 1) / numChunks
	totalTxs := numChunks

	// Add transactions for each chunk recursively
	for i := 0; i < numChunks; i++ {
		start := i * chunkSize
		end := min(start+chunkSize, n)
		if start < n {
			chunkN := end - start
			totalTxs += calcTotalTxCount(chunkN)
		}
	}

	return totalTxs
}

// HierarchicalDistribute calls sendTx in a hierarchical manner where each child can be parallelized.
// rich -> richChild -> richChildChild -> ... -> accs[0..n]
func HierarchicalDistribute(accs []*Account, from *Account, value *big.Int, gasFee *big.Int, sendTx func(from, to *Account, value *big.Int)) {
	if len(accs) <= numChunks {
		// Base case: distribute directly
		for _, acc := range accs {
			sendTx(from, acc, value)
		}
		return
	}

	// Divide-and-conquer case
	chunkSize := (len(accs) + numChunks - 1) / numChunks
	var wg sync.WaitGroup

	for i := 0; i < numChunks; i++ {
		start := i * chunkSize
		end := min(start+chunkSize, len(accs))
		chunkAccs := accs[start:end]

		// Calculate total amount needed for child richAcc
		chunkAmount := new(big.Int).Mul(value, big.NewInt(int64(len(chunkAccs))))

		// Calculate total transactions needed for this chunk's hierarchy
		// For n accounts with branching factor 4: ceil((n-1)/3) internal nodes + n leaf transactions
		numTxs := calcTotalTxCount(len(chunkAccs))

		// Add gas fees for all transactions in this chunk's hierarchy
		totalGasFees := new(big.Int).Mul(gasFee, big.NewInt(int64(numTxs)))
		chunkAmount = new(big.Int).Add(chunkAmount, totalGasFees)

		// Create intermediate account and fund it
		richChild := NewAccount(0)
		sendTx(from, richChild, chunkAmount)

		// Process chunk in parallel
		wg.Add(1)
		go func(child *Account, accounts []*Account) {
			defer wg.Done()
			HierarchicalDistribute(chunkAccs, richChild, value, gasFee, sendTx)
		}(richChild, chunkAccs)
	}

	// Wait for all chunks to complete
	wg.Wait()
}

func (acc *Account) PrivateKey() *ecdsa.PrivateKey {
	return acc.privateKey
}
