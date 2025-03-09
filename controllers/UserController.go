package controllers

import (
	"errors"
	"log"
	"net/http"
	db "restfulapi/database"
	"restfulapi/models"
	"restfulapi/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func FindAllUsers(c *gin.Context) {
	users := models.FindAllUsers()
	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}

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

func RegisterUser(c *gin.Context) {
	user := models.User{}
	err := c.BindJSON(&user)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", utils.ErrInvalidInput)
		return
	}
	createUserId, createUserErr := models.CreateUser(user)

	if createUserErr != nil {
		switch {
		case errors.Is(createUserErr, models.ErrEmailExists):
			utils.ErrorResponse(c, http.StatusConflict, "Email already registered", utils.ErrEmailExists)
		case errors.Is(createUserErr, models.ErrPasswordHashFail):
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user", utils.ErrPasswordHashFail)
		case errors.Is(createUserErr, models.ErrDatabaseError):
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user", utils.ErrDatabaseError)
		default:
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user", utils.ErrInternalError)
		}
		return
	}

	user.ID = createUserId

	// **回應標準 JSON**
	utils.SuccessResponse(c, "User created successfully", gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

func LoginUser(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// 解析請求 JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", utils.ErrInvalidInput)
		return
	}

	// 查詢使用者
	var user models.User
	db.DBconnect.Where("email = ?", input.Email).First(&user)

	// 檢查密碼
	if !user.CheckPassword(input.Password) {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password", utils.ErrInvalidCredentials)
		return
	}

	// **回應標準 JSON**
	utils.SuccessResponse(c, "Login successful!", gin.H{})
}

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
