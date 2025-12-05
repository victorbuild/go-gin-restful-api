package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"restfulapi/internal/config"
	"restfulapi/internal/model"
	repository "restfulapi/internal/repository"
	service "restfulapi/internal/service"
	"restfulapi/internal/util"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// AdminUserItem 定義使用者列表的單筆資料（Swagger 用）
// @Description 使用者列表的單筆資料
// swagger:model AdminUserItem
type AdminUserItem struct {
	ID    uint   `json:"id" example:"1"`
	Name  string `json:"name" example:"Victor"`
	Email string `json:"email" example:"victor@example.com"`
	Role  string `json:"role" example:"user"`
}

// AdminUserListData 定義列表回應的 data 結構（Swagger 用）
// swagger:model AdminUserListData
type AdminUserListData struct {
	Items []AdminUserItem `json:"items"`
}

// AdminUserListSuccessResponse 定義成功回應（Swagger 用）
// swagger:model AdminUserListSuccessResponse
type AdminUserListSuccessResponse struct {
	Status  string            `json:"status" example:"success"`
	Message string            `json:"message" example:"All users retrieved successfully"`
	Data    AdminUserListData `json:"data"`
}

// FindAllUsers - 取得所有使用者
// @Summary 取得使用者列表
// @Description 取得所有使用者資料（需管理員權限）
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} handler.AdminUserListSuccessResponse "成功取得使用者列表"
// @Failure 401 {object} util.ErrorAPIResponseTokenMissing "Unauthorized，error_code: 1006"
// @Failure 403 {object} util.ErrorAPIResponseForbidden "Forbidden，error_code: 1010"
// @Failure 500 {object} util.ErrorAPIResponseInternalServerError "Internal server error，error_code: 4001"
// @Router /v1/admin/users [get]
func FindAllUsers(c *gin.Context) {
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	users := userService.GetAllUsers()

	// 以 DTO struct 確保欄位順序
	type AdminUserListItem struct {
		ID    uint64 `json:"id"`
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

	util.SuccessResponse(c, "All users retrieved successfully", AdminUserListResponse{Items: items})
}

// FindByUserId - 取得單一使用者
func FindByUserId(c *gin.Context) {
	userId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.NewBadRequestError("Invalid user ID", util.CodeInvalidInput, err))
		return
	}

	ctx := context.Background()
	cacheKey := "user:" + strconv.FormatUint(userId, 10)

	// 1 先查詢 Redis
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

	user, err := model.FindByUserId(userId)
	// 使用 `errors.Is()` 來區分不同的錯誤
	if errors.Is(err, util.ErrUserNotFound) {
		c.Error(util.NewNotFoundError("User not found", util.CodeUserNotFound, err))
		return
	} else if err != nil {
		c.Error(util.NewInternalServerError(
			"Internal server error",
			util.CodeDatabaseError,
			fmt.Errorf("failed to find user by id: userId=%d, error=%w", userId, err),
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

// DeleteUser - 刪除使用者
func DeleteUser(c *gin.Context) {
	userId, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.Error(util.NewBadRequestError("Invalid user ID", util.CodeInvalidInput, err))
		return
	}

	// 嘗試刪除使用者
	rowsAffected, deleteErr := model.DeleteUser(userId)
	if deleteErr != nil {
		c.Error(util.NewInternalServerError(
			"Internal server error",
			util.CodeDeleteUserFailed,
			fmt.Errorf("failed to delete user: userId=%d, error=%w", userId, deleteErr),
		))
		return
	}

	// 如果沒有影響任何行，代表用戶不存在
	if rowsAffected == 0 {
		c.Error(util.NewNotFoundError("User not found", util.CodeUserNotFound, nil))
		return
	}

	// 刪除成功
	util.SuccessResponse(c, "User deleted successfully", nil)
}

// UpdateUser - 更新使用者資訊
func UpdateUser(c *gin.Context) {
	// 解析 `id`
	userId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.NewBadRequestError("Invalid user ID", util.CodeInvalidInput, err))
		return
	}

	var input model.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(util.NewBadRequestError("Invalid input", util.CodeInvalidInput, err))
		return
	}

	user, err := model.FindByUserId(userId)
	if errors.Is(err, util.ErrUserNotFound) {
		c.Error(util.NewNotFoundError("User not found", util.CodeUserNotFound, err))
		return
	} else if err != nil {
		c.Error(util.NewInternalServerError(
			"Internal server error",
			util.CodeDatabaseError,
			fmt.Errorf("failed to find user by id: userId=%d, error=%w", userId, err),
		))
		return
	}

	// 檢查 email 是否已經被其他用戶使用
	if input.Email != "" && input.Email != user.Email {
		exists := model.IsEmailExists(input.Email, userId)
		if exists {
			c.Error(util.NewConflictError("Email already in use", util.CodeEmailInUse, nil))
			return
		}
	}

	// 更新會員
	result := model.UpdateUser(user, input)
	if result == 1 {
		updatedUser, _ := model.FindByUserId(userId) // 查詢最新的資料

		util.SuccessResponse(c, "User updated successfully", gin.H{"user": updatedUser})
		return
	}

	c.Error(util.NewInternalServerError(
		"Internal server error",
		util.CodeUpdateUserFailed,
		fmt.Errorf("failed to update user: userId=%d, result=%d", userId, result),
	))
}
