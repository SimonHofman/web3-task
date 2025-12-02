package models

import "time"

type AccountOperation struct {
	ID            uint64 `gorm:"primaryKey"`
	Address       string `gorm:"size:42;index;not null"`
	OperationType string `gorm:"size:50;index"`
	BlockNumber   uint64 `gorm:"index"`
	FromAddress   string `gorm:"size:42;index"`
	ToAddress     string `gorm:"size:42;index"`
	Value         string
	Timestamp     uint64 `gorm:"not null"`
	CreatedAt     time.Time
}
