package models

import (
	"time"
)

// TxTypeLegacy 		传统交易类型（以太坊早期版本的交易格式）
// TxTypeEIP1559 		EIP-1559引入的新交易类型，支持基础费用和优先费用
// TxTypeAccessList 	EIP-2930引入的访问列表交易类型，允许指定预访问的合约和存储槽
// TxTypeContractDeploy 合约部署交易类型
// TxTypeTransfer 		普通转账交易类型
// TxTypeContractCall 	合约调用交易类型
// TxTypeTokenTransfer 	代币转账交易类型（如ERC-20代币转移）
// TxTypeTokenApprove 	代币授权交易类型（如ERC-20代币授权）
// TxTypeSwap 			交易交换类型（如去中心化交易所中的代币兑换）
// TxTypeMulticall 		多重调用交易类型（一次交易执行多个操作）
// TxTypeDelegate 		委托交易类型（如权益委托）
const (
	TxTypeLegacy         = "Legacy"
	TxTypeEIP1559        = "EIP1559"
	TxTypeAccessList     = "AccessList"
	TxTypeContractDeploy = "ContractDeploy"
	TxTypeTransfer       = "Transfer"
	TxTypeContractCall   = "ContractCall"
	TxTypeTokenTransfer  = "TokenTransfer"
	TxTypeTokenApprove   = "TokenApprove"
	TxTypeSwap           = "Swap"
	TxTypeMulticall      = "Multicall"
	TxTypeDelegate       = "Delegate"
)

// MaxPriorityFeePerGas 是EIP-1559交易中的最大优先费用，表示用户愿意为矿工支付的小费上限，可以为空
// MaxFeePerGas 		是EIP-1559交易中的最大费用，表示用户愿意支付的每单位Gas的最高价格，可以为空
// EffectiveGasPrice 	是交易实际使用的Gas价格，对于EIP-1559交易是基础费用和优先费用之和，对于传统交易就是GasPrice
type Transaction struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`
	// 区块关联
	BlockID     uint64 `gorm:"index:idx_block_id:not null"`
	BlockNumber uint64 `gorm:"index:idx_block_number;not null"`

	// 交易基本信息
	TransactionHash  string `grom:"size:66;uniqueIndex:idx_tx_hash;not null"`
	TransactionIndex uint64 `gorm:"not null"`
	FromAddress      string `grom:"size:42;index:idx_from_address;not null"`
	ToAddress        string `grom:"size:42;index:idx_to_address"`
	Value            string `gorm:"type:text"`

	// Gas 相关
	GasPrice             uint64 `gorm:"not null"`
	GasLimit             uint64 `gorm:"not null"`
	GasUsed              uint64 `gorm:"not null"`
	MaxPriorityFeePerGas *uint64
	MaxFeePerGas         *uint64
	EffectiveGasPrice    uint64

	// 其他字段
	Nonce     uint64 `gorm:"not null"`
	InputData string `gorm:"type:longtext"`
	Status    uint64 `gorm:"default:1"` // 0:失败, 1:成功

	// 交易类型相关字段
	TxType           string `gorm:"size:50;index:idx_tx_type"`          // 交易类型
	IsContractDeploy bool   `gorm:"default:false"`                      // 是否是合约部署
	ContractAddress  string `gorm:"size:42;index:idx_contract_address"` // 部署的合约地址

	// 访问列表（EIP-2930）
	AccessList string `gorm:"type:json"`

	// 关联到所属的区块
	Block *Block `gorm:"foreignKey:BlockID"`

	CreatedAt time.Time `gorm:"index:idx_tx_created_at"`
	UpdatedAt time.Time
}
