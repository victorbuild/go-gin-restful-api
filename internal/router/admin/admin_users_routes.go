package router

import (
	"restfulapi/internal/handler"
	"restfulapi/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupAdminUserRoutes - 設定 Admin 管理使用者 API
func SetupAdminUserRoutes(r *gin.RouterGroup) {
	admin := r.Group("/admin/users")
	admin.Use(middleware.RequireAuth(), middleware.RequireAdmin())
	{
		admin.POST("/", handler.RegisterUser)
		admin.GET("/", handler.FindAllUsers)
		admin.GET("/:id", handler.FindByUserId)
		admin.DELETE("/:id", handler.DeleteUser)
		admin.PUT("/:id", handler.UpdateUser)
	}
}
