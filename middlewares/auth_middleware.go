package middlewares

import (
	"net/http"
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
			utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", 1006)
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		token, err := config.ValidateToken(tokenString, config.JWTConfig.AccessTokenSecret)
		if err != nil || !token.Valid {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid token", 1007)
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
