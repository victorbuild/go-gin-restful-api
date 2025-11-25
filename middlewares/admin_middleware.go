package middlewares

import (
	"restfulapi/utils"

	"github.com/gin-gonic/gin"
)

// RequireAdmin - 確保 `role == "admin"` 才能存取
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role.(string) != "admin" {
			c.Error(utils.NewForbiddenError("Forbidden", utils.CodeForbidden, nil))
			c.Abort()
			return
		}
		c.Next()
	}
}
