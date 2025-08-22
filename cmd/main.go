package main

import (
	"github.com/joho/godotenv"
	"log"
	"melodia/internal/server"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Unable to load file .env")
		os.Exit(1)
	}

	server.Start()
}
