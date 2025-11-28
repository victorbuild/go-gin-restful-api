package middlewares

import (
	"fmt"
	"restfulapi/pkg/logger"
	"restfulapi/utils"

	"github.com/gin-gonic/gin"
)

// ErrorHandler 統一錯誤處理 middleware
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			var appErr *utils.AppError
			if appError, ok := err.(*utils.AppError); ok {
				appErr = appError
			} else {
				appErr = utils.NewInternalServerError("Internal server error", utils.CodeInternalError, err)
			}

			// 取得 Trace ID
			traceID := GetTraceID(c)

			// 記錄錯誤（根據狀態碼決定 log 等級）
			if appErr.StatusCode >= 500 {
				// 500 錯誤：記錄詳細錯誤（用於 debug）
				logger.LogError(
					fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path),
					traceID,
					appErr.Err,
				)
			} else if appErr.StatusCode == 401 || appErr.StatusCode == 403 {
				// 401/403 錯誤：記錄警告（安全相關）
				logger.LogWarn(
					fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path),
					traceID,
					appErr.Message,
				)
			}

			utils.ErrorResponse(c, appErr.StatusCode, appErr.Message, appErr.Code)
			c.Errors = c.Errors[:0]
		}
	}
}
