package main

import (
	"task02/internal/config"
	//"task02/pkg/client"
	"task02/pkg/db"
	"task02/router"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig("./etc/config.yaml")
	db.InitDB()
	r := router.InitRouter()
	ginErr := r.Run(":" + config.GetConfig().Server.Port)
}
