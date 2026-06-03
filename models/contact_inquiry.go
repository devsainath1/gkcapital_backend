package models

import (
	"time"
)

type ContactInquiry struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	Email     string    `gorm:"size:255;not null" json:"email"`
	Phone     string    `gorm:"size:20" json:"phone"`
	Message   string    `gorm:"type:text;not null" json:"message"`
	Status    string    `gorm:"size:50;default:new" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ContactInquiry) TableName() string {
	return "contact_inquiries"
}
