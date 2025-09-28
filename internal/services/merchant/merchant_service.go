package merchant

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"api-customer-merchant/internal/api/dto"
//"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"

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
	return s.repo.GetByID(ctx, id)
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
		id = e.ID
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



func (s *MerchantService) UpdateBankDetails(ctx context.Context, merchantID string ,details  dto.BankDetailsRequest) error {
     // Similar, but use Save or Update
	 if details.BankName == "" {
		return errors.New("empty bank name")
	}

	if details.AccountNumber == "" {
		return errors.New("empty bank name")
	}
	

	

	err := s.repo.UpdateBankDetails(ctx ,merchantID, details)
	if err != nil {
		return  err
	}

	//payment.Status = models.PaymentStatus(status)
	

	return nil

  
 }