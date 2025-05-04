package routes

import (
	"github.com/alvinhmg/blog/api"
	"github.com/alvinhmg/blog/middleware"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/route"
)

// RegisterRoutes 注册所有API路由
func RegisterRoutes(h *server.Hertz) {
	// API路由组
	apiGroup := h.Group("/api") // Renamed to avoid conflict with package name

	// 注册各模块路由
	registerAuthRoutes(apiGroup)
	registerUserRoutes(apiGroup)
	registerPostRoutes(apiGroup)
	registerCategoryRoutes(apiGroup)
	registerTagRoutes(apiGroup)
	registerCommentRoutes(apiGroup)
}

// 认证相关路由
func registerAuthRoutes(group *route.RouterGroup) {
	auth := group.Group("/auth")

	// 连接认证控制器
	auth.POST("/register", api.Register) // 用户注册
	auth.POST("/login", api.Login)       // 用户登录
	auth.POST("/logout", api.Logout)     // 用户登出
}

// 用户相关路由
func registerUserRoutes(group *route.RouterGroup) {
	users := group.Group("/users", middleware.JWTAuth()) // Add JWTAuth middleware

	users.GET("", api.GetUsers)          // 获取用户列表
	users.GET("/:id", api.GetUser)       // 获取用户详情
	users.PUT("/:id", api.UpdateUser)    // 更新用户信息
	users.DELETE("/:id", api.DeleteUser) // 删除用户
}

// 分类相关路由
func registerCategoryRoutes(group *route.RouterGroup) {
	categories := group.Group("/categories", middleware.JWTAuth(), middleware.AdminAuth()) // Use AdminAuth
	categories.GET("", api.GetCategories)
	categories.GET("/:id", api.GetCategory)
	categories.POST("", api.CreateCategory)
	categories.PUT("/:id", api.UpdateCategory)
	categories.DELETE("/:id", api.DeleteCategory)
}

// 标签相关路由
func registerTagRoutes(group *route.RouterGroup) {
	tags := group.Group("/tags", middleware.JWTAuth(), middleware.AdminAuth()) // Use AdminAuth
	tags.GET("", api.GetTags)
	tags.GET("/:id", api.GetTag)
	tags.POST("", api.CreateTag)
	tags.PUT("/:id", api.UpdateTag)
	tags.DELETE("/:id", api.DeleteTag)
}

// 文章相关路由
func registerPostRoutes(group *route.RouterGroup) {
	posts := group.Group("/posts")
	posts.GET("", api.GetPosts)
	posts.GET("/:id", api.GetPost)

	adminPosts := posts.Group("", middleware.JWTAuth(), middleware.AdminAuth()) // Use AdminAuth
	adminPosts.POST("", api.CreatePost)
	adminPosts.PUT("/:id", api.UpdatePost)
	adminPosts.DELETE("/:id", api.DeletePost)
}

// 评论相关路由
func registerCommentRoutes(group *route.RouterGroup) {
	comments := group.Group("/comments", middleware.JWTAuth())
	comments.GET("/post/:postId", api.GetPostComments)
	comments.POST("/post/:postId", api.AddComment)
	comments.DELETE("/:id", api.DeleteComment)

	adminComments := comments.Group("/admin", middleware.JWTAuth(), middleware.AdminAuth()) // Add JWTAuth and use AdminAuth
	adminComments.GET("/pending", api.GetPendingComments)
	adminComments.POST("/approve/:id", api.ApproveComment)
}
