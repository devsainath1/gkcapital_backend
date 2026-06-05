package dto

// MediaAssetResponse represents a media asset in API responses
type MediaAssetResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	MimeType string `json:"mime_type"`
	Size     int64  `json:"size"`
	URL      string `json:"url"`
}
