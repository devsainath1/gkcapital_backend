package services

import (
	"gk-capital-backend/repository"
)

type DashboardService struct {
	serviceRepo     *repository.ServiceRepository
	testimonialRepo *repository.TestimonialRepository
	contactRepo     *repository.ContactRepository
	loanInquiryRepo *repository.LoanInquiryRepository
}

type DashboardStats struct {
	TotalServices         int64 `json:"total_services"`
	TotalTestimonials     int64 `json:"total_testimonials"`
	TotalContactInquiries int64 `json:"total_contact_inquiries"`
	TotalLoanInquiries    int64 `json:"total_loan_inquiries"`
}

func NewDashboardService(
	serviceRepo *repository.ServiceRepository,
	testimonialRepo *repository.TestimonialRepository,
	contactRepo *repository.ContactRepository,
	loanInquiryRepo *repository.LoanInquiryRepository,
) *DashboardService {
	return &DashboardService{
		serviceRepo:     serviceRepo,
		testimonialRepo: testimonialRepo,
		contactRepo:     contactRepo,
		loanInquiryRepo: loanInquiryRepo,
	}
}

func (s *DashboardService) GetStats() (*DashboardStats, error) {
	services, err := s.serviceRepo.Count()
	if err != nil {
		return nil, err
	}

	testimonials, err := s.testimonialRepo.Count()
	if err != nil {
		return nil, err
	}

	contacts, err := s.contactRepo.Count()
	if err != nil {
		return nil, err
	}

	loans, err := s.loanInquiryRepo.Count()
	if err != nil {
		return nil, err
	}

	return &DashboardStats{
		TotalServices:         services,
		TotalTestimonials:     testimonials,
		TotalContactInquiries: contacts,
		TotalLoanInquiries:    loans,
	}, nil
}
