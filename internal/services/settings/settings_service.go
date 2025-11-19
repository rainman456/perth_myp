package settings

import (
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"
	"context"
	"encoding/json"
	"errors"
)

type SettingsService struct {
	settingsRepo *repositories.SettingsRepository
}

func NewSettingsService(settingsRepo *repositories.SettingsRepository) *SettingsService {
	return &SettingsService{
		settingsRepo: settingsRepo,
	}
}

// GetSettings retrieves the global settings
func (s *SettingsService) GetSettings(ctx context.Context) (*models.Settings, error) {
	return s.settingsRepo.GetSettings(ctx)
}

// GetShippingCost calculates shipping cost based on shipping method
func (s *SettingsService) GetShippingCost(ctx context.Context, shippingMethod string) (float64, error) {
	settings, err := s.settingsRepo.GetSettings(ctx)
	if err != nil {
		return 0, err
	}

	var shippingOptions []models.ShippingOption
	if err := json.Unmarshal(settings.ShippingOptions, &shippingOptions); err != nil {
		return 0, err
	}

	for _, option := range shippingOptions {
		if option.Name == shippingMethod && option.Enabled {
			return option.Price, nil
		}
	}

	return 0, errors.New("invalid or disabled shipping method")
}

// UpdateSettings updates the global settings


func (s *SettingsService) GetPlatformFee(ctx context.Context) (float64, error) {
	settings, err := s.settingsRepo.GetSettings(ctx)
	if err != nil {
		return 0, err
	}
	return settings.Fees, nil
}