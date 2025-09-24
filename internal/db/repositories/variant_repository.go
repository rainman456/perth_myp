package repositories

import (
	//"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"

	"context"

	"gorm.io/gorm"
)

type VariantRepository struct{ db *gorm.DB }

func NewVariantRepository(db *gorm.DB) *VariantRepository { return &VariantRepository{db} }
func (r *VariantRepository) FindByID(ctx context.Context, id string) (*models.Variant, error) {
	var variant models.Variant
	return &variant, r.db.WithContext(ctx).Preload("Inventory").First(&variant, "id = ?", id).Error
}
