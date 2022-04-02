package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBconnect *gorm.DB

var err error

func DB() {
	dsn := "host=localhost user=postgres password=1234 dbname=testgoapi port=5432 sslmode=disable TimeZone=Asia/Taipei"
	DBconnect, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}
}
