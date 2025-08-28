package models

import "time"

// Playlist represents a playlist in the system
type Playlist struct {
	ID          uint           `json:"id" db:"id"`
	Name        string         `json:"name" db:"name"`
	Description string         `json:"description" db:"description"`
	IsPublished bool           `json:"is_published" db:"is_published"`
	PublishedAt *time.Time     `json:"published_at,omitempty" db:"published_at"`
	Songs       []PlaylistSong `json:"songs" db:"-"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" db:"updated_at"`
}

// PlaylistSong represents a song within a playlist
type PlaylistSong struct {
	ID      uint      `json:"id" db:"id"`
	Title   string    `json:"title" db:"title"`
	Artist  string    `json:"artist" db:"artist"`
	AddedAt time.Time `json:"added_at" db:"added_at"`
}

// CreatePlaylistRequest represents the request to create a playlist
type CreatePlaylistRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required,min=50,max=255"`
}

// PublishPlaylistRequest represents the request to publish a playlist
type PublishPlaylistRequest struct {
	// Empty struct as this endpoint doesn't require body parameters
}

// AddSongToPlaylistRequest represents the request to add a song to a playlist
type AddSongToPlaylistRequest struct {
	SongID uint `json:"song_id" binding:"required"`
}

// PlaylistResponse represents the response for playlist operations
type PlaylistResponse struct {
	Data Playlist `json:"data"`
}

// PlaylistsResponse represents the response for multiple playlists
type PlaylistsResponse struct {
	Data []Playlist `json:"data"`
}
