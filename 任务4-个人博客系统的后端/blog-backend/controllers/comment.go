package controllers

import (
	"net/http"

	"github.com/18850341851/blog-backend/db"
	"github.com/18850341851/blog-backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 评论创建请求体
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1"`
}

// 创建评论
func CreateComment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	postID := c.Param("postId")
	var post models.Post
	if err := db.DB.First(&post, postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询文章失败：" + err.Error()})
		}
		return
	}
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误：" + err.Error()})
		return
	}
	//创建评论

	comment := models.Comment{
		Content: req.Content,
		UserID:  userID.(uint),
		PostID:  post.ID,
	}
	if err := db.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建评论失败" + err.Error()})
		return

	}

	//关联查询用户信息
	db.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID, Username")
	}).First(&comment, comment.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "评论创建成功",
		"comment": comment,
	})
}

// 获取文章的所有评论
func GetPostComments(c *gin.Context) {
	postID := c.Param("id") //这个id要跟routes里的public.GET("/posts/:id/comments", controllers.GetPostComments) 的:id保持一致
	var comments []models.Comment
	//验证文章存在

	var post models.Post
	if err := db.DB.First(&post, postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询文章失败" + err.Error()})
		}
		return
	}

	//查询评论并关联用户
	if err := db.DB.Where("post_id = ?", postID).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID, Username")
	}).Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询评论失败：" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"comments": comments})
}
