package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"restfulapi/internal/config"
	db "restfulapi/internal/database"
	"restfulapi/internal/model"
	"restfulapi/internal/util"
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
