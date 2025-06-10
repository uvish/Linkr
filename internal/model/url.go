package model

import "time"

type URL struct {
	ID          uint   `gorm:"primaryKey"`
	OriginalURL string `gorm:"not null"`
	ShortCode   string `gorm:"unique;not null"`
	CreatedAt   time.Time
	ClickCount  int
}
