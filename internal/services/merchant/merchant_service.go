package merchant

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/bank"
	"api-customer-merchant/internal/utils"

	//"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	//"github.com/gray-adeyi/paystack"

	"golang.org/x/crypto/bcrypt"
)

/*
	type MerchantService struct {
		appRepo  *repositories.MerchantApplicationRepository
		repo     *repositories.MerchantRepository
		validate *validator.Validate
	}

	func NewMerchantService(appRepo *repositories.MerchantApplicationRepository, repo *repositories.MerchantRepository) *MerchantService {
		return &MerchantService{
			appRepo:  appRepo,
			repo:     repo,
			validate: validator.New(),
		}
	}

// SubmitApplication allows a prospective merchant to submit an application.

	func (s *MerchantService) SubmitApplication(ctx context.Context, app *models.MerchantApplication) (*models.MerchantApplication, error) {
		// Validate required fields
		if err := s.validate.Struct(app); err != nil {
			return nil, errors.New("validation failed: " + err.Error())
		}

		// Validate JSONB fields
		if len(app.PersonalAddress) == 0 {
			return nil, errors.New("personal_address cannot be empty")
		}
		if len(app.WorkAddress) == 0 {
			return nil, errors.New("work_address cannot be empty")
		}
		var temp map[string]interface{}
		if err := json.Unmarshal(app.PersonalAddress, &temp); err != nil {
			return nil, errors.New("invalid personal_address JSON")
		}
		if err := json.Unmarshal(app.WorkAddress, &temp); err != nil {
			return nil, errors.New("invalid work_address JSON")
		}




		// Set Status to pending and ensure ID is not set
		app.Status = "pending"
		app.ID = ""

		if err := s.appRepo.Create(ctx, app); err != nil {
			return nil, err
		}
		return app, nil
	}

// GetApplication returns an application by ID (for applicant to check their own status).

	func (s *MerchantService) GetApplication(ctx context.Context, id string) (*models.MerchantApplication, error) {
		if id == "" {
			return nil, errors.New("application ID cannot be empty")
		}
		return s.appRepo.GetByID(ctx, id)
	}

// GetMerchantByUserID returns an active merchant account for the authenticated user.

	func (s *MerchantService) GetMerchantByUserID(ctx context.Context, uid string) (*models.Merchant, error) {
		if uid == "" {
			return nil, errors.New("user ID cannot be empty")
		}
		return s.repo.GetByUserID(ctx, uid)
	}

// GetMerchantByID returns an active merchant by ID.

	func (s *MerchantService) GetMerchantByID(ctx context.Context, id string) (*models.Merchant, error) {
		if id == "" {
			return nil, errors.New("merchant ID cannot be empty")
		}
		return s.repo.GetByID(ctx, id)
	}
*/
type MerchantService struct {
	appRepo  *repositories.MerchantApplicationRepository
	repo     *repositories.MerchantRepository
	validate *validator.Validate
}

func NewMerchantService(appRepo *repositories.MerchantApplicationRepository, repo *repositories.MerchantRepository) *MerchantService {
	return &MerchantService{
		appRepo:  appRepo,
		repo:     repo,
		validate: validator.New(),
	}
}

// SubmitApplication allows a prospective merchant to submit an application.
func (s *MerchantService) SubmitApplication(ctx context.Context, app *models.MerchantApplication) (*models.MerchantApplication, error) {
	// Validate required fields
	if err := s.validate.Struct(app); err != nil {
		return nil, errors.New("validation failed: " + err.Error())
	}

	// Validate JSONB fields
	if len(app.PersonalAddress) == 0 {
		return nil, errors.New("personal_address cannot be empty")
	}
	if len(app.WorkAddress) == 0 {
		return nil, errors.New("work_address cannot be empty")
	}
	var temp map[string]interface{}
	if err := json.Unmarshal(app.PersonalAddress, &temp); err != nil {
		return nil, errors.New("invalid personal_address JSON")
	}
	if err := json.Unmarshal(app.WorkAddress, &temp); err != nil {
		return nil, errors.New("invalid work_address JSON")
	}

	// Set Status to pending and ensure ID is not set
	app.Status = "pending"
	app.ID = ""

	if err := s.appRepo.Create(ctx, app); err != nil {
		return nil, err
	}
	return app, nil
}

// GetApplication returns an application by ID (for applicant to check their own status).
func (s *MerchantService) GetApplication(ctx context.Context, id string) (*models.MerchantApplication, error) {
	if id == "" {
		return nil, errors.New("application ID cannot be empty")
	}
	return s.appRepo.GetByID(ctx, id)
}

// GetMerchantByUserID returns an active merchant account for the authenticated user.
func (s *MerchantService) GetMerchantByUserID(ctx context.Context, uid string) (*models.Merchant, error) {
	if uid == "" {
		return nil, errors.New("user ID cannot be empty")
	}
	return s.repo.GetByMerchantID(ctx, uid)
}

// GetMerchantByID returns an active merchant by ID.
func (s *MerchantService) GetMerchantByID(ctx context.Context, id string) (*models.Merchant, error) {
	if id == "" {
		return nil, errors.New("merchant ID cannot be empty")
	}
	return s.repo.GetByMerchantID(ctx, id)
}

func (s *MerchantService) LoginMerchant(ctx context.Context, work_email, password string) (*models.Merchant, error) {
	merchant, err := s.repo.GetByWorkEmail(ctx, work_email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(merchant.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return merchant, nil
}

func (s *MerchantService) GenerateJWT(entity interface{}) (string, error) {
	var id string
	var entityType string

	switch e := entity.(type) {
	case *models.Merchant:
		id = e.MerchantID
		entityType = "merchant"

	default:
		return "", errors.New("invalid entity type")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         id,
		"entityType": entityType,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET not set")
	}

	return token.SignedString([]byte(secret))
}

//func (s *MerchantService) AddBankDetails(merchantID string, details MerchantBankDetails) error {
// Validate inputs (e.g., bank code format)
//  details.MerchantID = merchantID
//     details.Status = "pending"

//     // Create Paystack recipient
//     verify_client := paystack.VerificationClient(config.PaystackSecretKey)  // Assume injected
// 	var response models.Response[models.BankAccountInfo]
// 	if err := client.Verification.ResolveAccount(context.TODO(), &response,p.WithQuery("account_number","0022728151"),p.WithQuery("bank_code","063")); err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(response)
// }
//     recipientReq := &verifyclient. {
//         Type:          "nuban",
//         Name:          details.AccountName,
//         AccountNumber: details.AccountNumber,
//         BankCode:      details.BankCode,
//         Currency:      details.Currency,
//     }
//     resp, err := paystack.Recipient.Create(recipientReq)
//     if err != nil {
//         return err
//     }
//     details.RecipientCode = resp.Data.RecipientCode
//     details.Status = "verified"  // If Paystack verifies

//     return db.DB.Create(&details).Error
// }

// func (s *MerchantService) AddBankDetails(ctx context.Context, merchantID string, details dto.BankDetailsRequest) error {
// 	// Validate bank name
// 	bankSvc := bank.GetBankService()
// 	if err := bankSvc.LoadBanks(); err != nil {
// 		return fmt.Errorf("failed to load banks: %w", err)
// 	}

// 	bankCode, err := bankSvc.GetBankCode(details.BankName)
// 	if err != nil {
// 		return fmt.Errorf("invalid bank name: %w", err)
// 	}

// 	// Override with validated code
// 	details.BankCode = bankCode

// 	// Persist via repository
// 	if err := s.repo.UpdateBankDetails(ctx, merchantID, details); err != nil {
// 		return fmt.Errorf("failed to save bank details: %w", err)
// 	}

// 	return nil
// }




// func (s *MerchantService) UpdateBankDetails(ctx context.Context, merchantID string, details dto.BankDetailsRequest) error {
// 	// Similar, but use Save or Update
// 	if details.BankName == "" {
// 		return errors.New("empty bank name")
// 	}

// 	if details.AccountNumber == "" {
// 		return errors.New("empty bank name")
// 	}

// 	err := s.repo.UpdateBankDetails(ctx, merchantID, details)
// 	if err != nil {
// 		return err
// 	}

// 	//payment.Status = models.PaymentStatus(status)

// 	return nil

// }




















// GeneratePasswordResetToken generates a secure reset token and stores it in Redis
func (s *MerchantService) GeneratePasswordResetToken(email string) (string, time.Time, error) {
	// Verify merchant exists
	merchant, err := s.repo.GetByWorkEmail(context.Background(), email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// For security, don't reveal if email exists
			return "", time.Time{}, nil
		}
		return "", time.Time{}, err
	}

	// Generate secure random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", time.Time{}, err
	}
	token := base64.URLEncoding.EncodeToString(tokenBytes)

	// Store token in Redis with 1-hour expiration
	expiresAt := time.Now().Add(1 * time.Hour)
	ctx := context.Background()
	
	if utils.RedisClient != nil {
		key := "merchant_password_reset:" + token
		err = utils.RedisClient.Set(ctx, key, merchant.WorkEmail, 1*time.Hour).Err()
		if err != nil {
			log.Printf("Failed to store reset token in Redis: %v", err)
			return "", time.Time{}, err
		}
	} else {
		return "", time.Time{}, errors.New("redis not available")
	}

	return token, expiresAt, nil
}

// ResetPasswordWithToken validates the token and resets the password
func (s *MerchantService) ResetPasswordWithToken(token, newPassword string) error {
	if token == "" || newPassword == "" {
		return errors.New("token and password are required")
	}

	// Retrieve email from Redis using token
	ctx := context.Background()
	key := "merchant_password_reset:" + token
	
	if utils.RedisClient == nil {
		return errors.New("redis not available")
	}

	email, err := utils.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return errors.New("invalid or expired reset token")
	}

	// Delete token from Redis (one-time use)
	utils.RedisClient.Del(ctx, key)

	// Get merchant and reset password
	merchant, err := s.repo.GetByWorkEmail(ctx, email)
	if err != nil {
		return err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password using repository
	updates := map[string]interface{}{
		"password": string(hashed),
	}
	
	return s.repo.UpdateMerchant(ctx, merchant.MerchantID, updates)
}























































// Add these methods to your MerchantService in merchant_service.go

// AddBankDetails creates bank account details for a merchant
func (s *MerchantService) AddBankDetails(ctx context.Context, merchantID string, details dto.BankDetailsRequest) (*dto.BankDetailsResponse, error) {
	// Validate bank name
	bankSvc := bank.GetBankService()
	if err := bankSvc.LoadBanks(); err != nil {
		return nil, fmt.Errorf("failed to load banks: %w", err)
	}

	bankCode, err := bankSvc.GetBankCode(details.BankName)
	if err != nil {
		return nil, fmt.Errorf("invalid bank name: %w", err)
	}

	// Override with validated code
	details.BankCode = bankCode

	// Create bank details via repository
	bankDetails, err := s.repo.CreateBankDetails(ctx, merchantID, details)
	if err != nil {
		return nil, fmt.Errorf("failed to save bank details: %w", err)
	}

	// Convert to response DTO
	response := &dto.BankDetailsResponse{
		ID:            bankDetails.ID,
		MerchantID:    bankDetails.MerchantID,
		BankName:      bankDetails.BankName,
		BankCode:      bankDetails.BankCode,
		AccountNumber: bankDetails.AccountNumber,
		AccountName:   bankDetails.AccountName,
		RecipientCode: bankDetails.RecipientCode,
		Currency:      bankDetails.Currency,
		//Status:        bankDetails.Status,
		CreatedAt:     bankDetails.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     bankDetails.UpdatedAt.Format(time.RFC3339),
	}

	return response, nil
}

// GetBankDetails retrieves bank account details for a merchant
func (s *MerchantService) GetBankDetails(ctx context.Context, merchantID string) (*dto.BankDetailsResponse, error) {
	bankDetails, err := s.repo.GetBankDetails(ctx, merchantID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve bank details: %w", err)
	}

	// Convert to response DTO
	response := &dto.BankDetailsResponse{
		ID:            bankDetails.ID,
		MerchantID:    bankDetails.MerchantID,
		BankName:      bankDetails.BankName,
		BankCode:      bankDetails.BankCode,
		AccountNumber: bankDetails.AccountNumber,
		AccountName:   bankDetails.AccountName,
		RecipientCode: bankDetails.RecipientCode,
		Currency:      bankDetails.Currency,
		//Status:        bankDetails.Status,
		CreatedAt:     bankDetails.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     bankDetails.UpdatedAt.Format(time.RFC3339),
	}

	return response, nil
}

// UpdateBankDetails updates bank account details for a merchant
func (s *MerchantService) UpdateBankDetails(ctx context.Context, merchantID string, details dto.BankDetailsRequest) (*dto.BankDetailsResponse, error) {
	// Validate inputs
	if details.BankName == "" {
		return nil, errors.New("bank name is required")
	}
	if details.AccountNumber == "" {
		return nil, errors.New("account number is required")
	}
	if details.AccountName == "" {
		return nil, errors.New("account name is required")
	}

	// Validate bank name
	bankSvc := bank.GetBankService()
	if err := bankSvc.LoadBanks(); err != nil {
		return nil, fmt.Errorf("failed to load banks: %w", err)
	}

	bankCode, err := bankSvc.GetBankCode(details.BankName)
	if err != nil {
		return nil, fmt.Errorf("invalid bank name: %w", err)
	}

	// Override with validated code
	details.BankCode = bankCode

	// Update bank details via repository
	bankDetails, err := s.repo.UpdateBankDetailsRecord(ctx, merchantID, details)
	if err != nil {
		return nil, fmt.Errorf("failed to update bank details: %w", err)
	}

	// Convert to response DTO
	response := &dto.BankDetailsResponse{
		ID:            bankDetails.ID,
		MerchantID:    bankDetails.MerchantID,
		BankName:      bankDetails.BankName,
		BankCode:      bankDetails.BankCode,
		AccountNumber: bankDetails.AccountNumber,
		AccountName:   bankDetails.AccountName,
		RecipientCode: bankDetails.RecipientCode,
		Currency:      bankDetails.Currency,
		//Status:        bankDetails.Status,
		CreatedAt:     bankDetails.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     bankDetails.UpdatedAt.Format(time.RFC3339),
	}

	return response, nil
}

// DeleteBankDetails removes bank account details for a merchant
func (s *MerchantService) DeleteBankDetails(ctx context.Context, merchantID string) error {
	if err := s.repo.DeleteBankDetails(ctx, merchantID); err != nil {
		return fmt.Errorf("failed to delete bank details: %w", err)
	}
	return nil
}





















// UpdateMerchantProfile updates a merchant's profile information
func (s *MerchantService) UpdateMerchantProfile(ctx context.Context, merchantID string, input dto.UpdateMerchantProfileInput) error {
	// Validate input
	if err := s.validate.Struct(input); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Prepare updates
	updates := make(map[string]interface{})

	if input.StoreName != nil {
		updates["store_name"] = *input.StoreName
	}
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.PersonalEmail != nil {
		updates["personal_email"] = *input.PersonalEmail
	}
	if input.WorkEmail != nil {
		updates["work_email"] = *input.WorkEmail
	}
	if input.PhoneNumber != nil {
		updates["phone_number"] = *input.PhoneNumber
	}
	if input.PersonalAddress != nil {
		addressBytes, err := json.Marshal(input.PersonalAddress)
		if err != nil {
			return fmt.Errorf("failed to marshal personal address: %w", err)
		}
		updates["personal_address"] = addressBytes
	}
	if input.WorkAddress != nil {
		addressBytes, err := json.Marshal(input.WorkAddress)
		if err != nil {
			return fmt.Errorf("failed to marshal work address: %w", err)
		}
		updates["work_address"] = addressBytes
	}
	if input.BusinessType != nil {
		updates["business_type"] = *input.BusinessType
	}
	if input.Website != nil {
		updates["website"] = *input.Website
	}
	if input.BusinessDescription != nil {
		updates["business_description"] = *input.BusinessDescription
	}
	if input.StoreLogoURL != nil {
		updates["store_logo_url"] = *input.StoreLogoURL
	}
	if input.Banner != nil {
		updates["banner"] = *input.Banner
	}

	// Update merchant
	if err := s.repo.UpdateMerchant(ctx, merchantID, updates); err != nil {
		return fmt.Errorf("failed to update merchant: %w", err)
	}

	return nil

}



