package api

import (
	"context"
	"strconv"

	"github.com/alvinhmg/blog/config"
	"github.com/alvinhmg/blog/models"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"gorm.io/gorm"
)

// GetPosts 获取文章列表 (支持分页、分类、标签过滤)
func GetPosts(c context.Context, ctx *app.RequestContext) {
	// 获取查询参数
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")
	categoryIDStr := ctx.Query("category_id")
	tagIDStr := ctx.Query("tag_id")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	// 构建查询
	db := config.DB.Model(&models.Post{}).Preload("Author").Preload("Categories").Preload("Tags")

	// 分类过滤
	if categoryIDStr != "" {
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err == nil {
			db = db.Joins("JOIN post_categories ON post_categories.post_id = post.id").Where("post_categories.category_id = ?", categoryID)
		}
	}

	// 标签过滤
	if tagIDStr != "" {
		tagID, err := strconv.Atoi(tagIDStr)
		if err == nil {
			db = db.Joins("JOIN post_tags ON post_tags.post_id = post.id").Where("post_tags.tag_id = ?", tagID)
		}
	}

	// 查询文章列表
	var posts []models.Post
	var total int64

	// 计算总数
	db.Count(&total)

	// 查询文章列表，包含作者信息
	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}
	result := db.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&posts)
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

	// 增加浏览量 (使用事务保证原子性)
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&post).Update("view_count", gorm.Expr("view_count + ?", 1)).Error
	})

	if err != nil {
		// 记录日志，但不阻塞响应
		println("增加浏览量失败:", err.Error())
	}

	// 重新加载文章数据以获取最新的浏览量
	config.DB.Preload("Author").Preload("Categories").Preload("Tags").Preload("Comments").Preload("Comments.User").First(&post, id)

	// 返回文章详情
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "获取文章详情成功",
		"data":    post,
	})
}

// SearchPosts 全文搜索文章
func SearchPosts(c context.Context, ctx *app.RequestContext) {
	query := ctx.Query("q")
	if query == "" {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "缺少搜索关键词 'q'",
		})
		return
	}

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

	var posts []models.Post
	var total int64

	db := config.DB.Model(&models.Post{}).Preload("Author").Preload("Categories").Preload("Tags")

	// 在标题和内容中进行不区分大小写的模糊搜索
	searchQuery := "%" + query + "%"
	db = db.Where("LOWER(title) LIKE LOWER(?) OR LOWER(content) LIKE LOWER(?) ", searchQuery, searchQuery)

	// 计算总数
	db.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}
	result := db.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&posts)
	if result.Error != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "搜索文章失败",
			"error":   result.Error.Error(),
		})
		return
	}

	// 计算总页数
	totalPage := int64(0)
	if pageSize > 0 {
		totalPage = (total + int64(pageSize) - 1) / int64(pageSize)
	}

	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "搜索文章成功",
		"data": map[string]interface{}{
			"posts":      posts,
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"total_page": totalPage,
		},
	})
}

// LikePost 点赞文章
func LikePost(c context.Context, ctx *app.RequestContext) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "无效的文章ID",
		})
		return
	}

	var post models.Post
	// 检查文章是否存在
	if err := config.DB.First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(consts.StatusNotFound, map[string]interface{}{
				"code":    404,
				"message": "文章不存在",
			})
		} else {
			ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
				"code":    500,
				"message": "查询文章失败",
				"error":   err.Error(),
			})
		}
		return
	}

	// 增加点赞数 (使用事务保证原子性)
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&post).Update("like_count", gorm.Expr("like_count + ?", 1)).Error
	})

	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "点赞失败",
			"error":   err.Error(),
		})
		return
	}

	// 返回成功信息和更新后的点赞数
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "点赞成功",
		"data": map[string]interface{}{
			"like_count": post.LikeCount + 1, // 返回预期的点赞数
		},
	})
}
