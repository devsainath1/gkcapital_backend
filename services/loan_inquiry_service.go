package services

import (
	"gk-capital-backend/dto"
	"gk-capital-backend/models"
	"gk-capital-backend/repository"
)

type LoanInquiryService struct {
	repo *repository.LoanInquiryRepository
}

func NewLoanInquiryService(repo *repository.LoanInquiryRepository) *LoanInquiryService {
	return &LoanInquiryService{repo: repo}
}

func (s *LoanInquiryService) Submit(req dto.LoanInquiryRequest) (*models.LoanInquiry, error) {
	inquiry := &models.LoanInquiry{
		FullName:      req.FullName,
		Email:         req.Email,
		Phone:         req.Phone,
		MonthlyIncome: req.MonthlyIncome,
		City:          req.City,
		Service:       req.Service,
		Status:        "new",
	}

	err := s.repo.Create(inquiry)
	if err != nil {
		return nil, err
	}

	return inquiry, nil
}

func (s *LoanInquiryService) GetAll(page, pageSize int, status, search string) ([]models.LoanInquiry, int64, error) {
	return s.repo.FindAll(page, pageSize, status, search)
}

func (s *LoanInquiryService) GetByID(id uint) (*models.LoanInquiry, error) {
	return s.repo.FindByID(id)
}

func (s *LoanInquiryService) UpdateStatus(id uint, status string) error {
	validStatuses := map[string]bool{
		"new":       true,
		"contacted": true,
		"approved":  true,
		"rejected":  true,
	}

	if !validStatuses[status] {
		return nil
	}

	return s.repo.UpdateStatus(id, status)
}
