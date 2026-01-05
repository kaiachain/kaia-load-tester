// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package tetherContractTC

import (
	"errors"
	"math/big"
	"strings"

	"github.com/kaiachain/kaia"
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

// TetherContractTCMetaData contains all meta data concerning the TetherContractTC contract.
var TetherContractTCMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"BlockPlaced\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"BlockReleased\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_blockedUser\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_balance\",\"type\":\"uint256\"}],\"name\":\"DestroyedBlockedFunds\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_destination\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"Mint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"Redeem\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DOMAIN_SEPARATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"addToBlockedList\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_blockedUser\",\"type\":\"address\"}],\"name\":\"destroyBlockedFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_symbol\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"_decimals\",\"type\":\"uint8\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"isBlocked\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"isTrusted\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_destination\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_recipients\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_values\",\"type\":\"uint256[]\"}],\"name\":\"multiTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"permit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"redeem\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"removeFromBlockedList\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b50611ca38061001c5f395ff3fe608060405234801561000f575f5ffd5b5060043610610187575f3560e01c806370a08231116100d9578063a457c2d711610093578063db006a751161006e578063db006a7514610358578063dd62ed3e1461036b578063f2fde38b1461037e578063fbac395114610391575f5ffd5b8063a457c2d71461031f578063a9059cbb14610332578063d505accf14610345575f5ffd5b806370a0823114610296578063715018a6146102be5780637ecebe00146102c65780638da5cb5b146102d957806395d89b41146102f457806396d64879146102fc575f5ffd5b80631e89d545116101445780633644e5151161011f5780633644e51514610255578063395093511461025d5780633c7c9b901461027057806340c10f1914610283575f5ffd5b80631e89d5451461021957806323b872dd1461022c578063313ce5671461023f575f5ffd5b806306fdde031461018b578063095ea7b3146101a95780630e27a385146101cc5780631624f6c6146101e157806318160ddd146101f45780631a14f44914610206575b5f5ffd5b6101936103b3565b6040516101a09190611753565b60405180910390f35b6101bc6101b73660046117a3565b610443565b60405190151581526020016101a0565b6101df6101da3660046117cb565b61045c565b005b6101df6101ef36600461189a565b610539565b6035545b6040519081526020016101a0565b6101df6102143660046117cb565b61066a565b6101df610227366004611956565b6106ba565b6101bc61023a3660046119c2565b610777565b6101005460405160ff90911681526020016101a0565b6101f86107e9565b6101bc61026b3660046117a3565b6107f7565b6101df61027e3660046117cb565b610818565b6101df6102913660046117a3565b61086b565b6101f86102a43660046117cb565b6001600160a01b03165f9081526033602052604090205490565b6101df6108b8565b6101f86102d43660046117cb565b6108cb565b60cc546040516001600160a01b0390911681526020016101a0565b6101936108e8565b6101bc61030a3660046117cb565b60ff60208190525f9182526040909120541681565b6101bc61032d3660046117a3565b6108f7565b6101bc6103403660046117a3565b61097c565b6101df6103533660046119fc565b610989565b6101df610366366004611a62565b610aea565b6101f8610379366004611a79565b610b43565b6101df61038c3660046117cb565b610b6d565b6101bc61039f3660046117cb565b60fe6020525f908152604090205460ff1681565b6060603680546103c290611aaa565b80601f01602080910402602001604051908101604052809291908181526020018280546103ee90611aaa565b80156104395780601f1061041057610100808354040283529160200191610439565b820191905f5260205f20905b81548152906001019060200180831161041c57829003601f168201915b5050505050905090565b5f33610450818585610be6565b60019150505b92915050565b610464610d0a565b6001600160a01b0381165f90815260fe602052604090205460ff166104d05760405162461bcd60e51b815260206004820181905260248201527f546574686572546f6b656e3a2075736572206973206e6f7420626c6f636b656460448201526064015b60405180910390fd5b6001600160a01b0381165f908152603360205260409020546104f28282610d64565b816001600160a01b03167f6a2859ae7902313752498feb80a014e6e7275fe964c79aa965db815db1c7f1e98260405161052d91815260200190565b60405180910390a25050565b5f54610100900460ff161580801561055757505f54600160ff909116105b806105705750303b15801561057057505f5460ff166001145b6105d35760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b60648201526084016104c7565b5f805460ff1916600117905580156105f4575f805461ff0019166101001790555b610100805460ff191660ff841617905561060c610e9e565b6106168484610ecc565b61061f84610f00565b8015610664575f805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b50505050565b610672610d0a565b6001600160a01b0381165f81815260fe6020526040808220805460ff19169055517f665918c9e02eb2fd85acca3969cb054fc84c138e60ec4af22ab6ef2fd4c93c279190a250565b8281146107155760405162461bcd60e51b815260206004820152602360248201527f546574686572546f6b656e3a206d756c74695472616e73666572206d69736d616044820152620e8c6d60eb1b60648201526084016104c7565b5f5b838110156107705761076785858381811061073457610734611adc565b905060200201602081019061074991906117cb565b84848481811061075b5761075b611adc565b9050602002013561097c565b50600101610717565b5050505050565b335f90815260fe602052604081205460ff16156107d65760405162461bcd60e51b815260206004820152601e60248201527f426c6f636b65643a206d73672e73656e64657220697320626c6f636b6564000060448201526064016104c7565b6107e1848484610f49565b949350505050565b5f6107f2610f61565b905090565b5f336104508185856108098383610b43565b6108139190611af0565b610be6565b610820610d0a565b6001600160a01b0381165f81815260fe6020526040808220805460ff19166001179055517f406bbf2d8d145125adf1198d2cf8a67c66cc4bb0ab01c37dccd4f7c0aae1e7c79190a250565b610873610d0a565b61087d8282610fda565b816001600160a01b03167f0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d41213968858260405161052d91815260200190565b6108c0610d0a565b6108c95f6110a4565b565b6001600160a01b0381165f90815260996020526040812054610456565b6060603780546103c290611aaa565b5f33816109048286610b43565b9050838110156109645760405162461bcd60e51b815260206004820152602560248201527f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f77604482015264207a65726f60d81b60648201526084016104c7565b6109718286868403610be6565b506001949350505050565b5f336104508185856110f5565b834211156109d95760405162461bcd60e51b815260206004820152601d60248201527f45524332305065726d69743a206578706972656420646561646c696e6500000060448201526064016104c7565b5f7f6e71edae12b1b97f4d1f60370fef10105fa2faae0126114a169c64845d6126c9888888610a078c6112a9565b6040805160208101969096526001600160a01b0394851690860152929091166060840152608083015260a082015260c0810186905260e0016040516020818303038152906040528051906020012090505f610a61826112d0565b90505f610a708287878761131c565b9050896001600160a01b0316816001600160a01b031614610ad35760405162461bcd60e51b815260206004820152601e60248201527f45524332305065726d69743a20696e76616c6964207369676e6174757265000060448201526064016104c7565b610ade8a8a8a610be6565b50505050505050505050565b610af2610d0a565b610b0d610b0760cc546001600160a01b031690565b82610d64565b6040518181527f702d5967f45f6513a38ffc42d6ba9bf230bd40e8f53b16363c7eb4fd2deb9a449060200160405180910390a150565b6001600160a01b039182165f90815260346020908152604080832093909416825291909152205490565b610b75610d0a565b6001600160a01b038116610bda5760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b60648201526084016104c7565b610be3816110a4565b50565b6001600160a01b038316610c485760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b60648201526084016104c7565b6001600160a01b038216610ca95760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b60648201526084016104c7565b6001600160a01b038381165f8181526034602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92591015b60405180910390a3505050565b60cc546001600160a01b031633146108c95760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016104c7565b6001600160a01b038216610dc45760405162461bcd60e51b815260206004820152602160248201527f45524332303a206275726e2066726f6d20746865207a65726f206164647265736044820152607360f81b60648201526084016104c7565b610dcf825f83611342565b6001600160a01b0382165f9081526033602052604090205481811015610e425760405162461bcd60e51b815260206004820152602260248201527f45524332303a206275726e20616d6f756e7420657863656564732062616c616e604482015261636560f01b60648201526084016104c7565b6001600160a01b0383165f8181526033602090815260408083208686039055603580548790039055518581529192917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9101610cfd565b505050565b5f54610100900460ff16610ec45760405162461bcd60e51b81526004016104c790611b0f565b6108c961142d565b5f54610100900460ff16610ef25760405162461bcd60e51b81526004016104c790611b0f565b610efc828261145c565b5050565b5f54610100900460ff16610f265760405162461bcd60e51b81526004016104c790611b0f565b610be381604051806040016040528060018152602001603160f81b81525061149b565b5f33610f568582856114db565b6109718585856110f5565b5f6107f27f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f610f8f60655490565b6066546040805160208101859052908101839052606081018290524660808201523060a08201525f9060c0016040516020818303038152906040528051906020012090509392505050565b6001600160a01b0382166110305760405162461bcd60e51b815260206004820152601f60248201527f45524332303a206d696e7420746f20746865207a65726f20616464726573730060448201526064016104c7565b61103b5f8383611342565b8060355f82825461104c9190611af0565b90915550506001600160a01b0382165f818152603360209081526040808320805486019055518481527fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef910160405180910390a35050565b60cc80546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a35050565b6001600160a01b0383166111595760405162461bcd60e51b815260206004820152602560248201527f45524332303a207472616e736665722066726f6d20746865207a65726f206164604482015264647265737360d81b60648201526084016104c7565b6001600160a01b0382166111bb5760405162461bcd60e51b815260206004820152602360248201527f45524332303a207472616e7366657220746f20746865207a65726f206164647260448201526265737360e81b60648201526084016104c7565b6111c6838383611342565b6001600160a01b0383165f908152603360205260409020548181101561123d5760405162461bcd60e51b815260206004820152602660248201527f45524332303a207472616e7366657220616d6f756e7420657863656564732062604482015265616c616e636560d01b60648201526084016104c7565b6001600160a01b038085165f8181526033602052604080822086860390559286168082529083902080548601905591517fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9061129c9086815260200190565b60405180910390a3610664565b6001600160a01b0381165f9081526099602052604090208054600181018255905b50919050565b5f6104566112dc610f61565b8360405161190160f01b602082015260228101839052604281018290525f9060620160405160208183030381529060405280519060200120905092915050565b5f5f5f61132b8787878761154d565b915091506113388161160a565b5095945050505050565b6001600160a01b0383165f90815260fe602052604090205460ff161580611373575060cc546001600160a01b031633145b6113bf5760405162461bcd60e51b815260206004820152601c60248201527f546574686572546f6b656e3a2066726f6d20697320626c6f636b65640000000060448201526064016104c7565b306001600160a01b03831603610e995760405162461bcd60e51b815260206004820152602d60248201527f546574686572546f6b656e3a207472616e7366657220746f2074686520636f6e60448201526c7472616374206164647265737360981b60648201526084016104c7565b5f54610100900460ff166114535760405162461bcd60e51b81526004016104c790611b0f565b6108c9336110a4565b5f54610100900460ff166114825760405162461bcd60e51b81526004016104c790611b0f565b603661148e8382611b9e565b506037610e998282611b9e565b5f54610100900460ff166114c15760405162461bcd60e51b81526004016104c790611b0f565b815160209283012081519190920120606591909155606655565b5f6114e68484610b43565b90505f19811461066457818110156115405760405162461bcd60e51b815260206004820152601d60248201527f45524332303a20696e73756666696369656e7420616c6c6f77616e636500000060448201526064016104c7565b6106648484848403610be6565b5f807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a083111561158257505f90506003611601565b604080515f8082526020820180845289905260ff881692820192909252606081018690526080810185905260019060a0016020604051602081039080840390855afa1580156115d3573d5f5f3e3d5ffd5b5050604051601f1901519150506001600160a01b0381166115fb575f60019250925050611601565b91505f90505b94509492505050565b5f81600481111561161d5761161d611c59565b036116255750565b600181600481111561163957611639611c59565b036116865760405162461bcd60e51b815260206004820152601860248201527f45434453413a20696e76616c6964207369676e6174757265000000000000000060448201526064016104c7565b600281600481111561169a5761169a611c59565b036116e75760405162461bcd60e51b815260206004820152601f60248201527f45434453413a20696e76616c6964207369676e6174757265206c656e6774680060448201526064016104c7565b60038160048111156116fb576116fb611c59565b03610be35760405162461bcd60e51b815260206004820152602260248201527f45434453413a20696e76616c6964207369676e6174757265202773272076616c604482015261756560f01b60648201526084016104c7565b602081525f82518060208401528060208501604085015e5f604082850101526040601f19601f83011684010191505092915050565b80356001600160a01b038116811461179e575f5ffd5b919050565b5f5f604083850312156117b4575f5ffd5b6117bd83611788565b946020939093013593505050565b5f602082840312156117db575f5ffd5b6117e482611788565b9392505050565b634e487b7160e01b5f52604160045260245ffd5b5f82601f83011261180e575f5ffd5b813567ffffffffffffffff811115611828576118286117eb565b604051601f8201601f19908116603f0116810167ffffffffffffffff81118282101715611857576118576117eb565b60405281815283820160200185101561186e575f5ffd5b816020850160208301375f918101602001919091529392505050565b803560ff8116811461179e575f5ffd5b5f5f5f606084860312156118ac575f5ffd5b833567ffffffffffffffff8111156118c2575f5ffd5b6118ce868287016117ff565b935050602084013567ffffffffffffffff8111156118ea575f5ffd5b6118f6868287016117ff565b9250506119056040850161188a565b90509250925092565b5f5f83601f84011261191e575f5ffd5b50813567ffffffffffffffff811115611935575f5ffd5b6020830191508360208260051b850101111561194f575f5ffd5b9250929050565b5f5f5f5f60408587031215611969575f5ffd5b843567ffffffffffffffff81111561197f575f5ffd5b61198b8782880161190e565b909550935050602085013567ffffffffffffffff8111156119aa575f5ffd5b6119b68782880161190e565b95989497509550505050565b5f5f5f606084860312156119d4575f5ffd5b6119dd84611788565b92506119eb60208501611788565b929592945050506040919091013590565b5f5f5f5f5f5f5f60e0888a031215611a12575f5ffd5b611a1b88611788565b9650611a2960208901611788565b95506040880135945060608801359350611a456080890161188a565b9699959850939692959460a0840135945060c09093013592915050565b5f60208284031215611a72575f5ffd5b5035919050565b5f5f60408385031215611a8a575f5ffd5b611a9383611788565b9150611aa160208401611788565b90509250929050565b600181811c90821680611abe57607f821691505b6020821081036112ca57634e487b7160e01b5f52602260045260245ffd5b634e487b7160e01b5f52603260045260245ffd5b8082018082111561045657634e487b7160e01b5f52601160045260245ffd5b6020808252602b908201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960408201526a6e697469616c697a696e6760a81b606082015260800190565b601f821115610e9957805f5260205f20601f840160051c81016020851015611b7f5750805b601f840160051c820191505b81811015610770575f8155600101611b8b565b815167ffffffffffffffff811115611bb857611bb86117eb565b611bcc81611bc68454611aaa565b84611b5a565b6020601f821160018114611bfe575f8315611be75750848201515b5f19600385901b1c1916600184901b178455610770565b5f84815260208120601f198516915b82811015611c2d5787850151825560209485019460019092019101611c0d565b5084821015611c4a57868401515f19600387901b60f8161c191681555b50505050600190811b01905550565b634e487b7160e01b5f52602160045260245ffdfea264697066735822122048ee3968f0d0ce829e3ad80695b891308612ab4aed2b27279f6ef28a6c7d322864736f6c634300081f0033",
}

// TetherContractTCABI is the input ABI used to generate the binding from.
// Deprecated: Use TetherContractTCMetaData.ABI instead.
var TetherContractTCABI = TetherContractTCMetaData.ABI

// TetherContractTCBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TetherContractTCMetaData.Bin instead.
var TetherContractTCBin = TetherContractTCMetaData.Bin

// DeployTetherContractTC deploys a new kaia contract, binding an instance of TetherContractTC to it.
func DeployTetherContractTC(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TetherContractTC, error) {
	parsed, err := TetherContractTCMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TetherContractTCBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TetherContractTC{TetherContractTCCaller: TetherContractTCCaller{contract: contract}, TetherContractTCTransactor: TetherContractTCTransactor{contract: contract}, TetherContractTCFilterer: TetherContractTCFilterer{contract: contract}}, nil
}

// TetherContractTC is an auto generated Go binding around an kaia contract.
type TetherContractTC struct {
	TetherContractTCCaller     // Read-only binding to the contract
	TetherContractTCTransactor // Write-only binding to the contract
	TetherContractTCFilterer   // Log filterer for contract events
}

// TetherContractTCCaller is an auto generated read-only Go binding around an kaia contract.
type TetherContractTCCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TetherContractTCTransactor is an auto generated write-only Go binding around an kaia contract.
type TetherContractTCTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TetherContractTCFilterer is an auto generated log filtering Go binding around an kaia contract events.
type TetherContractTCFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TetherContractTCSession is an auto generated Go binding around an kaia contract,
// with pre-set call and transact options.
type TetherContractTCSession struct {
	Contract     *TetherContractTC // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TetherContractTCCallerSession is an auto generated read-only Go binding around an kaia contract,
// with pre-set call options.
type TetherContractTCCallerSession struct {
	Contract *TetherContractTCCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// TetherContractTCTransactorSession is an auto generated write-only Go binding around an kaia contract,
// with pre-set transact options.
type TetherContractTCTransactorSession struct {
	Contract     *TetherContractTCTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// TetherContractTCRaw is an auto generated low-level Go binding around an kaia contract.
type TetherContractTCRaw struct {
	Contract *TetherContractTC // Generic contract binding to access the raw methods on
}

// TetherContractTCCallerRaw is an auto generated low-level read-only Go binding around an kaia contract.
type TetherContractTCCallerRaw struct {
	Contract *TetherContractTCCaller // Generic read-only contract binding to access the raw methods on
}

// TetherContractTCTransactorRaw is an auto generated low-level write-only Go binding around an kaia contract.
type TetherContractTCTransactorRaw struct {
	Contract *TetherContractTCTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTetherContractTC creates a new instance of TetherContractTC, bound to a specific deployed contract.
func NewTetherContractTC(address common.Address, backend bind.ContractBackend) (*TetherContractTC, error) {
	contract, err := bindTetherContractTC(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TetherContractTC{TetherContractTCCaller: TetherContractTCCaller{contract: contract}, TetherContractTCTransactor: TetherContractTCTransactor{contract: contract}, TetherContractTCFilterer: TetherContractTCFilterer{contract: contract}}, nil
}

// NewTetherContractTCCaller creates a new read-only instance of TetherContractTC, bound to a specific deployed contract.
func NewTetherContractTCCaller(address common.Address, caller bind.ContractCaller) (*TetherContractTCCaller, error) {
	contract, err := bindTetherContractTC(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TetherContractTCCaller{contract: contract}, nil
}

// NewTetherContractTCTransactor creates a new write-only instance of TetherContractTC, bound to a specific deployed contract.
func NewTetherContractTCTransactor(address common.Address, transactor bind.ContractTransactor) (*TetherContractTCTransactor, error) {
	contract, err := bindTetherContractTC(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TetherContractTCTransactor{contract: contract}, nil
}

// NewTetherContractTCFilterer creates a new log filterer instance of TetherContractTC, bound to a specific deployed contract.
func NewTetherContractTCFilterer(address common.Address, filterer bind.ContractFilterer) (*TetherContractTCFilterer, error) {
	contract, err := bindTetherContractTC(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TetherContractTCFilterer{contract: contract}, nil
}

// bindTetherContractTC binds a generic wrapper to an already deployed contract.
func bindTetherContractTC(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TetherContractTCMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TetherContractTC *TetherContractTCRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TetherContractTC.Contract.TetherContractTCCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TetherContractTC *TetherContractTCRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TetherContractTC.Contract.TetherContractTCTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TetherContractTC *TetherContractTCRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TetherContractTC.Contract.TetherContractTCTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TetherContractTC *TetherContractTCCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TetherContractTC.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TetherContractTC *TetherContractTCTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TetherContractTC.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TetherContractTC *TetherContractTCTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TetherContractTC.Contract.contract.Transact(opts, method, params...)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_TetherContractTC *TetherContractTCCaller) DOMAINSEPARATOR(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TetherContractTC.contract.Call(opts, &out, "DOMAIN_SEPARATOR")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_TetherContractTC *TetherContractTCSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _TetherContractTC.Contract.DOMAINSEPARATOR(&_TetherContractTC.CallOpts)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_TetherContractTC *TetherContractTCCallerSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _TetherContractTC.Contract.DOMAINSEPARATOR(&_TetherContractTC.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TetherContractTC *TetherContractTCCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TetherContractTC.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TetherContractTC *TetherContractTCSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _TetherContractTC.Contract.Allowance(&_TetherContractTC.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TetherContractTC *TetherContractTCCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _TetherContractTC.Contract.Allowance(&_TetherContractTC.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TetherContractTC *TetherContractTCCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TetherContractTC.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TetherContractTC *TetherContractTCSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _TetherContractTC.Contract.BalanceOf(&_TetherContractTC.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TetherContractTC *TetherContractTCCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _TetherContractTC.Contract.BalanceOf(&_TetherContractTC.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TetherContractTC *TetherContractTCCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _TetherContractTC.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TetherContractTC *TetherContractTCSession) Decimals() (uint8, error) {
	return _TetherContractTC.Contract.Decimals(&_TetherContractTC.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TetherContractTC *TetherContractTCCallerSession) Decimals() (uint8, error) {
	return _TetherContractTC.Contract.Decimals(&_TetherContractTC.CallOpts)
}

// IsBlocked is a free data retrieval call binding the contract method 0xfbac3951.
//
// Solidity: function isBlocked(address ) view returns(bool)
func (_TetherContractTC *TetherContractTCCaller) IsBlocked(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _TetherContractTC.contract.Call(opts, &out, "isBlocked", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsBlocked is a free data retrieval call binding the contract method 0xfbac3951.
//
// Solidity: function isBlocked(address ) view returns(bool)
func (_TetherContractTC *TetherContractTCSession) IsBlocked(arg0 common.Address) (bool, error) {
	return _TetherContractTC.Contract.IsBlocked(&_TetherContractTC.CallOpts, arg0)
}

// IsBlocked is a free data retrieval call binding the contract method 0xfbac3951.
//
// Solidity: function isBlocked(address ) view returns(bool)
func (_TetherContractTC *TetherContractTCCallerSession) IsBlocked(arg0 common.Address) (bool, error) {
	return _TetherContractTC.Contract.IsBlocked(&_TetherContractTC.CallOpts, arg0)
}

// IsTrusted is a free data retrieval call binding the contract method 0x96d64879.
//
// Solidity: function isTrusted(address ) view returns(bool)
func (_TetherContractTC *TetherContractTCCaller) IsTrusted(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _TetherContractTC.contract.Call(opts, &out, "isTrusted", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsTrusted is a free data retrieval call binding the contract method 0x96d64879.
//
// Solidity: function isTrusted(address ) view returns(bool)
func (_TetherContractTC *TetherContractTCSession) IsTrusted(arg0 common.Address) (bool, error) {
	return _TetherContractTC.Contract.IsTrusted(&_TetherContractTC.CallOpts, arg0)
}

// IsTrusted is a free data retrieval call binding the contract method 0x96d64879.
//
// Solidity: function isTrusted(address ) view returns(bool)
func (_TetherContractTC *TetherContractTCCallerSession) IsTrusted(arg0 common.Address) (bool, error) {
	return _TetherContractTC.Contract.IsTrusted(&_TetherContractTC.CallOpts, arg0)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TetherContractTC *TetherContractTCCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TetherContractTC.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TetherContractTC *TetherContractTCSession) Name() (string, error) {
	return _TetherContractTC.Contract.Name(&_TetherContractTC.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TetherContractTC *TetherContractTCCallerSession) Name() (string, error) {
	return _TetherContractTC.Contract.Name(&_TetherContractTC.CallOpts)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_TetherContractTC *TetherContractTCCaller) Nonces(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TetherContractTC.contract.Call(opts, &out, "nonces", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_TetherContractTC *TetherContractTCSession) Nonces(owner common.Address) (*big.Int, error) {
	return _TetherContractTC.Contract.Nonces(&_TetherContractTC.CallOpts, owner)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_TetherContractTC *TetherContractTCCallerSession) Nonces(owner common.Address) (*big.Int, error) {
	return _TetherContractTC.Contract.Nonces(&_TetherContractTC.CallOpts, owner)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TetherContractTC *TetherContractTCCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TetherContractTC.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TetherContractTC *TetherContractTCSession) Owner() (common.Address, error) {
	return _TetherContractTC.Contract.Owner(&_TetherContractTC.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TetherContractTC *TetherContractTCCallerSession) Owner() (common.Address, error) {
	return _TetherContractTC.Contract.Owner(&_TetherContractTC.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TetherContractTC *TetherContractTCCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TetherContractTC.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TetherContractTC *TetherContractTCSession) Symbol() (string, error) {
	return _TetherContractTC.Contract.Symbol(&_TetherContractTC.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TetherContractTC *TetherContractTCCallerSession) Symbol() (string, error) {
	return _TetherContractTC.Contract.Symbol(&_TetherContractTC.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TetherContractTC *TetherContractTCCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TetherContractTC.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TetherContractTC *TetherContractTCSession) TotalSupply() (*big.Int, error) {
	return _TetherContractTC.Contract.TotalSupply(&_TetherContractTC.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TetherContractTC *TetherContractTCCallerSession) TotalSupply() (*big.Int, error) {
	return _TetherContractTC.Contract.TotalSupply(&_TetherContractTC.CallOpts)
}

// AddToBlockedList is a paid mutator transaction binding the contract method 0x3c7c9b90.
//
// Solidity: function addToBlockedList(address _user) returns()
func (_TetherContractTC *TetherContractTCTransactor) AddToBlockedList(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _TetherContractTC.contract.Transact(opts, "addToBlockedList", _user)
}

// AddToBlockedList is a paid mutator transaction binding the contract method 0x3c7c9b90.
//
// Solidity: function addToBlockedList(address _user) returns()
func (_TetherContractTC *TetherContractTCSession) AddToBlockedList(_user common.Address) (*types.Transaction, error) {
	return _TetherContractTC.Contract.AddToBlockedList(&_TetherContractTC.TransactOpts, _user)
}

// AddToBlockedList is a paid mutator transaction binding the contract method 0x3c7c9b90.
//
// Solidity: function addToBlockedList(address _user) returns()
func (_TetherContractTC *TetherContractTCTransactorSession) AddToBlockedList(_user common.Address) (*types.Transaction, error) {
	return _TetherContractTC.Contract.AddToBlockedList(&_TetherContractTC.TransactOpts, _user)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_TetherContractTC *TetherContractTCTransactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TetherContractTC.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_TetherContractTC *TetherContractTCSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TetherContractTC.Contract.Approve(&_TetherContractTC.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_TetherContractTC *TetherContractTCTransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TetherContractTC.Contract.Approve(&_TetherContractTC.TransactOpts, spender, amount)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_TetherContractTC *TetherContractTCTransactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _TetherContractTC.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_TetherContractTC *TetherContractTCSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _TetherContractTC.Contract.DecreaseAllowance(&_TetherContractTC.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_TetherContractTC *TetherContractTCTransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _TetherContractTC.Contract.DecreaseAllowance(&_TetherContractTC.TransactOpts, spender, subtractedValue)
}

// DestroyBlockedFunds is a paid mutator transaction binding the contract method 0x0e27a385.
//
// Solidity: function destroyBlockedFunds(address _blockedUser) returns()
func (_TetherContractTC *TetherContractTCTransactor) DestroyBlockedFunds(opts *bind.TransactOpts, _blockedUser common.Address) (*types.Transaction, error) {
	return _TetherContractTC.contract.Transact(opts, "destroyBlockedFunds", _blockedUser)
}

// DestroyBlockedFunds is a paid mutator transaction binding the contract method 0x0e27a385.
//
// Solidity: function destroyBlockedFunds(address _blockedUser) returns()
func (_TetherContractTC *TetherContractTCSession) DestroyBlockedFunds(_blockedUser common.Address) (*types.Transaction, error) {
	return _TetherContractTC.Contract.DestroyBlockedFunds(&_TetherContractTC.TransactOpts, _blockedUser)
}

// DestroyBlockedFunds is a paid mutator transaction binding the contract method 0x0e27a385.
//
// Solidity: function destroyBlockedFunds(address _blockedUser) returns()
func (_TetherContractTC *TetherContractTCTransactorSession) DestroyBlockedFunds(_blockedUser common.Address) (*types.Transaction, error) {
	return _TetherContractTC.Contract.DestroyBlockedFunds(&_TetherContractTC.TransactOpts, _blockedUser)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_TetherContractTC *TetherContractTCTransactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _TetherContractTC.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_TetherContractTC *TetherContractTCSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _TetherContractTC.Contract.IncreaseAllowance(&_TetherContractTC.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_TetherContractTC *TetherContractTCTransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _TetherContractTC.Contract.IncreaseAllowance(&_TetherContractTC.TransactOpts, spender, addedValue)
}

// Initialize is a paid mutator transaction binding the contract method 0x1624f6c6.
//
// Solidity: function initialize(string _name, string _symbol, uint8 _decimals) returns()
func (_TetherContractTC *TetherContractTCTransactor) Initialize(opts *bind.TransactOpts, _name string, _symbol string, _decimals uint8) (*types.Transaction, error) {
	return _TetherContractTC.contract.Transact(opts, "initialize", _name, _symbol, _decimals)
}

// Initialize is a paid mutator transaction binding the contract method 0x1624f6c6.
//
// Solidity: function initialize(string _name, string _symbol, uint8 _decimals) returns()
func (_TetherContractTC *TetherContractTCSession) Initialize(_name string, _symbol string, _decimals uint8) (*types.Transaction, error) {
	return _TetherContractTC.Contract.Initialize(&_TetherContractTC.TransactOpts, _name, _symbol, _decimals)
}

// Initialize is a paid mutator transaction binding the contract method 0x1624f6c6.
//
// Solidity: function initialize(string _name, string _symbol, uint8 _decimals) returns()
func (_TetherContractTC *TetherContractTCTransactorSession) Initialize(_name string, _symbol string, _decimals uint8) (*types.Transaction, error) {
	return _TetherContractTC.Contract.Initialize(&_TetherContractTC.TransactOpts, _name, _symbol, _decimals)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _destination, uint256 _amount) returns()
func (_TetherContractTC *TetherContractTCTransactor) Mint(opts *bind.TransactOpts, _destination common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TetherContractTC.contract.Transact(opts, "mint", _destination, _amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _destination, uint256 _amount) returns()
func (_TetherContractTC *TetherContractTCSession) Mint(_destination common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TetherContractTC.Contract.Mint(&_TetherContractTC.TransactOpts, _destination, _amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _destination, uint256 _amount) returns()
func (_TetherContractTC *TetherContractTCTransactorSession) Mint(_destination common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TetherContractTC.Contract.Mint(&_TetherContractTC.TransactOpts, _destination, _amount)
}

// MultiTransfer is a paid mutator transaction binding the contract method 0x1e89d545.
//
// Solidity: function multiTransfer(address[] _recipients, uint256[] _values) returns()
func (_TetherContractTC *TetherContractTCTransactor) MultiTransfer(opts *bind.TransactOpts, _recipients []common.Address, _values []*big.Int) (*types.Transaction, error) {
	return _TetherContractTC.contract.Transact(opts, "multiTransfer", _recipients, _values)
}

// MultiTransfer is a paid mutator transaction binding the contract method 0x1e89d545.
//
// Solidity: function multiTransfer(address[] _recipients, uint256[] _values) returns()
func (_TetherContractTC *TetherContractTCSession) MultiTransfer(_recipients []common.Address, _values []*big.Int) (*types.Transaction, error) {
	return _TetherContractTC.Contract.MultiTransfer(&_TetherContractTC.TransactOpts, _recipients, _values)
}

// MultiTransfer is a paid mutator transaction binding the contract method 0x1e89d545.
//
// Solidity: function multiTransfer(address[] _recipients, uint256[] _values) returns()
func (_TetherContractTC *TetherContractTCTransactorSession) MultiTransfer(_recipients []common.Address, _values []*big.Int) (*types.Transaction, error) {
	return _TetherContractTC.Contract.MultiTransfer(&_TetherContractTC.TransactOpts, _recipients, _values)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_TetherContractTC *TetherContractTCTransactor) Permit(opts *bind.TransactOpts, owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _TetherContractTC.contract.Transact(opts, "permit", owner, spender, value, deadline, v, r, s)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_TetherContractTC *TetherContractTCSession) Permit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _TetherContractTC.Contract.Permit(&_TetherContractTC.TransactOpts, owner, spender, value, deadline, v, r, s)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_TetherContractTC *TetherContractTCTransactorSession) Permit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _TetherContractTC.Contract.Permit(&_TetherContractTC.TransactOpts, owner, spender, value, deadline, v, r, s)
}

// Redeem is a paid mutator transaction binding the contract method 0xdb006a75.
//
// Solidity: function redeem(uint256 _amount) returns()
func (_TetherContractTC *TetherContractTCTransactor) Redeem(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _TetherContractTC.contract.Transact(opts, "redeem", _amount)
}

// Redeem is a paid mutator transaction binding the contract method 0xdb006a75.
//
// Solidity: function redeem(uint256 _amount) returns()
func (_TetherContractTC *TetherContractTCSession) Redeem(_amount *big.Int) (*types.Transaction, error) {
	return _TetherContractTC.Contract.Redeem(&_TetherContractTC.TransactOpts, _amount)
}

// Redeem is a paid mutator transaction binding the contract method 0xdb006a75.
//
// Solidity: function redeem(uint256 _amount) returns()
func (_TetherContractTC *TetherContractTCTransactorSession) Redeem(_amount *big.Int) (*types.Transaction, error) {
	return _TetherContractTC.Contract.Redeem(&_TetherContractTC.TransactOpts, _amount)
}

// RemoveFromBlockedList is a paid mutator transaction binding the contract method 0x1a14f449.
//
// Solidity: function removeFromBlockedList(address _user) returns()
func (_TetherContractTC *TetherContractTCTransactor) RemoveFromBlockedList(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _TetherContractTC.contract.Transact(opts, "removeFromBlockedList", _user)
}

// RemoveFromBlockedList is a paid mutator transaction binding the contract method 0x1a14f449.
//
// Solidity: function removeFromBlockedList(address _user) returns()
func (_TetherContractTC *TetherContractTCSession) RemoveFromBlockedList(_user common.Address) (*types.Transaction, error) {
	return _TetherContractTC.Contract.RemoveFromBlockedList(&_TetherContractTC.TransactOpts, _user)
}

// RemoveFromBlockedList is a paid mutator transaction binding the contract method 0x1a14f449.
//
// Solidity: function removeFromBlockedList(address _user) returns()
func (_TetherContractTC *TetherContractTCTransactorSession) RemoveFromBlockedList(_user common.Address) (*types.Transaction, error) {
	return _TetherContractTC.Contract.RemoveFromBlockedList(&_TetherContractTC.TransactOpts, _user)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TetherContractTC *TetherContractTCTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TetherContractTC.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TetherContractTC *TetherContractTCSession) RenounceOwnership() (*types.Transaction, error) {
	return _TetherContractTC.Contract.RenounceOwnership(&_TetherContractTC.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TetherContractTC *TetherContractTCTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _TetherContractTC.Contract.RenounceOwnership(&_TetherContractTC.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_TetherContractTC *TetherContractTCTransactor) Transfer(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TetherContractTC.contract.Transact(opts, "transfer", to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_TetherContractTC *TetherContractTCSession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TetherContractTC.Contract.Transfer(&_TetherContractTC.TransactOpts, to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_TetherContractTC *TetherContractTCTransactorSession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TetherContractTC.Contract.Transfer(&_TetherContractTC.TransactOpts, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _sender, address _recipient, uint256 _amount) returns(bool)
func (_TetherContractTC *TetherContractTCTransactor) TransferFrom(opts *bind.TransactOpts, _sender common.Address, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TetherContractTC.contract.Transact(opts, "transferFrom", _sender, _recipient, _amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _sender, address _recipient, uint256 _amount) returns(bool)
func (_TetherContractTC *TetherContractTCSession) TransferFrom(_sender common.Address, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TetherContractTC.Contract.TransferFrom(&_TetherContractTC.TransactOpts, _sender, _recipient, _amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _sender, address _recipient, uint256 _amount) returns(bool)
func (_TetherContractTC *TetherContractTCTransactorSession) TransferFrom(_sender common.Address, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TetherContractTC.Contract.TransferFrom(&_TetherContractTC.TransactOpts, _sender, _recipient, _amount)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TetherContractTC *TetherContractTCTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TetherContractTC.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TetherContractTC *TetherContractTCSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TetherContractTC.Contract.TransferOwnership(&_TetherContractTC.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TetherContractTC *TetherContractTCTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TetherContractTC.Contract.TransferOwnership(&_TetherContractTC.TransactOpts, newOwner)
}

// TetherContractTCApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the TetherContractTC contract.
type TetherContractTCApprovalIterator struct {
	Event *TetherContractTCApproval // Event containing the contract specifics and raw log

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
func (it *TetherContractTCApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TetherContractTCApproval)
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
		it.Event = new(TetherContractTCApproval)
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
func (it *TetherContractTCApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TetherContractTCApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TetherContractTCApproval represents a Approval event raised by the TetherContractTC contract.
type TetherContractTCApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TetherContractTC *TetherContractTCFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*TetherContractTCApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _TetherContractTC.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &TetherContractTCApprovalIterator{contract: _TetherContractTC.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TetherContractTC *TetherContractTCFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *TetherContractTCApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _TetherContractTC.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TetherContractTCApproval)
				if err := _TetherContractTC.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TetherContractTC *TetherContractTCFilterer) ParseApproval(log types.Log) (*TetherContractTCApproval, error) {
	event := new(TetherContractTCApproval)
	if err := _TetherContractTC.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TetherContractTCBlockPlacedIterator is returned from FilterBlockPlaced and is used to iterate over the raw logs and unpacked data for BlockPlaced events raised by the TetherContractTC contract.
type TetherContractTCBlockPlacedIterator struct {
	Event *TetherContractTCBlockPlaced // Event containing the contract specifics and raw log

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
func (it *TetherContractTCBlockPlacedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TetherContractTCBlockPlaced)
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
		it.Event = new(TetherContractTCBlockPlaced)
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
func (it *TetherContractTCBlockPlacedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TetherContractTCBlockPlacedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TetherContractTCBlockPlaced represents a BlockPlaced event raised by the TetherContractTC contract.
type TetherContractTCBlockPlaced struct {
	User common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterBlockPlaced is a free log retrieval operation binding the contract event 0x406bbf2d8d145125adf1198d2cf8a67c66cc4bb0ab01c37dccd4f7c0aae1e7c7.
//
// Solidity: event BlockPlaced(address indexed _user)
func (_TetherContractTC *TetherContractTCFilterer) FilterBlockPlaced(opts *bind.FilterOpts, _user []common.Address) (*TetherContractTCBlockPlacedIterator, error) {

	var _userRule []interface{}
	for _, _userItem := range _user {
		_userRule = append(_userRule, _userItem)
	}

	logs, sub, err := _TetherContractTC.contract.FilterLogs(opts, "BlockPlaced", _userRule)
	if err != nil {
		return nil, err
	}
	return &TetherContractTCBlockPlacedIterator{contract: _TetherContractTC.contract, event: "BlockPlaced", logs: logs, sub: sub}, nil
}

// WatchBlockPlaced is a free log subscription operation binding the contract event 0x406bbf2d8d145125adf1198d2cf8a67c66cc4bb0ab01c37dccd4f7c0aae1e7c7.
//
// Solidity: event BlockPlaced(address indexed _user)
func (_TetherContractTC *TetherContractTCFilterer) WatchBlockPlaced(opts *bind.WatchOpts, sink chan<- *TetherContractTCBlockPlaced, _user []common.Address) (event.Subscription, error) {

	var _userRule []interface{}
	for _, _userItem := range _user {
		_userRule = append(_userRule, _userItem)
	}

	logs, sub, err := _TetherContractTC.contract.WatchLogs(opts, "BlockPlaced", _userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TetherContractTCBlockPlaced)
				if err := _TetherContractTC.contract.UnpackLog(event, "BlockPlaced", log); err != nil {
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

// ParseBlockPlaced is a log parse operation binding the contract event 0x406bbf2d8d145125adf1198d2cf8a67c66cc4bb0ab01c37dccd4f7c0aae1e7c7.
//
// Solidity: event BlockPlaced(address indexed _user)
func (_TetherContractTC *TetherContractTCFilterer) ParseBlockPlaced(log types.Log) (*TetherContractTCBlockPlaced, error) {
	event := new(TetherContractTCBlockPlaced)
	if err := _TetherContractTC.contract.UnpackLog(event, "BlockPlaced", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TetherContractTCBlockReleasedIterator is returned from FilterBlockReleased and is used to iterate over the raw logs and unpacked data for BlockReleased events raised by the TetherContractTC contract.
type TetherContractTCBlockReleasedIterator struct {
	Event *TetherContractTCBlockReleased // Event containing the contract specifics and raw log

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
func (it *TetherContractTCBlockReleasedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TetherContractTCBlockReleased)
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
		it.Event = new(TetherContractTCBlockReleased)
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
func (it *TetherContractTCBlockReleasedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TetherContractTCBlockReleasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TetherContractTCBlockReleased represents a BlockReleased event raised by the TetherContractTC contract.
type TetherContractTCBlockReleased struct {
	User common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterBlockReleased is a free log retrieval operation binding the contract event 0x665918c9e02eb2fd85acca3969cb054fc84c138e60ec4af22ab6ef2fd4c93c27.
//
// Solidity: event BlockReleased(address indexed _user)
func (_TetherContractTC *TetherContractTCFilterer) FilterBlockReleased(opts *bind.FilterOpts, _user []common.Address) (*TetherContractTCBlockReleasedIterator, error) {

	var _userRule []interface{}
	for _, _userItem := range _user {
		_userRule = append(_userRule, _userItem)
	}

	logs, sub, err := _TetherContractTC.contract.FilterLogs(opts, "BlockReleased", _userRule)
	if err != nil {
		return nil, err
	}
	return &TetherContractTCBlockReleasedIterator{contract: _TetherContractTC.contract, event: "BlockReleased", logs: logs, sub: sub}, nil
}

// WatchBlockReleased is a free log subscription operation binding the contract event 0x665918c9e02eb2fd85acca3969cb054fc84c138e60ec4af22ab6ef2fd4c93c27.
//
// Solidity: event BlockReleased(address indexed _user)
func (_TetherContractTC *TetherContractTCFilterer) WatchBlockReleased(opts *bind.WatchOpts, sink chan<- *TetherContractTCBlockReleased, _user []common.Address) (event.Subscription, error) {

	var _userRule []interface{}
	for _, _userItem := range _user {
		_userRule = append(_userRule, _userItem)
	}

	logs, sub, err := _TetherContractTC.contract.WatchLogs(opts, "BlockReleased", _userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TetherContractTCBlockReleased)
				if err := _TetherContractTC.contract.UnpackLog(event, "BlockReleased", log); err != nil {
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

// ParseBlockReleased is a log parse operation binding the contract event 0x665918c9e02eb2fd85acca3969cb054fc84c138e60ec4af22ab6ef2fd4c93c27.
//
// Solidity: event BlockReleased(address indexed _user)
func (_TetherContractTC *TetherContractTCFilterer) ParseBlockReleased(log types.Log) (*TetherContractTCBlockReleased, error) {
	event := new(TetherContractTCBlockReleased)
	if err := _TetherContractTC.contract.UnpackLog(event, "BlockReleased", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TetherContractTCDestroyedBlockedFundsIterator is returned from FilterDestroyedBlockedFunds and is used to iterate over the raw logs and unpacked data for DestroyedBlockedFunds events raised by the TetherContractTC contract.
type TetherContractTCDestroyedBlockedFundsIterator struct {
	Event *TetherContractTCDestroyedBlockedFunds // Event containing the contract specifics and raw log

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
func (it *TetherContractTCDestroyedBlockedFundsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TetherContractTCDestroyedBlockedFunds)
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
		it.Event = new(TetherContractTCDestroyedBlockedFunds)
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
func (it *TetherContractTCDestroyedBlockedFundsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TetherContractTCDestroyedBlockedFundsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TetherContractTCDestroyedBlockedFunds represents a DestroyedBlockedFunds event raised by the TetherContractTC contract.
type TetherContractTCDestroyedBlockedFunds struct {
	BlockedUser common.Address
	Balance     *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterDestroyedBlockedFunds is a free log retrieval operation binding the contract event 0x6a2859ae7902313752498feb80a014e6e7275fe964c79aa965db815db1c7f1e9.
//
// Solidity: event DestroyedBlockedFunds(address indexed _blockedUser, uint256 _balance)
func (_TetherContractTC *TetherContractTCFilterer) FilterDestroyedBlockedFunds(opts *bind.FilterOpts, _blockedUser []common.Address) (*TetherContractTCDestroyedBlockedFundsIterator, error) {

	var _blockedUserRule []interface{}
	for _, _blockedUserItem := range _blockedUser {
		_blockedUserRule = append(_blockedUserRule, _blockedUserItem)
	}

	logs, sub, err := _TetherContractTC.contract.FilterLogs(opts, "DestroyedBlockedFunds", _blockedUserRule)
	if err != nil {
		return nil, err
	}
	return &TetherContractTCDestroyedBlockedFundsIterator{contract: _TetherContractTC.contract, event: "DestroyedBlockedFunds", logs: logs, sub: sub}, nil
}

// WatchDestroyedBlockedFunds is a free log subscription operation binding the contract event 0x6a2859ae7902313752498feb80a014e6e7275fe964c79aa965db815db1c7f1e9.
//
// Solidity: event DestroyedBlockedFunds(address indexed _blockedUser, uint256 _balance)
func (_TetherContractTC *TetherContractTCFilterer) WatchDestroyedBlockedFunds(opts *bind.WatchOpts, sink chan<- *TetherContractTCDestroyedBlockedFunds, _blockedUser []common.Address) (event.Subscription, error) {

	var _blockedUserRule []interface{}
	for _, _blockedUserItem := range _blockedUser {
		_blockedUserRule = append(_blockedUserRule, _blockedUserItem)
	}

	logs, sub, err := _TetherContractTC.contract.WatchLogs(opts, "DestroyedBlockedFunds", _blockedUserRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TetherContractTCDestroyedBlockedFunds)
				if err := _TetherContractTC.contract.UnpackLog(event, "DestroyedBlockedFunds", log); err != nil {
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

// ParseDestroyedBlockedFunds is a log parse operation binding the contract event 0x6a2859ae7902313752498feb80a014e6e7275fe964c79aa965db815db1c7f1e9.
//
// Solidity: event DestroyedBlockedFunds(address indexed _blockedUser, uint256 _balance)
func (_TetherContractTC *TetherContractTCFilterer) ParseDestroyedBlockedFunds(log types.Log) (*TetherContractTCDestroyedBlockedFunds, error) {
	event := new(TetherContractTCDestroyedBlockedFunds)
	if err := _TetherContractTC.contract.UnpackLog(event, "DestroyedBlockedFunds", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TetherContractTCInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the TetherContractTC contract.
type TetherContractTCInitializedIterator struct {
	Event *TetherContractTCInitialized // Event containing the contract specifics and raw log

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
func (it *TetherContractTCInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TetherContractTCInitialized)
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
		it.Event = new(TetherContractTCInitialized)
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
func (it *TetherContractTCInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TetherContractTCInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TetherContractTCInitialized represents a Initialized event raised by the TetherContractTC contract.
type TetherContractTCInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_TetherContractTC *TetherContractTCFilterer) FilterInitialized(opts *bind.FilterOpts) (*TetherContractTCInitializedIterator, error) {

	logs, sub, err := _TetherContractTC.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &TetherContractTCInitializedIterator{contract: _TetherContractTC.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_TetherContractTC *TetherContractTCFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *TetherContractTCInitialized) (event.Subscription, error) {

	logs, sub, err := _TetherContractTC.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TetherContractTCInitialized)
				if err := _TetherContractTC.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_TetherContractTC *TetherContractTCFilterer) ParseInitialized(log types.Log) (*TetherContractTCInitialized, error) {
	event := new(TetherContractTCInitialized)
	if err := _TetherContractTC.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TetherContractTCMintIterator is returned from FilterMint and is used to iterate over the raw logs and unpacked data for Mint events raised by the TetherContractTC contract.
type TetherContractTCMintIterator struct {
	Event *TetherContractTCMint // Event containing the contract specifics and raw log

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
func (it *TetherContractTCMintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TetherContractTCMint)
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
		it.Event = new(TetherContractTCMint)
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
func (it *TetherContractTCMintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TetherContractTCMintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TetherContractTCMint represents a Mint event raised by the TetherContractTC contract.
type TetherContractTCMint struct {
	Destination common.Address
	Amount      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterMint is a free log retrieval operation binding the contract event 0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885.
//
// Solidity: event Mint(address indexed _destination, uint256 _amount)
func (_TetherContractTC *TetherContractTCFilterer) FilterMint(opts *bind.FilterOpts, _destination []common.Address) (*TetherContractTCMintIterator, error) {

	var _destinationRule []interface{}
	for _, _destinationItem := range _destination {
		_destinationRule = append(_destinationRule, _destinationItem)
	}

	logs, sub, err := _TetherContractTC.contract.FilterLogs(opts, "Mint", _destinationRule)
	if err != nil {
		return nil, err
	}
	return &TetherContractTCMintIterator{contract: _TetherContractTC.contract, event: "Mint", logs: logs, sub: sub}, nil
}

// WatchMint is a free log subscription operation binding the contract event 0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885.
//
// Solidity: event Mint(address indexed _destination, uint256 _amount)
func (_TetherContractTC *TetherContractTCFilterer) WatchMint(opts *bind.WatchOpts, sink chan<- *TetherContractTCMint, _destination []common.Address) (event.Subscription, error) {

	var _destinationRule []interface{}
	for _, _destinationItem := range _destination {
		_destinationRule = append(_destinationRule, _destinationItem)
	}

	logs, sub, err := _TetherContractTC.contract.WatchLogs(opts, "Mint", _destinationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TetherContractTCMint)
				if err := _TetherContractTC.contract.UnpackLog(event, "Mint", log); err != nil {
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

// ParseMint is a log parse operation binding the contract event 0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885.
//
// Solidity: event Mint(address indexed _destination, uint256 _amount)
func (_TetherContractTC *TetherContractTCFilterer) ParseMint(log types.Log) (*TetherContractTCMint, error) {
	event := new(TetherContractTCMint)
	if err := _TetherContractTC.contract.UnpackLog(event, "Mint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TetherContractTCOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TetherContractTC contract.
type TetherContractTCOwnershipTransferredIterator struct {
	Event *TetherContractTCOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *TetherContractTCOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TetherContractTCOwnershipTransferred)
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
		it.Event = new(TetherContractTCOwnershipTransferred)
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
func (it *TetherContractTCOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TetherContractTCOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TetherContractTCOwnershipTransferred represents a OwnershipTransferred event raised by the TetherContractTC contract.
type TetherContractTCOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TetherContractTC *TetherContractTCFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TetherContractTCOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TetherContractTC.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TetherContractTCOwnershipTransferredIterator{contract: _TetherContractTC.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TetherContractTC *TetherContractTCFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TetherContractTCOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TetherContractTC.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TetherContractTCOwnershipTransferred)
				if err := _TetherContractTC.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TetherContractTC *TetherContractTCFilterer) ParseOwnershipTransferred(log types.Log) (*TetherContractTCOwnershipTransferred, error) {
	event := new(TetherContractTCOwnershipTransferred)
	if err := _TetherContractTC.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TetherContractTCRedeemIterator is returned from FilterRedeem and is used to iterate over the raw logs and unpacked data for Redeem events raised by the TetherContractTC contract.
type TetherContractTCRedeemIterator struct {
	Event *TetherContractTCRedeem // Event containing the contract specifics and raw log

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
func (it *TetherContractTCRedeemIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TetherContractTCRedeem)
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
		it.Event = new(TetherContractTCRedeem)
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
func (it *TetherContractTCRedeemIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TetherContractTCRedeemIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TetherContractTCRedeem represents a Redeem event raised by the TetherContractTC contract.
type TetherContractTCRedeem struct {
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRedeem is a free log retrieval operation binding the contract event 0x702d5967f45f6513a38ffc42d6ba9bf230bd40e8f53b16363c7eb4fd2deb9a44.
//
// Solidity: event Redeem(uint256 _amount)
func (_TetherContractTC *TetherContractTCFilterer) FilterRedeem(opts *bind.FilterOpts) (*TetherContractTCRedeemIterator, error) {

	logs, sub, err := _TetherContractTC.contract.FilterLogs(opts, "Redeem")
	if err != nil {
		return nil, err
	}
	return &TetherContractTCRedeemIterator{contract: _TetherContractTC.contract, event: "Redeem", logs: logs, sub: sub}, nil
}

// WatchRedeem is a free log subscription operation binding the contract event 0x702d5967f45f6513a38ffc42d6ba9bf230bd40e8f53b16363c7eb4fd2deb9a44.
//
// Solidity: event Redeem(uint256 _amount)
func (_TetherContractTC *TetherContractTCFilterer) WatchRedeem(opts *bind.WatchOpts, sink chan<- *TetherContractTCRedeem) (event.Subscription, error) {

	logs, sub, err := _TetherContractTC.contract.WatchLogs(opts, "Redeem")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TetherContractTCRedeem)
				if err := _TetherContractTC.contract.UnpackLog(event, "Redeem", log); err != nil {
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

// ParseRedeem is a log parse operation binding the contract event 0x702d5967f45f6513a38ffc42d6ba9bf230bd40e8f53b16363c7eb4fd2deb9a44.
//
// Solidity: event Redeem(uint256 _amount)
func (_TetherContractTC *TetherContractTCFilterer) ParseRedeem(log types.Log) (*TetherContractTCRedeem, error) {
	event := new(TetherContractTCRedeem)
	if err := _TetherContractTC.contract.UnpackLog(event, "Redeem", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TetherContractTCTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the TetherContractTC contract.
type TetherContractTCTransferIterator struct {
	Event *TetherContractTCTransfer // Event containing the contract specifics and raw log

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
func (it *TetherContractTCTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TetherContractTCTransfer)
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
		it.Event = new(TetherContractTCTransfer)
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
func (it *TetherContractTCTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TetherContractTCTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TetherContractTCTransfer represents a Transfer event raised by the TetherContractTC contract.
type TetherContractTCTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TetherContractTC *TetherContractTCFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TetherContractTCTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TetherContractTC.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &TetherContractTCTransferIterator{contract: _TetherContractTC.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TetherContractTC *TetherContractTCFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *TetherContractTCTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TetherContractTC.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TetherContractTCTransfer)
				if err := _TetherContractTC.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TetherContractTC *TetherContractTCFilterer) ParseTransfer(log types.Log) (*TetherContractTCTransfer, error) {
	event := new(TetherContractTCTransfer)
	if err := _TetherContractTC.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
