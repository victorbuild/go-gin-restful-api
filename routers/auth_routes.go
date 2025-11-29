package routers

import (
	"restfulapi/internal/handler"
	"restfulapi/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes - 設定身份驗證 API
func SetupAuthRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	auth.POST("/register", handler.RegisterUser)
	auth.POST("/login", handler.LoginUser)
	auth.POST("/refresh", handler.RefreshToken)
	auth.POST("/logout", middleware.RequireAuth(), handler.LogoutUser)
}
