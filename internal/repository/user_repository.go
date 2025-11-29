package repository

import (
	db "restfulapi/internal/database"
	"restfulapi/internal/model"
)

// IUserRepository 定義存取 User 資料的方法
type IUserRepository interface {
	FindAllUsers() []model.User
}

// UserRepository 真實的 Repository（實作 IUserRepository）
type UserRepository struct{}

// NewUserRepository 建立 UserRepository 實例
func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// FindAllUsers 取得所有使用者（透過 GORM 存取 DB）
func (r *UserRepository) FindAllUsers() []model.User {
	var users []model.User
	db.DbConnect.Find(&users)
	return users
}
