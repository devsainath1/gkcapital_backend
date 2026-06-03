package repository

import (
	"gk-capital-backend/models"

	"gorm.io/gorm"
)

type TestimonialRepository struct {
	db *gorm.DB
}

func NewTestimonialRepository(db *gorm.DB) *TestimonialRepository {
	return &TestimonialRepository{db: db}
}

func (r *TestimonialRepository) FindAllActive() ([]models.Testimonial, error) {
	var testimonials []models.Testimonial
	err := r.db.Where("is_active = ?", true).Order("sort_order ASC").Find(&testimonials).Error
	return testimonials, err
}

func (r *TestimonialRepository) FindAllAdmin(page, pageSize int) ([]models.Testimonial, int64, error) {
	var testimonials []models.Testimonial
	var total int64

	r.db.Model(&models.Testimonial{}).Count(&total)

	err := r.db.Order("sort_order ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&testimonials).Error

	return testimonials, total, err
}

func (r *TestimonialRepository) FindByID(id uint) (*models.Testimonial, error) {
	var testimonial models.Testimonial
	err := r.db.First(&testimonial, id).Error
	if err != nil {
		return nil, err
	}
	return &testimonial, nil
}

func (r *TestimonialRepository) Create(testimonial *models.Testimonial) error {
	return r.db.Create(testimonial).Error
}

func (r *TestimonialRepository) Update(testimonial *models.Testimonial) error {
	return r.db.Save(testimonial).Error
}

func (r *TestimonialRepository) Delete(id uint) error {
	return r.db.Delete(&models.Testimonial{}, id).Error
}

func (r *TestimonialRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Testimonial{}).Count(&count).Error
	return count, err
}
