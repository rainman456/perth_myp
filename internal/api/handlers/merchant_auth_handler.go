package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/services/merchant"
	"api-customer-merchant/internal/utils"

	"github.com/gin-gonic/gin"
)


type MerchantHandler struct {
	service *merchant.MerchantService
}

func NewMerchantAuthHandler(s *merchant.MerchantService) *MerchantHandler {
	return &MerchantHandler{service: s}
}

// Apply godoc
// @Summary Submit a new merchant application
// @Description Allows a prospective merchant to submit an application with personal, business, and address information
// @Tags Merchant
// @Accept json
// @Produce json
// @Param body body object{first_name=string,last_name=string,email=string,phone=string,personal_address=object{street=string,city=string,state=string,postal_code=string,country=string},work_address=object{street=string,city=string,state=string,postal_code=string,country=string},business_name=string,business_type=string,tax_id=string,documents=object{business_license=string,identification=string}} true "Merchant application details"
// @Success 201 {object} object{id=string,first_name=string,last_name=string,email=string,phone=string,personal_address=object{street=string,city=string,state=string,postal_code=string,country=string},work_address=object{street=string,city=string,state=string,postal_code=string,country=string},business_name=string,business_type=string,tax_id=string,documents=object{business_license=string,identification=string},status=string,created_at=string} "Created application"
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
            StoreName: req.StoreName,
            Name: req.Name,
            PersonalEmail: req.PersonalEmail,
            WorkEmail: req.WorkEmail,
            PhoneNumber: req.PhoneNumber,
        },
        MerchantAddress: models.MerchantAddress{
            PersonalAddress: personalAddressJSON,
            WorkAddress: workAddressJSON,
        },
        MerchantBusinessInfo: models.MerchantBusinessInfo{
            BusinessType: req.BusinessType,
            Website: req.Website,
            BusinessDescription: req.BusinessDescription,
            BusinessRegistrationNumber: req.BusinessRegistrationNumber,
        },
        MerchantDocuments: models.MerchantDocuments{
            StoreLogoURL: req.StoreLogoURL,
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
        ID: app.ID,
        StoreName: app.StoreName,
        Name: app.Name,
        PersonalEmail: app.PersonalEmail,
        WorkEmail: app.WorkEmail,
        PhoneNumber: app.PhoneNumber,
        PersonalAddress: req.PersonalAddress,
        WorkAddress: req.WorkAddress,
        BusinessType: app.BusinessType,
        Website: app.Website,
        BusinessDescription: app.BusinessDescription,
        BusinessRegistrationNumber: app.BusinessRegistrationNumber,
        StoreLogoURL: app.StoreLogoURL,
        BusinessRegistrationCertificate: app.BusinessRegistrationCertificate,
        Status: app.Status,
        CreatedAt: app.CreatedAt.Format(time.RFC3339),
        UpdatedAt: app.UpdatedAt.Format(time.RFC3339),
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

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// GetApplication godoc
// @Summary Retrieve a merchant application by ID
// @Description Fetches the details and status of a submitted merchant application (e.g., 'pending', 'approved', 'rejected')
// @Tags Merchant
// @Produce json
// @Param id path string true "Application ID (UUID)"
// @Success 200 {object} object{id=string,status=string,merchant_basic_info=object{name=string,store_name=string,personal_email=string,work_email=string,password=string},personal_address=object{street=string,city=string,state=string,zip=string,country=string},work_address=object{street=string,city=string,state=string,zip=string,country=string},merchant_business_info=object{business_type=string,years_in_business=integer,annual_revenue=number,tax_id=string},merchant_documents=object{id_proof=string,business_license=string,bank_statement=string},created_at=string,updated_at=string} "Application details retrieved"
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
// @Success 200 {object} object{id=string,user_id=string,business_name=string,business_type=string,tax_id=string,personal_address=object{street=string,city=string,state=string,zip=string,country=string},work_address=object{street=string,city=string,state=string,zip=string,country=string},status=string,created_at=string,updated_at=string} "Merchant account details"
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
	err:=utils.Add(tokenString)
	if  err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":"failed to logoout"})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
