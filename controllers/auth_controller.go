package controllers

import (
	"errors"
	"net/http"
	"restfulapi/config"
	db "restfulapi/database"
	"restfulapi/models"
	"restfulapi/utils"

	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
	user := models.User{}
	err := c.BindJSON(&user)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", utils.ErrInvalidInput)
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
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user", utils.ErrPasswordHashFail)
		case errors.Is(createUserErr, models.ErrDatabaseError):
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user", utils.ErrDatabaseError)
		default:
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user", utils.ErrInternalError)
		}
		return
	}

	user.ID = createUserId

	// **回應標準 JSON**
	utils.SuccessResponse(c, "User created successfully", gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

func LoginUser(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// 解析請求 JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", utils.ErrInvalidInput)
		return
	}

	// 查詢使用者
	var user models.User
	db.DBconnect.Where("email = ?", input.Email).First(&user)

	// 檢查密碼
	if !user.CheckPassword(input.Password) {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password", utils.ErrInvalidCredentials)
		return
	}

	// 產生 Access Token & Refresh Token
	accessToken, refreshToken, err := config.GenerateTokens(user.ID, user.Email, user.Role)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token", utils.ErrTokenGenerationFailed)
		return
	}

	// 回應標準 JSON
	utils.SuccessResponse(c, "Login successful!", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// LogoutUser - 使用者登出（JWT 無狀態登出）
func LogoutUser(c *gin.Context) {
	// 只需要回應成功，前端會移除 Token
	utils.SuccessResponse(c, "Logout successful", gin.H{})
}
