package middleware

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

// RegisterStaticFileServer 注册静态文件服务中间件
func RegisterStaticFileServer(h *server.Hertz) {
	// 为上传的图片文件提供静态文件服务
	// 修改静态文件服务的路径映射，指向项目根目录
	h.Static("/uploads", "..")

	// 添加中间件处理图片URL
	h.Use(func(c context.Context, ctx *app.RequestContext) {
		path := string(ctx.Path())
		// 检查是否是图片请求
		if strings.HasPrefix(path, "/uploads/") {
			ext := strings.ToLower(filepath.Ext(path))
			// 设置适当的Content-Type
			switch ext {
			case ".jpg", ".jpeg":
				ctx.Header("Content-Type", "image/jpeg")
			case ".png":
				ctx.Header("Content-Type", "image/png")
			case ".gif":
				ctx.Header("Content-Type", "image/gif")
			case ".webp":
				ctx.Header("Content-Type", "image/webp")
			}
			// 禁用缓存以解决304问题
			ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			ctx.Header("Pragma", "no-cache")
			ctx.Header("Expires", "0")
		}
		ctx.Next(c)
	})
}
