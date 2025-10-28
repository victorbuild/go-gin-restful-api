package routers

import (
	"restfulapi/controllers"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes - 設定身份驗證 API
func SetupAuthRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	auth.POST("/register", controllers.RegisterUser)
	auth.POST("/login", controllers.LoginUser)
	auth.POST("/logout", controllers.LogoutUser)
}
