package config

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

// RedisClient 全局 Redis 變數
var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", GetEnv("REDIS_HOST", "localhost"), GetEnv("REDIS_PORT", "6379")),
		Password: GetEnv("REDIS_PASSWORD", ""), // 若無密碼則預設空值
		DB:       GetEnvInt("REDIS_DB", 0),     // Redis 預設 DB 0
	})

	// 測試 Redis 連線
	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("❌ Failed to connect to Redis:", err)
	}

	log.Println("🚀 Redis connected successfully!")
}
