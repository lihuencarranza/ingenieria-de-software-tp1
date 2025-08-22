package server

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func Start() {
	host := os.Getenv("HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()

	// Endpoint de prueba
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	addr := host + ":" + port
	log.Printf("Server running on %s", addr)
	r.Run(addr)
}
