package models

import (
	"testing"
	"time"
)

func TestSong(t *testing.T) {
	now := time.Now()

	song := Song{
		ID:        1,
		Title:     "Test Song",
		Artist:    "Test Artist",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if song.ID != 1 {
		t.Errorf("Expected ID to be 1, got %d", song.ID)
	}

	if song.Title != "Test Song" {
		t.Errorf("Expected Title to be 'Test Song', got %s", song.Title)
	}

	if song.Artist != "Test Artist" {
		t.Errorf("Expected Artist to be 'Test Artist', got %s", song.Artist)
	}
}

func TestCreateSongRequest(t *testing.T) {
	req := CreateSongRequest{
		Title:  "New Song",
		Artist: "New Artist",
	}

	if req.Title != "New Song" {
		t.Errorf("Expected Title to be 'New Song', got %s", req.Title)
	}

	if req.Artist != "New Artist" {
		t.Errorf("Expected Artist to be 'New Artist', got %s", req.Artist)
	}
}

func TestUpdateSongRequest(t *testing.T) {
	req := UpdateSongRequest{
		Title:  "Updated Song",
		Artist: "Updated Artist",
	}

	if req.Title != "Updated Song" {
		t.Errorf("Expected Title to be 'Updated Song', got %s", req.Title)
	}

	if req.Artist != "Updated Artist" {
		t.Errorf("Expected Artist to be 'Updated Artist', got %s", req.Artist)
	}
}

func TestSongResponse(t *testing.T) {
	now := time.Now()
	song := Song{
		ID:        1,
		Title:     "Test Song",
		Artist:    "Test Artist",
		CreatedAt: now,
		UpdatedAt: now,
	}

	response := SongResponse{
		Data: song,
	}

	if response.Data.ID != song.ID {
		t.Errorf("Expected response Data ID to be %d, got %d", song.ID, response.Data.ID)
	}
}

func TestSongsResponse(t *testing.T) {
	now := time.Now()
	songs := []Song{
		{
			ID:        1,
			Title:     "Test Song 1",
			Artist:    "Test Artist 1",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        2,
			Title:     "Test Song 2",
			Artist:    "Test Artist 2",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	response := SongsResponse{
		Data: songs,
	}

	if len(response.Data) != 2 {
		t.Errorf("Expected response Data length to be 2, got %d", len(response.Data))
	}
}
