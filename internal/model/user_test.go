package model

import (
	"strings"
	"testing"

	db "restfulapi/internal/database"
	"restfulapi/internal/util"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestUser_CheckPassword 測試密碼驗證功能
func TestUser_CheckPassword(t *testing.T) {
	// 原始的密碼
	password := "testPassword123"
	// hash原始密碼
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err, "產生加密密碼不應該失敗")

	// 建立一個 User 實例，設定加密後的密碼
	user := User{
		ID:       1,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: string(hashedPassword), // 轉換成 string
		Role:     "user",
	}

	tests := []struct {
		name           string
		inputPassword  string
		expectedResult bool
	}{
		{
			name:           "正確的密碼",
			inputPassword:  "testPassword123", // 使用上面定義的原始密碼
			expectedResult: true,
		},
		{
			name:           "錯誤的密碼",
			inputPassword:  "wrongPassword",
			expectedResult: false,
		},
		{
			name:           "空密碼",
			inputPassword:  "",
			expectedResult: false,
		},
		{
			name:           "大小寫不同的密碼",
			inputPassword:  "TestPassword123", // 注意大小寫
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := user.CheckPassword(tt.inputPassword)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

// TestUser_CheckPassword_EmptyHash 測試邊界情況：空密碼雜湊值
// 確保當 User.Password 是空字串時，CheckPassword 不會崩潰並正確回傳 false
func TestUser_CheckPassword_EmptyHash(t *testing.T) {
	user := User{
		ID:       1,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "", // 空密碼雜湊
		Role:     "user",
	}

	// 測試任何密碼都應該回傳 false
	result := user.CheckPassword("anyPassword")
	assert.False(t, result, "空密碼雜湊應該回傳 false")

	// 測試空輸入密碼也應該回傳 false
	result2 := user.CheckPassword("")
	assert.False(t, result2, "空輸入密碼應該回傳 false")
}

// TestUser_CheckPassword_SQLInjection 測試 SQL 注入攻擊嘗試
// 確保系統不會被 SQL 注入攻擊，即使攻擊者嘗試使用 SQL 注入語法作為密碼
func TestUser_CheckPassword_SQLInjection(t *testing.T) {
	// 準備正常的使用者
	password := "testPassword123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err, "產生加密密碼不應該失敗")

	user := User{
		ID:       1,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Role:     "user",
	}

	// SQL 注入攻擊嘗試
	sqlInjectionAttempts := []string{
		"' OR '1'='1",
		"'; DROP TABLE users; --",
		"' OR 1=1--",
		"admin'--",
		"' UNION SELECT * FROM users--",
		"1' OR '1'='1",
		"admin' OR '1'='1'--",
	}

	for _, attempt := range sqlInjectionAttempts {
		t.Run("SQL注入嘗試: "+attempt, func(t *testing.T) {
			result := user.CheckPassword(attempt)
			// 應該回傳 false，不應該被當作正確密碼
			assert.False(t, result, "SQL 注入嘗試應該失敗，不應該被當作正確密碼")
		})
	}

	// 確保正確的密碼仍然可以通過驗證
	result := user.CheckPassword(password)
	assert.True(t, result, "正確的密碼應該能通過驗證")
}

// TestUser_CheckPassword_VeryLongPassword 測試非常長的密碼（DoS 攻擊嘗試）
// 確保系統不會因為非常長的密碼而崩潰
// 注意：bcrypt 有 72 字元的限制，超過會回傳錯誤
func TestUser_CheckPassword_VeryLongPassword(t *testing.T) {
	testCases := []struct {
		name          string
		length        int
		shouldSucceed bool
		description   string
	}{
		{"正常長度", 20, true, "正常長度的密碼應該能正常加密和驗證"},
		{"接近限制", 70, true, "接近 72 字元限制的密碼應該能正常處理"},
		{"達到限制", 72, true, "達到 72 字元限制的密碼應該能正常處理"},
		{"超過限制", 100, false, "超過 72 字元的密碼 bcrypt 會回傳錯誤"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			password := strings.Repeat("a", tt.length)
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

			if tt.shouldSucceed {
				// 應該能成功加密
				assert.NoError(t, err, tt.description)

				user := User{
					Password: string(hashedPassword),
				}

				// 正確的密碼應該能通過驗證
				result := user.CheckPassword(password)
				assert.True(t, result, "正確的長密碼應該能通過驗證")

				// 錯誤的密碼應該失敗（使用不同的錯誤密碼，避免超過 72 字元限制）
				wrongPassword := password[:len(password)/2] + "wrong"
				wrongResult := user.CheckPassword(wrongPassword)
				assert.False(t, wrongResult, "錯誤的長密碼應該失敗")
			} else {
				// bcrypt 會回傳錯誤（超過 72 字元）
				assert.Error(t, err, "超過 72 字元的密碼 bcrypt 應該回傳錯誤")
				assert.Contains(t, err.Error(), "password length exceeds 72 bytes", "錯誤訊息應該提到長度限制")

				// 即使加密失敗，CheckPassword 也應該能處理（不會 panic）
				// 測試空雜湊值的情況
				user := User{
					Password: "", // 空雜湊值
				}
				result := user.CheckPassword(password)
				assert.False(t, result, "空雜湊值應該回傳 false，不會 panic")
			}
		})
	}
}

// TestUser_CheckPassword_SpecialCharacters 測試特殊字元和 Unicode 密碼
// 確保系統能正確處理包含特殊字元、Unicode、emoji 的密碼
func TestUser_CheckPassword_SpecialCharacters(t *testing.T) {
	specialPasswords := []string{
		"!@#$%^&*()_+-=[]{}|;:,.<>?",
		"密碼123🔐",
		"パスワード123",
		"пароль123",
		"كلمة المرور123",
		"password with spaces",
		"password\twith\ttabs",
		"password\nwith\nnewlines",
	}

	for _, password := range specialPasswords {
		t.Run("特殊字元: "+password, func(t *testing.T) {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			assert.NoError(t, err, "特殊字元密碼應該能正常加密")

			user := User{
				Password: string(hashedPassword),
			}

			// 正確的密碼應該能驗證通過
			result := user.CheckPassword(password)
			assert.True(t, result, "特殊字元密碼應該能正確驗證")

			// 錯誤的密碼應該失敗
			wrongResult := user.CheckPassword(password + "wrong")
			assert.False(t, wrongResult, "錯誤的特殊字元密碼應該失敗")
		})
	}
}

// ==========================================
// CreateUser 測試
// ==========================================

// setupTestDB 設置測試用的 SQLite 記憶體資料庫
func setupTestDB(t *testing.T) *gorm.DB {
	// 使用 SQLite 記憶體資料庫（測試後自動清除）
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// 定義測試用的 User 結構（SQLite 不支援 bigserial，改用 autoIncrement）
	// 結構與 User 相同，只是 ID 欄位使用 autoIncrement 而不是 bigserial
	type TestUser struct {
		ID       uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
		Name     string `gorm:"size:255" json:"name"`
		Password string `gorm:"size:255" json:"password"`
		Email    string `gorm:"uniqueIndex;size:255" json:"email"`
		Role     string `gorm:"size:255" json:"role"`
	}

	// 自動遷移 User 表（使用 TestUser 結構，但表名還是 users）
	err = database.Table("users").AutoMigrate(&TestUser{})
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

	return database
}

// TestCreateUser_Success 測試成功建立使用者
func TestCreateUser_Success(t *testing.T) {
	setupTestDB(t)

	user := User{
		Name:     "Test User",
		Email:    "test-create-success@example.com",
		Password: "testPassword123",
		Role:     "user",
	}

	userID, err := CreateUser(user)

	// 驗證結果
	assert.NoError(t, err, "建立使用者應該成功")
	assert.Greater(t, userID, uint64(0), "使用者 ID 應該大於 0")

	// 驗證資料庫中真的有這個使用者
	var createdUser User
	result := db.DbConnect.Where("id = ?", userID).First(&createdUser)
	assert.NoError(t, result.Error, "應該能在資料庫中找到建立的使用者")
	assert.Equal(t, user.Name, createdUser.Name, "姓名應該一致")
	assert.Equal(t, user.Email, createdUser.Email, "Email 應該一致")
	assert.Equal(t, user.Role, createdUser.Role, "角色應該一致")
	// 驗證密碼已經被加密（不是原始密碼）
	assert.NotEqual(t, "testPassword123", createdUser.Password, "密碼應該被加密")
	assert.NotEmpty(t, createdUser.Password, "密碼不應該為空")
}

// TestCreateUser_EmailExists 測試 Email 已存在的情況
func TestCreateUser_EmailExists(t *testing.T) {
	setupTestDB(t)

	// 先建立第一個使用者
	user1 := User{
		Name:     "First User",
		Email:    "test-email-exists@example.com",
		Password: "password123",
		Role:     "user",
	}
	userID1, err1 := CreateUser(user1)
	assert.NoError(t, err1, "第一個使用者應該建立成功")
	assert.Greater(t, userID1, uint64(0))

	// 嘗試用相同 Email 建立第二個使用者
	user2 := User{
		Name:     "Second User",
		Email:    "test-email-exists@example.com", // 相同的 Email
		Password: "password456",
		Role:     "user",
	}
	userID2, err2 := CreateUser(user2)

	// 應該返回錯誤
	assert.Error(t, err2, "Email 已存在時應該返回錯誤")
	assert.Equal(t, util.ErrEmailExists, err2, "錯誤應該是 ErrEmailExists")
	assert.Equal(t, uint64(0), userID2, "使用者 ID 應該是 0")

	// 驗證資料庫中只有一個使用者（第一個）
	var count int64
	db.DbConnect.Model(&User{}).Where("email = ?", "test-email-exists@example.com").Count(&count)
	assert.Equal(t, int64(1), count, "資料庫中應該只有一個使用者")
}

// TestCreateUser_PasswordEncryption 測試密碼加密功能
func TestCreateUser_PasswordEncryption(t *testing.T) {
	setupTestDB(t)

	user := User{
		Name:     "Test User",
		Email:    "test-password-encryption@example.com",
		Password: "originalPassword123",
		Role:     "user",
	}

	userID, err := CreateUser(user)
	assert.NoError(t, err)
	assert.Greater(t, userID, uint64(0))

	// 從資料庫讀取使用者
	var createdUser User
	db.DbConnect.Where("id = ?", userID).First(&createdUser)

	// 驗證密碼已經被加密
	assert.NotEqual(t, "originalPassword123", createdUser.Password, "密碼應該被加密，不應該是原始密碼")
	assert.NotEmpty(t, createdUser.Password, "密碼不應該為空")

	// 驗證可以使用 CheckPassword 驗證原始密碼
	assert.True(t, createdUser.CheckPassword("originalPassword123"), "應該可以用原始密碼驗證")
	assert.False(t, createdUser.CheckPassword("wrongPassword"), "錯誤的密碼應該驗證失敗")
}

// TestCreateUser_MultipleUsers 測試建立多個使用者
func TestCreateUser_MultipleUsers(t *testing.T) {
	setupTestDB(t)

	users := []User{
		{
			Name:     "User 1",
			Email:    "user1@example.com",
			Password: "password1",
			Role:     "user",
		},
		{
			Name:     "User 2",
			Email:    "user2@example.com",
			Password: "password2",
			Role:     "admin",
		},
		{
			Name:     "User 3",
			Email:    "user3@example.com",
			Password: "password3",
			Role:     "user",
		},
	}

	// 建立多個使用者
	userIDs := make([]uint64, len(users))
	for i, user := range users {
		userID, err := CreateUser(user)
		assert.NoError(t, err, "建立使用者 %d 應該成功", i+1)
		assert.Greater(t, userID, uint64(0), "使用者 ID 應該大於 0")
		userIDs[i] = userID
	}

	// 驗證所有使用者 ID 都不同
	for i := 0; i < len(userIDs); i++ {
		for j := i + 1; j < len(userIDs); j++ {
			assert.NotEqual(t, userIDs[i], userIDs[j], "使用者 ID 應該不同")
		}
	}

	// 驗證資料庫中有正確數量的使用者
	var count int64
	db.DbConnect.Model(&User{}).Count(&count)
	assert.Equal(t, int64(3), count, "資料庫中應該有 3 個使用者")
}

// TestCreateUser_PasswordHashFail 測試密碼加密失敗的情況
// 當密碼超過 72 字元時，bcrypt.GenerateFromPassword 會回傳錯誤
func TestCreateUser_PasswordHashFail(t *testing.T) {
	setupTestDB(t)

	// 建立一個密碼超過 72 字元的使用者
	user := User{
		Name:     "Test User",
		Email:    "test-password-hash-fail@example.com",
		Password: strings.Repeat("a", 73), // 超過 72 字元限制
		Role:     "user",
	}

	userID, err := CreateUser(user)

	// 應該返回錯誤
	assert.Error(t, err, "密碼超過 72 字元時應該返回錯誤")
	assert.Equal(t, util.ErrPasswordHashFail, err, "錯誤應該是 ErrPasswordHashFail")
	assert.Equal(t, uint64(0), userID, "使用者 ID 應該是 0")

	// 驗證資料庫中沒有建立這個使用者
	var count int64
	db.DbConnect.Model(&User{}).Where("email = ?", "test-password-hash-fail@example.com").Count(&count)
	assert.Equal(t, int64(0), count, "資料庫中不應該有這個使用者")
}

// TestCreateUser_DatabaseError 測試資料庫錯誤的情況
// 通過關閉資料庫連線來模擬資料庫錯誤
func TestCreateUser_DatabaseError(t *testing.T) {
	setupTestDB(t)

	// 關閉資料庫連線來模擬資料庫錯誤
	sqlDB, _ := db.DbConnect.DB()
	sqlDB.Close()

	user := User{
		Name:     "Test User",
		Email:    "test-database-error@example.com",
		Password: "testPassword123",
		Role:     "user",
	}

	userID, err := CreateUser(user)

	// 應該返回錯誤
	assert.Error(t, err, "資料庫錯誤時應該返回錯誤")
	assert.Equal(t, util.ErrDatabaseError, err, "錯誤應該是 ErrDatabaseError")
	assert.Equal(t, uint64(0), userID, "使用者 ID 應該是 0")
}
