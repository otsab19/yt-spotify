package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"
)

func TestExtractSongArtistWithOllama(t *testing.T) {
	// Define test cases with real YouTube titles
	testCases := []struct {
		videoTitle     string
		expectedFormat string // Expected format: "Song: <song_name>, Artist: <artist_name>"
	}{
		{"The Weeknd - Blinding Lights (Official Music Video)", "Song: Blinding Lights, Artist: The Weeknd"},
		{"Eminem | Lose Yourself (Lyrics)", "Song: Lose Yourself, Artist: Eminem"},
		{"Taylor Swift - Love Story (Official Video)", "Song: Love Story, Artist: Taylor Swift"},
	}

	// Real Ollama API URL
	apiURL := "http://localhost:11434/api/generate"

	for _, tc := range testCases {
		// Call Ollama API
		song, artist, err := ExtractSongArtistFromOllama(tc.videoTitle, apiURL)
		if err != nil {
			t.Errorf("Unexpected error for '%s': %v", tc.videoTitle, err)
			continue
		}

		// Validate response format
		responseFormat := fmt.Sprintf("Song: %s, Artist: %s", song, artist)
		if !validateResponseFormat(responseFormat) {
			t.Errorf("For title '%s', expected format '%s', but got '%s'",
				tc.videoTitle, tc.expectedFormat, responseFormat)
		} else {
			fmt.Printf("Test passed for '%s': %s\n", tc.videoTitle, responseFormat)
		}
	}
}

// Function to validate the response format
func validateResponseFormat(response string) bool {
	re := regexp.MustCompile(`(?i)Song:\s*(.*?),\s*Artist:\s*(.*)`)
	return re.MatchString(response)
}

// Calls Ollama to extract song and artist
func ExtractSongArtistFromOllama(videoTitle, apiURL string) (string, string, error) {
	// Define the prompt
	prompt := fmt.Sprintf("Extract the song title and artist from this YouTube video title: '%s'. Return it in the format: Song: <song_name>, Artist: <artist_name>.", videoTitle)

	// Prepare JSON request body
	requestBody, _ := json.Marshal(map[string]string{
		"service": "llama3.2", // Ensure this matches the running service
		"prompt":  prompt,
	})

	// Send HTTP request to Ollama
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error sending request to Ollama:", err)
		return "", "", err
	}
	defer resp.Body.Close()

	// Accumulate the full response
	var fullResponse strings.Builder
	decoder := json.NewDecoder(resp.Body)

	for decoder.More() {
		var chunk struct {
			Response string `json:"response"`
			Done     bool   `json:"done"`
		}

		// Decode each JSON chunk
		err := decoder.Decode(&chunk)
		if err != nil {
			fmt.Println("Error decoding Ollama JSON chunk:", err)
			return "", "", err
		}

		// Append chunk response to the full response
		fullResponse.WriteString(chunk.Response)

		// Stop reading if Ollama is done
		if chunk.Done {
			break
		}
	}

	// Convert response to string
	responseText := strings.TrimSpace(fullResponse.String())

	fmt.Println("Full Ollama Response:", responseText)

	// Extract song and artist from response
	song, artist := parseOllamaResponse(responseText)

	fmt.Println("Extracted Song:", song)
	fmt.Println("Extracted Artist:", artist)

	return song, artist, nil
}

// Parses the Ollama response to extract song and artist
func parseOllamaResponse(response string) (string, string) {
	// Trim whitespace
	response = strings.TrimSpace(response)

	// Define regex pattern: Song: <song_name>, Artist: <artist_name>
	re := regexp.MustCompile(`(?i)Song:\s*(.*?),\s*Artist:\s*(.*)`)

	// Extract matches
	matches := re.FindStringSubmatch(response)
	if len(matches) == 3 {
		song := strings.TrimSpace(matches[1])
		artist := strings.TrimSpace(matches[2])
		return song, artist
	}

	// Return empty values if parsing fails
	return "", ""
}
