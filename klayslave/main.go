package main

//go:generate abigen --sol cpuHeavyTC/CPUHeavy.sol --pkg cpuHeavyTC --out cpuHeavyTC/CPUHeavy.go
//go:generate abigen --sol userStorageTC/UserStorage.sol --pkg userStorageTC --out userStorageTC/UserStorage.go

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/kaiachain/kaia-load-tester/klayslave/account"
	"github.com/kaiachain/kaia-load-tester/klayslave/config"
	"github.com/kaiachain/kaia-load-tester/testcase"
	"github.com/myzhan/boomer"
	"github.com/urfave/cli"
)

var app = cli.NewApp()

func init() {
	app.Name = filepath.Base(os.Args[0])
	app.Usage = "This is for kaia load testing."
	app.Version = config.GetVersionWithCommit() // To see the version, run 'klayslave -v'
	app.HideVersion = false
	app.Copyright = "Copyright 2024 Kaia-load-tester authors"
	app.Flags = append(config.Flags, config.BoomerFlags...)

	// This app doesn't provide any subcommand
	//		app.Commands = []*cli.Command{}
	//		sort.Sort(cli.CommandsByName(app.Commands))
	//		app.CommandNotFound = nodecmd.CommandNotExist
	// app.OnUsageError = nodecmd.OnUsageError
	app.Before = func(cli *cli.Context) error {
		// runtime.GOMAXPROCS(runtime.NumCPU())
		if runtime.GOOS == "darwin" {
			return nil
		}
		return config.SetRLimit()
	}
	app.Action = RunAction
	app.After = func(cli *cli.Context) error {
		return nil
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func RunAction(ctx *cli.Context) {
	cfg := config.NewConfig(ctx)
	accGrp := account.NewAccGroup(cfg.GetChainID(), cfg.GetGasPrice(), cfg.GetBaseFee(), cfg.GetBatchSize(), cfg.InTheTcList("transferUnsignedTx"))
	var nUserForGaslessRevertTx, nUserForGaslessApproveTx int = 0, 0
	if cfg.InTheTcList("gaslessRevertTransactionTC") {
		nUserForGaslessRevertTx = cfg.GetNUserForSigned() // same as nUserForSignedTx
	}
	if cfg.InTheTcList("gaslessOnlyApproveTC") {
		nUserForGaslessApproveTx = cfg.GetNUserForSigned() // same as nUserForSignedTx
	}
	accGrp.CreateAccountsPerAccGrp(cfg.GetNUserForSigned(), cfg.GetNUserForUnsigned(), cfg.GetNUserForNewAccounts(), nUserForGaslessRevertTx, nUserForGaslessApproveTx, cfg.GetTcStrList(), cfg.GetGEndpoint())

	createTestAccGroupsAndPrepareContracts(cfg, accGrp)
	tasks := cfg.GetExtendedTasks()
	initializeTasks(cfg, accGrp, tasks)
	boomer.Run(toBoomerTasks(tasks)...)
}

// TODO-kaia-load-tester: remove global variables in the tc packages
func setSmartContractAddressPerPackage(a *account.AccGroup) {
}

// createTestAccGroupsAndPrepareContracts do every init steps before task.Init
// those steps are about deploying test contracts and
func createTestAccGroupsAndPrepareContracts(cfg *config.Config, accGrp *account.AccGroup) *account.Account {
	if len(cfg.GetChargeValue().Bits()) == 0 {
		return nil
	}

	// 1. Import global reservoir Account and create local reservoir account
	globalReservoirAccount := account.GetAccountFromKey(0, cfg.GetRichWalletPrivateKey())
	localReservoirAccount := account.NewAccount(0)

	// 2. charge local reservoir
	_ = globalReservoirAccount.GetNonce(cfg.GetGCli())
	revertGroupChargeValue := new(big.Int).Mul(cfg.GetChargeValue(), big.NewInt(int64(len(accGrp.GetAccListByName(account.AccListForGaslessRevertTx)))))
	approveGroupChargeValue := new(big.Int).Mul(cfg.GetChargeValue(), big.NewInt(int64(len(accGrp.GetAccListByName(account.AccListForGaslessApproveTx)))))
	tx := globalReservoirAccount.TransferSignedTxWithGuaranteeRetry(cfg.GetGCli(), localReservoirAccount, new(big.Int).Add(cfg.GetTotalChargeValue(), new(big.Int).Add(revertGroupChargeValue, approveGroupChargeValue)))
	receipt, err := bind.WaitMined(context.Background(), cfg.GetGCli(), tx)
	if err != nil {
		log.Fatalf("receipt failed, err:%v", err.Error())
	}
	if receipt.Status != 1 {
		log.Fatalf("transfer for reservoir failed, localReservoir")
	}

	// for _, token := range []string{"2", "3", "4", "5", "6", "7", "8", "9", "10"} {
	// 	tx = globalReservoirAccount.TransferTokenSignedTxWithGuaranteeRetry(cfg.GetGCli(), localReservoirAccount, new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1e18)), token)
	// 	receipt, err = bind.WaitMined(context.Background(), cfg.GetGCli(), tx)
	// 	if err != nil {
	// 		log.Fatalf("receipt failed, err:%v", err.Error())
	// 	}
	// 	if receipt.Status != 1 {
	// 		log.Fatalf("transfer for reservoir failed, localReservoir")
	// 	}
	// }

	// 3. charge KAIA
	if cfg.InTheTcList("transferTxTC") {
		log.Printf("Start charging KLAY to test accounts because transferTxTC is enabled")
		accs := accGrp.GetValidAccGrp()
		accs = append(accs, accGrp.GetAccListByName(account.AccListForGaslessRevertTx)...)  // for avoid validation
		accs = append(accs, accGrp.GetAccListByName(account.AccListForGaslessApproveTx)...) // for avoid validation
		account.HierarchicalDistribute(accs, localReservoirAccount, cfg.GetChargeValue(), func(from, to *account.Account, value *big.Int) {
			from.TransferSignedTxWithGuaranteeRetry(cfg.GetGCli(), to, value)
		})
		log.Printf("Finished charging KLAY to %d test account(s)\n", len(accs))
	} else if cfg.InTheTcList("tokenTransferTxTC") {
		log.Printf("Start charging all Tokens [2-10] to test accounts")
		accs := accGrp.GetValidAccGrp()
		accs = append(accs, accGrp.GetAccListByName(account.AccListForGaslessRevertTx)...)  // for avoid validation
		accs = append(accs, accGrp.GetAccListByName(account.AccListForGaslessApproveTx)...) // for avoid validation
		for _, token := range []string{"2", "3", "4", "5", "6", "7", "8", "9", "10"} {
			account.HierarchicalDistribute(accs, localReservoirAccount, new(big.Int).Mul(big.NewInt(1e10), big.NewInt(1e18)), func(from, to *account.Account, value *big.Int) {
				from.TransferTokenSignedTxWithGuaranteeRetry(cfg.GetGCli(), to, value, token)
			})
			log.Printf("Finished charging Token \"%s\" to %d test account(s)\n", token, len(accs))
		}
	} else {
		log.Printf("Start charging Token [2-3] to test accounts")
		accs := accGrp.GetValidAccGrp()
		accs = append(accs, accGrp.GetAccListByName(account.AccListForGaslessRevertTx)...)  // for avoid validation
		accs = append(accs, accGrp.GetAccListByName(account.AccListForGaslessApproveTx)...) // for avoid validation
		for _, token := range []string{"2", "3"} {
			account.HierarchicalDistribute(accs, localReservoirAccount, new(big.Int).Mul(big.NewInt(1e10), big.NewInt(1e18)), func(from, to *account.Account, value *big.Int) {
				from.TransferTokenSignedTxWithGuaranteeRetry(cfg.GetGCli(), to, value, token)
			})
			log.Printf("Finished charging Token \"%s\" to %d test account(s)\n", token, len(accs))
		}
	}

	// Wait, charge KAIA happen in 100% of all created test accounts
	// But, from here including prepareTestContracts like MintERC721, only 20% of account happens
	accGrp.SetAccGrpByActivePercent(cfg.GetActiveUserPercent())

	// Set SmartContractAddress value in each packages if needed
	setSmartContractAddressPerPackage(accGrp)
	return localReservoirAccount
}

func initializeTasks(cfg *config.Config, accGrp *account.AccGroup, tasks []*testcase.ExtendedTask) {
	println("Initializing tasks")

	// Tc package initializes the task
	for _, extendedTask := range tasks {
		accs := accGrp.GetAccListByName(account.AccListForSignedTx)
		if extendedTask.Name == "transferUnsignedTx" {
			accs = accGrp.GetAccListByName(account.AccListForUnsignedTx)
		} else if extendedTask.Name == "gaslessRevertTransactionTC" {
			accs = accGrp.GetAccListByName(account.AccListForGaslessRevertTx)
		} else if extendedTask.Name == "gaslessOnlyApproveTC" {
			accs = accGrp.GetAccListByName(account.AccListForGaslessApproveTx)
		}
		extendedTask.Init(accs, cfg.GetGEndpoint(), cfg.GetGasPrice())
		println("=> " + extendedTask.Name + " extendedTask is initialized.")
	}
}

func toBoomerTasks(tasks []*testcase.ExtendedTask) []*boomer.Task {
	var boomerTasks []*boomer.Task
	for _, task := range tasks {
		boomerTasks = append(boomerTasks, &boomer.Task{Weight: task.Weight, Fn: task.Fn, Name: task.Name})
	}
	return boomerTasks
}
