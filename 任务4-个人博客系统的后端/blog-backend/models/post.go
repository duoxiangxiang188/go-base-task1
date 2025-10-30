package models

import (
	"github.com/18850341851/blog-backend/db"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model           // 内置字段：ID、CreatedAt、UpdatedAt、DeletedAt
	Title      string    `gorm:"size:200;not null" json:"title"`              //文章标题
	Content    string    `gorm:"type:text;not null" json:"content"`           //文章内容
	UserID     uint      `gorm:"not null" json:"user_id"`                     //关联的用户ID（外键)
	User       User      `gorm:"foreignKey:UserID" json:"user"`               //关联用户信息（一对一）
	Comments   []Comment `gorm:"foreignKey:PostID" json:"comments,omitempty"` //关联的评论（一对多）
}

// 迁移文章表
func MigratePosts() error {
	return db.DB.AutoMigrate(&Post{})
}
