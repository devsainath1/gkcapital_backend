package services

import (
	"gk-capital-backend/dto"
	"gk-capital-backend/models"
	"gk-capital-backend/repository"
)

type TestimonialService struct {
	repo *repository.TestimonialRepository
}

func NewTestimonialService(repo *repository.TestimonialRepository) *TestimonialService {
	return &TestimonialService{repo: repo}
}

func (s *TestimonialService) GetAllActive() ([]models.Testimonial, error) {
	return s.repo.FindAllActive()
}

func (s *TestimonialService) GetAllAdmin(page, pageSize int) ([]models.Testimonial, int64, error) {
	return s.repo.FindAllAdmin(page, pageSize)
}

func (s *TestimonialService) GetByID(id uint) (*models.Testimonial, error) {
	return s.repo.FindByID(id)
}

func (s *TestimonialService) Create(req dto.CreateTestimonialRequest) (*models.Testimonial, error) {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	rating := req.Rating
	if rating == 0 {
		rating = 5
	}

	testimonial := &models.Testimonial{
		Name:          req.Name,
		LegacyName:    req.Name,
		Designation:   req.Designation,
		Company:       req.Company,
		Content:       req.Content,
		LegacyContent: req.Content,
		Rating:        rating,
		Image:         req.Image,
		IsActive:      isActive,
		SortOrder:     req.SortOrder,
	}

	err := s.repo.Create(testimonial)
	if err != nil {
		return nil, err
	}

	return testimonial, nil
}

func (s *TestimonialService) Update(id uint, req dto.UpdateTestimonialRequest) (*models.Testimonial, error) {
	testimonial, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		testimonial.Name = req.Name
		testimonial.LegacyName = req.Name
	}
	if req.Designation != "" {
		testimonial.Designation = req.Designation
	}
	if req.Company != "" {
		testimonial.Company = req.Company
	}
	if req.Content != "" {
		testimonial.Content = req.Content
		testimonial.LegacyContent = req.Content
	}
	if req.Rating != 0 {
		testimonial.Rating = req.Rating
	}
	if req.Image != "" {
		testimonial.Image = req.Image
	}
	if req.IsActive != nil {
		testimonial.IsActive = *req.IsActive
	}
	if req.SortOrder != 0 {
		testimonial.SortOrder = req.SortOrder
	}

	err = s.repo.Update(testimonial)
	if err != nil {
		return nil, err
	}

	return testimonial, nil
}

func (s *TestimonialService) Delete(id uint) error {
	return s.repo.Delete(id)
}
