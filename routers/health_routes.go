package routers

import (
	"restfulapi/controllers"

	"github.com/gin-gonic/gin"
)

func SetupHealthRoutes(r *gin.Engine) {
	r.GET("/", controllers.Root)
	r.GET("/health", controllers.HealthCheck)
}
