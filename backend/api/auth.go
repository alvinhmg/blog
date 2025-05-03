package api

import (
	"github.com/alvinhmg/blog/config"
	"github.com/alvinhmg/blog/middleware"
	"github.com/alvinhmg/blog/models"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname"`
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register 用户注册
func Register(ctx *app.RequestContext) {
	var req RegisterRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	result := config.DB.Where("username = ?", req.Username).First(&existingUser)
	if result.RowsAffected > 0 {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "用户名已存在",
		})
		return
	}

	// 检查邮箱是否已存在
	result = config.DB.Where("email = ?", req.Email).First(&existingUser)
	if result.RowsAffected > 0 {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "邮箱已被注册",
		})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "服务器内部错误",
			"error":   err.Error(),
		})
		return
	}

	// 创建用户
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Role:     "user", // 默认为普通用户
	}

	// 保存用户到数据库
	result = config.DB.Create(&user)
	if result.Error != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "创建用户失败",
			"error":   result.Error.Error(),
		})
		return
	}

	// 生成JWT令牌
	token, err := middleware.GenerateToken(user)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "生成令牌失败",
			"error":   err.Error(),
		})
		return
	}

	// 返回用户信息和令牌
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "注册成功",
		"data": map[string]interface{}{
			"user": map[string]interface{}{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
				"nickname": user.Nickname,
				"role":     user.Role,
			},
			"token": token,
		},
	})
}

// Login 用户登录
func Login(ctx *app.RequestContext) {
	var req LoginRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 查询用户
	var user models.User
	result := config.DB.Where("username = ?", req.Username).First(&user)
	if result.Error != nil {
		ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "用户名或密码错误",
		})
		return
	}

	// 验证密码
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "用户名或密码错误",
		})
		return
	}

	// 生成JWT令牌
	token, err := middleware.GenerateToken(user)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "生成令牌失败",
			"error":   err.Error(),
		})
		return
	}

	// 返回用户信息和令牌
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "登录成功",
		"data": map[string]interface{}{
			"user": map[string]interface{}{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
				"nickname": user.Nickname,
				"role":     user.Role,
			},
			"token": token,
		},
	})
}

// Logout 用户登出
func Logout(ctx *app.RequestContext) {
	// 由于使用JWT，服务端不需要做特殊处理，客户端只需要删除本地存储的令牌即可
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "登出成功",
	})
}
