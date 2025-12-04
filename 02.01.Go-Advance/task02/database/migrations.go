package database

import (
	"task02/models"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	// 关闭外键约束
	db.Exec("set foreign_key_checks = 0")

	// 自动迁移所有模型
	err := db.AutoMigrate(
		&models.Block{},
		&models.Transaction{},
		&models.Contract{},
		&models.AccountOperation{},
		&models.ContractOperation{},
	)
	if err != nil {
		return err
	}

	// 启动外键约束
	db.Exec("set foreign_key_checks = 1")

	return nil
}
