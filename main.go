package main

// @title Go RESTful API Example
// @version 1.0
// @description 這是一個示範用的 API，包含 Kafka、JWT、RabbitMQ、Prometheus、Swagger 文件。
// @contact.name Victor
// @contact.email victor@email.com
// @host localhost:8000
// @BasePath /v1
import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"restfulapi/config"
	"restfulapi/database"
	"restfulapi/middlewares"
	"restfulapi/models"
	"restfulapi/routers"
	adminRoutes "restfulapi/routers/admin"
	"restfulapi/workers"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "restfulapi/docs"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "API 請求總數",
		},
		[]string{"method", "endpoint", "status_code"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "API 請求延遲",
		},
		[]string{"method", "endpoint"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

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
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// 加入 Prometheus `/metrics` 端點
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// **加入 Logger Middleware，確保所有請求都記錄到 Kafka**
	r.Use(middlewares.LoggerMiddleware())

	r.Use(PrometheusMiddleware())

	v1 := r.Group("/v1")
	// 管理員 API
	adminRoutes.SetupAdminUserRoutes(v1)

	// 一般使用者
	//routers.SetupUserRoutes(v1) // 一般用戶 API
	routers.SetupAuthRoutes(v1) // 註冊 / 登入 API

	// 啟動伺服器
	r.Run(":8000")
}

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		timer := prometheus.NewTimer(httpRequestDuration.WithLabelValues(c.Request.Method, c.FullPath()))
		c.Next()
		timer.ObserveDuration()

		httpRequestsTotal.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
			fmt.Sprintf("%d", c.Writer.Status()),
		).Inc()
	}
}
