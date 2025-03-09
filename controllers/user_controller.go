package controllers

import (
	"log"
	"net/http"
	"restfulapi/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FindAllUsers 所有會員清單
func FindAllUsers(c *gin.Context) {
	users := models.FindAllUsers()
	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}

// FindByUserId 單一會員資訊
func FindByUserId(c *gin.Context) {
	user := models.FindByUserId(c.Param("id"))
	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "not found",
		})
		return
	}
	log.Println("User ->", user)
	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

// DeleteUser 刪除會員
func DeleteUser(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("id"))

	result := models.DeleteUser(userId)
	if result == 1 {
		c.JSON(http.StatusNoContent, nil)
		return
	}

	c.JSON(http.StatusNotFound, gin.H{
		"message": "not found",
	})
}

// PutUser 更新會員資料
func PutUser(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.FindByUserId(c.Param("id"))

	result := models.UpdateUser(user, input)

	if result == 1 {
		user = models.FindByUserId(c.Param("id"))
		c.JSON(http.StatusOK, gin.H{
			"data": user,
		})
		return
	}

	c.JSON(http.StatusNotFound, gin.H{
		"message": "not found",
	})
}
