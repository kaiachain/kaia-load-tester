package testcase

import (
	"math/big"
	"math/rand"

	"github.com/kaiachain/kaia-load-tester/klayslave/account"
	"github.com/kaiachain/kaia/client"
	"github.com/myzhan/boomer"
)

// SmartContractTxFunc represents a smart contract transaction function signature
type SmartContractTxFunc = func(*client.Client, *account.Account, *account.Account) (interface{}, *big.Int, error)

// RunBaseWithContract creates a closure that executes a test case with contract account
func RunBaseWithContract(config *TCConfig, txFunc SmartContractTxFunc) func() {
	return func() {
		cli := config.CliPool.Alloc().(*client.Client)
		defer config.CliPool.Free(cli)

		from := config.AccGrp.GetAccountRandomly()
		to := config.SmartContractAccounts[config.TestContracts[0]]

		start := boomer.Now()
		_, _, err := txFunc(cli, from, to)
		elapsed := boomer.Now() - start

		if err == nil {
			boomer.Events.Publish("request_success", "http", config.Name+" to "+config.EndPoint, elapsed, int64(10))
		} else {
			boomer.Events.Publish("request_failure", "http", config.Name+" to "+config.EndPoint, elapsed, err.Error())
		}
	}
}

func RunNewSmartContractExecutionTC(config *TCConfig) func() {
	txFunc := func(cli *client.Client, from *account.Account, to *account.Account) (interface{}, *big.Int, error) {
		data := account.TestContractInfos[account.ContractGeneral].GenData(from.GetAddress(), nil)
		return from.TransferNewSmartContractExecutionTx(cli, to, nil, data)
	}
	return RunBaseWithContract(config, txFunc)
}

func RunNewFeeDelegatedSmartContractExecutionTC(config *TCConfig) func() {
	txFunc := func(cli *client.Client, from *account.Account, to *account.Account) (interface{}, *big.Int, error) {
		return from.TransferNewFeeDelegatedSmartContractExecutionTx(cli, to, big.NewInt(0))
	}
	return RunBaseWithContract(config, txFunc)
}

func RunNewFeeDelegatedSmartContractExecutionWithRatioTC(config *TCConfig) func() {
	txFunc := func(cli *client.Client, from *account.Account, to *account.Account) (interface{}, *big.Int, error) {
		return from.TransferNewFeeDelegatedSmartContractExecutionWithRatioTx(cli, to, big.NewInt(0))
	}
	return RunBaseWithContract(config, txFunc)
}

func RunCpuHeavyTC(config *TCConfig) func() {
	txFunc := func(cli *client.Client, from *account.Account, to *account.Account) (interface{}, *big.Int, error) {
		cpuHeavyValue := big.NewInt(100)
		cpuHeavyData := account.TestContractInfos[account.ContractCPUHeavy].GenData(from.GetAddress(), cpuHeavyValue)
		return from.TransferNewSmartContractExecutionTx(cli, to, big.NewInt(0), cpuHeavyData)
	}
	return RunBaseWithContract(config, txFunc)
}

func RunLargeMemoTC(config *TCConfig) func() {
	txFunc := func(cli *client.Client, from *account.Account, to *account.Account) (interface{}, *big.Int, error) {
		// Generate random memo size between 50 and 2000
		memoSize := big.NewInt(int64(50 + rand.Intn(1951)))
		memoData := account.TestContractInfos[account.ContractLargeMemo].GenData(from.GetAddress(), memoSize)
		return from.TransferNewSmartContractExecutionTx(cli, to, big.NewInt(0), memoData)
	}
	return RunBaseWithContract(config, txFunc)
}

func RunErc20TransferTC(config *TCConfig) func() {
	txFunc := func(cli *client.Client, from *account.Account, to *account.Account) (interface{}, *big.Int, error) {
		erc20Value := big.NewInt(int64(rand.Int() % 3))
		receiver := config.AccGrp.GetAccountRandomly()
		erc20Data := account.TestContractInfos[account.ContractErc20].GenData(receiver.GetAddress(), erc20Value)
		return from.TransferNewSmartContractExecutionTx(cli, to, nil, erc20Data)
	}
	return RunBaseWithContract(config, txFunc)
}

func RunErc20TransferWithBlockedListTC(config *TCConfig) func() {
	txFunc := func(cli *client.Client, from *account.Account, to *account.Account) (interface{}, *big.Int, error) {
		tetherValue := big.NewInt(int64(rand.Int() % 3))
		receiver := config.AccGrp.GetAccountRandomly()
		tetherData := account.TestContractInfos[account.ContractTetherProxy].GenData(receiver.GetAddress(), tetherValue)
		return from.TransferNewSmartContractExecutionTx(cli, to, nil, tetherData)
	}
	return RunBaseWithContract(config, txFunc)
}

func RunErc721TransferTC(config *TCConfig) func() {
	return func() {
		cli := config.CliPool.Alloc().(*client.Client)
		defer config.CliPool.Free(cli)

		// Find an account with available tokens
		var (
			fromAcc *account.Account
			tokenId *big.Int
			toAcc   = config.AccGrp.GetAccountRandomly()
		)

		// Try multiple accounts to find one with tokens
		candidateIdx := rand.Intn(config.AccGrp.Len())

		// limit the number of attempts to find a token
		for i := 0; i < config.AccGrp.Len(); i++ {
			candidateAcc := config.AccGrp.GetAccountIndex(candidateIdx)
			tokenId = account.ERC721Ledger.RemoveToken(candidateAcc.GetAddress())
			if tokenId != nil {
				fromAcc = candidateAcc
				break
			}
			candidateIdx = (candidateIdx + 1) % config.AccGrp.Len()
		}

		if tokenId == nil {
			// No tokens available in any account
			boomer.Events.Publish("request_failure", "http", config.Name+" to "+config.EndPoint, int64(0), "No tokens available")
			return
		}

		start := boomer.Now()

		data := account.PackContractCall(account.TestContractInfos[account.ContractErc721].Abi, "transferFrom", fromAcc.GetAddress(), toAcc.GetAddress(), tokenId)
		_, _, err := fromAcc.TransferNewSmartContractExecutionTx(cli, config.SmartContractAccounts[account.ContractErc721], nil, data)

		elapsed := boomer.Now() - start

		if err == nil {
			boomer.Events.Publish("request_success", "http", config.Name+" to "+config.EndPoint, elapsed, int64(10))
			// Transfer successful, add token to destination account
			account.ERC721Ledger.PutToken(toAcc.GetAddress(), tokenId)
		} else {
			boomer.Events.Publish("request_failure", "http", config.Name+" to "+config.EndPoint, elapsed, err.Error())
			// Transfer failed, put token back to original owner
			account.ERC721Ledger.PutToken(fromAcc.GetAddress(), tokenId)
		}
	}
}

// RunGaslessTransactionTC creates a closure for gasless transaction test case
func RunGaslessTransactionTC(config *TCConfig) func() {
	txFunc := func(cli *client.Client, from *account.Account, to *account.Account) (interface{}, *big.Int, error) {
		testTokenAccount := config.SmartContractAccounts[account.ContractGaslessToken]
		gsrAccount := config.SmartContractAccounts[account.ContractGaslessSwapRouter]
		_, _, _, err := from.TransferNewGaslessTx(cli, testTokenAccount, gsrAccount)
		return nil, nil, err
	}
	return RunBaseWithContract(config, txFunc)
}

// RunGaslessRevertTransactionTC creates a closure for gasless revert transaction test case
func RunGaslessRevertTransactionTC(config *TCConfig) func() {
	txFunc := func(cli *client.Client, from *account.Account, to *account.Account) (interface{}, *big.Int, error) {
		testTokenAccount := config.SmartContractAccounts[account.ContractGaslessToken]
		gsrAccount := config.SmartContractAccounts[account.ContractGaslessSwapRouter]
		_, _, _, err := from.TransferNewGaslessTx(cli, testTokenAccount, gsrAccount)
		return nil, nil, err
	}
	return RunBaseWithContract(config, txFunc)
}

// RunGaslessOnlyApproveTC creates a closure for gasless only approve test case
func RunGaslessOnlyApproveTC(config *TCConfig) func() {
	txFunc := func(cli *client.Client, from *account.Account, to *account.Account) (interface{}, *big.Int, error) {
		testTokenAccount := config.SmartContractAccounts[account.ContractGaslessToken]
		gsrAccount := config.SmartContractAccounts[account.ContractGaslessSwapRouter]
		_, _, err := from.TransferNewGaslessApproveTx(cli, testTokenAccount, gsrAccount)
		return nil, nil, err
	}
	return RunBaseWithContract(config, txFunc)
}

// RunInternalTxTC creates a closure for internal transaction test case
func RunInternalTxTC(config *TCConfig) func() {
	txFunc := func(cli *client.Client, from *account.Account, to *account.Account) (interface{}, *big.Int, error) {
		// Get main contract account
		mainContractAccount := config.SmartContractAccounts[account.ContractInternalTxMain]

		// Generate data for sendRewards function call
		data := account.TestContractInfos[account.ContractInternalTxMain].GenData(from.GetAddress(), big.NewInt(100))

		return from.TransferNewSmartContractExecutionTx(cli, mainContractAccount, big.NewInt(100), data)
	}
	return RunBaseWithContract(config, txFunc)
}

// RunMintNFTTC creates a closure for mint NFT test case
func RunMintNFTTC(config *TCConfig) func() {
	txFunc := func(cli *client.Client, from *account.Account, to *account.Account) (interface{}, *big.Int, error) {
		// Get KIP17 contract account
		kip17ContractAccount := config.SmartContractAccounts[account.ContractInternalTxKIP17]

		// Generate data for mintCard function call
		data := account.TestContractInfos[account.ContractInternalTxKIP17].GenData(from.GetAddress(), nil)

		return from.TransferNewSmartContractExecutionTx(cli, kip17ContractAccount, nil, data)
	}
	return RunBaseWithContract(config, txFunc)
}

// RunStorageTrieWriteTC creates a closure for storage trie write test case
func RunStorageTrieWriteTC(config *TCConfig) func() {
	txFunc := func(cli *client.Client, from *account.Account, to *account.Account) (interface{}, *big.Int, error) {
		// Get storage trie contract account
		storageTrieContractAccount := config.SmartContractAccounts[account.ContractStorageTrie]

		// Generate random value between 0 and 2
		value := big.NewInt(int64(rand.Int() % 3))

		return from.ExecuteStorageTrieStore(cli, storageTrieContractAccount, value)
	}
	return RunBaseWithContract(config, txFunc)
}

// RunUserStorageSetTC creates a closure for user storage set test case
func RunUserStorageSetTC(config *TCConfig) func() {
	txFunc := func(cli *client.Client, from *account.Account, to *account.Account) (interface{}, *big.Int, error) {
		value := big.NewInt(1)
		data := account.TestContractInfos[account.ContractUserStorage].GenData(from.GetAddress(), value)

		return from.TransferNewSmartContractExecutionTx(cli, to, nil, data)
	}
	return RunBaseWithContract(config, txFunc)
}

// RunUserStorageSetGetTC creates a closure for user storage set and get test case
func RunUserStorageSetGetTC(config *TCConfig) func() {
	return func() {
		cli := config.CliPool.Alloc().(*client.Client)
		defer config.CliPool.Free(cli)
		from := config.AccGrp.GetAccountRandomly()

		start := boomer.Now()

		// Get user storage contract account
		userStorageContractAccount := config.SmartContractAccounts[account.ContractUserStorage]

		// First, call set function
		setValue := big.NewInt(1)
		setData := account.TestContractInfos[account.ContractUserStorage].GenData(from.GetAddress(), setValue)
		_, _, setErr := from.TransferNewSmartContractExecutionTx(cli, userStorageContractAccount, nil, setData)

		if setErr != nil {
			elapsed := boomer.Now() - start
			boomer.Events.Publish("request_failure", "http", "userStorageSetGet to "+config.EndPoint, elapsed, setErr.Error())
			return
		}

		// Then, call get function
		getValue := big.NewInt(0)
		getData := account.TestContractInfos[account.ContractUserStorage].GenData(from.GetAddress(), getValue)
		_, _, getErr := from.TransferNewSmartContractExecutionTx(cli, userStorageContractAccount, nil, getData)

		elapsed := boomer.Now() - start
		if getErr == nil {
			boomer.Events.Publish("request_success", "http", "userStorageSetGet to "+config.EndPoint, elapsed, int64(10))
		} else {
			boomer.Events.Publish("request_failure", "http", "userStorageSetGet to "+config.EndPoint, elapsed, getErr.Error())
		}
	}
}
