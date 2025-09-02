package services

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"api-customer-merchant/internal/shared/auth/models"
	"api-customer-merchant/internal/shared/auth/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthService struct {
	repo      *repositories.UserRepository
	oauthConf *oauth2.Config
}

func NewAuthService() *AuthService {
	return &AuthService{
		repo: repositories.NewUserRepository(),
		oauthConf: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  "", // Set dynamically in GoogleLogin
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

func (s *AuthService) GetOAuthConfig() *oauth2.Config {
	return s.oauthConf
}

func (s *AuthService) Register(email, name, password, country string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:    email,
		Name:     name,
		Password: string(hashedPassword),
		Country:  country,
		Role:     "customer", // Default role
	}

	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Login(email, password string) (*models.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}

func (s *AuthService) GenerateJWT(user *models.User) (string, error) {
	secret := []byte("your_super_secret_key_here") // Load from env
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func (s *AuthService) GoogleLogin(code, baseURL string) (*models.User, string, error) {
	s.oauthConf.RedirectURL = baseURL + "/customer/auth/google/callback"
	token, err := s.oauthConf.Exchange(context.Background(), code)
	if err != nil {
		return nil, "", errors.New("failed to exchange code")
	}

	client := s.oauthConf.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, "", errors.New("failed to get user info")
	}
	defer resp.Body.Close()

	var userInfo struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, "", errors.New("failed to decode user info")
	}

	user, err := s.repo.FindByGoogleID(userInfo.ID)
	if err != nil {
		// User doesn't exist, create one
		user = &models.User{
			Email:    userInfo.Email,
			Name:     userInfo.Name,
			Role:     "customer", // Default role
			GoogleID: userInfo.ID,
			// Country left empty (optional)
		}
		if err := s.repo.Create(user); err != nil {
			return nil, "", err
		}
	}

	jwtToken, err := s.GenerateJWT(user)
	if err != nil {
		return nil, "", err
	}

	return user, jwtToken, nil
}