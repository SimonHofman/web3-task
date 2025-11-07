package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

//题目1：使用SQL扩展库进行查询
//假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
//要求 ：
//编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
//编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。

// Employee 结构体定义
type Employee struct {
	ID         int    `db:"id"`
	Name       string `db:"name"`
	Department string `db:"department"`
	Salary     int    `db:"salary"`
}

// insertTestData 插入测试数据
func insertTestData(db *sqlx.DB) error {
	// 创建 employees 表（如果不存在）
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS employees (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    name TEXT NOT NULL,
	    department TEXT NOT NULL,
	    salary INTEGER NOT NULL
	)`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	// 检查是否已有数据，避免重复插入
	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM employees")
	if err != nil {
		return err
	}
	if count > 0 {
		fmt.Println("数据已经存在，跳过插入测试数据")
		return nil
	}

	// 插入测试数据
	testData := []Employee{
		{Name: "张三", Department: "技术部", Salary: 15000},
		{Name: "李四", Department: "技术部", Salary: 18000},
		{Name: "王五", Department: "销售部", Salary: 12000},
		{Name: "赵六", Department: "技术部", Salary: 20000},
		{Name: "钱七", Department: "人事部", Salary: 10000},
	}

	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO employees (name, department, salary) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, emp := range testData {
		_, err := stmt.Exec(emp.Name, emp.Department, emp.Salary)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	fmt.Println("测试数据插入成功")
	return nil
}

// queryTechEmployees 查询技术部所有员工
func queryTechEmployess(db *sqlx.DB) ([]Employee, error) {
	var employees []Employee
	query := "SELECT id, name, department, salary FROM employees Where department = ?"
	err := db.Select(&employees, query, "技术部")
	if err != nil {
		return nil, err
	}
	return employees, nil
}

// queryHighestPaidEmployee 查询工资最高的员工
func queryHighestPaidEmployee(db *sqlx.DB) (*Employee, error) {
	var employee Employee
	query := "SELECT id, name, department, salary from employees ORDER BY salary DESC LIMIT 1"
	err := db.Get(&employee, query)
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func main() {
	// 连接 SQLite 数据库
	db, err := sqlx.Connect("sqlite3", "test.db")
	if err != nil {
		fmt.Println("数据库连接失败:", err)
	}
	defer db.Close()

	// 插入测试数据
	err = insertTestData(db)
	if err != nil {
		fmt.Println("插入测试数据失败:", err)
	}

	// 查询技术部员工
	techEmployees, err := queryTechEmployess(db)
	if err != nil {
		fmt.Println("查询技术部员工失败:", err)
	}
	fmt.Println("技术部员工")
	for _, emp := range techEmployees {
		fmt.Printf("员工ID: %d, 姓名: %s, 部门: %s, 工资: %d\n", emp.ID, emp.Name, emp.Department, emp.Salary)
	}

	// 查询工资最高的员工
	highestPaid, err := queryHighestPaidEmployee(db)
	if err != nil {
		fmt.Println("查询工资最高的员工失败:", err)
	}
	fmt.Printf("工资最高的员工: ID: %d, 姓名: %s, 部门: %s, 工资: %d", highestPaid.ID, highestPaid.Name, highestPaid.Department, highestPaid.Salary)
}
