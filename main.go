package main

import (
	"fmt"
	"os"
	"yt-spotify/config"
)

func main() {
	//init config
	config.LoadConfig()
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "yt-spotify":
			YouTubeToSpotify()
		case "songs-spotify":
			SongsToSpotify()
		default:
			fmt.Println("Invalid argument. Use 'yt-spotify' or 'songs-spotify'.")
		}
	} else {
		fmt.Println("Choose an option: \n1. Convert YouTube playlist to Spotify (yt-to-spotify)\n2. Convert songs from list to Spotify (songs-to-spotify)")
		var choice int
		fmt.Scanln(&choice)
		switch choice {
		case 1:
			YouTubeToSpotify()
		case 2:
			SongsToSpotify()
		default:
			fmt.Println("Invalid choice")
		}
	}
}
