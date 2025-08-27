package router

import (
	"melodia/internal/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes() *gin.Engine {
	router := gin.Default()

	// Initialize controllers
	songController := controllers.NewSongController()
	playlistController := controllers.NewPlaylistController()

	// Songs routes
	songs := router.Group("/songs")
	{
		songs.POST("", songController.CreateSong)
		songs.GET("", songController.GetSongs)
		songs.GET("/:id", songController.GetSong)
		songs.PUT("/:id", songController.UpdateSong)
		songs.DELETE("/:id", songController.DeleteSong)
	}

	// Playlists routes
	playlists := router.Group("/playlists")
	{
		playlists.POST("", playlistController.CreatePlaylist)
		playlists.GET("", playlistController.GetPlaylists)
		playlists.GET("/:id", playlistController.GetPlaylist)
		playlists.DELETE("/:id", playlistController.DeletePlaylist)
		playlists.POST("/:id/songs", playlistController.AddSongToPlaylist)
	}

	return router
}
