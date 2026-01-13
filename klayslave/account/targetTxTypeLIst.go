package account

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/kaiachain/kaia/accounts/abi"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/client"
	"github.com/kaiachain/kaia/common"
)

type TargetTxType struct {
	Description string
	GenerateTx  func(c *client.KaiaClient, account, tmpAccount *Account, nonce uint64, suggestedGasPrice *big.Int) *types.Transaction
	PreSendBid  func(c *client.KaiaClient, account, tmpAccount *Account, nonce uint64, suggestedGasPrice *big.Int) error
	PostSendBid func(c *client.KaiaClient, account, tmpAccount *Account, nonce uint64, suggestedGasPrice *big.Int, blockNumber *big.Int)
}

// TargetTxTypeList defines the list of auction target tx types.
var TargetTxTypeList = map[string]*TargetTxType{
	"VT": {
		Description: "VT creates a ValueTransfer that sends a 1kei to itself.",
		GenerateTx: func(c *client.KaiaClient, account, _ *Account, nonce uint64, suggestedGasPrice *big.Int) *types.Transaction {
			signer := types.NewEIP155Signer(chainID)
			tx, err := types.NewTransactionWithMap(types.TxTypeValueTransfer, map[types.TxValueKeyType]interface{}{
				types.TxValueKeyNonce:    nonce,
				types.TxValueKeyTo:       account.address,
				types.TxValueKeyAmount:   common.Big1,
				types.TxValueKeyGasLimit: uint64(100000),
				types.TxValueKeyGasPrice: suggestedGasPrice,
				types.TxValueKeyFrom:     account.address,
			})
			if err != nil {
				return nil
			}

			err = tx.SignWithKeys(signer, account.privateKey)
			if err != nil {
				return nil
			}
			return tx
		},
		PreSendBid: func(c *client.KaiaClient, account, _ *Account, nonce uint64, suggestedGasPrice *big.Int) error {
			return nil
		},
		PostSendBid: func(c *client.KaiaClient, account, _ *Account, nonce uint64, suggestedGasPrice *big.Int, blockNumber *big.Int) {
			account.nonce++
			account.updateLastBlocknumSentTx(blockNumber.Uint64())
		},
	},
	"SC": {
		Description: "SC creates a SmartContractCall that calls the counter contract.",
		GenerateTx: func(c *client.KaiaClient, account, _ *Account, nonce uint64, suggestedGasPrice *big.Int) *types.Transaction {
			signer := types.NewEIP155Signer(chainID)
			tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractExecution, map[types.TxValueKeyType]interface{}{
				types.TxValueKeyNonce:    nonce,
				types.TxValueKeyTo:       TestContractInfos[ContractCounterForTestAuction].GetAddress(c, CounterForTestAuctionDeployer),
				types.TxValueKeyData:     TestContractInfos[ContractCounterForTestAuction].GenData(common.Address{}, common.Big1), // 1 means calling incForSC()
				types.TxValueKeyAmount:   common.Big0,
				types.TxValueKeyGasLimit: uint64(5000000),
				types.TxValueKeyGasPrice: suggestedGasPrice,
				types.TxValueKeyFrom:     account.address,
			})
			if err != nil {
				return nil
			}

			err = tx.SignWithKeys(signer, account.privateKey)
			if err != nil {
				return nil
			}
			return tx
		},
		PreSendBid: func(c *client.KaiaClient, account, _ *Account, nonce uint64, suggestedGasPrice *big.Int) error {
			return nil
		},
		PostSendBid: func(c *client.KaiaClient, account, _ *Account, nonce uint64, suggestedGasPrice *big.Int, blockNumber *big.Int) {
			account.nonce++
			account.updateLastBlocknumSentTx(blockNumber.Uint64())
		},
	},
	"rSC": {
		Description: "rSC creates a reverted SmartContractCall fails when incorrect data is entered.",
		GenerateTx: func(c *client.KaiaClient, account, _ *Account, nonce uint64, suggestedGasPrice *big.Int) *types.Transaction {
			signer := types.NewEIP155Signer(chainID)
			tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractExecution, map[types.TxValueKeyType]interface{}{
				types.TxValueKeyNonce:    nonce,
				types.TxValueKeyTo:       TestContractInfos[ContractCounterForTestAuction].GetAddress(c, CounterForTestAuctionDeployer),
				types.TxValueKeyData:     TestContractInfos[ContractCounterForTestAuction].GenData(common.Address{}, common.Big2), // 2 means calling intendedRevert()
				types.TxValueKeyAmount:   common.Big0,
				types.TxValueKeyGasLimit: uint64(5000000),
				types.TxValueKeyGasPrice: suggestedGasPrice,
				types.TxValueKeyFrom:     account.address,
			})
			if err != nil {
				return nil
			}

			err = tx.SignWithKeys(signer, account.privateKey)
			if err != nil {
				return nil
			}
			return tx
		},
		PreSendBid: func(c *client.KaiaClient, account, _ *Account, nonce uint64, suggestedGasPrice *big.Int) error {
			return nil
		},
		PostSendBid: func(c *client.KaiaClient, account, _ *Account, nonce uint64, suggestedGasPrice *big.Int, blockNumber *big.Int) {
			account.nonce++
			account.updateLastBlocknumSentTx(blockNumber.Uint64())
		},
	},
	"GAA": {
		Description: "GAA creates a Gasless Approve.",
		GenerateTx: func(c *client.KaiaClient, account, _ *Account, nonce uint64, suggestedGasPrice *big.Int) *types.Transaction {
			approveTx := types.NewTransaction(
				nonce,
				TestContractInfos[ContractGaslessToken].GetAddress(c, GaslessTokenDeployer),
				common.Big0,
				100000,
				suggestedGasPrice,
				TestContractInfos[ContractGaslessToken].GenData(TestContractInfos[ContractGaslessSwapRouter].GetAddress(c, GaslessSwapRouterDeployer), abi.MaxUint256)) // Approve maximum amount
			signApproveTx, err := types.SignTx(approveTx, types.NewEIP155Signer(chainID), account.privateKey[0])
			if err != nil {
				return nil
			}
			return signApproveTx
		},
		PreSendBid: func(c *client.KaiaClient, account, _ *Account, nonce uint64, suggestedGasPrice *big.Int) error {
			return nil
		},
		PostSendBid: func(c *client.KaiaClient, account, _ *Account, nonce uint64, suggestedGasPrice *big.Int, blockNumber *big.Int) {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
			defer cancel()

			swapTx := types.NewTransaction(
				nonce+1,
				TestContractInfos[ContractGaslessSwapRouter].GetAddress(c, GaslessSwapRouterDeployer),
				common.Big0,
				500000,
				suggestedGasPrice,
				TestContractInfos[ContractGaslessSwapRouter].GenData(TestContractInfos[ContractGaslessToken].GetAddress(c, GaslessTokenDeployer), suggestedGasPrice))
			signSwapTx, _ := types.SignTx(swapTx, types.NewEIP155Signer(chainID), account.privateKey[0])

			_, err := c.SendRawTransaction(ctx, signSwapTx)
			if err != nil {
				if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
					fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", account.GetAddress().String(), nonce, err)
					fmt.Printf("Account(%v) nonce is added to %v\n", account.GetAddress().String(), nonce+1)
					account.nonce++
				} else {
					fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", account.GetAddress().String(), nonce, err)
				}
				return
			}

			account.nonce += 2
			account.updateLastBlocknumSentTx(blockNumber.Uint64())
		},
	},
	"GAS": {
		Description: "GAS creates a Gasless Swap.",
		GenerateTx: func(c *client.KaiaClient, account, _ *Account, nonce uint64, suggestedGasPrice *big.Int) *types.Transaction {
			swapTx := types.NewTransaction(
				nonce+1,
				TestContractInfos[ContractGaslessSwapRouter].GetAddress(c, GaslessSwapRouterDeployer),
				common.Big0,
				500000,
				suggestedGasPrice,
				TestContractInfos[ContractGaslessSwapRouter].GenData(TestContractInfos[ContractGaslessToken].GetAddress(c, GaslessTokenDeployer), suggestedGasPrice))
			signSwapTx, err := types.SignTx(swapTx, types.NewEIP155Signer(chainID), account.privateKey[0])
			if err != nil {
				return nil
			}
			return signSwapTx
		},
		PreSendBid: func(c *client.KaiaClient, account, _ *Account, nonce uint64, suggestedGasPrice *big.Int) error {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
			defer cancel()

			approveTx := types.NewTransaction(
				nonce,
				TestContractInfos[ContractGaslessToken].GetAddress(c, GaslessTokenDeployer),
				common.Big0,
				100000,
				suggestedGasPrice,
				TestContractInfos[ContractGaslessToken].GenData(TestContractInfos[ContractGaslessSwapRouter].GetAddress(c, GaslessSwapRouterDeployer), abi.MaxUint256)) // Approve maximum amount
			signApproveTx, err := types.SignTx(approveTx, types.NewEIP155Signer(chainID), account.privateKey[0])
			if err != nil {
				return err
			}

			_, err = c.SendRawTransaction(ctx, signApproveTx)
			if err != nil {
				if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
					fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", account.GetAddress().String(), nonce, err)
					fmt.Printf("Account(%v) nonce is added to %v\n", account.GetAddress().String(), nonce+1)
					account.nonce++
				} else {
					fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", account.GetAddress().String(), nonce, err)
				}
				return err
			}
			return nil
		},
		PostSendBid: func(c *client.KaiaClient, account, _ *Account, nonce uint64, suggestedGasPrice *big.Int, blockNumber *big.Int) {
			// Since gasless swap is the target, gasless bundle will be executed.
			// Therefore, gasless swap is performed normally and the nonce is incremented by 2.
			account.nonce += 2
			account.updateLastBlocknumSentTx(blockNumber.Uint64())
		},
	},
	"rGAA": {
		Description: "rGAA creates a reverted Gasless which occurs due to lack of balance during swap.",
		GenerateTx: func(c *client.KaiaClient, _, tmpAccount *Account, _ uint64, suggestedGasPrice *big.Int) *types.Transaction {
			// tmpAccount should have nonce 0
			approveTx := types.NewTransaction(
				0,
				TestContractInfos[ContractGaslessToken].GetAddress(c, GaslessTokenDeployer),
				common.Big0,
				100000,
				suggestedGasPrice,
				TestContractInfos[ContractGaslessToken].GenData(TestContractInfos[ContractGaslessSwapRouter].GetAddress(c, GaslessSwapRouterDeployer), abi.MaxUint256)) // Approve maximum amount
			signApproveTx, err := types.SignTx(approveTx, types.NewEIP155Signer(chainID), tmpAccount.privateKey[0])
			if err != nil {
				return nil
			}
			return signApproveTx
		},
		PreSendBid: func(c *client.KaiaClient, _, tmpAccount *Account, nonce uint64, suggestedGasPrice *big.Int) error {
			return nil
		},
		PostSendBid: func(c *client.KaiaClient, _, tmpAccount *Account, _ uint64, suggestedGasPrice *big.Int, blockNumber *big.Int) {
			// tmpAccount doesn't need nonce management.
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
			defer cancel()

			swapTx := types.NewTransaction(
				1,
				TestContractInfos[ContractGaslessSwapRouter].GetAddress(c, GaslessSwapRouterDeployer),
				common.Big0,
				500000,
				suggestedGasPrice,
				TestContractInfos[ContractGaslessSwapRouter].GenData(TestContractInfos[ContractGaslessToken].GetAddress(c, GaslessTokenDeployer), suggestedGasPrice))
			signSwapTx, _ := types.SignTx(swapTx, types.NewEIP155Signer(chainID), tmpAccount.privateKey[0])

			c.SendRawTransaction(ctx, signSwapTx)
		},
	},
	"rGAS": {
		Description: "rGAS creates a reverted Gasless which occurs due to lack of balance during swap.",
		GenerateTx: func(c *client.KaiaClient, _, tmpAccount *Account, _ uint64, suggestedGasPrice *big.Int) *types.Transaction {
			// tmpAccount should have nonce 1
			swapTx := types.NewTransaction(
				1,
				TestContractInfos[ContractGaslessSwapRouter].GetAddress(c, GaslessSwapRouterDeployer),
				common.Big0,
				500000,
				suggestedGasPrice,
				TestContractInfos[ContractGaslessSwapRouter].GenData(TestContractInfos[ContractGaslessToken].GetAddress(c, GaslessTokenDeployer), suggestedGasPrice))
			signSwapTx, err := types.SignTx(swapTx, types.NewEIP155Signer(chainID), tmpAccount.privateKey[0])
			if err != nil {
				return nil
			}
			return signSwapTx
		},
		PreSendBid: func(c *client.KaiaClient, _, tmpAccount *Account, _ uint64, suggestedGasPrice *big.Int) error {
			// tmpAccount doesn't need nonce management.
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
			defer cancel()

			approveTx := types.NewTransaction(
				0,
				TestContractInfos[ContractGaslessToken].GetAddress(c, GaslessTokenDeployer),
				common.Big0,
				100000,
				suggestedGasPrice,
				TestContractInfos[ContractGaslessToken].GenData(TestContractInfos[ContractGaslessSwapRouter].GetAddress(c, GaslessSwapRouterDeployer), abi.MaxUint256)) // Approve maximum amount
			signApproveTx, err := types.SignTx(approveTx, types.NewEIP155Signer(chainID), tmpAccount.privateKey[0])
			if err != nil {
				return err
			}

			_, err = c.SendRawTransaction(ctx, signApproveTx)
			if err != nil {
				return err
			}
			return nil
		},
		PostSendBid: func(c *client.KaiaClient, _, tmpAccount *Account, _ uint64, suggestedGasPrice *big.Int, blockNumber *big.Int) {
		},
	},
}
