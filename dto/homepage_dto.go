package dto

// CreateHomepageSectionRequest represents request to create a homepage section
type CreateHomepageSectionRequest struct {
	SectionKey  string      `json:"section_key" binding:"required"`
	Title       string      `json:"title"`
	Subtitle    string      `json:"subtitle"`
	Description string      `json:"description"`
	Image       string      `json:"image"`
	Content     interface{} `json:"content"`
	IsActive    *bool       `json:"is_active"`
	SortOrder   int         `json:"sort_order"`
}

// UpdateHomepageSectionRequest represents request to update a homepage section
type UpdateHomepageSectionRequest struct {
	SectionKey  string      `json:"section_key"`
	Title       string      `json:"title"`
	Subtitle    string      `json:"subtitle"`
	Description string      `json:"description"`
	Image       string      `json:"image"`
	Content     interface{} `json:"content"`
	IsActive    *bool       `json:"is_active"`
	SortOrder   int         `json:"sort_order"`
}
