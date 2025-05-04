package middleware

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alvinhmg/blog/config"
	"github.com/alvinhmg/blog/models"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuth JWT认证中间件
func JWTAuth() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 从请求头获取token
		tokenString := string(ctx.GetHeader("Authorization"))
		if tokenString == "" {
			ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "未提供认证令牌",
				"data":    nil,
			})
			ctx.Abort()
			return
		}

		// 移除Bearer前缀
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// 解析JWT令牌
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 验证签名算法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("无效的签名方法: %v", token.Header["alg"])
			}
			// 返回密钥
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "无效的认证令牌",
				"error":   err.Error(),
				"data":    nil,
			})
			ctx.Abort()
			return
		}

		// 验证令牌有效性
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "无效的认证令牌",
				"data":    nil,
			})
			ctx.Abort()
			return
		}

		// 获取用户ID
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "无效的用户ID",
				"data":    nil,
			})
			ctx.Abort()
			return
		}
		userID := uint(userIDFloat)

		// 查询用户信息
		var user models.User
		result := config.DB.First(&user, userID)
		if result.Error != nil {
			ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "用户不存在",
				"data":    nil,
			})
			ctx.Abort()
			return
		}

		// 将用户信息存储到上下文中
		ctx.Set("user", user)
		ctx.Set("userID", userID)
		ctx.Next(c)
	}
}

// GenerateToken 生成JWT令牌
func GenerateToken(user models.User) (string, error) {
	// 创建JWT声明
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // 7天过期
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// AdminAuth 管理员权限中间件
func AdminAuth() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 获取用户信息
		userInterface, exists := ctx.Get("user")
		if !exists {
			ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "未认证",
				"data":    nil,
			})
			ctx.Abort()
			return
		}

		// 类型断言
		user, ok := userInterface.(models.User)
		if !ok {
			ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
				"code":    500,
				"message": "用户信息类型错误",
				"data":    nil,
			})
			ctx.Abort()
			return
		}

		// 检查用户角色
		if user.Role != "admin" {
			ctx.JSON(consts.StatusForbidden, map[string]interface{}{
				"code":    403,
				"message": "权限不足",
				"data":    nil,
			})
			ctx.Abort()
			return
		}

		ctx.Next(c)
	}
}
