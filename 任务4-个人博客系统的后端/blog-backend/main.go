package main

import (
	"log"
	"net/http"

	"github.com/18850341851/blog-backend/db"
	"github.com/18850341851/blog-backend/middleware"
	"github.com/18850341851/blog-backend/models"
	"github.com/18850341851/blog-backend/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	//初始化数据库连接
	if err := db.InitDB(); err != nil {
		log.Fatalf("数据库连接失败：%v", err)

	}
	log.Println("数据库连接成功")

	//执行数据库迁移
	if err := models.MigrateAll(); err != nil {
		log.Fatalf("数据库迁移失败：%v", err)
	}
	log.Println("数据库迁移成功")

	//初始化Gin引擎
	r := gin.Default()
	//应用错误处理中间件
	r.Use(middleware.ErrorHandler())

	//初始化路由
	routes.InitRoutes(r)

	//启动服务
	log.Println("服务器启动，监听端口：8080")
	if err := r.Run(":8080"); err != nil && err != http.ErrServerClosed {
		log.Fatalf("服务器启动失败：%v", err)
	}
}
