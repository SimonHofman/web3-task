package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//假设有一个名为 students 的表，包含字段 id （主键，自增）、 name （学生姓名，字符串类型）、 age （学生年龄，整数类型）、 grade （学生年级，字符串类型）。
//要求 ：
//编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。
//编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
//编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
//编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。

type strudents struct {
	ID    int `gorm:"primaryKey"`
	Name  string
	Age   int
	Grade string
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败：" + err.Error())
	}

	////创建表
	//err = db.AutoMigrate(&strudents{})
	//if err != nil {
	//	panic("自动迁移失败：" + err.Error())
	//}
	//
	////写入初始数据
	//db.Create(&strudents{
	//	Name:  "张三",
	//	Age:   20,
	//	Grade: "三年级",
	//})

	////批量查询
	//var student_arr []strudents
	//err = db.Where("age > ?", 18).Find(&student_arr).Error
	//if err != nil {
	//	panic("查询失败：" + err.Error())
	//}
	//fmt.Println("批量查询结果：")
	//for _, item := range student_arr {
	//	fmt.Printf("Name: %s, Age: %d, Grade: %s\n", item.Name, item.Age, item.Grade)
	//}

	//// 更新
	//err = db.Model(&strudents{}).Where("name = ?", "张三").Update("grade", "四年级").Error
	//if err != nil {
	//	panic("更新失败：" + err.Error())
	//}

	err = db.Where("age < ? ", 15).Delete(&strudents{}).Error
	if err != nil {
		panic("删除失败：" + err.Error())
	}
}
