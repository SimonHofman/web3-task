package model

import "time"

type Block struct {
	ID          uint64 `gorm:"primaryKey"`
	BlockNumber uint64 `gorm:"uniqueIndex;not null"`
	BlockHash   string `grom:"size:66;uniqueIndex;not null"`
	ParentHash  string `grom:"size:66;not null"`
	Timestamp   uint64 `grom:"not null"`
	Nonce       uint64
	Difficulty  string
	GasLimit    uint64
	GasUsed     uint64
	Miner       string `gorm:"size:42"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
