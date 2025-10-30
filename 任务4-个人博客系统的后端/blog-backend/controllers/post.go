package controllers

import (
	"net/http"

	"github.com/18850341851/blog-backend/db"
	"github.com/18850341851/blog-backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 文章创建请求体
type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required,min=1"`
}

// 文章更新请求体
type UpdatePostRequest struct {
	Title   string `json:"title" binding:"omitempty,min=1,max=200"`
	Content string `json:"content" binding:"omitempty,min=1"`
}

// 创建文章
func CreatePost(c *gin.Context) {
	//从上下文获得当前用户ID
	userID, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误：" + err.Error()})
		return
	}

	//创建文章
	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID.(uint),
	}
	if err := db.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文章失败" + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "文章创建成功",
		"post":    post,
	})
}

// 获取所有文章列表
func GetAllPosts(c *gin.Context) {
	var posts []models.Post
	//关联查询用户信息（不包含密码）
	if err := db.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID, Username, Email, created_at")
	}).Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询文章失败：" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

// 获取单篇文章详情
func GetPost(c *gin.Context) {
	postID := c.Param("id")
	var post models.Post

	//关联查询用户和评论（评论关联用户）
	if err := db.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID, Username")
	}).Preload("Comments.User", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID, Username")
	}).First(&post, postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询文章失败：" + err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"post": post})
}

// 更新文章
func UpdatePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	postID := c.Param("id")
	var post models.Post
	if err := db.DB.First(&post, postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "文章查询失败：" + err.Error()})
		}
		return
	}

	//验证作者身份
	if post.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权修改此文章"})
		return
	}
	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误：" + err.Error()})
		return
	}

	//更新文章
	if err := db.DB.Model(&post).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新文章失败：" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "文章更新成功",
		"post":    post,
	})
}

// 删除文章（仅作者可操作）
func DeletePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	postID := c.Param("id")
	var post models.Post
	if err := db.DB.First(&post, postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})

		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询文章失败：" + err.Error()})
		}
		return
	}

	//验证作者身份
	if post.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权删除此文章"})
		return
	}
	//删除文章（软删除）
	if err := db.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文章失败：" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "文章删除成功"})
}
