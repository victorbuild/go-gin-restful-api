package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// MetaData 提供額外資訊
type MetaData struct {
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

	// Meta 資訊 (例如 API 版本、時間戳等)
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

// ErrorAPIResponseTokenMissing 定義 401 Token 缺失錯誤回應結構
// @Description Token 缺失錯誤回應
type ErrorAPIResponseTokenMissing struct {
	// Status 狀態: "error"
	Status string `json:"status" example:"error"`

	// Message 訊息描述
	Message string `json:"message" example:"Unauthorized"`

	// ErrorCode 錯誤代碼: 1006
	ErrorCode int `json:"error_code" example:"1006" swaggertype:"integer"`

	// Meta Meta 資訊 (例如 API 版本、請求 ID 等)
	Meta MetaData `json:"meta"`
}

// ErrorAPIResponseTokenInvalid 定義 401 Token 無效錯誤回應結構
// @Description Token 無效或過期錯誤回應
type ErrorAPIResponseTokenInvalid struct {
	// Status 狀態: "error"
	Status string `json:"status" example:"error"`

	// Message 訊息描述
	Message string `json:"message" example:"Invalid token"`

	// ErrorCode 錯誤代碼: 1007
	ErrorCode int `json:"error_code" example:"1007" swaggertype:"integer"`

	// Meta Meta 資訊 (例如 API 版本、請求 ID 等)
	Meta MetaData `json:"meta"`
}

// ErrorAPIResponseRefreshTokenInvalid 定義 401 Refresh Token 無效錯誤回應結構
// @Description Refresh Token 無效或過期錯誤回應
type ErrorAPIResponseRefreshTokenInvalid struct {
	// Status 狀態: "error"
	Status string `json:"status" example:"error"`

	// Message 訊息描述
	Message string `json:"message" example:"Invalid or expired refresh token"`

	// ErrorCode 錯誤代碼: 1008
	ErrorCode int `json:"error_code" example:"1008" swaggertype:"integer"`

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

// ErrorAPIResponseForbidden 定義 403 禁止存取錯誤回應結構
// @Description 禁止存取（權限不足）錯誤回應
// swagger:model ErrorAPIResponseForbidden
type ErrorAPIResponseForbidden struct {
	// Status 狀態: "error"
	Status string `json:"status" example:"error"`

	// Message 訊息描述
	Message string `json:"message" example:"Forbidden"`

	// ErrorCode 錯誤代碼: 1010
	ErrorCode int `json:"error_code" example:"1010" swaggertype:"integer"`

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
func generateMetaData(c *gin.Context) MetaData {
	return MetaData{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// SuccessResponse 成功回應（單筆資料）
func SuccessResponse(c *gin.Context, message string, data interface{}) {
	response := APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
		Meta:    generateMetaData(c),
	}
	c.JSON(http.StatusOK, response)
}

// CreatedResponse 成功創建回應（201）
func CreatedResponse(c *gin.Context, message string, data interface{}) {
	response := APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
		Meta:    generateMetaData(c),
	}
	c.JSON(http.StatusCreated, response)
}

// ErrorResponse 失敗回應
func ErrorResponse(c *gin.Context, statusCode int, message string, errorCode int) {
	response := APIResponse{
		Status:    "error",
		Message:   message,
		ErrorCode: errorCode,
		Meta:      generateMetaData(c),
	}
	c.JSON(statusCode, response)
}

// FormatValidationError 格式化驗證錯誤訊息
func FormatValidationError(err error) string {
	var messages []string

	// 檢查是否為 validator.ValidationErrors
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			// 取得欄位名稱（使用 json tag，如果沒有則使用欄位名稱）
			field := e.Field()

			// 根據驗證標籤產生錯誤訊息
			var msg string
			switch e.Tag() {
			case "required":
				msg = fmt.Sprintf("%s 為必填欄位", field)
			case "email":
				msg = fmt.Sprintf("%s 格式不正確", field)
			case "min":
				msg = fmt.Sprintf("%s 長度不足（最少 %s 個字元）", field, e.Param())
			case "max":
				msg = fmt.Sprintf("%s 長度過長（最多 %s 個字元）", field, e.Param())
			default:
				msg = fmt.Sprintf("%s 驗證失敗（%s）", field, e.Tag())
			}
			messages = append(messages, msg)
		}
		return strings.Join(messages, "; ")
	}

	// 處理 JSON 解析錯誤（如 EOF、unexpected end of JSON input 等）
	errMsg := err.Error()
	if errMsg == "EOF" || strings.Contains(errMsg, "unexpected end of JSON input") || strings.Contains(errMsg, "invalid character") {
		return "Invalid JSON format"
	}

	// 如果不是驗證錯誤，返回原始錯誤訊息
	return errMsg
}

// ValidationErrorResponse 驗證錯誤回應（使用統一錯誤處理）
func ValidationErrorResponse(c *gin.Context, err error) {
	message := FormatValidationError(err)
	c.Error(NewBadRequestError(message, CodeInvalidInput, err))
}
