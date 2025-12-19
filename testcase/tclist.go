package testcase

import (
	"github.com/kaiachain/kaia-load-tester/klayslave/account"
)

// Test case name constants
const (
	NewValueTransferTCName                               = "newValueTransferTC"
	NewValueTransferMemoTCName                           = "newValueTransferMemoTC"
	NewSmartContractExecutionTCName                      = "newSmartContractExecutionTC"
	Erc20TransferTCName                                  = "erc20TransferTC"
	Erc20TransferWithBlockedListTCName                   = "erc20TransferWithBlockedListTC"
	CpuHeavyTCName                                       = "cpuHeavyTC"
	NewFeeDelegatedValueTransferTCName                   = "newFeeDelegatedValueTransferTC"
	NewFeeDelegatedValueTransferWithRatioTCName          = "newFeeDelegatedValueTransferWithRatioTC"
	NewFeeDelegatedValueTransferMemoTCName               = "newFeeDelegatedValueTransferMemoTC"
	NewFeeDelegatedValueTransferMemoWithRatioTCName      = "newFeeDelegatedValueTransferMemoWithRatioTC"
	NewFeeDelegatedSmartContractDeployTCName             = "newFeeDelegatedSmartContractDeployTC"
	NewFeeDelegatedSmartContractDeployWithRatioTCName    = "newFeeDelegatedSmartContractDeployWithRatioTC"
	NewFeeDelegatedSmartContractExecutionTCName          = "newFeeDelegatedSmartContractExecutionTC"
	NewFeeDelegatedSmartContractExecutionWithRatioTCName = "newFeeDelegatedSmartContractExecutionWithRatioTC"
	NewValueTransferWithCancelTCName                     = "newValueTransferWithCancelTC"
	NewValueTransferLargeMemoTCName                      = "newValueTransferLargeMemoTC"
	NewValueTransferSmallMemoTCName                      = "newValueTransferSmallMemoTC"
	NewCancelTCName                                      = "newCancelTC"
	NewFeeDelegatedCancelTCName                          = "newFeeDelegatedCancelTC"
	NewFeeDelegatedCancelWithRatioTCName                 = "newFeeDelegatedCancelWithRatioTC"
	NewSmartContractDeployTCName                         = "newSmartContractDeployTC"
	LargeMemoTCName                                      = "largeMemoTC"
	Erc721TransferTCName                                 = "erc721TransferTC"
	AuctionBidTCName                                     = "auctionBidTC"
	AuctionRevertedBidTCName                             = "auctionRevertedBidTC"
	GaslessTransactionTCName                             = "gaslessTransactionTC"
	GaslessRevertTransactionTCName                       = "gaslessRevertTransactionTC"
	GaslessOnlyApproveTCName                             = "gaslessOnlyApproveTC"
	ReadGasPriceTCName                                   = "readGasPrice"
	ReadBlockNumberTCName                                = "readBlockNumber"
	ReadGetBlockByNumberTCName                           = "readGetBlockByNumber"
	ReadGetAccountTCName                                 = "readGetAccount"
	ReadGetBlockWithConsensusInfoByNumberTCName          = "readGetBlockWithConsensusInfoByNumber"
	ReadGetStorageAtTCName                               = "readGetStorageAt"
	ReadCallTCName                                       = "readCall"
	ReadEstimateGasTCName                                = "readEstimateGas"
	InternalTxTCName                                     = "internalTxTC"
	MintNFTTCName                                        = "mintNFTTC"
	StorageTrieWriteTCName                               = "storageTrieWriteTC"
	UserStorageSetTCName                                 = "userStorageSetTC"
	UserStorageSetGetTCName                              = "userStorageSetGetTC"
	NewAccountUpdateTCName                               = "newAccountUpdateTC"
	NewFeeDelegatedAccountUpdateTCName                   = "newFeeDelegatedAccountUpdateTC"
	NewFeeDelegatedAccountUpdateWithRatioTCName          = "newFeeDelegatedAccountUpdateWithRatioTC"
	TransferSignedTCName                                 = "transferSignedTx"
	TransferUnsignedTCName                               = "transferUnsignedTx"
	ReceiptCheckTCName                                   = "receiptCheckTx"
	TransferSignedWithCheckTCName                        = "transferSignedWithCheckTx"
	EthereumTxLegacyTCName                               = "ethereumTxLegacyTC"
	EthereumTxAccessListTCName                           = "ethereumTxAccessListTC"
	EthereumTxDynamicFeeTCName                           = "ethereumTxDynamicFeeTC"
	NewEthereumAccessListTCName                          = "newEthereumAccessListTC"
	NewEthereumDynamicFeeTCName                          = "newEthereumDynamicFeeTC"
)

// ExtendedTask represents a test case
type ExtendedTask struct {
	Name          string
	Weight        int
	Init          func(accGrp *account.AccGroup, endpoint string, testContracts []account.TestContract, tcName string, targetTxTypeList []string) *TCConfig
	Run           func(config *TCConfig) func()
	TestContracts []account.TestContract // Required test contracts for this task
}

// TcList contains test cases
var TcList = map[string]*ExtendedTask{
	NewValueTransferTCName: {
		Name:          NewValueTransferTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunNewValueTransferTC,
		TestContracts: []account.TestContract{}, // No specific contract needed
	},
	NewValueTransferMemoTCName: {
		Name:          NewValueTransferMemoTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunNewValueTransferMemoTC,
		TestContracts: []account.TestContract{}, // No specific contract needed
	},
	NewSmartContractExecutionTCName: {
		Name:          NewSmartContractExecutionTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunNewSmartContractExecutionTC,
		TestContracts: []account.TestContract{account.ContractGeneral},
	},
	Erc20TransferTCName: {
		Name:          Erc20TransferTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunErc20TransferTC,
		TestContracts: []account.TestContract{account.ContractErc20},
	},
	Erc20TransferWithBlockedListTCName: {
		Name:          Erc20TransferWithBlockedListTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunErc20TransferWithBlockedListTC,
		TestContracts: []account.TestContract{account.ContractTetherProxy, account.ContractTetherLogic},
	},
	CpuHeavyTCName: {
		Name:          CpuHeavyTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunCpuHeavyTC,
		TestContracts: []account.TestContract{account.ContractCPUHeavy},
	},
	NewFeeDelegatedValueTransferTCName: {
		Name:          NewFeeDelegatedValueTransferTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunNewFeeDelegatedValueTransferTC,
		TestContracts: []account.TestContract{},
	},
	NewFeeDelegatedValueTransferWithRatioTCName: {
		Name:          NewFeeDelegatedValueTransferWithRatioTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunNewFeeDelegatedValueTransferWithRatioTC,
		TestContracts: []account.TestContract{},
	},
	NewFeeDelegatedValueTransferMemoTCName: {
		Name:          NewFeeDelegatedValueTransferMemoTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunNewFeeDelegatedValueTransferMemoTC,
		TestContracts: []account.TestContract{},
	},
	NewFeeDelegatedValueTransferMemoWithRatioTCName: {
		Name:          NewFeeDelegatedValueTransferMemoWithRatioTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunNewFeeDelegatedValueTransferMemoWithRatioTC,
		TestContracts: []account.TestContract{},
	},
	NewFeeDelegatedSmartContractDeployTCName: {
		Name:          NewFeeDelegatedSmartContractDeployTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunNewFeeDelegatedSmartContractDeployTC,
		TestContracts: []account.TestContract{},
	},
	NewFeeDelegatedSmartContractDeployWithRatioTCName: {
		Name:          NewFeeDelegatedSmartContractDeployWithRatioTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunNewFeeDelegatedSmartContractDeployWithRatioTC,
		TestContracts: []account.TestContract{},
	},
	NewFeeDelegatedSmartContractExecutionTCName: {
		Name:          NewFeeDelegatedSmartContractExecutionTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunNewFeeDelegatedSmartContractExecutionTC,
		TestContracts: []account.TestContract{account.ContractGeneral},
	},
	NewFeeDelegatedSmartContractExecutionWithRatioTCName: {
		Name:          NewFeeDelegatedSmartContractExecutionWithRatioTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunNewFeeDelegatedSmartContractExecutionWithRatioTC,
		TestContracts: []account.TestContract{account.ContractGeneral},
	},
	NewValueTransferWithCancelTCName: {
		Name:          NewValueTransferWithCancelTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunNewValueTransferWithCancelTC,
		TestContracts: []account.TestContract{},
	},
	NewValueTransferLargeMemoTCName: {
		Name:          NewValueTransferLargeMemoTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunNewValueTransferLargeMemoTC,
		TestContracts: []account.TestContract{},
	},
	NewValueTransferSmallMemoTCName: {
		Name:          NewValueTransferSmallMemoTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunNewValueTransferSmallMemoTC,
		TestContracts: []account.TestContract{},
	},
	NewCancelTCName: {
		Name:          NewCancelTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunNewCancelTC,
		TestContracts: []account.TestContract{},
	},
	NewFeeDelegatedCancelTCName: {
		Name:          NewFeeDelegatedCancelTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunNewFeeDelegatedCancelTC,
		TestContracts: []account.TestContract{},
	},
	NewFeeDelegatedCancelWithRatioTCName: {
		Name:          NewFeeDelegatedCancelWithRatioTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunNewFeeDelegatedCancelWithRatioTC,
		TestContracts: []account.TestContract{},
	},
	NewSmartContractDeployTCName: {
		Name:          NewSmartContractDeployTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunNewSmartContractDeployTC,
		TestContracts: []account.TestContract{},
	},
	LargeMemoTCName: {
		Name:          LargeMemoTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunLargeMemoTC,
		TestContracts: []account.TestContract{account.ContractLargeMemo},
	},
	Erc721TransferTCName: {
		Name:          Erc721TransferTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunErc721TransferTC,
		TestContracts: []account.TestContract{account.ContractErc721},
	},
	AuctionBidTCName: {
		Name:          AuctionBidTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunAuctionBidTC,
		TestContracts: []account.TestContract{account.ContractAuctionEntryPoint, account.ContractCounterForTestAuction},
	},
	AuctionRevertedBidTCName: {
		Name:          AuctionRevertedBidTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunAuctionRevertedBidTC,
		TestContracts: []account.TestContract{account.ContractAuctionEntryPoint, account.ContractCounterForTestAuction},
	},
	GaslessTransactionTCName: {
		Name:          GaslessTransactionTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunGaslessTransactionTC,
		TestContracts: []account.TestContract{account.ContractGaslessToken, account.ContractGaslessSwapRouter},
	},
	GaslessRevertTransactionTCName: {
		Name:          GaslessRevertTransactionTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunGaslessRevertTransactionTC,
		TestContracts: []account.TestContract{account.ContractGaslessToken, account.ContractGaslessSwapRouter},
	},
	GaslessOnlyApproveTCName: {
		Name:          GaslessOnlyApproveTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunGaslessOnlyApproveTC,
		TestContracts: []account.TestContract{account.ContractGaslessToken, account.ContractGaslessSwapRouter},
	},
	ReadGasPriceTCName: {
		Name:          ReadGasPriceTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunGasPrice,
		TestContracts: []account.TestContract{},
	},
	ReadBlockNumberTCName: {
		Name:          ReadBlockNumberTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunBlockNumber,
		TestContracts: []account.TestContract{},
	},
	ReadGetBlockByNumberTCName: {
		Name:          ReadGetBlockByNumberTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunGetBlockByNumber,
		TestContracts: []account.TestContract{},
	},
	ReadGetAccountTCName: {
		Name:          ReadGetAccountTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunGetAccount,
		TestContracts: []account.TestContract{},
	},
	ReadGetBlockWithConsensusInfoByNumberTCName: {
		Name:          ReadGetBlockWithConsensusInfoByNumberTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunGetBlockWithConsensusInfoByNumber,
		TestContracts: []account.TestContract{},
	},
	ReadGetStorageAtTCName: {
		Name:          ReadGetStorageAtTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunGetStorageAt,
		TestContracts: []account.TestContract{account.ContractReadApiCallContract},
	},
	ReadCallTCName: {
		Name:          ReadCallTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunCall,
		TestContracts: []account.TestContract{account.ContractReadApiCallContract},
	},
	ReadEstimateGasTCName: {
		Name:          ReadEstimateGasTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunEstimateGas,
		TestContracts: []account.TestContract{account.ContractReadApiCallContract},
	},
	InternalTxTCName: {
		Name:          InternalTxTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunInternalTxTC,
		TestContracts: []account.TestContract{account.ContractInternalTxMain},
	},
	MintNFTTCName: {
		Name:          MintNFTTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunMintNFTTC,
		TestContracts: []account.TestContract{account.ContractInternalTxKIP17},
	},
	StorageTrieWriteTCName: {
		Name:          StorageTrieWriteTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunStorageTrieWriteTC,
		TestContracts: []account.TestContract{account.ContractStorageTrie},
	},
	UserStorageSetTCName: {
		Name:          UserStorageSetTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunUserStorageSetTC,
		TestContracts: []account.TestContract{account.ContractUserStorage},
	},
	UserStorageSetGetTCName: {
		Name:          UserStorageSetGetTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunUserStorageSetGetTC,
		TestContracts: []account.TestContract{account.ContractUserStorage},
	},
	NewAccountUpdateTCName: {
		Name:          NewAccountUpdateTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunNewAccountUpdateTC,
		TestContracts: []account.TestContract{},
	},
	NewFeeDelegatedAccountUpdateTCName: {
		Name:          NewFeeDelegatedAccountUpdateTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunNewFeeDelegatedAccountUpdateTC,
		TestContracts: []account.TestContract{},
	},
	NewFeeDelegatedAccountUpdateWithRatioTCName: {
		Name:          NewFeeDelegatedAccountUpdateWithRatioTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunNewFeeDelegatedAccountUpdateWithRatioTC,
		TestContracts: []account.TestContract{},
	},
	TransferSignedTCName: {
		Name:          TransferSignedTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunTransferSignedTC,
		TestContracts: []account.TestContract{},
	},
	TransferUnsignedTCName: {
		Name:          TransferUnsignedTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunTransferUnsignedTC,
		TestContracts: []account.TestContract{},
	},
	ReceiptCheckTCName: {
		Name:          ReceiptCheckTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunReceiptCheckTC,
		TestContracts: []account.TestContract{},
	},
	TransferSignedWithCheckTCName: {
		Name:          TransferSignedWithCheckTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunTransferSignedWithCheckTC,
		TestContracts: []account.TestContract{},
	},
	EthereumTxLegacyTCName: {
		Name:          EthereumTxLegacyTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunEthereumTxLegacyTC,
		TestContracts: []account.TestContract{account.ContractGeneral},
	},
	EthereumTxAccessListTCName: {
		Name:          EthereumTxAccessListTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunEthereumTxAccessListTC,
		TestContracts: []account.TestContract{account.ContractGeneral},
	},
	EthereumTxDynamicFeeTCName: {
		Name:          EthereumTxDynamicFeeTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunEthereumTxDynamicFeeTC,
		TestContracts: []account.TestContract{account.ContractGeneral},
	},
	NewEthereumAccessListTCName: {
		Name:          NewEthereumAccessListTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunNewEthereumAccessListTC,
		TestContracts: []account.TestContract{account.ContractGeneral},
	},
	NewEthereumDynamicFeeTCName: {
		Name:          NewEthereumDynamicFeeTCName,
		Weight:        10,
		Init:          Init,
		Run:           RunNewEthereumDynamicFeeTC,
		TestContracts: []account.TestContract{account.ContractGeneral},
	},
}
