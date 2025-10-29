package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"restfulapi/config"
	db "restfulapi/database"
	"restfulapi/models"
	"restfulapi/pkg/logger"
	"restfulapi/utils"
	"strings"

	"github.com/gin-gonic/gin"
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

	// Meta Meta 資訊 (例如 API 版本、請求 ID 等)
	Meta utils.MetaData `json:"meta"`
}

// RegisterUser
// @Summary 會員註冊
// @Description 註冊新使用者
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   user  body RegisterUserInput  true  "使用者資訊"
// @Success 201 {object} utils.SuccessAPIResponse{data=RegisterUserResponse} "註冊成功，回傳使用者資訊"
// @Failure 400 {object} utils.ErrorAPIResponseInvalidInput "無效的輸入格式，error_code: 1001"
// @Failure 400 {object} utils.ErrorAPIResponseMissingFields "缺少必填欄位，error_code: 1002"
// @Failure 409 {object} utils.ErrorAPIResponseEmailExists "Email 已經被註冊，error_code: 1003"
// @Failure 415 {object} utils.ErrorAPIResponseUnsupportedMediaType "不支援的媒體類型，error_code: 1000"
// @Failure 500 {object} utils.ErrorAPIResponseInternalServerError "伺服器內部錯誤，error_code: 4001"
// @Router /auth/register [post]
func RegisterUser(c *gin.Context) {
	// 檢查 Content-Type 是否為 application/json（支援帶參數的格式，如 application/json; charset=utf-8）
	contentType := c.GetHeader("Content-Type")
	if contentType != "" && !strings.HasPrefix(contentType, "application/json") {
		utils.ErrorResponse(c, http.StatusUnsupportedMediaType, "Unsupported media type. Expected application/json", utils.ErrUnsupportedMediaType)
		return
	}

	user := models.User{}
	err := c.BindJSON(&user)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", utils.ErrInvalidInput)
		return
	}

	// 驗證必填欄位
	if user.Name == "" || user.Email == "" || user.Password == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Missing required fields", utils.ErrMissingRequiredFields)
		return
	}

	// 強制設定 role 為 "user"，防止惡意註冊成 admin
	user.Role = "user"

	createUserId, createUserErr := models.CreateUser(user)

	if createUserErr != nil {
		switch {
		case errors.Is(createUserErr, models.ErrEmailExists):
			utils.ErrorResponse(c, http.StatusConflict, "Email already registered", utils.ErrEmailExists)
		case errors.Is(createUserErr, models.ErrPasswordHashFail):
			logger.LogError("register", createUserErr)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user", utils.ErrPasswordHashFail)
		case errors.Is(createUserErr, models.ErrDatabaseError):
			logger.LogError("register", createUserErr)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user", utils.ErrDatabaseError)
		default:
			logger.LogError("register", createUserErr)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user", utils.ErrInternalError)
		}
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
	utils.CreatedResponse(c, "User created successfully", gin.H{
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
// @Success 200 {object} utils.SuccessAPIResponse{data=LoginUserResponse} "登入成功，回傳 Access Token 和 Refresh Token"
// @Failure 400 {object} utils.ErrorAPIResponseInvalidInput "無效的輸入格式，error_code: 1001"
// @Failure 401 {object} utils.ErrorAPIResponseInvalidCredentials "帳號或密碼錯誤，error_code: 1004"
// @Failure 415 {object} utils.ErrorAPIResponseUnsupportedMediaType "不支援的媒體類型，error_code: 1000"
// @Failure 500 {object} utils.ErrorAPIResponseInternalServerError "伺服器內部錯誤，error_code: 4001"
// @Router /auth/login [post]
func LoginUser(c *gin.Context) {
	// 檢查 Content-Type 是否為 application/json（支援帶參數的格式，如 application/json; charset=utf-8）
	contentType := c.GetHeader("Content-Type")
	if contentType != "" && !strings.HasPrefix(contentType, "application/json") {
		utils.ErrorResponse(c, http.StatusUnsupportedMediaType, "Unsupported media type. Expected application/json", utils.ErrUnsupportedMediaType)
		return
	}

	var input LoginUserInput

	// 解析請求 JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", utils.ErrInvalidInput)
		return
	}

	// 查詢使用者
	var user models.User
	result := db.DbConnect.Where("email = ?", input.Email).First(&user)

	// 檢查使用者是否存在
	if result.RowsAffected == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password", utils.ErrInvalidCredentials)
		return
	}

	// 檢查資料庫查詢錯誤（使用者存在但查詢過程發生其他錯誤）
	if result.Error != nil {
		logger.LogError("login", result.Error)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to login", utils.ErrDatabaseError)
		return
	}

	// 檢查密碼
	if !user.CheckPassword(input.Password) {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password", utils.ErrInvalidCredentials)
		return
	}

	// 產生 Access Token & Refresh Token
	accessToken, refreshToken, err := config.GenerateTokens(user.ID, user.Email, user.Role)
	if err != nil {
		logger.LogError("login", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token", utils.ErrTokenGenerationFailed)
		return
	}

	// 回應標準 JSON
	utils.SuccessResponse(c, "Login successful!", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// LogoutUser
// @Summary 使用者登出
// @Description 使用者登出，前端需移除本地儲存的 Token（JWT 無狀態登出）
// @Tags Auth
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} LogoutSuccessResponse "登出成功"
// @Failure 401 {object} utils.ErrorAPIResponseTokenMissing "Token 缺失，error_code: 1006"
// @Failure 401 {object} utils.ErrorAPIResponseTokenInvalid "Token 無效或過期，error_code: 1007"
// @Router /auth/logout [post]
func LogoutUser(c *gin.Context) {
	// TODO: 實作 Token 黑名單機制
	// JWT 無狀態登出：目前只需回應成功，前端會移除本地儲存的 Token
	// 注意：目前 Token 在過期前仍有效，實作黑名單後可立即失效
	utils.SuccessResponse(c, "Logout successful", gin.H{})
}
