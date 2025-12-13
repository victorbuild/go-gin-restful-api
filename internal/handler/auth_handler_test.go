package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"restfulapi/internal/config"
	db "restfulapi/internal/database"
	"restfulapi/internal/middleware"
	"restfulapi/internal/model"
	"restfulapi/internal/util"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRegisterUser_Success(t *testing.T) {
	// Step 1: 設置測試資料庫（SQLite）
	setupTestDB(t)

	// Step 2: Mock RabbitMQ（避免真的發送訊息）
	var publishedMessages []string
	originalFunc := config.PublishMessage
	config.SetPublishMessageFunc(func(queueName string, message string) {
		publishedMessages = append(publishedMessages, message)
	})
	defer config.SetPublishMessageFunc(originalFunc)

	// Step 3: 準備 JSON 請求 Body
	requestBody := map[string]string{
		"name":     "Test User",
		"email":    "test@example.com",
		"password": "123456",
	}
	jsonBody, _ := json.Marshal(requestBody)

	// Step 4: 建立 POST 請求的 Context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	// Step 5: 執行 RegisterUser
	RegisterUser(c)

	// Step 6: 驗證 HTTP 狀態碼（應該是 201 Created）
	assert.Equal(t, http.StatusCreated, w.Code)

	// Step 7: 驗證 JSON 回應結構
	var response util.SuccessAPIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "應該能成功解析 JSON 回應")
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "User created successfully", response.Message)

	// Step 8: 驗證使用者資料
	data, ok := response.Data.(map[string]interface{})
	assert.True(t, ok, "回應的 data 應該是 map")
	assert.Equal(t, "Test User", data["name"])
	assert.Equal(t, "test@example.com", data["email"])
	assert.NotNil(t, data["id"], "ID 應該存在")

	// Step 9: 驗證 RabbitMQ 訊息有被發送（Mock 版本）
	assert.Len(t, publishedMessages, 1, "應該有發送一個 RabbitMQ 訊息")
	if len(publishedMessages) > 0 {
		var messageData map[string]interface{}
		err = json.Unmarshal([]byte(publishedMessages[0]), &messageData)
		assert.NoError(t, err, "RabbitMQ 訊息應該是有效的 JSON")
		assert.Equal(t, "Test User", messageData["name"])
		assert.Equal(t, "test@example.com", messageData["email"])
	}
}

// TestRegisterUser_EmailExists 測試 Email 已存在的情況
func TestRegisterUser_EmailExists(t *testing.T) {
	// Step 1: 設置測試資料庫
	setupTestDB(t)

	// Step 2: 先建立一個使用者（使用 CreateUser，它會自動 hash 密碼）
	existingUser := model.User{
		Name:     "Existing User",
		Email:    "existing@example.com",
		Password: "123456", // CreateUser 會自動 hash
		Role:     "user",
	}
	_, err := model.CreateUser(existingUser)
	assert.NoError(t, err, "應該能成功建立第一個使用者")

	// Step 3: Mock RabbitMQ
	var publishedMessages []string
	originalFunc := config.PublishMessage
	config.SetPublishMessageFunc(func(queueName string, message string) {
		publishedMessages = append(publishedMessages, message)
	})
	defer config.SetPublishMessageFunc(originalFunc)

	// Step 4: 準備 JSON 請求 Body（使用相同的 email）
	requestBody := map[string]string{
		"name":     "New User",
		"email":    "existing@example.com", // 重複的 email
		"password": "123456",
	}
	jsonBody, _ := json.Marshal(requestBody)

	// Step 5: 建立 POST 請求的 Context（需要設置 TraceID 和 ErrorHandler）
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	// 設置 TraceID（ErrorHandler 需要）
	middleware.TraceIDMiddleware()(c)

	// Step 6: 執行 RegisterUser
	RegisterUser(c)

	// Step 7: 處理錯誤（如果有的話）
	if len(c.Errors) > 0 {
		middleware.ErrorHandler()(c)
	}

	// Step 8: 驗證 HTTP 狀態碼（應該是 409 Conflict）
	assert.Equal(t, http.StatusConflict, w.Code, "Email 已存在時應該返回 409")

	// Step 9: 驗證錯誤回應
	var response util.ErrorAPIResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "應該能成功解析錯誤回應")
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, util.CodeEmailExists, response.ErrorCode)

	// Step 9: 驗證 RabbitMQ 訊息沒有被發送（因為註冊失敗）
	assert.Len(t, publishedMessages, 0, "註冊失敗時不應該發送 RabbitMQ 訊息")
}

// TestRegisterUser_InvalidJSON 測試 JSON 格式錯誤的情況
func TestRegisterUser_InvalidJSON(t *testing.T) {
	// Step 1: 設置測試資料庫
	setupTestDB(t)

	// Step 2: Mock RabbitMQ
	var publishedMessages []string
	originalFunc := config.PublishMessage
	config.SetPublishMessageFunc(func(queueName string, message string) {
		publishedMessages = append(publishedMessages, message)
	})
	defer config.SetPublishMessageFunc(originalFunc)

	// Step 3: 準備無效的 JSON（缺少引號）
	invalidJSON := `{"name":"Test User","email":"test@example.com","password":"123456"` // 缺少結尾的 }

	// Step 4: 建立 POST 請求的 Context（需要設置 TraceID 和 ErrorHandler）
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBuffer([]byte(invalidJSON)))
	c.Request.Header.Set("Content-Type", "application/json")

	// 設置 TraceID（ErrorHandler 需要）
	middleware.TraceIDMiddleware()(c)

	// Step 5: 執行 RegisterUser
	RegisterUser(c)

	// Step 6: 處理錯誤（如果有的話）
	if len(c.Errors) > 0 {
		middleware.ErrorHandler()(c)
	}

	// Step 7: 驗證 HTTP 狀態碼（應該是 400 Bad Request）
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Step 7: 驗證錯誤回應
	var response util.ErrorAPIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "應該能成功解析錯誤回應")
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, util.CodeInvalidInput, response.ErrorCode)

	// Step 8: 驗證 RabbitMQ 訊息沒有被發送
	assert.Len(t, publishedMessages, 0, "JSON 錯誤時不應該發送 RabbitMQ 訊息")
}

// TestRegisterUser_PasswordTooLong 測試密碼太長的情況（超過 72 字元）
func TestRegisterUser_PasswordTooLong(t *testing.T) {
	// Step 1: 設置測試資料庫
	setupTestDB(t)

	// Step 2: Mock RabbitMQ
	var publishedMessages []string
	originalFunc := config.PublishMessage
	config.SetPublishMessageFunc(func(queueName string, message string) {
		publishedMessages = append(publishedMessages, message)
	})
	defer config.SetPublishMessageFunc(originalFunc)

	// Step 3: 準備 JSON 請求 Body（密碼超過 72 字元）
	longPassword := strings.Repeat("a", 73) // 73 個字元
	requestBody := map[string]string{
		"name":     "Test User",
		"email":    "test@example.com",
		"password": longPassword,
	}
	jsonBody, _ := json.Marshal(requestBody)

	// Step 4: 建立 POST 請求的 Context（需要設置 TraceID 和 ErrorHandler）
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	// 設置 TraceID（ErrorHandler 需要）
	middleware.TraceIDMiddleware()(c)

	// Step 5: 執行 RegisterUser
	RegisterUser(c)

	// Step 6: 處理錯誤（如果有的話）
	if len(c.Errors) > 0 {
		middleware.ErrorHandler()(c)
	}

	// Step 7: 驗證 HTTP 狀態碼（應該是 400 Bad Request）
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Step 7: 驗證錯誤回應（密碼太長應該返回錯誤）
	// 注意：這裡的錯誤處理可能不同，需要根據實際實現調整
	var response util.ErrorAPIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "應該能成功解析錯誤回應")
	assert.Equal(t, "error", response.Status)

	// Step 8: 驗證 RabbitMQ 訊息沒有被發送
	assert.Len(t, publishedMessages, 0, "密碼太長時不應該發送 RabbitMQ 訊息")
}

// TestRegisterUser_MissingFields 測試缺少必填欄位的情況
func TestRegisterUser_MissingFields(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  map[string]string
		missingField string
	}{
		{
			name: "缺少 name",
			requestBody: map[string]string{
				"email":    "test@example.com",
				"password": "123456",
			},
			missingField: "name",
		},
		{
			name: "缺少 email",
			requestBody: map[string]string{
				"name":     "Test User",
				"password": "123456",
			},
			missingField: "email",
		},
		{
			name: "缺少 password",
			requestBody: map[string]string{
				"name":  "Test User",
				"email": "test@example.com",
			},
			missingField: "password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Step 1: 設置測試資料庫
			setupTestDB(t)

			// Step 2: Mock RabbitMQ
			var publishedMessages []string
			originalFunc := config.PublishMessage
			config.SetPublishMessageFunc(func(queueName string, message string) {
				publishedMessages = append(publishedMessages, message)
			})
			defer config.SetPublishMessageFunc(originalFunc)

			// Step 3: 準備 JSON 請求 Body（缺少欄位）
			jsonBody, _ := json.Marshal(tt.requestBody)

			// Step 4: 建立 POST 請求的 Context（需要設置 TraceID 和 ErrorHandler）
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBuffer(jsonBody))
			c.Request.Header.Set("Content-Type", "application/json")

			// 設置 TraceID（ErrorHandler 需要）
			middleware.TraceIDMiddleware()(c)

			// Step 5: 執行 RegisterUser
			RegisterUser(c)

			// Step 6: 處理錯誤（如果有的話）
			if len(c.Errors) > 0 {
				middleware.ErrorHandler()(c)
			}

			// Step 7: 驗證 HTTP 狀態碼（應該是 400 Bad Request）
			assert.Equal(t, http.StatusBadRequest, w.Code, "缺少 %s 時應該返回 400", tt.missingField)

			// Step 7: 驗證錯誤回應
			var response util.ErrorAPIResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "應該能成功解析錯誤回應")
			assert.Equal(t, "error", response.Status)
			assert.Equal(t, util.CodeInvalidInput, response.ErrorCode, "缺少必填欄位應該是 CodeInvalidInput")

			// Step 8: 驗證 RabbitMQ 訊息沒有被發送
			assert.Len(t, publishedMessages, 0, "缺少欄位時不應該發送 RabbitMQ 訊息")
		})
	}
}

// setupTestDB 設置測試用的 SQLite 記憶體資料庫
func setupTestDB(t *testing.T) {
	t.Helper()
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// 自動遷移 User 表
	err = database.AutoMigrate(&model.User{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// 替換全域變數（測試用）
	originalDB := db.DbConnect
	db.DbConnect = database

	// 清理函數：測試結束後恢復原始資料庫連線
	t.Cleanup(func() {
		sqlDB, _ := database.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
		db.DbConnect = originalDB
	})
}
