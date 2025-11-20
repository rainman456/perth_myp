// Create new file: internal/db/repositories/order_merchant_split_repository.go

package repositories

import (
	"context"
	//"log"

	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	
	"gorm.io/gorm"
)

type OrderMerchantSplitRepository struct {
	db *gorm.DB
}

func NewOrderMerchantSplitRepository() *OrderMerchantSplitRepository {
	return &OrderMerchantSplitRepository{db: db.DB}
}

// FindByOrderID retrieves all splits for an order
func (r *OrderMerchantSplitRepository) FindByOrderID(ctx context.Context, orderID uint) ([]models.OrderMerchantSplit, error) {
	var splits []models.OrderMerchantSplit
	err := r.db.WithContext(ctx).
		Preload("Merchant").
		Preload("Order").
		Where("order_id = ?", orderID).
		Find(&splits).Error
	return splits, err
}

// FindByMerchantID retrieves all splits for a merchant
func (r *OrderMerchantSplitRepository) FindByMerchantID(ctx context.Context, merchantID string) ([]models.OrderMerchantSplit, error) {
	var splits []models.OrderMerchantSplit
	err := r.db.WithContext(ctx).
		Preload("Order").
		Where("merchant_id = ?", merchantID).
		Find(&splits).Error
	return splits, err
}

// UpdateStatus updates the status of splits
func (r *OrderMerchantSplitRepository) UpdateStatus(ctx context.Context, orderID uint, oldStatus, newStatus models.OrderMerchantSplitStatus) error {
	return r.db.WithContext(ctx).
		Model(&models.OrderMerchantSplit{}).
		Where("order_id = ? AND status = ?", orderID, oldStatus).
		Update("status", newStatus).Error
}

// UpdateStatusByMerchantAndStatus updates splits status for a merchant
func (r *OrderMerchantSplitRepository) UpdateStatusByMerchantAndStatus(ctx context.Context, merchantID string, oldStatus, newStatus models.OrderMerchantSplitStatus) error {
	return r.db.WithContext(ctx).
		Model(&models.OrderMerchantSplit{}).
		Where("merchant_id = ? AND status = ?", merchantID, oldStatus).
		Update("status", newStatus).Error
}