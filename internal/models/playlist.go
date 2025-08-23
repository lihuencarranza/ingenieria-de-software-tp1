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

