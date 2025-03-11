package repository

import (
	db "restfulapi/database"
	"restfulapi/models"
)

// UserRepository 定義存取 User 資料的方法
type UserRepository struct{}

// NewUserRepository 建立 UserRepository 實例
func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// FindAllUsers 取得所有使用者
func (r *UserRepository) FindAllUsers() []models.User {
	var users []models.User
	db.DbConnect.Find(&users)
	return users
}
