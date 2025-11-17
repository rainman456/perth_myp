package handlers

import (
	"net/http"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/services/merchant"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type MerchantBankHandler struct {
	service  *merchant.MerchantService
	validate *validator.Validate
}

func NewMerchantBankHandler(service *merchant.MerchantService) *MerchantBankHandler {
	return &MerchantBankHandler{
		service:  service,
		validate: validator.New(),
	}
}

// CreateBankDetails godoc
// @Summary Add bank details
// @Description Add bank account details for the authenticated merchant
// @Tags Merchant Bank Details
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.BankDetailsRequest true "Bank details"
// @Success 201 {object} dto.BankDetailsResponse "Bank details created successfully"
// @Failure 400 {object} object{error=string} "Invalid request body"
// @Failure 401 {object} object{error=string} "Unauthorized"
// @Failure 500 {object} object{error=string} "Failed to add bank details"
// @Router /merchant/bank-details [post]
func (h *MerchantBankHandler) CreateBankDetails(c *gin.Context) {
	merchantID, exists := c.Get("merchantID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.BankDetailsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bankDetails, err := h.service.AddBankDetails(c.Request.Context(), merchantID.(string), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add bank details: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, bankDetails)
}

// GetBankDetails godoc
// @Summary Get bank details
// @Description Retrieve bank account details for the authenticated merchant
// @Tags Merchant Bank Details
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.BankDetailsResponse "Bank details retrieved successfully"
// @Failure 401 {object} object{error=string} "Unauthorized"
// @Failure 404 {object} object{error=string} "Bank details not found"
// @Failure 500 {object} object{error=string} "Failed to retrieve bank details"
// @Router /merchant/bank-details [get]
func (h *MerchantBankHandler) GetBankDetails(c *gin.Context) {
	merchantID, exists := c.Get("merchantID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	bankDetails, err := h.service.GetBankDetails(c.Request.Context(), merchantID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "bank details not found: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, bankDetails)
}

// UpdateBankDetails godoc
// @Summary Update bank details
// @Description Update bank account details for the authenticated merchant
// @Tags Merchant Bank Details
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.BankDetailsRequest true "Bank details"
// @Success 200 {object} dto.BankDetailsResponse "Bank details updated successfully"
// @Failure 400 {object} object{error=string} "Invalid request body"
// @Failure 401 {object} object{error=string} "Unauthorized"
// @Failure 500 {object} object{error=string} "Failed to update bank details"
// @Router /merchant/bank-details [put]
func (h *MerchantBankHandler) UpdateBankDetails(c *gin.Context) {
	merchantID, exists := c.Get("merchantID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.BankDetailsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bankDetails, err := h.service.UpdateBankDetails(c.Request.Context(), merchantID.(string), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update bank details: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, bankDetails)
}

// DeleteBankDetails godoc
// @Summary Delete bank details
// @Description Delete bank account details for the authenticated merchant
// @Tags Merchant Bank Details
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{message=string} "Bank details deleted successfully"
// @Failure 401 {object} object{error=string} "Unauthorized"
// @Failure 404 {object} object{error=string} "Bank details not found"
// @Failure 500 {object} object{error=string} "Failed to delete bank details"
// @Router /merchant/bank-details [delete]
func (h *MerchantBankHandler) DeleteBankDetails(c *gin.Context) {
	merchantID, exists := c.Get("merchantID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.service.DeleteBankDetails(c.Request.Context(), merchantID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete bank details: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "bank details deleted successfully"})
}