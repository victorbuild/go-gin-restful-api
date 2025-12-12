package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// TestUser_CheckPassword 測試密碼驗證功能
func TestUser_CheckPassword(t *testing.T) {
	// 原始的密碼
	password := "testPassword123"
	// hash原始密碼
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err, "產生加密密碼不應該失敗")
	// 疑問：DefaultCost 是要幹嘛？

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
// 額外挑戰（可選）
// ==========================================

// TestUser_CheckPassword_SpecialCharacters 測試特殊字元密碼
// 提示：測試包含特殊字元的密碼是否能正確驗證
// 例如：密碼包含 !@#$%^&*() 等特殊字元

// TestUser_CheckPassword_LongPassword 測試長密碼
// 提示：測試非常長的密碼（例如 100 個字元）

// TestUser_CheckPassword_Unicode 測試 Unicode 字元密碼
// 提示：測試包含中文、日文、emoji 等 Unicode 字元的密碼
