package repositories

import (
	"database/sql"
	"fmt"
	"melodia/internal/database"
	"melodia/internal/models"
	"time"
)

// PlaylistRepository handles database operations for playlists
type PlaylistRepository struct {
	db *sql.DB
}

// NewPlaylistRepository creates a new playlist repository
func NewPlaylistRepository() *PlaylistRepository {
	return &PlaylistRepository{
		db: database.DB,
	}
}

// CreatePlaylist creates a new playlist in the database
func (r *PlaylistRepository) CreatePlaylist(playlist *models.Playlist) error {
	query := `
		INSERT INTO playlists (name, description, is_published, published_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	err := r.db.QueryRow(query,
		playlist.Name,
		playlist.Description,
		playlist.IsPublished,
		playlist.PublishedAt,
		now,
		now,
	).Scan(&playlist.ID, &playlist.CreatedAt, &playlist.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error creating playlist: %v", err)
	}

	return nil
}

// GetPlaylists retrieves playlists with optional published filter
func (r *PlaylistRepository) GetPlaylists(published *bool) ([]models.Playlist, error) {
	var query string

	if published == nil || *published {
		// Default: only published playlists, ordered by publishedAt desc
		query = `
			SELECT id, name, description, is_published, published_at, created_at, updated_at 
			FROM playlists 
			WHERE is_published = true 
			ORDER BY published_at DESC
		`
	} else {
		// All playlists, ordered by created_at desc (most recent first)
		query = `
			SELECT id, name, description, is_published, published_at, created_at, updated_at 
			FROM playlists 
			ORDER BY created_at DESC
		`
	}

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying playlists: %v", err)
	}
	defer rows.Close()

	var playlists []models.Playlist
	for rows.Next() {
		var playlist models.Playlist
		err := rows.Scan(
			&playlist.ID,
			&playlist.Name,
			&playlist.Description,
			&playlist.IsPublished,
			&playlist.PublishedAt,
			&playlist.CreatedAt,
			&playlist.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning playlist: %v", err)
		}

		// Load songs for this playlist
		songsQuery := `
			SELECT s.id, s.title, s.artist, ps.added_at
			FROM playlist_songs ps
			JOIN songs s ON ps.song_id = s.id
			WHERE ps.playlist_id = $1
			ORDER BY ps.added_at DESC
		`

		songRows, err := r.db.Query(songsQuery, playlist.ID)
		if err != nil {
			return nil, fmt.Errorf("error querying playlist songs: %v", err)
		}

		var songs []models.PlaylistSong
		for songRows.Next() {
			var song models.PlaylistSong
			err := songRows.Scan(&song.ID, &song.Title, &song.Artist, &song.AddedAt)
			if err != nil {
				songRows.Close()
				return nil, fmt.Errorf("error scanning playlist song: %v", err)
			}
			songs = append(songs, song)
		}

		songRows.Close()
		if err = songRows.Err(); err != nil {
			return nil, fmt.Errorf("error iterating playlist songs: %v", err)
		}

		playlist.Songs = songs
		playlists = append(playlists, playlist)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating playlists: %v", err)
	}

	return playlists, nil
}

// GetPlaylistByID retrieves a playlist by its ID with songs ordered by addedAt desc
func (r *PlaylistRepository) GetPlaylistByID(id uint) (*models.Playlist, error) {
	// First get the playlist
	playlistQuery := `
		SELECT id, name, description, is_published, published_at, created_at, updated_at 
		FROM playlists 
		WHERE id = $1
	`

	var playlist models.Playlist
	err := r.db.QueryRow(playlistQuery, id).Scan(
		&playlist.ID,
		&playlist.Name,
		&playlist.Description,
		&playlist.IsPublished,
		&playlist.PublishedAt,
		&playlist.CreatedAt,
		&playlist.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("playlist not found")
		}
		return nil, fmt.Errorf("error querying playlist: %v", err)
	}

	// Then get the songs for this playlist
	songsQuery := `
		SELECT s.id, s.title, s.artist, ps.added_at
		FROM playlist_songs ps
		JOIN songs s ON ps.song_id = s.id
		WHERE ps.playlist_id = $1
		ORDER BY ps.added_at DESC
	`

	songRows, err := r.db.Query(songsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("error querying playlist songs: %v", err)
	}
	defer songRows.Close()

	var songs []models.PlaylistSong
	for songRows.Next() {
		var song models.PlaylistSong
		err := songRows.Scan(&song.ID, &song.Title, &song.Artist, &song.AddedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning playlist song: %v", err)
		}
		songs = append(songs, song)
	}

	if err = songRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating playlist songs: %v", err)
	}

	playlist.Songs = songs
	return &playlist, nil
}

// DeletePlaylist deletes a playlist from the database
func (r *PlaylistRepository) DeletePlaylist(id uint) error {
	query := `DELETE FROM playlists WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting playlist: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("playlist not found")
	}

	return nil
}

// AddSongToPlaylist adds a song to a playlist
func (r *PlaylistRepository) AddSongToPlaylist(playlistID, songID uint) error {
	// First check if the song exists
	songQuery := `SELECT id FROM songs WHERE id = $1`
	var songExists uint
	err := r.db.QueryRow(songQuery, songID).Scan(&songExists)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("song not found")
		}
		return fmt.Errorf("error checking song: %v", err)
	}

	// Then check if the playlist exists
	playlistQuery := `SELECT id FROM playlists WHERE id = $1`
	var playlistExists uint
	err = r.db.QueryRow(playlistQuery, playlistID).Scan(&playlistExists)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("playlist not found")
		}
		return fmt.Errorf("error checking playlist: %v", err)
	}

	// Add the song to the playlist
	insertQuery := `
		INSERT INTO playlist_songs (playlist_id, song_id, added_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (playlist_id, song_id) DO NOTHING
	`

	now := time.Now()
	_, err = r.db.Exec(insertQuery, playlistID, songID, now)
	if err != nil {
		return fmt.Errorf("error adding song to playlist: %v", err)
	}

	return nil
}

// PublishPlaylist publishes a playlist by setting isPublished=true and publishedAt=now()
func (r *PlaylistRepository) PublishPlaylist(id uint) error {
	query := `
		UPDATE playlists 
		SET is_published = true, published_at = $1, updated_at = $1
		WHERE id = $2
	`

	now := time.Now()
	result, err := r.db.Exec(query, now, id)
	if err != nil {
		return fmt.Errorf("error publishing playlist: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("playlist not found")
	}

	return nil
}
