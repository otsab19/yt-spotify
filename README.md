# YouTube to Spotify Playlist Converter

A Go-based project that allows users to import YouTube playlists or a list of songs into Spotify playlists. This tool fetches tracks from a YouTube playlist or a provided song list, searches for the corresponding tracks on Spotify, and creates a Spotify playlist with the matched tracks.

## Features

- Fetches tracks from a public YouTube playlist.
- Allows importing songs from a text file (`songs.txt`) placed in the `inputFiles` folder.
- Cleans and matches track titles and artist names to improve Spotify search results.
- Creates a new playlist on Spotify and populates it with matched tracks.
- Handles mismatched or missing tracks gracefully with error logging.

---

## Prerequisites

### 1. **Spotify Developer Account**
- Register an app in the [Spotify Developer Dashboard](https://developer.spotify.com/dashboard/).
- Obtain the `Client ID`, `Client Secret`, and set a redirect URI for authentication.

### 2. **Google Cloud Project**
- Enable the **YouTube Data API v3** in the [Google Cloud Console](https://console.cloud.google.com/).
- Obtain the API Key.

### 3. **Go Environment**
- Install [Go](https://golang.org/dl/) if not already installed.

---

## Environment Variables

Create a `.env` file in the root of the project and populate it with the following variables:

```plaintext
YOUTUBE_API_KEY=your_youtube_api_key
SPOTIFY_CLIENT_ID=your_spotify_client_id
SPOTIFY_CLIENT_SECRET=your_spotify_client_secret
SPOTIFY_REDIRECT_URI=your_redirect_uri
PLAYLISTS=["your_playlist_id_1", "your_playlist_id_2"]
PLAYLIST_NAME_TO_SAVE=name
```

---

## Input Methods

### 1. **YouTube Playlist Import**
- Provide YouTube playlist IDs in the `PLAYLISTS` environment variable.
- The tool will fetch the playlist tracks and attempt to match them on Spotify.

### 2. **Text File Import (Songs List)**
- Create a folder named `inputFiles` in the root directory.
- Add a text file (`songs.txt`) containing song names and artist names, one per line.
- Multiple files are supported if the filename starts with `songs` (e.g., `songs1.txt`, `songs2.txt`).
- The tool will process these files and search for corresponding tracks on Spotify.

---

## Usage

1. Clone the repository and navigate to the project folder.
2. Set up the `.env` file with your API keys and configurations.
3. Place song lists in the `inputFiles` folder (if using the text file method).
4. Run the Go program:
   ```sh
   go run main.go
   ```
5. The matched tracks will be added to a new Spotify playlist.

---



## Ollama and Mistral AI Integration for Song and Artist Name Extraction
This project now supports Ollama, a local AI-powered service, to improve track matching accuracy by extracting clean song titles and artist names.

### How It Works
The tool utilizes Ollama to intelligently extract the correct song title and artist.
Ensures better search results when querying the Spotify API.
### Setup for Ollama
Ensure that Ollama is installed and running on your local machine. If it's running it will take the LLM in use as consideration.
### Setup for Mistral AI
Provide API key in env.

Set model to use either mistrlai or ollama

---

## Error Handling

- Logs are generated for any tracks that couldn't be matched.
- Mismatched or missing tracks are gracefully handled and reported in the console.

---

## Future Improvements

- Improve track matching accuracy with fuzzy search.
- Support for private YouTube playlists via OAuth authentication.
- Enhance error handling and logging.

