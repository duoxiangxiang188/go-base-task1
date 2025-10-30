package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //PostgreSQL驱动
)

// Employee 结构体定义，与数据表字段映射
type Employee struct {
	ID         int    `db:"id"`
	Name       string `db:"name"`
	Department string `db:"department"`
	Salary     int    `db:"salary"`
}

func main() {
	//PostgreSQL连接字符串
	//格式:hose=主机名 port = 端口  user = 用户名  password = 密码  dbname = 密码  sslmode = 是否启动ssl
	connStr := "host=localhost port=5432 user=postgres password=888888 dbname=runoobdb sslmode=disable"
	//连接数据库
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("无法打开数据库连接: %v", err)
	}
	defer db.Close()
	//测试连接
	err = db.Ping()
	if err != nil {
		log.Fatalf("数据库连接失败：%v", err)
	}
	//创建employees表（如果不存在）
	err = createEmployeesTable(db)
	if err != nil {
		log.Fatalf("创建表失败： %v", err)
	}
	fmt.Println("表检查/创建完成")

	//检查是否有数据，若无则插入示例数据
	err = insertSampleDataIfEmpty(db)
	if err != nil {
		log.Printf("插入示例数据失败：%v", err)
	} else {
		fmt.Println("插入示例数据完成")
	}

	//1.查询所有部门为“技术部”的员工
	techEmployees, err := getTechDepartmentEmployees(db)
	if err != nil {
		log.Printf("查询技术部员工失败：%v", err)

	} else {
		fmt.Println("\n技术部员工列表：")
		for _, emp := range techEmployees {
			fmt.Printf("ID: %d, 姓名：%s, 部门：%s, 工资：%d\n", emp.ID, emp.Name, emp.Department, emp.Salary)
		}
	}
	//2.查询工资最高的员工
	topSalaryEmp, err := getHighestSalaryEmployees(db)
	if err != nil {
		log.Fatalf("查询工资最高 员工失败：%v", err)
	} else {
		fmt.Printf("\n工资最高员工: ID: %d, 姓名： %s, 部门：%s, 工资：%d\n", topSalaryEmp.ID, topSalaryEmp.Name, topSalaryEmp.Department, topSalaryEmp.Salary)
	}
}

// 创建employees表（如果不存在）
func createEmployeesTable(db *sqlx.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS employees (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		department VARCHAR(50) NOT NULL,
		salary INTEGER NOT NULL
	
	);`
	_, err := db.Exec(createTableSQL)
	return err
}

// insertSampleDataIfEmpty当表为空时插入示例数据
func insertSampleDataIfEmpty(db *sqlx.DB) error {
	//先检查表中是否有数据
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM employees")
	if err != nil {
		return err
	}
	if count > 0 {
		fmt.Println("表中已经有数据，不插入示例数据")
		return nil
	}
	//准备示例数据
	sampleEmployees := []Employee{
		{Name: "张三", Department: "技术部", Salary: 8000},
		{Name: "李四", Department: "技术部", Salary: 9500},
		{Name: "王五", Department: "市场部", Salary: 7500},
		{Name: "赵六", Department: "技术部", Salary: 12000},
		{Name: "孙七", Department: "人事部", Salary: 6800},
		{Name: "周八", Department: "市场部", Salary: 8800},
	}
	//开始事务批量插入
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback() //确保事务回滚
	//插入每条数据
	for _, emp := range sampleEmployees {
		_, err := tx.NamedExec(`
		INSERT INTO employees (name, department, salary)
		VALUES (:name, :department, :salary)`, emp)
		if err != nil {
			return err
		}
	}
	//提交事务
	return tx.Commit()

}

// 查询所有技术部员工
func getTechDepartmentEmployees(db *sqlx.DB) ([]Employee, error) {
	var employees []Employee
	err := db.Select(&employees, "SELECT id, name, department, salary FROM employees WHERE department = $1", "技术部")
	return employees, err
}

// 查询工资最高的员工
func getHighestSalaryEmployees(db *sqlx.DB) (Employee, error) {
	var emp Employee
	err := db.Get(&emp, `
		SELECT id, name, department, salary
		FROM employees
		WHERE salary = (SELECT MAX(salary) FROM employees)
		LIMIT 1
	`)
	return emp, err
}
