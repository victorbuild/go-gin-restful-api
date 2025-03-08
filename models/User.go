package models

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	db "restfulapi/database"
)

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"` // 改為 uint，並標記為主鍵
	Name     string `gorm:"size:255" json:"name"`
	Password string `gorm:"size:255" json:"password"`
	Email    string `gorm:"uniqueIndex;size:255" json:"email"`
}

type UpdateUserInput struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func FindAllUsers() []User {
	var users []User
	db.DBconnect.Find(&users)
	return users
}

func FindByUserId(userId string) User {
	var user User
	db.DBconnect.Where("id = ?", userId).First(&user)
	return user
}

// 加密密碼
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CreateUser(user User) (uint, error) {
	// 加密密碼
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		log.Println("❌ Failed to hash password:", err)
		return 0, err
	}

	// 將加密後的密碼存入 User 結構
	user.Password = hashedPassword

	// 存入資料庫
	result := db.DBconnect.Create(&user)

	if result.Error != nil {
		log.Println("❌ Failed to create user:", result.Error)
		return 0, result.Error
	}

	return user.ID, nil
}

func DeleteUser(userId int) int64 {
	result := db.DBconnect.Delete(&User{}, userId)
	if result.Error != nil {
		log.Panic(result.Error)
	}

	return result.RowsAffected
}

func UpdateUser(user User, input User) int64 {
	result := db.DBconnect.Model(&user).Updates(input)

	if result.Error != nil {
		log.Panic(result.Error)
	}

	return result.RowsAffected
}

func Migrate() {
	err := db.DBconnect.AutoMigrate(&User{})
	if err != nil {
		log.Panic("❌ Database migration failed:", err)
	} else {
		log.Println("✅ Database migration completed!")
	}
}
