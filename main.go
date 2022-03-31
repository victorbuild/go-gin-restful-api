package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "first test get api",
		})
	})

	router.POST("/users/:id", func(c *gin.Context) {
		userId := c.Param("id")

		c.JSON(http.StatusOK, gin.H{
			"id": userId,
		})
	})

	router.Run(":8000")
}
