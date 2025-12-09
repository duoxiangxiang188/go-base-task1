package models

import (
	"github.com/18850341851/blog-backend/db"
	"gorm.io/gorm"
)

// Comment 评论模型（对应comment表）
type Comment struct {
	gorm.Model        // 内置字段：ID、CreatedAt、UpdatedAt、DeletedAt
	Content    string `gorm:"type:text;not null" json:"content"`       //评论内容
	UserID     uint   `gorm:"not null" json:"user_id"`                 //关联的用户ID（外键）
	User       User   `gorm:"foreignKey:UserID" json:"user"`           //关联的用户信息（一对一）
	PostID     uint   `gorm:"not null" json:"post_id"`                 //关联的文章ID（外键）
	Post       Post   `gorm:"foreignKey:PostID" json:"post,omitempty"` //关联的文章（一对一，可选）

}

func MigrateComments() error {
	return db.DB.AutoMigrate(&Comment{})
}
