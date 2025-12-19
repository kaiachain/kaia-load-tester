package account

import (
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"sync"
	"time"

	"github.com/kaiachain/kaia/client"
)

// AccList defines the enum for accList
type AccList int

const (
	AccListForSignedTx AccList = iota
	AccListForUnsignedTx
	AccListForNewAccounts
	AccListForGaslessRevertTx
	AccListForGaslessApproveTx
	AccListForTC
	AccListEnd
)

// TestContract defines the enum for TestContract
type TestContract int

const (
	ContractErc20 TestContract = iota
	ContractErc721
	ContractStorageTrie
	ContractGeneral
	ContractGaslessToken
	ContractWKaia
	ContractUniswapV2Factory
	ContractUniswapV2Router
	ContractGaslessSwapRouter
	ContractCounterForTestAuction
	ContractAuctionFeeVault
	ContractAuctionDepositVault
	ContractAuctionEntryPoint
	ContractCPUHeavy
	ContractLargeMemo
	ContractReadApiCallContract
	ContractUserStorage
	ContractInternalTxKIP17
	ContractInternalTxMain
	ContractTetherLogic
	ContractTetherProxy
	ContractEnd
)

type AccLoader func(*AccGroup)

type AccGroup struct {
	containsUnsignedAccGrp bool

	accLists  [][]*Account
	contracts []*Account
}

func NewAccGroup(chainId *big.Int, gasPrice *big.Int, baseFee *big.Int, contains bool) *AccGroup {
	SetChainID(chainId)
	SetGasPrice(gasPrice)
	SetBaseFee(baseFee)

	return &AccGroup{
		containsUnsignedAccGrp: contains,
		accLists:               make([][]*Account, AccListEnd),
		contracts:              make([]*Account, ContractEnd),
	}
}
func (a *AccGroup) Load(loader AccLoader) { loader(a) }

func (a *AccGroup) GetTestContractList() []*Account               { return a.contracts }
func (a *AccGroup) GetTestContractByName(t TestContract) *Account { return a.contracts[t] }
func (a *AccGroup) GetAccListByName(t AccList) []*Account         { return a.accLists[t] }

func (a *AccGroup) SetTestContractByName(c *Account, t TestContract) { a.contracts[t] = c }
func (a *AccGroup) SetAccListByName(accs []*Account, t AccList) {
	for _, acc := range accs {
		a.accLists[t] = append(a.accLists[t], acc)
	}
}
func (a *AccGroup) AddAccToListByName(acc *Account, t AccList) {
	a.accLists[t] = append(a.accLists[t], acc)
}
func (a *AccGroup) CreateAccountsPerAccGrp(nUserForSignedTx int, nUserForUnsignedTx int, nUserForNewAccounts int, nUserForGaslessRevertTx int, nUserForGaslessApproveTx int, nUserForTC int, tcStrList []string, gEndpoint string) {
	for idx, nUser := range []int{nUserForSignedTx, nUserForUnsignedTx, nUserForNewAccounts, nUserForGaslessRevertTx, nUserForGaslessApproveTx, nUserForTC} {
		println(idx, " Account Group Preparation...")
		for i := 0; i < nUser; i++ {
			account := NewAccount(i)
			a.AddAccToListByName(account, AccList(idx))
			fmt.Printf("%v\n", account.address.String())
		}
	}

	// Unlock AccGrpForUnsignedTx if needed
	for _, tcName := range tcStrList {
		if tcName != "transferUnsignedTx" {
			continue
		}
		// If at least one task needs unlocking, unlock the accGrp for unsignedTx.
		for _, acc := range a.GetAccListByName(AccListForUnsignedTx) {
			acc.ImportUnLockAccount(gEndpoint)
		}
		break
	}
}

func (a *AccGroup) SetAccGrpByActivePercent(activeUserPercent int) {
	for i, accGrp := range a.accLists {
		numActiveAccGrpForSignedTx := len(accGrp) * activeUserPercent / 100
		if numActiveAccGrpForSignedTx == 0 {
			a.accLists[i] = nil
			continue
		}

		a.accLists[i] = accGrp[:numActiveAccGrpForSignedTx]
	}
}

func (a *AccGroup) GetValidAccGrp() []*Account {
	var accGrp []*Account
	for _, acc := range a.GetAccListByName(AccListForSignedTx) {
		accGrp = append(accGrp, acc)
	}
	//if !a.cfg.InTheTcList("transferUnsignedTx") {
	if !a.containsUnsignedAccGrp {
		return accGrp
	}
	for _, acc := range a.GetAccListByName(AccListForUnsignedTx) {
		accGrp = append(accGrp, acc)
	}
	return accGrp
}

func ContainsAnyInList(list []string, targets []string) bool {
	for _, target := range targets {
		for _, item := range list {
			if item == target {
				return true
			}
		}
	}
	return false
}

func (a *AccGroup) DeployTestContracts(gCli *client.Client, chargeValue *big.Int, maxConcurrency int, tcList []string, targetTxTypeList []string, localReservoir *Account, globalReservoir *Account, isLeader bool) {
	ctx := &AdditionalWorkContext{
		GCli:             gCli,
		LocalReservoir:   localReservoir,
		GlobalReservoir:  globalReservoir,
		ChargeValue:      chargeValue,
		IsLeader:         isLeader,
		MaxConcurrency:   maxConcurrency,
		AccGrp:           a,
		TcList:           tcList,
		TargetTxTypeList: targetTxTypeList,
	}

	// Phase 1: Deploy contracts and run setup work (leader only)
	for idx, info := range TestContractInfos {
		testContractType := TestContract(idx)
		if testContractType != ContractGeneral && !ContainsAnyInList(tcList, info.testNames) && !ContainsAnyInList(targetTxTypeList, info.auctionTargetTxTypeList) {
			continue
		}

		if isLeader {
			isAlreadyDeployed := info.IsDeployed(gCli, info.deployer)
			if !isAlreadyDeployed {
				localReservoir.TransferSignedTxWithGuaranteeRetry(gCli, info.deployer, chargeValue)
				a.contracts[idx] = info.deployer.SmartContractDeployWithGuaranteeRetry(gCli, info.GetBytecodeWithConstructorParam(info.Bytecode, a.contracts, info.deployer), info.contractName, true)
			} else {
				a.contracts[idx] = NewKaiaAccountWithAddr(0, info.GetAddress(gCli, info.deployer))
			}

			// Setup work - only leader executes (e.g., GSR registration, Auction registration)
			if info.DoSetupWork != nil {
				info.DoSetupWork(ctx)
			}
		}
	}

	// Phase 2: Wait for deploy/setup completion, get contract address, and run charging work
	for idx, info := range TestContractInfos {
		testContractType := TestContract(idx)
		if testContractType != ContractGeneral && !ContainsAnyInList(tcList, info.testNames) && !ContainsAnyInList(targetTxTypeList, info.auctionTargetTxTypeList) {
			continue
		}

		if !isLeader {
			// Wait for deploy completion (use WaitForSetup if available, otherwise use IsDeployed)
			if info.WaitForSetup != nil {
				WaitForSetupCompletion(gCli, info.WaitForSetup, info.contractName)
			} else {
				WaitForSetupCompletion(gCli, func(gCli *client.Client) bool {
					return info.IsDeployed(gCli, info.deployer)
				}, info.contractName)
			}

			// Get contract address after deploy completion
			a.contracts[idx] = NewKaiaAccountWithAddr(0, info.GetAddress(gCli, info.deployer))
		}

		// Charging work - all slaves execute (e.g., token charging, NFT minting, deposit)
		if info.DoChargingWork != nil {
			info.DoChargingWork(ctx)
		}
	}
}

// WaitForSetupCompletion polls on-chain state until setup is complete
func WaitForSetupCompletion(gCli *client.Client, waitFn func(*client.Client) bool, contractName string) {
	const maxRetries = 60
	const retryInterval = 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		if waitFn(gCli) {
			log.Printf("Setup for %s is complete, proceeding with charging work", contractName)
			return
		}
		log.Printf("Waiting for %s setup to complete... (attempt %d/%d)", contractName, i+1, maxRetries)
		time.Sleep(retryInterval)
	}
	log.Printf("WARNING: Setup for %s did not complete within timeout, proceeding anyway", contractName)
}

type AccountSet struct {
	accounts        []*Account
	mu              sync.Mutex
	roundRobinIndex int
}

func NewAccountSet(accounts []*Account) *AccountSet {
	return &AccountSet{
		accounts:        accounts,
		mu:              sync.Mutex{},
		roundRobinIndex: 0,
	}
}

func (a *AccountSet) Len() int {
	return len(a.accounts)
}

func (a *AccountSet) Add(acc *Account) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.accounts = append(a.accounts, acc)
}

func (a *AccountSet) GetAccountRandomly() *Account {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.Len() == 0 {
		return nil
	}
	return a.accounts[rand.Int()%a.Len()]
}

func (a *AccountSet) GetAccountIndex(index int) *Account {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.accounts[index]
}

func (a *AccountSet) GetAccountRoundRobin() *Account {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.Len() == 0 {
		return nil
	}
	acc := a.accounts[a.roundRobinIndex]
	a.roundRobinIndex++
	if a.roundRobinIndex >= a.Len() {
		a.roundRobinIndex = 0
	}
	return acc
}
