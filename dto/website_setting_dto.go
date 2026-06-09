package dto

type WebsiteSettingResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type UpdateWebsiteSettingsRequest struct {
	Settings []WebsiteSettingResponse `json:"settings" binding:"required,dive"`
}
