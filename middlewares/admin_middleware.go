package middlewares

import (
	"net/http"
	"restfulapi/utils"

	"github.com/gin-gonic/gin"
)

// RequireAdmin - 確保 `role == "admin"` 才能存取
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role.(string) != "admin" {
			utils.ErrorResponse(c, http.StatusForbidden, "Forbidden", 1010)
			c.Abort()
			return
		}
		c.Next()
	}
}
