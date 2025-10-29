package utils

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

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
	ErrorCode int `json:"error_code,omitempty" example:"1000"`

	// Meta 資訊 (例如 API 版本、請求 ID 等)
	Meta MetaData `json:"meta"`
}

// SuccessAPIResponse 定義成功回應結構（包含 data 欄位）
type SuccessAPIResponse struct {
	// Status 狀態: "success"
	Status string `json:"status" example:"success"`

	// Message 訊息描述
	Message string `json:"message" example:"User created successfully"`

	// Data 主要的回應數據
	Data interface{} `json:"data"`

	// Meta Meta 資訊 (例如 API 版本、請求 ID 等)
	Meta MetaData `json:"meta"`
}

// ErrorAPIResponse 定義錯誤回應結構
// @Description 錯誤回應格式，error_code 視情況而定
type ErrorAPIResponse struct {
	// Status 狀態: "error"
	Status string `json:"status" example:"error"`

	// Message 訊息描述
	Message string `json:"message" example:"Email already registered"`

	// ErrorCode 錯誤代碼
	ErrorCode int `json:"error_code" swaggertype:"integer"`

	// Meta Meta 資訊 (例如 API 版本、請求 ID 等)
	Meta MetaData `json:"meta"`
}

// ErrorAPIResponseUnsupportedMediaType 定義 415 不支援媒體類型錯誤回應結構
// @Description 不支援的媒體類型錯誤回應
type ErrorAPIResponseUnsupportedMediaType struct {
	// Status 狀態: "error"
	Status string `json:"status" example:"error"`

	// Message 訊息描述
	Message string `json:"message" example:"Unsupported media type. Expected application/json"`

	// ErrorCode 錯誤代碼: 1000
	ErrorCode int `json:"error_code" example:"1000" swaggertype:"integer"`

	// Meta Meta 資訊 (例如 API 版本、請求 ID 等)
	Meta MetaData `json:"meta"`
}

// ErrorAPIResponseEmailExists 定義 409 Email 重複錯誤回應結構
// @Description Email 已經被註冊錯誤回應
type ErrorAPIResponseEmailExists struct {
	// Status 狀態: "error"
	Status string `json:"status" example:"error"`

	// Message 訊息描述
	Message string `json:"message" example:"Email already registered"`

	// ErrorCode 錯誤代碼: 1003
	ErrorCode int `json:"error_code" example:"1003" swaggertype:"integer"`

	// Meta Meta 資訊 (例如 API 版本、請求 ID 等)
	Meta MetaData `json:"meta"`
}

// ErrorAPIResponseInvalidInput 定義 400 無效輸入錯誤回應結構
// @Description 無效的輸入格式（JSON 格式錯誤）錯誤回應
type ErrorAPIResponseInvalidInput struct {
	// Status 狀態: "error"
	Status string `json:"status" example:"error"`

	// Message 訊息描述
	Message string `json:"message" example:"Invalid input"`

	// ErrorCode 錯誤代碼: 1001
	ErrorCode int `json:"error_code" example:"1001" swaggertype:"integer"`

	// Meta Meta 資訊 (例如 API 版本、請求 ID 等)
	Meta MetaData `json:"meta"`
}

// ErrorAPIResponseMissingFields 定義 400 缺少必填欄位錯誤回應結構
// @Description 缺少必填欄位錯誤回應
type ErrorAPIResponseMissingFields struct {
	// Status 狀態: "error"
	Status string `json:"status" example:"error"`

	// Message 訊息描述
	Message string `json:"message" example:"Missing required fields"`

	// ErrorCode 錯誤代碼: 1002
	ErrorCode int `json:"error_code" example:"1002" swaggertype:"integer"`

	// Meta Meta 資訊 (例如 API 版本、請求 ID 等)
	Meta MetaData `json:"meta"`
}

// ErrorAPIResponseInvalidCredentials 定義 401 帳號密碼錯誤回應結構
// @Description 帳號或密碼錯誤回應
type ErrorAPIResponseInvalidCredentials struct {
	// Status 狀態: "error"
	Status string `json:"status" example:"error"`

	// Message 訊息描述
	Message string `json:"message" example:"Invalid email or password"`

	// ErrorCode 錯誤代碼: 1004
	ErrorCode int `json:"error_code" example:"1004" swaggertype:"integer"`

	// Meta Meta 資訊 (例如 API 版本、請求 ID 等)
	Meta MetaData `json:"meta"`
}

// ErrorAPIResponseInternalServerError 定義 500 伺服器內部錯誤回應結構
// @Description 伺服器內部錯誤回應
type ErrorAPIResponseInternalServerError struct {
	// Status 狀態: "error"
	Status string `json:"status" example:"error"`

	// Message 訊息描述
	Message string `json:"message" example:"Failed to create user"`

	// ErrorCode 錯誤代碼: 4001
	ErrorCode int `json:"error_code" example:"4001" swaggertype:"integer"`

	// Meta Meta 資訊 (例如 API 版本、請求 ID 等)
	Meta MetaData `json:"meta"`
}

// MarshalJSON 自訂 JSON 序列化以確保欄位順序
func (e ErrorAPIResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Status    string   `json:"status"`
		Message   string   `json:"message"`
		ErrorCode int      `json:"error_code"`
		Meta      MetaData `json:"meta"`
	}{
		Status:    e.Status,
		Message:   e.Message,
		ErrorCode: e.ErrorCode,
		Meta:      e.Meta,
	})
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

// CreatedResponse 成功創建回應（201）
func CreatedResponse(c *gin.Context, message string, data interface{}) {
	response := APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
		Meta:    generateMetaData(),
	}
	c.JSON(http.StatusCreated, response)
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
