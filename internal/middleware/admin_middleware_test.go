package middleware

import (
	"net/http"
	"testing"

	"restfulapi/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestRequireAdmin_WithAdminRole 測試當 role 為 admin 時，應該允許通過
func TestRequireAdmin_WithAdminRole(t *testing.T) {
	c, w := setupTestContext()

	// 設定 role 為 admin
	c.Set("role", "admin")

	// 設定一個 handler 來驗證 middleware 是否允許通過
	called := false
	handler := func(c *gin.Context) {
		called = true
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}

	// 執行 middleware 和 handler
	RequireAdmin()(c)
	if !c.IsAborted() {
		handler(c)
	}

	// 驗證 handler 被執行（表示 middleware 允許通過）
	assert.True(t, called, "當 role 為 admin 時，應該允許通過")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 0, len(c.Errors), "不應該有錯誤")
}

// TestRequireAdmin_WithoutRole 測試當 role 不存在時，應該返回 403 Forbidden
func TestRequireAdmin_WithoutRole(t *testing.T) {
	c, _ := setupTestContext()

	// 不設定 role

	// 執行 middleware
	RequireAdmin()(c)

	// 驗證請求被中止
	assert.True(t, c.IsAborted(), "當 role 不存在時，應該中止請求")

	// 驗證有錯誤被設定
	assert.True(t, len(c.Errors) > 0, "應該有錯誤被設定")

	// 驗證錯誤類型
	err := c.Errors[0].Err
	appErr, ok := err.(*util.AppError)
	assert.True(t, ok, "錯誤應該是 AppError 類型")
	assert.Equal(t, http.StatusForbidden, appErr.StatusCode)
	assert.Equal(t, util.CodeForbidden, appErr.Code)
	assert.Equal(t, "Forbidden", appErr.Message)
}

// TestRequireAdmin_WithNonAdminRole 測試當 role 不是 admin 時，應該返回 403 Forbidden
func TestRequireAdmin_WithNonAdminRole(t *testing.T) {
	tests := []struct {
		name string
		role string
	}{
		{
			name: "role 為 user",
			role: "user",
		},
		{
			name: "role 為空字串",
			role: "",
		},
		{
			name: "role 為 guest",
			role: "guest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := setupTestContext()

			// 設定 role 為非 admin
			c.Set("role", tt.role)

			// 執行 middleware
			RequireAdmin()(c)

			// 驗證請求被中止
			assert.True(t, c.IsAborted(), "當 role 不是 admin 時，應該中止請求")

			// 驗證有錯誤被設定
			assert.True(t, len(c.Errors) > 0, "應該有錯誤被設定")

			// 驗證錯誤類型
			err := c.Errors[0].Err
			appErr, ok := err.(*util.AppError)
			assert.True(t, ok, "錯誤應該是 AppError 類型")
			assert.Equal(t, http.StatusForbidden, appErr.StatusCode)
			assert.Equal(t, util.CodeForbidden, appErr.Code)
			assert.Equal(t, "Forbidden", appErr.Message)
		})
	}
}
