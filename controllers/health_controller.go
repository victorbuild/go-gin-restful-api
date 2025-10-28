package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
)

// HealthCheck - 健康檢查端點
// @Summary 健康檢查
// @Description 檢查 API 服務健康狀態
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "服務健康"
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "user-api",
		"version":   "1.0",
	})
}

func Root(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
