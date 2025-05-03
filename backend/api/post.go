package api

import (
	"github.com/alvinhmg/blog/config"
	"github.com/alvinhmg/blog/models"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// GetPosts 获取文章列表
func GetPosts(ctx *app.RequestContext) {
	// 获取分页参数
	page := ctx.DefaultQuery("page", "1")
	pageSize := ctx.DefaultQuery("page_size", "10")

	// 查询文章列表
	var posts []models.Post
	var total int64

	// 计算总数
	config.DB.Model(&models.Post{}).Count(&total)

	// 查询文章列表，包含作者信息
	result := config.DB.Preload("Author").Preload("Categories").Preload("Tags").Order("created_at DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&posts)
	if result.Error != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取文章列表失败",
			"error":   result.Error.Error(),
		})
		return
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
			"total_page": (total + pageSize - 1) / pageSize,
		},
	})
}

// GetPost 获取文章详情
func GetPost(ctx *app.RequestContext) {
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
