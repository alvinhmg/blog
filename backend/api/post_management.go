package api

import (
	"context"
	"strconv"
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
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}
	if req.Excerpt != "" {
		post.Excerpt = req.Excerpt
	}
	if req.CoverImage != "" {
		post.CoverImage = req.CoverImage
	}
	if req.Status != "" {
		post.Status = req.Status
	}

	// 更新slug
	if req.Slug != "" {
		post.Slug = req.Slug
	} else {
		post.Slug = slug.Make(post.Title)
	}

	// 开始事务
	tx := config.DB.Begin()

	// 保存更新
	if err := tx.Save(&post).Error; err != nil {
		tx.Rollback()
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "更新文章失败",
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

	// 删除文章
	if err := config.DB.Delete(&post).Error; err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "删除文章失败",
			"error":   err.Error(),
		})
		return
	}

	// 返回成功信息
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "删除文章成功",
	})
}

// LikePost 点赞文章
func LikePost(ctx *app.RequestContext) {
	// 获取文章ID
	id := ctx.Param("id")

	// 获取当前用户ID
	_, exists := ctx.Get("userID")
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

	// 增加点赞数
	post.LikeCount++
	if err := config.DB.Save(&post).Error; err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "点赞失败",
			"error":   err.Error(),
		})
		return
	}

	// 返回成功信息
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "点赞成功",
		"data":    map[string]int{"like_count": post.LikeCount},
	})
}

// CreateCategoryRequest 创建分类请求
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

// UpdateCategoryRequest 更新分类请求
type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

// CreateTagRequest 创建标签请求
type CreateTagRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug"`
}

// UpdateTagRequest 更新标签请求
type UpdateTagRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// SearchPosts 搜索文章
func SearchPosts(c context.Context, ctx *app.RequestContext) {
	// 获取搜索关键词
	keyword := ctx.Query("keyword")
	if keyword == "" {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "搜索关键词不能为空",
		})
		return
	}

	// 获取分页参数
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

	// 构建查询条件
	query := config.DB.Model(&models.Post{}).Where(
		"title LIKE ? OR content LIKE ? OR excerpt LIKE ?",
		"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%",
	)

	// 计算总数
	query.Count(&total)

	// 查询文章列表
	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}
	result := query.Preload("Author").Preload("Categories").Preload("Tags").Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&posts)
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

	// 返回文章列表
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "搜索文章成功",
		"data": map[string]interface{}{
			"posts":      posts,
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"total_page": totalPage,
			"keyword":    keyword,
		},
	})
}

// CreateCategory 创建分类
func CreateCategory(c context.Context, ctx *app.RequestContext) {
	// 解析请求体
	var req CreateCategoryRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 生成slug
	categorySlug := req.Slug
	if categorySlug == "" {
		categorySlug = slug.Make(req.Name)
	}

	// 检查slug是否已存在
	var existingCategory models.Category
	result := config.DB.Where("slug = ?", categorySlug).First(&existingCategory)
	if result.RowsAffected > 0 {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "分类别名已存在",
		})
		return
	}

	// 创建分类
	category := models.Category{
		Name:        req.Name,
		Slug:        categorySlug,
		Description: req.Description,
	}

	if err := config.DB.Create(&category).Error; err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "创建分类失败",
			"error":   err.Error(),
		})
		return
	}

	// 返回分类信息
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "创建分类成功",
		"data":    category,
	})
}

// UpdateCategory 更新分类
func UpdateCategory(c context.Context, ctx *app.RequestContext) {
	// 获取分类ID
	id := ctx.Param("id")

	// 解析请求体
	var req UpdateCategoryRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 查询分类
	var category models.Category
	result := config.DB.First(&category, id)
	if result.Error != nil {
		ctx.JSON(consts.StatusNotFound, map[string]interface{}{
			"code":    404,
			"message": "分类不存在",
			"error":   result.Error.Error(),
		})
		return
	}

	// 更新字段
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}

	// 更新slug
	if req.Slug != "" {
		// 检查新slug是否已存在 (排除当前分类)
		var existingCategory models.Category
		result := config.DB.Where("slug = ? AND id != ?", req.Slug, category.ID).First(&existingCategory)
		if result.RowsAffected > 0 {
			ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
				"code":    400,
				"message": "分类别名已存在",
			})
			return
		}
		category.Slug = req.Slug
	} else if req.Name != "" {
		// 如果没有提供slug但更新了名称，则重新生成slug
		newSlug := slug.Make(category.Name)
		if newSlug != category.Slug {
			// 检查新slug是否已存在 (排除当前分类)
			var existingCategory models.Category
			result := config.DB.Where("slug = ? AND id != ?", newSlug, category.ID).First(&existingCategory)
			if result.RowsAffected > 0 {
				ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
					"code":    400,
					"message": "根据新名称生成的分类别名已存在",
				})
				return
			}
			category.Slug = newSlug
		}
	}

	// 保存更新
	if err := config.DB.Save(&category).Error; err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "更新分类失败",
			"error":   err.Error(),
		})
		return
	}

	// 返回分类信息
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "更新分类成功",
		"data":    category,
	})
}

// DeleteCategory 删除分类
func DeleteCategory(c context.Context, ctx *app.RequestContext) {
	// 获取分类ID
	id := ctx.Param("id")

	// 查询分类
	var category models.Category
	result := config.DB.First(&category, id)
	if result.Error != nil {
		ctx.JSON(consts.StatusNotFound, map[string]interface{}{
			"code":    404,
			"message": "分类不存在",
			"error":   result.Error.Error(),
		})
		return
	}

	// 删除分类
	if err := config.DB.Delete(&category).Error; err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "删除分类失败",
			"error":   err.Error(),
		})
		return
	}

	// 返回成功信息
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "删除分类成功",
	})
}

// CreateTag 创建标签
func CreateTag(c context.Context, ctx *app.RequestContext) {
	// 解析请求体
	var req CreateTagRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 生成slug
	tagSlug := req.Slug
	if tagSlug == "" {
		tagSlug = slug.Make(req.Name)
	}

	// 检查slug是否已存在
	var existingTag models.Tag
	result := config.DB.Where("slug = ?", tagSlug).First(&existingTag)
	if result.RowsAffected > 0 {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "标签别名已存在",
		})
		return
	}

	// 创建标签
	tag := models.Tag{
		Name: req.Name,
		Slug: tagSlug,
	}

	if err := config.DB.Create(&tag).Error; err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "创建标签失败",
			"error":   err.Error(),
		})
		return
	}

	// 返回标签信息
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "创建标签成功",
		"data":    tag,
	})
}

// UpdateTag 更新标签
func UpdateTag(c context.Context, ctx *app.RequestContext) {
	// 获取标签ID
	id := ctx.Param("id")

	// 解析请求体
	var req UpdateTagRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 查询标签
	var tag models.Tag
	result := config.DB.First(&tag, id)
	if result.Error != nil {
		ctx.JSON(consts.StatusNotFound, map[string]interface{}{
			"code":    404,
			"message": "标签不存在",
			"error":   result.Error.Error(),
		})
		return
	}

	// 更新字段
	if req.Name != "" {
		tag.Name = req.Name
	}

	// 更新slug
	if req.Slug != "" {
		// 检查新slug是否已存在 (排除当前标签)
		var existingTag models.Tag
		result := config.DB.Where("slug = ? AND id != ?", req.Slug, tag.ID).First(&existingTag)
		if result.RowsAffected > 0 {
			ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
				"code":    400,
				"message": "标签别名已存在",
			})
			return
		}
		tag.Slug = req.Slug
	} else if req.Name != "" {
		// 如果没有提供slug但更新了名称，则重新生成slug
		newSlug := slug.Make(tag.Name)
		if newSlug != tag.Slug {
			// 检查新slug是否已存在 (排除当前标签)
			var existingTag models.Tag
			result := config.DB.Where("slug = ? AND id != ?", newSlug, tag.ID).First(&existingTag)
			if result.RowsAffected > 0 {
				ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
					"code":    400,
					"message": "根据新名称生成的标签别名已存在",
				})
				return
			}
			tag.Slug = newSlug
		}
	}

	// 保存更新
	if err := config.DB.Save(&tag).Error; err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "更新标签失败",
			"error":   err.Error(),
		})
		return
	}

	// 返回标签信息
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "更新标签成功",
		"data":    tag,
	})
}

// DeleteTag 删除标签
func DeleteTag(c context.Context, ctx *app.RequestContext) {
	// 获取标签ID
	id := ctx.Param("id")

	// 查询标签
	var tag models.Tag
	result := config.DB.First(&tag, id)
	if result.Error != nil {
		ctx.JSON(consts.StatusNotFound, map[string]interface{}{
			"code":    404,
			"message": "标签不存在",
			"error":   result.Error.Error(),
		})
		return
	}

	// 删除标签
	if err := config.DB.Delete(&tag).Error; err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "删除标签失败",
			"error":   err.Error(),
		})
		return
	}

	// 返回成功信息
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "删除标签成功",
	})
}

// GetCategories 获取分类列表
func GetCategories(c context.Context, ctx *app.RequestContext) {
	var categories []models.Category
	if err := config.DB.Find(&categories).Error; err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取分类列表失败",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "获取分类列表成功",
		"data":    categories,
	})
}

// GetCategory 获取分类详情
func GetCategory(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	var category models.Category
	if err := config.DB.First(&category, id).Error; err != nil {
		ctx.JSON(consts.StatusNotFound, map[string]interface{}{
			"code":    404,
			"message": "分类不存在",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "获取分类详情成功",
		"data":    category,
	})
}

// GetTags 获取标签列表
func GetTags(c context.Context, ctx *app.RequestContext) {
	var tags []models.Tag
	if err := config.DB.Find(&tags).Error; err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取标签列表失败",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "获取标签列表成功",
		"data":    tags,
	})
}

// GetTag 获取标签详情
func GetTag(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	var tag models.Tag
	if err := config.DB.First(&tag, id).Error; err != nil {
		ctx.JSON(consts.StatusNotFound, map[string]interface{}{
			"code":    404,
			"message": "标签不存在",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "获取标签详情成功",
		"data":    tag,
	})
}

// GetUsers 获取用户列表 (示例，需要权限控制)
func GetUsers(c context.Context, ctx *app.RequestContext) {
	var users []models.User
	if err := config.DB.Find(&users).Error; err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取用户列表失败",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "获取用户列表成功",
		"data":    users,
	})
}

// GetUser 获取用户详情 (示例，需要权限控制)
func GetUser(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		ctx.JSON(consts.StatusNotFound, map[string]interface{}{
			"code":    404,
			"message": "用户不存在",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "获取用户详情成功",
		"data":    user,
	})
}

// UpdateUser 更新用户信息 (示例，需要权限控制)
func UpdateUser(c context.Context, ctx *app.RequestContext) {
	// 实现更新用户逻辑...
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "更新用户信息成功 (待实现)",
	})
}

// DeleteUser 删除用户 (示例，需要权限控制)
func DeleteUser(c context.Context, ctx *app.RequestContext) {
	// 实现删除用户逻辑...
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "删除用户成功 (待实现)",
	})
}
