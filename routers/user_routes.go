package routers

import (
	"restfulapi/controllers"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.RouterGroup) {
	user := r.Group("/users")
	user.GET("/", controllers.FindAllUsers)
	user.GET("/:id", controllers.FindByUserId)
	user.POST("/", controllers.RegisterUser)
	user.DELETE("/:id", controllers.DeleteUser)
	user.PUT("/:id", controllers.PutUser)
}
