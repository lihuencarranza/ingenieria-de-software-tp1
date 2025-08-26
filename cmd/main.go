package main

import (
	"log"
	"melodia/internal/server"
	"os"

	"github.com/joho/godotenv"
	_ "melodia/docs" // Importar docs generados por swag
)

// @title           Melodía API
// @version         1.0
// @description     API para la plataforma de streaming musical Melodía
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @tag.name songs
// @tag.description Operaciones relacionadas con canciones

// @tag.name playlists
// @tag.description Operaciones relacionadas con playlists

func main() {
	// Cargar variables de entorno
	err := godotenv.Load()
	if err != nil {
		log.Println("Unable to load file .env")
		os.Exit(1)
	}

	// Iniciar el servidor
	server.Start()
}
