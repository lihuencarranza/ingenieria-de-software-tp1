package main

import (
	"log"
	"os"

	"melodia/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Unable to load file .env")
		os.Exit(1)
	}

	server.Start()
}
