package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	
	v1 := r.Group("/v1")
	{
		api.SetRoutes(v1)
	}
	
	r.Run(":8080")
}
