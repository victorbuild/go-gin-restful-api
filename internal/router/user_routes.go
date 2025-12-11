package router

import (
	"restfulapi/internal/handler"
	"restfulapi/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes - 設定一般使用者 API
func SetupUserRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")
	users.Use(middleware.RequireAuth())
	{
		users.GET("/me", handler.GetMyProfile)
	}
}
