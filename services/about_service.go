package services

import (
	"gk-capital-backend/dto"
	"gk-capital-backend/models"
	"gk-capital-backend/repository"
)

type AboutService struct {
	repo *repository.AboutRepository
}

func NewAboutService(repo *repository.AboutRepository) *AboutService {
	return &AboutService{repo: repo}
}

func (s *AboutService) GetAllActive() ([]models.AboutSection, error) {
	return s.repo.FindAllActive()
}

func (s *AboutService) GetAllAdmin(page, pageSize int) ([]models.AboutSection, int64, error) {
	return s.repo.FindAllAdmin(page, pageSize)
}

func (s *AboutService) GetByID(id uint) (*models.AboutSection, error) {
	return s.repo.FindByID(id)
}

func (s *AboutService) Create(req dto.CreateAboutSectionRequest) (*models.AboutSection, error) {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	section := &models.AboutSection{
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

func (s *AboutService) Update(id uint, req dto.UpdateAboutSectionRequest) (*models.AboutSection, error) {
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

func (s *AboutService) Delete(id uint) error {
	return s.repo.Delete(id)
}
