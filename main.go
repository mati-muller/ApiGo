package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Conectado"})
	})

	// Ensure SetupRoutes and SetupPostRoutes are accessible
	SetupRoutes(r)
	SetupPostRoutes(r)
	SetupUserRoutes(r)
	SetupProcesosRoutes(r)
	// Use environment variable PORT or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to 8080 if PORT is not set
	}

	r.Run(":" + port) // Listen on the configured port
}
