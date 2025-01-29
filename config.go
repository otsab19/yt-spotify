package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type AppContext struct {
	YouTubeAPIKey       string
	SpotifyClientID     string
	SpotifyClientSecret string
	SpotifyRedirectURI  string
	Playlists           []string
	PlayListsNameToSave string
}

func LoadConfig() (*AppContext, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	var playlists []string
	rawPlaylists := os.Getenv("PLAYLISTS")
	if rawPlaylists != "" {
		err := json.Unmarshal([]byte(rawPlaylists), &playlists)
		if err != nil {
			return nil, fmt.Errorf("error parsing PLAYLISTS environment variable: %w", err)
		}
	}
	var playListsName = os.Getenv("PLAYLIST_NAME_TO_SAVE")
	if playListsName == "" {
		playListsName = "Playlist"
	}

	return &AppContext{
		YouTubeAPIKey:       os.Getenv("YOUTUBE_API_KEY"),
		SpotifyClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		SpotifyClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		SpotifyRedirectURI:  os.Getenv("SPOTIFY_REDIRECT_URI"),
		PlayListsNameToSave: playListsName,
		Playlists:           playlists,
	}, nil
}
