package payout

import (
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"
	"context"
	"errors"
	//"fmt"
	"time"

	"api-customer-merchant/internal/db"
	"github.com/shopspring/decimal"
)

type PayoutService struct {
	payoutRepo *repositories.PayoutRepository
}

func NewPayoutService(payoutRepo *repositories.PayoutRepository) *PayoutService {
	return &PayoutService{
		payoutRepo: payoutRepo,
	}
}

// CreatePayout creates a payout for a merchant
func (s *PayoutService) CreatePayout(merchantID string, amount float64) (*models.Payout, error) {
	if merchantID == "" {
		return nil, errors.New("invalid merchant ID")
	}
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	// Simulate payout processing (placeholder for Stripe)
	payout := &models.Payout{
		MerchantID: merchantID,
		Amount:     amount,
		Status:     models.PayoutStatusPending,
	}
	if err := s.payoutRepo.Create(context.Background(), payout); err != nil {
		return nil, err
	}

	// Simulate successful payout
	payout.Status = models.PayoutStatusCompleted
	if err := s.payoutRepo.Update(context.Background(), payout); err != nil {
		return nil, err
	}

	return s.payoutRepo.FindByID(context.Background(), payout.ID)
}

// GetPayoutByID retrieves a payout by ID
func (s *PayoutService) GetPayoutByID(id uint) (*models.Payout, error) {
	if id == 0 {
		return nil, errors.New("invalid payout ID")
	}
	return s.payoutRepo.FindByID(context.Background(), id)
}

// GetPayoutsByMerchantID retrieves all payouts for a merchant
func (s *PayoutService) GetPayoutsByMerchantID(merchantID string) ([]models.Payout, error) {
	if merchantID == "" {
		return nil, errors.New("invalid merchant ID")
	}
	return s.payoutRepo.FindByMerchantID(context.Background(), merchantID)
}

// RequestPayout requests a payout for a merchant
func (s *PayoutService) RequestPayout(ctx context.Context, merchantID string) (*models.Payout, error) {
    var sumStr string
    err := db.DB.Model(&models.OrderMerchantSplit{}).
        Where("merchant_id = ? AND status = 'pending' AND hold_until < ?", merchantID, time.Now()).
        Pluck("COALESCE(SUM(amount_due), '0')", &sumStr).Error
    if err != nil {
        return nil, err
    }

    totalDue, err := decimal.NewFromString(sumStr)
    if err != nil || totalDue.LessThanOrEqual(decimal.Zero) {
        return nil, errors.New("no eligible balance")
    }

    payout := &models.Payout{
        MerchantID: merchantID,
        Amount:     totalDue.InexactFloat64(),
        Status:    models.PayoutStatusPending,
    }

    if err := s.payoutRepo.Create(ctx, payout); err != nil {
        return nil, err
    }

    // Update splits
    db.DB.Model(&models.OrderMerchantSplit{}).
        Where("merchant_id = ? AND status = 'pending' AND hold_until < ?", merchantID, time.Now()).
        Update("status", "payout_requested")

    return payout, nil
}
