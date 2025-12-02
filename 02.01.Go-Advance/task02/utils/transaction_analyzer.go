package utils

import (
	"encoding/hex"
	"fmt"
	"strings"
	"task02/models"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type TransactionAnalyzer struct {
	methodSignatures map[string]string // 方法前面到名称的映射
}

// 交易分析结果
type TransactionAnanlysis struct {
	BasicType       string
	BusinessType    string
	MethodName      string
	MethodSignature string
	IsContractCall  bool
	ContractAddress string
}

func NewTransactionAnalyzer() *TransactionAnalyzer {
	analyzer := &TransactionAnalyzer{
		methodSignatures: make(map[string]string),
	}
	analyzer.initMethodSignatures()
	return analyzer
}

// 初始化常见的方法签名
func (ta *TransactionAnalyzer) initMethodSignatures() {
	// ERC20 标准方法
	ta.methodSignatures["a9059cbb"] = "transfer(address,uint256)"             // transfer
	ta.methodSignatures["23b872dd"] = "transferFrom(address,address,uint256)" // transferFrom
	ta.methodSignatures["095ea7b3"] = "approve(address,uint256)"              // approve
	ta.methodSignatures["70a08231"] = "balanceOf(address)"                    // balanceOf

	// ERC721 标准方法
	ta.methodSignatures["42842e0e"] = "safeTransferFrom(address,address,uint256)"
	ta.methodSignatures["b88d4fde"] = "safeTransferFrom(address,address,uint256,bytes)"

	// Uniswap 相关方法
	ta.methodSignatures["f305d719"] = "addLiquidityETH"       // 添加流动性
	ta.methodSignatures["fb3bdb41"] = "swapETHForExactTokens" // 兑换
	ta.methodSignatures["7ff36ab5"] = "swapExactETHForTokens" // 精确ETH兑换代币
	ta.methodSignatures["18cbafe5"] = "swapExactTokensForETH" // 精确代币兑换ETH

	// 其他常见方法
	ta.methodSignatures["d0e30db0"] = "deposit()"         // 存款
	ta.methodSignatures["2e1a7d4d"] = "withdraw(uint256)" // 取款
}

// 分析交易类型
func (ta *TransactionAnalyzer) AnalyzeTransactionType(tx *types.Transaction, receipt *types.Receipt) *TransactionAnanlysis {
	analysis := &TransactionAnanlysis{
		BasicType:      ta.getBasicTransactionType(tx),
		BusinessType:   "Unkown",
		MethodName:     "",
		IsContractCall: len(tx.Data()) > 0,
	}

	// 检查是否是合约部署
	if tx.To() == nil {
		analysis.BusinessType = models.TxTypeContractDeploy
		if receipt != nil && receipt.ContractAddress != (common.Address{}) {
			analysis.ContractAddress = receipt.ContractAddress.Hex()
		}
		return analysis
	}

	// 分析输入数据
	if len(tx.Data()) >= 4 {
		methodSig := hex.EncodeToString(tx.Data()[:4])
		analysis.MethodSignature = methodSig

		// 查找方法名称
		if name, exists := ta.methodSignatures[methodSig]; exists {
			analysis.MethodName = name
			analysis.BusinessType = ta.getBusinessTypeFromMethod(methodSig, name)
		} else {
			analysis.MethodName = fmt.Sprintf("0x%s", methodSig)
			analysis.BusinessType = models.TxTypeContractCall
		}
	}

	if len(tx.Data()) == 0 && tx.Value().Sign() > 0 {
		analysis.BusinessType = models.TxTypeTransfer
	}

	// 分析日志事件（用于识别代表转账等）
	if receipt != nil {

	}

	return analysis
}

func (ta *TransactionAnalyzer) getBasicTransactionType(tx *types.Transaction) string {
	// 检查交易类型（EIP-2718）
	switch tx.Type() {
	case types.LegacyTxType:
		return models.TxTypeLegacy
	case types.AccessListTxType:
		return models.TxTypeAccessList
	case types.DynamicFeeTxType:
		return models.TxTypeEIP1559
	default:
		return "Unkown"
	}
}

func (ta *TransactionAnalyzer) getBusinessTypeFromMethod(methodSig, methodName string) string {
	switch methodSig {
	case "a9059cbb", "23b872dd":
		return models.TxTypeTokenTransfer
	case "095ea7b3":
		return models.TxTypeTokenApprove
	case "f305d719", "fb3bdb41", "7ff36ab5", "18cbafe5":
		return models.TxTypeSwap
	default:
		if strings.Contains(methodName, "transfer") {
			return models.TxTypeTokenTransfer
		}
		return models.TxTypeContractCall

	}
}

func (ta *TransactionAnalyzer) analyzeLogs(analysis *TransactionAnanlysis, logs []*types.Log) {
	for _, log := range logs {
		if len(log.Topics) > 0 {
			// ERC20 Transfer 事件
			if log.Topics[0] == common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef") {
				if analysis.BusinessType == models.TxTypeContractCall {
					analysis.BusinessType = models.TxTypeTokenTransfer
				}
			}
			// ERC20 Approval 事件
			if log.Topics[0] == common.HexToHash("0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925") {
				analysis.BusinessType = models.TxTypeTokenApprove
			}
		}
	}
}
