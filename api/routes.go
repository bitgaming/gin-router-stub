package api

import (
	"github.com/gin-gonic/gin"
)

func SetRoutes(rg *gin.RouterGroup) {
	var r_users = rg.Group("/users")
	{
		r_users.GET("", GetUsers)
		r_users.GET("/:id", GetUser)
		r_users.POST("", PostUser)
		r_users.PUT("/:id", PutUser)
		r_users.DELETE("/:id", DeleteUser)
	}
}

func GetUser(c *gin.Context) {
	c.JSON(200, gin.H{"username": "tedwen"})
}

func GetUsers(c *gin.Context) {
	users := []map[string]interface{}{
		map[string]interface{}{
			"username": "tedwen",
		},
		map[string]interface{}{
			"username": "adamalton",
		},
	}
	c.JSON(200, users)
}

func PostUser(c *gin.Context) {
	c.JSON(201, gin.H{"id": 1})
}

func PutUser(c *gin.Context) {
	c.JSON(200, gin.H{"status": "OK"})
}

func DeleteUser(c *gin.Context) {
	c.AbortWithStatus(405)
}
