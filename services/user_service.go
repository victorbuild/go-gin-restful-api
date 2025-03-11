package service

import (
	"restfulapi/models"
	repository "restfulapi/repositories"
)

type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService 建立一個新的 UserService 實例
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetAllUsers 取得所有使用者，未來可在此加入其他邏輯或調用 repository 層
func (s *UserService) GetAllUsers() []models.User {
	// 這裡暫時直接呼叫 models 層，後續可調整為調用 repository
	return s.userRepo.FindAllUsers()
}
