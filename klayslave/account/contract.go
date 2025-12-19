package account

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"strings"
	"time"

	// TODO: Change the import source to Kaia repo ---------
	auctionDepositVaultContracts "github.com/kaiachain/kaia-load-tester/klayslave/account/contracts/auctionDepositVault"
	auctionEntryPointContracts "github.com/kaiachain/kaia-load-tester/klayslave/account/contracts/auctionEntryPoint"
	auctionFeeVaultContracts "github.com/kaiachain/kaia-load-tester/klayslave/account/contracts/auctionFeeVault"
	tetherContracts "github.com/kaiachain/kaia-load-tester/klayslave/account/contracts/tetherContract"

	// ----------------

	"github.com/kaiachain/kaia/accounts/abi"
	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/client"
	"github.com/kaiachain/kaia/common"
	uniswapFactoryContracts "github.com/kaiachain/kaia/contracts/contracts/libs/uniswap/factory"
	uniswapRouterContracts "github.com/kaiachain/kaia/contracts/contracts/libs/uniswap/router"
	kip149contract "github.com/kaiachain/kaia/contracts/contracts/system_contracts/kip149"
	gaslessContract "github.com/kaiachain/kaia/contracts/contracts/system_contracts/kip247"
	testingContracts "github.com/kaiachain/kaia/contracts/contracts/testing/system_contracts"
	testingGaslessContracts "github.com/kaiachain/kaia/contracts/contracts/testing/system_contracts/gasless"
	"github.com/kaiachain/kaia/crypto"
	gaslessImpl "github.com/kaiachain/kaia/kaiax/gasless/impl"
	"github.com/kaiachain/kaia/params"
)

// ABI constants for contract calls
const (
	erc20ABI                     = `[{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"value","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"sender","type":"address"},{"name":"recipient","type":"address"},{"name":"amount","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"addedValue","type":"uint256"}],"name":"increaseAllowance","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"account","type":"address"},{"name":"amount","type":"uint256"}],"name":"mint","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"account","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"account","type":"address"}],"name":"addMinter","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"renounceMinter","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"subtractedValue","type":"uint256"}],"name":"decreaseAllowance","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"recipient","type":"address"},{"name":"amount","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"account","type":"address"}],"name":"isMinter","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"owner","type":"address"},{"name":"spender","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"name":"account","type":"address"}],"name":"MinterAdded","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"account","type":"address"}],"name":"MinterRemoved","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":true,"name":"spender","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Approval","type":"event"}]`
	erc721PerformanceABI         = `[{"constant":true,"inputs":[{"name":"interfaceId","type":"bytes4"}],"name":"supportsInterface","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"tokenId","type":"uint256"}],"name":"getApproved","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"tokenId","type":"uint256"}],"name":"approve","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"from","type":"address"},{"name":"to","type":"address"},{"name":"tokenId","type":"uint256"}],"name":"transferFrom","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"from","type":"address"},{"name":"to","type":"address"},{"name":"tokenId","type":"uint256"}],"name":"safeTransferFrom","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"tokenId","type":"uint256"},{"name":"tokenURI","type":"string"}],"name":"mintWithTokenURI","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"tokenId","type":"uint256"}],"name":"ownerOf","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_user","type":"address"},{"name":"_startID","type":"uint256"},{"name":"_endID","type":"uint256"}],"name":"registerBulk","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"approved","type":"bool"}],"name":"setApprovalForAll","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"from","type":"address"},{"name":"to","type":"address"},{"name":"tokenId","type":"uint256"},{"name":"_data","type":"bytes"}],"name":"safeTransferFrom","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"tokenId","type":"uint256"}],"name":"tokenURI","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"owner","type":"address"},{"name":"operator","type":"address"}],"name":"isApprovedForAll","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[{"name":"name","type":"string"},{"name":"symbol","type":"string"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":true,"name":"tokenId","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":true,"name":"approved","type":"address"},{"indexed":true,"name":"tokenId","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":true,"name":"operator","type":"address"},{"indexed":false,"name":"approved","type":"bool"}],"name":"ApprovalForAll","type":"event"}]`
	generalPurposeABI            = `[{"constant":true,"inputs":[],"name":"totalAmount","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"receiver","type":"address"}],"name":"reward","outputs":[],"payable":true,"stateMutability":"payable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"safeWithdrawal","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"inputs":[],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"payable":true,"stateMutability":"payable","type":"fallback"}]`
	gaslessApproveABI            = `[{"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"}]`
	gaslessSwapABI               = `[{"inputs":[{"internalType":"address","name":"token","type":"address"},{"internalType":"uint256","name":"amountIn","type":"uint256"},{"internalType":"uint256","name":"minAmountOut","type":"uint256"},{"internalType":"uint256","name":"amountRepay","type":"uint256"},{"internalType":"uint256","name":"deadline","type":"uint256"}],"name":"swapForGas","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
	counterAuctionABI            = `[{"inputs":[],"name":"getForAuction","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getForSC","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"incForAuction","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"incForSC","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"intendedRevert","outputs":[],"stateMutability":"pure","type":"function"}]`
	cpuHeavySortABI              = `[{"inputs":[{"internalType":"uint256","name":"size","type":"uint256"},{"internalType":"uint256","name":"signature","type":"uint256"}],"name":"sort","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
	largeMemoABI                 = `[{"inputs":[{"internalType":"string","name":"_str","type":"string"}],"name":"setName","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
	readApiGetABI                = `[{"inputs":[],"name":"get","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`
	readApiSetABI                = `[{"inputs":[],"name":"set","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
	userStorageGetABI            = `[{"inputs":[],"name":"get","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`
	userStorageSetABI            = `[{"inputs":[{"internalType":"uint256","name":"x","type":"uint256"}],"name":"set","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
	userStorageGetUserDataABI    = `[{"inputs":[{"internalType":"address","name":"user","type":"address"}],"name":"getUserData","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`
	mintCardABI                  = `[{"inputs":[],"name":"mintCard","outputs":[],"stateMutability":"payable","type":"function"}]`
	sendRewardsABI               = `[{"inputs":[{"internalType":"address payable","name":"invitee","type":"address"},{"internalType":"address payable","name":"host","type":"address"}],"name":"sendRewards","outputs":[],"stateMutability":"payable","type":"function"}]`
	internalTxMainConstructorABI = `[{"inputs":[{"internalType":"address","name":"cardContractAddress","type":"address"},{"internalType":"uint256","name":"rewardForHost","type":"uint256"},{"internalType":"uint256","name":"rewardForInvitee","type":"uint256"}],"stateMutability":"nonpayable","type":"constructor"}]`
)

// PackContractCall packs ABI-encoded call data for a contract method
func PackContractCall(abiStr string, method string, args ...interface{}) []byte {
	abii, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		log.Fatalf("failed to abi.JSON: %v", err)
	}
	data, err := abii.Pack(method, args...)
	if err != nil {
		log.Fatalf("failed to abi.Pack: %v", err)
	}
	return data
}

// Contract deployer accounts
var (
	ERC20Deployer                 = GetAccountFromKey(0, "eb2c84d41c639178ff26a81f488c196584d678bb1390cc20a3aeb536f3969a98")
	ERC721Deployer                = GetAccountFromKey(0, "45c40d95c9b7898a21e073b5bf952bcb05f2e70072e239a8bbd87bb74a53355e")
	StorageTrieDeployer           = GetAccountFromKey(0, "3737c381633deaaa4c0bdbc64728f6ef7d381b17e1d30bbb74665839cec942b8")
	GeneralPurposeDeployer        = GetAccountFromKey(0, "c0cd1721f60535cb7779e5db43a94390aff9ead01ee3d654abffcb0453bdc927")
	GaslessTokenDeployer          = GetAccountFromKey(0, "e095e5fdfc55ce9002edc26fdf402b8ece64586e9673f09b0a91dde39ccc8abe")
	WKaiaDeployer                 = GetAccountFromKey(0, "56f48de8c67737661df6b66d968e2597754051e9967dc665d901bc2e7aa2ee39")
	UniswapFactoryDeployer        = GetAccountFromKey(0, "c6f61a31be1ca48b7774568bd47eacd03aed8fa9265d9eeb64a97136ea8e411a")
	UniswapRouterDeployer         = GetAccountFromKey(0, "780d71b4ee7121673cf28492a3185bf97fcd9fe280c72d6df77d99648ba74541")
	GaslessSwapRouterDeployer     = GetAccountFromKey(0, "5a212da24b990b2164a2cbe070d15e8f2948b636cb224f9f72979faa564ef42f")
	CounterForTestAuctionDeployer = GetAccountFromKey(0, "caa1c841d600a7a37f08ca38478480a3ee254cefc30324048efd16a4f473d26d")
	AuctionFeeVaultDeployer       = GetAccountFromKey(0, "34e4baf3acf5fe6eeda59ffbe7e7c525c835aa58c671889cab95c3210b27f2cf")
	AuctionDepositVaultDeployer   = GetAccountFromKey(0, "ae2a792e63ffe80c098d70f9d486e775790b6abb74db54dc5f0894965c3d9578")
	AuctionEntryPointDeployer     = GetAccountFromKey(0, "2702a7f5f21a17ced7edd46f1d930e5b0d36d36278cfd343973b69a3673bce6c")
	GSRSetupManager               = GetAccountFromKey(0, "76a5b8060388e1f83f7b3bdbcc5248d13b3c7e9771e8445afb2754ad5e192237")
	Auctioneer                    = GetAccountFromKey(0, "b7dce0e6f88e4591bb8dc8c0f4d5082a10e38b4836e06a4a7f9d92cdb4a6f671")
	CPUHeavyDeployer              = GetAccountFromKey(0, "f8d7eccaf0d2bb863daf2301f5e6a29626cadb8746f29565aaabfdf6e0f0c073")
	LargeMemoDeployer             = GetAccountFromKey(0, "b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef12345678")
	ReadApiCallContractDeployer   = GetAccountFromKey(0, "d1e2c3f4a5b6789012345678901234567890abcdef1234567890abcdef123456")
	UserStorageDeployer           = GetAccountFromKey(0, "c3d4e5f6789012345678901234567890abcdef1234567890abcdef1234567890")
	InternalTxKIP17Deployer       = GetAccountFromKey(0, "f5a6b7c890123456789012345678901234567890abcdef1234567890abcdef12")
	InternalTxMainDeployer        = GetAccountFromKey(0, "e4f5a6b7c890123456789012345678901234567890abcdef1234567890abcdef")
	TetherLogicDeployer           = GetAccountFromKey(0, "a1b2c3d4e5f678901234567890abcdef1234567890abcdef1234567890abcd01")
	TetherProxyDeployer           = GetAccountFromKey(0, "a1b2c3d4e5f678901234567890abcdef1234567890abcdef1234567890abcd02")
)

// TestContractInfo represents a test contract configuration
type TestContractInfo struct {
	testNames               []string
	auctionTargetTxTypeList []string
	Bytecode                []byte
	deployer                *Account
	contractName            string
	// TODO: make GenData array or use go wrapper file
	Abi                             string
	GenData                         func(addr common.Address, value *big.Int) []byte
	GetBytecodeWithConstructorParam func(bin []byte, contracts []*Account, deployer *Account) []byte
	IsDeployed                      func(gCli *client.Client, deployer *Account) bool
	GetAddress                      func(gCli *client.Client, deployer *Account) common.Address
	// DoSetupWork is executed by leader only (e.g., GSR registration, Auction registration)
	DoSetupWork func(ctx *AdditionalWorkContext)
	// DoChargingWork is executed by all slaves after setup is complete (e.g., token charging, NFT minting)
	DoChargingWork func(ctx *AdditionalWorkContext)
	// WaitForSetup returns true when setup is complete (used by followers to wait for leader)
	WaitForSetup func(gCli *client.Client) bool
}

// AdditionalWorkContext contains all context needed for additional work after contract deployment
type AdditionalWorkContext struct {
	GCli             *client.Client
	LocalReservoir   *Account
	GlobalReservoir  *Account
	ChargeValue      *big.Int
	IsLeader         bool
	MaxConcurrency   int
	AccGrp           *AccGroup
	TcList           []string
	TargetTxTypeList []string
}

// TestContractInfos stores some dedicated and fixed private key used to deploy a smart contracts for TCs.
// Whether or not to deploy a Gasless-related contract is determined by whether gsr is registered in the registry.
// GenData is a data which
// (1) ERC20 value transfer
// (2) ERC721 value transfer
// (3) storage trie write performance test
// (4) other general contract calls
// (5) TestToken approve
// (6) -
// (7) -
// (8) -
// (9) Gasless Swap Router swap
// TODO-kaia-load-tester: register GenData at TcList extendedTask or find other way to register it.
var TestContractInfos = []TestContractInfo{
	createERC20ContractInfo(),
	createERC721ContractInfo(),
	createStorageTrieContractInfo(),
	createGeneralPurposeContractInfo(),
	createGaslessTokenContractInfo(),
	createWKaiaContractInfo(),
	createUniswapFactoryContractInfo(),
	createUniswapRouterContractInfo(),
	createGaslessSwapRouterContractInfo(),
	createCounterForTestAuctionContractInfo(),
	createAuctionFeeVaultContractInfo(),
	createAuctionDepositVaultContractInfo(),
	createAuctionEntryPointContractInfo(),
	createCPUHeavyContractInfo(),
	createLargeMemoContractInfo(),
	createReadApiCallContractInfo(),
	createUserStorageContractInfo(),
	createInternalTxKIP17ContractInfo(),
	createInternalTxMainContractInfo(),
	createTetherLogicContractInfo(),
	createTetherProxyContractInfo(),
}

func createERC20ContractInfo() TestContractInfo {
	return TestContractInfo{
		testNames:               []string{"erc20TransferTC"},
		auctionTargetTxTypeList: []string{},
		Bytecode:                common.FromHex("60806040523480156200001157600080fd5b506200002c3362000053640100000000026401000000009004565b6200004c3364e8d4a51000620000bd640100000000026401000000009004565b5062000642565b620000778160036200019964010000000002620013cc179091906401000000009004565b8073ffffffffffffffffffffffffffffffffffffffff167f6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f660405160405180910390a250565b6000620000d93362000288640100000000026401000000009004565b151562000174576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260308152602001807f4d696e746572526f6c653a2063616c6c657220646f6573206e6f74206861766581526020017f20746865204d696e74657220726f6c650000000000000000000000000000000081525060400191505060405180910390fd5b6200018f8383620002b5640100000000026401000000009004565b6001905092915050565b620001b4828262000493640100000000026401000000009004565b1515156200022a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601f8152602001807f526f6c65733a206163636f756e7420616c72656164792068617320726f6c650081525060200191505060405180910390fd5b60018260000160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055505050565b6000620002ae8260036200049364010000000002620012a9179091906401000000009004565b9050919050565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16141515156200035b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601f8152602001807f45524332303a206d696e7420746f20746865207a65726f20616464726573730081525060200191505060405180910390fd5b6200038081600254620005b76401000000000262000fae179091906401000000009004565b600281905550620003e7816000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054620005b76401000000000262000fae179091906401000000009004565b6000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040518082815260200191505060405180910390a35050565b60008073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415151562000560576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260228152602001807f526f6c65733a206163636f756e7420697320746865207a65726f20616464726581526020017f737300000000000000000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b8260000160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16905092915050565b600080828401905083811015151562000638576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f536166654d6174683a206164646974696f6e206f766572666c6f77000000000081525060200191505060405180910390fd5b8091505092915050565b6115d780620006526000396000f3006080604052600436106100ba576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063095ea7b3146100bf57806318160ddd1461012457806323b872dd1461014f57806339509351146101d457806340c10f191461023957806370a082311461029e578063983b2d56146102f55780639865027514610338578063a457c2d71461034f578063a9059cbb146103b4578063aa271e1a14610419578063dd62ed3e14610474575b600080fd5b3480156100cb57600080fd5b5061010a600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291905050506104eb565b604051808215151515815260200191505060405180910390f35b34801561013057600080fd5b50610139610502565b6040518082815260200191505060405180910390f35b34801561015b57600080fd5b506101ba600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019092919050505061050c565b604051808215151515815260200191505060405180910390f35b3480156101e057600080fd5b5061021f600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291905050506105bd565b604051808215151515815260200191505060405180910390f35b34801561024557600080fd5b50610284600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610662565b604051808215151515815260200191505060405180910390f35b3480156102aa57600080fd5b506102df600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061071b565b6040518082815260200191505060405180910390f35b34801561030157600080fd5b50610336600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610763565b005b34801561034457600080fd5b5061034d610812565b005b34801561035b57600080fd5b5061039a600480360381019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019092919050505061081d565b604051808215151515815260200191505060405180910390f35b3480156103c057600080fd5b506103ff600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291905050506108c2565b604051808215151515815260200191505060405180910390f35b34801561042557600080fd5b5061045a600480360381019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506108d9565b604051808215151515815260200191505060405180910390f35b34801561048057600080fd5b506104d5600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506108f6565b6040518082815260200191505060405180910390f35b60006104f833848461097d565b6001905092915050565b6000600254905090565b6000610519848484610bfe565b6105b284336105ad85600160008a73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610f2490919063ffffffff16565b61097d565b600190509392505050565b6000610658338461065385600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610fae90919063ffffffff16565b61097d565b6001905092915050565b600061066d336108d9565b1515610707576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260308152602001807f4d696e746572526f6c653a2063616c6c657220646f6573206e6f74206861766581526020017f20746865204d696e74657220726f6c650000000000000000000000000000000081525060400191505060405180910390fd5b6107118383611038565b6001905092915050565b60008060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b61076c336108d9565b1515610806576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260308152602001807f4d696e746572526f6c653a2063616c6c657220646f6573206e6f74206861766581526020017f20746865204d696e74657220726f6c650000000000000000000000000000000081525060400191505060405180910390fd5b61080f816111f5565b50565b61081b3361124f565b565b60006108b833846108b385600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610f2490919063ffffffff16565b61097d565b6001905092915050565b60006108cf338484610bfe565b6001905092915050565b60006108ef8260036112a990919063ffffffff16565b9050919050565b6000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1614151515610a48576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260248152602001807f45524332303a20617070726f76652066726f6d20746865207a65726f2061646481526020017f726573730000000000000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614151515610b13576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260228152602001807f45524332303a20617070726f766520746f20746865207a65726f20616464726581526020017f737300000000000000000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b80600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925836040518082815260200191505060405180910390a3505050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1614151515610cc9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260258152602001807f45524332303a207472616e736665722066726f6d20746865207a65726f20616481526020017f647265737300000000000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614151515610d94576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260238152602001807f45524332303a207472616e7366657220746f20746865207a65726f206164647281526020017f657373000000000000000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b610de5816000808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610f2490919063ffffffff16565b6000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550610e78816000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610fae90919063ffffffff16565b6000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040518082815260200191505060405180910390a3505050565b600080838311151515610f9f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f536166654d6174683a207375627472616374696f6e206f766572666c6f77000081525060200191505060405180910390fd5b82840390508091505092915050565b600080828401905083811015151561102e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f536166654d6174683a206164646974696f6e206f766572666c6f77000000000081525060200191505060405180910390fd5b8091505092915050565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16141515156110dd576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601f8152602001807f45524332303a206d696e7420746f20746865207a65726f20616464726573730081525060200191505060405180910390fd5b6110f281600254610fae90919063ffffffff16565b600281905550611149816000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610fae90919063ffffffff16565b6000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040518082815260200191505060405180910390a35050565b6112098160036113cc90919063ffffffff16565b8073ffffffffffffffffffffffffffffffffffffffff167f6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f660405160405180910390a250565b6112638160036114a990919063ffffffff16565b8073ffffffffffffffffffffffffffffffffffffffff167fe94479a9f7e1952cc78f2d6baab678adc1b772d936c6583def489e524cb6669260405160405180910390a250565b60008073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614151515611375576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260228152602001807f526f6c65733a206163636f756e7420697320746865207a65726f20616464726581526020017f737300000000000000000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b8260000160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16905092915050565b6113d682826112a9565b15151561144b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601f8152602001807f526f6c65733a206163636f756e7420616c72656164792068617320726f6c650081525060200191505060405180910390fd5b60018260000160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055505050565b6114b382826112a9565b151561154d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260218152602001807f526f6c65733a206163636f756e7420646f6573206e6f74206861766520726f6c81526020017f650000000000000000000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b60008260000160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff02191690831515021790555050505600a165627a7a72305820577de674f02c621a82595da1d61a932e3fd2a3286a9a4e9dbf48df7002e9b5010029"),
		deployer:                ERC20Deployer,
		contractName:            "ERC20 Performance Test Contract",
		Abi:                     erc20ABI,
		GenData: func(recipientAddr common.Address, value *big.Int) []byte {
			if value.Cmp(big.NewInt(1e11)) == 0 {
				return PackContractCall(erc20ABI, "mint", recipientAddr, value)
			}
			return PackContractCall(erc20ABI, "transfer", recipientAddr, value)
		},
		GetBytecodeWithConstructorParam: returnBinAsIs,
		IsDeployed:                      isDeployerNonceNotZero,
		GetAddress:                      getNonce0ContractAddress,
		DoChargingWork: func(ctx *AdditionalWorkContext) {
			log.Printf("Start erc20 token charging to the test account group")
			contract := ctx.AccGrp.GetTestContractByName(ContractErc20)
			ERC20Deployer.SmartContractExecutionWithGuaranteeRetry(ctx.GCli, contract, nil, PackContractCall(erc20ABI, "mint", ctx.LocalReservoir.address, big.NewInt(1e11)))
			ConcurrentTransactionSend(ctx.AccGrp.GetValidAccGrp(), ctx.MaxConcurrency, func(_ int, acc *Account) {
				ctx.LocalReservoir.SmartContractExecutionWithGuaranteeRetry(ctx.GCli, contract, nil, PackContractCall(erc20ABI, "transfer", acc.address, big.NewInt(1e4)))
			})
		},
	}
}

func createERC721ContractInfo() TestContractInfo {
	return TestContractInfo{
		testNames:                       []string{"erc721TransferTC"},
		auctionTargetTxTypeList:         []string{},
		Bytecode:                        common.FromHex("60806040523480156200001157600080fd5b506040516200231038038062002310833981018060405260408110156200003757600080fd5b8101908080516401000000008111156200005057600080fd5b828101905060208101848111156200006757600080fd5b81518560018202830111640100000000821117156200008557600080fd5b50509291906020018051640100000000811115620000a257600080fd5b82810190506020810184811115620000b957600080fd5b8151856001820283011164010000000082111715620000d757600080fd5b5050929190505050620000f76301ffc9a760e01b6200016160201b60201c565b6200010f6380ac58cd60e01b6200016160201b60201c565b8160059080519060200190620001279291906200026a565b508060069080519060200190620001409291906200026a565b5062000159635b5e139f60e01b6200016160201b60201c565b505062000319565b63ffffffff60e01b817bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19161415620001fe576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601c8152602001807f4552433136353a20696e76616c696420696e746572666163652069640000000081525060200191505060405180910390fd5b6001600080837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200190815260200160002060006101000a81548160ff02191690831515021790555050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10620002ad57805160ff1916838001178555620002de565b82800160010185558215620002de579182015b82811115620002dd578251825591602001919060010190620002c0565b5b509050620002ed9190620002f1565b5090565b6200031691905b8082111562000312576000816000905550600101620002f8565b5090565b90565b611fe780620003296000396000f3fe608060405234801561001057600080fd5b50600436106100f55760003560e01c80636352211e11610097578063a22cb46511610066578063a22cb46514610618578063b88d4fde14610668578063c87b56dd1461076d578063e985e9c514610814576100f5565b80636352211e1461047757806370a08231146104e55780637a9adac61461053d57806395d89b4114610595576100f5565b8063095ea7b3116100d3578063095ea7b31461025057806323b872dd1461029e57806342842e0e1461030c57806350bb4e7f1461037a576100f5565b806301ffc9a7146100fa57806306fdde031461015f578063081812fc146101e2575b600080fd5b6101456004803603602081101561011057600080fd5b8101908080357bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19169060200190929190505050610890565b604051808215151515815260200191505060405180910390f35b6101676108f7565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156101a757808201518184015260208101905061018c565b50505050905090810190601f1680156101d45780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b61020e600480360360208110156101f857600080fd5b8101908080359060200190929190505050610999565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b61029c6004803603604081101561026657600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610a34565b005b61030a600480360360608110156102b457600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610c0d565b005b6103786004803603606081101561032257600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610c7c565b005b61045d6004803603606081101561039057600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190803590602001906401000000008111156103d757600080fd5b8201836020820111156103e957600080fd5b8035906020019184600183028401116401000000008311171561040b57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610c9c565b604051808215151515815260200191505060405180910390f35b6104a36004803603602081101561048d57600080fd5b8101908080359060200190929190505050610cbd565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610527600480360360208110156104fb57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610d85565b6040518082815260200191505060405180910390f35b6105936004803603606081101561055357600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019092919080359060200190929190505050610e5a565b005b61059d610ebc565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156105dd5780820151818401526020810190506105c2565b50505050905090810190601f16801561060a5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6106666004803603604081101561062e57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803515159060200190929190505050610f5e565b005b61076b6004803603608081101561067e57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190803590602001906401000000008111156106e557600080fd5b8201836020820111156106f757600080fd5b8035906020019184600183028401116401000000008311171561071957600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050611101565b005b6107996004803603602081101561078357600080fd5b8101908080359060200190929190505050611173565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156107d95780820151818401526020810190506107be565b50505050905090810190601f1680156108065780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6108766004803603604081101561082a57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050611286565b604051808215151515815260200191505060405180910390f35b6000806000837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200190815260200160002060009054906101000a900460ff169050919050565b606060058054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561098f5780601f106109645761010080835404028352916020019161098f565b820191906000526020600020905b81548152906001019060200180831161097257829003601f168201915b5050505050905090565b60006109a48261131a565b6109f9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602c815260200180611eba602c913960400191505060405180910390fd5b6002600083815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050919050565b6000610a3f82610cbd565b90508073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415610ac6576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526021815260200180611f6a6021913960400191505060405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161480610b065750610b058133611286565b5b610b5b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526038815260200180611e2f6038913960400191505060405180910390fd5b826002600084815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550818373ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92560405160405180910390a4505050565b610c17338261138c565b610c6c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526031815260200180611f8b6031913960400191505060405180910390fd5b610c77838383611480565b505050565b610c9783838360405180602001604052806000815250611101565b505050565b6000610ca884846116db565b610cb283836118f3565b600190509392505050565b6000806001600084815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415610d7c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526029815260200180611e916029913960400191505060405180910390fd5b80915050919050565b60008073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415610e0c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602a815260200180611e67602a913960400191505060405180910390fd5b610e53600360008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002061197d565b9050919050565b60008290505b81811015610eb657610ea884826040518060400160405280600781526020017f7465737455524900000000000000000000000000000000000000000000000000815250610c9c565b508080600101915050610e60565b50505050565b606060068054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610f545780601f10610f2957610100808354040283529160200191610f54565b820191906000526020600020905b815481529060010190602001808311610f3757829003601f168201915b5050505050905090565b3373ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415611000576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260198152602001807f4552433732313a20617070726f766520746f2063616c6c65720000000000000081525060200191505060405180910390fd5b80600460003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055508173ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c3183604051808215151515815260200191505060405180910390a35050565b61110c848484610c0d565b6111188484848461198b565b61116d576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526032815260200180611dad6032913960400191505060405180910390fd5b50505050565b606061117e8261131a565b6111d3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602f815260200180611f3b602f913960400191505060405180910390fd5b600760008381526020019081526020016000208054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561127a5780601f1061124f5761010080835404028352916020019161127a565b820191906000526020600020905b81548152906001019060200180831161125d57829003601f168201915b50505050509050919050565b6000600460008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16905092915050565b6000806001600084815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415915050919050565b60006113978261131a565b6113ec576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602c815260200180611e03602c913960400191505060405180910390fd5b60006113f783610cbd565b90508073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff16148061146657508373ffffffffffffffffffffffffffffffffffffffff1661144e84610999565b73ffffffffffffffffffffffffffffffffffffffff16145b8061147757506114768185611286565b5b91505092915050565b8273ffffffffffffffffffffffffffffffffffffffff166114a082610cbd565b73ffffffffffffffffffffffffffffffffffffffff161461150c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526029815260200180611f126029913960400191505060405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415611592576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526024815260200180611ddf6024913960400191505060405180910390fd5b61159b81611b74565b6115e2600360008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020611c32565b611629600360008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020611c55565b816001600083815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550808273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef60405160405180910390a4505050565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16141561177e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f4552433732313a206d696e7420746f20746865207a65726f206164647265737381525060200191505060405180910390fd5b6117878161131a565b156117fa576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601c8152602001807f4552433732313a20746f6b656e20616c7265616479206d696e7465640000000081525060200191505060405180910390fd5b816001600083815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550611893600360008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020611c55565b808273ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef60405160405180910390a45050565b6118fc8261131a565b611951576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602c815260200180611ee6602c913960400191505060405180910390fd5b80600760008481526020019081526020016000209080519060200190611978929190611d07565b505050565b600081600001549050919050565b60006119ac8473ffffffffffffffffffffffffffffffffffffffff16611c6b565b6119b95760019050611b6c565b60008473ffffffffffffffffffffffffffffffffffffffff1663150b7a02338887876040518563ffffffff1660e01b8152600401808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200183815260200180602001828103825283818151815260200191508051906020019080838360005b83811015611a94578082015181840152602081019050611a79565b50505050905090810190601f168015611ac15780820380516001836020036101000a031916815260200191505b5095505050505050602060405180830381600087803b158015611ae357600080fd5b505af1158015611af7573d6000803e3d6000fd5b505050506040513d6020811015611b0d57600080fd5b8101908080519060200190929190505050905063150b7a0260e01b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916817bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916149150505b949350505050565b600073ffffffffffffffffffffffffffffffffffffffff166002600083815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614611c2f5760006002600083815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505b50565b611c4a60018260000154611c7e90919063ffffffff16565b816000018190555050565b6001816000016000828254019250508190555050565b600080823b905060008111915050919050565b600082821115611cf6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f536166654d6174683a207375627472616374696f6e206f766572666c6f77000081525060200191505060405180910390fd5b600082840390508091505092915050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10611d4857805160ff1916838001178555611d76565b82800160010185558215611d76579182015b82811115611d75578251825591602001919060010190611d5a565b5b509050611d839190611d87565b5090565b611da991905b80821115611da5576000816000905550600101611d8d565b5090565b9056fe4552433732313a207472616e7366657220746f206e6f6e20455243373231526563656976657220696d706c656d656e7465724552433732313a207472616e7366657220746f20746865207a65726f20616464726573734552433732313a206f70657261746f7220717565727920666f72206e6f6e6578697374656e7420746f6b656e4552433732313a20617070726f76652063616c6c6572206973206e6f74206f776e6572206e6f7220617070726f76656420666f7220616c6c4552433732313a2062616c616e636520717565727920666f7220746865207a65726f20616464726573734552433732313a206f776e657220717565727920666f72206e6f6e6578697374656e7420746f6b656e4552433732313a20617070726f76656420717565727920666f72206e6f6e6578697374656e7420746f6b656e4552433732314d657461646174613a2055524920736574206f66206e6f6e6578697374656e7420746f6b656e4552433732313a207472616e73666572206f6620746f6b656e2074686174206973206e6f74206f776e4552433732314d657461646174613a2055524920717565727920666f72206e6f6e6578697374656e7420746f6b656e4552433732313a20617070726f76616c20746f2063757272656e74206f776e65724552433732313a207472616e736665722063616c6c6572206973206e6f74206f776e6572206e6f7220617070726f766564a165627a7a723058203dc2cf31fcae73ad33476512294a22e95c89669971faa65afa79ce39770638df0029"),
		deployer:                        ERC721Deployer,
		contractName:                    "ERC721 Performance Test Contract",
		Abi:                             erc721PerformanceABI,
		GenData:                         nil,
		GetBytecodeWithConstructorParam: returnBinAsIs,
		IsDeployed:                      isDeployerNonceNotZero,
		GetAddress:                      getNonce0ContractAddress,
		DoChargingWork: func(ctx *AdditionalWorkContext) {
			log.Printf("Start erc721 nft minting to the test account group(similar to erc20 token charging)")
			contract := ctx.AccGrp.GetTestContractByName(ContractErc721)
			// Use timestamp-based offset to avoid token ID collision on re-runs
			baseOffset := time.Now().Unix() * 1000000
			ConcurrentTransactionSend(ctx.AccGrp.GetValidAccGrp(), ctx.MaxConcurrency, func(idx int, acc *Account) {
				ERC721Ledger.InitializeAccount(acc.address)

				startTokenId, endTokenId := baseOffset+int64(idx*5), baseOffset+int64((idx+1)*5)
				ctx.LocalReservoir.SmartContractExecutionWithGuaranteeRetry(ctx.GCli, contract, nil, PackContractCall(erc721PerformanceABI, "registerBulk", acc.address, big.NewInt(startTokenId), big.NewInt(endTokenId)))

				for tokenId := startTokenId; tokenId < endTokenId; tokenId++ {
					ERC721Ledger.PutToken(acc.address, big.NewInt(tokenId))
				}
			})

			log.Println("End MintERC721ToTestAccounts")
		},
	}
}

func createStorageTrieContractInfo() TestContractInfo {
	return TestContractInfo{
		testNames:                       []string{"storageTrieWriteTC"},
		auctionTargetTxTypeList:         []string{},
		Bytecode:                        common.FromHex("0x608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550610f76806100606000396000f30060806040526004361061008e576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806301c0ae49146100935780630a29ae6f146101235780631fde075b146102715780636bda98c3146102da5780638da5cb5b14610389578063b912b308146103e0578063bf951c68146104d5578063f09fdbef1461053e575b600080fd5b34801561009f57600080fd5b506100a8610620565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100e85780820151818401526020810190506100cd565b50505050905090810190601f1680156101155780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561012f57600080fd5b5061018a600480360381019080803590602001908201803590602001908080601f01602080910402602001604051908101604052809392919081815260200183838082843782019150505050505091929192905050506106be565b604051808060200180602001838103835285818151815260200191508051906020019080838360005b838110156101ce5780820151818401526020810190506101b3565b50505050905090810190601f1680156101fb5780820380516001836020036101000a031916815260200191505b50838103825284818151815260200191508051906020019080838360005b83811015610234578082015181840152602081019050610219565b50505050905090810190601f1680156102615780820380516001836020036101000a031916815260200191505b5094505050505060405180910390f35b34801561027d57600080fd5b506102d8600480360381019080803590602001908201803590602001908080601f01602080910402602001604051908101604052809392919081815260200183838082843782019150505050505091929192905050506108af565b005b3480156102e657600080fd5b50610387600480360381019080803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290505050610994565b005b34801561039557600080fd5b5061039e610a93565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156103ec57600080fd5b506104d3600480360381019080803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290505050610ab8565b005b3480156104e157600080fd5b5061053c600480360381019080803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290505050610baa565b005b34801561054a57600080fd5b506105a5600480360381019080803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290505050610ca6565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156105e55780820151818401526020810190506105ca565b50505050905090810190601f1680156106125780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b60018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156106b65780601f1061068b576101008083540402835291602001916106b6565b820191906000526020600020905b81548152906001019060200180831161069957829003601f168201915b505050505081565b6060806106c9610dc3565b600084511115156106d957600080fd5b6003846040518082805190602001908083835b60208310151561071157805182526020820191506020810190506020830392506106ec565b6001836020036101000a0380198251168184511680821785525050505050509050019150509081526020016040518091039020604080519081016040529081600082018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156107e85780601f106107bd576101008083540402835291602001916107e8565b820191906000526020600020905b8154815290600101906020018083116107cb57829003601f168201915b50505050508152602001600182018054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561088a5780601f1061085f5761010080835404028352916020019161088a565b820191906000526020600020905b81548152906001019060200180831161086d57829003601f168201915b5050505050815250509050806000015181602001518191508090509250925050915091565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561090a57600080fd5b6000815111151561091a57600080fd5b6002816040518082805190602001908083835b602083101515610952578051825260208201915060208101905060208303925061092d565b6001836020036101000a038019825116818451168082178552505050505050905001915050908152602001604051809103902060006109919190610ddd565b50565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415156109ef57600080fd5b600082511115156109ff57600080fd5b60008151111515610a0f57600080fd5b806002836040518082805190602001908083835b602083101515610a485780518252602082019150602081019050602083039250610a23565b6001836020036101000a03801982511681845116808217855250505050505090500191505090815260200160405180910390209080519060200190610a8e929190610e25565b505050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60008351111515610ac857600080fd5b60008251111515610ad857600080fd5b60008151111515610ae857600080fd5b6040805190810160405280838152602001828152506003846040518082805190602001908083835b602083101515610b355780518252602082019150602081019050602083039250610b10565b6001836020036101000a03801982511681845116808217855250505050505090500191505090815260200160405180910390206000820151816000019080519060200190610b84929190610ea5565b506020820151816001019080519060200190610ba1929190610ea5565b50905050505050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141515610c0557600080fd5b60008151111515610c1557600080fd5b6003816040518082805190602001908083835b602083101515610c4d5780518252602082019150602081019050602083039250610c28565b6001836020036101000a038019825116818451168082178552505050505050905001915050908152602001604051809103902060008082016000610c919190610ddd565b600182016000610ca19190610ddd565b505050565b606060008251111515610cb857600080fd5b6002826040518082805190602001908083835b602083101515610cf05780518252602082019150602081019050602083039250610ccb565b6001836020036101000a03801982511681845116808217855250505050505090500191505090815260200160405180910390208054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610db75780601f10610d8c57610100808354040283529160200191610db7565b820191906000526020600020905b815481529060010190602001808311610d9a57829003601f168201915b50505050509050919050565b604080519081016040528060608152602001606081525090565b50805460018160011615610100020316600290046000825580601f10610e035750610e22565b601f016020900490600052602060002090810190610e219190610f25565b5b50565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10610e6657805160ff1916838001178555610e94565b82800160010185558215610e94579182015b82811115610e93578251825591602001919060010190610e78565b5b509050610ea19190610f25565b5090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10610ee657805160ff1916838001178555610f14565b82800160010185558215610f14579182015b82811115610f13578251825591602001919060010190610ef8565b5b509050610f219190610f25565b5090565b610f4791905b80821115610f43576000816000905550600101610f2b565b5090565b905600a165627a7a7230582089a867aeaa08bec696937a378160fadb7e3ffe65cc89c1e648dec0b1359cd4e00029"),
		deployer:                        StorageTrieDeployer,
		contractName:                    "Storage Trie Performance Test Contract",
		Abi:                             "",
		GenData:                         nil,
		GetBytecodeWithConstructorParam: returnBinAsIs,
		IsDeployed:                      isDeployerNonceNotZero,
		GetAddress:                      getNonce0ContractAddress,
	}
}

func createGeneralPurposeContractInfo() TestContractInfo {
	return TestContractInfo{
		testNames:               []string{"newSmartContractExecutionTC", " newFeeDelegatedSmartContractExecutionTC", "newFeeDelegatedSmartContractExecutionWithRatioTC", "ethereumTxLegacyTC", "ethereumTxAccessListTC", "newEthereumAccessListTC", "newEthereumDynamicFeeTC"},
		auctionTargetTxTypeList: []string{},
		Bytecode:                common.FromHex("0x608060405234801561001057600080fd5b506101de806100206000396000f3006080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416631a39d8ef81146100805780636353586b146100a757806370a08231146100ca578063fd6b7ef8146100f8575b3360009081526001602052604081208054349081019091558154019055005b34801561008c57600080fd5b5061009561010d565b60408051918252519081900360200190f35b6100c873ffffffffffffffffffffffffffffffffffffffff60043516610113565b005b3480156100d657600080fd5b5061009573ffffffffffffffffffffffffffffffffffffffff60043516610147565b34801561010457600080fd5b506100c8610159565b60005481565b73ffffffffffffffffffffffffffffffffffffffff1660009081526001602052604081208054349081019091558154019055565b60016020526000908152604090205481565b336000908152600160205260408120805490829055908111156101af57604051339082156108fc029083906000818181858888f193505050501561019c576101af565b3360009081526001602052604090208190555b505600a165627a7a72305820627ca46bb09478a015762806cc00c431230501118c7c26c30ac58c4e09e51c4f0029"),
		deployer:                GeneralPurposeDeployer,
		contractName:            "General Purpose Test Smart Contract",
		Abi:                     generalPurposeABI,
		GenData: func(fromAddress common.Address, _ *big.Int) []byte {
			return PackContractCall(generalPurposeABI, "reward", fromAddress)
		},
		GetBytecodeWithConstructorParam: returnBinAsIs,
		IsDeployed:                      isDeployerNonceNotZero,
		GetAddress:                      getNonce0ContractAddress,
	}
}

func createGaslessTokenContractInfo() TestContractInfo {
	return TestContractInfo{
		testNames:               []string{"gaslessTransactionTC", "gaslessRevertTransactionTC", "gaslessOnlyApproveTC"},
		auctionTargetTxTypeList: []string{"GAA", "GAS", "rGAA", "rGAS"},
		// testingGaslessContracts.TestTokenBin to maxUint256
		Bytecode:     common.FromHex("0x608060405234801561000f575f80fd5b50604051610b35380380610b3583398101604081905261002e91610167565b604051806040016040528060098152602001682a32b9ba2a37b5b2b760b91b81525060405180604001604052806002815260200161151560f21b815250816003908161007a919061022b565b506004610087828261022b565b50505061009b815f196100a160201b60201c565b5061030f565b6001600160a01b0382166100fb5760405162461bcd60e51b815260206004820152601f60248201527f45524332303a206d696e7420746f20746865207a65726f206164647265737300604482015260640160405180910390fd5b8060025f82825461010c91906102ea565b90915550506001600160a01b0382165f81815260208181526040808320805486019055518481527fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef910160405180910390a35050565b505050565b5f60208284031215610177575f80fd5b81516001600160a01b038116811461018d575f80fd5b9392505050565b634e487b7160e01b5f52604160045260245ffd5b600181811c908216806101bc57607f821691505b6020821081036101da57634e487b7160e01b5f52602260045260245ffd5b50919050565b601f82111561016257805f5260205f20601f840160051c810160208510156102055750805b601f840160051c820191505b81811015610224575f8155600101610211565b5050505050565b81516001600160401b0381111561024457610244610194565b6102588161025284546101a8565b846101e0565b602080601f83116001811461028b575f84156102745750858301515b5f19600386901b1c1916600185901b1785556102e2565b5f85815260208120601f198616915b828110156102b95788860151825594840194600190910190840161029a565b50858210156102d657878501515f19600388901b60f8161c191681555b505060018460011b0185555b505050505050565b8082018082111561030957634e487b7160e01b5f52601160045260245ffd5b92915050565b6108198061031c5f395ff3fe608060405234801561000f575f80fd5b50600436106100a6575f3560e01c8063395093511161006e578063395093511461011f57806370a082311461013257806395d89b411461015a578063a457c2d714610162578063a9059cbb14610175578063dd62ed3e14610188575f80fd5b806306fdde03146100aa578063095ea7b3146100c857806318160ddd146100eb57806323b872dd146100fd578063313ce56714610110575b5f80fd5b6100b261019b565b6040516100bf919061068a565b60405180910390f35b6100db6100d63660046106da565b61022b565b60405190151581526020016100bf565b6002545b6040519081526020016100bf565b6100db61010b366004610702565b610244565b604051601281526020016100bf565b6100db61012d3660046106da565b610267565b6100ef61014036600461073b565b6001600160a01b03165f9081526020819052604090205490565b6100b2610288565b6100db6101703660046106da565b610297565b6100db6101833660046106da565b610316565b6100ef61019636600461075b565b610323565b6060600380546101aa9061078c565b80601f01602080910402602001604051908101604052809291908181526020018280546101d69061078c565b80156102215780601f106101f857610100808354040283529160200191610221565b820191905f5260205f20905b81548152906001019060200180831161020457829003601f168201915b5050505050905090565b5f3361023881858561034d565b60019150505b92915050565b5f33610251858285610470565b61025c8585856104e8565b506001949350505050565b5f336102388185856102798383610323565b61028391906107c4565b61034d565b6060600480546101aa9061078c565b5f33816102a48286610323565b9050838110156103095760405162461bcd60e51b815260206004820152602560248201527f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f77604482015264207a65726f60d81b60648201526084015b60405180910390fd5b61025c828686840361034d565b5f336102388185856104e8565b6001600160a01b039182165f90815260016020908152604080832093909416825291909152205490565b6001600160a01b0383166103af5760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b6064820152608401610300565b6001600160a01b0382166104105760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b6064820152608401610300565b6001600160a01b038381165f8181526001602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b5f61047b8484610323565b90505f1981146104e257818110156104d55760405162461bcd60e51b815260206004820152601d60248201527f45524332303a20696e73756666696369656e7420616c6c6f77616e63650000006044820152606401610300565b6104e2848484840361034d565b50505050565b6001600160a01b03831661054c5760405162461bcd60e51b815260206004820152602560248201527f45524332303a207472616e736665722066726f6d20746865207a65726f206164604482015264647265737360d81b6064820152608401610300565b6001600160a01b0382166105ae5760405162461bcd60e51b815260206004820152602360248201527f45524332303a207472616e7366657220746f20746865207a65726f206164647260448201526265737360e81b6064820152608401610300565b6001600160a01b0383165f90815260208190526040902054818110156106255760405162461bcd60e51b815260206004820152602660248201527f45524332303a207472616e7366657220616d6f756e7420657863656564732062604482015265616c616e636560d01b6064820152608401610300565b6001600160a01b038481165f81815260208181526040808320878703905593871680835291849020805487019055925185815290927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef910160405180910390a36104e2565b602081525f82518060208401528060208501604085015e5f604082850101526040601f19601f83011684010191505092915050565b80356001600160a01b03811681146106d5575f80fd5b919050565b5f80604083850312156106eb575f80fd5b6106f4836106bf565b946020939093013593505050565b5f805f60608486031215610714575f80fd5b61071d846106bf565b925061072b602085016106bf565b9150604084013590509250925092565b5f6020828403121561074b575f80fd5b610754826106bf565b9392505050565b5f806040838503121561076c575f80fd5b610775836106bf565b9150610783602084016106bf565b90509250929050565b600181811c908216806107a057607f821691505b6020821081036107be57634e487b7160e01b5f52602260045260245ffd5b50919050565b8082018082111561023e57634e487b7160e01b5f52601160045260245ffdfea2646970667358221220f370ffb70ad18e3e54aed1704d02a0cc3e4334c3ad21efc5056b3b733731c1ff64736f6c63430008190033"),
		deployer:     GaslessTokenDeployer,
		contractName: "ERC20 Test Token for gasless swap",
		Abi:          gaslessApproveABI,
		GenData: func(gsrAddress common.Address, approveAmount *big.Int) []byte {
			return PackContractCall(gaslessApproveABI, "approve", gsrAddress, approveAmount)
		},
		GetBytecodeWithConstructorParam: func(bin []byte, _ []*Account, initialHolder *Account) []byte {
			return append(bin, PackContractCall(testingGaslessContracts.TestTokenMetaData.ABI, "", initialHolder.address)...)
		},
		IsDeployed: IsGSRExistInRegistry,
		GetAddress: getGaslessTokenAddress,
		WaitForSetup: func(gCli *client.Client) bool {
			return IsGSRExistInRegistry(gCli, nil)
		},
		DoChargingWork: func(ctx *AdditionalWorkContext) {
			if !ContainsAnyInList(ctx.TcList, []string{"gaslessTransactionTC", "gaslessOnlyApproveTC"}) && !ContainsAnyInList(ctx.TargetTxTypeList, []string{"GAA", "GAS"}) {
				return
			}
			log.Printf("Start gasless test token charging to the test account group")
			contract := ctx.AccGrp.GetTestContractByName(ContractGaslessToken)
			lenValidAccGrp := big.NewInt(int64(len(ctx.AccGrp.GetValidAccGrp())))
			lenGaslessApproveAccGrp := big.NewInt(int64(len(ctx.AccGrp.GetAccListByName(AccListForGaslessApproveTx))))
			totalChargeValue := new(big.Int).Mul(ctx.ChargeValue, new(big.Int).Add(lenValidAccGrp, lenGaslessApproveAccGrp))
			GaslessTokenDeployer.SmartContractExecutionWithGuaranteeRetry(ctx.GCli, contract, nil, PackContractCall(erc20ABI, "transfer", ctx.LocalReservoir.address, totalChargeValue))

			// accounts(validAccGrp + gaslessApproveAccGrp) should be charged.
			accounts := ctx.AccGrp.GetValidAccGrp()
			accounts = append(accounts, ctx.AccGrp.GetAccListByName(AccListForGaslessApproveTx)...)
			ConcurrentTransactionSend(accounts, ctx.MaxConcurrency, func(_ int, acc *Account) {
				ctx.LocalReservoir.SmartContractExecutionWithGuaranteeRetry(ctx.GCli, contract, nil, PackContractCall(erc20ABI, "transfer", acc.address, ctx.ChargeValue))
			})
		},
	}
}

func createWKaiaContractInfo() TestContractInfo {
	return TestContractInfo{
		testNames:                       []string{"gaslessTransactionTC", "gaslessRevertTransactionTC", "gaslessOnlyApproveTC"},
		auctionTargetTxTypeList:         []string{"GAA", "GAS", "rGAA", "rGAS"},
		Bytecode:                        common.FromHex(testingContracts.WKAIABin),
		deployer:                        WKaiaDeployer,
		contractName:                    "Wrapped Kaia Contract",
		Abi:                             "",
		GenData:                         nil,
		GetBytecodeWithConstructorParam: returnBinAsIs,
		IsDeployed:                      IsGSRExistInRegistry,
		GetAddress:                      getNonce0ContractAddress,
	}
}

func createUniswapFactoryContractInfo() TestContractInfo {
	return TestContractInfo{
		testNames:               []string{"gaslessTransactionTC", "gaslessRevertTransactionTC", "gaslessOnlyApproveTC"},
		auctionTargetTxTypeList: []string{"GAA", "GAS", "rGAA", "rGAS"},
		Bytecode:                common.FromHex(uniswapFactoryContracts.UniswapV2FactoryBin),
		deployer:                UniswapFactoryDeployer,
		contractName:            "Uniswap V2 Factory Contract",
		Abi:                     "",
		GenData:                 nil,
		GetBytecodeWithConstructorParam: func(bin []byte, _ []*Account, feeToSetter *Account) []byte {
			return append(bin, PackContractCall(uniswapFactoryContracts.UniswapV2FactoryMetaData.ABI, "", feeToSetter.address)...)
		},
		IsDeployed: IsGSRExistInRegistry,
		GetAddress: getNonce0ContractAddress,
	}
}

func createUniswapRouterContractInfo() TestContractInfo {
	return TestContractInfo{
		testNames:               []string{"gaslessTransactionTC", "gaslessRevertTransactionTC", "gaslessOnlyApproveTC"},
		auctionTargetTxTypeList: []string{"GAA", "GAS", "rGAA", "rGAS"},
		Bytecode:                common.FromHex(uniswapRouterContracts.UniswapV2Router02Bin),
		deployer:                UniswapRouterDeployer,
		contractName:            "Uniswap V2 Router Contract",
		Abi:                     "",
		GenData:                 nil,
		GetBytecodeWithConstructorParam: func(bin []byte, contracts []*Account, _ *Account) []byte {
			return append(bin, PackContractCall(uniswapRouterContracts.UniswapV2Router02MetaData.ABI, "", contracts[ContractUniswapV2Factory].address, contracts[ContractWKaia].address)...)
		},
		IsDeployed: IsGSRExistInRegistry,
		GetAddress: getNonce0ContractAddress,
	}
}

func createGaslessSwapRouterContractInfo() TestContractInfo {
	return TestContractInfo{
		testNames:               []string{"gaslessTransactionTC", "gaslessRevertTransactionTC", "gaslessOnlyApproveTC"},
		auctionTargetTxTypeList: []string{"GAA", "GAS", "rGAA", "rGAS"},
		Bytecode:                common.FromHex(gaslessContract.GaslessSwapRouterBin),
		deployer:                GaslessSwapRouterDeployer,
		contractName:            "Gasless Swap Router for testing GA",
		Abi:                     gaslessSwapABI,
		GenData: func(testTokenAddress common.Address, suggestedGasPrice *big.Int) []byte {
			R1 := new(big.Int).Mul(big.NewInt(21000), suggestedGasPrice)
			R2 := new(big.Int).Mul(big.NewInt(100000), suggestedGasPrice)
			R3 := new(big.Int).Mul(big.NewInt(500000), suggestedGasPrice)
			amountRepay := new(big.Int).Add(R1, new(big.Int).Add(R2, R3))
			amountIn := new(big.Int).Mul(amountRepay, big.NewInt(10))
			minAmountOut := amountRepay
			deadline := big.NewInt(time.Now().Add(1 * time.Hour).Unix())
			return PackContractCall(gaslessSwapABI, "swapForGas", testTokenAddress, amountIn, minAmountOut, amountRepay, deadline)
		},
		GetBytecodeWithConstructorParam: func(bin []byte, contracts []*Account, _ *Account) []byte {
			return append(bin, PackContractCall(gaslessContract.GaslessSwapRouterMetaData.ABI, "", contracts[ContractWKaia].address)...)
		},
		IsDeployed: IsGSRExistInRegistry,
		GetAddress: getGSRAddress,
		WaitForSetup: func(gCli *client.Client) bool {
			return IsGSRExistInRegistry(gCli, nil)
		},
		DoSetupWork: func(ctx *AdditionalWorkContext) {
			if IsGSRExistInRegistry(ctx.GCli, nil) {
				return
			}
			log.Printf("GSR does not exist in registry, setting up liquidity and registering GSR...")

			// Charge KAIA and gasless tokens to GSRSetupManager
			ctx.LocalReservoir.TransferSignedTxWithGuaranteeRetry(
				ctx.GCli,
				GSRSetupManager,
				new(big.Int).Add(ctx.ChargeValue, GetInitialLiquidity()),
			)
			GaslessTokenDeployer.SmartContractExecutionWithGuaranteeRetry(
				ctx.GCli,
				ctx.AccGrp.GetTestContractByName(ContractGaslessToken),
				nil,
				PackContractCall(erc20ABI, "transfer", GSRSetupManager.GetAddress(), GetInitialLiquidity()),
			)

			// Setup liquidity and Register GSR
			SetupLiquidity(ctx.GCli, ctx.AccGrp)
			RegisterGSR(ctx.GCli, ctx.AccGrp, ctx.GlobalReservoir)
		},
	}
}

// Target contract example code
// // SPDX-License-Identifier: LGPL-3.0-only
// pragma solidity ^0.8.25;
//
//	contract CounterForAuction {
//	    uint256 counterForAuction;
//	    uint256 counterForSC;
//
//	    function incForAuction() public {
//	        counterForAuction += 1;
//	    }
//
//	    function incForSC() public {
//	        counterForSC += 1;
//	    }
//
//	    function intendedRevert() public pure {
//	        revert("CounterForAuction: intended revert");
//	    }
//
//	    function getForAuction() public view returns (uint256) {
//	        return counterForAuction;
//	    }
//
//	    function getForSC() public view returns (uint256) {
//	        return counterForSC;
//	    }
//	}
func createCounterForTestAuctionContractInfo() TestContractInfo {
	return TestContractInfo{
		testNames:               []string{"auctionBidTC", "auctionRevertedBidTC"},
		auctionTargetTxTypeList: []string{},
		Bytecode:                common.FromHex("0x6080604052348015600e575f80fd5b506101678061001c5f395ff3fe608060405234801561000f575f80fd5b5060043610610055575f3560e01c80633176e1f714610059578063740b87931461006357806382595f8e1461006b578063d19f6c6514610081578063f5bfd22414610088575b5f80fd5b610061610090565b005b6100616100a8565b6001545b60405190815260200160405180910390f35b5f5461006f565b6100616100ff565b60015f808282546100a1919061010c565b9091555050565b60405162461bcd60e51b815260206004820152602260248201527f436f756e746572466f7241756374696f6e3a20696e74656e64656420726576656044820152611c9d60f21b606482015260840160405180910390fd5b6001805f8282546100a191905b8082018082111561012b57634e487b7160e01b5f52601160045260245ffd5b9291505056fea2646970667358221220ade4515edf1930c8d936019c56b111d9db63a8a5f96595f1d6ac91942140d10464736f6c63430008190033"),
		deployer:                CounterForTestAuctionDeployer,
		contractName:            "Counter for Test Auction",
		Abi:                     counterAuctionABI,
		GenData: func(_ common.Address, methodNum *big.Int) []byte {
			if methodNum == nil {
				methodNum = big.NewInt(0)
			}
			switch methodNum.Int64() {
			case 0:
				return PackContractCall(counterAuctionABI, "incForAuction")
			case 1:
				return PackContractCall(counterAuctionABI, "incForSC")
			case 2:
				return PackContractCall(counterAuctionABI, "intendedRevert")
			default:
				log.Fatalf("unknown method number: %v", methodNum)
				return nil
			}
		},
		GetBytecodeWithConstructorParam: returnBinAsIs,
		IsDeployed:                      IsAuctionEntryPointExistInRegistry,
		GetAddress:                      getNonce0ContractAddress,
	}
}

func createAuctionFeeVaultContractInfo() TestContractInfo {
	return TestContractInfo{
		testNames:               []string{"auctionBidTC", "auctionRevertedBidTC"},
		auctionTargetTxTypeList: []string{},
		Bytecode:                common.FromHex(auctionFeeVaultContracts.AuctionFeeVaultBin),
		deployer:                AuctionFeeVaultDeployer,
		contractName:            "Auction Fee Vault",
		Abi:                     "",
		GenData:                 nil,
		GetBytecodeWithConstructorParam: func(bin []byte, _ []*Account, _ *Account) []byte {
			return append(bin, PackContractCall(auctionFeeVaultContracts.AuctionFeeVaultMetaData.ABI, "", Auctioneer.address, big.NewInt(1000), big.NewInt(1000))...)
		},
		IsDeployed: IsAuctionEntryPointExistInRegistry,
		GetAddress: getNonce0ContractAddress,
	}
}

func createAuctionDepositVaultContractInfo() TestContractInfo {
	return TestContractInfo{
		testNames:               []string{"auctionBidTC", "auctionRevertedBidTC"},
		auctionTargetTxTypeList: []string{},
		Bytecode:                common.FromHex(auctionDepositVaultContracts.AuctionDepositVaultBin),
		deployer:                AuctionDepositVaultDeployer,
		contractName:            "Auction Deposit Vault",
		Abi:                     "",
		GenData: func(searcher common.Address, _ *big.Int) []byte {
			abii, err := auctionDepositVaultContracts.AuctionDepositVaultMetaData.GetAbi()
			if err != nil {
				log.Fatalf("failed to get ABI: %v", err)
			}
			data, err := abii.Pack("depositFor", searcher)
			if err != nil {
				log.Fatalf("failed to pack deposit data: %v", err)
			}
			return data
		},
		GetBytecodeWithConstructorParam: func(bin []byte, contracts []*Account, _ *Account) []byte {
			return append(bin, PackContractCall(auctionDepositVaultContracts.AuctionDepositVaultMetaData.ABI, "", Auctioneer.address, contracts[ContractAuctionFeeVault].address)...)
		},
		IsDeployed: IsAuctionEntryPointExistInRegistry,
		GetAddress: getNonce0ContractAddress,
	}
}

func createAuctionEntryPointContractInfo() TestContractInfo {
	return TestContractInfo{
		testNames:               []string{"auctionBidTC", "auctionRevertedBidTC"},
		auctionTargetTxTypeList: []string{},
		Bytecode:                common.FromHex(auctionEntryPointContracts.AuctionEntryPointBin),
		deployer:                AuctionEntryPointDeployer,
		contractName:            "Auction Entry Point",
		Abi:                     "",
		GenData:                 nil,
		GetBytecodeWithConstructorParam: func(bin []byte, contracts []*Account, _ *Account) []byte {
			return append(bin, PackContractCall(auctionEntryPointContracts.AuctionEntryPointMetaData.ABI, "", Auctioneer.address, contracts[ContractAuctionDepositVault].address, Auctioneer.address)...)
		},
		IsDeployed: IsAuctionEntryPointExistInRegistry,
		GetAddress: getAuctionEntryPointAddress,
		WaitForSetup: func(gCli *client.Client) bool {
			return IsAuctionEntryPointExistInRegistry(gCli, nil)
		},
		DoSetupWork: func(ctx *AdditionalWorkContext) {
			if IsAuctionEntryPointExistInRegistry(ctx.GCli, nil) {
				return
			}
			log.Printf("Auction Entry Point does not exist in registry, registering...")
			RegisterAuctionEntryPoint(ctx.GCli, ctx.AccGrp, ctx.GlobalReservoir)
		},
		DoChargingWork: func(ctx *AdditionalWorkContext) {
			// Deposit to the Auction Contract for each account
			log.Printf("Start depositing to the Auction Contract for each account")
			abii, err := auctionDepositVaultContracts.AuctionDepositVaultMetaData.GetAbi()
			if err != nil {
				log.Fatalf("failed to get ABI: %v", err)
			}
			ConcurrentTransactionSend(ctx.AccGrp.GetValidAccGrp(), ctx.MaxConcurrency, func(_ int, acc *Account) {
				data, err := abii.Pack("depositFor", acc.GetAddress())
				if err != nil {
					log.Fatalf("failed to pack deposit data: %v", err)
				}
				ctx.LocalReservoir.SmartContractExecutionWithGuaranteeRetry(
					ctx.GCli,
					ctx.AccGrp.GetTestContractByName(ContractAuctionDepositVault),
					ctx.ChargeValue,
					data,
				)
			})
		},
	}
}

func createCPUHeavyContractInfo() TestContractInfo {
	return TestContractInfo{
		testNames:               []string{"cpuHeavyTC"},
		auctionTargetTxTypeList: []string{},
		Bytecode:                common.FromHex("0x608060405234801561000f575f80fd5b506106e78061001d5f395ff3fe608060405234801561000f575f80fd5b506004361061004a575f3560e01c80637b395ec21461004e578063a21d942f14610063578063e71c6c821461007f578063f2a75fe414610092575b5f80fd5b61006161005c3660046105e5565b61009a565b005b61006b610171565b604051901515815260200160405180910390f35b61006161008d366004610605565b6101c2565b610061610249565b5f8267ffffffffffffffff8111156100b4576100b461061c565b6040519080825280602002602001820160405280156100dd578160200160208202803683370190505b5090505f5b815181101561011a576100f58185610644565b8282815181106101075761010761065d565b60209081029190910101526001016100e2565b50610133815f6001845161012e9190610644565b61025e565b60408051848152602081018490527fd596fdad182d29130ce218f4c1590c4b5ede105bee36690727baa6592bd2bfc8910160405180910390a1505050565b600180545f9182905b60148110156101b857600181601481106101965761019661065d565b01549150818311156101ab575f935050505090565b909150819060010161017a565b5060019250505090565b60145f5b60148110156101f6576101d98183610644565b600182601481106101ec576101ec61065d565b01556001016101c6565b5061020c5f61020760016014610644565b61043f565b60408051828152602081018490527fd596fdad182d29130ce218f4c1590c4b5ede105bee36690727baa6592bd2bfc8910160405180910390a15050565b5f8054908061025783610671565b9190505550565b5f805f838510156104375750839150829050815b8183101561039b575b85818151811061028d5761028d61065d565b60200260200101518684815181106102a7576102a761065d565b6020026020010151111580156102bc57508383105b156102d357826102cb81610671565b93505061027b565b8581815181106102e5576102e561065d565b60200260200101518683815181106102ff576102ff61065d565b6020026020010151111561031f578161031781610689565b9250506102d3565b81831015610396578582815181106103395761033961065d565b60200260200101518684815181106103535761035361065d565b602002602001015187858151811061036d5761036d61065d565b602002602001018885815181106103865761038661065d565b6020908102919091010191909152525b610272565b8582815181106103ad576103ad61065d565b60200260200101518682815181106103c7576103c761065d565b60200260200101518783815181106103e1576103e161065d565b602002602001018885815181106103fa576103fa61065d565b602090810291909101019190915252600182111561042257610422868661012e600186610644565b6104378661043184600161069e565b8661025e565b505050505050565b5f805f838510156105de5750839150829050815b81831015610556575b6001816014811061046f5761046f61065d565b0154600184601481106104845761048461065d565b01541115801561049357508383105b156104aa57826104a281610671565b93505061045c565b600181601481106104bd576104bd61065d565b0154600183601481106104d2576104d261065d565b015411156104ec57816104e481610689565b9250506104aa565b8183101561055157600182601481106105075761050761065d565b01546001846014811061051c5761051c61065d565b0154600185601481106105315761053161065d565b015f600186601481106105465761054661065d565b019290925591909155505b610453565b600182601481106105695761056961065d565b01546001826014811061057e5761057e61065d565b0154600183601481106105935761059361065d565b015f600186601481106105a8576105a861065d565b0192909255919091555060018211156105ca576105ca85610207600185610644565b6105de6105d883600161069e565b8561043f565b5050505050565b5f80604083850312156105f6575f80fd5b50508035926020909101359150565b5f60208284031215610615575f80fd5b5035919050565b634e487b7160e01b5f52604160045260245ffd5b634e487b7160e01b5f52601160045260245ffd5b8181038181111561065757610657610630565b92915050565b634e487b7160e01b5f52603260045260245ffd5b5f6001820161068257610682610630565b5060010190565b5f8161069757610697610630565b505f190190565b808201808211156106575761065761063056fea2646970667358221220c2afba695068cfa8270096d1ec1d5be5842a7db5d78cff7c762ae65a884d38fd64736f6c63430008180033"),
		deployer:                CPUHeavyDeployer,
		contractName:            "CPU Heavy Performance Test Contract",
		Abi:                     cpuHeavySortABI,
		GenData: func(_ common.Address, value *big.Int) []byte {
			return PackContractCall(cpuHeavySortABI, "sort", value, big.NewInt(1))
		},
		GetBytecodeWithConstructorParam: returnBinAsIs,
		IsDeployed:                      isDeployerNonceNotZero,
		GetAddress:                      getNonce0ContractAddress,
	}
}

func createLargeMemoContractInfo() TestContractInfo {
	return TestContractInfo{
		testNames:               []string{"largeMemoTC"},
		auctionTargetTxTypeList: []string{},
		Bytecode:                common.FromHex("0x608060405234801561000f575f80fd5b5060408051808201909152600c81526b12195b1b1bcb0815dbdc9b1960a21b60208201525f9061003f90826100dd565b5061019c565b634e487b7160e01b5f52604160045260245ffd5b600181811c9082168061006d57607f821691505b60208210810361008b57634e487b7160e01b5f52602260045260245ffd5b50919050565b601f8211156100d857805f5260205f20601f840160051c810160208510156100b65750805b601f840160051c820191505b818110156100d5575f81556001016100c2565b50505b505050565b81516001600160401b038111156100f6576100f6610045565b61010a816101048454610059565b84610091565b602080601f83116001811461013d575f84156101265750858301515b5f19600386901b1c1916600185901b178555610194565b5f85815260208120601f198616915b8281101561016b5788860151825594840194600190910190840161014c565b508582101561018857878501515f19600388901b60f8161c191681555b505060018460011b0185555b505050505050565b61042c806101a95f395ff3fe608060405234801561000f575f80fd5b506004361061003f575f3560e01c8063c040622614610043578063c15bae8414610061578063c47f002714610069575b5f80fd5b61004b61007e565b60405161005891906101a7565b60405180910390f35b61004b61010d565b61007c610077366004610207565b610198565b005b60605f805461008c906102b2565b80601f01602080910402602001604051908101604052809291908181526020018280546100b8906102b2565b80156101035780601f106100da57610100808354040283529160200191610103565b820191905f5260205f20905b8154815290600101906020018083116100e657829003601f168201915b5050505050905090565b5f8054610119906102b2565b80601f0160208091040260200160405190810160405280929190818152602001828054610145906102b2565b80156101905780601f1061016757610100808354040283529160200191610190565b820191905f5260205f20905b81548152906001019060200180831161017357829003601f168201915b505050505081565b5f6101a38282610336565b5050565b5f602080835283518060208501525f5b818110156101d3578581018301518582016040015282016101b7565b505f604082860101526040601f19601f8301168501019250505092915050565b634e487b7160e01b5f52604160045260245ffd5b5f60208284031215610217575f80fd5b813567ffffffffffffffff8082111561022e575f80fd5b818401915084601f830112610241575f80fd5b813581811115610253576102536101f3565b604051601f8201601f19908116603f0116810190838211818310171561027b5761027b6101f3565b81604052828152876020848701011115610293575f80fd5b826020860160208301375f928101602001929092525095945050505050565b600181811c908216806102c657607f821691505b6020821081036102e457634e487b7160e01b5f52602260045260245ffd5b50919050565b601f82111561033157805f5260205f20601f840160051c8101602085101561030f5750805b601f840160051c820191505b8181101561032e575f815560010161031b565b50505b505050565b815167ffffffffffffffff811115610350576103506101f3565b6103648161035e84546102b2565b846102ea565b602080601f831160018114610397575f84156103805750858301515b5f19600386901b1c1916600185901b1785556103ee565b5f85815260208120601f198616915b828110156103c5578886015182559484019460019091019084016103a6565b50858210156103e257878501515f19600388901b60f8161c191681555b505060018460011b0185555b50505050505056fea2646970667358221220c896c345552a36eea6b73606ff322b52d5dbea45843bad29a9b88b474afae6a564736f6c63430008180033"),
		deployer:                LargeMemoDeployer,
		contractName:            "Large Memo Performance Test Contract",
		Abi:                     largeMemoABI,
		GenData: func(_ common.Address, value *big.Int) []byte {
			return PackContractCall(largeMemoABI, "setName", "test_memo_"+value.String())
		},
		GetBytecodeWithConstructorParam: returnBinAsIs,
		IsDeployed:                      isDeployerNonceNotZero,
		GetAddress:                      getNonce0ContractAddress,
	}
}

func createReadApiCallContractInfo() TestContractInfo {
	return TestContractInfo{
		testNames:               []string{"readCall", "readGetStorageAt", "readEstimateGas"},
		auctionTargetTxTypeList: []string{},
		Bytecode:                common.FromHex("0x608060405260045f553480156012575f80fd5b5060898061001f5f395ff3fe6080604052348015600e575f80fd5b50600436106030575f3560e01c80636d4ce63c146034578063b8e010de146048575b5f80fd5b5f5460405190815260200160405180910390f35b60516008600155565b00fea2646970667358221220df126a7401c0e4325514b30acabd5739aa3044200494e562de86408f9223952f64736f6c63430008180033"),
		deployer:                ReadApiCallContractDeployer,
		contractName:            "Read API Call Contract Test",
		Abi:                     readApiSetABI,
		GenData: func(_ common.Address, value *big.Int) []byte {
			if value.Cmp(big.NewInt(0)) == 0 {
				return PackContractCall(readApiGetABI, "get")
			}
			return PackContractCall(readApiSetABI, "set")
		},
		GetBytecodeWithConstructorParam: returnBinAsIs,
		IsDeployed:                      isDeployerNonceNotZero,
		GetAddress:                      getNonce0ContractAddress,
	}
}

func createUserStorageContractInfo() TestContractInfo {
	return TestContractInfo{
		testNames:               []string{"userStorageSetTC", "userStorageSetGetTC"},
		auctionTargetTxTypeList: []string{},
		Bytecode:                common.FromHex("0x608060405234801561000f575f80fd5b5061011a8061001d5f395ff3fe6080604052348015600e575f80fd5b5060043610603a575f3560e01c806360fe47b114603e5780636d4ce63c14605c578063ffc9896b14607e575b5f80fd5b605a604936600460a3565b335f90815260208190526040902055565b005b335f908152602081905260409020545b60405190815260200160405180910390f35b606c608936600460b9565b6001600160a01b03165f9081526020819052604090205490565b5f6020828403121560b2575f80fd5b5035919050565b5f6020828403121560c8575f80fd5b81356001600160a01b038116811460dd575f80fd5b939250505056fea2646970667358221220eef08bdea8d0cfec01948a092eac4a88e1c151d1cb0775449f3335f0d288d9d164736f6c63430008180033"),
		deployer:                UserStorageDeployer,
		contractName:            "User Storage Test Contract",
		Abi:                     userStorageSetABI,
		GenData: func(addr common.Address, value *big.Int) []byte {
			switch value.Int64() {
			case 0:
				return PackContractCall(userStorageGetABI, "get")
			case 1:
				return PackContractCall(userStorageSetABI, "set", big.NewInt(rand.Int63n(1000)))
			default:
				return PackContractCall(userStorageGetUserDataABI, "getUserData", addr)
			}
		},
		GetBytecodeWithConstructorParam: returnBinAsIs,
		IsDeployed:                      isDeployerNonceNotZero,
		GetAddress:                      getNonce0ContractAddress,
	}
}

func createInternalTxKIP17ContractInfo() TestContractInfo {
	return TestContractInfo{
		testNames:               []string{"internalTxTC", "mintNFTTC"},
		auctionTargetTxTypeList: []string{},
		Bytecode:                common.FromHex("0x60806040523480156200001157600080fd5b506200002a6301ffc9a760e01b620000a160201b60201c565b620000426380ac58cd60e01b620000a160201b60201c565b6200005a63780e9d6360e01b620000a160201b60201c565b33600960006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550620001aa565b63ffffffff60e01b817bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614156200013e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f4b495031333a20696e76616c696420696e74657266616365206964000000000081525060200191505060405180910390fd5b6001600080837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200190815260200160002060006101000a81548160ff02191690831515021790555050565b61239080620001ba6000396000f3fe6080604052600436106100f35760003560e01c806342842e0e1161008a5780638da5cb5b116100595780638da5cb5b14610534578063a22cb4651461058b578063b88d4fde146105e8578063e985e9c5146106fa576100f3565b806342842e0e1461038a5780634f6ccce7146104055780636352211e1461045457806370a08231146104cf576100f3565b806323a5a65d116100c657806323a5a65d1461026b57806323b872dd146102965780632f745c59146103115780633993c22014610380576100f3565b806301ffc9a7146100f8578063081812fc1461016a578063095ea7b3146101e557806318160ddd14610240575b600080fd5b34801561010457600080fd5b506101506004803603602081101561011b57600080fd5b8101908080357bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19169060200190929190505050610783565b604051808215151515815260200191505060405180910390f35b34801561017657600080fd5b506101a36004803603602081101561018d57600080fd5b81019080803590602001909291905050506107ea565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156101f157600080fd5b5061023e6004803603604081101561020857600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610885565b005b34801561024c57600080fd5b50610255610a7b565b6040518082815260200191505060405180910390f35b34801561027757600080fd5b50610280610a88565b6040518082815260200191505060405180910390f35b3480156102a257600080fd5b5061030f600480360360608110156102b957600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610a8e565b005b34801561031d57600080fd5b5061036a6004803603604081101561033457600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610afd565b6040518082815260200191505060405180910390f35b610388610bbc565b005b34801561039657600080fd5b50610403600480360360608110156103ad57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610bd6565b005b34801561041157600080fd5b5061043e6004803603602081101561042857600080fd5b8101908080359060200190929190505050610bf6565b6040518082815260200191505060405180910390f35b34801561046057600080fd5b5061048d6004803603602081101561047757600080fd5b8101908080359060200190929190505050610c76565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156104db57600080fd5b5061051e600480360360208110156104f257600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610d3e565b6040518082815260200191505060405180910390f35b34801561054057600080fd5b50610549610e13565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561059757600080fd5b506105e6600480360360408110156105ae57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803515159060200190929190505050610e39565b005b3480156105f457600080fd5b506106f86004803603608081101561060b57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291908035906020019064010000000081111561067257600080fd5b82018360208201111561068457600080fd5b803590602001918460018302840111640100000000831117156106a657600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610fdc565b005b34801561070657600080fd5b506107696004803603604081101561071d57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061104e565b604051808215151515815260200191505060405180910390f35b6000806000837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200190815260200160002060009054906101000a900460ff169050919050565b60006107f5826110e2565b61084a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602b8152602001806122d8602b913960400191505060405180910390fd5b6002600083815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050919050565b600061089082610c76565b90508073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415610934576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f4b495031373a20617070726f76616c20746f2063757272656e74206f776e657281525060200191505060405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614806109745750610973813361104e565b5b6109c9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260378152602001806123036037913960400191505060405180910390fd5b826002600084815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550818373ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92560405160405180910390a4505050565b6000600780549050905090565b600a5481565b610a983382611154565b610aed576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260308152602001806121fc6030913960400191505060405180910390fd5b610af8838383611248565b505050565b6000610b0883610d3e565b8210610b5f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602a8152602001806121aa602a913960400191505060405180910390fd5b600560008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208281548110610ba957fe5b9060005260206000200154905092915050565b610bc833600a5461126c565b6001600a5401600a81905550565b610bf183838360405180602001604052806000815250610fdc565b505050565b6000610c00610a7b565b8210610c57576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602b8152602001806122ad602b913960400191505060405180910390fd5b60078281548110610c6457fe5b90600052602060002001549050919050565b6000806001600084815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415610d35576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260288152602001806121d46028913960400191505060405180910390fd5b80915050919050565b60008073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415610dc5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602981526020018061225c6029913960400191505060405180910390fd5b610e0c600360008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002061128d565b9050919050565b600960009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b3373ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415610edb576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260188152602001807f4b495031373a20617070726f766520746f2063616c6c6572000000000000000081525060200191505060405180910390fd5b80600460003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055508173ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c3183604051808215151515815260200191505060405180910390a35050565b610fe7848484610a8e565b610ff38484848461129b565b611048576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252603081526020018061222c6030913960400191505060405180910390fd5b50505050565b6000600460008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16905092915050565b6000806001600084815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415915050919050565b600061115f826110e2565b6111b4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602b81526020018061233a602b913960400191505060405180910390fd5b60006111bf83610c76565b90508073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff16148061122e57508373ffffffffffffffffffffffffffffffffffffffff16611216846107ea565b73ffffffffffffffffffffffffffffffffffffffff16145b8061123f575061123e818561104e565b5b91505092915050565b6112538383836117fd565b61125d8382611a58565b6112678282611bf6565b505050565b6112768282611cbd565b6112808282611bf6565b61128981611ed5565b5050565b600081600001549050919050565b60008060606112bf8673ffffffffffffffffffffffffffffffffffffffff16611f21565b6112ce576001925050506117f5565b8573ffffffffffffffffffffffffffffffffffffffff1663150b7a0260e01b33898888604051602401808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200183815260200180602001828103825283818151815260200191508051906020019080838360005b8381101561139e578082015181840152602081019050611383565b50505050905090810190601f1680156113cb5780820380516001836020036101000a031916815260200191505b5095505050505050604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050506040518082805190602001908083835b602083106114635780518252602082019150602081019050602083039250611440565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d80600081146114c5576040519150601f19603f3d011682016040523d82523d6000602084013e6114ca565b606091505b508092508193505050600081511415801561154e575063150b7a0260e01b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191681806020019051602081101561151c57600080fd5b81019080805190602001909291905050507bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916145b1561155e576001925050506117f5565b8573ffffffffffffffffffffffffffffffffffffffff16636745782b60e01b33898888604051602401808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200183815260200180602001828103825283818151815260200191508051906020019080838360005b8381101561162e578082015181840152602081019050611613565b50505050905090810190601f16801561165b5780820380516001836020036101000a031916815260200191505b5095505050505050604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050506040518082805190602001908083835b602083106116f357805182526020820191506020810190506020830392506116d0565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114611755576040519150601f19603f3d011682016040523d82523d6000602084013e61175a565b606091505b50809250819350505060008151141580156117de5750636745782b60e01b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168180602001905160208110156117ac57600080fd5b81019080805190602001909291905050507bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916145b156117ee576001925050506117f5565b6000925050505b949350505050565b8273ffffffffffffffffffffffffffffffffffffffff1661181d82610c76565b73ffffffffffffffffffffffffffffffffffffffff1614611889576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260288152602001806122856028913960400191505060405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16141561190f576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260238152602001806121876023913960400191505060405180910390fd5b61191881611f34565b61195f600360008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020611ff2565b6119a6600360008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020612015565b816001600083815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550808273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef60405160405180910390a4505050565b6000611ab06001600560008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208054905061202b90919063ffffffff16565b9050600060066000848152602001908152602001600020549050818114611b9d576000600560008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208381548110611b1d57fe5b9060005260206000200154905080600560008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208381548110611b7557fe5b9060005260206000200181905550816006600083815260200190815260200160002081905550505b600560008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020805480919060019003611bef9190612135565b5050505050565b600560008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020805490506006600083815260200190815260200160002081905550600560008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190806001815401808255809150509060018203906000526020600020016000909192909190915055505050565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415611d60576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601f8152602001807f4b495031373a206d696e7420746f20746865207a65726f20616464726573730081525060200191505060405180910390fd5b611d69816110e2565b15611ddc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f4b495031373a20746f6b656e20616c7265616479206d696e746564000000000081525060200191505060405180910390fd5b816001600083815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550611e75600360008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020612015565b808273ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef60405160405180910390a45050565b6007805490506008600083815260200190815260200160002081905550600781908060018154018082558091505090600182039060005260206000200160009091929091909150555050565b600080823b905060008111915050919050565b600073ffffffffffffffffffffffffffffffffffffffff166002600083815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614611fef5760006002600083815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505b50565b61200a6001826000015461202b90919063ffffffff16565b816000018190555050565b6001816000016000828254019250508190555050565b600061206d83836040518060400160405280601e81526020017f536166654d6174683a207375627472616374696f6e206f766572666c6f770000815250612075565b905092915050565b6000838311158290612122576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825283818151815260200191508051906020019080838360005b838110156120e75780820151818401526020810190506120cc565b50505050905090810190601f1680156121145780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b5060008385039050809150509392505050565b81548183558181111561215c5781836000526020600020918201910161215b9190612161565b5b505050565b61218391905b8082111561217f576000816000905550600101612167565b5090565b9056fe4b495031373a207472616e7366657220746f20746865207a65726f20616464726573734b49503137456e756d657261626c653a206f776e657220696e646578206f7574206f6620626f756e64734b495031373a206f776e657220717565727920666f72206e6f6e6578697374656e7420746f6b656e4b495031373a207472616e736665722063616c6c6572206973206e6f74206f776e6572206e6f7220617070726f7665644b495031373a207472616e7366657220746f206e6f6e204b49503137526563656976657220696d706c656d656e7465724b495031373a2062616c616e636520717565727920666f7220746865207a65726f20616464726573734b495031373a207472616e73666572206f6620746f6b656e2074686174206973206e6f74206f776e4b49503137456e756d657261626c653a20676c6f62616c20696e646578206f7574206f6620626f756e64734b495031373a20617070726f76656420717565727920666f72206e6f6e6578697374656e7420746f6b656e4b495031373a20617070726f76652063616c6c6572206973206e6f74206f776e6572206e6f7220617070726f76656420666f7220616c6c4b495031373a206f70657261746f7220717565727920666f72206e6f6e6578697374656e7420746f6b656ea165627a7a7230582090b422ce9bac7707c7845dfb25d4972d61661943d7faf5c74646a17143da65370029"),
		deployer:                InternalTxKIP17Deployer,
		contractName:            "Internal Transaction Test Kip17 Contract",
		Abi:                     mintCardABI,
		GenData: func(_ common.Address, _ *big.Int) []byte {
			return PackContractCall(mintCardABI, "mintCard")
		},
		GetBytecodeWithConstructorParam: returnBinAsIs,
		IsDeployed:                      isDeployerNonceNotZero,
		GetAddress:                      getNonce0ContractAddress,
	}
}

func createInternalTxMainContractInfo() TestContractInfo {
	return TestContractInfo{
		testNames:               []string{"internalTxTC", "mintNFTTC"},
		auctionTargetTxTypeList: []string{},
		Bytecode:                common.FromHex("0x608060405234801561001057600080fd5b506040516060806111b58339810180604052606081101561003057600080fd5b81019080805190602001909291908051906020019092919080519060200190929190505050826000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555081600181905550806002819055506002546001540160038190555033600560006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055507f40b14d6cde858ffed04e16150145bbf7e871a7aa2f50d1aa25dd9d18281c8d626000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff16604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390a17fc3038753043c1f04562d483ceb40f4e93ed15236e3c54f4e86bdb9e7f8818715826040518082815260200191505060405180910390a17f283721f6c362a0ca643c543100aa225b5d935cddffed47f683ce873ac7ac1abb816040518082815260200191505060405180910390a13373ffffffffffffffffffffffffffffffffffffffff167f84022644ce39de434e8f39c4398a3628815893104188274c37205d41c2d5096760405160405180910390a2505050610f7d806102386000396000f3fe6080604052600436106100915760003560e01c80638da5cb5b116100595780638da5cb5b14610463578063b5af960d146104ba578063c3197cfd14610525578063c8333bb21461057c578063cd1cc1a7146105a757610091565b8063150b7a02146100e15780633a850850146102455780635516885f146102705780636745782b1461029b578063719cc42b146103ff575b3373ffffffffffffffffffffffffffffffffffffffff167fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c346040518082815260200191505060405180910390a2005b3480156100ed57600080fd5b506101f16004803603608081101561010457600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291908035906020019064010000000081111561016b57600080fd5b82018360208201111561017d57600080fd5b8035906020019184600183028401116401000000008311171561019f57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505091929192905050506105d2565b60405180827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200191505060405180910390f35b34801561025157600080fd5b5061025a6105f7565b6040518082815260200191505060405180910390f35b34801561027c57600080fd5b506102856105fd565b6040518082815260200191505060405180910390f35b3480156102a757600080fd5b506103ab600480360360808110156102be57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291908035906020019064010000000081111561032557600080fd5b82018360208201111561033757600080fd5b8035906020019184600183028401116401000000008311171561035957600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610603565b60405180827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200191505060405180910390f35b6104616004803603604081101561041557600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610628565b005b34801561046f57600080fd5b50610478610e7d565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156104c657600080fd5b50610509600480360360208110156104dd57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610ea3565b604051808260ff1660ff16815260200191505060405180910390f35b34801561053157600080fd5b5061053a610ec3565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561058857600080fd5b50610591610ee8565b6040518082815260200191505060405180910390f35b3480156105b357600080fd5b506105bc610eee565b6040518082815260200191505060405180910390f35b60006040518080610ef5602f9139602f01905060405180910390209050949350505050565b60025481565b60015481565b60006040518080610f24602e9139602e01905060405180910390209050949350505050565b6003543073ffffffffffffffffffffffffffffffffffffffff163110156106b7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601f8152602001807f6e6f7420656e6f756768204b4c415920696e2074686520636f6e74726163740081525060200191505060405180910390fd5b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166040516024016040516020818303038152906040527f3993c220000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050506040518082805190602001908083835b602083106107a75780518252602082019150602081019050602083039250610784565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114610809576040519150601f19603f3d011682016040523d82523d6000602084013e61080e565b606091505b5050905080610885576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260178152602001807f6661696c20746f2063616c6c206d696e7443617264282900000000000000000081525060200191505060405180910390fd5b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b15801561092557600080fd5b505afa158015610939573d6000803e3d6000fd5b505050506040513d602081101561094f57600080fd5b8101908080519060200190929190505050116109d3576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f6d73672073656e6465722068617665206e6f204e46540000000000000000000081525060200191505060405180910390fd5b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166342842e0e30856000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16632f745c593060016000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b158015610af057600080fd5b505afa158015610b04573d6000803e3d6000fd5b505050506040513d6020811015610b1a57600080fd5b8101908080519060200190929190505050036040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060206040518083038186803b158015610b9357600080fd5b505afa158015610ba7573d6000803e3d6000fd5b505050506040513d6020811015610bbd57600080fd5b81019080805190602001909291905050506040518463ffffffff1660e01b8152600401808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019350505050600060405180830381600087803b158015610c6a57600080fd5b505af1158015610c7e573d6000803e3d6000fd5b505050508273ffffffffffffffffffffffffffffffffffffffff166108fc6002549081150290604051600060405180830381858888f19350505050610d2b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f6661696c20746f2073656e642072657761726420746f20696e7669746565000081525060200191505060405180910390fd5b6002548373ffffffffffffffffffffffffffffffffffffffff167f46db8e69822e768089db10bb036ff46bf588f534130ec318d910663075d5fd9e60405160405180910390a38173ffffffffffffffffffffffffffffffffffffffff166108fc6001549081150290604051600060405180830381858888f19350505050610db157600080fd5b6001548273ffffffffffffffffffffffffffffffffffffffff167fb9dfee5af539a65bf2b9a2cabfe2b60a37ce170a2a42fae2968bd55eddf29b3460405160405180910390a3600660008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600081819054906101000a900460ff168092919060010191906101000a81548160ff021916908360ff16021790555050600460008154809291906001019190505550505050565b600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60066020528060005260406000206000915054906101000a900460ff1681565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60035481565b6004548156fe6f6e455243373231526563656976656428616464726573732c616464726573732c75696e743235362c6279746573296f6e4b49503137526563656976656428616464726573732c616464726573732c75696e743235362c627974657329a165627a7a723058205deb3ff6cf36a0fb596868bdf7c66cbde9e02a1a89c26bc46274bb06d0d55e220029"),
		deployer:                InternalTxMainDeployer,
		contractName:            "Internal Transaction Test Main Contract",
		Abi:                     sendRewardsABI,
		GenData: func(_ common.Address, _ *big.Int) []byte {
			inviteeAddr := common.HexToAddress("0x" + fmt.Sprintf("%040x", rand.Int63()))
			hostAddr := common.HexToAddress("0x" + fmt.Sprintf("%040x", rand.Int63()))
			return PackContractCall(sendRewardsABI, "sendRewards", inviteeAddr, hostAddr)
		},
		GetBytecodeWithConstructorParam: func(bin []byte, contracts []*Account, _ *Account) []byte {
			if len(contracts) > 0 && contracts[ContractInternalTxKIP17] != nil {
				return append(bin, PackContractCall(internalTxMainConstructorABI, "", contracts[ContractInternalTxKIP17].address, big.NewInt(0), big.NewInt(0))...)
			}
			return bin
		},
		IsDeployed: isDeployerNonceNotZero,
		GetAddress: getNonce0ContractAddress,
	}
}

// createTetherLogicContractInfo creates the TetherToken logic contract (implementation)
// This is deployed first, then the proxy contract points to this address.
func createTetherLogicContractInfo() TestContractInfo {
	return TestContractInfo{
		testNames:                       []string{"erc20TransferWithBlockedListTC"},
		auctionTargetTxTypeList:         []string{},
		Bytecode:                        common.FromHex(tetherContracts.TetherContractTCBin),
		deployer:                        TetherLogicDeployer,
		contractName:                    "Tether Token Logic Contract",
		GenData:                         nil,
		GetBytecodeWithConstructorParam: returnBinAsIs,
		IsDeployed:                      isDeployerNonceNotZero,
		GetAddress:                      getNonce0ContractAddress,
	}
}

// createTetherProxyContractInfo creates the TransparentUpgradeableProxy contract
// Constructor params: (address _logic, address admin_, bytes memory _data)
// _data is the encoded initialize("Tether USD", "USD₮", 6) call
func createTetherProxyContractInfo() TestContractInfo {
	return TestContractInfo{
		testNames:               []string{"erc20TransferWithBlockedListTC"},
		auctionTargetTxTypeList: []string{},
		Bytecode:                common.FromHex(tetherContracts.TetherProxyBin),
		deployer:                TetherProxyDeployer,
		contractName:            "Tether Token Proxy Contract",
		GenData: func(recipientAddr common.Address, value *big.Int) []byte {
			// Use TetherToken ABI for function calls (transfer, balanceOf, etc.)
			abii, err := abi.JSON(strings.NewReader(tetherContracts.TetherContractTCABI))
			if err != nil {
				log.Fatalf("failed to abi.JSON: %v", err)
			}
			data, err := abii.Pack("transfer", recipientAddr, value)
			if err != nil {
				log.Fatalf("failed to abi.Pack: %v", err)
			}
			return data
		},
		GetBytecodeWithConstructorParam: func(bin []byte, contracts []*Account, deployer *Account) []byte {
			// Get the TetherLogic contract address
			logicAddr := contracts[ContractTetherLogic].address
			// Use TetherLogicDeployer as proxy admin (not the deployer)
			// This allows TetherProxyDeployer (the tx sender) to be the TetherToken owner
			// and still call logic contract functions through the proxy
			adminAddr := TetherLogicDeployer.address

			// Encode initialize("Tether USD", "USD₮", 6) for _data parameter
			tetherAbi, err := abi.JSON(strings.NewReader(tetherContracts.TetherContractTCABI))
			if err != nil {
				log.Fatalf("failed to parse TetherToken ABI: %v", err)
			}
			initData, err := tetherAbi.Pack("initialize", "Tether USD", "USD₮", uint8(6))
			if err != nil {
				log.Fatalf("failed to pack initialize data: %v", err)
			}

			// Pack proxy constructor params: (address _logic, address admin_, bytes memory _data)
			proxyAbi, err := abi.JSON(strings.NewReader(tetherContracts.TetherProxyABI))
			if err != nil {
				log.Fatalf("failed to parse TetherProxy ABI: %v", err)
			}
			constructorData, err := proxyAbi.Pack("", logicAddr, adminAddr, initData)
			if err != nil {
				log.Fatalf("failed to pack proxy constructor data: %v", err)
			}

			return append(bin, constructorData...)
		},
		IsDeployed: isDeployerNonceNotZero,
		GetAddress: getNonce0ContractAddress,
		DoChargingWork: func(ctx *AdditionalWorkContext) {
			tetherAbi, err := abi.JSON(strings.NewReader(tetherContracts.TetherContractTCABI))
			if err != nil {
				log.Fatalf("failed to parse TetherToken ABI: %v", err)
			}

			// 1. Mint Tether tokens to test accounts
			log.Printf("Start minting Tether tokens to test accounts")
			proxyContract := ctx.AccGrp.GetTestContractByName(ContractTetherProxy)
			ConcurrentTransactionSend(ctx.AccGrp.GetValidAccGrp(), ctx.MaxConcurrency, func(_ int, acc *Account) {
				data, err := tetherAbi.Pack("mint", acc.address, big.NewInt(1e18))
				if err != nil {
					log.Fatalf("failed to pack mint data: %v", err)
				}
				TetherProxyDeployer.SmartContractExecutionWithGuaranteeRetry(ctx.GCli, proxyContract, nil, data)
			})
			log.Printf("Finished minting Tether tokens")
		},
	}
}

func IsGSRExistInRegistry(gCli *client.Client, _ *Account) bool {
	return getGSRAddressInRegistry(gCli, nil) != common.Address{}
}

func IsAuctionEntryPointExistInRegistry(gCli *client.Client, _ *Account) bool {
	return getAuctionEntryPointAddressInRegistry(gCli, nil) != common.Address{}
}

func GetInitialLiquidity() *big.Int {
	return new(big.Int).Mul(big.NewInt(10000000000), big.NewInt(params.KAIA))
}

// SetupLiquidity issues six transactions.
// 1. Create a pair (from GSRSetupManager, nonce0)
// 2. Deposit KAIA to WKAIA (from GSRSetupManager, nonce1)
// 3. Approve(TestToken) (from GSRSetupManager, nonce2)
// 4. Approve(WKAIA) (from GSRSetupManager, nonce3)
// 5. Add liquidity (from GSRSetupManager, nonce4)
// 6. Add token to GSR (from GaslessSwapRouterDeployer, nonce1)
// Each nonce is fixed, and if that nonce is used, the setup is considered complete and is skipped.
func SetupLiquidity(gCli *client.Client, accGrp *AccGroup) {
	log.Printf("SetupLiquidity started...")
	/* ------------- contract initialization  ------------- */
	var (
		testTokenAddr    = accGrp.contracts[ContractGaslessToken].address
		wkaiaAddr        = accGrp.contracts[ContractWKaia].address
		factoryAddr      = accGrp.contracts[ContractUniswapV2Factory].address
		routerAddr       = accGrp.contracts[ContractUniswapV2Router].address
		gsrAddr          = accGrp.contracts[ContractGaslessSwapRouter].address
		initialLiquidity = GetInitialLiquidity()
	)

	testTokenContract, err := testingGaslessContracts.NewTestToken(testTokenAddr, gCli)
	if err != nil {
		log.Fatalf("failed to get test token contract: %v", err)
	}
	wkaiaContract, err := testingContracts.NewWKAIA(wkaiaAddr, gCli)
	if err != nil {
		log.Fatalf("failed to get wkaia contract: %v", err)
	}
	factoryContract, err := uniswapFactoryContracts.NewUniswapV2Factory(factoryAddr, gCli)
	if err != nil {
		log.Fatalf("failed to get factory contract: %v", err)
	}
	routerContract, err := uniswapRouterContracts.NewUniswapV2Router02(routerAddr, gCli)
	if err != nil {
		log.Fatalf("failed to get router contract: %v", err)
	}
	gsrContract, err := gaslessContract.NewGaslessSwapRouter(gsrAddr, gCli)
	if err != nil {
		log.Fatalf("failed to get gsr contract: %v", err)
	}

	/* ------------- create pair ------------- */
	GSRSetupManager.TryRunTxSendFunctionWithGuaranteeRetry(gCli, []error{blockchain.ErrNonceTooLow}, func(_ *client.Client, sender *Account) (*types.Transaction, error) {
		transactOpts := bind.NewKeyedTransactor(sender.privateKey[0])
		transactOpts.Nonce = big.NewInt(0)
		transactOpts.GasLimit = 3000000
		tx, err := factoryContract.CreatePair(transactOpts, testTokenAddr, wkaiaAddr)
		return tx, err
	})

	/* ------------- deposit ------------- */
	GSRSetupManager.TryRunTxSendFunctionWithGuaranteeRetry(gCli, []error{blockchain.ErrNonceTooLow}, func(_ *client.Client, sender *Account) (*types.Transaction, error) {
		transactOpts := bind.NewKeyedTransactor(sender.privateKey[0])
		transactOpts.Nonce = big.NewInt(1)
		transactOpts.Value = initialLiquidity
		transactOpts.GasLimit = 3000000
		tx, err := wkaiaContract.Deposit(transactOpts)
		return tx, err
	})

	/* ------------- approve(TestToken) ------------- */
	GSRSetupManager.TryRunTxSendFunctionWithGuaranteeRetry(gCli, []error{blockchain.ErrNonceTooLow}, func(_ *client.Client, sender *Account) (*types.Transaction, error) {
		// The nonce for GaslessTokenDeployer is as is.
		transactOpts := bind.NewKeyedTransactor(sender.privateKey[0])
		transactOpts.Nonce = big.NewInt(2)
		transactOpts.GasLimit = 3000000
		tx, err := testTokenContract.Approve(transactOpts, routerAddr, initialLiquidity)
		return tx, err
	})

	/* ------------- approve(WKAIA) ------------- */
	GSRSetupManager.TryRunTxSendFunctionWithGuaranteeRetry(gCli, []error{blockchain.ErrNonceTooLow}, func(_ *client.Client, sender *Account) (*types.Transaction, error) {
		transactOpts := bind.NewKeyedTransactor(sender.privateKey[0])
		transactOpts.Nonce = big.NewInt(3)
		transactOpts.GasLimit = 3000000
		tx, err := wkaiaContract.Approve(transactOpts, routerAddr, initialLiquidity)
		return tx, err
	})

	/* ------------- add liquidity ------------- */
	GSRSetupManager.TryRunTxSendFunctionWithGuaranteeRetry(gCli, []error{blockchain.ErrNonceTooLow}, func(_ *client.Client, sender *Account) (*types.Transaction, error) {
		transactOpts := bind.NewKeyedTransactor(sender.privateKey[0])
		transactOpts.Nonce = big.NewInt(4)
		transactOpts.GasLimit = 3000000
		deadline := time.Now().Unix() + 60*20
		tx, err := routerContract.AddLiquidity(transactOpts, testTokenAddr, wkaiaAddr,
			initialLiquidity, initialLiquidity, common.Big0, common.Big0, sender.address, big.NewInt(deadline))
		return tx, err
	})

	/* ------------- add token to gsr ------------- */
	// Because the AddToken can be called by only the owner, need to use the deployer account.
	GaslessSwapRouterDeployer.TryRunTxSendFunctionWithGuaranteeRetry(gCli, []error{blockchain.ErrNonceTooLow}, func(_ *client.Client, sender *Account) (*types.Transaction, error) {
		transactOpts := bind.NewKeyedTransactor(sender.privateKey[0])
		transactOpts.Nonce = big.NewInt(1)
		transactOpts.GasLimit = 3000000
		tx, err := gsrContract.AddToken(transactOpts, testTokenAddr, factoryAddr, routerAddr)
		return tx, err
	})
	log.Printf("SetupLiquidity finished...")
}

// RegisterGSR registers a GsrAddress from a globalReservoirAccount.
// The GsrAddress is always 0x8a9af77d180CE8377437f82504f739bFe4074839 because it is determined that it is created by the nonce0 of GaslessSwapRouterDeployer.
// Therefore, an address that conflicts with another slave will not be registered.
func RegisterGSR(gCli *client.Client, accGrp *AccGroup, globalReservoirAccount *Account) {
	log.Printf("RegisterGSR started...")
	registry, err := kip149contract.NewRegistry(system.RegistryAddr, gCli)
	if err != nil {
		return
	}
	ctx := context.Background()
	blockNum, err := gCli.BlockNumber(ctx)
	if err != nil {
		return
	}

	targetBlockNum := new(big.Int).Add(blockNum, big.NewInt(10))
	globalReservoirAccount.TryRunTxSendFunctionWithGuaranteeRetry(gCli, []error{}, func(_ *client.Client, sender *Account) (*types.Transaction, error) {
		transactOpts := bind.NewKeyedTransactor(sender.privateKey[0])
		transactOpts.GasLimit = 3000000
		return registry.Register(transactOpts, gaslessImpl.GaslessSwapRouterName, accGrp.contracts[ContractGaslessSwapRouter].address, targetBlockNum)
	})

	// Wait until desired targetBlockNum plus 10 seconds margin
	timeoutSec := targetBlockNum.Uint64() - blockNum.Uint64() + 10
	timeout := time.NewTimer(time.Duration(timeoutSec) * time.Second)
	defer timeout.Stop()

	for {
		select {
		case <-timeout.C:
			log.Fatalf("Timeout waiting for target block %d", targetBlockNum.Uint64())
			return
		default:
			time.Sleep(1 * time.Second)
			blockNum, err := gCli.BlockNumber(ctx)
			if err != nil {
				continue
			}
			if blockNum.Cmp(targetBlockNum) >= 0 {
				log.Printf("Registered GaslessSwapRouter address %s at block %d", accGrp.contracts[ContractGaslessSwapRouter].address.String(), targetBlockNum.Uint64())
				return
			}
			log.Printf("Waiting for target block %d, current block %d", targetBlockNum.Uint64(), blockNum.Uint64())
		}
	}
}

// RegisterAuctionEntryPoint registers a AuctionEntryPointAddress from a globalReservoirAccount.
// The AuctionEntryPointAddress is always 0x259c74F5aBbc66D6015EfD15C2A80E8e10a1b435 because it is determined that it is created by the nonce0 of AuctionEntryPointDeployer.
// Therefore, an address that conflicts with another slave will not be registered.
func RegisterAuctionEntryPoint(gCli *client.Client, accGrp *AccGroup, globalReservoirAccount *Account) {
	log.Printf("RegisterAuctionEntryPoint started...")
	registry, err := kip149contract.NewRegistry(system.RegistryAddr, gCli)
	if err != nil {
		return
	}
	ctx := context.Background()
	blockNum, err := gCli.BlockNumber(ctx)
	if err != nil {
		return
	}

	targetBlockNum := new(big.Int).Add(blockNum, big.NewInt(10))
	globalReservoirAccount.TryRunTxSendFunctionWithGuaranteeRetry(gCli, []error{}, func(_ *client.Client, sender *Account) (*types.Transaction, error) {
		transactOpts := bind.NewKeyedTransactor(sender.privateKey[0])
		transactOpts.GasLimit = 3000000
		return registry.Register(transactOpts, system.AuctionEntryPointName, accGrp.contracts[ContractAuctionEntryPoint].address, targetBlockNum)
	})

	// Wait until desired targetBlockNum plus 10 seconds margin
	timeoutSec := targetBlockNum.Uint64() - blockNum.Uint64() + 10
	timeout := time.NewTimer(time.Duration(timeoutSec) * time.Second)
	defer timeout.Stop()

	for {
		select {
		case <-timeout.C:
			log.Fatalf("Timeout waiting for target block %d", targetBlockNum.Uint64())
			return
		default:
			time.Sleep(1 * time.Second)
			blockNum, err := gCli.BlockNumber(ctx)
			if err != nil {
				continue
			}
			if blockNum.Cmp(targetBlockNum) >= 0 {
				log.Printf("Registered AuctionEntryPoint address %s at block %d", accGrp.contracts[ContractAuctionEntryPoint].address.String(), targetBlockNum.Uint64())
				return
			}
			log.Printf("Waiting for target block %d, current block %d", targetBlockNum.Uint64(), blockNum.Uint64())
		}
	}
}

func returnBinAsIs(bin []byte, _ []*Account, _ *Account) []byte {
	return bin
}

func getNonce0ContractAddress(_ *client.Client, deployer *Account) common.Address {
	if deployer == nil {
		return common.Address{}
	}
	return crypto.CreateAddress(deployer.GetAddress(), 0)
}

// isDeployerNonceNotZero checks if deployer has already deployed (nonce > 0)
func isDeployerNonceNotZero(gCli *client.Client, deployer *Account) bool {
	return deployer.GetNonce(gCli) > 0
}

func getGSRAddressInRegistry(gCli *client.Client, _ *Account) common.Address {
	registry, err := kip149contract.NewRegistry(system.RegistryAddr, gCli)
	if err != nil {
		return common.Address{}
	}
	addr, err := registry.GetActiveAddr(&bind.CallOpts{}, gaslessImpl.GaslessSwapRouterName)
	if err != nil {
		return common.Address{}
	}
	return addr
}

func getAuctionEntryPointAddressInRegistry(gCli *client.Client, _ *Account) common.Address {
	registry, err := kip149contract.NewRegistry(system.RegistryAddr, gCli)
	if err != nil {
		return common.Address{}
	}
	addr, err := registry.GetActiveAddr(&bind.CallOpts{}, system.AuctionEntryPointName)
	if err != nil {
		return common.Address{}
	}
	return addr
}

// Gives priority to data obtained from the chain.
func getGSRAddress(gCli *client.Client, _ *Account) common.Address {
	addressExpectedFromDeployer := getNonce0ContractAddress(gCli, GaslessSwapRouterDeployer)
	addr := getGSRAddressInRegistry(gCli, nil)
	if addr != (common.Address{}) {
		return addr
	}
	return addressExpectedFromDeployer
}

// Gives priority to data obtained from the chain.
func getAuctionEntryPointAddress(gCli *client.Client, _ *Account) common.Address {
	addressExpectedFromDeployer := getNonce0ContractAddress(gCli, AuctionEntryPointDeployer)
	addr := getAuctionEntryPointAddressInRegistry(gCli, nil)
	if addr != (common.Address{}) {
		return addr
	}
	return addressExpectedFromDeployer
}

func getEntrypointNonce(gCli *client.Client, searcher common.Address) *big.Int {
	auctionEntryPoint, err := auctionEntryPointContracts.NewAuctionEntryPoint(getAuctionEntryPointAddress(gCli, nil), gCli)
	if err != nil {
		return big.NewInt(0)
	}
	nonce, err := auctionEntryPoint.Nonces(&bind.CallOpts{}, searcher)
	if err != nil {
		return big.NewInt(0)
	}
	return nonce
}

// Gives priority to data obtained from the chain.
func getGaslessTokenAddress(gCli *client.Client, deployer *Account) common.Address {
	// If a supported token cannot be obtained from the registry, it will be assumed that the test token deployer has deployed it.
	addressExpectedFromDeployer := getNonce0ContractAddress(gCli, GaslessTokenDeployer)
	addr := getGSRAddress(gCli, GaslessSwapRouterDeployer)
	if addr == (common.Address{}) {
		return addressExpectedFromDeployer
	}
	gaslessSwapRouter, err := gaslessContract.NewGaslessSwapRouter(addr, gCli)
	if err != nil {
		return addressExpectedFromDeployer
	}
	supportedTokens, err := gaslessSwapRouter.GetSupportedTokens(&bind.CallOpts{})
	if err != nil {
		return addressExpectedFromDeployer
	}
	if len(supportedTokens) == 0 {
		return addressExpectedFromDeployer
	}
	return supportedTokens[0]
}
