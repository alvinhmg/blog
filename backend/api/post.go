package api

import (
	"context"
	"strconv"

	"github.com/alvinhmg/blog/config"
	"github.com/alvinhmg/blog/models"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// GetPosts 获取文章列表 (支持分页、分类、标签过滤)
func GetPosts(c context.Context, ctx *app.RequestContext) {
	// 获取查询参数
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	// 查询文章列表
	var posts []models.Post
	var total int64

	// 计算总数
	config.DB.Model(&models.Post{}).Count(&total)

	// 查询文章列表，包含作者信息
	// 注意：Offset需要非负数，Limit需要正数
	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}
	result := config.DB.Preload("Author").Preload("Categories").Preload("Tags").Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&posts)
	if result.Error != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取文章列表失败",
			"error":   result.Error.Error(),
		})
		return
	}

	// 计算总页数
	totalPage := int64(0)
	if pageSize > 0 {
		totalPage = (total + int64(pageSize) - 1) / int64(pageSize)
	}

	// 返回文章列表
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "获取文章列表成功",
		"data": map[string]interface{}{
			"posts":      posts,
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"total_page": totalPage,
		},
	})
}

// GetPost 获取单篇文章详情
func GetPost(c context.Context, ctx *app.RequestContext) {
	// 获取文章ID
	id := ctx.Param("id")

	// 查询文章详情
	var post models.Post
	result := config.DB.Preload("Author").Preload("Categories").Preload("Tags").Preload("Comments").Preload("Comments.User").First(&post, id)
	if result.Error != nil {
		ctx.JSON(consts.StatusNotFound, map[string]interface{}{
			"code":    404,
			"message": "文章不存在",
			"error":   result.Error.Error(),
		})
		return
	}

	// 增加浏览量
	config.DB.Model(&post).Update("view_count", post.ViewCount+1)

	// 返回文章详情
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "获取文章详情成功",
		"data":    post,
	})
}
