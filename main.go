package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"yt-spotify/spotify"
	"yt-spotify/youtube"
)

type AppContext struct {
	YouTubeAPIKey       string
	SpotifyClientID     string
	SpotifyClientSecret string
	SpotifyRedirectURI  string
	Playlists           []string
}

func loadEnv() (*AppContext, error) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	// Parse the PLAYLISTS environment variable (JSON array)
	var playlists []string
	rawPlaylists := os.Getenv("PLAYLISTS")
	if rawPlaylists != "" {
		err := json.Unmarshal([]byte(rawPlaylists), &playlists)
		if err != nil {
			return nil, fmt.Errorf("error parsing PLAYLISTS environment variable: %w", err)
		}
	}

	return &AppContext{
		YouTubeAPIKey:       os.Getenv("YOUTUBE_API_KEY"),
		SpotifyClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		SpotifyClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		SpotifyRedirectURI:  os.Getenv("SPOTIFY_REDIRECT_URI"),
		Playlists:           playlists,
	}, nil
}

func main() {
	// Load environment variables into context
	appCtx, err := loadEnv()
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}
	var playlistID string
	if len(appCtx.Playlists) == 0 {
		// Get YouTube playlist ID from user
		fmt.Print("Enter YouTube Playlist ID: ")
		fmt.Scanln(&playlistID)
	} else {
		//todo: search all
		playlistID = appCtx.Playlists[0]
	}

	// Fetch YouTube playlist items
	youtubeService, err := youtube.NewService(appCtx.YouTubeAPIKey)
	if err != nil {
		log.Fatalf("Unable to create YouTube service: %v", err)
	}

	playlistItems, err := youtube.FetchPlaylistItems(youtubeService, playlistID)
	if err != nil {
		log.Fatalf("Unable to fetch YouTube playlist items: %v", err)
	}

	// Authenticate with Spotify
	spotifyClient, err := spotify.Authenticate(appCtx.SpotifyClientID, appCtx.SpotifyClientSecret, appCtx.SpotifyRedirectURI)
	if err != nil {
		log.Fatalf("Unable to authenticate with Spotify: %v", err)
	}

	spotifyPlaylistID, err := spotify.CreatePlaylist(spotifyClient, "My YouTube Playlist")
	if err != nil {
		log.Fatalf("Unable to create Spotify playlist: %v", err)
	}

	// Add tracks to Spotify playlist
	for _, item := range playlistItems {
		trackName := item.Snippet.Title
		artistName := item.Snippet.VideoOwnerChannelTitle

		trackID, err := spotify.SearchTrack(spotifyClient, trackName, artistName)
		if err != nil {
			log.Printf("Unable to find track '%s' by '%s' on Spotify: %v", trackName, artistName, err)
			continue
		}

		err = spotify.AddTrackToPlaylist(spotifyClient, spotifyPlaylistID, trackID)
		if err != nil {
			log.Printf("Unable to add track '%s' to Spotify playlist: %v", trackName, err)
			continue
		}

		fmt.Printf("Added '%s' by '%s' to Spotify playlist\n", trackName, artistName)
	}
}
