package testcase

import (
	"math/big"
	"math/rand"

	"github.com/kaiachain/kaia-load-tester/klayslave/account"
	"github.com/kaiachain/kaia/client"
	"github.com/myzhan/boomer"
)

// ValueTransferTxFunc represents a value transfer transaction function signature
type ValueTransferTxFunc = func(*client.KaiaClient, *account.Account, *account.Account, *big.Int) (interface{}, *big.Int, error)

// RunBaseValueTransfer creates a closure that executes a test case with common logic
func RunBaseValueTransfer(config *TCConfig, txFunc ValueTransferTxFunc) func() {
	return func() {
		cli := config.CliPool.Alloc().(*client.KaiaClient)
		defer config.CliPool.Free(cli)

		from := config.AccGrp.GetAccountRandomly()
		to := config.AccGrp.GetAccountRandomly()
		value := big.NewInt(int64(rand.Int() % 3))

		start := boomer.Now()
		_, _, err := txFunc(cli, from, to, value)
		elapsed := boomer.Now() - start

		if err == nil {
			boomer.Events.Publish("request_success", "http", config.Name+" to "+config.EndPoint, elapsed, int64(10))
		} else {
			boomer.Events.Publish("request_failure", "http", config.Name+" to "+config.EndPoint, elapsed, err.Error())
		}
	}
}

// Run functions for each test case
func RunNewValueTransferTC(config *TCConfig) func() {
	txFunc := func(cli *client.KaiaClient, from *account.Account, to *account.Account, value *big.Int) (interface{}, *big.Int, error) {
		return from.TransferNewValueTransferTx(cli, to, value)
	}
	return RunBaseValueTransfer(config, txFunc)
}

func RunNewFeeDelegatedValueTransferTC(config *TCConfig) func() {
	txFunc := func(cli *client.KaiaClient, from *account.Account, to *account.Account, value *big.Int) (interface{}, *big.Int, error) {
		return from.TransferNewFeeDelegatedValueTransferTx(cli, to, value)
	}
	return RunBaseValueTransfer(config, txFunc)
}

func RunNewFeeDelegatedValueTransferWithRatioTC(config *TCConfig) func() {
	txFunc := func(cli *client.KaiaClient, from *account.Account, to *account.Account, value *big.Int) (interface{}, *big.Int, error) {
		return from.TransferNewFeeDelegatedValueTransferWithRatioTx(cli, to, value)
	}
	return RunBaseValueTransfer(config, txFunc)
}

func RunNewCancelTC(config *TCConfig) func() {
	txFunc := func(cli *client.KaiaClient, from *account.Account, to *account.Account, value *big.Int) (interface{}, *big.Int, error) {
		return from.TransferNewCancelTx(cli, to, value)
	}
	return RunBaseValueTransfer(config, txFunc)
}

func RunNewValueTransferWithCancelTC(config *TCConfig) func() {
	txFunc := func(cli *client.KaiaClient, from *account.Account, to *account.Account, value *big.Int) (interface{}, *big.Int, error) {
		return from.TransferNewValueTransferWithCancelTx(cli, to, value)
	}
	return RunBaseValueTransfer(config, txFunc)
}

func RunNewFeeDelegatedCancelTC(config *TCConfig) func() {
	txFunc := func(cli *client.KaiaClient, from *account.Account, to *account.Account, value *big.Int) (interface{}, *big.Int, error) {
		return from.TransferNewFeeDelegatedCancelTx(cli, to, value)
	}
	return RunBaseValueTransfer(config, txFunc)
}

func RunNewFeeDelegatedCancelWithRatioTC(config *TCConfig) func() {
	txFunc := func(cli *client.KaiaClient, from *account.Account, to *account.Account, value *big.Int) (interface{}, *big.Int, error) {
		return from.TransferNewFeeDelegatedCancelWithRatioTx(cli, to, value)
	}
	return RunBaseValueTransfer(config, txFunc)
}

func RunNewSmartContractDeployTC(config *TCConfig) func() {
	txFunc := func(cli *client.KaiaClient, from *account.Account, to *account.Account, value *big.Int) (interface{}, *big.Int, error) {
		newTo := account.NewKaiaAccount(0)
		_, _, _, err := from.TransferNewSmartContractDeployTx(cli, newTo, value, account.TestContractInfos[account.ContractGeneral].Bytecode, false)
		return nil, nil, err
	}
	return RunBaseValueTransfer(config, txFunc)
}

func RunNewFeeDelegatedSmartContractDeployTC(config *TCConfig) func() {
	txFunc := func(cli *client.KaiaClient, from *account.Account, to *account.Account, value *big.Int) (interface{}, *big.Int, error) {
		return from.TransferNewFeeDelegatedSmartContractDeployTx(cli, to, value)
	}
	return RunBaseValueTransfer(config, txFunc)
}

func RunNewFeeDelegatedSmartContractDeployWithRatioTC(config *TCConfig) func() {
	txFunc := func(cli *client.KaiaClient, from *account.Account, to *account.Account, value *big.Int) (interface{}, *big.Int, error) {
		return from.TransferNewFeeDelegatedSmartContractDeployWithRatioTx(cli, to, value)
	}
	return RunBaseValueTransfer(config, txFunc)
}

func RunNewValueTransferMemoTC(config *TCConfig) func() {
	txFunc := func(cli *client.KaiaClient, from *account.Account, to *account.Account, value *big.Int) (interface{}, *big.Int, error) {
		return from.TransferNewValueTransferMemoTx(cli, to, value)
	}
	return RunBaseValueTransfer(config, txFunc)
}

func RunNewFeeDelegatedValueTransferMemoTC(config *TCConfig) func() {
	txFunc := func(cli *client.KaiaClient, from *account.Account, to *account.Account, value *big.Int) (interface{}, *big.Int, error) {
		return from.TransferNewFeeDelegatedValueTransferMemoTx(cli, to, value)
	}
	return RunBaseValueTransfer(config, txFunc)
}

func RunNewFeeDelegatedValueTransferMemoWithRatioTC(config *TCConfig) func() {
	txFunc := func(cli *client.KaiaClient, from *account.Account, to *account.Account, value *big.Int) (interface{}, *big.Int, error) {
		return from.TransferNewFeeDelegatedValueTransferMemoWithRatioTx(cli, to, value)
	}
	return RunBaseValueTransfer(config, txFunc)
}

func RunNewValueTransferLargeMemoTC(config *TCConfig) func() {
	txFunc := func(cli *client.KaiaClient, from *account.Account, to *account.Account, value *big.Int) (interface{}, *big.Int, error) {
		return from.TransferNewValueTransferLargeMemoTx(cli, to, value)
	}
	return RunBaseValueTransfer(config, txFunc)
}

func RunNewValueTransferSmallMemoTC(config *TCConfig) func() {
	txFunc := func(cli *client.KaiaClient, from *account.Account, to *account.Account, value *big.Int) (interface{}, *big.Int, error) {
		return from.TransferNewValueTransferSmallMemoTx(cli, to, value)
	}
	return RunBaseValueTransfer(config, txFunc)
}

func RunNewAccountUpdateTC(config *TCConfig) func() {
	txFunc := func(cli *client.KaiaClient, from *account.Account, to *account.Account, value *big.Int) (interface{}, *big.Int, error) {
		return from.TransferNewAccountUpdateTx(cli, to, value)
	}
	return RunBaseValueTransfer(config, txFunc)
}

func RunNewFeeDelegatedAccountUpdateTC(config *TCConfig) func() {
	txFunc := func(cli *client.KaiaClient, from *account.Account, to *account.Account, value *big.Int) (interface{}, *big.Int, error) {
		return from.TransferNewFeeDelegatedAccountUpdateTx(cli, to, value)
	}
	return RunBaseValueTransfer(config, txFunc)
}

func RunNewFeeDelegatedAccountUpdateWithRatioTC(config *TCConfig) func() {
	txFunc := func(cli *client.KaiaClient, from *account.Account, to *account.Account, value *big.Int) (interface{}, *big.Int, error) {
		return from.TransferNewFeeDelegatedAccountUpdateWithRatioTx(cli, to, value)
	}
	return RunBaseValueTransfer(config, txFunc)
}

func RunTransferSignedTC(config *TCConfig) func() {
	txFunc := func(cli *client.KaiaClient, from *account.Account, to *account.Account, value *big.Int) (interface{}, *big.Int, error) {
		return from.TransferSignedTx(cli, to, value, nil)
	}
	return RunBaseValueTransfer(config, txFunc)
}

func RunTransferUnsignedTC(config *TCConfig) func() {
	txFunc := func(cli *client.KaiaClient, from *account.Account, to *account.Account, value *big.Int) (interface{}, *big.Int, error) {
		_, err := from.TransferUnsignedTx(cli, to, value)
		return nil, nil, err
	}
	return RunBaseValueTransfer(config, txFunc)
}
