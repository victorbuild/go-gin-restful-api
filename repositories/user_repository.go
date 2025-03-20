package repositories

import (
	db "restfulapi/database"
	"restfulapi/models"
)

// IUserRepository 定義存取 User 資料的方法
type IUserRepository interface {
	FindAllUsers() []models.User
}

// UserRepository 真實的 Repository（實作 IUserRepository）
type UserRepository struct{}

// NewUserRepository 建立 UserRepository 實例
func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// FindAllUsers 取得所有使用者（透過 GORM 存取 DB）
func (r *UserRepository) FindAllUsers() []models.User {
	var users []models.User
	db.DbConnect.Find(&users)
	return users
}
