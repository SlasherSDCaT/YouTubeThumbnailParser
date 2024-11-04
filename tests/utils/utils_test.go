package utils_test

import (
	"YouTubeParser/internal/utils"
	"os"
	"testing"
)

func TestReadURLsFromFile(t *testing.T) {
	filePath := "test_urls.txt"
	content := "https://youtube.com/watch?v=test_video1\nhttps://youtube.com/watch?v=test_video2\n"
	os.WriteFile(filePath, []byte(content), 0644)
	defer os.Remove(filePath)

	urls, err := utils.ReadURLsFromFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read URLs: %v", err)
	}
	expected := []string{"https://youtube.com/watch?v=test_video1", "https://youtube.com/watch?v=test_video2"}
	if len(urls) != len(expected) {
		t.Errorf("Expected %d URLs, got %d", len(expected), len(urls))
	}
	for i, url := range urls {
		if url != expected[i] {
			t.Errorf("Expected %v, got %v", expected[i], url)
		}
	}
}
