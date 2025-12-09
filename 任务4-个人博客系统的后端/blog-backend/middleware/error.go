package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 错误处理中间件（捕捉panic并返回统一格式错误）
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				//记录错误日志
				log.Printf("Panic recovered: %v", err)
				//返回500错误
				c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
				c.Abort()
			}
		}()
		c.Next()
	}
}
