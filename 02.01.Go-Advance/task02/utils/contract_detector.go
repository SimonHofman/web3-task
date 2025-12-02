package utils

import (
	"fmt"
	"strings"
	"task02/services"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

type ContractDetector struct {
	eth *services.EthereumService
}

func NewContractDetector(eth *services.EthereumService) *ContractDetector {
	return &ContractDetector{
		eth: eth,
	}
}

// DetectProxyContract 检测代理合约并返回实现地址
func (cd *ContractDetector) DetectProxyContract(address string) (bool, string, string, error) {
	contractAddr := common.HexToAddress(address)

	// 检查 EIP-1967 代理
	if implementation, err := cd.checkEIP1967(contractAddr); err == nil && implementation != common.HexToAddress("0") {
		return true, implementation.Hex(), "EIP1967", nil
	}

	// 检查 EIP-1822 代理
	if implementation, err := cd.checkEIP1822(contractAddr); err == nil && implementation != common.HexToAddress("0") {
		return true, implementation.Hex(), "EIP1822", nil
	}

	// 检查 OpenZeppelin 代理
	if implementation, err := cd.checkOpenZeppelinProxy(contractAddr); err == nil && implementation != common.HexToAddress("0") {
		return true, implementation.Hex(), "OpenZeppelinProxy", nil
	}

	return false, "", "", nil
}

func (cd *ContractDetector) checkEIP1967(contractAddr common.Address) (common.Address, error) {
	// EIP-1967 逻辑合约地址存储位置
	slots := []string{
		"0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc", // 逻辑合约地址
		"0x7050c9e0f4ca769c69bd3a8ef740bc37934f8e2c036e5a723fd8ee048ed3f8c3", // Beacon
	}

	for _, slot := range slots {
		implementation, err := cd.getStorageAt(contractAddr, slot)
		if err == nil && implementation != common.HexToAddress("0") {
			return implementation, nil
		}
	}

	return common.Address{}, fmt.Errorf("not EIP1967 proxy")
}

func (cd *ContractDetector) checkEIP1822(contractAddr common.Address) (common.Address, error) {
	// EIP-1822 实现合约管理
	proxyManager := common.HexToAddress("0xa2ca1241f01d0b638ccac5dd4f4b71482f73db79")

	data := common.Hex2Bytes("c5f16f0f" + strings.Repeat("0", 24) + contractAddr.Hex()[2:])

	result, err := cd.eth.CallContract(ethereum.CallMsg{
		To:   &proxyManager,
		Data: data,
	})

	if err == nil && len(result) >= 32 {
		implementation := common.BytesToAddress(result[12:32])
		if implementation != common.HexToAddress("0") {
			return implementation, nil
		}
	}

	return common.Address{}, fmt.Errorf("not EIP1822 proxy")
}

// OpenZeppenlin 代理检测
func (cd *ContractDetector) checkOpenZeppelinProxy(contractAddr common.Address) (common.Address, error) {
	// 尝试调用 implementation() 方法
	data := common.Hex2Bytes("5c60da1b") // implementation()函数选择器

	result, err := cd.eth.CallContract(ethereum.CallMsg{
		To:   &contractAddr,
		Data: data,
	})

	if err == nil && len(result) >= 32 {
		implementation := common.BytesToAddress(result[12:32])
		if implementation != (common.Address{}) {
			return implementation, nil
		}
	}

	return common.Address{}, fmt.Errorf("no OpenZeppelin proxy")
}

func (cd *ContractDetector) getStorageAt(addr common.Address, slot string) (common.Address, error) {
	result, err := cd.eth.client.StorageAt(nil, addr, common.HexToHash(slot), nil)
	if err != nil {
		return common.Address{}, err
	}

	if len(result) >= 20 {
		return common.BytesToAddress(result[12:32]), nil
	}

	return common.Address{}, fmt.Errorf("invalid storage data")
}
