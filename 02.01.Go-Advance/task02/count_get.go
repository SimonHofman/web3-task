package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
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

	//privateKey, err := crypto.HexToECDSA("ed0e5a7070c2d4d5dde2608ea8711dacf7dd6f1fc2b48643683c3b8edce8bbf5")
	//if err != nil {
	//	panic(err)
	//}

	//publicKey := privateKey.Public()
	//publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	//if !ok {
	//	panic("error casting public key to ECDSA")
	//}

	//fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	//nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	//if err != nil {
	//	panic(err)
	//}

	//gasPrice, err := client.SuggestGasPrice(context.Background())
	//if err != nil {
	//	panic(err)
	//}

	contractABI, err := abi.JSON(strings.NewReader(countAbi))
	if err != nil {
		panic(err)
	}

	//chainID := big.NewInt(int64(11155111))

	// 调用get函数
	methodName := "get"
	callInput, err := contractABI.Pack(methodName)
	if err != nil {
		panic(err)
	}
	to := common.HexToAddress(contractAddress)
	callMsg := ethereum.CallMsg{
		To:   &to,
		Data: callInput,
	}

	result, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		log.Fatal(err)
	}

	// 方法一：通过
	var unpacked *big.Int
	err = contractABI.UnpackIntoInterface(&unpacked, "get", result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("num is: %s\n", unpacked.String())

	// 方法二：通过
	//var unpacked interface{}
	//err = contractABI.UnpackIntoInterface(&unpacked, "get", result)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//// 类型断言转换为*big.Int
	//if num, ok := unpacked.(*big.Int); ok {
	//	fmt.Printf("num is: %s\n", num.String())
	//}

	// 方法三：通过
	//results, err := contractABI.Unpack("get", result)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//if len(results) > 0 {
	//	if num, ok := results[0].(*big.Int); ok {
	//		fmt.Printf("num is: %s\n", num.String())
	//	}
	//}
}
