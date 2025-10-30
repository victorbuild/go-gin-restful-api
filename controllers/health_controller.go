package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
)

// HealthResponse 定義 /health 標準回應格式
// @Description 健康檢查標準格式
// swagger:model HealthResponse
type HealthResponse struct {
	Service   string `json:"service" example:"user-api"`
	Status    string `json:"status" example:"healthy"`
	Timestamp string `json:"timestamp" example:"2025-10-30T17:52:05+08:00"`
	Version   string `json:"version" example:"1.0"`
}

// HealthCheck - 健康檢查端點
// @Summary 健康檢查
// @Description 檢查 API 服務健康狀態
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} controllers.HealthResponse
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(200, HealthResponse{
		Service:   "user-api",
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   "1.0",
	})
}

// HealthRootResponse 定義 / 路由健康檢查標準格式
// @Description 服務根目錄健康回應
// swagger:model HealthRootResponse
type HealthRootResponse struct {
	Status string `json:"status" example:"ok"`
}

// Root - 服務根路徑
// @Summary 服務根路徑
// @Description 顯示 API 狀態（根目錄）
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} controllers.HealthRootResponse "服務啟動"
// @Router / [get]
func Root(c *gin.Context) {
	c.JSON(200, HealthRootResponse{Status: "ok"})
}
