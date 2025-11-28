package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	contractAddress = "0x8E60307716931638FE322F11789AedC8c3aaECc1"
	countAbi        = `[{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"num","type":"uint256"}],"name":"Counted","type":"event"},{"inputs":[],"name":"get","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"inc","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
)

func main() {
	client, err := ethclient.Dial("https://ethereum-sepolia-rpc.publicnode.com")
	if err != nil {
		panic(err)
	}

	privateKey, err := crypto.HexToECDSA("ed0e5a7070c2d4d5dde2608ea8711dacf7dd6f1fc2b48643683c3b8edce8bbf5")
	if err != nil {
		panic(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		panic(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		panic(err)
	}

	contractABI, err := abi.JSON(strings.NewReader(countAbi))
	if err != nil {
		panic(err)
	}

	chainID := big.NewInt(int64(11155111))

	// 调用inc函数
	methodName := "inc"
	input, err := contractABI.Pack(methodName)

	tx := types.NewTransaction(nonce, common.HexToAddress(contractAddress), big.NewInt(0), 3000000, gasPrice, input)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Transaction send: %s\n", signedTx.Hash().Hex())
}
