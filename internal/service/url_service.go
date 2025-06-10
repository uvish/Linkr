package service

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/uvish/url-shortener/internal/cache"
	"github.com/uvish/url-shortener/internal/model"
	"github.com/uvish/url-shortener/internal/repository"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Shorten(originalURL string) (string, error) {
	const maxAttempts = 5
	for attempt := 0; attempt < maxAttempts; attempt++ {
		shortCode := generateShortCode(6)
		url := &model.URL{
			OriginalURL: originalURL,
			ShortCode:   shortCode,
			CreatedAt:   time.Now(),
		}

		err := repository.CreateURL(url)
		if err == nil {
			cache.AddURL(url)
			return shortCode, nil
		}
		// If it's a duplicate key error, try again
		if err.Error() != "pq: duplicate key value violates unique constraint" {
			return "", err
		}
	}
	return "", fmt.Errorf("failed to generate unique short code after %d attempts", maxAttempts)
}

func GetOriginalURL(shortCode string) (*model.URL, error) {
	// Check cache first
	if url, ok := cache.GetURL(shortCode); ok {
		log.Printf("Cache hit for shortCode: %s", shortCode)
		return url, nil
	}
	// Fallback to database
	log.Printf("Cache miss for shortCode: %s", shortCode)
	url, err := repository.GetURLByShortCode(shortCode)
	if err != nil {
		return nil, err
	}
	// Add to cache
	cache.AddURL(url)
	return url, nil
}

func GetAllURLs(page, pageSize int) ([]model.URL, int64, error) {
	return repository.GetAllURLs(page, pageSize)
}

func DeleteURL(shortCode string) error {
	cache.URLCache.Remove(shortCode) // Remove from cache first
	return repository.DeleteURL(shortCode)
}

func generateShortCode(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func IncrementClickCount(shortCode string) {
	url, err := repository.GetURLByShortCode(shortCode)
	if err != nil {
		log.Printf("Error finding URL to increment click count: %v", err)
		return
	}
	cache.AddURL(url)
	repository.IncrementClickCount(url)
}
