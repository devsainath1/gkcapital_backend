package services

import (
	"gk-capital-backend/dto"
	"gk-capital-backend/models"
	"gk-capital-backend/repository"
)

type HomepageService struct {
	repo *repository.HomepageRepository
}

func NewHomepageService(repo *repository.HomepageRepository) *HomepageService {
	return &HomepageService{repo: repo}
}

func (s *HomepageService) GetAllActive() ([]models.HomepageSection, error) {
	return s.repo.FindAllActive()
}

func (s *HomepageService) GetAllAdmin(page, pageSize int) ([]models.HomepageSection, int64, error) {
	return s.repo.FindAllAdmin(page, pageSize)
}

func (s *HomepageService) GetByID(id uint) (*models.HomepageSection, error) {
	return s.repo.FindByID(id)
}

func (s *HomepageService) Create(req dto.CreateHomepageSectionRequest) (*models.HomepageSection, error) {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	section := &models.HomepageSection{
		SectionKey:  req.SectionKey,
		Title:       req.Title,
		Subtitle:    req.Subtitle,
		Description: req.Description,
		Image:       req.Image,
		Content:     req.Content,
		IsActive:    isActive,
		SortOrder:   req.SortOrder,
	}

	err := s.repo.Create(section)
	if err != nil {
		return nil, err
	}

	return section, nil
}

func (s *HomepageService) Update(id uint, req dto.UpdateHomepageSectionRequest) (*models.HomepageSection, error) {
	section, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.SectionKey != "" {
		section.SectionKey = req.SectionKey
	}
	if req.Title != "" {
		section.Title = req.Title
	}
	if req.Subtitle != "" {
		section.Subtitle = req.Subtitle
	}
	if req.Description != "" {
		section.Description = req.Description
	}
	if req.Image != "" {
		section.Image = req.Image
	}
	if req.Content != nil {
		section.Content = req.Content
	}
	if req.IsActive != nil {
		section.IsActive = *req.IsActive
	}
	if req.SortOrder != 0 {
		section.SortOrder = req.SortOrder
	}

	err = s.repo.Update(section)
	if err != nil {
		return nil, err
	}

	return section, nil
}

func (s *HomepageService) Delete(id uint) error {
	return s.repo.Delete(id)
}
