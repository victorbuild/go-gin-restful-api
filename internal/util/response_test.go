package util

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

// setupTestContext 建立測試用的 Gin Context
func setupTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	return c, w
}

// TestSuccessResponse 測試成功回應
func TestSuccessResponse(t *testing.T) {
	tests := []struct {
		name    string
		message string
		data    interface{}
	}{
		{
			name:    "成功回應（有資料）",
			message: "User created successfully",
			data:    map[string]string{"id": "1", "name": "Victor"},
		},
		{
			name:    "成功回應（無資料）",
			message: "Operation completed",
			data:    nil,
		},
		{
			name:    "成功回應（空字串訊息）",
			message: "",
			data:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := setupTestContext()
			SuccessResponse(c, tt.message, tt.data)

			assert.Equal(t, http.StatusOK, w.Code)

			var response APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "success", response.Status)
			assert.Equal(t, tt.message, response.Message)
			// 對於 data，使用更寬鬆的比較（因為 JSON 反序列化會改變類型）
			if tt.data != nil {
				assert.NotNil(t, response.Data)
			} else {
				assert.Nil(t, response.Data)
			}
		})
	}
}

// TestCreatedResponse 測試創建回應（201）
func TestCreatedResponse(t *testing.T) {
	tests := []struct {
		name    string
		message string
		data    interface{}
	}{
		{
			name:    "創建回應（有資料）",
			message: "User created",
			data:    map[string]interface{}{"id": 1, "email": "test@example.com"},
		},
		{
			name:    "創建回應（無資料）",
			message: "Resource created",
			data:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := setupTestContext()
			CreatedResponse(c, tt.message, tt.data)

			assert.Equal(t, http.StatusCreated, w.Code)

			var response APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "success", response.Status)
			assert.Equal(t, tt.message, response.Message)
			// 對於 data，使用更寬鬆的比較（因為 JSON 反序列化會改變類型）
			if tt.data != nil {
				assert.NotNil(t, response.Data)
			} else {
				assert.Nil(t, response.Data)
			}
		})
	}
}

// TestErrorResponse 測試錯誤回應
func TestErrorResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		message    string
		errorCode  int
	}{
		{
			name:       "400 錯誤回應",
			statusCode: http.StatusBadRequest,
			message:    "Invalid input",
			errorCode:  CodeInvalidInput,
		},
		{
			name:       "404 錯誤回應",
			statusCode: http.StatusNotFound,
			message:    "User not found",
			errorCode:  CodeUserNotFound,
		},
		{
			name:       "500 錯誤回應",
			statusCode: http.StatusInternalServerError,
			message:    "Internal server error",
			errorCode:  CodeInternalError,
		},
		{
			name:       "401 錯誤回應",
			statusCode: http.StatusUnauthorized,
			message:    "Unauthorized",
			errorCode:  CodeTokenMissing,
		},
		{
			name:       "403 錯誤回應",
			statusCode: http.StatusForbidden,
			message:    "Forbidden",
			errorCode:  CodeForbidden,
		},
		{
			name:       "409 錯誤回應",
			statusCode: http.StatusConflict,
			message:    "Email already registered",
			errorCode:  CodeEmailExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := setupTestContext()
			ErrorResponse(c, tt.statusCode, tt.message, tt.errorCode)

			assert.Equal(t, tt.statusCode, w.Code)

			var response APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "error", response.Status)
			assert.Equal(t, tt.message, response.Message)
			assert.Equal(t, tt.errorCode, response.ErrorCode)
		})
	}
}

// TestFormatValidationError 測試格式化驗證錯誤
func TestFormatValidationError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "驗證錯誤 - 空陣列",
			err:      validator.ValidationErrors{},
			expected: "", // 空的 ValidationErrors 會回傳空字串
		},
		{
			name:     "非驗證錯誤 - 一般錯誤",
			err:      assert.AnError,
			expected: assert.AnError.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatValidationError(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestFormatValidationError_JSONErrors 測試 JSON 相關錯誤
func TestFormatValidationError_JSONErrors(t *testing.T) {
	// 使用 errors.New 來模擬 JSON 錯誤訊息
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "EOF 錯誤",
			err:      errors.New("EOF"),
			expected: "Invalid JSON format",
		},
		{
			name:     "unexpected end of JSON input",
			err:      errors.New("unexpected end of JSON input"),
			expected: "Invalid JSON format",
		},
		{
			name:     "invalid character",
			err:      errors.New("invalid character 'x' looking for beginning of value"),
			expected: "Invalid JSON format",
		},
		{
			name:     "一般錯誤訊息",
			err:      assert.AnError,
			expected: assert.AnError.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatValidationError(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestValidationErrorResponse 測試驗證錯誤回應
func TestValidationErrorResponse(t *testing.T) {
	c, _ := setupTestContext()

	// 建立一個測試用的錯誤
	testErr := assert.AnError
	ValidationErrorResponse(c, testErr)

	// 驗證是否有設定錯誤到 context
	assert.True(t, len(c.Errors) > 0, "應該有錯誤被設定到 context")

	// 驗證錯誤類型
	err := c.Errors[0].Err
	appErr, ok := err.(*AppError)
	assert.True(t, ok, "錯誤應該是 AppError 類型")
	assert.Equal(t, http.StatusBadRequest, appErr.StatusCode)
	assert.Equal(t, CodeInvalidInput, appErr.Code)
}

// TestErrorAPIResponse_MarshalJSON 測試 ErrorAPIResponse 的 JSON 序列化
func TestErrorAPIResponse_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		response ErrorAPIResponse
	}{
		{
			name: "完整錯誤回應",
			response: ErrorAPIResponse{
				Status:    "error",
				Message:   "Test error",
				ErrorCode: CodeInvalidInput,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.response)
			assert.NoError(t, err)

			var result ErrorAPIResponse
			err = json.Unmarshal(data, &result)
			assert.NoError(t, err)
			assert.Equal(t, tt.response.Status, result.Status)
			assert.Equal(t, tt.response.Message, result.Message)
			assert.Equal(t, tt.response.ErrorCode, result.ErrorCode)
		})
	}
}
