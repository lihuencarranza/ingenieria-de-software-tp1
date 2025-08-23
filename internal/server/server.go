package server

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func Start() {
	host := os.Getenv("HOST")
	if host == "" {
		log.Fatal("HOST environment variable not set")
		os.Exit(1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
		os.Exit(1)
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
