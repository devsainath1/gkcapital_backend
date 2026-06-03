package services

import (
	"gk-capital-backend/dto"
	"gk-capital-backend/models"
	"gk-capital-backend/repository"
)

type FAQService struct {
	repo *repository.FAQRepository
}

func NewFAQService(repo *repository.FAQRepository) *FAQService {
	return &FAQService{repo: repo}
}

func (s *FAQService) GetAllActive() ([]models.FAQ, error) {
	return s.repo.FindAllActive()
}

func (s *FAQService) GetAllAdmin(page, pageSize int) ([]models.FAQ, int64, error) {
	return s.repo.FindAllAdmin(page, pageSize)
}

func (s *FAQService) GetByID(id uint) (*models.FAQ, error) {
	return s.repo.FindByID(id)
}

func (s *FAQService) Create(req dto.CreateFAQRequest) (*models.FAQ, error) {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	faq := &models.FAQ{
		Question:  req.Question,
		Answer:    req.Answer,
		Category:  req.Category,
		IsActive:  isActive,
		SortOrder: req.SortOrder,
	}

	err := s.repo.Create(faq)
	if err != nil {
		return nil, err
	}

	return faq, nil
}

func (s *FAQService) Update(id uint, req dto.UpdateFAQRequest) (*models.FAQ, error) {
	faq, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Question != "" {
		faq.Question = req.Question
	}
	if req.Answer != "" {
		faq.Answer = req.Answer
	}
	if req.Category != "" {
		faq.Category = req.Category
	}
	if req.IsActive != nil {
		faq.IsActive = *req.IsActive
	}
	if req.SortOrder != 0 {
		faq.SortOrder = req.SortOrder
	}

	err = s.repo.Update(faq)
	if err != nil {
		return nil, err
	}

	return faq, nil
}

func (s *FAQService) Delete(id uint) error {
	return s.repo.Delete(id)
}
