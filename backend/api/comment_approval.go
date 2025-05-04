package api

import (
	"context"
	"strconv"

	"github.com/alvinhmg/blog/config"
	"github.com/alvinhmg/blog/models"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// GetPendingComments 获取待审核评论列表
func GetPendingComments(c context.Context, ctx *app.RequestContext) {
	var comments []models.Comment
	if err := config.DB.Where("status = ?", "pending").Preload("User").Preload("Post").Find(&comments).Error; err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "查询待审核评论失败",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "获取待审核评论成功",
		"data":    comments,
	})
}

// ApproveComment 批准评论
func ApproveComment(c context.Context, ctx *app.RequestContext) {
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

	// 查询评论
	var comment models.Comment
	if err := config.DB.First(&comment, commentID).Error; err != nil {
		ctx.JSON(consts.StatusNotFound, map[string]interface{}{
			"code":    404,
			"message": "评论未找到",
		})
		return
	}

	// 更新评论状态
	if err := config.DB.Model(&comment).Update("status", "approved").Error; err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "批准评论失败",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "批准评论成功",
		"data":    comment,
	})
}
