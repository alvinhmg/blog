package api

import (
	"context"
	"net/http"

	"github.com/alvinhmg/blog/config"
	"github.com/alvinhmg/blog/models"
	"github.com/cloudwego/hertz/pkg/app"
)

// GetHomePageData 获取首页所需数据 (示例：最新文章、热门文章、热门标签/分类)
func GetHomePageData(c context.Context, ctx *app.RequestContext) {
	var latestPosts []models.Post
	var hotPosts []models.Post
	var hotCategories []struct {
		models.Category
		PostCount int `json:"post_count"`
	}
	var hotTags []struct {
		models.Tag
		PostCount int `json:"post_count"`
	}

	// 获取最新文章 (示例：取5篇)
	if err := config.DB.Preload("Author").Preload("Categories").Preload("Tags").Order("created_at DESC").Limit(5).Find(&latestPosts).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"code": 500, "message": "获取最新文章失败", "error": err.Error()})
		return
	}

	// 获取热门文章 (示例：按浏览量排序，取5篇)
	if err := config.DB.Preload("Author").Preload("Categories").Preload("Tags").Order("view_count DESC").Limit(5).Find(&hotPosts).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"code": 500, "message": "获取热门文章失败", "error": err.Error()})
		return
	}

	// 获取热门分类 (示例：按关联文章数量排序，取5)
	if err := config.DB.Table("category").
		Select("category.*, count(post_categories.post_id) as post_count").
		Joins("left join post_categories on post_categories.category_id = category.id").
		Where("category.deleted_at IS NULL").
		Group("category.id").
		Order("post_count DESC").
		Limit(5).
		Scan(&hotCategories).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"code": 500, "message": "获取热门分类失败", "error": err.Error()})
		return
	}

	// 获取热门标签 (示例：按关联文章数量排序，取10)
	if err := config.DB.Table("tag").
		Select("tag.*, count(post_tags.post_id) as post_count").
		Joins("left join post_tags on post_tags.tag_id = tag.id").
		Where("tag.deleted_at IS NULL").
		Group("tag.id").
		Order("post_count DESC").
		Limit(10).
		Scan(&hotTags).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"code": 500, "message": "获取热门标签失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "获取首页数据成功",
		"data": map[string]interface{}{
			"latest_posts":   latestPosts,
			"hot_posts":      hotPosts,
			"hot_categories": hotCategories,
			"hot_tags":       hotTags,
		},
	})
}

// GetArchiveData 获取归档数据 (按年月分组)
func GetArchiveData(c context.Context, ctx *app.RequestContext) {
	type ArchiveItem struct {
		YearMonth string        `json:"year_month"` // 格式如 "2023-10"
		Posts     []models.Post `json:"posts"`
	}

	var results []struct {
		YearMonth string
		models.Post
	}

	// 查询所有已发布文章，并按年月分组
	// 注意：数据库特定的日期格式化函数
	// MySQL: DATE_FORMAT(created_at, '%Y-%m')
	// PostgreSQL: TO_CHAR(created_at, 'YYYY-MM')
	// SQLite: strftime('%Y-%m', created_at)
	// 假设使用 MySQL
	dbResult := config.DB.Model(&models.Post{}).
		Select("DATE_FORMAT(created_at, '%Y-%m') as year_month, id, title, slug, created_at"). // 选择需要的字段
		Where("status = ?", "published").
		Order("created_at DESC").
		Scan(&results)

	if dbResult.Error != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取归档数据失败",
			"error":   dbResult.Error.Error(),
		})
		return
	}

	// 按年月组织数据
	archiveMap := make(map[string][]models.Post)
	var orderedKeys []string // 保持年月顺序

	for _, r := range results {
		if _, exists := archiveMap[r.YearMonth]; !exists {
			archiveMap[r.YearMonth] = []models.Post{}
			orderedKeys = append(orderedKeys, r.YearMonth) // 记录年月顺序
		}
		// 创建一个只包含所需字段的 Post 结构
		postSummary := models.Post{
			ID:        r.Post.ID,
			Title:     r.Post.Title,
			Slug:      r.Post.Slug,
			CreatedAt: r.Post.CreatedAt,
		}
		archiveMap[r.YearMonth] = append(archiveMap[r.YearMonth], postSummary)
	}

	// 转换为最终的列表格式，保持顺序
	var archiveData []ArchiveItem
	for _, key := range orderedKeys {
		archiveData = append(archiveData, ArchiveItem{
			YearMonth: key,
			Posts:     archiveMap[key],
		})
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "获取归档数据成功",
		"data":    archiveData,
	})
}
