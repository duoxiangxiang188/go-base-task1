package main

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// User 模型定义（用户）
type User struct {
	gorm.Model        // 包含ID、CreatedAt、UpdatedAt、DeletedAt字段
	Username   string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"` // 用户名
	Email      string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`   // 邮箱
	Posts      []Post `gorm:"foreignKey:UserID" json:"posts"`                        // 关联的文章列表
}

// Post 模型定义（文章）
type Post struct {
	gorm.Model           // 包含基础字段
	Title      string    `gorm:"type:varchar(200);not null" json:"title"` // 文章标题
	Content    string    `gorm:"type:text" json:"content"`                // 文章内容
	UserID     uint      `gorm:"not null" json:"user_id"`                 // 外键：关联用户ID
	User       User      `gorm:"foreignKey:UserID" json:"user"`           // 反向关联用户
	Comments   []Comment `gorm:"foreignKey:PostID" json:"comments"`       // 关联的评论列表
}

// Comment 模型定义（评论）
type Comment struct {
	gorm.Model        // 包含基础字段
	Content    string `gorm:"type:varchar(500);not null" json:"content"` // 评论内容
	PostID     uint   `gorm:"not null" json:"post_id"`                   // 外键：关联文章ID
	Post       Post   `gorm:"foreignKey:PostID" json:"post"`             // 反向关联文章
	UserID     uint   `gorm:"not null" json:"user_id"`                   // 外键：关联用户ID（评论作者）
	User       User   `gorm:"foreignKey:UserID" json:"user"`             // 反向关联用户
}

func main() {
	// MySQL 连接配置
	// 格式：用户名:密码@tcp(主机:端口)/数据库名?charset=utf8mb4&parseTime=True&loc=Local
	dsn := "root:8888888@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"

	// 连接 MySQL 数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 打印SQL日志
	})
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}

	// 获取底层sql.DB对象，配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		panic("获取数据库连接池失败: " + err.Error())
	}
	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大存活时间

	// 自动迁移创建表结构
	// 会根据模型定义创建users、posts、comments表，并自动添加外键约束
	err = db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		panic("表结构迁移失败: " + err.Error())
	}
}
