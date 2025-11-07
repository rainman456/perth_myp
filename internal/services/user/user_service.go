package user

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	//"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"

	//"google.golang.org/api/oauth2/v2"

	//"api-customer-merchant/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

type googleUserInfo struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// func (s *AuthService) RegisterUser(email, name, password, country string,address []string) (*models.User, error) {
// 	_, err := s.userRepo.FindByEmail(email)
// 	if err == nil {
// 		return nil, errors.New("email already exists")
// 	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
// 		return nil, err
// 	}

// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return nil, err
// 	}

// 	user := &models.User{
// 		Email:    email,
// 		Name:     name,
// 		Password: string(hashedPassword),
// 		Country:  country,
// 		Address:  address,
// 	}

// 	return user, nil
// }


// RegisterUser registers a new user
func (s *AuthService) RegisterUser(email, name, password, country string) (*models.User, error) {
	_, err := s.userRepo.FindByEmail(email)
	if err == nil {
		return nil, errors.New("email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Convert []string to []Addresslist
	// addressList := make([]models.UserAddress, len(addresses))
	// for i, addr := range addresses {
	// 	addressList[i] = models.UserAddress{
	// 		Address: addr,
	// 	}
	// }

	user := &models.User{
		Email:    email,
		Name:     name,
		Password: string(hashedPassword),
		Country:  country,
		// Change Address to the correct field name
		//Addresses: addressList, // Use the converted slice
	}



	err = s.userRepo.Create(user)
if err != nil {
    return nil, err
}
return user, nil
}

func (s *AuthService) LoginUser(email, password string) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *AuthService) GenerateJWT(entity interface{}) (string, error) {
	var id uint
	var entityType string
	var name string
	var email string

	switch e := entity.(type) {
	case *models.User:
		id = e.ID
		entityType = "customer"
		name = e.Name
		email = e.Email

	default:
		return "", errors.New("invalid entity type")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         float64(id),
		"entityType": entityType,
		"name":       name,
		"email":      email,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET not set")
	}

	return token.SignedString([]byte(secret))
}

func (s *AuthService) GetOAuthConfig(entityType string) *oauth2.Config {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		log.Println("BASE_URL environment variable is not set")
		return nil
	}
	redirectURL := baseURL + "/customer/auth/google/callback"
	log.Printf("OAuth redirect URL: %s", redirectURL)
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid"},
		Endpoint: google.Endpoint,
	}
}

func (s *AuthService) GoogleLogin(code, baseURL, entityType string) (*models.User, string, error) {
	if entityType != "customer" {
		log.Printf("Invalid entityType for OAuth: %s", entityType)
		return nil, "", errors.New("OAuth only supported for customers")
	}

	// Get OAuth config
	oauthConfig := s.GetOAuthConfig(entityType)
	if oauthConfig == nil || oauthConfig.ClientID == "" || oauthConfig.ClientSecret == "" {
		log.Println("Google OAuth credentials not set")
		return nil, "", errors.New("OAuth configuration error")
	}

	// Exchange code for access token
	ctx := context.Background()
	token, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		log.Printf("Failed to exchange code: %v", err)
		return nil, "", errors.New("failed to exchange code")
	}

	// Fetch user info
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		log.Printf("Failed to create userinfo request: %v", err)
		return nil, "", errors.New("failed to create userinfo request")
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to get user info: %v", err)
		return nil, "", errors.New("failed to get user info")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Userinfo endpoint returned status: %d", resp.StatusCode)
		return nil, "", errors.New("failed to get user info")
	}

	var userInfo googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		log.Printf("Failed to decode user info: %v", err)
		return nil, "", errors.New("failed to decode user info")
	}

	// Validate email
	if userInfo.Email == "" {
		log.Println("No email provided by Google")
		return nil, "", errors.New("no email provided")
	}

	// Check if user exists
	user, err := s.userRepo.FindByEmail(userInfo.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Failed to find user: %v", err)
		return nil, "", err
	}
	if user == nil {
		// Register new user
		user = &models.User{
			Email:   userInfo.Email,
			Name:    userInfo.Name,
			Country: "", // Set default or prompt later
		}
		if err := s.userRepo.Create(user); err != nil {
			log.Printf("Failed to create user: %v", err)
			return nil, "", err
		}
	}

	// Generate JWT
	jwtToken, err := s.GenerateJWT(user)
	if err != nil {
		log.Printf("Failed to generate JWT: %v", err)
		return nil, "", errors.New("failed to generate JWT")
	}

	return user, jwtToken, nil
}

func (s *AuthService) UpdateProfile(ctx context.Context ,userID uint, name, country string, addresses []string) error {

	user, err := s.userRepo.FindByID(ctx,userID)
	if err != nil {
		return err
	}


	// addressList := make([]models.UserAddress, len(addresses))
	// for i, addr := range addresses {
	// 	addressList[i] = models.UserAddress{
	// 		Address: addr,
	// 	}
	// }
	user.Name = name
	user.Country = country
	//user.Addresses = addressList
	// Addresses as JSON; add field to User model if needed: Addresses jsonb
	return s.userRepo.Update(user)
}

func (s *AuthService) ResetPassword(email, newPassword string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)
	return s.userRepo.Update(user)
}



// GetOrder retrieves a single order by its ID.
func (s *AuthService) GetUser(ctx context.Context, id uint) (*models.User, error) {
	if id == 0 {
		return nil, errors.New("invalid order ID")
	}
	// The repository already preloads necessary associations.
	return s.userRepo.FindByID(ctx, id)
	//return s.orderRepo.FindByID(id)
}


