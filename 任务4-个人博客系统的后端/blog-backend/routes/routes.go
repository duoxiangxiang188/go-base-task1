package routes

import (
	"github.com/18850341851/blog-backend/controllers"
	"github.com/18850341851/blog-backend/middleware"
	"github.com/gin-gonic/gin"
)

// 初始化路由
func InitRoutes(r *gin.Engine) {
	//公开路由（无需认证）
	public := r.Group("/api")
	{
		//用户相关
		public.POST("/register", controllers.Register)
		public.POST("/login", controllers.Login)

		//文章相关（读操作公开）
		public.GET("/posts", controllers.GetAllPosts)
		public.GET("/posts/:id", controllers.GetPost)

		//评论相关（读操作公开）
		public.GET("/posts/:id/comments", controllers.GetPostComments)
	}

	//需要认证的路由

	protected := r.Group("/api")
	protected.Use(middleware.JWTAuthMiddleware()) //应用JWT认证中间件
	{
		//文章相关（写操作需要认证）
		protected.POST("/posts", controllers.CreatePost)
		protected.PUT("/posts/:id", controllers.UpdatePost)
		protected.DELETE("/posts/:id", controllers.DeletePost)

		//评论相关（写操作需要认证）
		protected.POST("/posts/:postId/comments", controllers.CreateComment)
	}
}
