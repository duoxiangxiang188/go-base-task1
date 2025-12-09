package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// 转账函数：从fromAcoount向toAccount转账amount金额
func transfer(db *sql.DB, fromAccountID, toAccountID int, amount float64) error {
	//开始事务
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("无法开始事务： %v", err)
	}
	defer func() {
		//发送错误时回滚事务
		if r := recover(); r != nil || err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("回滚事务失败 ： %v", rollbackErr)
			}
		}
	}()

	//1.检查转出账户余额是否充足
	var balance float64
	checkSQL := "SELECT balance FROM accounts WHERE id = $1 FOR UPDATE" // FOR UPDATE 加行锁 防止并发问题
	err = tx.QueryRow(checkSQL, fromAccountID).Scan(&balance)
	if err != nil {
		return fmt.Errorf("查询账户余额失败: %v", err)
	}
	//检查余额是否足够
	if balance < amount {
		return fmt.Errorf("余额不足，当前余额: %.2f, 需要转账: %.2f", balance, amount)
	}

	//2.从转出账户扣减余额
	updateFromSQL := "UPDATE accounts SET balance = balance - $1 WHERE id = $2"
	_, err = tx.Exec(updateFromSQL, amount, fromAccountID)
	if err != nil {
		return fmt.Errorf("扣减转出账户金额失败:%v", err)
	}

	//3.向转入账户增加金额
	updateToSQL := "UPDATE accounts SET balance = balance + $1 WHERE id = $2"
	_, err = tx.Exec(updateToSQL, amount, toAccountID)
	if err != nil {
		return fmt.Errorf("增加转入账户金额失败: %v", err)
	}

	//4.记录交易信息
	insertTxSQL := `
	INSERT INTO transactions
	(from_account_id, to_account_id, amount)
	VALUES($1, $2, $3)`
	_, err = tx.Exec(insertTxSQL, fromAccountID, toAccountID, amount)
	if err != nil {
		return fmt.Errorf("记录交易信息失败：%v", err)
	}

	//提交事务
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}
	return nil
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

	if err != nil {
		log.Fatalf("无法打开数据库连接: %v", err)
	}
	defer db.Close()

	// 验证连接
	if err := db.Ping(); err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}
	fmt.Println("成功连接到PostgreSQL数据库")
	//初始化表结构
	initTables(db)

	//示例：从账户1 向账户2转账100元
	fromAccountID := 1
	toAccountID := 2
	amount := 100.0
	//执行转账
	err = transfer(db, fromAccountID, toAccountID, amount)
	if err != nil {
		log.Printf("转账失败: %v", err)
	} else {
		fmt.Printf("成功从账户%d向账户%d转账%.2f元\n", fromAccountID, toAccountID, amount)
		//打印转账后的余额
		printAccountBalance(db, fromAccountID)
		printAccountBalance(db, toAccountID)
	}
}

// 初始化表结构
func initTables(db *sql.DB) {
	//创建账户表
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS accounts(
		id SERIAL PRIMARY KEY,
		balance NUMERIC(10, 2) NOT NULL DEFAULT 0.00 CHECK (balance >= 0)
	);`)
	if err != nil {
		log.Fatalf("创建账户表失败： %v", err)
	}
	//创建交易记录表
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS transactions (
		id SERIAL PRIMARY KEY,
		from_account_id INT NOT NULL,
		to_account_id INT NOT NULL,
		amount NUMERIC(10, 2) NOT NULL CHECK(amount >=0),
		create_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (from_account_id) REFERENCES accounts(id),
		FOREIGN KEY(to_account_id) REFERENCES accounts(id)
	);`)
	if err != nil {
		log.Fatalf("创建交易表失败： %v", err)
	}
	_, err = db.Exec(`INSERT INTO accounts (id, balance) VALUES(1, 800.00) ON CONFLICT(id) DO NOTHING`)
	if err != nil {
		log.Fatalf("插入测试账户1失败 ： %v", err)
	}
	_, err = db.Exec(`INSERT INTO accounts (id, balance) VALUES(2, 300.00)ON CONFLICT(id) DO NOTHING`)
	if err != nil {
		log.Fatalf("插入测试账户2失败 ： %v", err)
	}
}

// 打印账户余额
func printAccountBalance(db *sql.DB, accountID int) {
	var balance float64
	err := db.QueryRow("SELECT balance FROM accounts WHERE id = $1", accountID).Scan(&balance)
	if err != nil {
		log.Printf("查询账户 %d 余额失败 ： %v", accountID, err)
	}
	fmt.Printf("账户%d 当前余额：%.2f元\n", accountID, balance)
}

// 获取环境变量，不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
