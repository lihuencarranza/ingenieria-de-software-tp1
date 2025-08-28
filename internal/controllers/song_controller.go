package controllers

import (
	"net/http"
	"strconv"

	"melodia/internal/models"
	"melodia/internal/repositories"

	"github.com/gin-gonic/gin"
)

// SongController handles song-related HTTP requests
type SongController struct {
	songRepo *repositories.SongRepository
}

// NewSongController creates a new song controller
func NewSongController() *SongController {
	return &SongController{
		songRepo: repositories.NewSongRepository(),
	}
}

// CreateSong handles POST /songs
// @Summary Create a new song
// @Description Create a new song with title and artist
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.CreateSongRequest true "Song information"
// @Success 201 {object} models.SongResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /songs [post]
func (sc *SongController) CreateSong(c *gin.Context) {
	var req models.CreateSongRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errorResp := models.NewErrorResponse("Bad Request", 400, "Invalid request body", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	if req.Title == "" || req.Artist == "" {
		errorResp := models.NewErrorResponse("Bad Request", 400, "Title and artist are required", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	song := &models.Song{
		Title:  req.Title,
		Artist: req.Artist,
	}

	if err := sc.songRepo.CreateSong(song); err != nil {
		errorResp := models.NewErrorResponse("Bad Request", 400, "Failed to create song", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	response := models.SongResponse{
		Data: *song,
	}

	c.JSON(http.StatusCreated, response)
}

// GetSongs handles GET /songs
// @Summary Retrieve all songs
// @Description Get a list of all songs
// @Tags songs
// @Produce json
// @Success 200 {object} models.SongsResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /songs [get]
func (sc *SongController) GetSongs(c *gin.Context) {
	songs, err := sc.songRepo.GetSongs()
	if err != nil {
		errorResp := models.NewErrorResponse("Bad Request", 400, "Failed to retrieve songs", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	response := models.SongsResponse{
		Data: songs,
	}

	c.JSON(http.StatusOK, response)
}

// GetSong handles GET /songs/{id}
// @Summary Retrieve a song by ID
// @Description Song retrieved successfully
// @Tags songs
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {object} models.SongResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /songs/{id} [get]
func (sc *SongController) GetSong(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errorResp := models.NewErrorResponse("Bad Request", 400, "Invalid song ID", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	song, err := sc.songRepo.GetSongByID(uint(id))
	if err != nil {
		if err.Error() == "song not found" {
			errorResp := models.NewErrorResponse("Not Found", 404, "Song not found", c.Request.URL.Path)
			c.JSON(http.StatusNotFound, errorResp)
			return
		}
		errorResp := models.NewErrorResponse("Bad Request", 400, "Failed to retrieve song", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	response := models.SongResponse{
		Data: *song,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateSong handles PUT /songs/{id}
// @Summary Update a song by ID
// @Description Song updated successfully
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body models.UpdateSongRequest true "Updated song information"
// @Success 200 {object} models.SongResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /songs/{id} [put]
func (sc *SongController) UpdateSong(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errorResp := models.NewErrorResponse("Bad Request", 400, "Invalid song ID", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	var req models.UpdateSongRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResp := models.NewErrorResponse("Bad Request", 400, "Invalid request body", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	// Get existing song to check if it exists
	existingSong, err := sc.songRepo.GetSongByID(uint(id))
	if err != nil {
		if err.Error() == "song not found" {
			errorResp := models.NewErrorResponse("Not Found", 404, "Song not found", c.Request.URL.Path)
			c.JSON(http.StatusNotFound, errorResp)
			return
		}
		errorResp := models.NewErrorResponse("Bad Request", 400, "Failed to retrieve song", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	// Update song fields
	existingSong.Title = req.Title
	existingSong.Artist = req.Artist

	// Save updated song to database
	if err := sc.songRepo.UpdateSong(existingSong); err != nil {
		errorResp := models.NewErrorResponse("Bad Request", 400, "Failed to update song", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	response := models.SongResponse{
		Data: *existingSong,
	}

	c.JSON(http.StatusOK, response)
}

// DeleteSong handles DELETE /songs/{id}
// @Summary Delete a song by ID
// @Description Song deleted successfully
// @Tags songs
// @Param id path int true "Song ID"
// @Success 204 "Song deleted successfully"
// @Failure 404 {object} models.ErrorResponse
// @Router /songs/{id} [delete]
func (sc *SongController) DeleteSong(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errorResp := models.NewErrorResponse("Bad Request", 400, "Invalid song ID", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	// Delete song from database
	if err := sc.songRepo.DeleteSong(uint(id)); err != nil {
		if err.Error() == "song not found" {
			errorResp := models.NewErrorResponse("Not Found", 404, "Song not found", c.Request.URL.Path)
			c.JSON(http.StatusNotFound, errorResp)
			return
		}
		errorResp := models.NewErrorResponse("Bad Request", 400, "Failed to delete song", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	c.Status(http.StatusNoContent)
}
