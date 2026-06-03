package models

import (
	"time"
)

type HomepageSection struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	SectionKey  string      `gorm:"size:255;not null" json:"section_key"`
	Title       string      `gorm:"size:500" json:"title"`
	Subtitle    string      `gorm:"size:500" json:"subtitle"`
	Description string      `gorm:"type:text" json:"description"`
	Image       string      `gorm:"size:500" json:"image"`
	Content     interface{} `gorm:"serializer:json" json:"content"`
	IsActive    bool        `gorm:"default:true" json:"is_active"`
	SortOrder   int         `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

func (HomepageSection) TableName() string {
	return "homepage_sections"
}
