package cache

import (
	lru "github.com/hashicorp/golang-lru"
	"github.com/uvish/url-shortener/internal/model"
)

var URLCache *lru.Cache

func InitCache(size int) error {
	cache, err := lru.New(size)
	if err != nil {
		return err
	}
	URLCache = cache
	return nil
}

func GetURL(shortCode string) (*model.URL, bool) {
	value, ok := URLCache.Get(shortCode)
	if !ok {
		return nil, false
	}
	url, ok := value.(model.URL)
	if !ok {
		return nil, false
	}
	return &url, true
}

func AddURL(url *model.URL) {
	URLCache.Add(url.ShortCode, *url)
}
