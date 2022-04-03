package models

import (
	"log"
	db "restfulapi/database"
)

type User struct {
	ID       int    `json:"id"`
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

func CreateUser(user User) int {
	result := db.DBconnect.Create(&user)

	if result.Error != nil {
		log.Panic(result.Error)
	}

	return user.ID
}
