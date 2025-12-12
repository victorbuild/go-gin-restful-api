package util

import (
	"errors"
	"net/http"
)

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
	CodePasswordTooLong  = 4004 // 密碼太長（超過 bcrypt 72 字元限制）
)

// 自訂錯誤
var (
	ErrUserNotFound     = errors.New("user not found")
	ErrEmailExists      = errors.New("email already registered")
	ErrPasswordHashFail = errors.New("password encryption failed")
	ErrDatabaseError    = errors.New("database error")
)

// AppError 應用程式錯誤結構
type AppError struct {
	StatusCode int    // HTTP 狀態碼
	Code       int    // 應用程式錯誤碼
	Message    string // 錯誤訊息
	Err        error  // 原始錯誤
}

// Error 實作 error 介面
func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

// NewNotFoundError 建立 404 錯誤
func NewNotFoundError(message string, code int, err error) *AppError {
	return &AppError{
		StatusCode: http.StatusNotFound,
		Code:       code,
		Message:    message,
		Err:        err,
	}
}

// NewBadRequestError 建立 400 錯誤
func NewBadRequestError(message string, code int, err error) *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		Code:       code,
		Message:    message,
		Err:        err,
	}
}

// NewInternalServerError 建立 500 錯誤
func NewInternalServerError(message string, code int, err error) *AppError {
	return &AppError{
		StatusCode: http.StatusInternalServerError,
		Code:       code,
		Message:    message,
		Err:        err,
	}
}

// NewForbiddenError 建立 403 錯誤
func NewForbiddenError(message string, code int, err error) *AppError {
	return &AppError{
		StatusCode: http.StatusForbidden,
		Code:       code,
		Message:    message,
		Err:        err,
	}
}

// NewUnauthorizedError 建立 401 錯誤
func NewUnauthorizedError(message string, code int, err error) *AppError {
	return &AppError{
		StatusCode: http.StatusUnauthorized,
		Code:       code,
		Message:    message,
		Err:        err,
	}
}

// NewConflictError 建立 409 錯誤
func NewConflictError(message string, code int, err error) *AppError {
	return &AppError{
		StatusCode: http.StatusConflict,
		Code:       code,
		Message:    message,
		Err:        err,
	}
}

// NewUnsupportedMediaTypeError 建立 415 錯誤
func NewUnsupportedMediaTypeError(message string, code int, err error) *AppError {
	return &AppError{
		StatusCode: http.StatusUnsupportedMediaType,
		Code:       code,
		Message:    message,
		Err:        err,
	}
}

// ValidatePasswordLength 驗證密碼長度（bcrypt 限制 72 字元）
// 如果密碼為空或超過 72 字元，回傳對應的錯誤
// 如果密碼長度符合要求，回傳 nil
func ValidatePasswordLength(password string) *AppError {
	// 空密碼的情況由其他驗證處理（required 標籤），這裡不處理
	if password == "" {
		return nil
	}

	// bcrypt 限制 72 字元，避免 CVE-2025-22228 相關問題
	if len(password) > 72 {
		return NewBadRequestError(
			"Password too long (maximum 72 characters)",
			CodePasswordTooLong,
			nil,
		)
	}

	return nil
}
