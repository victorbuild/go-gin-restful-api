package main

// @title Go RESTful API Example
// @version 1.0
// @description 這是一個示範用的 API，包含 Kafka、JWT、RabbitMQ、Prometheus、Swagger 文件。
// @contact.name Victor
// @contact.email victor@email.com
// @host localhost:8000
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 輸入 "Bearer {token}" 進行身份驗證。可以在登入 API 取得 Token。
import (
	"fmt"
	"restfulapi/config"
	"restfulapi/database"
	"restfulapi/middlewares"
	"restfulapi/models"
	"restfulapi/routers"
	adminRoutes "restfulapi/routers/admin"
	"restfulapi/workers"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	_ "restfulapi/docs"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// 初始化 Redis 連線
	config.InitRedis()

	config.InitKafka()

	// 初始化資料庫
	database.DB()

	if database.DbConnect == nil {
		fmt.Println("❌ Database connection failed, cannot run migration!")
		return
	}

	fmt.Println("📌 Running database migration...")
	database.DbConnect.AutoMigrate(&models.User{})
	fmt.Println("✅ Database migration completed!")

	// **啟動 RabbitMQ Worker（使用 Goroutine）**
	go workers.StartUserWorker()
	// **啟動 Kafka Log Worker**
	go workers.StartLogWorker()

	// 初始化 Gin 路由
	r := gin.Default()

	// 設定健康檢查路由（根目錄和 /health）
	routers.SetupHealthRoutes(r)

	// Swagger 文件
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Prometheus metrics
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// **加入 Trace ID Middleware
	r.Use(middlewares.TraceIDMiddleware())

	// **加入 Logger Middleware，確保所有請求都記錄到 Kafka**
	r.Use(middlewares.LoggerMiddleware())

	// **加入 Prometheus Middleware，記錄 API 指標**
	r.Use(middlewares.PrometheusMiddleware())

	// 錯誤處理 middleware（必須放在最後）
	r.Use(middlewares.ErrorHandler())

	v1 := r.Group("/v1")
	// 管理員 API
	adminRoutes.SetupAdminUserRoutes(v1)

	// 一般使用者
	//routers.SetupUserRoutes(v1) // 一般用戶 API
	routers.SetupAuthRoutes(v1) // 註冊 / 登入 API

	// 啟動伺服器
	r.Run(":8000")
}
