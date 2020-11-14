package service

import (
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
	"github.com/zhashkevych/jewelry-shop-backend/pkg/repository"
)

type SettingsService struct {
	repo repository.Settings
}

func NewSettingsService(repo repository.Settings) *SettingsService {
	return &SettingsService{repo: repo}
}

func (s *SettingsService) GetSettings() (jewerly.Settings, error) {
	var settings jewerly.Settings

	images, err := s.GetImages()
	if err != nil {
		return settings, err
	}

	settings.Images = images

	textBlocks, err := s.GetTextBlocks()
	if err != nil {
		return settings, err
	}

	settings.TextBlocks = textBlocks

	return settings, nil
}

func (s *SettingsService) GetImages() ([]jewerly.HomepageImage, error) {
	return s.repo.GetImages()
}

func (s *SettingsService) CreateImage(imageID int) error {
	return s.repo.CreateImage(imageID)
}

func (s *SettingsService) UpdateImage(id, imageID int) error {
	return s.repo.UpdateImage(id, imageID)
}

func (s *SettingsService) GetTextBlocks() ([]jewerly.TextBlock, error) {
	return s.repo.GetTextBlocks()
}

func (s *SettingsService) GetTextBlockById(id int) (jewerly.TextBlock, error) {
	return s.repo.GetTextBlockById(id)
}

func (s *SettingsService) CreateTextBlock(block jewerly.TextBlock) error {
	return s.repo.CreateTextBlock(block)
}

func (s *SettingsService) UpdateTextBlock(id int, block jewerly.UpdateTextBlockInput) error {
	return s.repo.UpdateTextBlock(id, block)
}
