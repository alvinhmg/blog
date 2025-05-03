package api

import (
	"github.com/alvinhmg/blog/config"
	"github.com/alvinhmg/blog/models"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// GetPostComments 获取文章评论列表
func GetPostComments(ctx *app.RequestContext) {
	// 获取文章ID
	postID := ctx.Param("postId")

	// 查询评论列表
	var comments []models.Comment
	result := config.DB.Where("post_id = ?", postID).Preload("User").Order("created_at DESC").Find(&comments)
	if result.Error != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取评论列表失败",
			"error":   result.Error.Error(),
		})
		return
	}

	// 返回评论列表
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "获取评论列表成功",
		"data":    comments,
	})
}

// AddComment 添加评论
func AddComment(ctx *app.RequestContext) {
	// 获取文章ID
	postID := ctx.Param("postId")

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
	type CommentRequest struct {
		Content string `json:"content"`
	}
	var req CommentRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 验证评论内容
	if req.Content == "" {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "评论内容不能为空",
		})
		return
	}

	// 创建评论
	comment := models.Comment{
		Content: req.Content,
		PostID:  uint(ctx.ParamInt64("postId")),
		UserID:  uint(userID.(uint)),
	}

	// 保存评论
	result := config.DB.Create(&comment)
	if result.Error != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "添加评论失败",
			"error":   result.Error.Error(),
		})
		return
	}

	// 加载用户信息
	config.DB.Preload("User").First(&comment, comment.ID)

	// 返回评论信息
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "添加评论成功",
		"data":    comment,
	})
}

// DeleteComment 删除评论
func DeleteComment(ctx *app.RequestContext) {
	// 获取评论ID
	commentID := ctx.Param("id")

	// 获取当前用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "请先登录",
		})
		return
	}

	// 查询评论
	var comment models.Comment
	result := config.DB.First(&comment, commentID)
	if result.Error != nil {
		ctx.JSON(consts.StatusNotFound, map[string]interface{}{
			"code":    404,
			"message": "评论不存在",
			"error":   result.Error.Error(),
		})
		return
	}

	// 检查权限（只有评论作者或管理员可以删除评论）
	if comment.UserID != uint(userID.(uint)) {
		// 检查是否为管理员
		var user models.User
		config.DB.First(&user, userID)
		if user.Role != "admin" {
			ctx.JSON(consts.StatusForbidden, map[string]interface{}{
				"code":    403,
				"message": "无权删除该评论",
			})
			return
		}
	}

	// 删除评论
	result = config.DB.Delete(&comment)
	if result.Error != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "删除评论失败",
			"error":   result.Error.Error(),
		})
		return
	}

	// 返回成功信息
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "删除评论成功",
	})
}
