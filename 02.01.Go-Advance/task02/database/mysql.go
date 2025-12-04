package database

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDatabase(dsn string) *gorm.DB {
	var err error
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败", err)
	}
	return db
}
