package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// LoadEnv 讀取 .env 檔案
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("找不到 .env 檔案，使用系統環境變數")
	}
}

// getEnv 讀取環境變數，若沒有則回傳預設值
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvInt 讀取環境變數並轉換為 int
func getEnvInt(key string, defaultValue int) int {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
