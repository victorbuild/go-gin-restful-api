package database

import (
	"fmt"
	"log"
	"restfulapi/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DbConnect *gorm.DB

func DB() {
	var err error
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Taipei",
		config.DatabaseConfig.Host,
		config.DatabaseConfig.User,
		config.DatabaseConfig.Password,
		config.DatabaseConfig.DBName,
		config.DatabaseConfig.Port,
	)

	DbConnect, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	fmt.Println("🚀 Database connected successfully!")
}
