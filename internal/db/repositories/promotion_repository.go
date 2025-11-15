package repositories

import (
	"context"
	"time"

	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"

	"gorm.io/gorm"
)

type PromotionRepository struct {
	db *gorm.DB
}

func NewPromotionRepository() *PromotionRepository {
	return &PromotionRepository{db: db.DB}
}

// Create adds a new promotion
func (r *PromotionRepository) Create(ctx context.Context, promotion *models.Promotion) error {
	return r.db.WithContext(ctx).Create(promotion).Error
}

// FindByID retrieves a promotion by ID with associated Merchant and Products
func (r *PromotionRepository) FindByID(ctx context.Context, id string) (*models.Promotion, error) {
	var promotion models.Promotion
	err := r.db.WithContext(ctx).Preload("Merchant").Preload("Products").First(&promotion, "id = ?", id).Error
	return &promotion, err
}

// FindByMerchantID retrieves all promotions for a merchant
func (r *PromotionRepository) FindByMerchantID(ctx context.Context, merchantID string) ([]models.Promotion, error) {
	var promotions []models.Promotion
	err := r.db.WithContext(ctx).Preload("Merchant").Preload("Products").Where("merchant_id = ?", merchantID).Find(&promotions).Error
	return promotions, err
}

// FindActiveByMerchantID retrieves all active promotions for a merchant
func (r *PromotionRepository) FindActiveByMerchantID(ctx context.Context, merchantID string) ([]models.Promotion, error) {
	var promotions []models.Promotion
	err := r.db.WithContext(ctx).Preload("Merchant").Preload("Products").
		Where("merchant_id = ? AND status = ? AND start_date <= ? AND end_date >= ?", 
			merchantID, models.PromotionStatusActive, time.Now(), time.Now()).
		Find(&promotions).Error
	return promotions, err
}

// Update modifies an existing promotion
func (r *PromotionRepository) Update(ctx context.Context, promotion *models.Promotion) error {
	return r.db.WithContext(ctx).Save(promotion).Error
}

// Delete removes a promotion by ID
func (r *PromotionRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Promotion{}, "id = ?", id).Error
}

// AddProductsToPromotion associates products with a promotion
func (r *PromotionRepository) AddProductsToPromotion(ctx context.Context, promotionID string, productIDs []string) error {
	// First, remove existing associations
	if err := r.db.WithContext(ctx).Exec("DELETE FROM promotion_products WHERE promotion_id = ?", promotionID).Error; err != nil {
		return err
	}

	// Then add new associations
	if len(productIDs) > 0 {
		var associations []map[string]interface{}
		for _, productID := range productIDs {
			associations = append(associations, map[string]interface{}{
				"promotion_id": promotionID,
				"product_id":   productID,
			})
		}
		return r.db.WithContext(ctx).Table("promotion_products").Create(associations).Error
	}
	return nil
}

// RemoveProductsFromPromotion removes products from a promotion
func (r *PromotionRepository) RemoveProductsFromPromotion(ctx context.Context, promotionID string, productIDs []string) error {
	if len(productIDs) == 0 {
		return nil
	}
	
	return r.db.WithContext(ctx).Exec("DELETE FROM promotion_products WHERE promotion_id = ? AND product_id IN ?", 
		promotionID, productIDs).Error
}