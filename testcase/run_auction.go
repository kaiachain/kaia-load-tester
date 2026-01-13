package testcase

import (
	"math/big"
	"math/rand"

	"github.com/kaiachain/kaia-load-tester/klayslave/account"
	"github.com/kaiachain/kaia/client"
	"github.com/kaiachain/kaia/common"
	"github.com/myzhan/boomer"
)

// AuctionTxFunc represents an auction transaction function signature
type AuctionTxFunc = func(*account.Account, *client.KaiaClient, *account.Account, *account.Account, string) (common.Hash, common.Hash, *big.Int, error)

// RunBaseWithAuction creates a closure that executes an auction test case with common logic
func RunBaseWithAuction(config *TCConfig, auctionTxFunc AuctionTxFunc) func() {
	return func() {
		cli := config.CliPool.Alloc().(*client.KaiaClient)
		defer config.CliPool.Free(cli)

		// Use round robin to avoid the same account being used too often
		from := config.AccGrp.GetAccountRoundRobin()

		auctionEntryPoint := config.SmartContractAccounts[account.ContractAuctionEntryPoint]
		targetContract := config.SmartContractAccounts[account.ContractCounterForTestAuction]

		// Select a targetTxType randomly from the list
		targetTxTypeKey := config.AuctionTargetTxTypeList[rand.Int()%len(config.AuctionTargetTxTypeList)]

		start := boomer.Now()
		_, _, _, err := auctionTxFunc(from, cli, auctionEntryPoint, targetContract, targetTxTypeKey)
		elapsed := boomer.Now() - start

		if err == nil {
			boomer.Events.Publish("request_success", "http", config.Name+" to "+config.EndPoint, elapsed, int64(10))
		} else {
			boomer.Events.Publish("request_failure", "http", config.Name+" to "+config.EndPoint, elapsed, err.Error())
		}
	}
}

func RunAuctionBidTC(config *TCConfig) func() {
	auctionTxFunc := func(from *account.Account, cli *client.KaiaClient, auctionEntryPoint, targetContract *account.Account, targetTxTypeKey string) (common.Hash, common.Hash, *big.Int, error) {
		return from.AuctionBid(cli, auctionEntryPoint, targetContract, targetTxTypeKey)
	}
	return RunBaseWithAuction(config, auctionTxFunc)
}

func RunAuctionRevertedBidTC(config *TCConfig) func() {
	auctionTxFunc := func(from *account.Account, cli *client.KaiaClient, auctionEntryPoint, targetContract *account.Account, targetTxTypeKey string) (common.Hash, common.Hash, *big.Int, error) {
		return from.AuctionRevertedBid(cli, auctionEntryPoint, targetContract, targetTxTypeKey)
	}
	return RunBaseWithAuction(config, auctionTxFunc)
}
