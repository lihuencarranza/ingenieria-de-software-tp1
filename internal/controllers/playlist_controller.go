package controllers

import (
	"net/http"
	"strconv"
	"time"

	"melodia/internal/models"

	"github.com/gin-gonic/gin"
)

// PlaylistController handles playlist-related HTTP requests
type PlaylistController struct {
	// TODO: Add playlist service/repository dependency
}

// NewPlaylistController creates a new playlist controller
func NewPlaylistController() *PlaylistController {
	return &PlaylistController{}
}

// CreatePlaylist handles POST /playlists
// @Summary Create a new playlist
// @Description Create a new playlist with name and description
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

	// TODO: Validate request
	if req.Name == "" || req.Description == "" {
		errorResp := models.NewErrorResponse("Bad Request", 400, "Name and description are required", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	// TODO: Save to database
	now := time.Now()
	playlist := models.Playlist{
		ID:          1, // TODO: Generate proper ID
		Name:        req.Name,
		Description: req.Description,
		IsPublished: true, // Playlists are created as published by default
		PublishedAt: &now, // Published at creation time
		Songs:       []models.PlaylistSong{},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	response := models.PlaylistResponse{
		Data: playlist,
	}

	c.JSON(http.StatusCreated, response)
}

// GetPlaylists handles GET /playlists
// @Summary Retrieve published playlists
// @Description Get a list of published playlists ordered by publishedAt desc
// @Tags playlists
// @Produce json
// @Success 200 {object} models.PlaylistsResponse
// @Router /playlists [get]
func (pc *PlaylistController) GetPlaylists(c *gin.Context) {
	// TODO: Get from database, ordered by publishedAt desc
	playlists := []models.Playlist{} // Empty for now

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

	// TODO: Get from database by ID
	// For now, return a mock playlist
	now := time.Now()
	playlist := models.Playlist{
		ID:          uint(id),
		Name:        "Sample Playlist",
		Description: "A sample playlist description",
		IsPublished: true,
		PublishedAt: &now,
		Songs:       []models.PlaylistSong{},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	response := models.PlaylistResponse{
		Data: playlist,
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
	_, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errorResp := models.NewErrorResponse("Bad Request", 400, "Invalid playlist ID", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	// TODO: Delete from database
	c.Status(http.StatusNoContent)
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

	var req models.AddSongToPlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResp := models.NewErrorResponse("Bad Request", 400, "Invalid request body", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	// TODO: Validate song exists and add to playlist
	// TODO: Get updated playlist from database
	now := time.Now()
	playlist := models.Playlist{
		ID:          uint(playlistID),
		Name:        "Sample Playlist",
		Description: "A sample playlist description",
		IsPublished: true,
		PublishedAt: &now,
		Songs:       []models.PlaylistSong{},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	response := models.PlaylistResponse{
		Data: playlist,
	}

	c.JSON(http.StatusOK, response)
}
