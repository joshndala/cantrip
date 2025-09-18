package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joshndala/cantrip/router"
)

func main() {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create router
	r := gin.Default()

	// Disable automatic trailing slash redirect to prevent CORS issues
	r.RedirectTrailingSlash = false

	// Setup CORS middleware with proper configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD", "PATCH"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Accept", "Cache-Control", "X-Requested-With"}
	config.AllowCredentials = true
	config.MaxAge = 12 * 3600 // 12 hours
	config.AllowWildcard = true

	r.Use(cors.New(config))

	// Additional CORS middleware for debugging
	r.Use(func(c *gin.Context) {
		log.Printf("Request: %s %s from %s", c.Request.Method, c.Request.URL.Path, c.Request.Header.Get("Origin"))
		c.Next()
	})

	// Setup routes
	router.SetupRoutes(r)

	// Start server
	log.Println("Starting CanTrip API server on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
