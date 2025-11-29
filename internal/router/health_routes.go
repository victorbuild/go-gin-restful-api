package router

import (
	"restfulapi/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupHealthRoutes(r *gin.Engine) {
	r.GET("/", handler.Root)
	r.GET("/health", handler.HealthCheck)
}
