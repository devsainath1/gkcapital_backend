package repository

import (
	"gk-capital-backend/models"

	"gorm.io/gorm"
)

type ServiceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}

func (r *ServiceRepository) FindAll() ([]models.Service, error) {
	var services []models.Service
	err := r.db.Where("is_active = ?", true).Order("sort_order ASC").Find(&services).Error
	return services, err
}

func (r *ServiceRepository) FindAllAdmin(page, pageSize int) ([]models.Service, int64, error) {
	var services []models.Service
	var total int64

	r.db.Model(&models.Service{}).Count(&total)

	err := r.db.Order("sort_order ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&services).Error

	return services, total, err
}

func (r *ServiceRepository) FindByID(id uint) (*models.Service, error) {
	var service models.Service
	err := r.db.First(&service, id).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *ServiceRepository) Create(service *models.Service) error {
	return r.db.Create(service).Error
}

func (r *ServiceRepository) Update(service *models.Service) error {
	return r.db.Save(service).Error
}

func (r *ServiceRepository) Delete(id uint) error {
	return r.db.Delete(&models.Service{}, id).Error
}

func (r *ServiceRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Service{}).Count(&count).Error
	return count, err
}
