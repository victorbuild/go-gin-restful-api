package utils

const (
	// Auth 錯誤碼（1000-1099）
	ErrUnsupportedMediaType  = 1000 // 不支援的媒體類型（Content-Type 錯誤）
	ErrInvalidInput          = 1001 // JSON 格式錯誤
	ErrMissingRequiredFields = 1002 // 缺少必填欄位
	ErrEmailExists           = 1003 // Email 已存在
	ErrInvalidCredentials    = 1004 // 帳號密碼錯誤
	ErrTokenGenerationFailed = 1005 // Token 產生失敗

	// Server 錯誤碼（4000-4999）
	ErrInternalError    = 4001 // 伺服器內部錯誤
	ErrPasswordHashFail = 4002 // 密碼加密錯誤
	ErrDatabaseError    = 4003 // 資料庫錯誤
)
