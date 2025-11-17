package handlers

import (
	//"fmt"
	"fmt"
	"log"
	"net/http"
	"net/url"

	//"net/url"
	"os"
	"strconv"

	//"os"
	"strings"

	//"api-customer-merchant/internal/db/models"
	//"api-customer-merchant/internal/db/repositories"
	"api-customer-merchant/internal/api/dto"
	//"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/services/email"
	services "api-customer-merchant/internal/services/user"
	"api-customer-merchant/internal/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	service      *services.AuthService
	emailService *email.EmailService
}

// In customer/handlers/auth_handler.go AND merchant/handlers/auth_handler.go
func NewAuthHandler(s *services.AuthService , emailSvc *email.EmailService) *AuthHandler {
	return &AuthHandler{
		service: s,
		emailService: emailSvc,
	}
}

// Register godoc
// @Summary Register a new customer
// @Description Creates a new customer account with email, name, password, and optional country
// @Tags Customer
// @Accept json
// @Produce json
// @Param body body dto.RegisterRequest true "Customer registration details"
// @Success 200 {object} object{token=string} "JWT token"
// @Failure 400 {object} object{error=string} "Invalid request"
// @Failure 500 {object} object{error=string} "Server error"
// @Router /customer/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	// var req struct {
	// 	Email    string `json:"email" binding:"required,email"`
	// 	Name     string `json:"name" binding:"required"`
	// 	Password string `json:"password" binding:"required,min=6"`
	// 	Country  string `json:"country"`
	// }
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.RegisterUser(req.Email, req.Name, req.Password, req.Country)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go func() {
		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:3000"
		}
		
		emailData := map[string]interface{}{
			"Name":           user.Name,
			"MarketplaceURL": frontendURL,
		}
		
		if err := h.emailService.SendWelcome(user.Email, emailData); err != nil {
			log.Printf("Failed to send welcome email to %s: %v", user.Email, err)
		}
	}()

	token, err := h.service.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token})
}

// Login godoc
// @Summary Customer login
// @Description Authenticates a customer using email and password
// @Tags Customer
// @Accept json
// @Produce json
// @Param body body object{email=string,password=string} true "Customer login credentials"
// @Success 200 {object} object{token=string} "JWT token"
// @Failure 400 {object} object{error=string} "Invalid request"
// @Failure 401 {object} object{error=string} "Unauthorized"
// @Failure 403 {object} object{error=string} "Invalid role"
// @Failure 500 {object} object{error=string} "Server error"
// @Router /customer/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	// var req struct {
	// 	Email    string `json:"email" binding:"required,email"`
	// 	Password string `json:"password" binding:"required"`
	// }
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.LoginUser(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// if user.Role != "customer" {
	// 	c.JSON(http.StatusForbidden, gin.H{"error": "Invalid role for this API"})
	// 	return
	// }

	token, err := h.service.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// GoogleAuth godoc
// @Summary Initiate Google OAuth for customer
// @Description Redirects to Google OAuth login page
// @Tags Customer
// @Produce json
// @Success 307 {object} object{} "Redirect to Google OAuth"
// @Router /customer/auth/google [get]
func (h *AuthHandler) GoogleAuth(c *gin.Context) {
	url := h.service.GetOAuthConfig("customer").AuthCodeURL("state-customer", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
	//c.JSON(http.StatusOK, gin.H{"url": url})
}

// GoogleCallback godoc
// @Summary Handle Google OAuth callback for customer
// @Description Processes Google OAuth callback and returns JWT token
// @Tags Customer
// @Produce json
// @Param code query string true "OAuth code"
// @Success 200 {object} object{token=string} "JWT token"
// @Failure 400 {object} object{error=string} "Code not provided"
// @Failure 500 {object} object{error=string} "Server error"
// @Router /customer/auth/google/callback [get]
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	if code == "" || state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code or state"})
		return
	}
	// Verify state (in production, check against stored value)
	if state != "state-customer" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state parameter"})
		return
	}

	user, token, err := h.service.GoogleLogin(code, os.Getenv("BASE_URL"), "customer")
	if err != nil {
		log.Printf("Google login failed: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}




	go func() {
		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:3000"
		}
		
		emailData := map[string]interface{}{
			"Name":           user.Name,
			"MarketplaceURL": frontendURL,
		}
		
		// Send welcome email (the service should handle duplicate prevention if needed)
		if err := h.emailService.SendWelcome(user.Email, emailData); err != nil {
			log.Printf("Failed to send welcome email to %s: %v", user.Email, err)
		}
	}()

	// // --- Determine redirect URL dynamically ---
	// var frontendURL string

	// //  Try to get from Origin header
	// origin := c.Request.Header.Get("Origin")

	// //  If Origin is missing, try Referer
	// if origin == "" {
	// 	referer := c.Request.Referer()
	// 	if referer != "" {
	// 		// Extract base (scheme + host) from referer
	// 		u, err := url.Parse(referer)
	// 		if err == nil {
	// 			origin = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	// 		}
	// 	}
	// }

	// // Fallback to environment or localhost
	// if origin != "" {
	// 	frontendURL = origin
	// } else if os.Getenv("FRONTEND_URL") != "" {
	// 	frontendURL = os.Getenv("FRONTEND_URL")
	// } else {
	// 	frontendURL = "http://localhost:3000"
	// }
	// Decode the frontend URL
	//  frontendURL, err := url.QueryUnescape(state)
	//  if err != nil {
	// 	 frontendURL = os.Getenv("FRONTEND_URL")
	//  }

	//c.JSON(http.StatusCreated, gin.H{"token": token, "user": user})

	// c.SetCookie(
	// 	"auth_token",        // name
	// 	token,               // value
	// 	3600*24,             // max age (1 day)
	// 	"/",                 // path
	// 	//".yourdomain.com",
	// 	 "localhost",                   // domain — change this (see notes below)
	// 	true,                // secure (HTTPS only)
	// 	true,                // httpOnly (not accessible to JS)
	// )
	// frontendURL := "http://localhost:3000"
	//redirectURL := fmt.Sprintf("%s/auth/success?token=%s", frontendURL, token)

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	isLocal := strings.Contains(frontendURL, "localhost")

	// --- Domain + Security ---
	domain := ""
	secure := false
	if !isLocal {
		u, err := url.Parse(frontendURL)
		if err == nil && u.Host != "" {
			parts := strings.Split(u.Host, ":")
			host := parts[0]
			if !strings.HasPrefix(host, ".") {
				host = "." + host
			}
			domain = host
		}
		secure = true
	}

	// --- Choose SameSite Policy ---
	// Use Lax in production (most secure)
	// Use None for localhost (since it's often cross-origin)
	sameSite := "Lax"
	if isLocal {
		sameSite = "None"
	}

	// --- Build Cookie ---
	cookie := fmt.Sprintf(
		"auth_token=%s; Path=/; Max-Age=%d; HttpOnly; Secure=%t; SameSite=%s",
		token, 3600*24, secure, sameSite,
	)
	if domain != "" {
		cookie += fmt.Sprintf("; Domain=%s", domain)
	}

	c.Writer.Header().Add("Set-Cookie", cookie)

	redirectURL := fmt.Sprintf("%s/auth/success", frontendURL)
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// Logout godoc
// @Summary Customer logout
// @Description Invalidates the customer's JWT token
// @Tags Customer
// @Security BearerAuth
// @Produce json
// @Success 200 {object} object{message=string} "Logout successful"
// @Failure 400 {object} object{error=string} "Authorization header required"
// @Router /customer/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header required"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	err := utils.Add(tokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to logoout"})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// Update godoc
// @Summary Customer Profile Update
// @Description Updates Customer profile
// @Tags Customer
// @Security BearerAuth
// @Param body body dto.UserUpdateRequest true "Update Customer Profile"
// @Produce json
// @Success 200 {object} object{message=string} "Update successful"
// @Failure 400 {object} object{error=string} "Authorization header required"
// @Router /customer/update [patch]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	ctx := c.Request.Context()
	userID, exists := c.Get("userID")

	if !exists {

		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user session"})
		return
	}
	userIDUint, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	// var req struct {
	// 	Name      string
	// 	Country   string
	// 	Addresses []string
	// }
	var req dto.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.UpdateProfile(ctx, uint(userIDUint), req.Name, req.Country, req.Addresses); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// GetProfile godoc
// @Summary Get customer profile
// @Description Retrieves the customer's profile information
// @Tags Customer
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.ProfileResponse "Profile details"
// @Failure 401 {object} object{error=string} "Unauthorized"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /customer/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	ctx := c.Request.Context()
	userID, exists := c.Get("userID")
	if !exists {

		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user session"})
		return
	}
	userIDUint, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := h.service.GetUser(ctx, uint(userIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//prof dto.Pr

	// addressList := make([]string, len(user.Addresses))
	// for i, addr := range user.Addresses {
	// 	addressList[i] = addr.Address // Access the Address field from the UserAddress struct
	// }

	resp := &dto.ProfileResponse{
		ID:      user.ID,
		Email:   user.Email,
		Name:    user.Name,
		Country: user.Country,
		//Addresses: addressList, // Assign the converted slice
	}
	// if err := utils.RespMap(user, resp); err != nil {

	// }
	c.JSON(http.StatusOK, resp)

	//c.JSON(http.StatusOK, user)
}
// RequestPasswordReset godoc
// @Summary Request password reset
// @Description Sends a password reset email with a secure token
// @Tags Customer
// @Accept json
// @Produce json
// @Param body body dto.RequestPasswordResetRequest true "Email address"
// @Success 200 {object} object{message=string} "Password reset email sent"
// @Failure 400 {object} object{error=string} "Invalid input"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /customer/request-password-reset [post]
func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	var req dto.RequestPasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate reset token
	token, expiresAt, err := h.service.GeneratePasswordResetToken(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send reset email
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, token)
	
	emailData := map[string]interface{}{
		"Name":      req.Email, // Use email if name not available
		"ResetLink": resetLink,
		"ExpiresAt": expiresAt.Format("January 2, 2006 at 3:04 PM"),
	}

	if err := h.emailService.SendPasswordReset(req.Email, emailData); err != nil {
		log.Printf("Failed to send password reset email to %s: %v", req.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send reset email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset email sent"})
}

// ResetPassword godoc
// @Summary Reset customer password
// @Description Resets the customer's password using a valid reset token
// @Tags Customer
// @Accept json
// @Produce json
// @Param body body dto.ResetPasswordRequest true "Reset details with token"
// @Success 200 {object} object{message=string} "Password reset successful"
// @Failure 400 {object} object{error=string} "Invalid input or token"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /customer/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate token and reset password
	err := h.service.ResetPasswordWithToken(req.Token, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
}