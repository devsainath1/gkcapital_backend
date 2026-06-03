package models

import (
	"time"
)

type LoanInquiry struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	FullName      string    `gorm:"size:255;not null" json:"full_name"`
	Email         string    `gorm:"size:255;not null" json:"email"`
	Phone         string    `gorm:"size:20;not null" json:"phone"`
	MonthlyIncome string    `gorm:"size:100" json:"monthly_income"`
	City          string    `gorm:"size:255" json:"city"`
	Service       string    `gorm:"size:255" json:"service"`
	Status        string    `gorm:"size:50;default:new" json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (LoanInquiry) TableName() string {
	return "loan_inquiries"
}
