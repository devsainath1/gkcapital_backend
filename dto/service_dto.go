package dto

// CreateServiceRequest represents request to create a service
type CreateServiceRequest struct {
	Title       string `json:"title" binding:"required"`
	Slug        string `json:"slug"`
	Description string `json:"description" binding:"required"`
	Image       string `json:"image"`
	Icon        string `json:"icon"`
	IsActive    *bool  `json:"is_active"`
	SortOrder   int    `json:"sort_order"`
}

// UpdateServiceRequest represents request to update a service
type UpdateServiceRequest struct {
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Icon        string `json:"icon"`
	IsActive    *bool  `json:"is_active"`
	SortOrder   int    `json:"sort_order"`
}

// ServiceResponse represents service data in responses
type ServiceResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Icon        string `json:"icon"`
	IsActive    bool   `json:"is_active"`
	SortOrder   int    `json:"sort_order"`
}
