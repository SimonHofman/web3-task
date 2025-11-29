package model

import (
	"time"
)

type Transaction struct {
	ID              uint64 `gorm:"primaryKey"`
	BlockNumber     uint64 `gorm:"index;not null"`
	TransactionHash string `grom:"size:66;uniqueIndex;not null"`
	FromAddress     string `grom:"size:42;index"`
	ToAddress       string `grom:"size:42;index"`
	Value           string
	GasPrice        uint64
	GasLimit        uint64
	GasUsed         uint64
	Nonce           uint64
	InputData       string `gorm:"type:text"`
	Status          uint64 // 0:失败, 1:成功
	CreatedAt       time.Time
}
