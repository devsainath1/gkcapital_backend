package repository

import (
	"gk-capital-backend/models"

	"gorm.io/gorm"
)

type SEORepository struct {
	db *gorm.DB
}

func NewSEORepository(db *gorm.DB) *SEORepository {
	return &SEORepository{db: db}
}

func (r *SEORepository) FindAll() ([]models.SEOPage, error) {
	var pages []models.SEOPage
	err := r.db.Order("page_slug ASC").Find(&pages).Error
	return pages, err
}

func (r *SEORepository) FindByID(id uint) (*models.SEOPage, error) {
	var page models.SEOPage
	err := r.db.First(&page, id).Error
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *SEORepository) FindBySlug(slug string) (*models.SEOPage, error) {
	var page models.SEOPage
	err := r.db.Where("page_slug = ?", slug).First(&page).Error
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *SEORepository) Create(page *models.SEOPage) error {
	return r.db.Create(page).Error
}

func (r *SEORepository) Update(page *models.SEOPage) error {
	return r.db.Save(page).Error
}

func (r *SEORepository) Delete(id uint) error {
	return r.db.Delete(&models.SEOPage{}, id).Error
}
