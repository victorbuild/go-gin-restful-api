package util

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAppError_Error 測試 AppError 的 Error() 方法
func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name     string
		appError *AppError
		expected string
	}{
		{
			name: "有原始錯誤時，回傳原始錯誤訊息",
			appError: &AppError{
				StatusCode: http.StatusNotFound,
				Code:       CodeUserNotFound,
				Message:    "User not found",
				Err:        errors.New("database query failed"),
			},
			expected: "database query failed",
		},
		{
			name: "沒有原始錯誤時，回傳 Message",
			appError: &AppError{
				StatusCode: http.StatusNotFound,
				Code:       CodeUserNotFound,
				Message:    "User not found",
				Err:        nil,
			},
			expected: "User not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.appError.Error()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestNewNotFoundError 測試建立 404 錯誤
func TestNewNotFoundError(t *testing.T) {
	tests := []struct {
		name           string
		message        string
		code           int
		err            error
		expectedStatus int
		expectedCode   int
		expectedMsg    string
	}{
		{
			name:           "建立 404 錯誤（有原始錯誤）",
			message:        "User not found",
			code:           CodeUserNotFound,
			err:            errors.New("database error"),
			expectedStatus: http.StatusNotFound,
			expectedCode:   CodeUserNotFound,
			expectedMsg:    "User not found",
		},
		{
			name:           "建立 404 錯誤（無原始錯誤）",
			message:        "Resource not found",
			code:           CodeUserNotFound,
			err:            nil,
			expectedStatus: http.StatusNotFound,
			expectedCode:   CodeUserNotFound,
			expectedMsg:    "Resource not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appError := NewNotFoundError(tt.message, tt.code, tt.err)

			assert.NotNil(t, appError)
			assert.Equal(t, tt.expectedStatus, appError.StatusCode)
			assert.Equal(t, tt.expectedCode, appError.Code)
			assert.Equal(t, tt.expectedMsg, appError.Message)
			assert.Equal(t, tt.err, appError.Err)
		})
	}
}

// TestNewBadRequestError 測試建立 400 錯誤
func TestNewBadRequestError(t *testing.T) {
	tests := []struct {
		name           string
		message        string
		code           int
		err            error
		expectedStatus int
		expectedCode   int
		expectedMsg    string
	}{
		{
			name:           "建立 400 錯誤（有原始錯誤）",
			message:        "Invalid input",
			code:           CodeInvalidInput,
			err:            errors.New("json decode error"),
			expectedStatus: http.StatusBadRequest,
			expectedCode:   CodeInvalidInput,
			expectedMsg:    "Invalid input",
		},
		{
			name:           "建立 400 錯誤（無原始錯誤）",
			message:        "Missing required fields",
			code:           CodeMissingRequiredFields,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   CodeMissingRequiredFields,
			expectedMsg:    "Missing required fields",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appError := NewBadRequestError(tt.message, tt.code, tt.err)

			assert.NotNil(t, appError)
			assert.Equal(t, tt.expectedStatus, appError.StatusCode)
			assert.Equal(t, tt.expectedCode, appError.Code)
			assert.Equal(t, tt.expectedMsg, appError.Message)
			assert.Equal(t, tt.err, appError.Err)
		})
	}
}

// TestNewInternalServerError 測試建立 500 錯誤
func TestNewInternalServerError(t *testing.T) {
	tests := []struct {
		name           string
		message        string
		code           int
		err            error
		expectedStatus int
		expectedCode   int
		expectedMsg    string
	}{
		{
			name:           "建立 500 錯誤（有原始錯誤）",
			message:        "Internal server error",
			code:           CodeInternalError,
			err:            errors.New("database connection failed"),
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   CodeInternalError,
			expectedMsg:    "Internal server error",
		},
		{
			name:           "建立 500 錯誤（無原始錯誤）",
			message:        "Password hash failed",
			code:           CodePasswordHashFail,
			err:            nil,
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   CodePasswordHashFail,
			expectedMsg:    "Password hash failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appError := NewInternalServerError(tt.message, tt.code, tt.err)

			assert.NotNil(t, appError)
			assert.Equal(t, tt.expectedStatus, appError.StatusCode)
			assert.Equal(t, tt.expectedCode, appError.Code)
			assert.Equal(t, tt.expectedMsg, appError.Message)
			assert.Equal(t, tt.err, appError.Err)
		})
	}
}

// TestNewForbiddenError 測試建立 403 錯誤
func TestNewForbiddenError(t *testing.T) {
	tests := []struct {
		name           string
		message        string
		code           int
		err            error
		expectedStatus int
		expectedCode   int
		expectedMsg    string
	}{
		{
			name:           "建立 403 錯誤（有原始錯誤）",
			message:        "Access forbidden",
			code:           CodeForbidden,
			err:            errors.New("insufficient permissions"),
			expectedStatus: http.StatusForbidden,
			expectedCode:   CodeForbidden,
			expectedMsg:    "Access forbidden",
		},
		{
			name:           "建立 403 錯誤（無原始錯誤）",
			message:        "Permission denied",
			code:           CodeForbidden,
			err:            nil,
			expectedStatus: http.StatusForbidden,
			expectedCode:   CodeForbidden,
			expectedMsg:    "Permission denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appError := NewForbiddenError(tt.message, tt.code, tt.err)

			assert.NotNil(t, appError)
			assert.Equal(t, tt.expectedStatus, appError.StatusCode)
			assert.Equal(t, tt.expectedCode, appError.Code)
			assert.Equal(t, tt.expectedMsg, appError.Message)
			assert.Equal(t, tt.err, appError.Err)
		})
	}
}

// TestNewUnauthorizedError 測試建立 401 錯誤
func TestNewUnauthorizedError(t *testing.T) {
	tests := []struct {
		name           string
		message        string
		code           int
		err            error
		expectedStatus int
		expectedCode   int
		expectedMsg    string
	}{
		{
			name:           "建立 401 錯誤（有原始錯誤）",
			message:        "Unauthorized",
			code:           CodeTokenInvalid,
			err:            errors.New("token expired"),
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   CodeTokenInvalid,
			expectedMsg:    "Unauthorized",
		},
		{
			name:           "建立 401 錯誤（無原始錯誤）",
			message:        "Token missing",
			code:           CodeTokenMissing,
			err:            nil,
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   CodeTokenMissing,
			expectedMsg:    "Token missing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appError := NewUnauthorizedError(tt.message, tt.code, tt.err)

			assert.NotNil(t, appError)
			assert.Equal(t, tt.expectedStatus, appError.StatusCode)
			assert.Equal(t, tt.expectedCode, appError.Code)
			assert.Equal(t, tt.expectedMsg, appError.Message)
			assert.Equal(t, tt.err, appError.Err)
		})
	}
}

// TestNewConflictError 測試建立 409 錯誤
func TestNewConflictError(t *testing.T) {
	tests := []struct {
		name           string
		message        string
		code           int
		err            error
		expectedStatus int
		expectedCode   int
		expectedMsg    string
	}{
		{
			name:           "建立 409 錯誤（有原始錯誤）",
			message:        "Email already exists",
			code:           CodeEmailExists,
			err:            errors.New("duplicate email"),
			expectedStatus: http.StatusConflict,
			expectedCode:   CodeEmailExists,
			expectedMsg:    "Email already exists",
		},
		{
			name:           "建立 409 錯誤（無原始錯誤）",
			message:        "Email in use",
			code:           CodeEmailInUse,
			err:            nil,
			expectedStatus: http.StatusConflict,
			expectedCode:   CodeEmailInUse,
			expectedMsg:    "Email in use",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appError := NewConflictError(tt.message, tt.code, tt.err)

			assert.NotNil(t, appError)
			assert.Equal(t, tt.expectedStatus, appError.StatusCode)
			assert.Equal(t, tt.expectedCode, appError.Code)
			assert.Equal(t, tt.expectedMsg, appError.Message)
			assert.Equal(t, tt.err, appError.Err)
		})
	}
}

// TestNewUnsupportedMediaTypeError 測試建立 415 錯誤
func TestNewUnsupportedMediaTypeError(t *testing.T) {
	tests := []struct {
		name           string
		message        string
		code           int
		err            error
		expectedStatus int
		expectedCode   int
		expectedMsg    string
	}{
		{
			name:           "建立 415 錯誤（有原始錯誤）",
			message:        "Unsupported media type",
			code:           CodeUnsupportedMediaType,
			err:            errors.New("invalid content-type"),
			expectedStatus: http.StatusUnsupportedMediaType,
			expectedCode:   CodeUnsupportedMediaType,
			expectedMsg:    "Unsupported media type",
		},
		{
			name:           "建立 415 錯誤（無原始錯誤）",
			message:        "Expected application/json",
			code:           CodeUnsupportedMediaType,
			err:            nil,
			expectedStatus: http.StatusUnsupportedMediaType,
			expectedCode:   CodeUnsupportedMediaType,
			expectedMsg:    "Expected application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appError := NewUnsupportedMediaTypeError(tt.message, tt.code, tt.err)

			assert.NotNil(t, appError)
			assert.Equal(t, tt.expectedStatus, appError.StatusCode)
			assert.Equal(t, tt.expectedCode, appError.Code)
			assert.Equal(t, tt.expectedMsg, appError.Message)
			assert.Equal(t, tt.err, appError.Err)
		})
	}
}
