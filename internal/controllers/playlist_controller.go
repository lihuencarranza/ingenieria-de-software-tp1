package controllers

import (
	"net/http"
	"strconv"

	"melodia/internal/models"
	"melodia/internal/repositories"

	"github.com/gin-gonic/gin"
)

// PlaylistController handles playlist-related HTTP requests
type PlaylistController struct {
	playlistRepo *repositories.PlaylistRepository
}

// NewPlaylistController creates a new playlist controller
func NewPlaylistController() *PlaylistController {
	return &PlaylistController{
		playlistRepo: repositories.NewPlaylistRepository(),
	}
}

// CreatePlaylist handles POST /playlists
// @Summary Create a new playlist
// @Description Playlist created successfully
// @Tags playlists
// @Accept json
// @Produce json
// @Param playlist body models.CreatePlaylistRequest true "Playlist information"
// @Success 201 {object} models.PlaylistResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /playlists [post]
func (pc *PlaylistController) CreatePlaylist(c *gin.Context) {
	var req models.CreatePlaylistRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errorResp := models.NewErrorResponse("Bad Request", 400, "Invalid request body", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	// Validate required fields and length constraints
	if req.Name == "" || req.Description == "" {
		errorResp := models.NewErrorResponse("Bad Request", 400, "Name and description are required", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	// Validate description length (50-255 characters)
	if len(req.Description) < 50 {
		errorResp := models.NewErrorResponse("Bad Request", 400, "Description must be at least 50 characters long", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	if len(req.Description) > 255 {
		errorResp := models.NewErrorResponse("Bad Request", 400, "Description cannot exceed 255 characters", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	// Create playlist object
	playlist := models.Playlist{
		Name:        req.Name,
		Description: req.Description,
		IsPublished: false, // Playlists are created as unpublished by default
		PublishedAt: nil,   // Not published yet
		Songs:       []models.PlaylistSong{},
	}

	// Save to database
	if err := pc.playlistRepo.CreatePlaylist(&playlist); err != nil {
		errorResp := models.NewErrorResponse("Bad Request", 400, "Failed to create playlist", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	response := models.PlaylistResponse{
		Data: playlist,
	}

	c.JSON(http.StatusCreated, response)
}

// GetPlaylists handles GET /playlists
// @Summary Retrieve playlists (filter by published)
// @Description By default returns only published playlists ordered by publishedAt desc. Use published=false to get all playlists.
// @Tags playlists
// @Produce json
// @Param published query bool false "Filter by published status (default: true)"
// @Success 200 {object} models.PlaylistsResponse
// @Router /playlists [get]
func (pc *PlaylistController) GetPlaylists(c *gin.Context) {
	// Parse query parameter for published filter
	publishedStr := c.Query("published")
	var published *bool

	if publishedStr != "" {
		if publishedStr == "true" {
			val := true
			published = &val
		} else if publishedStr == "false" {
			val := false
			published = &val
		}
	}
	// If publishedStr is empty or invalid, published remains nil (default behavior)

	// Get from database with filter
	playlists, err := pc.playlistRepo.GetPlaylists(published)
	if err != nil {
		errorResp := models.NewErrorResponse("Bad Request", 400, "Failed to retrieve playlists", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	response := models.PlaylistsResponse{
		Data: playlists,
	}

	c.JSON(http.StatusOK, response)
}

// GetPlaylist handles GET /playlists/{id}
// @Summary Retrieve a playlist by ID
// @Description Get a specific playlist by its ID with songs ordered by addedAt desc
// @Tags playlists
// @Produce json
// @Param id path int true "Playlist ID"
// @Success 200 {object} models.PlaylistResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /playlists/{id} [get]
func (pc *PlaylistController) GetPlaylist(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errorResp := models.NewErrorResponse("Bad Request", 400, "Invalid playlist ID", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	// Get from database by ID
	playlist, err := pc.playlistRepo.GetPlaylistByID(uint(id))
	if err != nil {
		if err.Error() == "playlist not found" {
			errorResp := models.NewErrorResponse("Not Found", 404, "Playlist not found", c.Request.URL.Path)
			c.JSON(http.StatusNotFound, errorResp)
			return
		}
		errorResp := models.NewErrorResponse("Bad Request", 400, "Failed to retrieve playlist", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	response := models.PlaylistResponse{
		Data: *playlist,
	}

	c.JSON(http.StatusOK, response)
}

// DeletePlaylist handles DELETE /playlists/{id}
// @Summary Delete a playlist by ID
// @Description Delete a specific playlist by its ID
// @Tags playlists
// @Param id path int true "Playlist ID"
// @Success 204 "No Content"
// @Failure 404 {object} models.ErrorResponse
// @Router /playlists/{id} [delete]
func (pc *PlaylistController) DeletePlaylist(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errorResp := models.NewErrorResponse("Bad Request", 400, "Invalid playlist ID", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	// Delete from database
	if err := pc.playlistRepo.DeletePlaylist(uint(id)); err != nil {
		if err.Error() == "playlist not found" {
			errorResp := models.NewErrorResponse("Not Found", 404, "Playlist not found", c.Request.URL.Path)
			c.JSON(http.StatusNotFound, errorResp)
			return
		}
		errorResp := models.NewErrorResponse("Bad Request", 400, "Failed to delete playlist", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	c.Status(http.StatusNoContent)
}

// PublishPlaylist handles POST /playlists/{id}/publish
// @Summary Publish a playlist (idempotent)
// @Description Sets isPublished=true and publishedAt=now() if not already published
// @Tags playlists
// @Accept json
// @Produce json
// @Param id path int true "Playlist ID"
// @Success 200 {object} models.PlaylistResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /playlists/{id}/publish [post]
func (pc *PlaylistController) PublishPlaylist(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errorResp := models.NewErrorResponse("Bad Request", 400, "Invalid playlist ID", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	// Get playlist from database
	playlist, err := pc.playlistRepo.GetPlaylistByID(uint(id))
	if err != nil {
		if err.Error() == "playlist not found" {
			errorResp := models.NewErrorResponse("Not Found", 404, "Playlist not found", c.Request.URL.Path)
			c.JSON(http.StatusNotFound, errorResp)
			return
		}
		errorResp := models.NewErrorResponse("Bad Request", 400, "Failed to retrieve playlist", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	// Publish playlist (idempotent - if already published, just return success)
	if !playlist.IsPublished {
		if err := pc.playlistRepo.PublishPlaylist(uint(id)); err != nil {
			errorResp := models.NewErrorResponse("Bad Request", 400, "Failed to publish playlist", c.Request.URL.Path)
			c.JSON(http.StatusBadRequest, errorResp)
			return
		}
		// Get updated playlist
		playlist, err = pc.playlistRepo.GetPlaylistByID(uint(id))
		if err != nil {
			errorResp := models.NewErrorResponse("Bad Request", 400, "Failed to retrieve updated playlist", c.Request.URL.Path)
			c.JSON(http.StatusBadRequest, errorResp)
			return
		}
	}

	response := models.PlaylistResponse{
		Data: *playlist,
	}

	c.JSON(http.StatusOK, response)
}

// AddSongToPlaylist handles POST /playlists/{id}/songs
// @Summary Add a song to a playlist
// @Description Add an existing song to a playlist
// @Tags playlists
// @Accept json
// @Produce json
// @Param id path int true "Playlist ID"
// @Param song body models.AddSongToPlaylistRequest true "Song to add"
// @Success 200 {object} models.PlaylistResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /playlists/{id}/songs [post]
func (pc *PlaylistController) AddSongToPlaylist(c *gin.Context) {
	idStr := c.Param("id")
	playlistID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errorResp := models.NewErrorResponse("Bad Request", 400, "Invalid playlist ID", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	// Accept both songId (contract) and song_id (legacy)
	type addSongBody struct {
		SongID      *uint `json:"songId"`
		SongIDSnake *uint `json:"song_id"`
	}
	var body addSongBody
	if err := c.ShouldBindJSON(&body); err != nil {
		errorResp := models.NewErrorResponse("Bad Request", 400, "Invalid request body", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	var songID uint
	if body.SongID != nil {
		songID = *body.SongID
	} else if body.SongIDSnake != nil {
		songID = *body.SongIDSnake
	} else {
		errorResp := models.NewErrorResponse("Bad Request", 400, "Missing required field: songId", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	// Add song to playlist
	if err := pc.playlistRepo.AddSongToPlaylist(uint(playlistID), songID); err != nil {
		if err.Error() == "playlist not found" {
			errorResp := models.NewErrorResponse("Not Found", 404, "Playlist not found", c.Request.URL.Path)
			c.JSON(http.StatusNotFound, errorResp)
			return
		}
		if err.Error() == "song not found" {
			errorResp := models.NewErrorResponse("Not Found", 404, "Song not found", c.Request.URL.Path)
			c.JSON(http.StatusNotFound, errorResp)
			return
		}
		errorResp := models.NewErrorResponse("Bad Request", 400, "Failed to add song to playlist", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	// Get updated playlist from database
	playlist, err := pc.playlistRepo.GetPlaylistByID(uint(playlistID))
	if err != nil {
		errorResp := models.NewErrorResponse("Bad Request", 400, "Failed to retrieve updated playlist", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	response := models.PlaylistResponse{
		Data: *playlist,
	}

	c.JSON(http.StatusOK, response)
}
