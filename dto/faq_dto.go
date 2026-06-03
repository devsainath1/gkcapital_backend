package dto

// CreateFAQRequest represents request to create a FAQ
type CreateFAQRequest struct {
	Question string `json:"question" binding:"required"`
	Answer   string `json:"answer" binding:"required"`
	Category string `json:"category"`
	IsActive *bool  `json:"is_active"`
	SortOrder int   `json:"sort_order"`
}

// UpdateFAQRequest represents request to update a FAQ
type UpdateFAQRequest struct {
	Question  string `json:"question"`
	Answer    string `json:"answer"`
	Category  string `json:"category"`
	IsActive  *bool  `json:"is_active"`
	SortOrder int    `json:"sort_order"`
}
