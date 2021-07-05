package modules

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/LeNgocPhuc99/blockchain-api/models"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// GetLatestBlock retrives thr latest block from the blockchain
func GetLatestBlock(client ethclient.Client) *models.Block {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	// query the lastest block from blockchain
	header, _ := client.HeaderByNumber(context.Background(), nil)
	blockNumber := big.NewInt(header.Number.Int64())
	block, err := client.BlockByNumber(context.Background(), blockNumber)

	if err != nil {
		log.Fatal(err)
	}

	// build block for response
	resBlock := &models.Block{
		BlockNumber:       block.Number().Int64(),
		Timestamp:         block.Time(),
		Difficulty:        block.Difficulty().Uint64(),
		Hash:              block.Hash().String(),
		TransactionsCount: len(block.Transactions()),
		Transactions:      []models.Transaction{},
	}

	for _, tx := range block.Transactions() {
		resBlock.Transactions = append(resBlock.Transactions, models.Transaction{
			Hash:     tx.Hash().String(),
			Value:    tx.Value().String(),
			Gas:      tx.Gas(),
			GasPrice: tx.GasPrice().Uint64(),
			Nonce:    tx.Nonce(),
			To:       tx.To().String(),
		})
	}

	return resBlock
}

// GetTxByHash retrive info about a transaction by given tx hash
func GetTxByHash(client ethclient.Client, hash common.Hash) *models.Transaction {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	// get tx from blockchain
	tx, pending, err := client.TransactionByHash(context.Background(), hash)
	if err != nil {
		fmt.Println(err)
	}

	return &models.Transaction{
		Hash:     tx.Hash().String(),
		Value:    tx.Value().String(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice().Uint64(),
		To:       tx.To().String(),
		Pending:  pending,
		Nonce:    tx.Nonce(),
	}
}

// GetAddressBalance return balance of given address
func GetAddressBalance(client ethclient.Client, address string) (string, error) {
	account := common.HexToAddress(address)
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return "0", err
	}

	return balance.String(), nil
}

// TransferEth send ether from one account to another
func TransferEth(client ethclient.Client, privKey string, to string, amount int64) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	privateKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		return "", err
	}

	// get public address of the account
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", err
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// get the nonce that uses for the account's transaction
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", err
	}

	value := big.NewInt(amount)
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}

	toAddress := common.HexToAddress(to)
	var data []byte

	// create transaction
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return "", err
	}

	// sign the transaction use the sender's private key
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", err
	}

	// broadcast the transaction to the entire network
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().String(), nil

}
