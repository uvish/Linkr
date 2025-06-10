package repository

import (
	"github.com/uvish/url-shortener/internal/db"
	"github.com/uvish/url-shortener/internal/model"
)

func CreateURL(url *model.URL) error {
	return db.DB.Create(url).Error
}

func GetURLByShortCode(code string) (*model.URL, error) {
	var url model.URL
	result := db.DB.Where("short_code = ?", code).First(&url)
	return &url, result.Error
}

func IncrementClickCount(url *model.URL) error {
	url.ClickCount++
	return db.DB.Save(url).Error
}

func GetAllURLs(page, pageSize int) ([]model.URL, int64, error) {
	var urls []model.URL
	var total int64

	offset := (page - 1) * pageSize
	result := db.DB.Model(&model.URL{}).Count(&total).Offset(offset).Limit(pageSize).Find(&urls)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return urls, total, nil
}

func DeleteURL(shortCode string) error {
	var url model.URL
	result := db.DB.Where("short_code = ?", shortCode).First(&url)
	if result.Error != nil {
		return result.Error
	}
	return db.DB.Delete(&url).Error
}
