package testcase

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	kaia "github.com/kaiachain/kaia"
	"github.com/kaiachain/kaia-load-tester/klayslave/account"
	"github.com/kaiachain/kaia/accounts/abi"
	"github.com/kaiachain/kaia/client"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/myzhan/boomer"
	"github.com/tidwall/gjson"
)

var (
	readApiCallMutex  sync.Mutex
	latestBlockNumber = big.NewInt(0)
	count             uint64

	retValOfCall      = big.NewInt(4)
	retValOfStorageAt = big.NewInt(4)
)

func sendBoomerEvent(tcName string, logString string, elapsed int64, err error, endpoint string) {
	if err == nil {
		boomer.Events.Publish("request_success", "http", tcName+" to "+endpoint, elapsed, int64(10))
	} else {
		boomer.Events.Publish("request_failure", "http", tcName+" to "+endpoint, elapsed, err.Error())
	}
}

func getRandomBlockNumber(cli *client.KaiaClient, ctx context.Context) *big.Int {
	readApiCallMutex.Lock()
	defer readApiCallMutex.Unlock()

	count %= 10000000
	if count%10000 == 0 {
		bn, err := cli.BlockNumber(ctx)
		if err != nil {
			log.Printf("Failed to update the current block number. err=%s\n", err)
		} else {
			log.Printf("Update the current block number. blockNumber=0x%s\n", bn.Text(16))
			latestBlockNumber.Set(bn)
		}
	}
	count += 1

	return big.NewInt(0).Rand(rand.New(rand.NewSource(time.Now().UnixNano())), latestBlockNumber)
}

// RunGasPrice creates a closure for gas price test case
func RunGasPrice(config *TCConfig) func() {
	return func() {
		ctx := context.Background()
		cli := config.CliPool.Alloc().(*client.KaiaClient)
		defer config.CliPool.Free(cli)

		start := boomer.Now()
		_, err := cli.SuggestGasPrice(ctx)
		elapsed := boomer.Now() - start
		sendBoomerEvent("readGasPrice", "Failed to call klay_gasPrice", elapsed, err, config.EndPoint)
	}
}

// RunBlockNumber creates a closure for block number test case
func RunBlockNumber(config *TCConfig) func() {
	return func() {
		ctx := context.Background()
		cli := config.CliPool.Alloc().(*client.KaiaClient)
		defer config.CliPool.Free(cli)

		start := boomer.Now()

		bn, err := cli.BlockNumber(ctx)
		if err == nil && bn.Cmp(big.NewInt(0)) != 1 {
			err = errors.New("wrong block number: 0x" + bn.Text(16) + ", answer: smaller than 0")
		}

		elapsed := boomer.Now() - start
		sendBoomerEvent("readBlockNumber", "Failed to call klay_blockNumber", elapsed, err, config.EndPoint)
	}
}

// RunGetBlockByNumber creates a closure for get block by number test case
func RunGetBlockByNumber(config *TCConfig) func() {
	return func() {
		ctx := context.Background()
		cli := config.CliPool.Alloc().(*client.KaiaClient)
		defer config.CliPool.Free(cli)

		ansBN := getRandomBlockNumber(cli, ctx)
		start := boomer.Now()

		block, err := cli.BlockByNumber(ctx, ansBN) // read the random block
		if err == nil && block.Header().Number.Cmp(ansBN) != 0 {
			err = errors.New("wrong block: 0x" + block.Header().Number.Text(16) + ", answer: 0x" + ansBN.Text(16))
		}

		elapsed := boomer.Now() - start
		sendBoomerEvent("readGetBlockByNumber", "Failed to call klay_getBlockByNumber", elapsed, err, config.EndPoint)
	}
}

// RunGetAccount creates a closure for get account test case
func RunGetAccount(config *TCConfig) func() {
	return func() {
		ctx := context.Background()
		rpcCli := config.RpcCliPool.Alloc().(*rpc.Client)
		defer config.RpcCliPool.Free(rpcCli)

		fromAccount := config.AccGrp.GetAccountRandomly()
		start := boomer.Now()

		var j json.RawMessage
		err := rpcCli.CallContext(ctx, &j, "klay_getAccount", fromAccount.GetAddress(), "latest")
		if err == nil {
			ret := gjson.Get(string(j), "accType").String()
			if ret != "1" {
				err = errors.New("wrong account type: " + ret + ", answer: 1")
			}
		}

		elapsed := boomer.Now() - start
		sendBoomerEvent("readGetAccount", "Failed to call klay_getAccount", elapsed, err, config.EndPoint)
	}
}

// RunGetBlockWithConsensusInfoByNumber creates a closure for get block with consensus info by number test case
func RunGetBlockWithConsensusInfoByNumber(config *TCConfig) func() {
	return func() {
		ctx := context.Background()
		cli := config.CliPool.Alloc().(*client.KaiaClient)
		defer config.CliPool.Free(cli)
		rpcCli := config.RpcCliPool.Alloc().(*rpc.Client)
		defer config.RpcCliPool.Free(rpcCli)

		ansBN := getRandomBlockNumber(cli, ctx)
		start := boomer.Now()

		var j json.RawMessage
		err := rpcCli.CallContext(ctx, &j, "klay_getBlockWithConsensusInfoByNumber", "0x"+ansBN.Text(16))
		if err == nil {
			ret := gjson.Get(string(j), "number").String()
			if !strings.Contains(ret, "0x"+ansBN.Text(16)) {
				err = errors.New("wrong block: " + ret + ", answer: " + "0x" + ansBN.Text(16))
			}
		}

		elapsed := boomer.Now() - start
		sendBoomerEvent("readGetBlockWithConsensusInfoByNumber",
			"Failed to call klay_GetBlockWithConsensusInfoByNumber", elapsed, err, config.EndPoint)
	}
}

// RunGetStorageAt creates a closure for get storage at test case
func RunGetStorageAt(config *TCConfig) func() {
	return func() {
		ctx := context.Background()
		cli := config.CliPool.Alloc().(*client.KaiaClient)
		defer config.CliPool.Free(cli)

		smartContractAccount := config.SmartContractAccounts[account.ContractReadApiCallContract]
		contractAddr := smartContractAccount.GetAddress()
		start := boomer.Now()
		ret, err := cli.StorageAt(ctx, contractAddr, common.Hash{}, nil)
		elapsed := boomer.Now() - start

		if err == nil && new(big.Int).SetBytes(ret).Cmp(retValOfStorageAt) != 0 {
			err = errors.New("wrong storage value: " + string(ret) + ", answer: " + retValOfStorageAt.String())
		}
		sendBoomerEvent("readGetStorageAt", "Failure to call klay_getStorageAt", elapsed, err, config.EndPoint)
	}
}

// RunCall creates a closure for call test case
func RunCall(config *TCConfig) func() {
	return func() {
		cli := config.CliPool.Alloc().(*client.KaiaClient)
		defer config.CliPool.Free(cli)

		fromAccount := config.AccGrp.GetAccountRandomly().GetAddress()
		contractAddr := config.SmartContractAccounts[account.ContractReadApiCallContract].GetAddress()
		data := account.TestContractInfos[account.ContractReadApiCallContract].GenData(fromAccount, big.NewInt(0))

		callMsg := kaia.CallMsg{
			From: fromAccount,
			To:   &contractAddr,
			Data: data,
		}

		start := boomer.Now()
		result, err := cli.CallContract(context.Background(), callMsg, nil)
		elapsed := boomer.Now() - start

		if err == nil {
			// Parse the result using the same ABI
			abiStr := `[{"inputs":[],"name":"get","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`
			parsedABI, err := abi.JSON(strings.NewReader(abiStr))
			if err == nil {
				var ret *big.Int
				err = parsedABI.UnpackIntoInterface(&ret, "get", result)
				if err == nil && ret.Cmp(retValOfCall) != 0 {
					err = errors.New("wrong call: " + ret.String() + ", answer: " + retValOfCall.String())
				}
			}
		}
		sendBoomerEvent("readCall", "Failed to call klay_call", elapsed, err, config.EndPoint)
	}
}

// RunEstimateGas creates a closure for estimate gas test case
func RunEstimateGas(config *TCConfig) func() {
	return func() {
		ctx := context.Background()
		cli := config.CliPool.Alloc().(*client.KaiaClient)
		defer config.CliPool.Free(cli)

		fromAccount := config.AccGrp.GetAccountRandomly().GetAddress()
		contractAddr := config.SmartContractAccounts[account.ContractReadApiCallContract].GetAddress()
		data := account.TestContractInfos[account.ContractReadApiCallContract].GenData(fromAccount, big.NewInt(1))

		callMsg := kaia.CallMsg{
			From:     fromAccount,
			To:       &contractAddr,
			Gas:      1100000,
			GasPrice: big.NewInt(75000000000), // Default gas price: 25 Gwei
			Value:    big.NewInt(0),
			Data:     data,
		}
		start := boomer.Now()
		ret, err := cli.EstimateGas(ctx, callMsg)
		elapsed := boomer.Now() - start

		if err == nil && ret == 0 {
			err = errors.New("wrong estimate gas: " + strconv.Itoa(int(ret)))
		}
		sendBoomerEvent("readEstimateGas", "Failed to call klay_estimateGas", elapsed, err, config.EndPoint)
	}
}
