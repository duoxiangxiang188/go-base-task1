package models

import (
	"github.com/18850341851/blog-backend/db"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 用户模型（对应users表）
type User struct {
	gorm.Model           // 内置字段：ID、CreatedAt、UpdatedAt、DeletedAt（软删除）
	Username   string    `gorm:"size:50;not null;unique" json:"username"`     //用户名（唯一）
	Password   string    `gorm:"size:100;not null" json:"-"`                  // 密码（加密存储，JSON返回时隐藏）
	Email      string    `gorm:"size:100;not null;unique" json:"email"`       // 邮箱（唯一）
	Posts      []Post    `gorm:"foreignKey:UserID" json:"posts,omitempty"`    // 关联的文章（一对多）
	Comments   []Comment `gorm:"foreignKey:UserID" json:"comments,omitempty"` // 关联的评论（一对多）
}

// 钩子方法：创建用户前加密密码
func (u *User) BeforeCreate(tx *gorm.DB) error {
	//使用bcrypt 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// 验证密码是否正确
func (u *User) CheckPassword(password string) bool {
	//对比明文密码与加密密码
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// 迁移用户表
func MigrateUsers() error {
	return db.DB.AutoMigrate(&User{})
}
