package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Student struct {
	ID    int
	Name  string
	Age   int
	Grade string
}

func main() {
	// 连接数据库

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "888888"),
		getEnv("DB_NAME", "runoobdb"),
	)
	db, err := sql.Open("postgres", connStr)
	// db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("无法打开数据库连接: %v", err)
	}
	defer db.Close()

	// 验证连接
	if err := db.Ping(); err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}
	fmt.Println("成功连接到PostgreSQL数据库")

	// 确保表存在
	if err := createStudentsTable(db); err != nil {
		log.Fatalf("创建表失败: %v", err)
	}

	// 1. 插入数据
	insertSQL := `INSERT INTO students (name, age, grade) VALUES ($1, $2, $3) RETURNING id`

	var newID int
	err = db.QueryRow(insertSQL, "张三", 20, "三年级").Scan(&newID)
	if err != nil {
		log.Printf("插入数据失败: %v", err)
	} else {
		fmt.Printf("插入成功，新纪录ID：%d\n", newID)
	}

	// 2. 查询年龄大于18岁的学生
	selectSQL := `SELECT id, name, age, grade FROM students WHERE age > $1`
	rows, err := db.Query(selectSQL, 18)
	if err != nil {
		log.Printf("查询失败: %v", err)
	} else {
		defer rows.Close()
		fmt.Println("年龄大于18岁的学生:")
		for rows.Next() {
			var s Student
			if err := rows.Scan(&s.ID, &s.Name, &s.Age, &s.Grade); err != nil {
				log.Printf("扫描结果失败：%v", err)
				continue
			}
			fmt.Printf("ID:%d, 姓名：%s, 年龄：%d, 年级： %s\n", s.ID, s.Name, s.Age, s.Grade)
		}
	}

	// 3. 更新张三的年级为四年级
	updateSQL := `UPDATE students SET grade = $1 WHERE name = $2`
	result, err := db.Exec(updateSQL, "四年级", "张三")
	if err != nil {
		log.Printf("更新失败：%v", err)
	} else {
		rowsAffected, _ := result.RowsAffected()
		fmt.Printf("更新成功，影响行数： %d\n", rowsAffected)
	}

	// 4. 删除年龄小于15岁的学生
	deleteSQL := `DELETE FROM students WHERE age < $1`
	result, err = db.Exec(deleteSQL, 15)
	if err != nil {
		log.Printf("删除失败：%v", err)
	} else {
		rowsAffected, _ := result.RowsAffected()
		fmt.Printf("删除成功，影响行数：%d\n", rowsAffected)
	}

}

// 创建students表
func createStudentsTable(db *sql.DB) error {
	createSQL := `
	CREATE TABLE IF NOT EXISTS students (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		age INT NOT NULL,
		grade VARCHAR(50) NOT NULL
	);`
	_, err := db.Exec(createSQL)
	return err
}

// 获取环境变量，不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
