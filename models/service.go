package models

import (
	"time"
)

type Service struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Slug        string    `gorm:"size:255;uniqueIndex" json:"slug"`
	Description string    `gorm:"type:text" json:"description"`
	Image       string    `gorm:"size:500" json:"image"`
	Icon        string    `gorm:"size:255" json:"icon"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Service) TableName() string {
	return "services"
}
