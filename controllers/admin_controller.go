package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"restfulapi/config"
	"restfulapi/models"
	repository "restfulapi/repositories"
	service "restfulapi/services"
	"restfulapi/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// FindAllUsers - 取得所有使用者
// @Summary 取得使用者列表
// @Description 取得所有使用者資料（需管理員權限）
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.SuccessAPIResponse "成功取得使用者列表"
// @Failure 401 {object} utils.ErrorAPIResponseTokenMissing "Unauthorized，error_code: 1006"
// @Router /v1/admin/users [get]
func FindAllUsers(c *gin.Context) {
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	users := userService.GetAllUsers()

	// 以 DTO struct 確保欄位順序
	type AdminUserListItem struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	items := make([]AdminUserListItem, 0, len(users))
	for _, u := range users {
		items = append(items, AdminUserListItem{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
			Role:  u.Role,
		})
	}

	// 用包裝 struct，而非 map，避免 key 順序隨機
	type AdminUserListResponse struct {
		Items []AdminUserListItem `json:"items"`
	}

	utils.SuccessResponse(c, "All users retrieved successfully", AdminUserListResponse{Items: items})
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
