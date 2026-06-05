package repository

import (
	"gk-capital-backend/models"

	"gorm.io/gorm"
)

type MediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository(db *gorm.DB) *MediaRepository {
	return &MediaRepository{db: db}
}

func (r *MediaRepository) Create(asset *models.MediaAsset) error {
	return r.db.Create(asset).Error
}

func (r *MediaRepository) GetByID(id uint) (*models.MediaAsset, error) {
	var asset models.MediaAsset
	err := r.db.First(&asset, id).Error
	return &asset, err
}

func (r *MediaRepository) GetByName(name string) (*models.MediaAsset, error) {
	var asset models.MediaAsset
	err := r.db.Where("name = ?", name).First(&asset).Error
	return &asset, err
}

func (r *MediaRepository) GetAll() ([]models.MediaAsset, error) {
	var assets []models.MediaAsset
	// Select without Data column to keep listing lightweight
	err := r.db.Select("id, name, mime_type, size, created_at, updated_at").
		Order("created_at DESC").
		Find(&assets).Error
	return assets, err
}

func (r *MediaRepository) Delete(id uint) error {
	return r.db.Delete(&models.MediaAsset{}, id).Error
}

func (r *MediaRepository) ExistsByName(name string) bool {
	var count int64
	r.db.Model(&models.MediaAsset{}).Where("name = ?", name).Count(&count)
	return count > 0
}
