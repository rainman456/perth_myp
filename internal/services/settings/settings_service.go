package settings

import (
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"
	"context"
	"encoding/json"
	"errors"
	"strconv"
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

	var shippingOptions map[string]interface{}
	if err := json.Unmarshal(settings.ShippingOptions, &shippingOptions); err != nil {
		return 0, err
	}

	priceVal, ok := shippingOptions[shippingMethod]
	if !ok {
		return 0, errors.New("invalid shipping method")
	}

	var price float64
	switch v := priceVal.(type) {
	case float64:
		price = v
	case string:
		price, err = strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, err
		}
	default:
		return 0, errors.New("invalid price type")
	}

	return price, nil
}
// UpdateSettings updates the global settings


func (s *SettingsService) GetPlatformFee(ctx context.Context) (float64, error) {
	settings, err := s.settingsRepo.GetSettings(ctx)
	if err != nil {
		return 0, err
	}
	return settings.Fees, nil
}