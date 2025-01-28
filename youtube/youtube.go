package youtube

import (
	"context"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// NewService creates a new YouTube service.
func NewService(apiKey string) (*youtube.Service, error) {
	service, err := youtube.NewService(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return service, nil
}

// FetchPlaylistItems fetches items from a YouTube playlist.
func FetchPlaylistItems(service *youtube.Service, playlistID string) ([]*youtube.PlaylistItem, error) {
	var items []*youtube.PlaylistItem
	nextPageToken := ""

	for {
		call := service.PlaylistItems.List([]string{"snippet"}).PlaylistId(playlistID).MaxResults(50).PageToken(nextPageToken)
		response, err := call.Do()
		if err != nil {
			return nil, err
		}

		items = append(items, response.Items...)
		nextPageToken = response.NextPageToken

		if nextPageToken == "" {
			break
		}
	}

	return items, nil
}
