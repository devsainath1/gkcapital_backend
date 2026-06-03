package services

import (
	"gk-capital-backend/dto"
	"gk-capital-backend/models"
	"gk-capital-backend/repository"
)

type ContactService struct {
	repo *repository.ContactRepository
}

func NewContactService(repo *repository.ContactRepository) *ContactService {
	return &ContactService{repo: repo}
}

func (s *ContactService) Submit(req dto.ContactRequest) (*models.ContactInquiry, error) {
	inquiry := &models.ContactInquiry{
		Name:    req.Name,
		Email:   req.Email,
		Phone:   req.Phone,
		Message: req.Message,
		Status:  "new",
	}

	err := s.repo.Create(inquiry)
	if err != nil {
		return nil, err
	}

	return inquiry, nil
}

func (s *ContactService) GetAll(page, pageSize int, status, search string) ([]models.ContactInquiry, int64, error) {
	return s.repo.FindAll(page, pageSize, status, search)
}

func (s *ContactService) GetByID(id uint) (*models.ContactInquiry, error) {
	return s.repo.FindByID(id)
}

func (s *ContactService) UpdateStatus(id uint, status string) error {
	validStatuses := map[string]bool{
		"new":       true,
		"read":      true,
		"responded": true,
	}

	if !validStatuses[status] {
		return nil
	}

	return s.repo.UpdateStatus(id, status)
}
