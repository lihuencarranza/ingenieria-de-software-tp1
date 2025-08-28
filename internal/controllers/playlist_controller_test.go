package controllers

// FALLAN TESTS!!!!
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

// Test CreatePlaylist endpoint
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
	assert.NotZero(t, response.Data.ID)
	assert.NotZero(t, response.Data.CreatedAt)
	assert.NotZero(t, response.Data.UpdatedAt)
	assert.NotNil(t, response.Data.Songs)
}

func TestCreatePlaylist_BadRequest_EmptyName(t *testing.T) {
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

	var errorResp models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "Bad Request", errorResp.Title)
	assert.Equal(t, 400, errorResp.Status)
}

func TestCreatePlaylist_BadRequest_EmptyDescription(t *testing.T) {
	router := setupPlaylistTestRouter()

	req := models.CreatePlaylistRequest{
		Name:        "Test Playlist",
		Description: "", // Empty description
	}

	jsonData, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/playlists", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreatePlaylist_BadRequest_InvalidJSON(t *testing.T) {
	router := setupPlaylistTestRouter()

	invalidJSON := `{"name": "Test Playlist", "description": "A test playlist description"` // Missing closing brace

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/playlists", bytes.NewBuffer([]byte(invalidJSON)))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Test GetPlaylists endpoint
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

// Test GetPlaylist endpoint
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
	assert.NotNil(t, response.Data.Songs)
}

func TestGetPlaylist_InvalidID_String(t *testing.T) {
	router := setupPlaylistTestRouter()

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/playlists/invalid", nil)

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errorResp models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "Bad Request", errorResp.Title)
	assert.Equal(t, 400, errorResp.Status)
}

func TestGetPlaylist_InvalidID_Zero(t *testing.T) {
	router := setupPlaylistTestRouter()

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/playlists/0", nil)

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetPlaylist_NotFound(t *testing.T) {
	router := setupPlaylistTestRouter()

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/playlists/999", nil)

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var errorResp models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "Not Found", errorResp.Title)
	assert.Equal(t, 404, errorResp.Status)
}

// Test DeletePlaylist endpoint
func TestDeletePlaylist_Success(t *testing.T) {
	router := setupPlaylistTestRouter()

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("DELETE", "/playlists/1", nil)

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeletePlaylist_InvalidID(t *testing.T) {
	router := setupPlaylistTestRouter()

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("DELETE", "/playlists/invalid", nil)

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeletePlaylist_NotFound(t *testing.T) {
	router := setupPlaylistTestRouter()

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("DELETE", "/playlists/999", nil)

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// Test AddSongToPlaylist endpoint
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

func TestAddSongToPlaylist_InvalidJSON(t *testing.T) {
	router := setupPlaylistTestRouter()

	invalidJSON := `{"songID": 1` // Missing closing brace

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/playlists/1/songs", bytes.NewBuffer([]byte(invalidJSON)))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddSongToPlaylist_PlaylistNotFound(t *testing.T) {
	router := setupPlaylistTestRouter()

	req := models.AddSongToPlaylistRequest{
		SongID: 1,
	}

	jsonData, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/playlists/999/songs", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// Test edge cases
func TestCreatePlaylist_LongName(t *testing.T) {
	router := setupPlaylistTestRouter()

	req := models.CreatePlaylistRequest{
		Name:        "This is a very long playlist name that should still be valid and accepted by the system",
		Description: "A test playlist description",
	}

	jsonData, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/playlists", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreatePlaylist_LongDescription(t *testing.T) {
	router := setupPlaylistTestRouter()

	req := models.CreatePlaylistRequest{
		Name:        "Test Playlist",
		Description: "This is a very long description that should still be valid and accepted by the system. It can contain multiple sentences and should handle long text properly without any issues.",
	}

	jsonData, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/playlists", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreatePlaylist_SpecialCharacters(t *testing.T) {
	router := setupPlaylistTestRouter()

	req := models.CreatePlaylistRequest{
		Name:        "Playlist with special chars: !@#$%^&*()_+-=[]{}|;':\",./<>?",
		Description: "Description with special chars: !@#$%^&*()_+-=[]{}|;':\",./<>?",
	}

	jsonData, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/playlists", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusCreated, w.Code)
}
