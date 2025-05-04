package api

import (
	"context"
	"strconv"

	"github.com/alvinhmg/blog/config"
	"github.com/alvinhmg/blog/models"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// GetPostComments 获取文章评论列表
func GetPostComments(c context.Context, ctx *app.RequestContext) {
	// 获取文章ID
	postID := ctx.Param("postId")

	// 查询评论列表（只返回已审核通过的评论）
	var comments []models.Comment
	result := config.DB.Where("post_id = ? AND (status = 'approved' OR status IS NULL)", postID).Preload("User").Order("created_at DESC").Find(&comments)
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
func AddComment(c context.Context, ctx *app.RequestContext) {
	// 获取文章ID
	postIDStr := ctx.Param("postId")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "无效的文章ID",
			"error":   err.Error(),
		})
		return
	}

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
		Content  string `json:"content"`
		ParentID *uint  `json:"parent_id,omitempty"` // 可选的父评论ID
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

	// 检查用户角色
	var user models.User
	config.DB.First(&user, userID)

	// 创建评论
	comment := models.Comment{
		Content:  req.Content,
		PostID:   uint(postID),
		UserID:   uint(userID.(uint)),
		ParentID: req.ParentID,
	}

	// 管理员发表的评论自动通过审核
	if user.Role == "admin" {
		comment.Status = "approved"
	} else {
		comment.Status = "pending" // 普通用户评论需要审核
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
	message := "评论已提交，等待审核"
	if user.Role == "admin" {
		message = "评论已发布"
	}
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": message,
		"data":    comment,
	})
}

// DeleteComment 删除评论 (需要权限检查: 评论作者或管理员)
func DeleteComment(c context.Context, ctx *app.RequestContext) {
	// 获取评论ID
	commentIDStr := ctx.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "无效的评论ID",
			"error":   err.Error(),
		})
		return
	}

	// 获取当前用户ID和角色
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
