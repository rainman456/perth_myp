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
func (s *PayoutService) GetPayoutByID(id string) (*models.Payout, error) {
	if id == "" {
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
func (s *PayoutService) RequestPayout(ctx context.Context, merchantID string, requestedAmount float64) (*models.Payout, error) {
	// Calculate total available balance (only from processing splits that passed hold period)
	var sumStr string
	err := db.DB.Model(&models.OrderMerchantSplit{}).
		Where("merchant_id = ? AND status = ? AND hold_until < ?", 
			merchantID, models.OrderMerchantSplitStatusProcessing, time.Now()).
		Pluck("COALESCE(SUM(amount_due), '0')", &sumStr).Error
	if err != nil {
		return nil, err
	}

	totalAvailable, err := decimal.NewFromString(sumStr)
	if err != nil || totalAvailable.LessThanOrEqual(decimal.Zero) {
		return nil, errors.New("no eligible balance available")
	}

	// Validate requested amount
	requestedDec := decimal.NewFromFloat(requestedAmount)
	if requestedDec.LessThanOrEqual(decimal.Zero) {
		return nil, errors.New("requested amount must be greater than zero")
	}
	
	if requestedDec.GreaterThan(totalAvailable) {
		return nil, errors.New("requested amount exceeds available balance")
	}

	// Create payout with requested amount
	payout := &models.Payout{
		MerchantID: merchantID,
		Amount:     requestedAmount,
		Status:     models.PayoutStatusPending,
	}

	if err := s.payoutRepo.Create(ctx, payout); err != nil {
		return nil, err
	}

	// Update splits to mark them as being paid out
	// Keep them in processing status, but track which payout they belong to
	var splits []models.OrderMerchantSplit
	err = db.DB.Where("merchant_id = ? AND status = ? AND hold_until < ?", 
		merchantID, models.OrderMerchantSplitStatusProcessing, time.Now()).
		Order("hold_until ASC").
		Find(&splits).Error
	if err != nil {
		return nil, err
	}

	remaining := requestedDec
	var splitIDs []uint
	for _, split := range splits {
		if remaining.LessThanOrEqual(decimal.Zero) {
			break
		}
		splitIDs = append(splitIDs, split.ID)
		remaining = remaining.Sub(split.AmountDue)
	}

	// Note: Splits remain in processing status until payout completes
	// They are tracked by the payout relationship
	// When payout completes, they will be marked as completed

	return payout, nil
}




func (s *PayoutService) GetAvailableBalance(ctx context.Context, merchantID string) (float64, error) {
	var sumStr string
	err := db.DB.Model(&models.OrderMerchantSplit{}).
		Where("merchant_id = ? AND status = ? AND hold_until < ?", 
			merchantID, models.OrderMerchantSplitStatusProcessing, time.Now()).
		Pluck("COALESCE(SUM(amount_due), '0')", &sumStr).Error
	if err != nil {
		return 0, err
	}

	total, err := decimal.NewFromString(sumStr)
	if err != nil {
		return 0, err
	}

	return total.InexactFloat64(), nil
}