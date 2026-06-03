package services

import (
	"gk-capital-backend/dto"
	"gk-capital-backend/models"
	"gk-capital-backend/repository"
)

type SEOService struct {
	repo *repository.SEORepository
}

func NewSEOService(repo *repository.SEORepository) *SEOService {
	return &SEOService{repo: repo}
}

func (s *SEOService) GetAll() ([]models.SEOPage, error) {
	return s.repo.FindAll()
}

func (s *SEOService) GetByID(id uint) (*models.SEOPage, error) {
	return s.repo.FindByID(id)
}

func (s *SEOService) GetBySlug(slug string) (*models.SEOPage, error) {
	return s.repo.FindBySlug(slug)
}

func (s *SEOService) Create(req dto.CreateSEOPageRequest) (*models.SEOPage, error) {
	page := &models.SEOPage{
		PageSlug:        req.PageSlug,
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
		MetaKeywords:    req.MetaKeywords,
		OGTitle:         req.OGTitle,
		OGDescription:   req.OGDescription,
		OGImage:         req.OGImage,
		CanonicalURL:    req.CanonicalURL,
	}

	err := s.repo.Create(page)
	if err != nil {
		return nil, err
	}

	return page, nil
}

func (s *SEOService) Update(id uint, req dto.UpdateSEOPageRequest) (*models.SEOPage, error) {
	page, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.PageSlug != "" {
		page.PageSlug = req.PageSlug
	}
	if req.MetaTitle != "" {
		page.MetaTitle = req.MetaTitle
	}
	if req.MetaDescription != "" {
		page.MetaDescription = req.MetaDescription
	}
	if req.MetaKeywords != "" {
		page.MetaKeywords = req.MetaKeywords
	}
	if req.OGTitle != "" {
		page.OGTitle = req.OGTitle
	}
	if req.OGDescription != "" {
		page.OGDescription = req.OGDescription
	}
	if req.OGImage != "" {
		page.OGImage = req.OGImage
	}
	if req.CanonicalURL != "" {
		page.CanonicalURL = req.CanonicalURL
	}

	err = s.repo.Update(page)
	if err != nil {
		return nil, err
	}

	return page, nil
}

func (s *SEOService) Delete(id uint) error {
	return s.repo.Delete(id)
}
