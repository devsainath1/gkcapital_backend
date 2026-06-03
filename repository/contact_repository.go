package repository

import (
	"gk-capital-backend/models"

	"gorm.io/gorm"
)

type ContactRepository struct {
	db *gorm.DB
}

func NewContactRepository(db *gorm.DB) *ContactRepository {
	return &ContactRepository{db: db}
}

func (r *ContactRepository) Create(inquiry *models.ContactInquiry) error {
	return r.db.Create(inquiry).Error
}

func (r *ContactRepository) FindAll(page, pageSize int, status, search string) ([]models.ContactInquiry, int64, error) {
	var inquiries []models.ContactInquiry
	var total int64

	query := r.db.Model(&models.ContactInquiry{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("name LIKE ? OR email LIKE ? OR phone LIKE ?", searchTerm, searchTerm, searchTerm)
	}

	query.Count(&total)

	err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&inquiries).Error

	return inquiries, total, err
}

func (r *ContactRepository) FindByID(id uint) (*models.ContactInquiry, error) {
	var inquiry models.ContactInquiry
	err := r.db.First(&inquiry, id).Error
	if err != nil {
		return nil, err
	}
	return &inquiry, nil
}

func (r *ContactRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.ContactInquiry{}).Where("id = ?", id).Update("status", status).Error
}

func (r *ContactRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.ContactInquiry{}).Count(&count).Error
	return count, err
}
