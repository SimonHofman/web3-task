package router

import (
	"task02/internal/handler"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.Use(gin.Recovery())

	apiV1 := r.Group("/api/v1")
	{
		blockGroup := apiV1.Group("/block")
		{
			blockGroup.GET("/searchByNumber", handler.BlockSearchByNumber)
			blockGroup.GET("/syncByNumber", handler.BlockSyncByNumber)
			blockGroup.GET("/searchByHash", handler.BlockSearchByHash)
			blockGroup.GET("/syncByHash", handler.BlockSyncByHash)
		}
	}
	return r
}
