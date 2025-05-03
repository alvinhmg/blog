package middleware

import (
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
	return func(ctx *app.RequestContext) {
		// 从请求头获取token
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "未提供认证令牌",
			})
			ctx.Abort()
			return
		}

		// 移除Bearer前缀
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

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
			})
			ctx.Abort()
			return
		}

		// 验证令牌有效性
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// 获取用户ID
			userID := uint(claims["user_id"].(float64))

			// 查询用户信息
			var user models.User
			result := config.DB.First(&user, userID)
			if result.Error != nil {
				ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
					"code":    401,
					"message": "用户不存在",
				})
				ctx.Abort()
				return
			}

			// 将用户信息存储到上下文中
			ctx.Set("user", user)
			ctx.Set("user_id", userID)
			ctx.Next()
		} else {
			ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "无效的认证令牌",
			})
			ctx.Abort()
			return
		}
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
	return func(ctx *app.RequestContext) {
		// 获取用户信息
		user, exists := ctx.Get("user")
		if !exists {
			ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "未认证",
			})
			ctx.Abort()
			return
		}

		// 检查用户角色
		if user.(models.User).Role != "admin" {
			ctx.JSON(consts.StatusForbidden, map[string]interface{}{
				"code":    403,
				"message": "权限不足",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
