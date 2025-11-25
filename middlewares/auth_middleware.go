package middlewares

import (
	"context"
	"restfulapi/config"
	"restfulapi/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// RequireAuth - 讓 API 驗證 `JWT Token`
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.Error(utils.NewUnauthorizedError("Unauthorized", utils.CodeTokenMissing, nil))
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// 檢查 Access Token 是否在黑名單中（登出時加入的）
		ctx := context.Background()
		blacklistKey := "access_token:blacklist:" + tokenString
		exists, err := config.RedisClient.Exists(ctx, blacklistKey).Result()
		if err == nil && exists > 0 {
			c.Error(utils.NewUnauthorizedError("Token has been revoked", utils.CodeTokenInvalid, nil))
			c.Abort()
			return
		}

		// 驗證 Token
		token, err := config.ValidateToken(tokenString, config.JWTConfig.AccessTokenSecret)
		if err != nil || !token.Valid {
			c.Error(utils.NewUnauthorizedError("Invalid token", utils.CodeTokenInvalid, err))
			c.Abort()
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)
		c.Set("user_id", claims["sub"])
		c.Set("email", claims["email"])
		c.Set("role", claims["role"])
		c.Next()
	}
}
