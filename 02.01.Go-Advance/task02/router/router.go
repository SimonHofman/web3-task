package router

import (
	"task02/handlers"
	"task02/logic"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB, blockSyncLogic *logic.BlockSyncLogic) *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.Use(gin.Recovery())

	blockHandler := handlers.NewBlockHandler(db, blockSyncLogic)

	apiV1 := r.Group("/api/v1")
	{
		blockGroup := apiV1.Group("/block")
		{
			//blockGroup.GET("/searchByNumber", services.BlockSearchByNumber)
			blockGroup.GET("/syncByNumber", blockHandler.BlockSyncByNumber)
			//blockGroup.GET("/searchByHash", handler.BlockSearchByHash)
			//blockGroup.GET("/syncByHash", handler.BlockSyncByHash)
		}
	}
	return r
}
