package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Book struct {
	ID     int     `db:"id"`
	Title  string  `db:"title"`
	Author string  `db:"author"`
	Price  float64 `db:"price"`
}

func main() {
	connStr := "host=localhost port=5432 user=postgres password=5627951 dbname=runoobdb sslmode=disable"
	//连接数据库
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("无法打开数据库:%v", err)
	}
	defer db.Close()

	//测试数据库
	err = db.Ping()
	if err != nil {
		log.Fatalf("数据库连接失败:%v", err)
	}

	//创建books表（如果不存在）
	err = createBooksTable(db)
	if err != nil {
		log.Fatalf("创建表失败：%v", err)
	}
	fmt.Println("检查/创建表完成")
	err = insertSampleBooks(db)
	if err != nil {
		log.Printf("插入示例书籍失败：%v", err)
	} else {
		fmt.Println("示例书籍数据插入完成")
	}

	//查询价格大于50元的价格
	expensiveBooks, err := getBooksPriceOver(db, 50.0)
	if err != nil {
		log.Printf("查询高价书籍失败：%v", err)
	} else {
		fmt.Println("\n价格大于50元的书籍 ：")
		for _, book := range expensiveBooks {
			fmt.Printf("ID: %d, 书名：%s, 作者：%s, 价格：%.2f\n", book.ID, book.Title, book.Author, book.Price)
		}
	}
}

// createBooksTable 创建books表（如果不存在）
func createBooksTable(db *sqlx.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS books (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		author VARCHAR(100) NOT NULL,
		price NUMERIC(10, 2) NOT NULL  --支持两位小数的价格
		);`
	_, err := db.Exec(createTableSQL)
	return err
}

func insertSampleBooks(db *sqlx.DB) error {
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM books")
	if err != nil {
		return err
	}
	if count > 0 {
		fmt.Println("books表中已有数据，不插入示例数据")
		return nil
	}

	//示例书籍数据
	sampleBooks := []Book{
		{Title: "Go语言编程", Author: "张三", Price: 89.00},
		{Title: "PostgresSQL实战", Author: "李四", Price: 79.50},
		{Title: "SQL入门", Author: "王五", Price: 45.00},
		{Title: "数据结构和算法", Author: "赵六", Price: 99.00},
		{Title: "Python教程", Author: "孙七", Price: 39.00},
	}
	//使用事务批量插入数据
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	//使用NameExec进行类型安全插入
	for _, book := range sampleBooks {
		_, err := tx.NamedExec(`
		INSERT INTO books (title, author, price)
		VALUES (:title, :author, :price)`, book)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// getBooksPriceOver 查询价格大于指定金额的书籍
// 确保类型安全：参数和返回值都是用明确的类型
func getBooksPriceOver(db *sqlx.DB, minPrice float64) ([]Book, error) {
	var books []Book
	//使用参数化查询，避免SQL注入并确保明确的类型正确映射
	query := "SELECT id, title, author, price FROM books WHERE price >$1 ORDER BY price DESC"

	//Select方法会自动进行类型映射，确保类型安全
	err := db.Select(&books, query, minPrice)
	return books, err

}
