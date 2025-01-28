package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/oauth2"
)

// Authenticate authenticates with Spotify and returns an HTTP client.
func Authenticate(clientID, clientSecret, redirectURI string) (*http.Client, error) {
	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURI,
		Scopes:       []string{"playlist-modify-public", "playlist-modify-private"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.spotify.com/authorize",
			TokenURL: "https://accounts.spotify.com/api/token",
		},
	}

	// Redirect user to Spotify authorization page
	authURL := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v\n", authURL)

	// After authorization, Spotify will redirect to the redirect URI with a code
	fmt.Print("Enter the code from the redirect URL: ")
	var code string
	fmt.Scanln(&code)

	// Exchange the code for a token
	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	// Create an HTTP client using the token
	client := conf.Client(context.Background(), token)

	return client, nil
}

// CreatePlaylist creates a new Spotify playlist and returns its ID.
func CreatePlaylist(client *http.Client, name string) (string, error) {
	userID, err := getSpotifyUserID(client)
	if err != nil {
		return "", err
	}

	reqBody := map[string]interface{}{
		"name":        name,
		"description": "Playlist imported from YouTube",
		"public":      false,
	}

	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists", userID), strings.NewReader(string(reqBodyJSON)))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result["id"].(string), nil
}

func getSpotifyUserID(client *http.Client) (string, error) {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result["id"].(string), nil
}

// cleanText removes unnecessary keywords from track or artist names.
func cleanText(input string) string {
	re := regexp.MustCompile(`(?i)(official video|remastered|4k|live|audio|topic|vevo|remix)`)
	cleaned := re.ReplaceAllString(input, "")
	return strings.TrimSpace(cleaned)
}

// SearchTrack searches for a track on Spotify and returns its ID.
func SearchTrack(client *http.Client, trackName, artistName string) (string, error) {
	// Clean track and artist names
	trackName = cleanText(trackName)
	artistName = cleanText(artistName)

	fmt.Printf("Searching for track: '%s' by artist: '%s'\n", trackName, artistName)

	// First, try an exact match with track and artist
	query := url.QueryEscape(fmt.Sprintf("track:%s artist:%s", trackName, artistName))
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track&limit=5", query), nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// Handle Spotify's response for the first query
	tracks := result["tracks"].(map[string]interface{})["items"].([]interface{})
	if len(tracks) > 0 {
		// Attempt to find the best match in the first result set
		for _, t := range tracks {
			track := t.(map[string]interface{})
			trackTitle := track["name"].(string)
			trackArtist := track["artists"].([]interface{})[0].(map[string]interface{})["name"].(string)

			// Fuzzy match: Check if Spotify track matches the cleaned inputs
			if strings.Contains(strings.ToLower(trackTitle), strings.ToLower(trackName)) &&
				strings.Contains(strings.ToLower(trackArtist), strings.ToLower(artistName)) {
				return track["id"].(string), nil
			}
		}
	}

	// If no match is found, try a broader search with just the track name
	fmt.Printf("Exact match failed for '%s' by '%s'. Trying broader search...\n", trackName, artistName)
	query = url.QueryEscape(fmt.Sprintf("track:%s", trackName))
	req, err = http.NewRequest("GET", fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track&limit=5", query), nil)
	if err != nil {
		return "", err
	}

	resp, err = client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// Handle Spotify's response for the broader query
	tracks = result["tracks"].(map[string]interface{})["items"].([]interface{})
	if len(tracks) > 0 {
		// Attempt to find the best match in the broader result set
		for _, t := range tracks {
			track := t.(map[string]interface{})
			trackTitle := track["name"].(string)

			// Loose match: Check if Spotify track matches the cleaned track name
			if strings.Contains(strings.ToLower(trackTitle), strings.ToLower(trackName)) {
				return track["id"].(string), nil
			}
		}
	}

	// If no matches were found at all
	return "", fmt.Errorf("no suitable tracks found for '%s' by '%s'", trackName, artistName)
}

// AddTrackToPlaylist adds a track to a Spotify playlist.
func AddTrackToPlaylist(client *http.Client, playlistID, trackID string) error {
	reqBody := map[string]interface{}{
		"uris": []string{fmt.Sprintf("spotify:track:%s", trackID)},
	}

	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playlistID), strings.NewReader(string(reqBodyJSON)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to add track to playlist: %s", resp.Status)
	}

	return nil
}
