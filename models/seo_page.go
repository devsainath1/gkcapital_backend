package models

import (
	"time"
)

type SEOPage struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	PageSlug        string    `gorm:"size:255;uniqueIndex;not null" json:"page_slug"`
	PageName        string    `gorm:"size:255;not null" json:"page_name"`
	MetaTitle       string    `gorm:"size:500" json:"meta_title"`
	MetaDescription string    `gorm:"type:text" json:"meta_description"`
	MetaKeywords    string    `gorm:"type:text" json:"meta_keywords"`
	OGTitle         string    `gorm:"size:500" json:"og_title"`
	OGDescription   string    `gorm:"type:text" json:"og_description"`
	OGImage         string    `gorm:"size:500" json:"og_image"`
	CanonicalURL    string    `gorm:"size:500" json:"canonical_url"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (SEOPage) TableName() string {
	return "seo_pages"
}
