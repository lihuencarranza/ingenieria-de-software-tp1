-- Create songs table
CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    artist VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create playlists table
CREATE TABLE IF NOT EXISTS playlists (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_published BOOLEAN DEFAULT true,
    published_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create playlist_songs junction table
CREATE TABLE IF NOT EXISTS playlist_songs (
    id SERIAL PRIMARY KEY,
    playlist_id INTEGER NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,
    song_id INTEGER NOT NULL REFERENCES songs(id) ON DELETE CASCADE,
    added_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(playlist_id, song_id)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_songs_title ON songs(title);
CREATE INDEX IF NOT EXISTS idx_songs_artist ON songs(artist);
CREATE INDEX IF NOT EXISTS idx_playlists_name ON playlists(name);
CREATE INDEX IF NOT EXISTS idx_playlists_published_at ON playlists(published_at);
CREATE INDEX IF NOT EXISTS idx_playlist_songs_playlist_id ON playlist_songs(playlist_id);
CREATE INDEX IF NOT EXISTS idx_playlist_songs_song_id ON playlist_songs(song_id);
CREATE INDEX IF NOT EXISTS idx_playlist_songs_added_at ON playlist_songs(added_at);
