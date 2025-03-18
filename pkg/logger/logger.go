package logger

import (
	"fmt"
	"log"
	"restfulapi/config"
	"time"
)

// LogError - 記錄錯誤日誌並發送到 Kafka
func LogError(api string, err error) {
	errorLog := fmt.Sprintf(`{"timestamp": "%s", "api": "%s", "level": "error", "message": "%s"}`,
		time.Now().Format(time.RFC3339), api, err.Error())

	log.Println("❌", errorLog)
	config.PublishKafkaMessage("error-logs", errorLog) // 發送到 Kafka
}
