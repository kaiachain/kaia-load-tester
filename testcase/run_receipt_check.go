package testcase

import (
	"context"
	"errors"
	"log"
	"math/big"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kaiachain/kaia-load-tester/klayslave/account"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/client"
	"github.com/kaiachain/kaia/common"
	"github.com/myzhan/boomer"
)

// Receipt check related variables
var (
	hashPoolSize      = 100 * 5 * 60 // for init 5min, if input send tps is 100Txs/Sec
	defaultInitSendTx = 1000 * 10    // for init 10sec, if input send TPS is 1000Txs/Sec

	tail     = 0
	isFull   = false
	hashPool = make([]common.Hash, hashPoolSize)
	rwMutex  sync.RWMutex

	ratioReadPerSend = 9 // read:send = ratioReadPerSend:1

	cnt      uint32
	initFlag = false
)

// addHash adds a hash to the hash pool
func addHash(hash common.Hash) {
	rwMutex.Lock()
	hashPool[tail] = hash

	tail = (tail + 1) % hashPoolSize
	if tail == 0 {
		isFull = true
	}

	rwMutex.Unlock()
}

// getHash gets a random hash from the hash pool
func getHash() common.Hash {
	rwMutex.RLock()
	defer rwMutex.RUnlock()
	if isFull {
		return hashPool[rand.Int()%hashPoolSize]
	}
	return hashPool[rand.Int()%tail]
}

// doubleLock locks two accounts in a consistent order to prevent deadlock
func doubleLock(to *account.Account, from *account.Account) {
	if from.GetAddress().String() == to.GetAddress().String() {
		from.Lock()
	} else if from.GetAddress().String() > to.GetAddress().String() {
		from.Lock()
		to.Lock()
	} else {
		to.Lock()
		from.Lock()
	}
}

// doubleUnlock unlocks two accounts in reverse order
func doubleUnlock(to *account.Account, from *account.Account) {
	if from.GetAddress().String() == to.GetAddress().String() {
		from.UnLock()
	} else if from.GetAddress().String() > to.GetAddress().String() {
		from.UnLock()
		to.UnLock()
	} else {
		to.UnLock()
		from.UnLock()
	}
}

// runReceiptCheckSendTx creates a closure for receipt check send transaction
func runReceiptCheckSendTx(config *TCConfig) func() {
	return func() {
		cli := config.CliPool.Alloc().(*client.KaiaClient)
		defer config.CliPool.Free(cli)

		from := config.AccGrp.GetAccountRandomly()
		to := config.AccGrp.GetAccountRandomly()
		value := big.NewInt(int64(rand.Int() % 3))

		start := boomer.Now()
		hash, _, err := from.TransferSignedTx(cli, to, value, nil)
		addHash(hash)
		elapsed := boomer.Now() - start

		if err == nil {
			boomer.Events.Publish("request_success", "receiptCheckTx", "send tx"+" to "+config.EndPoint, elapsed, int64(10))
		} else {
			boomer.Events.Publish("request_failure", "receiptCheckTx", "send tx"+" to "+config.EndPoint, elapsed, err.Error())
		}
	}
}

// runReceiptCheckReadTx creates a closure for receipt check read transaction
func runReceiptCheckReadTx(config *TCConfig) func() {
	return func() {
		cli := config.CliPool.Alloc().(*client.KaiaClient)
		defer config.CliPool.Free(cli)

		ctx := context.Background()
		hash := getHash()

		start := boomer.Now()

		receipt, err := cli.TransactionReceipt(ctx, hash)
		if err == nil {
			if rand.Int()%(1000*60) == 0 {
				log.Printf("pid(%v) : hash(%v) receipt checked\n", os.Getpid(), hash.String())
				log.Printf("%v", receipt)
			}
		} else {
			log.Printf("pid(%v) : hash(%v) receipt check err : %v\n", os.Getpid(), hash.String(), err)
		}

		elapsed := boomer.Now() - start

		if err == nil {
			boomer.Events.Publish("request_success", "receiptCheckTx", "read tx"+" to "+config.EndPoint, elapsed, int64(10))
		} else {
			boomer.Events.Publish("request_failure", "receiptCheckTx", "read tx"+" to "+config.EndPoint, elapsed, err.Error())
		}
	}
}

// RunReceiptCheckTC creates a closure for receipt check test case
func RunReceiptCheckTC(config *TCConfig) func() {
	return func() {
		nc := atomic.AddUint32(&cnt, 1)

		if !initFlag && nc < uint32(defaultInitSendTx) {
			runReceiptCheckSendTx(config)()
		} else {
			initFlag = true

			// following logic can control the ratio between send/read task
			nc = nc % uint32(ratioReadPerSend+1)

			if nc == uint32(ratioReadPerSend) {
				runReceiptCheckSendTx(config)()
			} else {
				runReceiptCheckReadTx(config)()
			}
		}
	}
}

// transferAndCheck performs a transfer and checks the balance
func transferAndCheck(cli *client.KaiaClient, to *account.Account, from *account.Account, value *big.Int) error {
	ctx := context.Background()

	doubleLock(to, from)
	defer doubleUnlock(to, from)
	// The reason of saving balance of current accounts is to comparing with later balance.
	fromFormerBalance, _ := from.GetBalance(cli)
	toFormerBalance, _ := to.GetBalance(cli)

	hash, gasPrice, err := from.TransferSignedTxWithoutLock(cli, to, value, nil)
	if err != nil {
		return err
	}
	startTime := time.Now().Unix()
	var receipt *types.Receipt
	for {
		receipt, _ = cli.TransactionReceipt(ctx, hash)
		if receipt != nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
		if time.Now().Unix()-startTime > 100 {
			return errors.New("Time out : It took more than 100 seconds to make a block ")
		}
	}

	if to.GetAddress() == from.GetAddress() {
		value.SetUint64(0)
	}

	fromFormerBalance.Sub(fromFormerBalance, value)
	gasUsed := big.NewInt((int64)(receipt.GasUsed))
	fee := new(big.Int).Mul(gasUsed, gasPrice)
	fromFormerBalance.Sub(fromFormerBalance, fee)
	toFormerBalance.Add(toFormerBalance, value)

	startTime = time.Now().Unix()
	for {
		errFrom := from.CheckBalance(fromFormerBalance, cli)
		if errFrom != nil {
			log.Printf("from account : %s", errFrom.Error())
			time.Sleep(100 * time.Millisecond)
			if time.Now().Unix()-startTime > 10 {
				return errors.New("Time out (from) : It took more than 10 seconds to retrieve the correct receipt ")
			}
		} else {
			break
		}
	}

	if from.GetAddress() == to.GetAddress() {
		return nil
	}

	startTime = time.Now().Unix()
	for {
		errTo := to.CheckBalance(toFormerBalance, cli)
		if errTo != nil {
			log.Printf("to account : %s", errTo.Error())
			time.Sleep(100 * time.Millisecond)
			if time.Now().Unix()-startTime > 10 {
				return errors.New("Time out (to) : It took more than 10 seconds to retrieve the correct receipt ")
			}
		} else {
			break
		}
	}

	return nil
}

// RunTransferSignedWithCheckTC creates a closure for transfer signed with check test case
func RunTransferSignedWithCheckTC(config *TCConfig) func() {
	return func() {
		cli := config.CliPool.Alloc().(*client.KaiaClient)
		defer config.CliPool.Free(cli)

		from := config.AccGrp.GetAccountRandomly()
		to := config.AccGrp.GetAccountRandomly()

		value := big.NewInt(int64(rand.Int() % 3))
		start := boomer.Now()

		err := transferAndCheck(cli, to, from, value)

		elapsed := boomer.Now() - start

		if err == nil {
			boomer.Events.Publish("request_success", "http", "signedtransfer_with_check"+" to "+config.EndPoint, elapsed, int64(10))
		} else {
			boomer.Events.Publish("request_failure", "http", "signedtransfer_with_check"+" to "+config.EndPoint, elapsed, err.Error())
		}
	}
}
