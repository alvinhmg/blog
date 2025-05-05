package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/alvinhmg/blog/config"
	"github.com/alvinhmg/blog/models"
	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"
)

// GetCategories 获取所有分类列表
func GetCategories(c context.Context, ctx *app.RequestContext) {
	var categories []models.Category
	if err := config.DB.Find(&categories).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取分类列表失败",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "获取分类列表成功",
		"data":    categories,
	})
}

// GetCategory 获取单个分类详情
func GetCategory(c context.Context, ctx *app.RequestContext) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{"code": 400, "message": "无效的分类ID"})
		return
	}

	var category models.Category
	if err := config.DB.First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, map[string]interface{}{"code": 404, "message": "分类不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"code": 500, "message": "查询分类失败", "error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{"code": 200, "message": "获取分类详情成功", "data": category})
}

// GetHotCategories 获取热门分类 (示例：按关联文章数量排序，取前5)
func GetHotCategories(c context.Context, ctx *app.RequestContext) {
	var categories []struct {
		models.Category
		PostCount int `json:"post_count"`
	}

	// 使用子查询计算每个分类的文章数，并排序
	// 注意：此查询可能因数据库类型而异，这里是MySQL/PostgreSQL的示例
	result := config.DB.Table("category").
		Select("category.*, count(post_categories.post_id) as post_count").
		Joins("left join post_categories on post_categories.category_id = category.id").
		Where("category.deleted_at IS NULL"). // 确保只查询未删除的分类
		Group("category.id").
		Order("post_count DESC").
		Limit(5).
		Scan(&categories)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取热门分类失败",
			"error":   result.Error.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "获取热门分类成功",
		"data":    categories,
	})
}

// --- 以下是管理员权限的操作 ---

// CreateCategory 创建分类
func CreateCategory(c context.Context, ctx *app.RequestContext) {
	var req struct {
		Name        string `json:"name" vd:"len($)>0 && len($)<=50"`
		Slug        string `json:"slug" vd:"len($)>0 && len($)<=50"`
		Description string `json:"description,omitempty" vd:"len($)<=255"`
	}

	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	category := models.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
	}

	if err := config.DB.Create(&category).Error; err != nil {
		// 考虑处理唯一约束冲突等错误
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"code": 500, "message": "创建分类失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, map[string]interface{}{"code": 201, "message": "创建分类成功", "data": category})
}

// UpdateCategory 更新分类
func UpdateCategory(c context.Context, ctx *app.RequestContext) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{"code": 400, "message": "无效的分类ID"})
		return
	}

	var req struct {
		Name        string `json:"name,omitempty" vd:"len($)<=50"`
		Slug        string `json:"slug,omitempty" vd:"len($)<=50"`
		Description string `json:"description,omitempty" vd:"len($)<=255"`
	}

	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	var category models.Category
	if err := config.DB.First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, map[string]interface{}{"code": 404, "message": "分类不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"code": 500, "message": "查询分类失败", "error": err.Error()})
		}
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Slug != "" {
		updates["slug"] = req.Slug
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}

	if len(updates) > 0 {
		if err := config.DB.Model(&category).Updates(updates).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"code": 500, "message": "更新分类失败", "error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{"code": 200, "message": "更新分类成功", "data": category})
}

// DeleteCategory 删除分类
func DeleteCategory(c context.Context, ctx *app.RequestContext) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{"code": 400, "message": "无效的分类ID"})
		return
	}

	// 检查分类是否存在
	var category models.Category
	if err := config.DB.First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, map[string]interface{}{"code": 404, "message": "分类不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"code": 500, "message": "查询分类失败", "error": err.Error()})
		}
		return
	}

	// 使用事务删除分类及其关联 (如果需要解除关联)
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		// 解除文章与该分类的关联
		if err := tx.Model(&category).Association("Posts").Clear(); err != nil {
			return err
		}
		// 删除分类本身 (软删除)
		if err := tx.Delete(&models.Category{}, id).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"code": 500, "message": "删除分类失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{"code": 200, "message": "删除分类成功"})
}
