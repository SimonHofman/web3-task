package model

import (
	"time"
)

const (
	TxLegacy         = "Legacy"
	TxEIP1559        = "EIP1559"
	TxAccessList     = "AccessList"
	TxContractDeploy = "ContractDeploy"
	TxTransfer       = "Transfer"
	TxContractCall   = "ContractCall"
	TxTokenTransfer  = "TokenTransfer"
	TxTokenApprove   = "TokenApprove"
	TxSwap           = "Swap"
)

type Transaction struct {
	ID               uint64 `gorm:"primaryKey"`
	BlockID          uint64 `gorm:"index:not null"`
	BlockNumber      uint64 `gorm:"index;not null"`
	TransactionHash  string `grom:"size:66;uniqueIndex;not null"`
	TransactionIndex uint64 `gorm:"not null"`
	FromAddress      string `grom:"size:42;index"`
	ToAddress        string `grom:"size:42;index"`
	Value            string
	GasPrice         uint64
	GasLimit         uint64
	GasUsed          uint64
	Nonce            uint64
	InputData        string `gorm:"type:text"`
	Status           uint64 // 0:失败, 1:成功

	// 交易类型相关字段
	TxType           string `gorm:"size:50;index"` // 交易类型
	IsContractDeploy bool   `gorm:"default:false"` // 是否是合约部署
	ContractAddress  string `gorm:"size:42"`       // 部署的合约地址

	// EIP-1559相关字段
	MaxPriorityFeePerGas uint64 // 最大优先费用
	MaxFeePerGas         uint64 // 最大费用
	EffectiveGasPrice    uint64 // 实际燃料价格

	// 关联到所属的区块
	Block *Block `gorm:"foreignKey:BlockID"`

	CreatedAt time.Time
}
