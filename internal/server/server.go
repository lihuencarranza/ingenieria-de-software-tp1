package server

import (
	"fmt"
	"log"
	"os"

	"melodia/internal/router"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Start initializes and starts the server
func Start() {
	// Load environment variables
	host := os.Getenv("HOST")
	if host == "" {
		log.Printf("Unable to load host from environment variables")
		os.Exit(1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Printf("Unable to load port from environment variables")
		os.Exit(1)
	}

	// Setup routes
	r := router.SetupRoutes()

	// Setup Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "melodia",
		})
	})

	// Start server
	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Server running on %s", addr)
	log.Printf("Swagger UI available at: http://%s/swagger/index.html", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
