package models

import (
	"time"
)

type Contract struct {
	ID              uint64 `gorm:"primaryKey"`
	ContractAddress string `gorm:"size:42;uniqueIndex;not null"`
	IsProxy         bool   `gorm:"default:false"`
	Implementation  string `gorm:"size:42"` // 真实实现地址
	ProxyType       string `gorm:"size:50"` // EIP1967, EIP1822
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type ContractOperation struct {
	ID              uint64 `gorm:"primaryKey"`
	ContractAddress string `gorm:"size:42;index;not null"`
	OperationType   string `gorm:"size:50;index"`
	TransactionHash string `gorm:"size:66;index"`
	BlockNumber     uint64 `gorm:"index"`
	FromAddress     string `gorm:"size:42:index"`
	InputData       string `gorm:"type:longtext"`
	EventName       string `gorm:"size:100"`
	EventData       string `gorm:"type:json"`
	CreatedAt       time.Time
}
