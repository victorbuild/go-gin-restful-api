package controllers

import (
	"restfulapi/utils"

	"github.com/gin-gonic/gin"
)

// LogoutUser - 使用者登出（JWT 無狀態登出）
func LogoutUser(c *gin.Context) {
	// 只需要回應成功，前端會移除 Token
	utils.SuccessResponse(c, "Logout successful", gin.H{})
}
