package routers

import (
	"restfulapi/controllers"
	"restfulapi/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes - 設定身份驗證 API
func SetupAuthRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	auth.POST("/register", controllers.RegisterUser)
	auth.POST("/login", controllers.LoginUser)
	auth.POST("/refresh", controllers.RefreshToken)
	auth.POST("/logout", middlewares.RequireAuth(), controllers.LogoutUser)
}
