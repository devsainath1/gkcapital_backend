package repository

import (
	"gk-capital-backend/models"

	"gorm.io/gorm"
)

type AboutRepository struct {
	db *gorm.DB
}

func NewAboutRepository(db *gorm.DB) *AboutRepository {
	return &AboutRepository{db: db}
}

func (r *AboutRepository) FindAllActive() ([]models.AboutSection, error) {
	var sections []models.AboutSection
	err := r.db.Where("is_active = ?", true).Order("sort_order ASC").Find(&sections).Error
	return sections, err
}

func (r *AboutRepository) FindAllAdmin(page, pageSize int) ([]models.AboutSection, int64, error) {
	var sections []models.AboutSection
	var total int64

	r.db.Model(&models.AboutSection{}).Count(&total)

	err := r.db.Order("sort_order ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&sections).Error

	return sections, total, err
}

func (r *AboutRepository) FindByID(id uint) (*models.AboutSection, error) {
	var section models.AboutSection
	err := r.db.First(&section, id).Error
	if err != nil {
		return nil, err
	}
	return &section, nil
}

func (r *AboutRepository) FindByKey(key string) (*models.AboutSection, error) {
	var section models.AboutSection
	err := r.db.Where("section_key = ?", key).First(&section).Error
	if err != nil {
		return nil, err
	}
	return &section, nil
}

func (r *AboutRepository) Create(section *models.AboutSection) error {
	return r.db.Create(section).Error
}

func (r *AboutRepository) Update(section *models.AboutSection) error {
	return r.db.Save(section).Error
}

func (r *AboutRepository) Delete(id uint) error {
	return r.db.Delete(&models.AboutSection{}, id).Error
}
