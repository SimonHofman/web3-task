package handlers

import (
	"strconv"
	"task02/logic"

	"github.com/gin-gonic/gin"
)

func BlockSyncByNumber(c *gin.Context) {
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

	err = logic.BlockLogic.BlockSyncByNumber(blockNum)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "sync block failed",
		})
		return
	}
}
