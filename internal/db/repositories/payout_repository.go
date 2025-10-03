package repositories

import (
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"

	"gorm.io/gorm"
)

type PayoutRepository struct {
	db *gorm.DB
}

func NewPayoutRepository() *PayoutRepository {
	return &PayoutRepository{db: db.DB}
}

// Create adds a new payout record
func (r *PayoutRepository) Create(payout *models.Payout) error {
	return r.db.Create(payout).Error
}

// FindByID retrieves a payout by ID with associated Merchant
func (r *PayoutRepository) FindByID(id uint) (*models.Payout, error) {
	var payout models.Payout
	err := r.db.Preload("Merchant").First(&payout, id).Error
	return &payout, err
}

// FindByMerchantID retrieves all payouts for a merchant
func (r *PayoutRepository) FindByMerchantID(merchantID uint) ([]models.Payout, error) {
	var payouts []models.Payout
	err := r.db.Preload("Merchant").Where("merchant_id = ?", merchantID).Find(&payouts).Error
	return payouts, err
}

// Update modifies an existing payout
func (r *PayoutRepository) Update(payout *models.Payout) error {
	return r.db.Save(payout).Error
}

// Delete removes a payout by ID
func (r *PayoutRepository) Delete(id uint) error {
	return r.db.Delete(&models.Payout{}, id).Error
}
