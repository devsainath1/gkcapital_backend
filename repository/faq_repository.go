package repository

import (
	"gk-capital-backend/models"

	"gorm.io/gorm"
)

type FAQRepository struct {
	db *gorm.DB
}

func NewFAQRepository(db *gorm.DB) *FAQRepository {
	return &FAQRepository{db: db}
}

func (r *FAQRepository) FindAllActive() ([]models.FAQ, error) {
	var faqs []models.FAQ
	err := r.db.Where("is_active = ?", true).Order("sort_order ASC").Find(&faqs).Error
	return faqs, err
}

func (r *FAQRepository) FindAllAdmin(page, pageSize int) ([]models.FAQ, int64, error) {
	var faqs []models.FAQ
	var total int64

	r.db.Model(&models.FAQ{}).Count(&total)

	err := r.db.Order("sort_order ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&faqs).Error

	return faqs, total, err
}

func (r *FAQRepository) FindByID(id uint) (*models.FAQ, error) {
	var faq models.FAQ
	err := r.db.First(&faq, id).Error
	if err != nil {
		return nil, err
	}
	return &faq, nil
}

func (r *FAQRepository) Create(faq *models.FAQ) error {
	return r.db.Create(faq).Error
}

func (r *FAQRepository) Update(faq *models.FAQ) error {
	return r.db.Save(faq).Error
}

func (r *FAQRepository) Delete(id uint) error {
	return r.db.Delete(&models.FAQ{}, id).Error
}
