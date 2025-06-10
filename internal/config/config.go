package config

import (
	"encoding/json"
	"log"
	"os"
)

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
	SSLMode  string `json:"sslmode"`
}
type CacheConfig struct {
	Size int `json:"size"`
}

type Config struct {
	Domain   string         `json:"domain"`
	Database DatabaseConfig `json:"database"`
	Cache    CacheConfig    `json:"cache"`
}

var Cfg Config

func LoadConfig(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	err = json.Unmarshal(data, &Cfg)
	if err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	log.Println("Config loaded successfully")
}
