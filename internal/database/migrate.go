package database

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations runs database migrations
func RunMigrations() error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	driver, err := postgres.WithInstance(DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("error creating postgres driver: %v", err)
	}

	// Get current working directory and construct migrations path
	workDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working directory: %v", err)
	}

	migrationsPath := fmt.Sprintf("file://%s/internal/database/migrations", workDir)
	log.Printf("Running migrations from: %s", migrationsPath)

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		return fmt.Errorf("error creating migrate instance: %v", err)
	}
	defer m.Close()

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error running migrations: %v", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// CreateTablesIfNotExist creates tables if they don't exist (fallback method)
func CreateTablesIfNotExist() error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	// Create songs table
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS songs (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			artist VARCHAR(255) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating songs table: %v", err)
	}

	// Create playlists table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS playlists (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			is_published BOOLEAN DEFAULT true,
			published_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating playlists table: %v", err)
	}

	// Create playlist_songs junction table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS playlist_songs (
			id SERIAL PRIMARY KEY,
			playlist_id INTEGER NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,
			song_id INTEGER NOT NULL REFERENCES songs(id) ON DELETE CASCADE,
			added_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(playlist_id, song_id)
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating playlist_songs table: %v", err)
	}

	log.Println("Database tables created successfully")
	return nil
}
