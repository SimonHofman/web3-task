package models

import (
	"time"
)

// Block 结构体表示区块链中的一个区块实体，用于存储和管理区块链数据
// 包含区块头信息、交易引用以及数据库元数据
// ID              是数据库中的唯一标识符，自动递增
// BlockNumber     是区块链中的区块序号，必须唯一且非空
// BlockHash       是区块的哈希值，必须唯一且非空，长度限制为66个字符
// ParentHash      是前一个区块的哈希值，非空，长度限制为66个字符
// Timestamp       是区块创建时的时间戳(Unix时间)，非空
// Nonce           是用于工作量证明的随机数，默认值为0
// Difficulty      是当前区块的挖矿难度，以文本形式存储以适应大数值
// Size            是区块的大小(以字节为单位)，默认值为0
// GasLimit        是本区块中所有交易可消耗的最大Gas数量，非空
// GasUsed         是本区块中所有交易实际消耗的Gas数量，非空
// Miner           是挖出此区块的矿工地址，长度限制为42个字符，建立索引以便快速查询
// ExtraData       是矿工添加的额外数据，以文本形式存储
// MixHash         是用于工作量证明验证的哈希值，长度限制为66个字符
// BaseFeePerGas   是EIP-1559引入的基础费用，每个Gas单位的价格，可以为空
// WithdrawalsRoot 是提款记录的Merkle树根哈希值，长度限制为66个字符
// Transaction     是与该区块关联的交易列表，通过外键BlockID关联到Block的ID
// CreatedAt       记录数据库中该记录的创建时间
// UpdatedAt       记录数据库中该记录的最后更新时间
type Block struct {
	ID              uint64 `gorm:"primaryKey;autoIncrement"`
	BlockNumber     uint64 `gorm:"uniqueIndex:idx_block_number;not null"`
	BlockHash       string `gorm:"size:66;uniqueIndex:idx_block_hash;not null"`
	ParentHash      string `gorm:"size:66;not null"`
	Timestamp       uint64 `gorm:"not null"`
	Nonce           uint64 `gorm:"default:0"`
	Difficulty      string `gorm:"type:text"`
	Size            uint64 `gorm:"default:0"`
	GasLimit        uint64 `gorm:"not null"`
	GasUsed         uint64 `gorm:"not null"`
	Miner           string `gorm:"size:42;index:idx_miner"`
	ExtraData       string `gorm:"type:text"`
	MixHash         string `gorm:"size:66"`
	BaseFeePerGas   *uint64
	WithdrawalsRoot string `gorm:"size:66"`

	Transaction []*Transaction `gorm:"foreignKey:BlockID;references:ID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
