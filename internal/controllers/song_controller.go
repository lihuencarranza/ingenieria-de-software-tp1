package controllers

import (
	"net/http"
	"strconv"
	"time"

	"melodia/internal/models"

	"github.com/gin-gonic/gin"
)

// SongController handles song-related HTTP requests
type SongController struct {
	// TODO: Add song service/repository dependency
}

// NewSongController creates a new song controller
func NewSongController() *SongController {
	return &SongController{}
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

	// TODO: Validate request
	if req.Title == "" || req.Artist == "" {
		errorResp := models.NewErrorResponse("Bad Request", 400, "Title and artist are required", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	// TODO: Save to database
	now := time.Now()
	song := models.Song{
		ID:        1, // TODO: Generate proper ID
		Title:     req.Title,
		Artist:    req.Artist,
		CreatedAt: now,
		UpdatedAt: now,
	}

	response := models.SongResponse{
		Data: song,
	}

	c.JSON(http.StatusCreated, response)
}

// GetSongs handles GET /songs
// @Summary Retrieve all songs
// @Description Get a list of all songs
// @Tags songs
// @Produce json
// @Success 200 {object} models.SongsResponse
// @Router /songs [get]
func (sc *SongController) GetSongs(c *gin.Context) {
	// TODO: Get from database
	songs := []models.Song{} // Empty for now

	response := models.SongsResponse{
		Data: songs,
	}

	c.JSON(http.StatusOK, response)
}

// GetSong handles GET /songs/{id}
// @Summary Retrieve a song by ID
// @Description Get a specific song by its ID
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

	// TODO: Get from database by ID
	// For now, return a mock song
	song := models.Song{
		ID:        uint(id),
		Title:     "Sample Song",
		Artist:    "Sample Artist",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	response := models.SongResponse{
		Data: song,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateSong handles PUT /songs/{id}
// @Summary Update a song by ID
// @Description Update an existing song's information
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

	// TODO: Update in database
	now := time.Now()
	song := models.Song{
		ID:        uint(id),
		Title:     req.Title,
		Artist:    req.Artist,
		CreatedAt: now,
		UpdatedAt: now,
	}

	response := models.SongResponse{
		Data: song,
	}

	c.JSON(http.StatusOK, response)
}

// DeleteSong handles DELETE /songs/{id}
// @Summary Delete a song by ID
// @Description Delete a specific song by its ID
// @Tags songs
// @Param id path int true "Song ID"
// @Success 204 "No Content"
// @Failure 404 {object} models.ErrorResponse
// @Router /songs/{id} [delete]
func (sc *SongController) DeleteSong(c *gin.Context) {
	idStr := c.Param("id")
	_, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errorResp := models.NewErrorResponse("Bad Request", 400, "Invalid song ID", c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	// TODO: Delete from database
	c.Status(http.StatusNoContent)
}
