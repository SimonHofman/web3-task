package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	contractBytecode = "60806040525f5f553480156011575f5ffd5b506101938061001f5f395ff3fe608060405234801561000f575f5ffd5b5060043610610034575f3560e01c8063371303c0146100385780636d4ce63c14610042575b5f5ffd5b610040610060565b005b61004a6100b0565b60405161005791906100d0565b60405180910390f35b5f5f81548092919061007190610116565b91905055507f15df0a6785153dfd625c8af51397704c7ddaff0690e5076243b5e52f4a0d54095f546040516100a691906100d0565b60405180910390a1565b5f5f54905090565b5f819050919050565b6100ca816100b8565b82525050565b5f6020820190506100e35f8301846100c1565b92915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f610120826100b8565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203610152576101516100e9565b5b60018201905091905056fea26469706673582212202b7ce02c41fd9ba419d6297052831c1772f48122fb97b0fb295e85063b979bee64736f6c634300081e0033"
)

func main() {
	client, err := ethclient.Dial("https://ethereum-sepolia-rpc.publicnode.com")
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA("ed0e5a7070c2d4d5dde2608ea8711dacf7dd6f1fc2b48643683c3b8edce8bbf5")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	data, err := hex.DecodeString(contractBytecode)
	if err != nil {
		log.Fatal(err)
	}

	tx := types.NewContractCreation(nonce, big.NewInt(0), 3000000, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Transaction sent: %s\n", signedTx.Hash().Hex())
}
