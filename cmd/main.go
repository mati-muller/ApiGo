package main

import (
	"ApiGo/routes"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Rate limiting middleware with a sliding window
func rateLimit() gin.HandlerFunc {
	limiter := make(map[string][]time.Time)
	var mu sync.RWMutex
	const maxRequests = 1000
	const window = time.Minute

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		now := time.Now()

		mu.Lock()
		defer mu.Unlock()

		// Clean up old requests outside the window
		requests := limiter[clientIP]
		filteredRequests := []time.Time{}
		for _, t := range requests {
			if now.Sub(t) <= window {
				filteredRequests = append(filteredRequests, t)
			}
		}
		limiter[clientIP] = filteredRequests

		// Check if the request limit is exceeded
		if len(filteredRequests) >= maxRequests {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}

		// Record the current request
		limiter[clientIP] = append(limiter[clientIP], now)
		c.Next()
	}
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	r := gin.Default()

	// Enable CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Adjust this to restrict origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Add rate limiting middleware
	r.Use(rateLimit())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Connected"})
	})

	// Setup routes
	routes.SetupPostRoutes(r)
	routes.SetupUserRoutes(r)
	routes.SetupProcesosRoutes(r)
	routes.SetupInventarioRoutes(r)
	routes.SetupEdits(r)
	routes.Reportes(r)
	routes.SetupUserDataRoutes(r)
	routes.SetupProcAppRoutes(r)
	routes.SetupEditsApp(r)

	// Use environment variable PORT or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to 8080 if PORT is not set
	}

	r.Run(":" + port) // Listen on the configured port
}
