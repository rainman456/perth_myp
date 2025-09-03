package identity

import (
	"errors"
	"os"

	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	//"api-customer-merchant/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) RegisterUser(email, name, password, country string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:    email,
		Name:     name,
		Password: string(hashedPassword),
		Country:  country,
	}

	if err := db.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) RegisterMerchant(email, name, password, country, storeName, personalEmail, workEmail, phoneNumber, streetAddress, city, zipCode, workAddress, businessType, website, businessDescription, storeLogoURL, businessRegistrationCertificate string) (*models.Merchant, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	merchant := &models.Merchant{
		MerchantBasicInfo: models.MerchantBasicInfo{
			Name:          name,
			StoreName:     storeName,
			PersonalEmail: personalEmail,
			WorkEmail:     workEmail,
			PhoneNumber:   phoneNumber,
			Password:      string(hashedPassword),
		},
		MerchantAddress: models.MerchantAddress{
			StreetAddress: streetAddress,
			City:          city,
			Country:       country,
			ZipCode:       zipCode,
			WorkAddress:   workAddress,
		},
		MerchantBusinessInfo: models.MerchantBusinessInfo{
			BusinessType:        businessType,
			Website:             website,
			BusinessDescription: businessDescription,
		},
		MerchantDocuments: models.MerchantDocuments{
			StoreLogoURL:                   storeLogoURL,
			BusinessRegistrationCertificate: businessRegistrationCertificate,
		},
		Status: models.MerchantStatusPending,
	}

	if err := db.DB.Create(merchant).Error; err != nil {
		return nil, err
	}

	return merchant, nil
}

func (s *AuthService) LoginUser(email, password string) (*models.User, error) {
	var user models.User
	if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}

func (s *AuthService) LoginMerchant(workEmail, password string) (*models.Merchant, error) {
	var merchant models.Merchant
	if err := db.DB.Where("work_email = ?", workEmail).First(&merchant).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(merchant.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	if merchant.Status != models.MerchantStatusApproved {
		return nil, errors.New("merchant account not approved")
	}

	return &merchant, nil
}

func (s *AuthService) GenerateJWT(entity interface{}) (string, error) {
	var id uint
	var entityType string

	switch e := entity.(type) {
	case *models.User:
		id = e.ID
		entityType = "user"
	case *models.Merchant:
		id = e.ID
		entityType = "merchant"
	default:
		return "", errors.New("invalid entity type")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         id,
		"entityType": entityType,
	})

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET not set")
	}

	return token.SignedString([]byte(secret))
}

func (s *AuthService) GetOAuthConfig(entityType string) *oauth2.Config {
	redirectURL := os.Getenv("BASE_URL") + "/" + entityType + "/auth/google/callback"
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  redirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

func (s *AuthService) GoogleLogin(code, baseURL, entityType string) (interface{}, string, error) {
	// Placeholder for Google OAuth logic
	return nil, "", errors.New("not implemented")
}