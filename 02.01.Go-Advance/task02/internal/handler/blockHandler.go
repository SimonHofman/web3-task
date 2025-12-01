package handler

import (
	"task02/internal/logic"
	"task02/internal/model"

	"github.com/gin-gonic/gin"
)

func BlockSearchByNumber(c *gin.Context) {
	blockNumber := c.GetQuery("blockNumber")
}

func BlockSyncByNumber(c *gin.Context) {

}

func BlockSearchByHash(c *gin.Context) {
	blockNumber := c.GetQuery("blockHash")
}

func BlockSyncByHash(c *gin.Context) {

}
