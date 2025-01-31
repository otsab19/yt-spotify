package service

type AiService interface {
	ExtractSongArtist(videoTitle string) (string, string, error)
}
