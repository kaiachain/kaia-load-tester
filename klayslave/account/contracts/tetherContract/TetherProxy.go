// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package tetherContractTC

import (
	"errors"
	"math/big"
	"strings"

	kaia "github.com/kaiachain/kaia"
	"github.com/kaiachain/kaia/accounts/abi"
	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = kaia.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// TetherProxyMetaData contains all meta data concerning the TetherProxy contract.
var TetherProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_logic\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"admin_\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"admin_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"changeAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"implementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"implementation_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x6080604052604051610d9e380380610d9e83398101604081905261002291610422565b828161004f60017f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbd6104f3565b5f516020610d575f395f51905f521461006a5761006a610512565b61007582825f6100cf565b506100a3905060017fb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d61046104f3565b5f516020610d375f395f51905f52146100be576100be610512565b6100c7826100fa565b505050610571565b6100d883610167565b5f825111806100e45750805b156100f5576100f383836101a6565b505b505050565b7f7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f6101395f516020610d375f395f51905f52546001600160a01b031690565b604080516001600160a01b03928316815291841660208301520160405180910390a1610164816101d4565b50565b6101708161026f565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b905f90a250565b60606101cb8383604051806060016040528060278152602001610d77602791396102e7565b90505b92915050565b6001600160a01b03811661023e5760405162461bcd60e51b815260206004820152602660248201527f455243313936373a206e65772061646d696e20697320746865207a65726f206160448201526564647265737360d01b60648201526084015b60405180910390fd5b805f516020610d375f395f51905f525b80546001600160a01b0319166001600160a01b039290921691909117905550565b803b6102d35760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b6064820152608401610235565b805f516020610d575f395f51905f5261024e565b6060833b6103465760405162461bcd60e51b815260206004820152602660248201527f416464726573733a2064656c65676174652063616c6c20746f206e6f6e2d636f6044820152651b9d1c9858dd60d21b6064820152608401610235565b5f5f856001600160a01b0316856040516103609190610526565b5f60405180830381855af49150503d805f8114610398576040519150601f19603f3d011682016040523d82523d5f602084013e61039d565b606091505b5090925090506103ae8282866103ba565b925050505b9392505050565b606083156103c95750816103b3565b8251156103d95782518084602001fd5b8160405162461bcd60e51b8152600401610235919061053c565b80516001600160a01b0381168114610409575f5ffd5b919050565b634e487b7160e01b5f52604160045260245ffd5b5f5f5f60608486031215610434575f5ffd5b61043d846103f3565b925061044b602085016103f3565b60408501519092506001600160401b03811115610466575f5ffd5b8401601f81018613610476575f5ffd5b80516001600160401b0381111561048f5761048f61040e565b604051601f8201601f19908116603f011681016001600160401b03811182821017156104bd576104bd61040e565b6040528181528282016020018810156104d4575f5ffd5b8160208401602083015e5f602083830101528093505050509250925092565b818103818111156101ce57634e487b7160e01b5f52601160045260245ffd5b634e487b7160e01b5f52600160045260245ffd5b5f82518060208501845e5f920191825250919050565b602081525f82518060208401528060208501604085015e5f604082850101526040601f19601f83011684010191505092915050565b6107b98061057e5f395ff3fe60806040526004361061004d575f3560e01c80633659cfe6146100645780634f1ef286146100835780635c60da1b146100965780638f283970146100c6578063f851a440146100e55761005c565b3661005c5761005a6100f9565b005b61005a6100f9565b34801561006f575f5ffd5b5061005a61007e36600461067a565b610113565b61005a610091366004610693565b61014e565b3480156100a1575f5ffd5b506100aa6101b4565b6040516001600160a01b03909116815260200160405180910390f35b3480156100d1575f5ffd5b5061005a6100e036600461067a565b6101e4565b3480156100f0575f5ffd5b506100aa610204565b610101610224565b61011161010c6102b9565b6102c2565b565b61011b6102e0565b6001600160a01b03163303610146576101438160405180602001604052805f8152505f610312565b50565b6101436100f9565b6101566102e0565b6001600160a01b031633036101ac576101a78383838080601f0160208091040260200160405190810160405280939291908181526020018383808284375f9201919091525060019250610312915050565b505050565b6101a76100f9565b5f6101bd6102e0565b6001600160a01b031633036101d9576101d46102b9565b905090565b6101e16100f9565b90565b6101ec6102e0565b6001600160a01b03163303610146576101438161033c565b5f61020d6102e0565b6001600160a01b031633036101d9576101d46102e0565b61022c6102e0565b6001600160a01b031633036101115760405162461bcd60e51b815260206004820152604260248201527f5472616e73706172656e745570677261646561626c6550726f78793a2061646d60448201527f696e2063616e6e6f742066616c6c6261636b20746f2070726f78792074617267606482015261195d60f21b608482015260a4015b60405180910390fd5b5f6101d4610390565b365f5f375f5f365f845af43d5f5f3e8080156102dc573d5ff35b3d5ffd5b5f7fb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d61035b546001600160a01b0316919050565b61031b836103b7565b5f825111806103275750805b156101a75761033683836103f6565b50505050565b7f7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f6103656102e0565b604080516001600160a01b03928316815291841660208301520160405180910390a161014381610422565b5f7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc610303565b6103c0816104cb565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b905f90a250565b606061041b838360405180606001604052806027815260200161075d60279139610556565b9392505050565b6001600160a01b0381166104875760405162461bcd60e51b815260206004820152602660248201527f455243313936373a206e65772061646d696e20697320746865207a65726f206160448201526564647265737360d01b60648201526084016102b0565b807fb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d61035b80546001600160a01b0319166001600160a01b039290921691909117905550565b803b61052f5760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b60648201526084016102b0565b807f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc6104aa565b6060833b6105b55760405162461bcd60e51b815260206004820152602660248201527f416464726573733a2064656c65676174652063616c6c20746f206e6f6e2d636f6044820152651b9d1c9858dd60d21b60648201526084016102b0565b5f5f856001600160a01b0316856040516105cf9190610711565b5f60405180830381855af49150503d805f8114610607576040519150601f19603f3d011682016040523d82523d5f602084013e61060c565b606091505b509150915061061c828286610626565b9695505050505050565b6060831561063557508161041b565b8251156106455782518084602001fd5b8160405162461bcd60e51b81526004016102b09190610727565b80356001600160a01b0381168114610675575f5ffd5b919050565b5f6020828403121561068a575f5ffd5b61041b8261065f565b5f5f5f604084860312156106a5575f5ffd5b6106ae8461065f565b9250602084013567ffffffffffffffff8111156106c9575f5ffd5b8401601f810186136106d9575f5ffd5b803567ffffffffffffffff8111156106ef575f5ffd5b866020828401011115610700575f5ffd5b939660209190910195509293505050565b5f82518060208501845e5f920191825250919050565b602081525f82518060208401528060208501604085015e5f604082850101526040601f19601f8301168401019150509291505056fe416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a2646970667358221220502d440fde437e75c6c41f0a27a53d87f511f45629f7ea6a300e1c19b1cda6e164736f6c634300081f0033b53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564",
}

// TetherProxyABI is the input ABI used to generate the binding from.
// Deprecated: Use TetherProxyMetaData.ABI instead.
var TetherProxyABI = TetherProxyMetaData.ABI

// TetherProxyBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TetherProxyMetaData.Bin instead.
var TetherProxyBin = TetherProxyMetaData.Bin

// DeployTetherProxy deploys a new Ethereum contract, binding an instance of TetherProxy to it.
func DeployTetherProxy(auth *bind.TransactOpts, backend bind.ContractBackend, _logic common.Address, admin_ common.Address, _data []byte) (common.Address, *types.Transaction, *TetherProxy, error) {
	parsed, err := TetherProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TetherProxyBin), backend, _logic, admin_, _data)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TetherProxy{TetherProxyCaller: TetherProxyCaller{contract: contract}, TetherProxyTransactor: TetherProxyTransactor{contract: contract}, TetherProxyFilterer: TetherProxyFilterer{contract: contract}}, nil
}

// TetherProxy is an auto generated Go binding around an Ethereum contract.
type TetherProxy struct {
	TetherProxyCaller     // Read-only binding to the contract
	TetherProxyTransactor // Write-only binding to the contract
	TetherProxyFilterer   // Log filterer for contract events
}

// TetherProxyCaller is an auto generated read-only Go binding around an Ethereum contract.
type TetherProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TetherProxyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TetherProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TetherProxyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TetherProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TetherProxySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TetherProxySession struct {
	Contract     *TetherProxy      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TetherProxyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TetherProxyCallerSession struct {
	Contract *TetherProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// TetherProxyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TetherProxyTransactorSession struct {
	Contract     *TetherProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// TetherProxyRaw is an auto generated low-level Go binding around an Ethereum contract.
type TetherProxyRaw struct {
	Contract *TetherProxy // Generic contract binding to access the raw methods on
}

// TetherProxyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TetherProxyCallerRaw struct {
	Contract *TetherProxyCaller // Generic read-only contract binding to access the raw methods on
}

// TetherProxyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TetherProxyTransactorRaw struct {
	Contract *TetherProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTetherProxy creates a new instance of TetherProxy, bound to a specific deployed contract.
func NewTetherProxy(address common.Address, backend bind.ContractBackend) (*TetherProxy, error) {
	contract, err := bindTetherProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TetherProxy{TetherProxyCaller: TetherProxyCaller{contract: contract}, TetherProxyTransactor: TetherProxyTransactor{contract: contract}, TetherProxyFilterer: TetherProxyFilterer{contract: contract}}, nil
}

// NewTetherProxyCaller creates a new read-only instance of TetherProxy, bound to a specific deployed contract.
func NewTetherProxyCaller(address common.Address, caller bind.ContractCaller) (*TetherProxyCaller, error) {
	contract, err := bindTetherProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TetherProxyCaller{contract: contract}, nil
}

// NewTetherProxyTransactor creates a new write-only instance of TetherProxy, bound to a specific deployed contract.
func NewTetherProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*TetherProxyTransactor, error) {
	contract, err := bindTetherProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TetherProxyTransactor{contract: contract}, nil
}

// NewTetherProxyFilterer creates a new log filterer instance of TetherProxy, bound to a specific deployed contract.
func NewTetherProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*TetherProxyFilterer, error) {
	contract, err := bindTetherProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TetherProxyFilterer{contract: contract}, nil
}

// bindTetherProxy binds a generic wrapper to an already deployed contract.
func bindTetherProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TetherProxyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TetherProxy *TetherProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TetherProxy.Contract.TetherProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TetherProxy *TetherProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TetherProxy.Contract.TetherProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TetherProxy *TetherProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TetherProxy.Contract.TetherProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TetherProxy *TetherProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TetherProxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TetherProxy *TetherProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TetherProxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TetherProxy *TetherProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TetherProxy.Contract.contract.Transact(opts, method, params...)
}

// Admin is a paid mutator transaction binding the contract method 0xf851a440.
//
// Solidity: function admin() returns(address admin_)
func (_TetherProxy *TetherProxyTransactor) Admin(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TetherProxy.contract.Transact(opts, "admin")
}

// Admin is a paid mutator transaction binding the contract method 0xf851a440.
//
// Solidity: function admin() returns(address admin_)
func (_TetherProxy *TetherProxySession) Admin() (*types.Transaction, error) {
	return _TetherProxy.Contract.Admin(&_TetherProxy.TransactOpts)
}

// Admin is a paid mutator transaction binding the contract method 0xf851a440.
//
// Solidity: function admin() returns(address admin_)
func (_TetherProxy *TetherProxyTransactorSession) Admin() (*types.Transaction, error) {
	return _TetherProxy.Contract.Admin(&_TetherProxy.TransactOpts)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdmin) returns()
func (_TetherProxy *TetherProxyTransactor) ChangeAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _TetherProxy.contract.Transact(opts, "changeAdmin", newAdmin)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdmin) returns()
func (_TetherProxy *TetherProxySession) ChangeAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _TetherProxy.Contract.ChangeAdmin(&_TetherProxy.TransactOpts, newAdmin)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdmin) returns()
func (_TetherProxy *TetherProxyTransactorSession) ChangeAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _TetherProxy.Contract.ChangeAdmin(&_TetherProxy.TransactOpts, newAdmin)
}

// Implementation is a paid mutator transaction binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() returns(address implementation_)
func (_TetherProxy *TetherProxyTransactor) Implementation(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TetherProxy.contract.Transact(opts, "implementation")
}

// Implementation is a paid mutator transaction binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() returns(address implementation_)
func (_TetherProxy *TetherProxySession) Implementation() (*types.Transaction, error) {
	return _TetherProxy.Contract.Implementation(&_TetherProxy.TransactOpts)
}

// Implementation is a paid mutator transaction binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() returns(address implementation_)
func (_TetherProxy *TetherProxyTransactorSession) Implementation() (*types.Transaction, error) {
	return _TetherProxy.Contract.Implementation(&_TetherProxy.TransactOpts)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_TetherProxy *TetherProxyTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _TetherProxy.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_TetherProxy *TetherProxySession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _TetherProxy.Contract.UpgradeTo(&_TetherProxy.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_TetherProxy *TetherProxyTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _TetherProxy.Contract.UpgradeTo(&_TetherProxy.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_TetherProxy *TetherProxyTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _TetherProxy.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_TetherProxy *TetherProxySession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _TetherProxy.Contract.UpgradeToAndCall(&_TetherProxy.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_TetherProxy *TetherProxyTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _TetherProxy.Contract.UpgradeToAndCall(&_TetherProxy.TransactOpts, newImplementation, data)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_TetherProxy *TetherProxyTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _TetherProxy.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_TetherProxy *TetherProxySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TetherProxy.Contract.Fallback(&_TetherProxy.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_TetherProxy *TetherProxyTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TetherProxy.Contract.Fallback(&_TetherProxy.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_TetherProxy *TetherProxyTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TetherProxy.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_TetherProxy *TetherProxySession) Receive() (*types.Transaction, error) {
	return _TetherProxy.Contract.Receive(&_TetherProxy.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_TetherProxy *TetherProxyTransactorSession) Receive() (*types.Transaction, error) {
	return _TetherProxy.Contract.Receive(&_TetherProxy.TransactOpts)
}

// TetherProxyAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the TetherProxy contract.
type TetherProxyAdminChangedIterator struct {
	Event *TetherProxyAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  kaia.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TetherProxyAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TetherProxyAdminChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TetherProxyAdminChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TetherProxyAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TetherProxyAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TetherProxyAdminChanged represents a AdminChanged event raised by the TetherProxy contract.
type TetherProxyAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_TetherProxy *TetherProxyFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*TetherProxyAdminChangedIterator, error) {

	logs, sub, err := _TetherProxy.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &TetherProxyAdminChangedIterator{contract: _TetherProxy.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_TetherProxy *TetherProxyFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *TetherProxyAdminChanged) (event.Subscription, error) {

	logs, sub, err := _TetherProxy.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TetherProxyAdminChanged)
				if err := _TetherProxy.contract.UnpackLog(event, "AdminChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_TetherProxy *TetherProxyFilterer) ParseAdminChanged(log types.Log) (*TetherProxyAdminChanged, error) {
	event := new(TetherProxyAdminChanged)
	if err := _TetherProxy.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TetherProxyBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the TetherProxy contract.
type TetherProxyBeaconUpgradedIterator struct {
	Event *TetherProxyBeaconUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  kaia.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TetherProxyBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TetherProxyBeaconUpgraded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TetherProxyBeaconUpgraded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TetherProxyBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TetherProxyBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TetherProxyBeaconUpgraded represents a BeaconUpgraded event raised by the TetherProxy contract.
type TetherProxyBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_TetherProxy *TetherProxyFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*TetherProxyBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _TetherProxy.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &TetherProxyBeaconUpgradedIterator{contract: _TetherProxy.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_TetherProxy *TetherProxyFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *TetherProxyBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _TetherProxy.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TetherProxyBeaconUpgraded)
				if err := _TetherProxy.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBeaconUpgraded is a log parse operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_TetherProxy *TetherProxyFilterer) ParseBeaconUpgraded(log types.Log) (*TetherProxyBeaconUpgraded, error) {
	event := new(TetherProxyBeaconUpgraded)
	if err := _TetherProxy.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TetherProxyUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the TetherProxy contract.
type TetherProxyUpgradedIterator struct {
	Event *TetherProxyUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  kaia.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TetherProxyUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TetherProxyUpgraded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TetherProxyUpgraded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TetherProxyUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TetherProxyUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TetherProxyUpgraded represents a Upgraded event raised by the TetherProxy contract.
type TetherProxyUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_TetherProxy *TetherProxyFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*TetherProxyUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _TetherProxy.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &TetherProxyUpgradedIterator{contract: _TetherProxy.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_TetherProxy *TetherProxyFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *TetherProxyUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _TetherProxy.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TetherProxyUpgraded)
				if err := _TetherProxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_TetherProxy *TetherProxyFilterer) ParseUpgraded(log types.Log) (*TetherProxyUpgraded, error) {
	event := new(TetherProxyUpgraded)
	if err := _TetherProxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
