package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://ethereum-sepolia-rpc.publicnode.com")
	if err != nil {
		log.Fatal(err)
	}
	blockNumber := big.NewInt(5671744)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(block.Number())
	fmt.Println(block.Hash().Hex())
	fmt.Println(block.Time())
	fmt.Println(len(block.Transactions()))

	privateKey, err := crypto.HexToECDSA("ed0e5a7070c2d4d5dde2608ea8711dacf7dd6f1fc2b48643683c3b8edce8bbf5")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *crypto.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(1000000000000000000)
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress("0x886f420125a9789de4eafb345e660156560118fc")
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	channID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTX, err := types.SignTx(tx, types.NewEIP155Signer(channID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTX)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", signedTX.Hash().Hex())
}
