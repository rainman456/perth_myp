package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/services/email"
	"api-customer-merchant/internal/services/merchant"
	"api-customer-merchant/internal/utils"

	"github.com/gin-gonic/gin"
)

type MerchantHandler struct {
	service      *merchant.MerchantService
	emailService *email.EmailService
}

func NewMerchantAuthHandler(s *merchant.MerchantService,emailSvc *email.EmailService) *MerchantHandler {
	return &MerchantHandler{	
		service:      s,
		emailService: emailSvc,}
}

// Apply godoc
// @Summary Submit a new merchant application
// @Description Allows a prospective merchant to submit an application with personal, business, and address information
// @Tags Merchant
// @Accept json
// @Produce json
// @Param body body dto.MerchantApplyDTO true "Merchant application details"
// @Success 201 {object} dto.MerchantApplyResponse "Created application"
// @Failure 400 {object} object{error=string} "Invalid request body or malformed JSON"
// @Failure 500 {object} object{error=string} "Failed to submit application"
// @Router /merchant/apply [post]
func (h *MerchantHandler) Apply(c *gin.Context) {
	var req dto.MerchantApplyDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}
	// Convert DTO to model
	personalAddressJSON, _ := json.Marshal(req.PersonalAddress)
	workAddressJSON, _ := json.Marshal(req.WorkAddress)
	app := &models.MerchantApplication{
		MerchantBasicInfo: models.MerchantBasicInfo{
			StoreName:     req.StoreName,
			Name:          req.Name,
			PersonalEmail: req.PersonalEmail,
			WorkEmail:     req.WorkEmail,
			PhoneNumber:   req.PhoneNumber,
		},
		MerchantAddress: models.MerchantAddress{
			PersonalAddress: personalAddressJSON,
			WorkAddress:     workAddressJSON,
		},
		MerchantBusinessInfo: models.MerchantBusinessInfo{
			BusinessType:               req.BusinessType,
			Website:                    req.Website,
			BusinessDescription:        req.BusinessDescription,
			BusinessRegistrationNumber: req.BusinessRegistrationNumber,
		},
		MerchantDocuments: models.MerchantDocuments{
			StoreLogoURL:                    req.StoreLogoURL,
			BusinessRegistrationCertificate: req.BusinessRegistrationCertificate,
		},
	}
	// Service call, response mapping
	app, err := h.service.SubmitApplication(c.Request.Context(), app)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to submit application: " + err.Error()})
		return
	}
	// Build response
	//resp := dto.MerchantApplyResponse
	resp := dto.MerchantApplyResponse{
		ID:                              app.ID,
		StoreName:                       app.StoreName,
		Name:                            app.Name,
		PersonalEmail:                   app.PersonalEmail,
		WorkEmail:                       app.WorkEmail,
		PhoneNumber:                     app.PhoneNumber,
		PersonalAddress:                 req.PersonalAddress,
		WorkAddress:                     req.WorkAddress,
		BusinessType:                    app.BusinessType,
		Website:                         app.Website,
		BusinessDescription:             app.BusinessDescription,
		BusinessRegistrationNumber:      app.BusinessRegistrationNumber,
		StoreLogoURL:                    app.StoreLogoURL,
		BusinessRegistrationCertificate: app.BusinessRegistrationCertificate,
		Status:                          app.Status,
		CreatedAt:                       app.CreatedAt.Format(time.RFC3339),
		UpdatedAt:                       app.UpdatedAt.Format(time.RFC3339),
	}
	c.JSON(http.StatusCreated, resp)
}

// Login godoc
// @Summary Merchant login
// @Description Authenticates a merchant using work email and password
// @Tags Merchant
// @Accept json
// @Produce json
// @Param body body dto.MerchantLogin true "Login credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Router /merchant/login [post]
func (h *MerchantHandler) Login(c *gin.Context) {
	// var req struct {
	// 	Work_Email string `json:"email" binding:"required,email"`
	// 	Password   string `json:"password" binding:"required"`
	// }
	var req dto.MerchantLogin

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	merchant, err := h.service.LoginMerchant(c.Request.Context(), req.Work_Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.GenerateJWT(merchant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	merchantResponse := gin.H{
		"id":       merchant.ID,        // Assuming merchant.ID is the field; adjust if different (e.g., merchant.Id)
		"email":    merchant.WorkEmail, // Map work_email to "email" in response
		"username": merchant.Name,      // Assuming merchant.Username field
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "merchant": merchantResponse})
}

// GetApplication godoc
// @Summary Retrieve a merchant application by ID
// @Description Fetches the details and status of a submitted merchant application (e.g., 'pending', 'approved', 'rejected')
// @Tags Merchant
// @Produce json
// @Param id path string true "Application ID (UUID)"
// @Success 200 {object} dto.MerchantApplyResponse "Application details retrieved"
// @Failure 404 {object} object{error=string} "Application not found"
// @Failure 500 {object} object{error=string} "Failed to retrieve application"
// @Router /merchant/application/{id} [get]
func (h *MerchantHandler) GetApplication(c *gin.Context) {
	id := c.Param("id")
	app, err := h.service.GetApplication(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, app)
}

// GetMyMerchant godoc
// @Summary Retrieve current merchant account
// @Description Fetches the merchant account details for the authenticated user, if their application has been approved
// @Tags Merchant
// @Produce json
// @Security BearerAuth
// @Success 200 {object}   dto.MerchantApplyResponse  "Merchant account details"
// @Failure 401 {object} object{error=string} "Unauthorized: Missing or invalid authentication"
// @Failure 404 {object} object{error=string} "Merchant account not found"
// @Router /merchant/me [get]
func (h *MerchantHandler) GetMyMerchant(c *gin.Context) {
	userID, ok := c.Get("merchantID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	m, err := h.service.GetMerchantByUserID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "merchant not found: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, m)
}

// UpdateProfile godoc
// @Summary Update merchant profile
// @Description Updates the merchant's profile information
// @Tags Merchant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.UpdateMerchantProfileInput true "Profile update details"
// @Success 200 {object} object{message=string} "Profile updated successfully"
// @Failure 400 {object} object{error=string} "Invalid request body"
// @Failure 401 {object} object{error=string} "Unauthorized"
// @Failure 500 {object} object{error=string} "Failed to update profile"
// @Router /merchant/profile [put]
func (h *MerchantHandler) UpdateProfile(c *gin.Context) {
	merchantID, exists := c.Get("merchantID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.UpdateMerchantProfileInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateMerchantProfile(c.Request.Context(), merchantID.(string), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile updated successfully"})
}

// Logout godoc
// @Summary Merchant logout
// @Description Invalidates the Merchant's JWT token
// @Tags Merchant
// @Security BearerAuth
// @Produce json
// @Success 200 {object} object{message=string} "Logout successful"
// @Failure 400 {object} object{error=string} "Authorization header required"
// @Router /merchant/logout [post]
func (h *MerchantHandler) Logout(c *gin.Context) {
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




















// RequestPasswordReset godoc
// @Summary Request password reset for merchant
// @Description Sends a password reset email with a secure token
// @Tags Merchant
// @Accept json
// @Produce json
// @Param body body dto.RequestPasswordResetRequest true "Work email address"
// @Success 200 {object} object{message=string} "Password reset email sent"
// @Failure 400 {object} object{error=string} "Invalid input"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /merchant/request-password-reset [post]
func (h *MerchantHandler) RequestPasswordReset(c *gin.Context) {
	var req dto.MerchRequestPasswordResetRequest
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

	resetLink := fmt.Sprintf("%s/merchant/reset-password?token=%s", frontendURL, token)
	
	emailData := map[string]interface{}{
		"Name":      req.Email,
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
// @Summary Reset merchant password
// @Description Resets the merchant's password using a valid reset token
// @Tags Merchant
// @Accept json
// @Produce json
// @Param body body dto.ResetPasswordRequest true "Reset details with token"
// @Success 200 {object} object{message=string} "Password reset successful"
// @Failure 400 {object} object{error=string} "Invalid input or token"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /merchant/reset-password [post]
func (h *MerchantHandler) ResetPassword(c *gin.Context) {
	var req dto.MerchResetPasswordRequest
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