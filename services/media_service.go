package services

import (
	"fmt"
	"gk-capital-backend/models"
	"gk-capital-backend/repository"
	"sync"
)

type MediaService struct {
	mediaRepo *repository.MediaRepository
	cacheMu   sync.RWMutex
	nameCache map[string]*models.MediaAsset
	idCache   map[uint]*models.MediaAsset
}

func NewMediaService(mediaRepo *repository.MediaRepository) *MediaService {
	return &MediaService{
		mediaRepo: mediaRepo,
		nameCache: make(map[string]*models.MediaAsset),
		idCache:   make(map[uint]*models.MediaAsset),
	}
}

// allowedMimeTypes lists the accepted image content types
var allowedMimeTypes = map[string]bool{
	"image/png":     true,
	"image/jpeg":    true,
	"image/jpg":     true,
	"image/gif":     true,
	"image/webp":    true,
	"image/svg+xml": true,
}

const maxUploadSize = 10 * 1024 * 1024 // 10 MB

func (s *MediaService) clearCache() {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()
	s.nameCache = make(map[string]*models.MediaAsset)
	s.idCache = make(map[uint]*models.MediaAsset)
}

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

	s.clearCache()
	return asset, nil
}

func (s *MediaService) GetByID(id uint) (*models.MediaAsset, error) {
	s.cacheMu.RLock()
	asset, found := s.idCache[id]
	s.cacheMu.RUnlock()
	if found {
		return asset, nil
	}

	asset, err := s.mediaRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	s.cacheMu.Lock()
	s.idCache[id] = asset
	s.nameCache[asset.Name] = asset
	s.cacheMu.Unlock()

	return asset, nil
}

func (s *MediaService) GetByName(name string) (*models.MediaAsset, error) {
	s.cacheMu.RLock()
	asset, found := s.nameCache[name]
	s.cacheMu.RUnlock()
	if found {
		return asset, nil
	}

	asset, err := s.mediaRepo.GetByName(name)
	if err != nil {
		return nil, err
	}

	s.cacheMu.Lock()
	s.idCache[asset.ID] = asset
	s.nameCache[name] = asset
	s.cacheMu.Unlock()

	return asset, nil
}

func (s *MediaService) GetAll() ([]models.MediaAsset, error) {
	// GetAll doesn't return raw binary data, so no need to cache this query
	return s.mediaRepo.GetAll()
}

func (s *MediaService) Delete(id uint) error {
	if err := s.mediaRepo.Delete(id); err != nil {
		return err
	}
	s.clearCache()
	return nil
}

func (s *MediaService) ExistsByName(name string) bool {
	s.cacheMu.RLock()
	_, found := s.nameCache[name]
	s.cacheMu.RUnlock()
	if found {
		return true
	}
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
	s.clearCache()
	return asset, nil
}
