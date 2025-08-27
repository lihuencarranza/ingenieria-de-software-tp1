package repositories

import (
	"database/sql"
	"fmt"
	"melodia/internal/database"
	"melodia/internal/models"
	"time"
)

// SongRepository handles database operations for songs
type SongRepository struct {
	db *sql.DB
}

// NewSongRepository creates a new song repository
func NewSongRepository() *SongRepository {
	return &SongRepository{
		db: database.DB,
	}
}

// CreateSong creates a new song in the database
func (r *SongRepository) CreateSong(song *models.Song) error {
	query := `
		INSERT INTO songs (title, artist, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	err := r.db.QueryRow(query, song.Title, song.Artist, now, now).
		Scan(&song.ID, &song.CreatedAt, &song.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error creating song: %v", err)
	}

	return nil
}

// GetSongs retrieves all songs from the database
func (r *SongRepository) GetSongs() ([]models.Song, error) {
	query := `SELECT id, title, artist, created_at, updated_at FROM songs ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying songs: %v", err)
	}
	defer rows.Close()

	var songs []models.Song
	for rows.Next() {
		var song models.Song
		err := rows.Scan(&song.ID, &song.Title, &song.Artist, &song.CreatedAt, &song.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning song: %v", err)
		}
		songs = append(songs, song)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating songs: %v", err)
	}

	return songs, nil
}

// GetSongByID retrieves a song by its ID
func (r *SongRepository) GetSongByID(id uint) (*models.Song, error) {
	query := `SELECT id, title, artist, created_at, updated_at FROM songs WHERE id = $1`

	var song models.Song
	err := r.db.QueryRow(query, id).
		Scan(&song.ID, &song.Title, &song.Artist, &song.CreatedAt, &song.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("song not found")
		}
		return nil, fmt.Errorf("error querying song: %v", err)
	}

	return &song, nil
}

// UpdateSong updates an existing song in the database
func (r *SongRepository) UpdateSong(song *models.Song) error {
	query := `
		UPDATE songs 
		SET title = $1, artist = $2, updated_at = $3
		WHERE id = $4
		RETURNING created_at, updated_at
	`

	now := time.Now()
	err := r.db.QueryRow(query, song.Title, song.Artist, now, song.ID).
		Scan(&song.CreatedAt, &song.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("song not found")
		}
		return fmt.Errorf("error updating song: %v", err)
	}

	song.UpdatedAt = now
	return nil
}

// DeleteSong deletes a song from the database
func (r *SongRepository) DeleteSong(id uint) error {
	query := `DELETE FROM songs WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting song: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("song not found")
	}

	return nil
}
