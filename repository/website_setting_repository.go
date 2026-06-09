package repository

import (
	"gk-capital-backend/models"

	"gorm.io/gorm"
)

type WebsiteSettingRepository struct {
	db *gorm.DB
}

func NewWebsiteSettingRepository(db *gorm.DB) *WebsiteSettingRepository {
	return &WebsiteSettingRepository{db: db}
}

func (r *WebsiteSettingRepository) FindAll() ([]models.WebsiteSetting, error) {
	var settings []models.WebsiteSetting
	err := r.db.Order("website_settings.key ASC").Find(&settings).Error
	return settings, err
}

func (r *WebsiteSettingRepository) FindByKey(key string) (*models.WebsiteSetting, error) {
	var setting models.WebsiteSetting
	err := r.db.Where("website_settings.key = ?", key).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *WebsiteSettingRepository) Upsert(setting *models.WebsiteSetting) error {
	var existing models.WebsiteSetting
	err := r.db.Where("website_settings.key = ?", setting.Key).First(&existing).Error
	if err == nil {
		existing.Value = setting.Value
		return r.db.Save(&existing).Error
	}
	return r.db.Create(setting).Error
}
