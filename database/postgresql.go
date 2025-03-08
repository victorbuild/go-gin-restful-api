package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBconnect *gorm.DB

func DB() {
	var err error

	dsn := "host=localhost user=postgres password=admin dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Taipei"
	DBconnect, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	fmt.Println("🚀 Database connected successfully!")
}
