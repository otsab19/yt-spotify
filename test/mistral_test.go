package test

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
	"yt-spotify/config"
	"yt-spotify/service"
	"yt-spotify/utils"

	"github.com/stretchr/testify/assert"
)

func GetConfig() (*config.AppContext, error) {
	err := godotenv.Load("../.env")
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

	var model string
	if os.Getenv("MODEL_TO_USE") == "mistral" {
		model = utils.MISTRAL
	} else if os.Getenv("MODEL_TO_USE") == "ollama" {
		model = utils.OLLAMA
	}

	return &config.AppContext{
		YouTubeAPIKey:       os.Getenv("YOUTUBE_API_KEY"),
		SpotifyClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		SpotifyClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		SpotifyRedirectURI:  os.Getenv("SPOTIFY_REDIRECT_URI"),
		PlayListsNameToSave: playListsName,
		Playlists:           playlists,
		MistralApiKey:       os.Getenv("MISTRAL_API_KEY"),
		ModelToUse:          model,
	}, nil
}

// Test Mistral API with real request
func TestExtractSongArtist_RealMistralAPI(t *testing.T) {
	appCtx, err := GetConfig()
	if err != nil {
		return
	}

	// Ensure MISTRAL_API_KEY is set
	if appCtx.MistralApiKey == "" {
		t.Fatal("MISTRAL_API_KEY is not set in the environment")
	}

	// Initialize real Mistral service
	mistralService, err := service.NewMistralService(appCtx)
	if err != nil {
		t.Fatalf("Failed to initialize Mistral service: %v", err)
	}

	// Call the real API
	song, artist, err := mistralService.ExtractSongArtist("The Weeknd - Blinding Lights (Official Video)")

	// Assertions
	assert.NoError(t, err, "Mistral API call should not fail")
	assert.NotEmpty(t, song, "Extracted song name should not be empty")
	assert.NotEmpty(t, artist, "Extracted artist name should not be empty")

	// Print results
	log.Printf("Extracted Song: %s, Artist: %s\n", song, artist)
}

// Test real Mistral API with an unusual input
func TestExtractSongArtist_UnusualInput(t *testing.T) {
	appCtx, err := GetConfig()
	if err != nil {
		return
	}
	// Ensure MISTRAL_API_KEY is set
	if appCtx.MistralApiKey == "" {
		t.Fatal("MISTRAL_API_KEY is not set in the environment")
	}

	// Initialize real Mistral service
	mistralService, err := service.NewMistralService(appCtx)
	if err != nil {
		t.Fatalf("Failed to initialize Mistral service: %v", err)
	}

	// Call the real API with an unusual input
	song, artist, err := mistralService.ExtractSongArtist("This is not a song - Just testing (Live 2025)")

	// Assertions
	assert.NoError(t, err, "Mistral API call should not fail")
	assert.NotEmpty(t, song, "Extracted song name should not be empty")
	assert.NotEmpty(t, artist, "Extracted artist name should not be empty")

	// Print results
	log.Printf("Extracted Song: %s, Artist: %s\n", song, artist)
}
