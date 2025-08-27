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

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	songController := NewSongController()

	router.POST("/songs", songController.CreateSong)
	router.GET("/songs", songController.GetSongs)
	router.GET("/songs/:id", songController.GetSong)
	router.PUT("/songs/:id", songController.UpdateSong)
	router.DELETE("/songs/:id", songController.DeleteSong)

	return router
}

func TestCreateSong_Success(t *testing.T) {
	router := setupTestRouter()

	req := models.CreateSongRequest{
		Title:  "Test Song",
		Artist: "Test Artist",
	}

	jsonData, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/songs", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.SongResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Test Song", response.Data.Title)
	assert.Equal(t, "Test Artist", response.Data.Artist)
}

func TestCreateSong_BadRequest(t *testing.T) {
	router := setupTestRouter()

	req := models.CreateSongRequest{
		Title:  "", // Empty title
		Artist: "Test Artist",
	}

	jsonData, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/songs", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetSongs_Success(t *testing.T) {
	router := setupTestRouter()

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/songs", nil)

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.SongsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Data)
}

func TestGetSong_Success(t *testing.T) {
	router := setupTestRouter()

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/songs/1", nil)

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.SongResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), response.Data.ID)
}

func TestGetSong_InvalidID(t *testing.T) {
	router := setupTestRouter()

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/songs/invalid", nil)

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateSong_Success(t *testing.T) {
	router := setupTestRouter()

	req := models.UpdateSongRequest{
		Title:  "Updated Song",
		Artist: "Updated Artist",
	}

	jsonData, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("PUT", "/songs/1", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.SongResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Song", response.Data.Title)
	assert.Equal(t, "Updated Artist", response.Data.Artist)
}

func TestDeleteSong_Success(t *testing.T) {
	router := setupTestRouter()

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("DELETE", "/songs/1", nil)

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
