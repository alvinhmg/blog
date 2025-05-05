package api

import (
	"context"
	"time"

	"github.com/alvinhmg/blog/config"
	"github.com/alvinhmg/blog/models"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/gosimple/slug"
)

// CreatePostRequest 创建文章请求
type CreatePostRequest struct {
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required"`
	Excerpt    string `json:"excerpt"`
	CoverImage string `json:"cover_image"`
	Status     string `json:"status" binding:"required,oneof=draft published"`
	Categories []uint `json:"categories"`
	Tags       []uint `json:"tags"`
	Slug       string `json:"slug"`
}

// UpdatePostRequest 更新文章请求
type UpdatePostRequest struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	Excerpt    string `json:"excerpt"`
	CoverImage string `json:"cover_image"`
	Status     string `json:"status" binding:"omitempty,oneof=draft published"`
	Categories []uint `json:"categories"`
	Tags       []uint `json:"tags"`
	Slug       string `json:"slug"`
}

// CreatePost 创建文章
func CreatePost(c context.Context, ctx *app.RequestContext) {
	// 获取当前用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "请先登录",
		})
		return
	}

	// 解析请求体
	var req CreatePostRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 生成文章slug
	postSlug := req.Slug
	if postSlug == "" {
		postSlug = slug.Make(req.Title)
	}

	// 检查slug是否已存在
	var existingPost models.Post
	result := config.DB.Where("slug = ?", postSlug).First(&existingPost)
	if result.RowsAffected > 0 {
		// 如果slug已存在，添加时间戳后缀
		postSlug = postSlug + "-" + time.Now().Format("20060102150405")
	}

	// 创建文章
	post := models.Post{
		Title:      req.Title,
		Slug:       postSlug,
		Content:    req.Content,
		Excerpt:    req.Excerpt,
		CoverImage: req.CoverImage,
		Status:     req.Status,
		AuthorID:   uint(userID.(uint)),
	}

	// 开始事务
	tx := config.DB.Begin()

	// 保存文章
	if err := tx.Create(&post).Error; err != nil {
		tx.Rollback()
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "创建文章失败",
			"error":   err.Error(),
		})
		return
	}

	// 处理分类
	if len(req.Categories) > 0 {
		var categories []models.Category
		if err := tx.Where("id IN ?", req.Categories).Find(&categories).Error; err != nil {
			tx.Rollback()
			ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
				"code":    500,
				"message": "获取分类失败",
				"error":   err.Error(),
			})
			return
		}
		if err := tx.Model(&post).Association("Categories").Replace(categories); err != nil {
			tx.Rollback()
			ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
				"code":    500,
				"message": "设置文章分类失败",
				"error":   err.Error(),
			})
			return
		}
	}

	// 处理标签
	if len(req.Tags) > 0 {
		var tags []models.Tag
		if err := tx.Where("id IN ?", req.Tags).Find(&tags).Error; err != nil {
			tx.Rollback()
			ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
				"code":    500,
				"message": "获取标签失败",
				"error":   err.Error(),
			})
			return
		}
		if err := tx.Model(&post).Association("Tags").Replace(tags); err != nil {
			tx.Rollback()
			ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
				"code":    500,
				"message": "设置文章标签失败",
				"error":   err.Error(),
			})
			return
		}
	}

	// 提交事务
	tx.Commit()

	// 加载关联数据
	config.DB.Preload("Author").Preload("Categories").Preload("Tags").First(&post, post.ID)

	// 返回文章信息
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "创建文章成功",
		"data":    post,
	})
}

// UpdatePost 更新文章
func UpdatePost(c context.Context, ctx *app.RequestContext) {
	// 获取文章ID
	id := ctx.Param("id")

	// 获取当前用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "请先登录",
		})
		return
	}

	// 查询文章
	var post models.Post
	result := config.DB.First(&post, id)
	if result.Error != nil {
		ctx.JSON(consts.StatusNotFound, map[string]interface{}{
			"code":    404,
			"message": "文章不存在",
			"error":   result.Error.Error(),
		})
		return
	}

	// 检查权限（只有作者或管理员可以更新文章）
	if post.AuthorID != uint(userID.(uint)) {
		// 检查是否为管理员
		var user models.User
		config.DB.First(&user, userID)
		if user.Role != "admin" {
			ctx.JSON(consts.StatusForbidden, map[string]interface{}{
				"code":    403,
				"message": "无权更新该文章",
			})
			return
		}
	}

	// 解析请求体
	var req UpdatePostRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 更新文章字段
	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Excerpt != "" {
		updates["excerpt"] = req.Excerpt
	}
	if req.CoverImage != "" {
		updates["cover_image"] = req.CoverImage
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	// 处理slug更新
	postSlug := req.Slug
	if postSlug != "" && postSlug != post.Slug {
		// 检查新slug是否已存在
		var existingPost models.Post
		result := config.DB.Where("slug = ? AND id != ?", postSlug, post.ID).First(&existingPost)
		if result.RowsAffected > 0 {
			// 如果slug已存在，添加时间戳后缀
			postSlug = postSlug + "-" + time.Now().Format("20060102150405")
		}
		updates["slug"] = postSlug
	}

	// 开始事务
	tx := config.DB.Begin()

	// 更新文章基本信息
	if len(updates) > 0 {
		if err := tx.Model(&post).Updates(updates).Error; err != nil {
			tx.Rollback()
			ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
				"code":    500,
				"message": "更新文章失败",
				"error":   err.Error(),
			})
			return
		}
	}

	// 更新分类
	if req.Categories != nil {
		var categories []models.Category
		if len(req.Categories) > 0 {
			if err := tx.Where("id IN ?", req.Categories).Find(&categories).Error; err != nil {
				tx.Rollback()
				ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
					"code":    500,
					"message": "获取分类失败",
					"error":   err.Error(),
				})
				return
			}
		}
		if err := tx.Model(&post).Association("Categories").Replace(categories); err != nil {
			tx.Rollback()
			ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
				"code":    500,
				"message": "更新文章分类失败",
				"error":   err.Error(),
			})
			return
		}
	}

	// 更新标签
	if req.Tags != nil {
		var tags []models.Tag
		if len(req.Tags) > 0 {
			if err := tx.Where("id IN ?", req.Tags).Find(&tags).Error; err != nil {
				tx.Rollback()
				ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
					"code":    500,
					"message": "获取标签失败",
					"error":   err.Error(),
				})
				return
			}
		}
		if err := tx.Model(&post).Association("Tags").Replace(tags); err != nil {
			tx.Rollback()
			ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
				"code":    500,
				"message": "更新文章标签失败",
				"error":   err.Error(),
			})
			return
		}
	}

	// 提交事务
	tx.Commit()

	// 加载更新后的关联数据
	config.DB.Preload("Author").Preload("Categories").Preload("Tags").First(&post, post.ID)

	// 返回更新后的文章信息
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "更新文章成功",
		"data":    post,
	})
}

// DeletePost 删除文章
func DeletePost(c context.Context, ctx *app.RequestContext) {
	// 获取文章ID
	id := ctx.Param("id")

	// 获取当前用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "请先登录",
		})
		return
	}

	// 查询文章
	var post models.Post
	result := config.DB.First(&post, id)
	if result.Error != nil {
		ctx.JSON(consts.StatusNotFound, map[string]interface{}{
			"code":    404,
			"message": "文章不存在",
			"error":   result.Error.Error(),
		})
		return
	}

	// 检查权限（只有作者或管理员可以删除文章）
	if post.AuthorID != uint(userID.(uint)) {
		// 检查是否为管理员
		var user models.User
		config.DB.First(&user, userID)
		if user.Role != "admin" {
			ctx.JSON(consts.StatusForbidden, map[string]interface{}{
				"code":    403,
				"message": "无权删除该文章",
			})
			return
		}
	}

	// 开始事务
	tx := config.DB.Begin()

	// 删除文章与分类的关联
	if err := tx.Model(&post).Association("Categories").Clear(); err != nil {
		tx.Rollback()
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "删除文章分类关联失败",
			"error":   err.Error(),
		})
		return
	}

	// 删除文章与标签的关联
	if err := tx.Model(&post).Association("Tags").Clear(); err != nil {
		tx.Rollback()
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "删除文章标签关联失败",
			"error":   err.Error(),
		})
		return
	}

	// 删除文章的评论
	if err := tx.Where("post_id = ?", post.ID).Delete(&models.Comment{}).Error; err != nil {
		tx.Rollback()
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "删除文章评论失败",
			"error":   err.Error(),
		})
		return
	}

	// 删除文章
	if err := tx.Delete(&post).Error; err != nil {
		tx.Rollback()
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "删除文章失败",
			"error":   err.Error(),
		})
		return
	}

	// 提交事务
	tx.Commit()

	// 返回成功信息
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "删除文章成功",
	})
}
