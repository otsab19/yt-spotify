package config

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"yt-spotify/utils"
)

type AppContext struct {
	YouTubeAPIKey       string
	SpotifyClientID     string
	SpotifyClientSecret string
	SpotifyRedirectURI  string
	Playlists           []string
	PlayListsNameToSave string
	MistralApiKey       string
	ModelToUse          string
}

var appContext *AppContext

func GetConfig() (*AppContext, error) {
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

	var model string
	if os.Getenv("MODEL_TO_USE") == "mistral" {
		model = utils.MISTRAL
	} else if os.Getenv("MODEL_TO_USE") == "ollama" {
		model = utils.OLLAMA
	}

	return &AppContext{
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

func GetAppContext() *AppContext {
	if appContext == nil {
		panic("AppContext is not initialized. Call LoadConfig() first.")
	}
	return appContext
}

func LoadConfig() {
	ctx, err := GetConfig()
	appContext = ctx
	if err != nil {
		fmt.Println("Failed to load config:", err)
		panic("Failed to load config.")
		return
	}

}
