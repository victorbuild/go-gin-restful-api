package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// TraceIDHeader 定義 Trace ID 的 HTTP Header 名稱
	TraceIDHeader = "X-Trace-ID"
	// TraceIDKey 定義在 gin.Context 中儲存 Trace ID 的 key
	TraceIDKey = "trace_id"
)

// TraceIDMiddleware 生成或讀取 Trace ID，並設定到 response header
func TraceIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 嘗試從 request header 讀取 Trace ID
		traceID := c.GetHeader(TraceIDHeader)

		// 如果沒有，則生成新的 UUID
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// 將 Trace ID 存入 context，供後續使用
		c.Set(TraceIDKey, traceID)

		// 設定 response header
		c.Header(TraceIDHeader, traceID)

		c.Next()
	}
}

// GetTraceID 從 gin.Context 取得 Trace ID（輔助函數）
func GetTraceID(c *gin.Context) string {
	if traceID, exists := c.Get(TraceIDKey); exists {
		if id, ok := traceID.(string); ok {
			return id
		}
	}
	return ""
}
