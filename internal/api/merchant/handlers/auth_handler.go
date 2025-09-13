package handlers

import (
	"encoding/json"
	"net/http"

	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/domain/merchant"

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
	var req struct {
		models.MerchantBasicInfo
		PersonalAddress            map[string]any `json:"personal_address" validate:"required"`
		WorkAddress                map[string]any `json:"work_address" validate:"required"`
		models.MerchantBusinessInfo
		models.MerchantDocuments
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	// Convert personal_address and work_address to JSONB
	personalAddressJSON, err := json.Marshal(req.PersonalAddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid personal_address JSON: " + err.Error()})
		return
	}
	workAddressJSON, err := json.Marshal(req.WorkAddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid work_address JSON: " + err.Error()})
		return
	}

	app := &models.MerchantApplication{
		MerchantBasicInfo:    req.MerchantBasicInfo,
		MerchantAddress:      models.MerchantAddress{PersonalAddress: personalAddressJSON, WorkAddress: workAddressJSON},
		MerchantBusinessInfo: req.MerchantBusinessInfo,
		MerchantDocuments:    req.MerchantDocuments,
	}

	app, err = h.service.SubmitApplication(c.Request.Context(), app)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to submit application: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, app)
}

// GetApplication godoc
// @Summary Retrieve a merchant application by ID
// @Description Allows an applicant to view the status and details of their merchant application
// @Tags Merchant
// @Produce json
// @Param id path string true "Application ID" format(uuid)
// @Success 200 {object} object{id=string,first_name=string,last_name=string,email=string,phone=string,personal_address=object{street=string,city=string,state=string,postal_code=string,country=string},work_address=object{street=string,city=string,state=string,postal_code=string,country=string},business_name=string,business_type=string,tax_id=string,documents=object{business_license=string,identification=string},status=string,created_at=string} "Application details"
// @Failure 400 {object} object{error=string} "Invalid application ID format"
// @Failure 404 {object} object{error=string} "Application not found"
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
// @Success 200 {object} object{id=string,user_id=string,business_name=string,business_type=string,tax_id=string,personal_address=object{street=string,city=string,state=string,postal_code=string,country=string},work_address=object{street=string,city=string,state=string,postal_code=string,country=string},status=string,created_at=string,updated_at=string} "Merchant account details"
// @Failure 401 {object} object{error=string} "Unauthorized: Missing or invalid authentication"
// @Failure 404 {object} object{error=string} "Merchant account not found"
// @Router /merchant/me [get]
func (h *MerchantHandler) GetMyMerchant(c *gin.Context) {
	userID, ok := c.Get("id")
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