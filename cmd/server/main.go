package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/uvish/url-shortener/internal/cache"
	"github.com/uvish/url-shortener/internal/config"
	"github.com/uvish/url-shortener/internal/db"
	"github.com/uvish/url-shortener/internal/handler"
)

func main() {
	config.LoadConfig("config.json")
	if err := db.Connect(config.Cfg.Database); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	// Cache size of 1000 URLs
	if err := cache.InitCache(config.Cfg.Cache.Size); err != nil {
		log.Fatalf("Failed to initialize cache: %v", err)
	}

	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:1234"},
		AllowMethods:     []string{"GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}))

	r.POST("/shorten", handler.ShortenURL)
	r.GET("/:shortCode", handler.Redirect)
	r.GET("/urls", handler.GetAllURLs)
	r.DELETE("/urls/:shortCode", handler.DeleteURL)
	r.Run(":8080")
}
