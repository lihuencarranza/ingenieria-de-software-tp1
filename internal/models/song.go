package models

import "time"

// Song represents a song in the system
type Song struct {
	ID        uint      `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Artist    string    `json:"artist" db:"artist"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateSongRequest represents the request to create a song
type CreateSongRequest struct {
	Title  string `json:"title" binding:"required"`
	Artist string `json:"artist" binding:"required"`
}

// UpdateSongRequest represents the request to update a song
type UpdateSongRequest struct {
	Title  string `json:"title" binding:"required"`
	Artist string `json:"artist" binding:"required"`
}

// SongResponse represents the response for song operations
type SongResponse struct {
	Data Song `json:"data"`
}

// SongsResponse represents the response for multiple songs
type SongsResponse struct {
	Data []Song `json:"data"`
}
