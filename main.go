package main

import (
	"fmt"
	"restfulapi/database"
	"restfulapi/models"
	"restfulapi/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	// **初始化資料庫**
	database.DB() // ✅ 確保這行執行後 `DBconnect` 不為 nil

	// **確保 `DBconnect` 已正確初始化後再執行 Migrate**
	if database.DBconnect == nil {
		fmt.Println("❌ Database connection failed, cannot run migration!")
		return
	}

	fmt.Println("📌 Running database migration...")
	database.DBconnect.AutoMigrate(&models.User{})
	fmt.Println("✅ Database migration completed!")

	// 初始化 Gin 路由
	r := gin.Default()
	v1 := r.Group("/v1")
	routers.SetupUserRoutes(v1) // 一般用戶 API
	routers.SetupAuthRoutes(v1) // 註冊 / 登入 API

	// 啟動伺服器
	r.Run(":8000")
}
