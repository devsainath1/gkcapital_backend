package services

import (
	"gk-capital-backend/dto"
	"gk-capital-backend/models"
	"gk-capital-backend/repository"
)

type WebsiteSettingService struct {
	settingRepo *repository.WebsiteSettingRepository
}

func NewWebsiteSettingService(settingRepo *repository.WebsiteSettingRepository) *WebsiteSettingService {
	return &WebsiteSettingService{settingRepo: settingRepo}
}

func (s *WebsiteSettingService) GetAll() ([]dto.WebsiteSettingResponse, error) {
	settings, err := s.settingRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var resp []dto.WebsiteSettingResponse
	for _, setting := range settings {
		resp = append(resp, dto.WebsiteSettingResponse{
			Key:   setting.Key,
			Value: setting.Value,
		})
	}

	if resp == nil {
		resp = []dto.WebsiteSettingResponse{}
	}

	return resp, nil
}

func (s *WebsiteSettingService) GetAllMap() (map[string]string, error) {
	settings, err := s.settingRepo.FindAll()
	if err != nil {
		return nil, err
	}

	resp := make(map[string]string)
	for _, setting := range settings {
		resp[setting.Key] = setting.Value
	}

	return resp, nil
}

func (s *WebsiteSettingService) BulkUpsert(req dto.UpdateWebsiteSettingsRequest) error {
	for _, setting := range req.Settings {
		ws := models.WebsiteSetting{
			Key:   setting.Key,
			Value: setting.Value,
		}
		if err := s.settingRepo.Upsert(&ws); err != nil {
			return err
		}
	}
	return nil
}
