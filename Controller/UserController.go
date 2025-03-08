package controller

import (
	"log"
	"net/http"
	"restfulapi/models"
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

func PostUser(c *gin.Context) {
	user := models.User{}
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Error",
		})
		return
	}
	createUserId, createUserErr := models.CreateUser(user)

	if createUserErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error",
		})
		return
	}

	user.ID = createUserId
	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
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
