package routes

import (
	"github.com/alvinhmg/blog/api"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/route"
)

// RegisterRoutes 注册所有API路由
func RegisterRoutes(h *server.Hertz) {
	// API路由组
	api := h.Group("/api")

	// 注册各模块路由
	registerAuthRoutes(api)
	registerUserRoutes(api)
	registerPostRoutes(api)
	registerCategoryRoutes(api)
	registerTagRoutes(api)
	registerCommentRoutes(api)
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
	users := group.Group("/users")

	// TODO: 实现用户控制器
	users.GET("", nil)        // 获取用户列表
	users.GET("/:id", nil)    // 获取用户详情
	users.PUT("/:id", nil)    // 更新用户信息
	users.DELETE("/:id", nil) // 删除用户
}

// 文章相关路由
func registerPostRoutes(group *route.RouterGroup) {
	posts := group.Group("/posts")

	// 连接文章控制器
	posts.GET("", api.GetPosts)    // 获取文章列表
	posts.GET("/:id", api.GetPost) // 获取文章详情
	posts.POST("", nil)            // 创建文章 (待实现)
	posts.PUT("/:id", nil)         // 更新文章 (待实现)
	posts.DELETE("/:id", nil)      // 删除文章 (待实现)
	posts.POST("/:id/like", nil)   // 点赞文章 (待实现)
	posts.GET("/search", nil)      // 搜索文章 (待实现)
}

// 分类相关路由
func registerCategoryRoutes(group *route.RouterGroup) {
	categories := group.Group("/categories")

	// TODO: 实现分类控制器
	categories.GET("", nil)           // 获取分类列表
	categories.GET("/:id", nil)       // 获取分类详情
	categories.POST("", nil)          // 创建分类
	categories.PUT("/:id", nil)       // 更新分类
	categories.DELETE("/:id", nil)    // 删除分类
	categories.GET("/:id/posts", nil) // 获取分类下的文章
}

// 标签相关路由
func registerTagRoutes(group *route.RouterGroup) {
	tags := group.Group("/tags")

	// TODO: 实现标签控制器
	tags.GET("", nil)           // 获取标签列表
	tags.GET("/:id", nil)       // 获取标签详情
	tags.POST("", nil)          // 创建标签
	tags.PUT("/:id", nil)       // 更新标签
	tags.DELETE("/:id", nil)    // 删除标签
	tags.GET("/:id/posts", nil) // 获取标签下的文章
}

// 评论相关路由
func registerCommentRoutes(group *route.RouterGroup) {
	comments := group.Group("/comments")

	// 连接评论控制器
	comments.GET("/post/:postId", api.GetPostComments) // 获取文章下的评论
	comments.DELETE("/:id", api.DeleteComment)         // 删除评论

	// 文章评论路由
	group.POST("/posts/:postId/comments", api.AddComment) // 添加评论到指定文章
}
