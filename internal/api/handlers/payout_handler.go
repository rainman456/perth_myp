package handlers

import (
	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/services/payout"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PayoutHandler struct {
	payoutService *payout.PayoutService
	logger        *zap.Logger
}

func NewPayoutHandler(payoutService *payout.PayoutService, logger *zap.Logger) *PayoutHandler {
	return &PayoutHandler{
		payoutService: payoutService,
		logger:        logger,
	}
}

// GetMerchantPayouts retrieves all payouts for a merchant
// @Summary Get merchant payouts
// @Description Retrieves all payouts for the authenticated merchant
// @Tags Merchant Payouts
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.PayoutResponse
// @Failure 401 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /merchant/payouts [get]
func (h *PayoutHandler) GetMerchantPayouts(c *gin.Context) {
	//ctx := c.Request.Context()

	// Get merchant ID from context
	merchantID, exists := c.Get("merchantID")
	if !exists {
		h.logger.Warn("Unauthorized access to GetMerchantPayouts")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	merchantIDStr, ok := merchantID.(string)
	if !ok || merchantIDStr == "" {
		h.logger.Warn("Invalid merchant ID in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid merchant ID"})
		return
	}

	payouts, err := h.payoutService.GetPayoutsByMerchantID(merchantIDStr)
	if err != nil {
		h.logger.Error("Failed to get merchant payouts", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve payouts"})
		return
	}

	// Convert to payout response DTOs
	var responses []dto.PayoutResponse
	for _, payout := range payouts {
		responses = append(responses, dto.PayoutResponse{
			ID:              payout.ID,
			MerchantID:      payout.MerchantID,
			Amount:          payout.Amount,
			Status:          string(payout.Status),
			PayoutAccountID: payout.PayoutAccountID,
			CreatedAt:       payout.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:       payout.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, responses)
}

// RequestPayout requests a new payout for a merchant
// @Summary Request merchant payout
// @Description Requests a new payout for the authenticated merchant
// @Tags Merchant Payouts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.PayoutRequest true "Payout request"
// @Success 200 {object} dto.PayoutResponse
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /merchant/payouts/request [post]
func (h *PayoutHandler) RequestPayout(c *gin.Context) {
	ctx := c.Request.Context()

	merchantID, exists := c.Get("merchantID")
	if !exists {
		h.logger.Warn("Unauthorized access to RequestPayout")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	merchantIDStr, ok := merchantID.(string)
	if !ok || merchantIDStr == "" {
		h.logger.Warn("Invalid merchant ID in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid merchant ID"})
		return
	}

	var req dto.PayoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	// NOW USING THE AMOUNT FROM DTO
	payout, err := h.payoutService.RequestPayout(ctx, merchantIDStr, req.Amount)
	if err != nil {
		h.logger.Error("Failed to request payout", zap.Error(err))
		
		// Better error messages
		if err.Error() == "no eligible balance available" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no eligible balance available"})
			return
		}
		if err.Error() == "requested amount exceeds available balance" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "requested amount exceeds available balance"})
			return
		}
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to request payout"})
		return
	}

	response := dto.PayoutResponse{
		ID:              payout.ID,
		MerchantID:      payout.MerchantID,
		Amount:          payout.Amount,
		Status:          string(payout.Status),
		PayoutAccountID: payout.PayoutAccountID,
		CreatedAt:       payout.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       payout.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, response)
}