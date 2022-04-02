package routers

import (
	controller "restfulapi/Controller"

	"github.com/gin-gonic/gin"
)

func AddUserRouter(r *gin.RouterGroup) {
	user := r.Group("/users")
	user.GET("/", controller.FindAllUsers)
	user.POST("/", controller.PostUser)
	user.DELETE("/:id", controller.DeleteUser)
	user.PUT("/:id", controller.PutUser)
}
