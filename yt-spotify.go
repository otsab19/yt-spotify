package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	"yt-spotify/config"
	"yt-spotify/service"
	"yt-spotify/spotify"
	"yt-spotify/utils"
	"yt-spotify/youtube"

	youtubeV3 "google.golang.org/api/youtube/v3"
)

func YouTubeToSpotify() {
	appCtx := config.GetAppContext()
	if len(appCtx.Playlists) == 0 {
		fmt.Print("Enter YouTube Playlist ID: ")
		var playlistID string
		fmt.Scanln(&playlistID)
		appCtx.Playlists = append(appCtx.Playlists, playlistID)
	}

	youtubeService, err := youtube.NewService(appCtx.YouTubeAPIKey)
	if err != nil {
		log.Fatalf("Unable to create YouTube service: %v", err)
	}

	spotifyClient, err := spotify.Authenticate(appCtx.SpotifyClientID, appCtx.SpotifyClientSecret, appCtx.SpotifyRedirectURI)
	if err != nil {
		log.Fatalf("Unable to authenticate with Spotify: %v", err)
	}

	var wg sync.WaitGroup
	for _, playlistID := range appCtx.Playlists {
		wg.Add(1)
		playlistID := playlistID
		go func() {
			defer wg.Done()
			processYouTubePlaylist(youtubeService, spotifyClient, playlistID, appCtx)
		}()
		time.Sleep(500 * time.Millisecond)
	}
	wg.Wait()
}

func processYouTubePlaylist(youtubeService *youtubeV3.Service, spotifyClient *http.Client, playlistID string, appCtx *config.AppContext) {
	playlistItems, err := youtube.FetchPlaylistItems(youtubeService, playlistID)
	if err != nil {
		log.Printf("Unable to fetch YouTube playlist items for %s: %v", playlistID, err)
		return
	}

	spotifyPlaylistID, err := spotify.CheckOrCreatePlaylist(spotifyClient, appCtx.PlayListsNameToSave)
	if err != nil {
		log.Printf("Unable to find or create Spotify playlist for %s: %v", playlistID, err)
		return
	}

	var aiService service.AiService

	switch appCtx.ModelToUse {
	case utils.MISTRAL:
		mistralService, err := service.NewMistralService(appCtx)
		if err != nil {
			log.Printf("Error initializing Mistral Service: %v", err)
		} else {
			aiService = mistralService
		}
	case utils.OLLAMA:
		ollamaService := service.NewOllamaService()
		if ollamaService.IsOllamaAvailable() {
			aiService = ollamaService
		} else {
			log.Println("Ollama API is not running. Falling back to raw metadata.")
		}
	default:
		log.Println("No valid AI model selected. Using raw metadata.")
	}

	for _, item := range playlistItems {
		trackName := item.Snippet.Title
		artistName := item.Snippet.VideoOwnerChannelTitle

		// Use LLM
		if aiService != nil {
			extractedTrack, extractedArtist, err := aiService.ExtractSongArtist(trackName)
			if err == nil {
				trackName = extractedTrack
				artistName = extractedArtist
			} else {
				log.Printf("AI extraction failed, using default metadata: %v", err)
			}
		}
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
