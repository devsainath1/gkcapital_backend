package dto

// CreateTestimonialRequest represents request to create a testimonial
type CreateTestimonialRequest struct {
	Name        string `json:"name" binding:"required"`
	Designation string `json:"designation"`
	Company     string `json:"company"`
	Content     string `json:"content" binding:"required"`
	Rating      int    `json:"rating"`
	Image       string `json:"image"`
	IsActive    *bool  `json:"is_active"`
	SortOrder   int    `json:"sort_order"`
}

// UpdateTestimonialRequest represents request to update a testimonial
type UpdateTestimonialRequest struct {
	Name        string `json:"name"`
	Designation string `json:"designation"`
	Company     string `json:"company"`
	Content     string `json:"content"`
	Rating      int    `json:"rating"`
	Image       string `json:"image"`
	IsActive    *bool  `json:"is_active"`
	SortOrder   int    `json:"sort_order"`
}
