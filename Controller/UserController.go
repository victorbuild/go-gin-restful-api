package controller

import (
	"net/http"
	models "restfulapi/Models"

	"github.com/gin-gonic/gin"
)

var users = []models.User{}

func FindAllUsers(c *gin.Context) {
	c.JSON(http.StatusOK, users)
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
	users = append(users, user)
	c.JSON(http.StatusOK, user)
}
