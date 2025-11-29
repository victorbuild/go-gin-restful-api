package middlewares

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"restfulapi/internal/config"
)

// LoggerMiddleware - 記錄 API 請求資訊並發送到 Kafka
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next() // 執行 API

		latency := time.Since(start).Milliseconds()
		statusCode := c.Writer.Status()
		traceID := GetTraceID(c)
		requestLog := fmt.Sprintf(
			`{"event":"log.api_request","timestamp":"%s","trace_id":"%s","method":"%s","endpoint":"%s","status_code":%d,"response_time":%d}`,
			start.Format(time.RFC3339), traceID, c.Request.Method, c.Request.URL.Path, statusCode, latency,
		)

		config.PublishKafkaMessage("log-topic", requestLog)
	}
}
