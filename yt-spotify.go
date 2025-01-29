package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	"yt-spotify/spotify"
	"yt-spotify/youtube"

	youtubeV3 "google.golang.org/api/youtube/v3"
)

func YouTubeToSpotify() {
	appCtx, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

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
			processYouTubePlaylist(youtubeService, spotifyClient, playlistID, appCtx.PlayListsNameToSave)
		}()
		time.Sleep(500 * time.Millisecond)
	}
	wg.Wait()
}

func processYouTubePlaylist(youtubeService *youtubeV3.Service, spotifyClient *http.Client, playlistID string, playsListName string) {
	playlistItems, err := youtube.FetchPlaylistItems(youtubeService, playlistID)
	if err != nil {
		log.Printf("Unable to fetch YouTube playlist items for %s: %v", playlistID, err)
		return
	}

	spotifyPlaylistID, err := spotify.CheckOrCreatePlaylist(spotifyClient, playsListName)
	if err != nil {
		log.Printf("Unable to find or create Spotify playlist for %s: %v", playlistID, err)
		return
	}

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
