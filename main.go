package main

import (
	routers "restfulapi/Routers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	v1 := r.Group("/v1")

	routers.AddUserRouter(v1)

	r.Run(":8000")
}
