package youtube_test

import (
	"YouTubeParser/internal/youtube"
	"testing"
)

func TestGetYouTubeVideoID(t *testing.T) {
	url := "https://youtube.com/watch?v=test_video"
	videoID, err := youtube.GetYouTubeVideoID(url)
	if err != nil {
		t.Fatalf("Failed to get video ID: %v", err)
	}
	expectedID := "test_video"
	if videoID != expectedID {
		t.Errorf("Expected %v, got %v", expectedID, videoID)
	}
}

func TestGetThumbnailURL_InvalidKey(t *testing.T) {
	videoID := "test_video"
	apiKey := "invalid_key"
	_, err := youtube.GetThumbnailURL(apiKey, videoID)
	if err == nil {
		t.Error("Expected error for invalid API key, got nil")
	}
}
