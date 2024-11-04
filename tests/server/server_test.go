package server_test

import (
	"YouTubeParser/internal/cache"
	"YouTubeParser/internal/proto"
	"YouTubeParser/internal/server"
	"context"
	"os"
	"testing"
)

func TestServer_GetThumbnail(t *testing.T) {
	dbPath := "test_server_cache.db"
	defer os.Remove(dbPath)

	cache := cache.NewCacheService(dbPath)
	server := server.NewServer("test_api_key", cache)

	req := &proto.ThumbnailRequest{VideoUrl: "https://youtube.com/watch?v=test_video"}
	_, err := server.GetThumbnail(context.Background(), req)
	if err == nil {
		t.Error("Expected error for invalid YouTube API key, got nil")
	}
}

func TestServer_HealthCheck(t *testing.T) {
	cache := cache.NewCacheService("test_server_cache.db")
	server := server.NewServer("test_api_key", cache)

	resp, err := server.HealthCheck(context.Background(), &proto.HealthRequest{})
	if err != nil {
		t.Fatalf("HealthCheck failed: %v", err)
	}
	if !resp.Status {
		t.Error("Expected HealthCheck status to be true")
	}
}
