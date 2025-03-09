package utils

const (
	// Auth 錯誤碼（1000-1999）
	ErrInvalidInput          = 1000
	ErrEmailExists           = 1001
	ErrInvalidCredentials    = 1005
	ErrTokenGenerationFailed = 1006

	// Server 錯誤碼（4000-4999）
	ErrInternalError    = 4001
	ErrPasswordHashFail = 4002 // ✅ 密碼加密錯誤
	ErrDatabaseError    = 4003 // ✅ 資料庫錯誤
)
