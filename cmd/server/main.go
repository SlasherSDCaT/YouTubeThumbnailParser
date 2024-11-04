package main

import (
	"YouTubeParser/internal/cache"
	"YouTubeParser/internal/proto"
	"YouTubeParser/internal/server"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net"
)

// Config - структура для конфигурации
type Config struct {
	YouTubeAPIKey string `yaml:"youtube_api_key"`
	CacheDBPath   string `yaml:"cache_db_path"`
}

// LoadConfig - загрузка конфигурации
func LoadConfig() (*Config, error) {
	var config Config
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &config)
	return &config, err
}

func main() {
	config, err := LoadConfig()
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации:", err)
	}

	// Инициализация кеша
	cacheService := cache.NewCacheService(config.CacheDBPath)

	// Создание gRPC сервера
	grpcServer := grpc.NewServer()
	server := server.NewServer(config.YouTubeAPIKey, cacheService)
	proto.RegisterThumbnailServiceServer(grpcServer, server)

	// Запуск сервера
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
	log.Println("gRPC сервер запущен на :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Ошибка работы gRPC сервера:", err)
	}
}
