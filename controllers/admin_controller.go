package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"restfulapi/config"
	"restfulapi/models"
	"restfulapi/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// FindAllUsers - 取得所有使用者
func FindAllUsers(c *gin.Context) {
	users := models.FindAllUsers()
	utils.SuccessResponse(c, "All users retrieved successfully", gin.H{"items": users})
}

// FindByUserId - 取得單一使用者
func FindByUserId(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", 1001)
		return
	}

	ctx := context.Background()
	cacheKey := "user:" + strconv.Itoa(userId)

	// 1 先查詢 Redis
	cachedUser, err := config.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var user models.User
		if json.Unmarshal([]byte(cachedUser), &user) == nil {
			utils.SuccessResponse(c, "User retrieved from cache", gin.H{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
				"role":  user.Role,
			})
			return
		}
	}

	user, err := models.FindByUserId(userId)
	// 使用 `errors.Is()` 來區分不同的錯誤
	if errors.Is(err, models.ErrUserNotFound) {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found", 1009)
		return
	} else if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Database error", 5000)
		return
	}

	// 存入 Redis，快取 5 分鐘
	userJSON, _ := json.Marshal(user)
	config.RedisClient.Set(ctx, cacheKey, userJSON, time.Minute*5)

	utils.SuccessResponse(c, "User retrieved successfully", gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	})
}

// DeleteUser - 刪除使用者
func DeleteUser(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", 1001)
		return
	}

	// 嘗試刪除使用者
	rowsAffected, deleteErr := models.DeleteUser(userId)
	if deleteErr != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete user", 1011)
		return
	}

	// 如果沒有影響任何行，代表用戶不存在
	if rowsAffected == 0 {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found", 1009)
		return
	}

	// 刪除成功
	utils.SuccessResponse(c, "User deleted successfully", nil)
}

// UpdateUser - 更新使用者資訊
func UpdateUser(c *gin.Context) {
	// 解析 `id`
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", 1001)
		return
	}

	var input models.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", 1000)
		return
	}

	user, err := models.FindByUserId(userId)
	if errors.Is(err, models.ErrUserNotFound) {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found", 1009)
		return
	} else if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Database error", 5000)
		return
	}

	// ✅ 檢查 `email` 是否已經被其他用戶使用
	if input.Email != "" && input.Email != user.Email {
		exists := models.IsEmailExists(input.Email, userId) // ✅ 確保不包含自己
		if exists {
			utils.ErrorResponse(c, http.StatusConflict, "Email already in use", 1013)
			return
		}
	}

	// 更新會員
	result := models.UpdateUser(user, input)
	if result == 1 {
		updatedUser, _ := models.FindByUserId(userId) // 查詢最新的資料

		utils.SuccessResponse(c, "User updated successfully", gin.H{"user": updatedUser})
		return
	}

	utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update user", 1012)
}
