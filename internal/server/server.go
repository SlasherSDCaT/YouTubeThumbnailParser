package server

import (
	"YouTubeParser/internal/cache"
	"YouTubeParser/internal/proto"
	"YouTubeParser/internal/youtube"
	"context"
	"log"
)

// Server - реализация сервиса ThumbnailService
type Server struct {
	proto.UnimplementedThumbnailServiceServer
	cache  *cache.CacheService
	apiKey string
}

// NewServer - создание нового gRPC сервера
func NewServer(apiKey string, cache *cache.CacheService) *Server {
	return &Server{apiKey: apiKey, cache: cache}
}

// GetThumbnail - метод для получения превью
func (s *Server) GetThumbnail(ctx context.Context, req *proto.ThumbnailRequest) (*proto.ThumbnailResponse, error) {
	log.Printf("Получен запрос на превью для видео URL: %s", req.VideoUrl)

	videoID, err := youtube.GetYouTubeVideoID(req.VideoUrl)
	if err != nil {
		log.Printf("Ошибка извлечения ID видео: %v", err)
		return nil, err
	}

	// Проверка кеша
	thumbnailURL, err := s.cache.GetThumbnail(videoID)
	if err != nil {
		log.Printf("Ошибка доступа к кешу: %v", err)
		return nil, err
	}
	if thumbnailURL == "" {
		thumbnailURL, err = youtube.GetThumbnailURL(s.apiKey, videoID)
		if err != nil {
			log.Printf("Ошибка получения превью с YouTube API: %v", err)
			return nil, err
		}
		// Сохранение в кеш
		if err := s.cache.SaveThumbnail(videoID, thumbnailURL); err != nil {
			log.Printf("Ошибка сохранения превью в кеш: %v", err)
		}
	}

	log.Printf("Отправка превью URL: %s для видео %s", thumbnailURL, req.VideoUrl)
	return &proto.ThumbnailResponse{ThumbnailUrl: thumbnailURL}, nil
}

// HealthCheck - проверка статуса сервера
func (s *Server) HealthCheck(ctx context.Context, req *proto.HealthRequest) (*proto.HealthResponse, error) {
	log.Println("Получен запрос HealthCheck")
	response := &proto.HealthResponse{Status: true}
	log.Println("Отправка ответа HealthCheck:", response.Status)
	return response, nil
}
