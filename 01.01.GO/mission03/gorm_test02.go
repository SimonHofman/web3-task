package main

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//假设有两个表： accounts 表（包含字段 id 主键， balance 账户余额）和 transactions 表（包含字段 id 主键， from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
//要求 ：
//编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。在事务中，需要先检查账户 A 的余额是否足够，如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。

type accounts struct {
	ID      int `gorm:"primaryKey"`
	Balance int
}

type transactions struct {
	ID              int `gorm:"primaryKey"`
	From_account_id int
	To_account_id   int
	Amount          int
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败：" + err.Error())
	}

	err = db.AutoMigrate(&accounts{}, &transactions{})
	if err != nil {
		panic("创建表失败：" + err.Error())
	}

	//err = db.Create(&accounts{ID: 1, Balance: 1000}).Error
	//err = db.Create(&accounts{ID: 2, Balance: 20}).Error
	//err = db.Create(&accounts{ID: 3, Balance: 80}).Error

	transactionDemo(db, 3, 2, 100)
}

func transactionDemo(db *gorm.DB, fromId, toId, transferAccount int) {
	err := db.Transaction(func(tx *gorm.DB) error {
		// 判断账户fromId余额是否充足
		var fromAccount accounts
		if err := tx.First(&fromAccount, fromId).Error; err != nil {
			panic("查询账户失败：" + err.Error())
		}
		if fromAccount.Balance < transferAccount {
			panic("账户余额不足")
		}

		// fromId账号减少100元
		if err := tx.Model(&accounts{}).
			Where("id = ?", fromId).
			Update("balance", gorm.Expr("balance - ?", transferAccount)).Error; err != nil {
			panic("更新账户失败：" + err.Error())
		}

		// toId账号中增加100元
		if err := tx.Model(&accounts{}).
			Where("id = ?", toId).
			Update("balance", gorm.Expr("balance + ?", transferAccount)).Error; err != nil {
			panic("更新账户失败：" + err.Error())
		}

		// 将转账记录写入transactions表
		transaction := transactions{
			From_account_id: fromId,
			To_account_id:   toId,
			Amount:          transferAccount,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			panic("写入转账记录失败：" + err.Error())
		}

		// 所有操作都成功了，提交事务
		return nil
	})

	if err != nil {
		// 如果有错误发生，回滚事务
		fmt.Println("事务执行失败，已回滚：", err)
	} else {
		// 没有错误发生，提交事务
		fmt.Println("事务执行成功，已提交：")
	}
}
