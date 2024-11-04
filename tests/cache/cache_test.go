package cache_test

import (
	"YouTubeParser/internal/cache"
	"os"
	"testing"
)

func TestNewCacheService(t *testing.T) {
	dbPath := "test_cache.db"
	defer os.Remove(dbPath)

	cache := cache.NewCacheService(dbPath)
	if cache == nil {
		t.Fatal("CacheService initialization failed")
	}
}

func TestCacheService_SaveAndGetThumbnail(t *testing.T) {
	dbPath := "test_cache.db"
	defer os.Remove(dbPath)

	cache := cache.NewCacheService(dbPath)
	videoID := "test_video"
	thumbnailURL := "http://example.com/thumbnail.jpg"

	err := cache.SaveThumbnail(videoID, thumbnailURL)
	if err != nil {
		t.Fatalf("Failed to save thumbnail: %v", err)
	}

	url, err := cache.GetThumbnail(videoID)
	if err != nil {
		t.Fatalf("Failed to get thumbnail: %v", err)
	}
	if url != thumbnailURL {
		t.Errorf("Expected %v, got %v", thumbnailURL, url)
	}
}
