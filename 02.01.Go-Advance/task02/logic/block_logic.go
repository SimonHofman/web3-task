package logic

import (
	"log"
	"task02/db"
	"task02/models"
)

type blockLogic struct{}

var BlockLogic = new(blockLogic)

func (bl *blockLogic) BlockSyncByNumber(blockNumber uint64) error {
	var search_block models.Block
	if err := db.DB.Where("number = ?", blockNumber).First(&search_block).Error; err != nil {
		log.Printf("Block %d already exits", blockNumber)
		return nil
	}

	return nil
}
