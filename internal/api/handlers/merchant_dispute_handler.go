package handlers

import (
	"net/http"

	//"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/services/dispute"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type MerchantDisputeHandler struct {
	service  *dispute.DisputeService
	validate *validator.Validate
}

func NewMerchantDisputeHandler(service *dispute.DisputeService) *MerchantDisputeHandler {
	return &MerchantDisputeHandler{
		service:  service,
		validate: validator.New(),
	}
}

// ListMerchantDisputes handles GET /merchant/disputes
// @Summary List merchant's disputes
// @Description Retrieve all disputes for the authenticated merchant, grouped by order
// @Tags Merchant Disputes
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.DisputeResponseDTO
// @Failure 401 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /merchant/disputes [get]
func (h *MerchantDisputeHandler) ListMerchantDisputes(c *gin.Context) {
	merchantID, exists := c.Get("merchantID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	disputes, err := h.service.GetMerchantDisputes(c.Request.Context(), merchantID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, disputes)
}

// UpdateDispute handles PUT /merchant/disputes/:id
// @Summary Update a dispute
// @Description Update the status and resolution of a dispute for the authenticated merchant
// @Tags Merchant Disputes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Dispute ID"
// @Param body body object{status=string,resolution=string} true "Dispute update details"
// @Success 200 {object} object{message=string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /merchant/disputes/{id} [put]
func (h *MerchantDisputeHandler) UpdateDispute(c *gin.Context) {
	merchantID, exists := c.Get("merchantID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	disputeID := c.Param("id")
	if disputeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dispute ID required"})
		return
	}

	var req struct {
		Status     string `json:"status" binding:"required,oneof=open resolved rejected"`
		Resolution string `json:"resolution"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.UpdateDisputeStatus(c.Request.Context(), disputeID, merchantID.(string), req.Status, req.Resolution)
	if err != nil {
		if err == dispute.ErrUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "dispute updated successfully"})
}
