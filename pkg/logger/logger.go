package logger

import (
	"fmt"
	"log"
	"restfulapi/config"
	"time"
)

// LogDebug 記錄除錯資訊（開發環境）
func LogDebug(api string, traceID string, message string) {
	debugLog := fmt.Sprintf(`{"timestamp": "%s", "api": "%s", "trace_id": "%s", "level": "debug", "message": "%s"}`,
		time.Now().Format(time.RFC3339), api, traceID, message)

	log.Println("[DEBUG]", debugLog)
	config.PublishKafkaMessage("log-topic", debugLog)
}

// LogInfo 記錄一般資訊（正常流程）
func LogInfo(api string, traceID string, message string) {
	infoLog := fmt.Sprintf(`{"timestamp": "%s", "api": "%s", "trace_id": "%s", "level": "info", "message": "%s"}`,
		time.Now().Format(time.RFC3339), api, traceID, message)

	log.Println("[INFO]", infoLog)
	config.PublishKafkaMessage("log-topic", infoLog)
}

// LogWarn 記錄警告（需要注意但不影響運作）
func LogWarn(api string, traceID string, message string) {
	warnLog := fmt.Sprintf(`{"timestamp": "%s", "api": "%s", "trace_id": "%s", "level": "warn", "message": "%s"}`,
		time.Now().Format(time.RFC3339), api, traceID, message)

	log.Println("[WARN]", warnLog)
	config.PublishKafkaMessage("log-topic", warnLog)
}

// LogError 記錄錯誤（需要處理）
func LogError(api string, traceID string, err error) {
	errorLog := fmt.Sprintf(`{"timestamp": "%s", "api": "%s", "trace_id": "%s", "level": "error", "message": "%s"}`,
		time.Now().Format(time.RFC3339), api, traceID, err.Error())

	log.Println("[ERROR]", errorLog)
	config.PublishKafkaMessage("log-topic", errorLog)
}

// LogFatal 記錄嚴重錯誤（程式會終止）
func LogFatal(api string, traceID string, err error) {
	fatalLog := fmt.Sprintf(`{"timestamp": "%s", "api": "%s", "trace_id": "%s", "level": "fatal", "message": "%s"}`,
		time.Now().Format(time.RFC3339), api, traceID, err.Error())

	log.Println("[FATAL]", fatalLog)
	config.PublishKafkaMessage("log-topic", fatalLog)
	log.Fatal(err)
}
