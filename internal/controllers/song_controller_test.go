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

// Test CreateSong endpoint
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
	assert.NotZero(t, response.Data.ID)
	assert.NotZero(t, response.Data.CreatedAt)
	assert.NotZero(t, response.Data.UpdatedAt)
}

func TestCreateSong_BadRequest_EmptyTitle(t *testing.T) {
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

	var errorResp models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "Bad Request", errorResp.Title)
	assert.Equal(t, 400, errorResp.Status)
}

func TestCreateSong_BadRequest_EmptyArtist(t *testing.T) {
	router := setupTestRouter()

	req := models.CreateSongRequest{
		Title:  "Test Song",
		Artist: "", // Empty artist
	}

	jsonData, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/songs", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateSong_BadRequest_InvalidJSON(t *testing.T) {
	router := setupTestRouter()

	invalidJSON := `{"title": "Test Song", "artist": "Test Artist"` // Missing closing brace

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/songs", bytes.NewBuffer([]byte(invalidJSON)))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Test GetSongs endpoint
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

// Test GetSong endpoint
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

func TestGetSong_InvalidID_String(t *testing.T) {
	router := setupTestRouter()

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/songs/invalid", nil)

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errorResp models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "Bad Request", errorResp.Title)
	assert.Equal(t, 400, errorResp.Status)
}

func TestGetSong_InvalidID_Zero(t *testing.T) {
	router := setupTestRouter()

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/songs/0", nil)

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetSong_NotFound(t *testing.T) {
	router := setupTestRouter()

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/songs/999", nil)

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var errorResp models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "Not Found", errorResp.Title)
	assert.Equal(t, 404, errorResp.Status)
}

// Test UpdateSong endpoint
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

func TestUpdateSong_BadRequest_EmptyTitle(t *testing.T) {
	router := setupTestRouter()

	req := models.UpdateSongRequest{
		Title:  "", // Empty title
		Artist: "Updated Artist",
	}

	jsonData, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("PUT", "/songs/1", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateSong_BadRequest_EmptyArtist(t *testing.T) {
	router := setupTestRouter()

	req := models.UpdateSongRequest{
		Title:  "Updated Song",
		Artist: "", // Empty artist
	}

	jsonData, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("PUT", "/songs/1", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateSong_BadRequest_InvalidJSON(t *testing.T) {
	router := setupTestRouter()

	invalidJSON := `{"title": "Updated Song", "artist": "Updated Artist"` // Missing closing brace

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("PUT", "/songs/1", bytes.NewBuffer([]byte(invalidJSON)))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateSong_InvalidID(t *testing.T) {
	router := setupTestRouter()

	req := models.UpdateSongRequest{
		Title:  "Updated Song",
		Artist: "Updated Artist",
	}

	jsonData, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("PUT", "/songs/invalid", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateSong_NotFound(t *testing.T) {
	router := setupTestRouter()

	req := models.UpdateSongRequest{
		Title:  "Updated Song",
		Artist: "Updated Artist",
	}

	jsonData, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("PUT", "/songs/999", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// Test DeleteSong endpoint
func TestDeleteSong_Success(t *testing.T) {
	router := setupTestRouter()

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("DELETE", "/songs/1", nil)

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteSong_InvalidID(t *testing.T) {
	router := setupTestRouter()

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("DELETE", "/songs/invalid", nil)

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteSong_NotFound(t *testing.T) {
	router := setupTestRouter()

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("DELETE", "/songs/999", nil)

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// Test edge cases
func TestCreateSong_LongTitle(t *testing.T) {
	router := setupTestRouter()

	req := models.CreateSongRequest{
		Title:  "This is a very long song title that should still be valid and accepted by the system",
		Artist: "Test Artist",
	}

	jsonData, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/songs", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateSong_SpecialCharacters(t *testing.T) {
	router := setupTestRouter()

	req := models.CreateSongRequest{
		Title:  "Song with special chars: !@#$%^&*()_+-=[]{}|;':\",./<>?",
		Artist: "Artist with special chars: !@#$%^&*()_+-=[]{}|;':\",./<>?",
	}

	jsonData, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/songs", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusCreated, w.Code)
}
