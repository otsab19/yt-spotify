package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"yt-spotify/spotify"
)

func SongsToSpotify() {
	appCtx, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	files, err := os.ReadDir("inputFiles")
	if err != nil {
		log.Fatalf("Unable to read input directory: %v", err)
	}

	var songLines []string
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "songs") {
			data, err := os.ReadFile("inputFiles/" + file.Name())
			if err != nil {
				log.Printf("Unable to read file %s: %v", file.Name(), err)
				continue
			}
			songLines = append(songLines, strings.Split(string(data), "\n")...)
		}
	}

	spotifyClient, err := spotify.Authenticate(appCtx.SpotifyClientID, appCtx.SpotifyClientSecret, appCtx.SpotifyRedirectURI)
	if err != nil {
		log.Fatalf("Unable to authenticate with Spotify: %v", err)
	}

	spotifyPlaylistID, err := spotify.CheckOrCreatePlaylist(spotifyClient, appCtx.PlayListsNameToSave)
	if err != nil {
		log.Fatalf("Unable to create Spotify playlist: %v", err)
	}

	for _, line := range songLines {
		trackID, err := spotify.SearchTrack(spotifyClient, string(line), "")
		if err != nil {
			log.Printf("Unable to find track '%s' on Spotify: %v", line, err)
			continue
		}
		err = spotify.AddTrackToPlaylist(spotifyClient, spotifyPlaylistID, trackID)
		if err != nil {
			log.Printf("Unable to add track '%s' to Spotify playlist: %v", line, err)
			continue
		}
		fmt.Printf("Added '%s' to Spotify playlist\n", line)
	}
}
