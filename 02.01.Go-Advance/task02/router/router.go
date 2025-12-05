package router

import (
	"task02/handlers"
	"task02/logic"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRouter(
	db *gorm.DB,
	blockSyncLogic *logic.BlockSyncLogic,
	addressLogic *logic.AddressLogic,
) *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.Use(gin.Recovery())

	blockHandler := handlers.NewBlockHandler(db, blockSyncLogic)
	addressHandle := handlers.NewAddressHandler(addressLogic)

	apiV1 := r.Group("/api/v1")
	{
		blockGroup := apiV1.Group("/block")
		{
			blockGroup.GET("/syncByNumber", blockHandler.BlockSyncByNumber)
		}
		addressGroup := apiV1.Group("/address")
		{
			addressGroup.GET("/search", addressHandle.GetAddressInfo)
		}
	}
	return r
}
