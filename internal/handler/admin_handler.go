package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

// AdminUserDetailData 定義單一使用者回應的 data 結構（Swagger 用）
// swagger:model AdminUserDetailData
type AdminUserDetailData struct {
	ID    uint64 `json:"id" example:"1"`
	Name  string `json:"name" example:"Victor"`
	Email string `json:"email" example:"victor@example.com"`
	Role  string `json:"role" example:"user"`
}

// AdminUserDetailSuccessResponse 定義單一使用者成功回應（Swagger 用）
// swagger:model AdminUserDetailSuccessResponse
type AdminUserDetailSuccessResponse struct {
	Status  string              `json:"status" example:"success"`
	Message string              `json:"message" example:"User retrieved successfully"`
	Data    AdminUserDetailData `json:"data"`
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
// @Summary 取得單一使用者資料
// @Description 根據使用者 ID 取得單一使用者資料（需管理員權限）
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "使用者 ID" example(1)
// @Success 200 {object} handler.AdminUserDetailSuccessResponse "成功取得使用者資料"
// @Failure 400 {object} util.ErrorAPIResponse "Bad Request，error_code: 1001（無效的使用者 ID）"
// @Failure 401 {object} util.ErrorAPIResponse "Unauthorized，error_code: 1006（Token 缺失）或 1007（Token 無效）"
// @Failure 403 {object} util.ErrorAPIResponse "Forbidden，error_code: 1010（權限不足）"
// @Failure 404 {object} util.ErrorAPIResponse "Not Found，error_code: 1005（使用者不存在）"
// @Failure 500 {object} util.ErrorAPIResponse "Internal Server Error，error_code: 4001（伺服器內部錯誤）"
// @Router /v1/admin/users/{id} [get]
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
// @Summary 刪除使用者
// @Description 根據使用者 ID 刪除使用者（需管理員權限）
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "使用者 ID" example(1)
// @Success 200 {object} util.SuccessAPIResponse "成功刪除使用者，data 為 null"
// @Failure 400 {object} util.ErrorAPIResponse "Bad Request，error_code: 1001（無效的使用者 ID）"
// @Failure 401 {object} util.ErrorAPIResponse "Unauthorized，error_code: 1006（Token 缺失）或 1007（Token 無效）"
// @Failure 403 {object} util.ErrorAPIResponse "Forbidden，error_code: 1010（權限不足）"
// @Failure 404 {object} util.ErrorAPIResponse "Not Found，error_code: 1009（使用者不存在）"
// @Failure 500 {object} util.ErrorAPIResponse "Internal Server Error，error_code: 1011（刪除使用者失敗）或 4001（伺服器內部錯誤）"
// @Router /v1/admin/users/{id} [delete]
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
// @Summary 更新使用者資訊
// @Description 根據使用者 ID 更新使用者資訊（需管理員權限）。所有欄位都是可選的，只會更新提供的欄位。
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "使用者 ID" example(1)
// @Param input body model.UpdateUserInput true "更新使用者資訊" example({"name":"Victor","email":"victor@example.com","password":"123456"})
// @Success 200 {object} handler.AdminUserDetailSuccessResponse "成功更新使用者資訊"
// @Failure 400 {object} util.ErrorAPIResponse "Bad Request，error_code: 1001（無效的使用者 ID 或輸入格式錯誤）"
// @Failure 401 {object} util.ErrorAPIResponse "Unauthorized，可能的 error_code: 1006（Token 缺失）、1007（Token 無效）。詳細說明請參考錯誤代碼文檔"
// @Failure 403 {object} util.ErrorAPIResponse "Forbidden，error_code: 1010（權限不足）。詳細說明請參考錯誤代碼文檔"
// @Failure 404 {object} util.ErrorAPIResponse "Not Found，error_code: 1009（使用者不存在）。詳細說明請參考錯誤代碼文檔"
// @Failure 409 {object} util.ErrorAPIResponse "Conflict，error_code: 1013（Email 已被其他使用者使用）。詳細說明請參考錯誤代碼文檔"
// @Failure 500 {object} util.ErrorAPIResponse "Internal Server Error，可能的 error_code: 4003（資料庫錯誤）、1012（更新使用者失敗）。詳細說明請參考錯誤代碼文檔"
// @Router /v1/admin/users/{id} [put]
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

// AdminCreateUserInput 定義管理員建立使用者輸入結構
// @Description 管理員建立使用者資訊的輸入結構，role 為可選欄位（預設為 "user"）
// swagger:model AdminCreateUserInput
type AdminCreateUserInput struct {
	Name     string `json:"name" example:"Victor" binding:"required"`                    // 姓名
	Email    string `json:"email" example:"victor@example.com" binding:"required,email"` // Email
	Password string `json:"password" example:"123456" binding:"required"`                // 密碼
	Role     string `json:"role,omitempty" example:"user"`                               // 角色（可選，預設為 "user"，可選值：user, admin）
}

// AdminCreateUser - 管理員建立使用者
// @Summary 管理員建立使用者
// @Description 管理員建立新使用者（需管理員權限）。與一般註冊 `/v1/auth/register` 不同，管理員可以指定使用者的 role（user 或 admin）。如果不提供 role，預設為 "user"。
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body handler.AdminCreateUserInput true "使用者資訊" example({"name":"Victor","email":"victor@example.com","password":"123456","role":"user"})
// @Success 201 {object} handler.AdminUserDetailSuccessResponse "成功建立使用者"
// @Failure 400 {object} util.ErrorAPIResponse "Bad Request，error_code: 1001（無效的輸入格式、無效的 role 或缺少必填欄位）或 1002（缺少必填欄位）"
// @Failure 401 {object} util.ErrorAPIResponse "Unauthorized，可能的 error_code: 1006（Token 缺失）、1007（Token 無效）。詳細說明請參考錯誤代碼文檔"
// @Failure 403 {object} util.ErrorAPIResponse "Forbidden，error_code: 1010（權限不足）。詳細說明請參考錯誤代碼文檔"
// @Failure 409 {object} util.ErrorAPIResponse "Conflict，error_code: 1003（Email 已經被註冊）。詳細說明請參考錯誤代碼文檔"
// @Failure 415 {object} util.ErrorAPIResponse "Unsupported Media Type，error_code: 1000（Content-Type 不是 application/json）。詳細說明請參考錯誤代碼文檔"
// @Failure 500 {object} util.ErrorAPIResponse "Internal Server Error，error_code: 4001（伺服器內部錯誤）。詳細說明請參考錯誤代碼文檔"
// @Router /v1/admin/users [post]
func AdminCreateUser(c *gin.Context) {
	// 驗證 Content-Type
	if !validateContentType(c) {
		return
	}

	var input AdminCreateUserInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		util.ValidationErrorResponse(c, err)
		return
	}

	// 驗證 role（如果提供）
	if input.Role != "" && input.Role != "user" && input.Role != "admin" {
		c.Error(util.NewBadRequestError("Invalid role. Role must be 'user' or 'admin'", util.CodeInvalidInput, nil))
		return
	}

	// 如果沒有提供 role，預設為 "user"
	if input.Role == "" {
		input.Role = "user"
	}

	user := model.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		Role:     input.Role, // Admin 可以指定 role
	}

	createUserId, createUserErr := model.CreateUser(user)
	if createUserErr != nil {
		handleCreateUserError(c, createUserErr, "admin create user")
		return
	}

	user.ID = createUserId

	// 定義 RabbitMQ 訊息結構（不包含密碼）
	type UserCreatedMessage struct {
		ID    uint64 `json:"id"`
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

	log.Println("管理員建立使用者成功，已發送 user_created 事件:", string(message))

	// 回應標準 JSON
	util.CreatedResponse(c, "User created successfully", gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	})
}
