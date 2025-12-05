package handlers

import (
	"strconv"
	"task02/logic"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BlockHandler struct {
	db             *gorm.DB
	blockSyncLogic *logic.BlockSyncLogic
}

func NewBlockHandler(db *gorm.DB, blockSyncLogic *logic.BlockSyncLogic) *BlockHandler {
	return &BlockHandler{
		db:             db,
		blockSyncLogic: blockSyncLogic,
	}
}

func (bh *BlockHandler) SyncBlockByNumber(c *gin.Context) {
	blockNumber, exist := c.GetQuery("blockNumber")
	if !exist {
		c.JSON(400, gin.H{
			"message": "blockNumber is required",
		})
		return
	}

	// 将字符串转换为uint64
	blockNum, err := strconv.ParseUint(blockNumber, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid blockNumber format",
		})
		return
	}

	err = bh.blockSyncLogic.SyncBlockByNumber(blockNum)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "sync block failed",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "sync block success",
		"block":   blockNumber,
	})
}

func (bl *BlockHandler) SyncLatestBlock(c *gin.Context) {
	blockNumber, err := bl.blockSyncLogic.SyncLatestBlock()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "sync block failed",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "sync block success",
		"block":   blockNumber,
	})
}
