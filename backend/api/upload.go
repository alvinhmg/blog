package api

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// 图片上传处理函数
func UploadImage(c context.Context, ctx *app.RequestContext) {
	// 获取当前用户ID，确保用户已登录
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "请先登录",
		})
		return
	}

	// 获取上传的文件
	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "请选择要上传的图片",
			"error":   err.Error(),
		})
		return
	}

	// 检查文件类型
	filename := file.Filename
	ext := strings.ToLower(filepath.Ext(filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	if !allowedExts[ext] {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "不支持的图片格式，请上传jpg、jpeg、png、gif或webp格式的图片",
		})
		return
	}

	// 创建上传目录
	// 使用backend目录下的uploads目录，与静态文件服务配置一致
	// 使用相对路径以便部署到服务器
	uploadsDir := filepath.Join("uploads", "images")
	// 确保目录存在
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "创建上传目录失败",
			"error":   err.Error(),
		})
		return
	}

	// 生成唯一文件名
	timestamp := time.Now().Format("20060102150405")
	newFilename := fmt.Sprintf("%d_%s%s", userID, timestamp, ext)
	filePath := filepath.Join(uploadsDir, newFilename)

	// 打印文件路径，用于调试
	fmt.Printf("保存图片到路径: %s\n", filePath)

	// 保存文件
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "保存图片失败",
			"error":   err.Error(),
		})
		return
	}

	// 返回图片URL
	// 静态文件服务配置为/uploads指向项目根目录的uploads文件夹
	// 所以URL路径应该保持一致
	imageURL := "/uploads/images/" + newFilename
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "图片上传成功",
		"data": map[string]interface{}{
			"url": imageURL,
		},
	})
}
