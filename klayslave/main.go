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

	"github.com/kaiachain/kaia-load-tester/klayslave/account"
	"github.com/kaiachain/kaia-load-tester/klayslave/config"
	"github.com/kaiachain/kaia-load-tester/testcase"
	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/api/debug"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/console"
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
		//runtime.GOMAXPROCS(runtime.NumCPU())
		if runtime.GOOS == "darwin" {
			return nil
		}
		return config.SetRLimit()
	}
	app.Action = RunAction
	app.After = func(cli *cli.Context) error {
		debug.Exit()
		console.Stdin.Close() // Resets terminal mode.
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
	tcList := cfg.GetTcStrList()
	accGrp := account.NewAccGroup(cfg.GetChainID(), cfg.GetGasPrice(), cfg.GetBaseFee(), account.ContainsAnyInList(tcList, []string{"transferUnsignedTx"}))
	var nUserForGaslessRevertTx, nUserForGaslessApproveTx int = 0, 0
	if account.ContainsAnyInList(tcList, []string{"gaslessRevertTransactionTC"}) {
		nUserForGaslessRevertTx = cfg.GetNUserForSigned() // same as nUserForSignedTx
	}
	if account.ContainsAnyInList(tcList, []string{"gaslessOnlyApproveTC"}) {
		nUserForGaslessApproveTx = cfg.GetNUserForSigned() // same as nUserForSignedTx
	}
	accGrp.CreateAccountsPerAccGrp(cfg.GetNUserForSigned(), cfg.GetNUserForUnsigned(), cfg.GetNUserForNewAccounts(), nUserForGaslessRevertTx, nUserForGaslessApproveTx, cfg.GetNUserForTC(), cfg.GetTcStrList(), cfg.GetGEndpoint())

	createTestAccGroupsAndPrepareContracts(cfg, accGrp)

	// Initialize refactored test cases (after contracts are deployed)
	boomerTasks := initializeTasks(cfg, accGrp, cfg.GetExtendedTasks())
	boomer.Run(boomerTasks...)
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
	forAuctionDepositChargeValue := new(big.Int).Mul(cfg.GetChargeValue(), big.NewInt(int64(len(accGrp.GetAccListByName(account.AccListForSignedTx)))))
	initialLiquidity := common.Big0
	if !account.IsGSRExistInRegistry(cfg.GetGCli(), nil) {
		// If GSR does not exist, charge initial liquidity to the local reservoir
		initialLiquidity = account.GetInitialLiquidity()
	}
	totalChargeValue := new(big.Int).Add(cfg.GetTotalChargeValue(), new(big.Int).Add(initialLiquidity, new(big.Int).Add(forAuctionDepositChargeValue, new(big.Int).Add(revertGroupChargeValue, approveGroupChargeValue))))
	tx := globalReservoirAccount.TransferSignedTxWithGuaranteeRetry(cfg.GetGCli(), localReservoirAccount, totalChargeValue)
	receipt, err := bind.WaitMined(context.Background(), cfg.GetGCli(), tx)
	if err != nil {
		log.Fatalf("receipt failed, err:%v", err.Error())
	}
	if receipt.Status != 1 {
		log.Fatalf("transfer for reservoir failed, localReservoir")
	}

	// 3. charge KAIA
	log.Printf("Start charging KLAY to test accounts")
	accs := accGrp.GetValidAccGrp()
	accs = append(accs, accGrp.GetAccListByName(account.AccListForGaslessRevertTx)...)  // for avoid validation
	accs = append(accs, accGrp.GetAccListByName(account.AccListForGaslessApproveTx)...) // for avoid validation
	account.ConcurrentTransactionSend(accs, cfg.GetChargeParallelNum(), func(_ int, acc *account.Account) {
		localReservoirAccount.TransferSignedTxWithGuaranteeRetry(cfg.GetGCli(), acc, cfg.GetChargeValue())
	})
	log.Printf("Finished charging KLAY to %d test account(s)\n", len(accs))

	// Wait, charge KAIA happen in 100% of all created test accounts
	// But, from here including prepareTestContracts like MintERC721, only 20% of account happens
	accGrp.SetAccGrpByActivePercent(cfg.GetActiveUserPercent())

	// 4. Deploy the test contracts which will be used in various TCs. If needed, charge tokens to test accounts.
	// GSR setup, Auction Entry Point registration, and Deposit are also done in DoAdditionalWork.
	// All slaves call DeployTestContracts, but only leader does additional work (token charging, minting, etc.)
	// Contracts are deployed only if deployer nonce is 0 (not yet deployed)
	accGrp.DeployTestContracts(cfg.GetGCli(), cfg.GetChargeValue(), cfg.GetChargeParallelNum(), cfg.GetTcStrList(), cfg.GetAuctionTargetTxTypeList(), localReservoirAccount, globalReservoirAccount, cfg.IsLeaderSlave())

	return localReservoirAccount
}

func initializeTasks(cfg *config.Config, accGrp *account.AccGroup, tasks []*testcase.ExtendedTask) []*boomer.Task {
	println("Initializing tasks")
	var boomerTasks []*boomer.Task

	// Tc package initializes the task
	for _, extendedTask := range tasks {
		config := extendedTask.Init(accGrp, cfg.GetGEndpoint(), extendedTask.TestContracts, extendedTask.Name, cfg.GetAuctionTargetTxTypeList())
		boomerTask := &boomer.Task{
			Name:   extendedTask.Name,
			Weight: extendedTask.Weight,
			Fn:     extendedTask.Run(config),
		}
		boomerTasks = append(boomerTasks, boomerTask)
		println("=> " + extendedTask.Name + " extendedTask is initialized.")
	}
	return boomerTasks
}
