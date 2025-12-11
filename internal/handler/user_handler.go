package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"restfulapi/internal/config"
	"restfulapi/internal/model"
	"restfulapi/internal/util"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// UserProfileData 定義使用者資料回應的 data 結構（Swagger 用）
// @Description 使用者資料結構
// swagger:model UserProfileData
type UserProfileData struct {
	ID    uint64 `json:"id" example:"1"`
	Name  string `json:"name" example:"Victor"`
	Email string `json:"email" example:"victor@example.com"`
	Role  string `json:"role" example:"user"`
}

// UserProfileSuccessResponse 定義使用者資料成功回應（Swagger 用）
// swagger:model UserProfileSuccessResponse
type UserProfileSuccessResponse struct {
	Status  string          `json:"status" example:"success"`
	Message string          `json:"message" example:"User retrieved successfully"`
	Data    UserProfileData `json:"data"`
}

// GetMyProfile - 查詢自己的資料
// @Summary 查詢自己的資料
// @Description 讓使用者可以查詢自己的資料（從 JWT Token 取得 userID）
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} handler.UserProfileSuccessResponse "成功取得使用者資料"
// @Failure 401 {object} util.ErrorAPIResponse "Unauthorized，error_code: 1006（Token 缺失）或 1007（Token 無效）"
// @Failure 404 {object} util.ErrorAPIResponse "Not Found，error_code: 1005（使用者不存在）"
// @Failure 500 {object} util.ErrorAPIResponse "Internal Server Error，error_code: 4001（伺服器內部錯誤）"
// @Router /v1/users/me [get]
func GetMyProfile(c *gin.Context) {
	// 從 JWT Token 取得 user_id（由 RequireAuth middleware 設置）
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.Error(util.NewUnauthorizedError("Unauthorized", util.CodeTokenMissing, nil))
		return
	}

	// 將 user_id 轉換為 uint64
	// JWT claims 中的數字可能是 float64 類型（JSON 數字解析）
	var userID uint64
	switch v := userIDInterface.(type) {
	case float64:
		userID = uint64(v)
	case uint64:
		userID = v
	case string:
		parsedID, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			c.Error(util.NewUnauthorizedError("Invalid user ID in token", util.CodeTokenInvalid, err))
			return
		}
		userID = parsedID
	default:
		c.Error(util.NewUnauthorizedError("Invalid user ID type in token", util.CodeTokenInvalid, nil))
		return
	}

	ctx := context.Background()
	cacheKey := "user:" + strconv.FormatUint(userID, 10)

	// 先查詢 Redis 快取
	cachedUser, err := config.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var user model.User
		if json.Unmarshal([]byte(cachedUser), &user) == nil {
			util.SuccessResponse(c, "User retrieved from cache", gin.H{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
				"role":  user.Role,
			})
			return
		}
	}

	// 查詢資料庫
	user, err := model.FindByUserId(userID)
	// 使用 `errors.Is()` 來區分不同的錯誤
	if errors.Is(err, util.ErrUserNotFound) {
		c.Error(util.NewNotFoundError("User not found", util.CodeUserNotFound, err))
		return
	} else if err != nil {
		c.Error(util.NewInternalServerError(
			"Internal server error",
			util.CodeDatabaseError,
			fmt.Errorf("failed to find user by id: userId=%d, error=%w", userID, err),
		))
		return
	}

	// 存入 Redis，快取 5 分鐘
	userJSON, _ := json.Marshal(user)
	config.RedisClient.Set(ctx, cacheKey, userJSON, time.Minute*5)

	util.SuccessResponse(c, "User retrieved successfully", gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	})
}
