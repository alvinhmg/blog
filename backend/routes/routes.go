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
	// registerUserRoutes(apiGroup) // User routes removed
	registerPostRoutes(apiGroup)
	registerCategoryRoutes(apiGroup)
	registerTagRoutes(apiGroup)
	registerCommentRoutes(apiGroup)
	registerMiscRoutes(apiGroup) // 添加杂项路由注册
}

// 认证相关路由
func registerAuthRoutes(group *route.RouterGroup) {
	auth := group.Group("/auth")

	// 连接认证控制器
	auth.POST("/register", api.Register) // 用户注册
	auth.POST("/login", api.Login)       // 用户登录
	auth.POST("/logout", api.Logout)     // 用户登出
}

// 用户相关路由 (示例 - 功能已移除)
// func registerUserRoutes(group *route.RouterGroup) {
// 	users := group.Group("/users", middleware.JWTAuth()) // Add JWTAuth middleware
//
// 	users.GET("", api.GetUsers)          // 获取用户列表
// 	users.GET("/:id", api.GetUser)       // 获取用户详情
// 	users.PUT("/:id", api.UpdateUser)    // 更新用户信息
// 	users.DELETE("/:id", api.DeleteUser) // 删除用户
// }

// 文章相关路由
func registerPostRoutes(group *route.RouterGroup) {
	posts := group.Group("/posts")

	posts.GET("", api.GetPosts)                                 // 获取文章列表
	posts.GET("/search", api.SearchPosts)                       // 全文搜索文章
	posts.GET("/:id", api.GetPost)                              // 获取文章详情
	posts.POST("/:id/like", middleware.JWTAuth(), api.LikePost) // 点赞文章

	// 管理员权限路由
	adminPosts := posts.Group("", middleware.JWTAuth(), middleware.AdminAuth())
	adminPosts.POST("", api.CreatePost)       // 创建文章
	adminPosts.PUT("/:id", api.UpdatePost)    // 更新文章
	adminPosts.DELETE("/:id", api.DeletePost) // 删除文章
}

// 分类相关路由
func registerCategoryRoutes(group *route.RouterGroup) {
	categories := group.Group("/categories")
	categories.GET("", api.GetCategories)        // 公开获取分类列表
	categories.GET("/hot", api.GetHotCategories) // 获取热门分类
	categories.GET("/:id", api.GetCategory)      // 公开获取分类详情

	// 管理员权限路由
	adminCategories := categories.Group("", middleware.JWTAuth(), middleware.AdminAuth())
	adminCategories.POST("", api.CreateCategory)
	adminCategories.PUT("/:id", api.UpdateCategory)
	adminCategories.DELETE("/:id", api.DeleteCategory)
}

// 标签相关路由
func registerTagRoutes(group *route.RouterGroup) {
	tags := group.Group("/tags")
	tags.GET("", api.GetTags)        // 公开获取标签列表
	tags.GET("/hot", api.GetHotTags) // 获取热门标签
	tags.GET("/:id", api.GetTag)     // 公开获取标签详情

	// 管理员权限路由
	adminTags := tags.Group("", middleware.JWTAuth(), middleware.AdminAuth())
	adminTags.POST("", api.CreateTag)
	adminTags.PUT("/:id", api.UpdateTag)
	adminTags.DELETE("/:id", api.DeleteTag)
}

// 评论相关路由
func registerCommentRoutes(group *route.RouterGroup) {
	comments := group.Group("/comments") // 评论相关路由

	// 创建评论 (需要登录)
	comments.POST("/post/:postId", middleware.JWTAuth(), api.AddComment)
	// comments.GET("", api.GetComments) // 获取评论列表 (通常在文章详情中获取)
	// comments.GET("/:id", api.GetComment) // 获取单个评论详情
	// comments.PUT("/:id", api.UpdateComment) // 更新评论 (通常不允许用户更新)
	// comments.DELETE("/:id", api.DeleteComment) // 删除评论 (用户或管理员)

	// 管理员权限路由
	adminComments := group.Group("/admin/comments", middleware.JWTAuth(), middleware.AdminAuth())
	adminComments.PUT("/:id/approve", api.ApproveComment) // 审核通过评论
	// adminComments.PUT("/:id/reject", api.RejectComment)   // TODO: 实现拒绝评论功能
}

// 杂项路由 (首页数据、归档等)
func registerMiscRoutes(group *route.RouterGroup) {
	group.GET("/home", api.GetHomePageData)   // 获取首页数据
	group.GET("/archive", api.GetArchiveData) // 获取归档数据
}
