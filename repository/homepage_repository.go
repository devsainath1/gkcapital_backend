package repository

import (
	"gk-capital-backend/models"

	"gorm.io/gorm"
)

type HomepageRepository struct {
	db *gorm.DB
}

func NewHomepageRepository(db *gorm.DB) *HomepageRepository {
	return &HomepageRepository{db: db}
}

func (r *HomepageRepository) FindAllActive() ([]models.HomepageSection, error) {
	var sections []models.HomepageSection
	err := r.db.Where("is_active = ?", true).Order("sort_order ASC").Find(&sections).Error
	return sections, err
}

func (r *HomepageRepository) FindAllAdmin(page, pageSize int) ([]models.HomepageSection, int64, error) {
	var sections []models.HomepageSection
	var total int64

	r.db.Model(&models.HomepageSection{}).Count(&total)

	err := r.db.Order("sort_order ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&sections).Error

	return sections, total, err
}

func (r *HomepageRepository) FindByID(id uint) (*models.HomepageSection, error) {
	var section models.HomepageSection
	err := r.db.First(&section, id).Error
	if err != nil {
		return nil, err
	}
	return &section, nil
}

func (r *HomepageRepository) FindByKey(key string) (*models.HomepageSection, error) {
	var section models.HomepageSection
	err := r.db.Where("section_key = ?", key).First(&section).Error
	if err != nil {
		return nil, err
	}
	return &section, nil
}

func (r *HomepageRepository) Create(section *models.HomepageSection) error {
	return r.db.Create(section).Error
}

func (r *HomepageRepository) Update(section *models.HomepageSection) error {
	return r.db.Save(section).Error
}

func (r *HomepageRepository) Delete(id uint) error {
	return r.db.Delete(&models.HomepageSection{}, id).Error
}
