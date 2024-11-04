package youtube

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// getYouTubeVideoID - извлечение ID видео из URL
func GetYouTubeVideoID(url string) (string, error) {
	if strings.Contains(url, "v=") {
		parts := strings.Split(url, "v=")
		if len(parts) > 1 {
			return parts[1], nil
		}
	}
	return "", errors.New("неверный URL видео")
}

// GetThumbnailURL - получает URL превью видео по его идентификатору
func GetThumbnailURL(apiKey, videoID string) (string, error) {
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?id=%s&key=%s&part=snippet", videoID, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("не удалось получить данные с YouTube API")
	}

	var result struct {
		Items []struct {
			Snippet struct {
				Thumbnails struct {
					Medium struct {
						URL string `json:"url"`
					} `json:"medium"`
				} `json:"thumbnails"`
			} `json:"snippet"`
		} `json:"items"`
	}

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	if len(result.Items) > 0 {
		return result.Items[0].Snippet.Thumbnails.Medium.URL, nil
	}
	return "", errors.New("превью не найдено")
}
