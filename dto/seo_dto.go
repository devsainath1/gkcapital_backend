package dto

// CreateSEOPageRequest represents request to create an SEO page entry
type CreateSEOPageRequest struct {
	PageSlug        string `json:"page_slug" binding:"required"`
	MetaTitle       string `json:"meta_title"`
	MetaDescription string `json:"meta_description"`
	MetaKeywords    string `json:"meta_keywords"`
	OGTitle         string `json:"og_title"`
	OGDescription   string `json:"og_description"`
	OGImage         string `json:"og_image"`
	CanonicalURL    string `json:"canonical_url"`
}

// UpdateSEOPageRequest represents request to update an SEO page entry
type UpdateSEOPageRequest struct {
	PageSlug        string `json:"page_slug"`
	MetaTitle       string `json:"meta_title"`
	MetaDescription string `json:"meta_description"`
	MetaKeywords    string `json:"meta_keywords"`
	OGTitle         string `json:"og_title"`
	OGDescription   string `json:"og_description"`
	OGImage         string `json:"og_image"`
	CanonicalURL    string `json:"canonical_url"`
}
