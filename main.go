package main

import (
	routers "restfulapi/Routers"
	db "restfulapi/database"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	v1 := r.Group("/v1")

	routers.AddUserRouter(v1)

	go func() {
		db.DB()
	}()

	r.Run(":8000")
}
