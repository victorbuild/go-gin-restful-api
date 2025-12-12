package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
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

// TestHealthCheck 測試健康檢查端點
func TestHealthCheck(t *testing.T) {
	c, w := setupTestContext()

	// 1. 調用 HealthCheck handler
	HealthCheck(c)

	// 2. 驗證 HTTP 狀態碼應該是 200
	assert.Equal(t, http.StatusOK, w.Code)

	// 3. 解析 JSON 回應到 HealthResponse 結構
	var response HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	//  4. 驗證回應內容：
	assert.Equal(t, "user-api", response.Service)
	assert.Equal(t, "healthy", response.Status)
	assert.Equal(t, "1.0", response.Version)

	assert.NotEmpty(t, response.Timestamp)
	// 5. 驗證 timestamp 格式
	_, err = time.Parse(time.RFC3339, response.Timestamp)
	assert.NoError(t, err, "timestamp 應該是有效的 RFC3339 格式")
}

// TestRoot 測試根路徑端點
func TestRoot(t *testing.T) {
	c, w := setupTestContext()

	// 1. 調用 Root handler
	Root(c)

	// 2. 驗證 HTTP 狀態碼應該是 200
	assert.Equal(t, http.StatusOK, w.Code)

	// 3. 解析 JSON 回應到 HealthRootResponse 結構
	var response HealthRootResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 4. 驗證 Status 應該是 "ok"
	assert.Equal(t, "ok", response.Status)
}
