package repository

import (
	"gk-capital-backend/models"

	"gorm.io/gorm"
)

type LoanInquiryRepository struct {
	db *gorm.DB
}

func NewLoanInquiryRepository(db *gorm.DB) *LoanInquiryRepository {
	return &LoanInquiryRepository{db: db}
}

func (r *LoanInquiryRepository) Create(inquiry *models.LoanInquiry) error {
	return r.db.Create(inquiry).Error
}

func (r *LoanInquiryRepository) FindAll(page, pageSize int, status, search string) ([]models.LoanInquiry, int64, error) {
	var inquiries []models.LoanInquiry
	var total int64

	query := r.db.Model(&models.LoanInquiry{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("full_name LIKE ? OR email LIKE ? OR phone LIKE ? OR city LIKE ?", searchTerm, searchTerm, searchTerm, searchTerm)
	}

	query.Count(&total)

	err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&inquiries).Error

	return inquiries, total, err
}

func (r *LoanInquiryRepository) FindByID(id uint) (*models.LoanInquiry, error) {
	var inquiry models.LoanInquiry
	err := r.db.First(&inquiry, id).Error
	if err != nil {
		return nil, err
	}
	return &inquiry, nil
}

func (r *LoanInquiryRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.LoanInquiry{}).Where("id = ?", id).Update("status", status).Error
}

func (r *LoanInquiryRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.LoanInquiry{}).Count(&count).Error
	return count, err
}
