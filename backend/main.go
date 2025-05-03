package main

import (
	"log"

	"github.com/alvinhmg/blog/config"
	"github.com/alvinhmg/blog/routes"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/joho/godotenv"
)

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

	// 注册路由
	routes.RegisterRoutes(h)

	// 启动服务器
	h.Spin()
}
