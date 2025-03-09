package routers

import (
	"restfulapi/controllers"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes - 設定身份驗證 API
func SetupAuthRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	auth.POST("/register", controllers.RegisterUser)
	auth.POST("/login", controllers.LoginUser) // ✅ 讓登入 API 獨立
	//auth.POST("/logout", controllers.LogoutUser)     // ✅ 讓登出 API 獨立
}
