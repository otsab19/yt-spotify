package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// OllamaService defines the interface for the Ollama API interaction.
type OllamaService interface {
	AiService
	IsOllamaAvailable() bool
}

// OllamaServiceImpl implements OllamaService.
type OllamaServiceImpl struct {
	apiURL string
	model  string
}

// OllamaRequest represents the request payload for Ollama.
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// OllamaResponse represents the response from Ollama.
type OllamaResponse struct {
	Response string `json:"response"`
}

// NewOllamaService initializes a new OllamaServiceImpl.
func NewOllamaService() OllamaService {
	//todo: check model
	return &OllamaServiceImpl{
		apiURL: "http://localhost:11434/api/generate",
		model:  "llama3.2",
	}
}

// IsOllamaAvailable checks if the Ollama API is running.
func (o *OllamaServiceImpl) IsOllamaAvailable() bool {
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://localhost:11434/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// ExtractSongArtist calls Ollama and extracts the song and artist from the response.
func (o *OllamaServiceImpl) ExtractSongArtist(videoTitle string) (string, string, error) {
	// Define the prompt
	prompt := fmt.Sprintf("Extract the song title and artist from this YouTube video title: '%s'. Return it in the format: Song: <song_name>, Artist: <artist_name>.", videoTitle)

	// Prepare JSON request body
	requestBody, _ := json.Marshal(OllamaRequest{
		Model:  o.model,
		Prompt: prompt,
	})

	// Send HTTP request to Ollama
	resp, err := http.Post(o.apiURL, "application/json", bytes.NewBuffer(requestBody))
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

// parseOllamaResponse extracts the song title and artist name from Ollama response.
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
