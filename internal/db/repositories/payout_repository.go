package repositories

import (
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	"context"

	"gorm.io/gorm"
)

type PayoutRepository struct {
	db *gorm.DB
}

func NewPayoutRepository() *PayoutRepository {
	return &PayoutRepository{db: db.DB}
}

// Create adds a new payout record
func (r *PayoutRepository) Create(ctx context.Context,payout *models.Payout) error {
	return r.db.WithContext(ctx).Create(payout).Error
}

// FindByID retrieves a payout by ID with associated Merchant
func (r *PayoutRepository) FindByID(ctx context.Context,id uint) (*models.Payout, error) {
	var payout models.Payout
	err := r.db.WithContext(ctx).Preload("Merchant").First(&payout, id).Error
	return &payout, err
}

// FindByMerchantID retrieves all payouts for a merchant
func (r *PayoutRepository) FindByMerchantID(ctx context.Context,merchantID string) ([]models.Payout, error) {
	var payouts []models.Payout
	err := r.db.WithContext(ctx).Preload("Merchant").Where("merchant_id = ?", merchantID).Find(&payouts).Error
	return payouts, err
}

// Update modifies an existing payout
func (r *PayoutRepository) Update(ctx context.Context,payout *models.Payout) error {
	return r.db.WithContext(ctx).Save(payout).Error
}

// Delete removes a payout by ID
func (r *PayoutRepository) Delete(ctx context.Context,id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Payout{}, id).Error
}


func (r *PayoutRepository) FindByPaystackTransferID(ctx context.Context, transferID string) (*models.Payout, error) {
	var payout models.Payout
	err := r.db.WithContext(ctx).Where("paystack_transfer_id = ?", transferID).First(&payout).Error
	return &payout, err
}