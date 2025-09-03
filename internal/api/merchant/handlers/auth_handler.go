package handlers

import (
	"net/http"
	"strings"

	 //"api-customer-merchant/internal/db/models"
	 services "api-customer-merchant/internal/domain/identity"
	"api-customer-merchant/internal/utils"

	"github.com/gin-gonic/gin"
	//"golang.org/x/oauth2"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{service: services.NewAuthService()}
}


// Register godoc
// @Summary Submit application for new merchant
// @Description Creates a new merchant account with required details (password set by admin console later)
// @Tags Merchant
// @Accept json
// @Produce json
// @Param body body object{email=string,name=string,password=string,country=string,store_name=string,personal_email=string,work_email=string,phone_number=string,street_address=string,city=string,zip_code=string,work_address=string,business_type=string,website=string,business_description=string,store_logo_url=string,business_registration_certificate=string} true "Merchant registration details"
// @Success 200 {object} object{message=string} "JWT token"
// @Failure 400 {object} object{error=string} "Invalid request"
// @Failure 500 {object} object{error=string} "Server error"
// @Router /merchant/submitApplication [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Email                        string `json:"email" binding:"required,email"`
		Name                         string `json:"name" binding:"required"`
		Password                     string `json:"password" binding:"required,min=6"`
		Country                      string `json:"country"`
		StoreName                    string `json:"store_name" binding:"required"`
		PersonalEmail                string `json:"personal_email" binding:"required,email"`
		WorkEmail                    string `json:"work_email" binding:"required,email"`
		PhoneNumber                  string `json:"phone_number"`
		StreetAddress                string `json:"street_address"`
		City                         string `json:"city"`
		ZipCode                      string `json:"zip_code"`
		WorkAddress                  string `json:"work_address"`
		BusinessType                 string `json:"business_type"`
		Website                      string `json:"website"`
		BusinessDescription          string `json:"business_description"`
		StoreLogoURL                 string `json:"store_logo_url"`
		BusinessRegistrationCertificate string `json:"business_registration_certificate"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create user with merchant role
	_, err := h.service.RegisterMerchant(
		
		req.Name,
		"", // Password set by admin console
		req.Country,
		req.StoreName,
		req.PersonalEmail,
		req.WorkEmail,
		req.PhoneNumber,
		req.StreetAddress,
		req.City,
		req.ZipCode,
		req.WorkAddress,
		req.BusinessType,
		req.Website,
		req.BusinessDescription,
		req.StoreLogoURL,
		req.BusinessRegistrationCertificate,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}




	c.JSON(http.StatusOK, gin.H{"message": "Merchant application submitted"})
}


// Login godoc
// @Summary Merchant login
// @Description Authenticates a merchant using work_email and password
// @Tags Merchant
// @Accept json
// @Produce json
// @Param body body object{work_email=string,password=string} true "Merchant login credentials"
// @Success 200 {object} object{token=string} "JWT token"
// @Failure 400 {object} object{error=string} "Invalid request"
// @Failure 401 {object} object{error=string} "Unauthorized"
// @Failure 403 {object} object{error=string} "Account not approved"
// @Failure 500 {object} object{error=string} "Server error"
// @Router /merchant/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		WorkEmail string `json:"work_email" binding:"required,email"`
		Password  string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	merchant, err := h.service.LoginMerchant(req.WorkEmail, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.GenerateJWT(merchant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}


// func (h *AuthHandler) GoogleAuth(c *gin.Context) {
// 	url := h.service.GetOAuthConfig().AuthCodeURL("state-customer", oauth2.AccessTypeOffline)
// 	c.Redirect(http.StatusTemporaryRedirect, url)
// }

// func (h *AuthHandler) GoogleCallback(c *gin.Context) {
// 	code := c.Query("code")
// 	if code == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Code not provided"})
// 		return
// 	}

// 	_, token, err := h.service.GoogleLogin(code, "http://localhost:8080")
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"token": token})
// }

// Logout godoc
// @Summary Merchant logout
// @Description Invalidates the merchant's JWT token
// @Tags Merchant
// @Security BearerAuth
// @Produce json
// @Success 200 {object} object{message=string} "Logout successful"
// @Failure 400 {object} object{error=string} "Authorization header required"
// @Router /merchant/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header required"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	utils.Add(tokenString)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}