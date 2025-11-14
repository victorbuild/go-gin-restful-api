package utils

import "errors"

// 錯誤碼常數（用於 API 回應）
const (
	// Auth 錯誤碼（1000-1099）
	CodeUnsupportedMediaType  = 1000 // 不支援的媒體類型（Content-Type 錯誤）
	CodeInvalidInput          = 1001 // JSON 格式錯誤
	CodeMissingRequiredFields = 1002 // 缺少必填欄位
	CodeEmailExists           = 1003 // Email 已存在
	CodeInvalidCredentials    = 1004 // 帳號密碼錯誤
	CodeTokenGenerationFailed = 1005 // Token 產生失敗
	CodeTokenMissing          = 1006 // Token 缺失
	CodeTokenInvalid          = 1007 // Token 無效或過期
	CodeRefreshTokenInvalid   = 1008 // Refresh Token 無效或過期
	CodeUserNotFound          = 1009 // 使用者不存在
	CodeForbidden             = 1010 // 禁止存取（權限不足）
	CodeDeleteUserFailed      = 1011 // 刪除使用者失敗
	CodeUpdateUserFailed      = 1012 // 更新使用者失敗
	CodeEmailInUse            = 1013 // Email 已被使用

	// Server 錯誤碼（4000-4999）
	CodeInternalError    = 4001 // 伺服器內部錯誤
	CodePasswordHashFail = 4002 // 密碼加密錯誤
	CodeDatabaseError    = 4003 // 資料庫錯誤
)

// 自訂錯誤
var (
	ErrUserNotFound     = errors.New("user not found")
	ErrEmailExists      = errors.New("email already registered")
	ErrPasswordHashFail = errors.New("password encryption failed")
	ErrDatabaseError    = errors.New("database error")
)
