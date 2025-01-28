# YouTube to Spotify Playlist Converter

A Go-based project that allows users to import YouTube playlists into Spotify playlists. This tool fetches tracks from a YouTube playlist, searches for the corresponding tracks on Spotify, and creates a Spotify playlist with the matched tracks.

## Features

- Fetches tracks from a public YouTube playlist.
- Cleans and matches track titles and artist names to improve Spotify search results.
- Creates a new playlist on Spotify and populates it with matched tracks.
- Handles mismatched or missing tracks gracefully with error logging.

---

## Prerequisites

1. **Spotify Developer Account**
   - Register an app in the [Spotify Developer Dashboard](https://developer.spotify.com/dashboard/).
   - Obtain the `Client ID`, `Client Secret`, and set a redirect URI for authentication.

2. **Google Cloud Project**
   - Enable the **YouTube Data API v3** in the [Google Cloud Console](https://console.cloud.google.com/).
   - Obtain the API Key.

3. **Go Environment**
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
