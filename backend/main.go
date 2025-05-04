package main

import (
	"context"
	"log"

	"github.com/alvinhmg/blog/config"
	"github.com/alvinhmg/blog/routes"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/joho/godotenv"
)

// 定义CORS中间件函数
func corsMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Response.Header.Set("Access-Control-Allow-Origin", "*")
		c.Response.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Response.Header.Set("Access-Control-Max-Age", "86400")

		// 处理 OPTIONS 预检请求
		if string(c.Method()) == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next(ctx)
	}
}

func main() {
	// 加载环境变量
	err := godotenv.Load()
	if err != nil {
		log.Println("未找到.env文件，将使用默认配置")
	}

	// 初始化数据库连接
	config.InitDB()

	// 创建Hertz服务器实例
	h := server.Default(
		server.WithHostPorts(":8080"),
	)
	h.Use(corsMiddleware())

	// 注册路由
	routes.RegisterRoutes(h)

	// 启动服务器
	h.Spin()
}
