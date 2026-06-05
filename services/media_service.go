package services

import (
	"fmt"
	"gk-capital-backend/models"
	"gk-capital-backend/repository"
)

type MediaService struct {
	mediaRepo *repository.MediaRepository
}

func NewMediaService(mediaRepo *repository.MediaRepository) *MediaService {
	return &MediaService{mediaRepo: mediaRepo}
}

// allowedMimeTypes lists the accepted image content types
var allowedMimeTypes = map[string]bool{
	"image/png":  true,
	"image/jpeg": true,
	"image/jpg":  true,
	"image/gif":  true,
	"image/webp": true,
	"image/svg+xml": true,
}

const maxUploadSize = 10 * 1024 * 1024 // 10 MB

func (s *MediaService) Upload(name string, mimeType string, data []byte) (*models.MediaAsset, error) {
	if !allowedMimeTypes[mimeType] {
		return nil, fmt.Errorf("unsupported file type: %s", mimeType)
	}
	if len(data) > maxUploadSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size of 10MB")
	}

	asset := &models.MediaAsset{
		Name:     name,
		MimeType: mimeType,
		Data:     data,
		Size:     int64(len(data)),
	}

	if err := s.mediaRepo.Create(asset); err != nil {
		return nil, err
	}
	return asset, nil
}

func (s *MediaService) GetByID(id uint) (*models.MediaAsset, error) {
	return s.mediaRepo.GetByID(id)
}

func (s *MediaService) GetByName(name string) (*models.MediaAsset, error) {
	return s.mediaRepo.GetByName(name)
}

func (s *MediaService) GetAll() ([]models.MediaAsset, error) {
	return s.mediaRepo.GetAll()
}

func (s *MediaService) Delete(id uint) error {
	return s.mediaRepo.Delete(id)
}

func (s *MediaService) ExistsByName(name string) bool {
	return s.mediaRepo.ExistsByName(name)
}

// CreateFromBytes directly stores binary data (used by seeder)
func (s *MediaService) CreateFromBytes(name string, mimeType string, data []byte) (*models.MediaAsset, error) {
	asset := &models.MediaAsset{
		Name:     name,
		MimeType: mimeType,
		Data:     data,
		Size:     int64(len(data)),
	}
	if err := s.mediaRepo.Create(asset); err != nil {
		return nil, err
	}
	return asset, nil
}
