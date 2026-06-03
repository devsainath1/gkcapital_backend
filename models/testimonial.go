package models

import (
	"time"
)

type Testimonial struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Designation string    `gorm:"size:255" json:"designation"`
	Company     string    `gorm:"size:255" json:"company"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	Rating      int       `gorm:"default:5" json:"rating"`
	Image       string    `gorm:"size:500" json:"image"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Testimonial) TableName() string {
	return "testimonials"
}
