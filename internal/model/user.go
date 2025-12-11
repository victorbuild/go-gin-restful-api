package model

import (
	"log"
	db "restfulapi/internal/database"
	"restfulapi/internal/util"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uint64 `gorm:"primaryKey;type:bigserial" json:"id"`
	Name     string `gorm:"size:255" json:"name"`
	Password string `gorm:"size:255" json:"password"`
	Email    string `gorm:"uniqueIndex;size:255" json:"email"`
	Role     string `gorm:"size:255" json:"role"`
}

// UpdateUserInput 用於更新使用者資訊
// @Description 更新使用者資訊的輸入結構，所有欄位都是可選的
// swagger:model UpdateUserInput
type UpdateUserInput struct {
	Name     string `json:"name,omitempty" example:"Victor"`              // 姓名（可選）
	Password string `json:"password,omitempty" example:"123456"`          // 密碼（可選）
	Email    string `json:"email,omitempty" example:"victor@example.com"` // email（可選）
}

func FindByUserId(userId uint64) (User, error) {
	var user User
	result := db.DbConnect.Where("id = ?", userId).First(&user)

	// 如果沒有找到使用者，回傳 `ErrUserNotFound`
	if result.RowsAffected == 0 {
		return User{}, util.ErrUserNotFound
	}

	// 其他錯誤（例如資料庫錯誤）
	if result.Error != nil {
		return User{}, result.Error
	}

	return user, nil
}

// 加密密碼
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// IsEmailExists - 確認 `email` 是否已存在（可選擇排除 `excludeUserID`）
func IsEmailExists(email string, excludeUserID uint64) bool {
	var count int64

	// 建立查詢條件
	query := db.DbConnect.Model(&User{}).Where("email = ?", email)

	// 如果 `excludeUserID` > 0，則排除該 ID
	if excludeUserID > 0 {
		query = query.Where("id != ?", excludeUserID)
	}

	query.Count(&count)
	return count > 0
}

func CreateUser(user User) (uint64, error) {
	// 檢查 Email 是否已存在
	if IsEmailExists(user.Email, 0) {
		return 0, util.ErrEmailExists
	}

	// 加密密碼
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		log.Println("Failed to hash password:", err)
		return 0, util.ErrPasswordHashFail
	}

	// 將加密後的密碼存入 User 結構
	user.Password = hashedPassword

	// 存入資料庫
	result := db.DbConnect.Create(&user)

	if result.Error != nil {
		log.Println("Failed to create user:", result.Error)
		return 0, util.ErrDatabaseError
	}

	return user.ID, nil
}

func DeleteUser(userId uint64) (int64, error) {
	result := db.DbConnect.Delete(&User{}, userId)

	if result.Error != nil {
		log.Println("刪除使用者失敗:", result.Error)
		return 0, result.Error
	}

	return result.RowsAffected, nil // 回傳影響的行數
}

func UpdateUser(user User, input UpdateUserInput) int64 {
	result := db.DbConnect.Model(&user).Updates(input)

	if result.Error != nil {
		log.Panic(result.Error)
	}

	return result.RowsAffected
}

// CheckPassword 驗證使用者輸入的密碼是否正確
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
