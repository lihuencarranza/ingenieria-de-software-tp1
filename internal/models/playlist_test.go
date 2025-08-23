package models

import (
	"testing"
	"time"
)

func TestPlaylist(t *testing.T) {
	now := time.Now()
	publishedAt := now.Add(-24 * time.Hour)

	playlist := Playlist{
		ID:          1,
		Name:        "My Playlist",
		Description: "A test playlist",
		IsPublished: true,
		PublishedAt: &publishedAt,
		Songs: []PlaylistSong{
			{
				ID:      1,
				Title:   "Song 1",
				Artist:  "Artist 1",
				AddedAt: now,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	if playlist.ID != 1 {
		t.Errorf("Expected ID to be 1, got %d", playlist.ID)
	}

	if playlist.Name != "My Playlist" {
		t.Errorf("Expected Name to be 'My Playlist', got %s", playlist.Name)
	}

	if playlist.Description != "A test playlist" {
		t.Errorf("Expected Description to be 'A test playlist', got %s", playlist.Description)
	}

	if !playlist.IsPublished {
		t.Error("Expected IsPublished to be true")
	}

	if playlist.PublishedAt.IsZero() {
		t.Error("Expected PublishedAt to be set")
	}

	if len(playlist.Songs) != 1 {
		t.Errorf("Expected 1 song, got %d", len(playlist.Songs))
	}
}

func TestCreatePlaylistRequest(t *testing.T) {
	req := CreatePlaylistRequest{
		Name:        "New Playlist",
		Description: "A new playlist description",
	}

	if req.Name != "New Playlist" {
		t.Errorf("Expected Name to be 'New Playlist', got %s", req.Name)
	}

	if req.Description != "A new playlist description" {
		t.Errorf("Expected Description to be 'A new playlist description', got %s", req.Description)
	}
}

func TestAddSongToPlaylistRequest(t *testing.T) {
	req := AddSongToPlaylistRequest{
		SongID: 1,
	}

	if req.SongID != 1 {
		t.Errorf("Expected SongID to be 1, got %d", req.SongID)
	}
}

func TestPlaylistResponse(t *testing.T) {
	now := time.Now()
	playlist := Playlist{
		ID:          1,
		Name:        "Test Playlist",
		Description: "Test Description",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	response := PlaylistResponse{
		Data: playlist,
	}

	if response.Data.ID != playlist.ID {
		t.Errorf("Expected response Data.ID to be %d, got %d", playlist.ID, response.Data.ID)
	}
}

func TestPlaylistsResponse(t *testing.T) {
	now := time.Now()
	playlists := []Playlist{
		{
			ID:          1,
			Name:        "Playlist 1",
			Description: "Description 1",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          2,
			Name:        "Playlist 2",
			Description: "Description 2",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	response := PlaylistsResponse{
		Data: playlists,
	}

	if len(response.Data) != 2 {
		t.Errorf("Expected response Data length to be 2, got %d", len(response.Data))
	}
}
