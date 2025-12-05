package main

import (
	"fmt"
	"log"
	"task02/config"
	"task02/database"
	"task02/logic"
	"task02/router"
	"task02/services"
	"task02/utils"
)

func main() {
	cfg := config.NewConfiguration("./etc/config.yaml")

	// 初始化数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.MySQL.User, cfg.MySQL.Password, cfg.MySQL.Host, cfg.MySQL.Port, cfg.MySQL.Database)
	db := database.InitDatabase(dsn)
	database.AutoMigrate(db)

	// 初始化以太坊服务
	ethService := services.NewEthereumService(cfg)

	// 初始化业务服务
	transactionAnalyzer := utils.NewTransactionAnalyzer()

	blockSyncLogic := logic.NewBlockSyncLogic(db, ethService, transactionAnalyzer)
	addressLogic := logic.NewAddressLogic(db, ethService)

	r := router.InitRouter(db, blockSyncLogic, addressLogic)
	err := r.Run(":" + cfg.Server.Port)
	if err != nil {
		log.Fatal(err)
	}
}
