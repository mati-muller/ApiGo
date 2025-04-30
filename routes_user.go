package main

import (
	"github.com/gin-gonic/gin"
)

func registerUserRoutes(r *gin.Engine) {
	userGroup := r.Group("/users")
	{
		userGroup.GET("/", getAllUsers)
		userGroup.POST("/", createUser)
	}
}

func getAllUsers(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Get all users"})
}

func createUser(c *gin.Context) {
	c.JSON(201, gin.H{"message": "Create user"})
}