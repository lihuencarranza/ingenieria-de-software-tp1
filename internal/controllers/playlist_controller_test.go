package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"melodia/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupPlaylistTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	playlistController := NewPlaylistController()

	router.POST("/playlists", playlistController.CreatePlaylist)
	router.GET("/playlists", playlistController.GetPlaylists)
	router.GET("/playlists/:id", playlistController.GetPlaylist)
	router.DELETE("/playlists/:id", playlistController.DeletePlaylist)
	router.POST("/playlists/:id/songs", playlistController.AddSongToPlaylist)

	return router
}

func TestCreatePlaylist_Success(t *testing.T) {
	router := setupPlaylistTestRouter()

	req := models.CreatePlaylistRequest{
		Name:        "Test Playlist",
		Description: "A test playlist description",
	}

	jsonData, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/playlists", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.PlaylistResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Test Playlist", response.Data.Name)
	assert.Equal(t, "A test playlist description", response.Data.Description)
	assert.True(t, response.Data.IsPublished)
	assert.NotNil(t, response.Data.PublishedAt)
}

func TestCreatePlaylist_BadRequest(t *testing.T) {
	router := setupPlaylistTestRouter()

	req := models.CreatePlaylistRequest{
		Name:        "", // Empty name
		Description: "A test playlist description",
	}

	jsonData, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/playlists", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetPlaylists_Success(t *testing.T) {
	router := setupPlaylistTestRouter()

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/playlists", nil)

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.PlaylistsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Data)
}

func TestGetPlaylist_Success(t *testing.T) {
	router := setupPlaylistTestRouter()

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/playlists/1", nil)

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.PlaylistResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), response.Data.ID)
	assert.True(t, response.Data.IsPublished)
}

func TestGetPlaylist_InvalidID(t *testing.T) {
	router := setupPlaylistTestRouter()

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/playlists/invalid", nil)

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeletePlaylist_Success(t *testing.T) {
	router := setupPlaylistTestRouter()

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("DELETE", "/playlists/1", nil)

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestAddSongToPlaylist_Success(t *testing.T) {
	router := setupPlaylistTestRouter()

	req := models.AddSongToPlaylistRequest{
		SongID: 1,
	}

	jsonData, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/playlists/1/songs", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.PlaylistResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), response.Data.ID)
}

func TestAddSongToPlaylist_InvalidPlaylistID(t *testing.T) {
	router := setupPlaylistTestRouter()

	req := models.AddSongToPlaylistRequest{
		SongID: 1,
	}

	jsonData, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/playlists/invalid/songs", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
