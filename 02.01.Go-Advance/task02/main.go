package main

import (
	"task02/config"
	"task02/db"
	"task02/router"
)

func main() {
	config.InitConfig("./etc/config.yaml")
	db.InitDB()
	r := router.InitRouter()
	r.Run(":" + config.GetConfig().Server.Port)
}
