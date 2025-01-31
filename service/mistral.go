package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"yt-spotify/config"
)

// MistralService defines the interface for extracting song and artist using Mistral AI.
type MistralService interface {
	AiService
	parseMistralResponse(response string) (string, string)
}

// MistralServiceImpl implements the MistralService interface.
type MistralServiceImpl struct {
	apiURL string
	model  string
	apiKey string
}

// NewMistralService initializes a new MistralServiceImpl with API Key from environment.
func NewMistralService(config *config.AppContext) (MistralService, error) {
	apiKey := config.MistralApiKey
	if apiKey == "" {
		return nil, fmt.Errorf("MISTRAL_API_KEY is not set")
	}

	return &MistralServiceImpl{
		apiURL: "https://api.mistral.ai/v1/chat/completions",
		model:  "mistral-large-latest",
		apiKey: apiKey,
	}, nil
}

// ExtractSongArtist calls Mistral AI API to get the song and artist.
func (m *MistralServiceImpl) ExtractSongArtist(videoTitle string) (string, string, error) {
	// Define the prompt
	prompt := fmt.Sprintf("Extract the song title and artist from this YouTube video title: '%s'. Return it in the format: Song: <song_name>, Artist: <artist_name>.", videoTitle)

	// Prepare JSON request body
	requestBody, _ := json.Marshal(map[string]interface{}{
		"model": m.model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"temperature": 0.7,
	})

	// Create HTTP request
	req, err := http.NewRequest("POST", m.apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", "", err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.apiKey)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	// Parse response
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}

	// Check response
	if len(result.Choices) == 0 {
		return "", "", fmt.Errorf("no response from Mistral AI")
	}

	responseText := strings.TrimSpace(result.Choices[0].Message.Content)

	// Debugging: Print full response
	fmt.Println("🟢 Mistral AI Response:", responseText)

	// Extract song and artist from response
	song, artist := m.parseMistralResponse(responseText)

	// Ensure we return all three values
	if song == "" || artist == "" {
		return "", "", fmt.Errorf("failed to extract song and artist from response")
	}

	return song, artist, nil
}

// parseMistralResponse extracts the song title and artist from the response
func (m *MistralServiceImpl) parseMistralResponse(response string) (string, string) {
	re := regexp.MustCompile(`(?i)Song:\s*(.*?),\s*Artist:\s*(.*)`)
	matches := re.FindStringSubmatch(response)
	if len(matches) == 3 {
		return strings.TrimSpace(matches[1]), strings.TrimSpace(matches[2])
	}
	return "", ""
}
