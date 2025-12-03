package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"restfulapi/internal/config"
	db "restfulapi/internal/database"
	"restfulapi/internal/middleware"
	"restfulapi/internal/model"
	"restfulapi/internal/util"
	"restfulapi/pkg/logger"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// RegisterUserInput 定義註冊使用者輸入結構
type RegisterUserInput struct {
	Name     string `json:"name" example:"Victor" binding:"required"`                    // 姓名
	Email    string `json:"email" example:"victor@example.com" binding:"required,email"` // Email
	Password string `json:"password" example:"123456" binding:"required"`                // 密碼
}

// RegisterUserResponse 定義註冊成功回傳的資料結構
type RegisterUserResponse struct {
	ID    uint   `json:"id" example:"18"`                    // 使用者 ID
	Name  string `json:"name" example:"Victor"`              // 姓名
	Email string `json:"email" example:"victor@example.com"` // Email
}

// LoginUserInput 定義登入使用者輸入結構
type LoginUserInput struct {
	Email    string `json:"email" example:"victor@example.com" binding:"required,email"` // Email
	Password string `json:"password" example:"123456" binding:"required"`                // 密碼
}

// LoginUserResponse 定義登入成功回傳的資料結構
type LoginUserResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`  // Access Token
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // Refresh Token
}

// LogoutUserResponse 定義登出成功回傳的資料結構（空對象）
type LogoutUserResponse struct {
}

// LogoutSuccessResponse 定義登出成功回應結構
type LogoutSuccessResponse struct {
	// Status 狀態: "success"
	Status string `json:"status" example:"success"`

	// Message 訊息描述
	Message string `json:"message" example:"Logout successful"`

	// Data 主要的回應數據（空對象）
	Data LogoutUserResponse `json:"data"`
}

// RefreshTokenInput 定義刷新 Token 輸入結構
type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." binding:"required"` // Refresh Token
}

// LogoutInput 定義登出輸入結構（可選的 Refresh Token）
type LogoutInput struct {
	RefreshToken string `json:"refresh_token,omitempty" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // Refresh Token（可選，如果有就註銷）
}

// RegisterUser
// @Summary 會員註冊
// @Description 註冊新使用者
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   user  body RegisterUserInput  true  "使用者資訊"
// @Success 201 {object} util.SuccessAPIResponse{data=RegisterUserResponse} "註冊成功，回傳使用者資訊"
// @Failure 400 {object} util.ErrorAPIResponseInvalidInput "無效的輸入格式，error_code: 1001"
// @Failure 400 {object} util.ErrorAPIResponseMissingFields "缺少必填欄位，error_code: 1002"
// @Failure 409 {object} util.ErrorAPIResponseEmailExists "Email 已經被註冊，error_code: 1003"
// @Failure 415 {object} util.ErrorAPIResponseUnsupportedMediaType "不支援的媒體類型，error_code: 1000"
// @Failure 500 {object} util.ErrorAPIResponseInternalServerError "伺服器內部錯誤，error_code: 4001"
// @Router /v1/auth/register [post]
func RegisterUser(c *gin.Context) {
	if !validateContentType(c) {
		return
	}

	input := RegisterUserInput{}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		util.ValidationErrorResponse(c, err)
		return
	}

	user := model.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		Role:     "user", // 強制設定 role 為 "user"，防止惡意註冊成 admin
	}

	createUserId, createUserErr := model.CreateUser(user)

	if createUserErr != nil {
		handleCreateUserError(c, createUserErr, "register")
		return
	}

	user.ID = createUserId

	// 定義 RabbitMQ 訊息結構（不包含密碼）
	type UserCreatedMessage struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	}

	// 轉換成不包含密碼的格式
	messageData := UserCreatedMessage{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	// 發送 RabbitMQ 事件
	message, _ := json.Marshal(messageData)
	config.PublishMessage("user_created", string(message))

	log.Println("註冊成功，已發送 user_created 事件:", string(message))

	// 回應標準 JSON
	util.CreatedResponse(c, "User created successfully", gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

// LoginUser
// @Summary 使用者登入
// @Description 使用者登入並取得 Access Token 和 Refresh Token
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   credentials  body LoginUserInput  true  "登入資訊"
// @Success 200 {object} util.SuccessAPIResponse{data=LoginUserResponse} "登入成功，回傳 Access Token 和 Refresh Token"
// @Failure 400 {object} util.ErrorAPIResponseInvalidInput "無效的輸入格式，error_code: 1001"
// @Failure 401 {object} util.ErrorAPIResponseInvalidCredentials "帳號或密碼錯誤，error_code: 1004"
// @Failure 415 {object} util.ErrorAPIResponseUnsupportedMediaType "不支援的媒體類型，error_code: 1000"
// @Failure 500 {object} util.ErrorAPIResponseInternalServerError "伺服器內部錯誤，error_code: 4001"
// @Router /v1/auth/login [post]
func LoginUser(c *gin.Context) {
	if !validateContentType(c) {
		return
	}

	var input LoginUserInput

	// 解析請求 JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(util.NewBadRequestError("Invalid input", util.CodeInvalidInput, err))
		return
	}

	// 查詢使用者
	var user model.User
	result := db.DbConnect.Where("email = ?", input.Email).First(&user)

	// 檢查使用者是否存在
	if result.RowsAffected == 0 {
		c.Error(util.NewUnauthorizedError("Invalid email or password", util.CodeInvalidCredentials, nil))
		return
	}

	// 檢查資料庫查詢錯誤（使用者存在但查詢過程發生其他錯誤）
	if result.Error != nil {
		c.Error(util.NewInternalServerError(
			"Internal server error",
			util.CodeDatabaseError,
			fmt.Errorf("failed to login: email=%s, error=%w", input.Email, result.Error),
		))
		return
	}

	// 檢查密碼
	if !user.CheckPassword(input.Password) {
		c.Error(util.NewUnauthorizedError("Invalid email or password", util.CodeInvalidCredentials, nil))
		return
	}

	// 產生 Access Token & Refresh Token
	accessToken, refreshToken, err := config.GenerateTokens(user.ID, user.Email, user.Role)
	if err != nil {
		c.Error(util.NewInternalServerError(
			"Internal server error",
			util.CodeTokenGenerationFailed,
			fmt.Errorf("failed to generate token: userId=%d, error=%w", user.ID, err),
		))
		return
	}

	// 回應標準 JSON
	util.SuccessResponse(c, "Login successful!", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// RefreshToken
// @Summary 刷新 Access Token
// @Description 使用 Refresh Token 來刷新 Access Token 和 Refresh Token
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   refresh_token  body RefreshTokenInput  true  "Refresh Token"
// @Success 200 {object} util.SuccessAPIResponse{data=LoginUserResponse} "刷新成功，回傳新的 Access Token 和 Refresh Token"
// @Failure 400 {object} util.ErrorAPIResponseInvalidInput "無效的輸入格式，error_code: 1001"
// @Failure 401 {object} util.ErrorAPIResponseRefreshTokenInvalid "Refresh Token 無效或過期，error_code: 1008"
// @Failure 415 {object} util.ErrorAPIResponseUnsupportedMediaType "不支援的媒體類型，error_code: 1000"
// @Failure 500 {object} util.ErrorAPIResponseInternalServerError "伺服器內部錯誤，error_code: 4001"
// @Router /v1/auth/refresh [post]
func RefreshToken(c *gin.Context) {
	if !validateContentType(c) {
		return
	}

	var input RefreshTokenInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(util.NewBadRequestError("Invalid input", util.CodeInvalidInput, err))
		return
	}

	// 檢查 Redis 是否已經黑名單（登出或rotation 放入的）
	ctx := context.Background()
	refreshTokenBlacklistKey := "refresh_token:blacklist:" + input.RefreshToken
	exists, err := config.RedisClient.Exists(ctx, refreshTokenBlacklistKey).Result()
	if err == nil && exists > 0 {
		c.Error(util.NewUnauthorizedError("Refresh token has been revoked", util.CodeRefreshTokenInvalid, nil))
		return
	}

	token, err := config.ValidateToken(input.RefreshToken, config.JWTConfig.RefreshTokenSecret)
	if err != nil || !token.Valid {
		c.Error(util.NewUnauthorizedError("Invalid or expired refresh token", util.CodeRefreshTokenInvalid, err))
		return
	}

	// 提取 exp，進行 Rotation，本次 refresh token 立即失效（寫入黑名單）
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.Error(util.NewInternalServerError(
			"Internal server error",
			util.CodeInternalError,
			errors.New("failed to parse token claims"),
		))
		return
	}

	refreshExp, ok := claims["exp"].(float64)
	if !ok {
		c.Error(util.NewInternalServerError(
			"Internal server error",
			util.CodeInternalError,
			errors.New("no exp in refresh token"),
		))
		return
	}
	refreshExpTime := time.Unix(int64(refreshExp), 0)
	refreshRemainingTime := time.Until(refreshExpTime)
	if refreshRemainingTime > 0 {
		// Rotation: 用過即失效
		err := config.RedisClient.Set(ctx, refreshTokenBlacklistKey, "1", refreshRemainingTime).Err()
		if err != nil {
			traceID := middleware.GetTraceID(c)
			logger.LogError("refresh-rotation", traceID, err)
		}
	}

	// 從 Token 中提取使用者資訊
	userID, ok := claims["sub"].(float64)
	if !ok {
		c.Error(util.NewInternalServerError(
			"Internal server error",
			util.CodeInternalError,
			errors.New("invalid user ID in token"),
		))
		return
	}

	email, ok := claims["email"].(string)
	if !ok {
		c.Error(util.NewInternalServerError(
			"Internal server error",
			util.CodeInternalError,
			errors.New("invalid email in token"),
		))
		return
	}

	// 查詢使用者以獲取 role
	var user model.User
	result := db.DbConnect.Where("id = ?", uint(userID)).First(&user)
	if result.RowsAffected == 0 {
		c.Error(util.NewUnauthorizedError("Invalid or expired refresh token", util.CodeRefreshTokenInvalid, errors.New("user not found")))
		return
	}

	if result.Error != nil {
		c.Error(util.NewInternalServerError(
			"Internal server error",
			util.CodeDatabaseError,
			fmt.Errorf("failed to refresh token: userId=%d, error=%w", uint(userID), result.Error),
		))
		return
	}

	// 產生新的 Access Token & Refresh Token
	accessToken, refreshToken, err := config.GenerateTokens(user.ID, email, user.Role)
	if err != nil {
		c.Error(util.NewInternalServerError(
			"Internal server error",
			util.CodeTokenGenerationFailed,
			fmt.Errorf("failed to generate token: userId=%d, error=%w", user.ID, err),
		))
		return
	}

	// 回應標準 JSON
	util.SuccessResponse(c, "Token refreshed successfully", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// LogoutUser
// @Summary 使用者登出
// @Description 使用者登出，將 Access Token 加入黑名單使其立即失效。如果提供了 Refresh Token，也會一併註銷。
// @Tags Auth
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param   refresh_token  body LogoutInput false "Refresh Token（可選，如果有就一併註銷）"
// @Success 200 {object} LogoutSuccessResponse "登出成功"
// @Failure 400 {object} util.ErrorAPIResponseInvalidInput "無效的輸入格式，error_code: 1001"
// @Failure 401 {object} util.ErrorAPIResponseTokenMissing "Token 缺失，error_code: 1006"
// @Failure 401 {object} util.ErrorAPIResponseTokenInvalid "Token 無效或過期，error_code: 1007"
// @Router /v1/auth/logout [post]
func LogoutUser(c *gin.Context) {
	// 從 Header 中提取 Access Token（可能缺/過期/無效，仍然繼續流程）
	tokenString := c.GetHeader("Authorization")
	if tokenString != "" {
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// 嘗試驗證 Token 並解析 Claims 以取得過期時間（即使無效也不提前 return，僅能做黑名單的做黑名單）
		token, err := config.ValidateToken(tokenString, config.JWTConfig.AccessTokenSecret)
		if err == nil && token.Valid {
			ctx := context.Background()
			claims, ok := token.Claims.(jwt.MapClaims)
			if ok {
				exp, ok := claims["exp"].(float64)
				if ok {
					expTime := time.Unix(int64(exp), 0)
					remainingTime := time.Until(expTime)
					if remainingTime > 0 {
						blacklistKey := "access_token:blacklist:" + tokenString
						err := config.RedisClient.Set(ctx, blacklistKey, "1", remainingTime).Err()
						if err != nil {
							traceID := middleware.GetTraceID(c)
							logger.LogError("logout", traceID, err)
						}
					}
				}
			}
		}
	}

	// 不論 access token 狀態如何，皆可處理 Refresh Token
	var input LogoutInput
	if err := c.ShouldBindJSON(&input); err == nil && input.RefreshToken != "" {
		refreshToken, err := config.ValidateToken(input.RefreshToken, config.JWTConfig.RefreshTokenSecret)
		if err == nil && refreshToken.Valid {
			ctx := context.Background()
			refreshClaims, ok := refreshToken.Claims.(jwt.MapClaims)
			if ok {
				refreshExp, ok := refreshClaims["exp"].(float64)
				if ok {
					refreshExpTime := time.Unix(int64(refreshExp), 0)
					refreshRemainingTime := time.Until(refreshExpTime)
					if refreshRemainingTime > 0 {
						refreshBlacklistKey := "refresh_token:blacklist:" + input.RefreshToken
						err := config.RedisClient.Set(ctx, refreshBlacklistKey, "1", refreshRemainingTime).Err()
						if err != nil {
							traceID := middleware.GetTraceID(c)
							logger.LogError("logout", traceID, err)
						}
					}
				}
			}
		}
		// 即使 Refresh Token 無效，也不影響登出流程（因為可能過期或被清除了）
	}

	util.SuccessResponse(c, "Logout successful", gin.H{})
}

// handleCreateUserError 處理創建用戶時的錯誤
func handleCreateUserError(c *gin.Context, err error, operation string) {
	switch {
	case errors.Is(err, util.ErrEmailExists):
		c.Error(util.NewConflictError("Email already registered", util.CodeEmailExists, err))
	case errors.Is(err, util.ErrPasswordHashFail):
		c.Error(util.NewInternalServerError(
			"Internal server error",
			util.CodePasswordHashFail,
			fmt.Errorf("failed to create user: operation=%s, error=%w", operation, err),
		))
	case errors.Is(err, util.ErrDatabaseError):
		c.Error(util.NewInternalServerError(
			"Internal server error",
			util.CodeDatabaseError,
			fmt.Errorf("failed to create user: operation=%s, error=%w", operation, err),
		))
	default:
		c.Error(util.NewInternalServerError(
			"Internal server error",
			util.CodeInternalError,
			fmt.Errorf("failed to create user: operation=%s, error=%w", operation, err),
		))
	}
}

func validateContentType(c *gin.Context) bool {
	contentType := c.GetHeader("Content-Type")
	if contentType == "" || strings.HasPrefix(contentType, "application/json") {
		return true
	}

	c.Error(util.NewUnsupportedMediaTypeError("Unsupported media type. Expected application/json", util.CodeUnsupportedMediaType, nil))
	return false
}
