package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

//题目2：实现类型安全映射
//假设有一个 books 表，包含字段 id 、 title 、 author 、 price 。
//要求 ：
//定义一个 Book 结构体，包含与 books 表对应的字段。
//编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。

// Book 结构体定义，与 books 表字段对应
type Book struct {
	Id     int     `db:"id"`
	Title  string  `db:"title"`
	Author string  `db:"author"`
	Price  float64 `db:"price"`
}

// insertBookTestData 插入书籍测试数据
func insertBookTestData(db *sqlx.DB) error {
	// 创建 books 表（如果不存在）
	createTableSQL := `CREATE TABLE IF NOT EXISTS books (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	title TEXT NOT NULL,
    	author TEXT NOT NULL,
    	price REAL NOT NULL
	)`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	// 检查是否已有数据，避免重复插入
	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM books")
	if err != nil {
		return err
	}
	if count > 0 {
		fmt.Println("书籍书籍已存在，跳过插入测试数据")
		return nil
	}

	// 插入测试数据
	testData := []Book{
		{Title: "Go语言编程", Author: "许式伟", Price: 79.00},
		{Title: "Python核心编程", Author: "Wesley Chun", Price: 128.00},
		{Title: "算法导论", Author: "Thomas Cormen", Price: 129.00},
		{Title: "设计模式", Author: "GoF", Price: 45.00},
		{Title: "重构", Author: "Martin Fowler", Price: 69.00},
		{Title: "代码大全", Author: "Steve McConnell", Price: 89.00},
		{Title: "程序员修炼之道", Author: "Andy Hunt", Price: 55.00},
	}

	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO books (title, author, price) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, book := range testData {
		_, err := stmt.Exec(book.Title, book.Author, book.Price)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	fmt.Println("书籍测试数据插入成功")
	return nil
}

// queryExpensiveBooks 查询价格大于指定金额的书籍
func queryExpensiveBooks(db *sqlx.DB, minPrice float64) ([]Book, error) {
	var books []Book
	query := "SELECT id, title, author, price FROM books Where price > ? ORDER BY price DESC"
	err := db.Select(&books, query, minPrice)
	if err != nil {
		return nil, err
	}
	return books, nil
}

// queryBooksByAuthor 根据作者查询书籍
func queryBooksByAuthor(db *sqlx.DB, author string) ([]Book, error) {
	var books []Book
	query := "SELECT id, title, author, price FROM books WHERE author LIKE ?"
	err := db.Select(&books, query, "%"+author+"%")
	if err != nil {
		return nil, err
	}
	return books, nil
}

// queryBookByID 根据ID查询书籍
func queryBookByID(db *sqlx.DB, id int) (*Book, error) {
	var book Book
	query := "SELECT id, title, author, price FROM books WHERE id = ?"
	err := db.Get(&book, query, id)
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func main() {
	// 连接 SQLite 书籍库
	db, err := sqlx.Open("sqlite3", "books.db")
	if err != nil {
		fmt.Println("连接数据库失败：", err)
		return
	}
	defer db.Close()

	// 插入书籍测试书籍
	err = insertBookTestData(db)
	if err != nil {
		fmt.Println("插入书籍测试数据失败：", err)
	}

	// 查询价格大于50元的书籍
	expensiveBooks, err := queryExpensiveBooks(db, 50)
	if err != nil {
		fmt.Println("查询价格大于50元的书籍失败：", err)
	}
	fmt.Println("价格大于50元的书籍：")
	for _, book := range expensiveBooks {
		fmt.Printf("ID: %d, Title: %s, Author: %s, Price: %.2f\n", book.Id, book.Title, book.Author, book.Price)
	}

	// 根据作者查询书籍
	authorBooks, err := queryBooksByAuthor(db, "GoF")
	if err != nil {
		fmt.Println("根据作者查询书籍失败：", err)
	}
	fmt.Println("作者为GoF的书籍：")
	for _, book := range authorBooks {
		fmt.Printf("ID: %d, Title: %s, Author: %s, Price: %.2f\n", book.Id, book.Title, book.Author, book.Price)
	}

	// 根据IP查询书籍
	specificBook, err := queryBookByID(db, 3)
	if err != nil {
		fmt.Println("根据ID查询书籍失败：", err)
	}
	fmt.Printf("ID为3的书籍：\nID: %d, Title: %s, Author: %s, Price: %.2f\n", specificBook.Id, specificBook.Title, specificBook.Author, specificBook.Price)
}
