package services

import (
	"restfulapi/models"
	"restfulapi/repositories"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockUserRepository 用來模擬 UserRepository，避免 DB 連線
type MockUserRepository struct{}

// 確保 MockUserRepository 符合 IUserRepository 介面
var _ repositories.IUserRepository = (*MockUserRepository)(nil)

// FindAllUsers 模擬回傳假資料
func (m *MockUserRepository) FindAllUsers() []models.User {
	return []models.User{
		{ID: 1, Name: "Alice", Email: "alice@example.com", Role: "user"},
		{ID: 2, Name: "Bob", Email: "bob@example.com", Role: "admin"},
	}
}

// TestGetAllUsers 測試 UserService 的 GetAllUsers 方法
func TestGetAllUsers(t *testing.T) {
	// 初始化 Mock Repository
	mockRepo := &MockUserRepository{}

	// 傳入 Mock Repository 給 UserService
	userService := NewUserService(mockRepo)

	// 執行測試
	users := userService.GetAllUsers()

	// 驗證結果
	assert.Equal(t, 2, len(users), "應該有 2 個使用者")
	assert.Equal(t, "Alice", users[0].Name, "第一個使用者應該是 Alice")
	assert.Equal(t, "Bob", users[1].Name, "第二個使用者應該是 Bob")
}
