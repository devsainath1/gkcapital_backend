package services

import (
	"gk-capital-backend/dto"
	"gk-capital-backend/models"
	"gk-capital-backend/repository"
	"strings"
)

type ServiceService struct {
	repo *repository.ServiceRepository
}

func NewServiceService(repo *repository.ServiceRepository) *ServiceService {
	return &ServiceService{repo: repo}
}

func (s *ServiceService) GetAll() ([]models.Service, error) {
	return s.repo.FindAll()
}

func (s *ServiceService) GetAllAdmin(page, pageSize int) ([]models.Service, int64, error) {
	return s.repo.FindAllAdmin(page, pageSize)
}

func (s *ServiceService) GetByID(id uint) (*models.Service, error) {
	return s.repo.FindByID(id)
}

func (s *ServiceService) Create(req dto.CreateServiceRequest) (*models.Service, error) {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	slug := req.Slug
	if slug == "" {
		slug = generateSlug(req.Title)
	}

	service := &models.Service{
		Title:       req.Title,
		Slug:        slug,
		Description: req.Description,
		Image:       req.Image,
		Icon:        req.Icon,
		IsActive:    isActive,
		SortOrder:   req.SortOrder,
	}

	err := s.repo.Create(service)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (s *ServiceService) Update(id uint, req dto.UpdateServiceRequest) (*models.Service, error) {
	service, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Title != "" {
		service.Title = req.Title
	}
	if req.Slug != "" {
		service.Slug = req.Slug
	}
	if req.Description != "" {
		service.Description = req.Description
	}
	if req.Image != "" {
		service.Image = req.Image
	}
	if req.Icon != "" {
		service.Icon = req.Icon
	}
	if req.IsActive != nil {
		service.IsActive = *req.IsActive
	}
	if req.SortOrder != 0 {
		service.SortOrder = req.SortOrder
	}

	err = s.repo.Update(service)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (s *ServiceService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func generateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, slug)
	return slug
}
