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

// GetTags 获取所有标签列表
func GetTags(c context.Context, ctx *app.RequestContext) {
	var tags []models.Tag
	if err := config.DB.Find(&tags).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取标签列表失败",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "获取标签列表成功",
		"data":    tags,
	})
}

// GetTag 获取单个标签详情
func GetTag(c context.Context, ctx *app.RequestContext) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{"code": 400, "message": "无效的标签ID"})
		return
	}

	var tag models.Tag
	if err := config.DB.First(&tag, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, map[string]interface{}{"code": 404, "message": "标签不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"code": 500, "message": "查询标签失败", "error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{"code": 200, "message": "获取标签详情成功", "data": tag})
}

// GetHotTags 获取热门标签 (示例：按关联文章数量排序，取前10)
func GetHotTags(c context.Context, ctx *app.RequestContext) {
	var tags []struct {
		models.Tag
		PostCount int `json:"post_count"`
	}

	// 使用子查询计算每个标签的文章数，并排序
	result := config.DB.Table("tag").
		Select("tag.*, count(post_tags.post_id) as post_count").
		Joins("left join post_tags on post_tags.tag_id = tag.id").
		Where("tag.deleted_at IS NULL"). // 确保只查询未删除的标签
		Group("tag.id").
		Order("post_count DESC").
		Limit(10).
		Scan(&tags)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取热门标签失败",
			"error":   result.Error.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "获取热门标签成功",
		"data":    tags,
	})
}

// --- 以下是管理员权限的操作 ---

// CreateTag 创建标签
func CreateTag(c context.Context, ctx *app.RequestContext) {
	var req struct {
		Name string `json:"name" vd:"len($)>0 && len($)<=50"`
		Slug string `json:"slug" vd:"len($)>0 && len($)<=50"`
	}

	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	tag := models.Tag{
		Name: req.Name,
		Slug: req.Slug,
	}

	if err := config.DB.Create(&tag).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"code": 500, "message": "创建标签失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, map[string]interface{}{"code": 201, "message": "创建标签成功", "data": tag})
}

// UpdateTag 更新标签
func UpdateTag(c context.Context, ctx *app.RequestContext) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{"code": 400, "message": "无效的标签ID"})
		return
	}

	var req struct {
		Name string `json:"name,omitempty" vd:"len($)<=50"`
		Slug string `json:"slug,omitempty" vd:"len($)<=50"`
	}

	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	var tag models.Tag
	if err := config.DB.First(&tag, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, map[string]interface{}{"code": 404, "message": "标签不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"code": 500, "message": "查询标签失败", "error": err.Error()})
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

	if len(updates) > 0 {
		if err := config.DB.Model(&tag).Updates(updates).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"code": 500, "message": "更新标签失败", "error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{"code": 200, "message": "更新标签成功", "data": tag})
}

// DeleteTag 删除标签
func DeleteTag(c context.Context, ctx *app.RequestContext) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{"code": 400, "message": "无效的标签ID"})
		return
	}

	var tag models.Tag
	if err := config.DB.First(&tag, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, map[string]interface{}{"code": 404, "message": "标签不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"code": 500, "message": "查询标签失败", "error": err.Error()})
		}
		return
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		// 解除文章与该标签的关联
		if err := tx.Model(&tag).Association("Posts").Clear(); err != nil {
			return err
		}
		// 删除标签本身 (软删除)
		if err := tx.Delete(&models.Tag{}, id).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"code": 500, "message": "删除标签失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{"code": 200, "message": "删除标签成功"})
}
