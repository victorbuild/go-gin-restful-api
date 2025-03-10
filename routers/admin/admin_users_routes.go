package routes

import (
	"restfulapi/controllers"
	"restfulapi/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupAdminUserRoutes - 設定 Admin 管理使用者 API
func SetupAdminUserRoutes(r *gin.RouterGroup) {
	admin := r.Group("/admin/users")
	admin.Use(middlewares.RequireAuth(), middlewares.RequireAdmin())
	{
		admin.POST("/", controllers.RegisterUser)
		admin.GET("/", controllers.FindAllUsers)
		admin.GET("/:id", controllers.FindByUserId)
		admin.DELETE("/:id", controllers.DeleteUser)
		admin.PUT("/:id", controllers.UpdateUser)
	}
}
