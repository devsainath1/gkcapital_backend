package models

import (
	"time"
)

// MediaAsset stores images as binary blobs in the database
type MediaAsset struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:255;not null;uniqueIndex" json:"name"`
	MimeType  string    `gorm:"size:100;not null" json:"mime_type"`
	Data      []byte    `gorm:"type:longblob;not null" json:"-"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (MediaAsset) TableName() string {
	return "media_assets"
}
