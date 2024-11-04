package cache

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

// CacheService - сервис кеша с SQLite
type CacheService struct {
	db *sql.DB
}

// NewCacheService - инициализация подключения к базе данных
func NewCacheService(dbPath string) *CacheService {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}

	// Создание таблицы, если она не существует
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS thumbnails (
        video_id TEXT PRIMARY KEY,
        thumbnail_url TEXT,
        created_at TIMESTAMP
    )`)
	if err != nil {
		log.Fatal("Ошибка при создании таблицы:", err)
	}

	return &CacheService{db: db}
}

// SaveThumbnail - сохраняет превью в базе данных
func (c *CacheService) SaveThumbnail(videoID, thumbnailURL string) error {
	_, err := c.db.Exec(`INSERT OR REPLACE INTO thumbnails (video_id, thumbnail_url, created_at)
                         VALUES (?, ?, ?)`, videoID, thumbnailURL, time.Now())
	return err
}

// GetThumbnail - возвращает URL превью из кеша, если он существует
func (c *CacheService) GetThumbnail(videoID string) (string, error) {
	var url string
	err := c.db.QueryRow(`SELECT thumbnail_url FROM thumbnails WHERE video_id = ?`, videoID).Scan(&url)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return url, err
}
