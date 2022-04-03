package controller

import (
	"log"
	"net/http"
	models "restfulapi/Models"
	"strconv"

	"github.com/gin-gonic/gin"
)

var users = []models.User{}

func FindAllUsers(c *gin.Context) {
	users := models.FindAllUsers()
	c.JSON(http.StatusOK, users)
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
	c.JSON(http.StatusOK, user)
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
	createUserId := models.CreateUser(user)
	user.ID = createUserId
	c.JSON(http.StatusOK, user)
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
	beforeUser := models.User{}
	err := c.BindJSON(&beforeUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "data error",
		})
		return
	}
	userId, _ := strconv.Atoi(c.Param("id"))
	for key, user := range users {
		if userId == user.ID {
			users[key] = beforeUser
			log.Println(users[key])
			c.JSON(http.StatusOK, users[key])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"message": "not found",
	})
}
