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
	"strings"
	"time"

	"github.com/klaytn/klaytn-load-tester/klayslave/account"
	"github.com/klaytn/klaytn-load-tester/task"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/api/debug"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/client"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/console"
	"github.com/klaytn/klaytn/crypto"
	"github.com/myzhan/boomer"
	"github.com/urfave/cli"
)

// Dedicated and fixed private key used to deploy a smart contract for ERC20 and ERC721 value transfer performance test.
var ERC20DeployPrivateKeyStr = "eb2c84d41c639178ff26a81f488c196584d678bb1390cc20a3aeb536f3969a98"
var ERC721DeployPrivateKeyStr = "45c40d95c9b7898a21e073b5bf952bcb05f2e70072e239a8bbd87bb74a53355e"

// prepareERC20Transfer sets up ERC20 transfer performance test.
func prepareERC20Transfer(cfg *task.Config) {
	if !cfg.InTheTcList("erc20TransferTC") {
		return
	}
	erc20DeployAcc := account.GetAccountFromKey(0, ERC20DeployPrivateKeyStr)
	log.Printf("prepareERC20Transfer", "addr", erc20DeployAcc.GetAddress().String())
	fmt.Println("localReservoir!!", cfg.GetLocalReservoir())
	cfg.GetLocalReservoir().TransferSignedTxWithGuaranteeRetry(cfg.GetGCli(), erc20DeployAcc, cfg.GetChargeValue())
	//chargeKLAYToTestAccounts(map[common.Address]*account.Account{erc20DeployAcc.GetAddress(): erc20DeployAcc})

	// A smart contract for ERC20 value transfer performance TC.
	erc20TransferTcSca := deploySingleSmartContract(erc20DeployAcc, erc20DeployAcc.DeployERC20, "ERC20 Performance Test Contract", cfg)

	// okay.. let's transfer erc20 tokens for further uses..
	cfg.GetLocalReservoir()
	firstChargeTokenToTestAccounts(erc20TransferTcSca.GetAddress(), erc20DeployAcc.TransferERC20, big.NewInt(1e11), cfg)
	chargeTokenToTestAccounts(erc20TransferTcSca.GetAddress(), cfg.GetLocalReservoir().TransferERC20, big.NewInt(1e4), cfg)
	task.SetErc20TransferTcContract(erc20TransferTcSca)
}

// prepareERC721Transfer sets up ERC721 transfer performance test.
func prepareERC721Transfer(cfg *task.Config) {
	if !cfg.InTheTcList("erc721TransferTC") {
		return
	}
	erc721DeployAcc := account.GetAccountFromKey(0, ERC721DeployPrivateKeyStr)
	log.Printf("prepareERC721Transfer", "addr", erc721DeployAcc.GetAddress().String())
	cfg.GetLocalReservoir().TransferSignedTxWithGuaranteeRetry(cfg.GetGCli(), erc721DeployAcc, cfg.GetChargeValue())
	//chargeKLAYToTestAccounts(map[common.Address]*account.Account{erc721DeployAcc.GetAddress(): erc721DeployAcc})

	// A smart contract for ERC721 value transfer performance TC.
	erc721TransferTcSca := deploySingleSmartContract(erc721DeployAcc, erc721DeployAcc.DeployERC721, "ERC721 Performance Test Contract", cfg)

	// Wait for reward tester to get started
	time.Sleep(30 * time.Second)

	cfg.GetLocalReservoir().MintERC721ToTestAccounts(cfg.GetGCli(), cfg.GetAllAccGrp(), erc721TransferTcSca.GetAddress(), 5)
	log.Println("MintERC721ToTestAccounts", "len(cfg.accGrp)", len(cfg.GetAllAccGrp()))

	task.SetErc721TransferTcContract(erc721TransferTcSca)
}

// Dedicated and fixed private key used to deploy a smart contract for storage trie write performance test.
var storageTrieDeployPrivateKeyStr = "3737c381633deaaa4c0bdbc64728f6ef7d381b17e1d30bbb74665839cec942b8"

// prepareStorageTrieWritePerformance sets up ERC20 storage trie write performance test.
func prepareStorageTrieWritePerformance(cfg *task.Config) {
	if !cfg.InTheTcList("storageTrieWriteTC") {
		return
	}
	storageTrieDeployAcc := account.GetAccountFromKey(0, storageTrieDeployPrivateKeyStr)

	log.Printf("prepareStorageTrieWritePerformance", "addr", storageTrieDeployAcc.GetAddress().String())
	cfg.GetLocalReservoir().TransferSignedTxWithGuaranteeRetry(cfg.GetGCli(), storageTrieDeployAcc, cfg.GetChargeValue())
	//chargeKLAYToTestAccounts(map[common.Address]*account.Account{storageTrieDeployAcc.GetAddress(): storageTrieDeployAcc})

	// A smart contract for storage trie store performance TC.
	storageTrieWriteTcSca := deploySingleSmartContract(storageTrieDeployAcc, storageTrieDeployAcc.DeployStorageTrieWrite, "Storage Trie Performance Test Contract", cfg)
	task.SetStorageTrieWriteTcContract(storageTrieWriteTcSca)
}

func prepareTestAccountsAndContracts(cfg *task.Config) {
	// accGrp map[common.Address]*account.Account,
	// First, charging KLAY to the test accounts.
	chargeKLAYToTestAccounts(cfg)

	// Second, deploy contracts used for some TCs.
	// If the test case is not on the list, corresponding contract won't be deployed.
	prepareERC20Transfer(cfg)
	prepareStorageTrieWritePerformance(cfg)

	// Third, deploy contracts for general tests.
	// A smart contract for general smart contract related TCs.
	GeneralSmartContract := deploySmartContract(cfg.GetLocalReservoir().TransferNewSmartContractDeployTxHumanReadable, "General Purpose Test Smart Contract", cfg)
	task.SetTcGeneralSmartContract(GeneralSmartContract)

	// Wait, charge needs to be happen in 100% of all created test accounts
	// But, from MintERC721, only 20% of account happens
	cfg.SetAccGrpByActivePercent()

	prepareERC721Transfer(cfg)
}

func chargeKLAYToTestAccounts(cfg *task.Config) {
	log.Printf("Start charging KLAY to test accounts")
	for _, acc := range cfg.GetAllAccGrp() {
		cfg.GetLocalReservoir().TransferSignedTxWithGuaranteeRetry(cfg.GetGCli(), acc, cfg.GetChargeValue())
	}
	log.Printf("Finished charging KLAY to %d test account(s)\n", len(cfg.GetAllAccGrp()))
}

type tokenChargeFunc func(initialCharge bool, c *client.Client, tokenContractAddr common.Address, recipient *account.Account, value *big.Int) (*types.Transaction, *big.Int, error)

// firstChargeTokenToTestAccounts charges initially generated tokens to local reservoir account for further testing.
// As this work is done simultaneously by different slaves, this should be done in "try and check" manner. // ??????? Anyway.. txpool can be full.. that's why it is done by "try and check" manner.
func firstChargeTokenToTestAccounts(tokenContractAddr common.Address, tokenChargeFn tokenChargeFunc, tokenChargeAmount *big.Int, cfg *task.Config) {
	log.Printf("Start initial token charging to the local revoir account")

	for {
		tx, _, err := tokenChargeFn(true, cfg.GetGCli(), tokenContractAddr, cfg.GetLocalReservoir(), tokenChargeAmount)
		for err != nil {
			log.Printf("Failed to execute %s: err %s", tx.Hash().String(), err.Error())
			time.Sleep(1 * time.Second) // Mostly, the err is `txpool is full`, retry after a while.
			continue
			//tx, _, err = tokenChargeFn(true, cfg.GetGCli(), tokenContractAddr, cfg.GetLocalReservoir(), tokenChargeAmount)
		}
		ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
		receipt, err := bind.WaitMined(ctx, cfg.GetGCli(), tx)
		cancelFn()
		if receipt != nil {
			break
		}
	}

	log.Printf("Finished initial ERC20 token charging to the local reservoir account")
}

// chargeTokenToTestAccounts charges default token to the test accounts for testing.
// As it is done independently among the slaves, it has simpler logic than firstChargeTokenToTestAccounts.
func chargeTokenToTestAccounts(tokenContractAddr common.Address, tokenChargeFn tokenChargeFunc, tokenChargeAmount *big.Int, cfg *task.Config) {
	log.Printf("Start charging tokens to test accounts")

	numChargedAcc := 0
	lastFailedNum := 0
	for _, recipientAccount := range cfg.GetAllAccGrp() {
		for {
			_, _, err := tokenChargeFn(false, cfg.GetGCli(), tokenContractAddr, recipientAccount, tokenChargeAmount)
			if err == nil {
				break // Success, move to next account.
			}
			numChargedAcc, lastFailedNum = estimateRemainingTime(cfg.GetAllAccGrp(), numChargedAcc, lastFailedNum)
		}
		numChargedAcc++
	}

	log.Printf("Finished charging tokens to %d test account(s), Total %d transactions are sent.\n", len(cfg.GetAllAccGrp()), numChargedAcc)
}

func estimateRemainingTime(accGrp []*account.Account, numChargedAcc, lastFailedNum int) (int, int) {
	if lastFailedNum > 0 {
		// Not 1st failed cases.
		TPS := (numChargedAcc - lastFailedNum) / 5 // TPS of only this slave during `txpool is full` situation.
		lastFailedNum = numChargedAcc

		if TPS <= 5 {
			log.Printf("Retry to charge test account #%d. But it is too slow. %d TPS\n", numChargedAcc, TPS)
		} else {
			remainTime := (len(accGrp) - numChargedAcc) / TPS
			remainHour := remainTime / 3600
			remainMinute := (remainTime % 3600) / 60

			log.Printf("Retry to charge test account #%d. Estimated remaining time: %d hours %d mins later\n", numChargedAcc, remainHour, remainMinute)
		}
	} else {
		// 1st failed case.
		lastFailedNum = numChargedAcc
		log.Printf("Retry to charge test account #%d.\n", numChargedAcc)
	}
	time.Sleep(5 * time.Second) // Mostly, the err is `txpool is full`, retry after a while.
	return numChargedAcc, lastFailedNum
}

type contractDeployFunc func(c *client.Client, to *account.Account, value *big.Int, humanReadable bool) (common.Address, *types.Transaction, *big.Int, error)

// deploySmartContract deploys smart contracts by the number of locust slaves.
// In other words, each slave owns its own contract for testing.
func deploySmartContract(contractDeployFn contractDeployFunc, contractName string, cfg *task.Config) *account.Account {
	addr, lastTx, _, err := contractDeployFn(cfg.GetGCli(), nil, common.Big0, false)
	for err != nil {
		log.Printf("Failed to deploy a %s: err %s", contractName, err.Error())
		time.Sleep(5 * time.Second) // Mostly, the err is `txpool is full`, retry after a while.
		addr, lastTx, _, err = contractDeployFn(cfg.GetGCli(), nil, common.Big0, false)
	}

	log.Printf("Start waiting the receipt of the %s tx(%v).\n", contractName, lastTx.Hash().String())
	bind.WaitMined(context.Background(), cfg.GetGCli(), lastTx)

	deployedContract := account.NewKlaytnAccountWithAddr(1, addr)
	log.Printf("%s has been deployed to : %s\n", contractName, addr.String())
	return deployedContract
}

// deploySingleSmartContract deploys only one smart contract among the slaves.
// It the contract is already deployed by other slave, it just calculates the address of the contract.
func deploySingleSmartContract(erc20DeployAcc *account.Account, contractDeployFn contractDeployFunc, contractName string, cfg *task.Config) *account.Account {
	addr, lastTx, _, err := contractDeployFn(cfg.GetGCli(), nil, common.Big0, false)
	for err != nil {
		if err == account.AlreadyDeployedErr {
			erc20Addr := crypto.CreateAddress(erc20DeployAcc.GetAddress(), 0)
			return account.NewKlaytnAccountWithAddr(1, erc20Addr)
		}
		if strings.HasPrefix(err.Error(), "known transaction") {
			erc20Addr := crypto.CreateAddress(erc20DeployAcc.GetAddress(), 0)
			return account.NewKlaytnAccountWithAddr(1, erc20Addr)
		}
		log.Printf("Failed to deploy a %s: err %s", contractName, err.Error())
		time.Sleep(5 * time.Second) // Mostly, the err is `txpool is full`, retry after a while.
		addr, lastTx, _, err = contractDeployFn(cfg.GetGCli(), nil, common.Big0, false)
	}

	log.Printf("Start waiting the receipt of the %s tx(%v).\n", contractName, lastTx.Hash().String())
	bind.WaitMined(context.Background(), cfg.GetGCli(), lastTx)

	deployedContract := account.NewKlaytnAccountWithAddr(1, addr)
	log.Printf("%s has been deployed to : %s\n", contractName, addr.String())
	return deployedContract
}

func createAndChargeLocalReservoirAccount(cfg *task.Config) {
	//totalChargeValue := new(big.Int)
	//totalChargeValue.Mul(chargeValue, big.NewInt(int64(nUserForUnsigned+nUserForSigned+nUserForNewAccounts+1)))

	// Import coinbase Account
	globalReservoirAccount := account.GetAccountFromKey(0, cfg.GetRichWalletPrivateKey())
	cfg.SetLocalReservoirAccount(account.NewAccount(0))

	if len(cfg.GetChargeValue().Bits()) != 0 {
		for {
			globalReservoirAccount.GetNonceFromBlock(cfg.GetGCli())
			hash, _, err := globalReservoirAccount.TransferSignedTx(cfg.GetGCli(), cfg.GetLocalReservoir(), cfg.GetTotalChargeValue())
			if err != nil {
				log.Printf("%v: charge newCoinbase fail: %v\n", os.Getpid(), err)
				time.Sleep(1000 * time.Millisecond)
				continue
			}

			log.Printf("%v : charge new local reservoir account: %v, Txhash=%v\n", os.Getpid(), cfg.GetLocalReservoir().GetAddress().String(), hash.String())

			getReceipt := false
			// After this loop waiting for 10 sec, It will retry to charge with new nonce.
			// it means another node stole the nonce.
			for i := 0; i < 5; i++ {
				time.Sleep(2000 * time.Millisecond)
				ctx := context.Background()

				//_, err := gCli.TransactionReceipt(ctx, hash)
				//if err != nil {
				//	getReceipt = true
				//	log.Printf("%v : charge newCoinbase success: %v\n", os.Getpid(), newCoinbase.GetAddress().String())
				//	break
				//}
				//log.Printf("%v : charge newCoinbase waiting: %v\n", os.Getpid(), newCoinbase.GetAddress().String())

				val, err := cfg.GetGCli().BalanceAt(ctx, cfg.GetLocalReservoir().GetAddress(), nil)
				if err == nil {
					if val.Cmp(big.NewInt(0)) == 1 {
						getReceipt = true
						log.Printf("%v : charge newCoinbase success: %v, balance=%v peb\n", os.Getpid(), cfg.GetLocalReservoir().GetAddress().String(), val.String())
						break
					}
					log.Printf("%v : charge newCoinbase waiting: %v\n", os.Getpid(), cfg.GetLocalReservoir().GetAddress().String())
				} else {
					log.Printf("%v : check balance err: %v\n", os.Getpid(), err)
				}
			}

			if getReceipt {
				break
			}
		}
	}
}

var app = cli.NewApp()

func init() {
	app.Name = filepath.Base(os.Args[0])
	app.Usage = "This is for kaia load testing."
	app.Version = task.GetVersionWithCommit() // To see the version, run 'klayslave -v'
	app.HideVersion = false
	app.Copyright = "Copyright 2024 Kaia-load-tester authors"
	app.Flags = append(task.Flags, task.BoomerFlags...)

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
		task.SetRLimit()
		return nil
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
	config := task.NewConfig(ctx)

	// Set coinbase & Create Test Account
	createAndChargeLocalReservoirAccount(config)
	config.InitAccGrp()
	config.UnlockAccGrpForUnsignedIfNeeded()

	if len(config.GetChargeValue().Bits()) != 0 {
		prepareTestAccountsAndContracts(config)
	}

	println("Initializing tasks")

	// Locust Slave Run
	config.InitTasks()
	boomer.Run(config.GetBoomerTasksList()...)
	//boomer.Run(cpuHeavyTx)
}
