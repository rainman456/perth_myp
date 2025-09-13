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
// @Description Allows a prospective merchant to submit an application with personal, business, and address information.
// @Tags Merchant
// @Accept json
// @Produce json
// @Param body body models.MerchantApplication true "Merchant application payload"
// @Success 201 {object} models.MerchantApplication "Successfully created application"
// @Failure 400 {object} map[string]string "Invalid request body or malformed JSON"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /merchant/apply [post]
// @Schema models.MerchantApplication
// {
//   "type": "object",
//   "properties": {
//     "id": { "type": "string", "format": "uuid", "description": "Unique application ID" },
//     "first_name": { "type": "string", "description": "Merchant's first name" },
//     "last_name": { "type": "string", "description": "Merchant's last name" },
//     "email": { "type": "string", "format": "email", "description": "Merchant's email address" },
//     "phone": { "type": "string", "description": "Merchant's phone number" },
//     "personal_address": {
//       "type": "object",
//       "description": "Merchant's personal address",
//       "properties": {
//         "street": { "type": "string" },
//         "city": { "type": "string" },
//         "state": { "type": "string" },
//         "postal_code": { "type": "string" },
//         "country": { "type": "string" }
//       },
//       "required": ["street", "city", "state", "postal_code", "country"]
//     },
//     "work_address": {
//       "type": "object",
//       "description": "Merchant's work address",
//       "properties": {
//         "street": { "type": "string" },
//         "city": { "type": "string" },
//         "state": { "type": "string" },
//         "postal_code": { "type": "string" },
//         "country": { "type": "string" }
//       },
//       "required": ["street", "city", "state", "postal_code", "country"]
//     },
//     "business_name": { "type": "string", "description": "Name of the business" },
//     "business_type": { "type": "string", "description": "Type of business (e.g., Retail, Service)" },
//     "tax_id": { "type": "string", "description": "Business tax identification number" },
//     "documents": {
//       "type": "object",
//       "description": "Merchant's documents",
//       "properties": {
//         "business_license": { "type": "string", "description": "Business license number" },
//         "identification": { "type": "string", "description": "Personal identification number" }
//       }
//     },
//     "status": { "type": "string", "description": "Application status (e.g., pending, approved, rejected)" },
//     "created_at": { "type": "string", "format": "date-time", "description": "Application creation timestamp" }
//   },
//   "required": ["first_name", "last_name", "email", "phone", "personal_address", "work_address", "business_name", "business_type", "tax_id"]
// }
// @Example request
// {
//   "first_name": "John",
//   "last_name": "Doe",
//   "email": "john.doe@example.com",
//   "phone": "+1234567890",
//   "personal_address": {
//     "street": "123 Main St",
//     "city": "Anytown",
//     "state": "CA",
//     "postal_code": "12345",
//     "country": "USA"
//   },
//   "work_address": {
//     "street": "456 Business Rd",
//     "city": "Anytown",
//     "state": "CA",
//     "postal_code": "12345",
//     "country": "USA"
//   },
//   "business_name": "Doe Enterprises",
//   "business_type": "Retail",
//   "tax_id": "12-3456789",
//   "documents": {
//     "business_license": "BL123456",
//     "identification": "ID789012"
//   }
// }
// @Example response 201
// {
//   "id": "123e4567-e89b-12d3-a456-426614174000",
//   "first_name": "John",
//   "last_name": "Doe",
//   "email": "john.doe@example.com",
//   "phone": "+1234567890",
//   "personal_address": {
//     "street": "123 Main St",
//     "city": "Anytown",
//     "state": "CA",
//     "postal_code": "12345",
//     "country": "USA"
//   },
//   "work_address": {
//     "street": "456 Business Rd",
//     "city": "Anytown",
//     "state": "CA",
//     "postal_code": "12345",
//     "country": "USA"
//   },
//   "business_name": "Doe Enterprises",
//   "business_type": "Retail",
//   "tax_id": "12-3456789",
//   "documents": {
//     "business_license": "BL123456",
//     "identification": "ID789012"
//   },
//   "status": "pending",
//   "created_at": "2025-09-13T03:45:00Z"
// }
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
// @Description Allows an applicant to view the status and details of their merchant application.
// @Tags Merchant
// @Produce json
// @Param id path string true "Application ID" format(uuid)
// @Success 200 {object} models.MerchantApplication "Application details"
// @Failure 400 {object} map[string]string "Invalid application ID format"
// @Failure 404 {object} map[string]string "Application not found"
// @Router /merchant/application/{id} [get]
// @Schema models.MerchantApplication
// {
//   "type": "object",
//   "properties": {
//     "id": { "type": "string", "format": "uuid", "description": "Unique application ID" },
//     "first_name": { "type": "string", "description": "Merchant's first name" },
//     "last_name": { "type": "string", "description": "Merchant's last name" },
//     "email": { "type": "string", "format": "email", "description": "Merchant's email address" },
//     "phone": { "type": "string", "description": "Merchant's phone number" },
//     "personal_address": {
//       "type": "object",
//       "description": "Merchant's personal address",
//       "properties": {
//         "street": { "type": "string" },
//         "city": { "type": "string" },
//         "state": { "type": "string" },
//         "postal_code": { "type": "string" },
//         "country": { "type": "string" }
//       },
//       "required": ["street", "city", "state", "postal_code", "country"]
//     },
//     "work_address": {
//       "type": "object",
//       "description": "Merchant's work address",
//       "properties": {
//         "street": { "type": "string" },
//         "city": { "type": "string" },
//         "state": { "type": "string" },
//         "postal_code": { "type": "string" },
//         "country": { "type": "string" }
//       },
//       "required": ["street", "city", "state", "postal_code", "country"]
//     },
//     "business_name": { "type": "string", "description": "Name of the business" },
//     "business_type": { "type": "string", "description": "Type of business (e.g., Retail, Service)" },
//     "tax_id": { "type": "string", "description": "Business tax identification number" },
//     "documents": {
//       "type": "object",
//       "description": "Merchant's documents",
//       "properties": {
//         "business_license": { "type": "string", "description": "Business license number" },
//         "identification": { "type": "string", "description": "Personal identification number" }
//       }
//     },
//     "status": { "type": "string", "description": "Application status (e.g., pending, approved, rejected)" },
//     "created_at": { "type": "string", "format": "date-time", "description": "Application creation timestamp" }
//   }
// }
// @Example response 200
// {
//   "id": "123e4567-e89b-12d3-a456-426614174000",
//   "first_name": "John",
//   "last_name": "Doe",
//   "email": "john.doe@example.com",
//   "phone": "+1234567890",
//   "personal_address": {
//     "street": "123 Main St",
//     "city": "Anytown",
//     "state": "CA",
//     "postal_code": "12345",
//     "country": "USA"
//   },
//   "work_address": {
//     "street": "456 Business Rd",
//     "city": "Anytown",
//     "state": "CA",
//     "postal_code": "12345",
//     "country": "USA"
//   },
//   "business_name": "Doe Enterprises",
//   "business_type": "Retail",
//   "tax_id": "12-3456789",
//   "documents": {
//     "business_license": "BL123456",
//     "identification": "ID789012"
//   },
//   "status": "pending",
//   "created_at": "2025-09-13T03:45:00Z"
// }
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
// @Description Fetches the merchant account details for the authenticated user, if their application has been approved.
// @Tags Merchant
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Merchant "Merchant account details"
// @Failure 401 {object} map[string]string "Unauthorized: Missing or invalid authentication"
// @Failure 404 {object} map[string]string "Merchant account not found"
// @Router /merchant/me [get]
// @Schema models.Merchant
// {
//   "type": "object",
//   "properties": {
//     "id": { "type": "string", "format": "uuid", "description": "Unique merchant ID" },
//     "user_id": { "type": "string", "format": "uuid", "description": "Associated user ID" },
//     "business_name": { "type": "string", "description": "Name of the business" },
//     "business_type": { "type": "string", "description": "Type of business (e.g., Retail, Service)" },
//     "tax_id": { "type": "string", "description": "Business tax identification number" },
//     "personal_address": {
//       "type": "object",
//       "description": "Merchant's personal address",
//       "properties": {
//         "street": { "type": "string" },
//         "city": { "type": "string" },
//         "state": { "type": "string" },
//         "postal_code": { "type": "string" },
//         "country": { "type": "string" }
//       },
//       "required": ["street", "city", "state", "postal_code", "country"]
//     },
//     "work_address": {
//       "type": "object",
//       "description": "Merchant's work address",
//       "properties": {
//         "street": { "type": "string" },
//         "city": { "type": "string" },
//         "state": { "type": "string" },
//         "postal_code": { "type": "string" },
//         "country": { "type": "string" }
//       },
//       "required": ["street", "city", "state", "postal_code", "country"]
//     },
//     "status": { "type": "string", "description": "Merchant status (e.g., active, inactive)" },
//     "created_at": { "type": "string", "format": "date-time", "description": "Merchant creation timestamp" },
//     "updated_at": { "type": "string", "format": "date-time", "description": "Merchant last updated timestamp" }
//   },
//   "required": ["id", "user_id", "business_name", "business_type", "tax_id", "personal_address", "work_address", "status"]
// }
// @Example response 200
// {
//   "id": "123e4567-e89b-12d3-a456-426614174000",
//   "user_id": "987e6543-e21b-12d3-a456-426614174000",
//   "business_name": "Doe Enterprises",
//   "business_type": "Retail",
//   "tax_id": "12-3456789",
//   "personal_address": {
//     "street": "123 Main St",
//     "city": "Anytown",
//     "state": "CA",
//     "postal_code": "12345",
//     "country": "USA"
//   },
//   "work_address": {
//     "street": "456 Business Rd",
//     "city": "Anytown",
//     "state": "CA",
//     "postal_code": "12345",
//     "country": "USA"
//   },
//   "status": "approved",
//   "created_at": "2025-09-13T03:45:00Z",
//   "updated_at": "2025-09-13T03:45:00Z"
// }
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