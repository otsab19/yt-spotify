package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

// Ollama API URL
const ollamaAPI = "http://localhost:11434/api/tags"

// Check if Ollama is running
func IsOllamaAvailable() bool {
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(ollamaAPI)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// ExtractSongArtist calls Ollama and extracts the song and artist from the response
func ExtractSongArtist(videoTitle string) (string, string, error) {
	// Define the prompt
	prompt := fmt.Sprintf("Extract the song title and artist from this YouTube video title: '%s'. Return it in the format: Song: <song_name>, Artist: <artist_name>.", videoTitle)

	// Prepare JSON request body
	requestBody, _ := json.Marshal(OllamaRequest{
		Model:  "llama3.2",
		Prompt: prompt,
	})

	// Send HTTP request to Ollama
	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error sending request to Ollama:", err)
		return "", "", err
	}
	defer resp.Body.Close()

	// Accumulate the full response
	var fullResponse strings.Builder
	decoder := json.NewDecoder(resp.Body)

	for decoder.More() {
		var chunk OllamaResponse
		err := decoder.Decode(&chunk)
		if err != nil {
			fmt.Println("Error decoding Ollama JSON chunk:", err)
			return "", "", err
		}

		// Append chunk response to the full response
		fullResponse.WriteString(chunk.Response)
	}

	// Convert response to string
	responseText := strings.TrimSpace(fullResponse.String())

	fmt.Println("🟢 Full Ollama Response:", responseText)

	// Extract song and artist from response
	song, artist := parseOllamaResponse(responseText)

	fmt.Println("🟢 Extracted Song:", song)
	fmt.Println("🟢 Extracted Artist:", artist)

	return song, artist, nil
}

// parseOllamaResponse extracts the song title and artist name from Ollama response
func parseOllamaResponse(response string) (string, string) {
	// Trim spaces and normalize response
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

// Example usage
func Run(videoTitle string) {
	song, artist, err := ExtractSongArtist(videoTitle)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("xtracted Song: %s, Artist: %s\n", song, artist)
	}
}
