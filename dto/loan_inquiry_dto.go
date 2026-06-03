package dto

// LoanInquiryRequest represents a loan application form submission
type LoanInquiryRequest struct {
	FullName      string `json:"full_name" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	Phone         string `json:"phone" binding:"required"`
	MonthlyIncome string `json:"monthly_income" binding:"required"`
	City          string `json:"city" binding:"required"`
	Service       string `json:"service" binding:"required"`
}

// UpdateInquiryStatusRequest represents a status update for inquiries
type UpdateInquiryStatusRequest struct {
	Status string `json:"status" binding:"required"`
}
