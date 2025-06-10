package handler

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uvish/url-shortener/internal/config"
	"github.com/uvish/url-shortener/internal/service"
)

type shortenRequest struct {
	URL string `json:"url" binding:"required"`
}

type urlResponse struct {
	ID          uint   `json:"id"`
	OriginalURL string `json:"original_url"`
	ShortCode   string `json:"short_code"`
	CreatedAt   string `json:"created_at"`
	ClickCount  int    `json:"click_count"`
}

func ShortenURL(c *gin.Context) {
	var req shortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if _, err := url.ParseRequestURI(req.URL); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL format"})
		return
	}

	shortCode, err := service.Shorten(req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not shorten URL"})
		return
	}

	shortURL := config.Cfg.Domain + "/" + shortCode
	c.JSON(http.StatusOK, gin.H{"short_url": shortURL})
}

func Redirect(c *gin.Context) {
	shortCode := c.Param("shortCode")

	url, err := service.GetOriginalURL(shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	service.IncrementClickCount(shortCode)
	c.Redirect(http.StatusFound, url.OriginalURL)
}

func GetAllURLs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	urls, total, err := service.GetAllURLs(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve URLs"})
		return
	}

	response := make([]urlResponse, len(urls))
	for i, u := range urls {
		response[i] = urlResponse{
			ID:          u.ID,
			OriginalURL: u.OriginalURL,
			ShortCode:   u.ShortCode,
			CreatedAt:   u.CreatedAt.Format(time.RFC3339),
			ClickCount:  u.ClickCount,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"urls":        response,
		"page":        page,
		"page_size":   pageSize,
		"total":       total,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

func DeleteURL(c *gin.Context) {
	shortCode := c.Param("shortCode")
	if err := service.DeleteURL(shortCode); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "URL deleted successfully"})
}
