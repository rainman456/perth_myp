package repositories

import (
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	"context"
	"gorm.io/gorm"
)

type SettingsRepository struct {
	db *gorm.DB
}

func NewSettingsRepository() *SettingsRepository {
	return &SettingsRepository{db: db.DB}
}

// GetSettings retrieves the global settings
func (r *SettingsRepository) GetSettings(ctx context.Context) (*models.Settings, error) {
	var settings models.Settings
	err := r.db.WithContext(ctx).Where("id = ?", "global").First(&settings).Error
	return &settings, err
}
