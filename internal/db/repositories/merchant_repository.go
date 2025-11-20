package repositories

import (
	"context"
	"fmt"
	"log"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"

	"github.com/shopspring/decimal"
	//"gorm.io/gorm"
)

// MerchantApplicationRepository handles CRUD for merchant applications
// Note: Admin service (in Express) will be responsible for updating status/approval.
type MerchantApplicationRepository struct{}

func NewMerchantApplicationRepository() *MerchantApplicationRepository {
	return &MerchantApplicationRepository{}
}

func (r *MerchantApplicationRepository) Create(ctx context.Context, m *models.MerchantApplication) error {
	err := db.DB.WithContext(ctx).Create(m).Error
	if err != nil {
		log.Printf("Failed to create merchant application: %v", err)
		return err
	}
	return nil
}

func (r *MerchantApplicationRepository) GetByID(ctx context.Context, id string) (*models.MerchantApplication, error) {
	var m models.MerchantApplication
	if err := db.DB.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		log.Printf("Failed to get merchant application by ID %s: %v", id, err)
		return nil, err
	}
	return &m, nil
}

func (r *MerchantApplicationRepository) GetByUserEmail(ctx context.Context, email string) (*models.MerchantApplication, error) {
	var m models.MerchantApplication
	if err := db.DB.WithContext(ctx).Where("personal_email = ? OR work_email = ?", email, email).First(&m).Error; err != nil {
		log.Printf("Failed to get merchant application by email %s: %v", email, err)
		return nil, err
	}
	return &m, nil
}

// MerchantRepository handles active merchants
type MerchantRepository struct{}

func NewMerchantRepository() *MerchantRepository {
	return &MerchantRepository{}
}

func (r *MerchantRepository) GetByID(ctx context.Context, id string) (*models.Merchant, error) {
	var m models.Merchant
	if err := db.DB.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		log.Printf("Failed to get merchant by ID %s: %v", id, err)
		return nil, err
	}
	return &m, nil
}

func (r *MerchantRepository) GetByMerchantID(ctx context.Context, uid string) (*models.Merchant, error) {
	var m models.Merchant
	if err := db.DB.WithContext(ctx).Where("merchant_id = ?", uid).First(&m).Error; err != nil {
		log.Printf("Failed to get merchant by user ID %s: %v", uid, err)
		return nil, err
	}
	return &m, nil
}

func (r *MerchantRepository) GetByWorkEmail(ctx context.Context, email string) (*models.Merchant, error) {
	var m models.Merchant
	if err := db.DB.WithContext(ctx).Where("personal_email = ? OR work_email = ?", email, email).First(&m).Error; err != nil {
		log.Printf("Failed to get merchant  by email %s: %v", email, err)
		return nil, err
	}
	return &m, nil
}








// func (r *MerchantRepository) UpdateBankDetails(ctx context.Context, merchantID string, details dto.BankDetailsRequest) error {
// 	// Use WithContext so DB operations respect request lifecycle
// 	if err := db.DB.WithContext(ctx).
// 		Model(&models.MerchantBankDetails{}).
// 		Where("merchant_id = ?", merchantID).
// 		Save(details).Error; err != nil {

// 		log.Printf("Failed to update bank details for merchant %s: %v", merchantID, err)
// 		return err
// 	}
// 	return nil
// }


















// Add these methods to your MerchantRepository in merchant_repository.go

// CreateBankDetails creates new bank details for a merchant
func (r *MerchantRepository) CreateBankDetails(ctx context.Context, merchantID string, details dto.BankDetailsRequest) (*models.MerchantBankDetails, error) {
	bankDetails := &models.MerchantBankDetails{
		MerchantID:    merchantID,
		BankName:      details.BankName,
		BankCode:      details.BankCode,
		AccountNumber: details.AccountNumber,
		AccountName:   details.AccountName,
		Currency:      details.Currency,
		//Status:        "pending",
	}

	if err := db.DB.WithContext(ctx).Create(bankDetails).Error; err != nil {
		log.Printf("Failed to create bank details for merchant %s: %v", merchantID, err)
		return nil, err
	}

	return bankDetails, nil
}

// GetBankDetails retrieves bank details for a merchant
func (r *MerchantRepository) GetBankDetails(ctx context.Context, merchantID string) (*models.MerchantBankDetails, error) {
	var bankDetails models.MerchantBankDetails
	if err := db.DB.WithContext(ctx).
		Where("merchant_id = ?", merchantID).
		First(&bankDetails).Error; err != nil {
		log.Printf("Failed to get bank details for merchant %s: %v", merchantID, err)
		return nil, err
	}
	return &bankDetails, nil
}

// UpdateBankDetailsRecord updates existing bank details for a merchant
func (r *MerchantRepository) UpdateBankDetailsRecord(ctx context.Context, merchantID string, details dto.BankDetailsRequest) (*models.MerchantBankDetails, error) {
	var bankDetails models.MerchantBankDetails
	
	// First, get the existing record
	if err := db.DB.WithContext(ctx).
		Where("merchant_id = ?", merchantID).
		First(&bankDetails).Error; err != nil {
		log.Printf("Failed to find bank details for merchant %s: %v", merchantID, err)
		return nil, err
	}

	// Update the fields
	updates := map[string]interface{}{
		"bank_name":      details.BankName,
		"bank_code":      details.BankCode,
		"account_number": details.AccountNumber,
		"account_name":   details.AccountName,
		"currency":       details.Currency,
		//"status":         "pending", // Reset to pending on update
	}

	if err := db.DB.WithContext(ctx).
		Model(&bankDetails).
		Updates(updates).Error; err != nil {
		log.Printf("Failed to update bank details for merchant %s: %v", merchantID, err)
		return nil, err
	}

	// Fetch the updated record
	if err := db.DB.WithContext(ctx).
		Where("merchant_id = ?", merchantID).
		First(&bankDetails).Error; err != nil {
		log.Printf("Failed to fetch updated bank details for merchant %s: %v", merchantID, err)
		return nil, err
	}

	return &bankDetails, nil
}

// DeleteBankDetails removes bank details for a merchant
func (r *MerchantRepository) DeleteBankDetails(ctx context.Context, merchantID string) error {
	if err := db.DB.WithContext(ctx).
		Where("merchant_id = ?", merchantID).
		Delete(&models.MerchantBankDetails{}).Error; err != nil {
		log.Printf("Failed to delete bank details for merchant %s: %v", merchantID, err)
		return err
	}
	return nil
}


















// UpdateMerchant updates a merchant's profile information
func (r *MerchantRepository) UpdateMerchant(ctx context.Context, merchantID string, updates map[string]interface{}) error {
	// Remove any fields that shouldn't be updated
	delete(updates, "id")
	delete(updates, "merchant_id")
	delete(updates, "application_id")
	delete(updates, "password")
	delete(updates, "status")
	delete(updates, "commission_tier")
	delete(updates, "commission_rate")
	delete(updates, "account_balance")
	delete(updates, "total_sales")
	delete(updates, "total_payouts")
	delete(updates, "payout_schedule")
	delete(updates, "last_payout_date")
	delete(updates, "created_at")
	delete(updates, "updated_at")

	if err := db.DB.WithContext(ctx).
		Model(&models.Merchant{}).
		Where("merchant_id = ?", merchantID).
		Updates(updates).Error; err != nil {
		log.Printf("Failed to update merchant %s: %v", merchantID, err)
		return err
	}
	return nil
}
























func (r *MerchantRepository) CalculateTotalSales(ctx context.Context, merchantID string) (float64, error) {
	var totalSales decimal.Decimal
	
	err := db.DB.WithContext(ctx).
		Model(&models.OrderMerchantSplit{}).
		Select("COALESCE(SUM(amount_due), 0)").
		Joins("JOIN orders ON orders.id = order_merchant_splits.order_id").
		Where("order_merchant_splits.merchant_id = ? AND orders.status = ?", 
			merchantID, models.OrderStatusCompleted).
		Scan(&totalSales).Error
	
	if err != nil {
		log.Printf("Failed to calculate total sales for merchant %s: %v", merchantID, err)
		return 0, err
	}
	
	return totalSales.InexactFloat64(), nil
}

// CalculateTotalPayouts calculates total completed payouts
func (r *MerchantRepository) CalculateTotalPayouts(ctx context.Context, merchantID string) (float64, error) {
	var totalPayouts float64
	
	err := db.DB.WithContext(ctx).
		Model(&models.Payout{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("merchant_id = ? AND status = ?", merchantID, models.PayoutStatusCompleted).
		Scan(&totalPayouts).Error
	
	if err != nil {
		log.Printf("Failed to calculate total payouts for merchant %s: %v", merchantID, err)
		return 0, err
	}
	
	return totalPayouts, nil
}

// UpdateMerchantFinancials updates merchant's financial fields
func (r *MerchantRepository) UpdateMerchantFinancials(ctx context.Context, merchantID string) error {
	totalSales, err := r.CalculateTotalSales(ctx, merchantID)
	if err != nil {
		return fmt.Errorf("failed to calculate total sales: %w", err)
	}
	
	totalPayouts, err := r.CalculateTotalPayouts(ctx, merchantID)
	if err != nil {
		return fmt.Errorf("failed to calculate total payouts: %w", err)
	}
	
	// Get merchant to access commission_rate (platform fee)
	var merchant models.Merchant
	if err := db.DB.WithContext(ctx).
		Where("merchant_id = ?", merchantID).
		First(&merchant).Error; err != nil {
		log.Printf("Failed to get merchant %s: %v", merchantID, err)
		return fmt.Errorf("failed to get merchant: %w", err)
	}
	
	// Calculate platform fee from total sales using merchant's commission rate
	platformFee := totalSales * (merchant.CommissionRate / 100)
	
	// Account balance = total sales - platform fee
	accountBalance := totalSales - platformFee
	
	updates := map[string]interface{}{
		"total_sales":     totalSales,
		"total_payouts":   totalPayouts,
		"account_balance": accountBalance,
	}
	
	if err := db.DB.WithContext(ctx).
		Model(&models.Merchant{}).
		Where("merchant_id = ?", merchantID).
		Updates(updates).Error; err != nil {
		log.Printf("Failed to update merchant financials for %s: %v", merchantID, err)
		return err
	}
	
	return nil
}