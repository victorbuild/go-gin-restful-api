package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// TestTraceIDMiddleware_GenerateNewTraceID 測試生成新的 Trace ID
func TestTraceIDMiddleware_GenerateNewTraceID(t *testing.T) {
	c, w := setupTestContext()

	// 執行 middleware
	TraceIDMiddleware()(c)

	// 驗證 response header 有 Trace ID
	traceID := w.Header().Get(TraceIDHeader)
	assert.NotEmpty(t, traceID, "應該生成 Trace ID")

	// 驗證 Trace ID 格式是 UUID
	_, err := uuid.Parse(traceID)
	assert.NoError(t, err, "Trace ID 應該是有效的 UUID")

	// 驗證 context 中有 Trace ID
	contextTraceID := GetTraceID(c)
	assert.Equal(t, traceID, contextTraceID, "Context 中的 Trace ID 應該與 header 一致")
}

// TestTraceIDMiddleware_UseExistingTraceID 測試使用現有的 Trace ID
func TestTraceIDMiddleware_UseExistingTraceID(t *testing.T) {
	c, w := setupTestContext()

	// 設定 request header 中的 Trace ID
	existingTraceID := "existing-trace-id-12345"
	c.Request.Header.Set(TraceIDHeader, existingTraceID)

	// 執行 middleware
	TraceIDMiddleware()(c)

	// 驗證使用現有的 Trace ID
	traceID := w.Header().Get(TraceIDHeader)
	assert.Equal(t, existingTraceID, traceID, "應該使用 request header 中的 Trace ID")

	// 驗證 context 中的 Trace ID
	contextTraceID := GetTraceID(c)
	assert.Equal(t, existingTraceID, contextTraceID, "Context 中的 Trace ID 應該與 header 一致")
}

// TestGetTraceID 測試 GetTraceID 輔助函數
func TestGetTraceID(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*gin.Context)
		expected string
	}{
		{
			name: "有 Trace ID",
			setup: func(c *gin.Context) {
				c.Set(TraceIDKey, "test-trace-id")
			},
			expected: "test-trace-id",
		},
		{
			name: "沒有 Trace ID",
			setup: func(c *gin.Context) {
				// 不設定 Trace ID
			},
			expected: "",
		},
		{
			name: "Trace ID 類型錯誤",
			setup: func(c *gin.Context) {
				c.Set(TraceIDKey, 123) // 設定錯誤的類型
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := setupTestContext()
			tt.setup(c)

			result := GetTraceID(c)
			assert.Equal(t, tt.expected, result)
		})
	}
}
