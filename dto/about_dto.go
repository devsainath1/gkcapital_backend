package dto

// CreateAboutSectionRequest represents request to create an about section
type CreateAboutSectionRequest struct {
	SectionKey  string      `json:"section_key" binding:"required"`
	Title       string      `json:"title"`
	Subtitle    string      `json:"subtitle"`
	Description string      `json:"description"`
	Image       string      `json:"image"`
	Content     interface{} `json:"content"`
	IsActive    *bool       `json:"is_active"`
	SortOrder   int         `json:"sort_order"`
}

// UpdateAboutSectionRequest represents request to update an about section
type UpdateAboutSectionRequest struct {
	SectionKey  string      `json:"section_key"`
	Title       string      `json:"title"`
	Subtitle    string      `json:"subtitle"`
	Description string      `json:"description"`
	Image       string      `json:"image"`
	Content     interface{} `json:"content"`
	IsActive    *bool       `json:"is_active"`
	SortOrder   int         `json:"sort_order"`
}
