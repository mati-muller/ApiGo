package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Conectado"})
	})

	setupRoutes(r) // Import the routes from routes.go
	setupPostRoutes(r) // Import the post routes from routes.go
	r.Run() // Por defecto en localhost:8080
}
