package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Pagination 定義分頁資訊
type Pagination struct {
	Total      int `json:"total"`
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalPages int `json:"total_pages"`
}

// MetaData 提供額外資訊
type MetaData struct {
	RequestID string `json:"request_id"`
	Timestamp string `json:"timestamp"`
}

// APIResponse 定義標準 API 回應結構
type APIResponse struct {
	// 狀態: "success" 或 "error"
	Status string `json:"status" example:"success"`

	// 訊息描述
	Message string `json:"message" example:"User created successfully"`

	// 主要的回應數據 (可省略)
	Data interface{} `json:"data,omitempty"`

	// 錯誤代碼 (成功時省略)
	ErrorCode int `json:"error_code,omitempty" example:"1001"`

	// Meta 資訊 (例如 API 版本、請求 ID 等)
	Meta MetaData `json:"meta"`
}

// 產生 MetaData
func generateMetaData() MetaData {
	return MetaData{
		RequestID: uuid.New().String(),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// SuccessResponse 成功回應（單筆資料）
func SuccessResponse(c *gin.Context, message string, data interface{}) {
	response := APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
		Meta:    generateMetaData(),
	}
	c.JSON(http.StatusOK, response)
}

// ListResponse 成功回應（列表 + 分頁）
func ListResponse(c *gin.Context, message string, items interface{}, pagination Pagination) {
	response := APIResponse{
		Status:  "success",
		Message: message,
		Data: gin.H{
			"items":      items,
			"pagination": pagination,
		},
		Meta: generateMetaData(),
	}
	c.JSON(http.StatusOK, response)
}

// ErrorResponse 失敗回應
func ErrorResponse(c *gin.Context, statusCode int, message string, errorCode int) {
	response := APIResponse{
		Status:    "error",
		Message:   message,
		ErrorCode: errorCode,
		Meta:      generateMetaData(),
	}
	c.JSON(statusCode, response)
}
