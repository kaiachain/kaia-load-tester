package newCancelTC

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"math/rand"

	"github.com/kaiachain/kaia-load-tester/klayslave/account"
	"github.com/kaiachain/kaia-load-tester/klayslave/clipool"
	"github.com/kaiachain/kaia/client"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/params"
	"github.com/myzhan/boomer"
)

const Name = "newCancelTC"

var (
	endPoint string
	nAcc     int
	accGrp   []*account.Account
	cliPool  clipool.ClientPool
	gasPrice *big.Int

	// multinode tester
	transferedValue *big.Int
	expectedFee     *big.Int

	fromAccount     *account.Account
	prevBalanceFrom *big.Int

	toAccount     *account.Account
	prevBalanceTo *big.Int
)

func Init(accs []*account.Account, endpoint string, gp *big.Int) {
	gasPrice = gp

	endPoint = endpoint

	cliCreate := func() interface{} {
		c, err := client.Dial(endPoint)
		if err != nil {
			log.Fatalf("Failed to connect RPC: %v", err)
		}
		return c
	}

	cliPool.Init(20, 300, cliCreate)

	for _, acc := range accs {
		accGrp = append(accGrp, acc)
	}

	nAcc = len(accGrp)
}

func Run() {
	cli := cliPool.Alloc().(*client.Client)

	from := accGrp[rand.Int()%nAcc]
	to := accGrp[rand.Int()%nAcc]
	value := big.NewInt(int64(rand.Int() % 3))

	start := boomer.Now()
	_, _, err := from.TransferNewCancelTx(cli, to, value)
	elapsed := boomer.Now() - start

	if err == nil {
		boomer.Events.Publish("request_success", "http", "transferNewCancelTx"+" to "+endPoint, elapsed, int64(10))
		cliPool.Free(cli)
	} else {
		boomer.Events.Publish("request_failure", "http", "transferNewCancelTx"+" to "+endPoint, elapsed, err.Error())
	}
}

func RunSingle() (txHash common.Hash, err error) {
	cli := cliPool.Alloc().(*client.Client)
	defer cliPool.Free(cli)

	fromIdx := rand.Int() % nAcc
	toIdx := (fromIdx + 1) % nAcc

	from := accGrp[fromIdx]
	to := accGrp[toIdx]
	value := big.NewInt(int64(rand.Int() % 3))
	fmt.Printf("[TC] transferSignedTc: %v, from:%v, to:%v, value:%v\n", endPoint, from.GetAddress().String(), to.GetAddress().String(), value)
	transferedValue = big.NewInt(value.Int64())
	expectedFee = big.NewInt(0).Mul(big.NewInt(25*params.Gkei), big.NewInt(21000))

	balance, err := from.GetBalance(cli)
	if err != nil {
		return common.Hash{}, err
	}
	fromAccount = from
	prevBalanceFrom = big.NewInt(balance.Int64())
	fmt.Printf("From:%v, balance:%v\n", fromAccount.GetAddress().String(), prevBalanceFrom.Int64())

	balance, err = to.GetBalance(cli)
	if err != nil {
		return common.Hash{}, err
	}
	toAccount = to
	prevBalanceTo = big.NewInt(balance.Int64())
	fmt.Printf("To:%v, balance:%v\n", toAccount.GetAddress().String(), prevBalanceTo.Int64())

	txHash, _, err = from.TransferNewCancelTx(cli, to, value)
	if err != nil {
		return common.Hash{}, err
	}

	return txHash, err
}

// CheckResult returns true and nil error, if expected results are observed.
// Otherewise returns false and error.
func CheckResult() (bool, error) {
	cli := cliPool.Alloc().(*client.Client)
	defer cliPool.Free(cli)

	balance, err := fromAccount.GetBalance(cli)
	if err != nil {
		return false, err
	}
	expectedBalance := big.NewInt(0)
	expectedBalance.Sub(prevBalanceFrom, transferedValue)
	expectedBalance.Sub(expectedBalance, expectedFee)
	// fmt.Printf("prevBalanceFrom=%v, transferedValue=%v, expectedFee=%v\n", prevBalanceFrom.Int64(), transferedValue.Int64(), expectedFee.Int64())

	if expectedBalance.Int64() != balance.Int64() {
		fmt.Printf("[FAILED] From account address=%v, Expected balance=%v, Actual balance=%v\n", fromAccount.GetAddress().String(), expectedBalance.Int64(), balance.Int64())
		return false, errors.New("Balance mismatched!")
	} else {
		fmt.Printf("[PASSED] From account address=%v, Expected balance=%v, Actual balance=%v\n", fromAccount.GetAddress().String(), expectedBalance.Int64(), balance.Int64())
	}

	balance, err = toAccount.GetBalance(cli)
	if err != nil {
		return false, err
	}
	expectedBalance = big.NewInt(0)
	expectedBalance.Add(prevBalanceTo, transferedValue)
	fmt.Printf("prevBalanceTo=%v, transferedValue=%v\n", prevBalanceTo.Int64(), transferedValue.Int64())
	if expectedBalance.Int64() != balance.Int64() {
		fmt.Printf("[FAILED] To account address=%v, Expected balance=%v, Actual balance=%v\n", toAccount.GetAddress().String(), expectedBalance.Int64(), balance.Int64())
		return false, errors.New("Balance mismatched!")
	} else {
		fmt.Printf("[PASSED] To account address=%v, Expected balance=%v, Actual balance=%v\n", toAccount.GetAddress().String(), expectedBalance.Int64(), balance.Int64())
	}

	return true, err
}
