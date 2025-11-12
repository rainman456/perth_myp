package handlers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"io"
	//"io/ioutil"
	"net/http"

	//"strconv"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/config"
	"api-customer-merchant/internal/services/payment"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PaymentHandler struct {
	service *payment.PaymentService
	config  *config.Config // Add config to access PaystackSecretKey
	logger  *zap.Logger    // Add logger if needed
}

func NewPaymentHandler(s *payment.PaymentService, conf *config.Config, logger *zap.Logger) *PaymentHandler {
	return &PaymentHandler{service: s, config: conf, logger: logger}
}

// Initialize payment
// @Summary Initialize payment
// @Description Starts Paystack checkout for an order
// @Tags Payments
// @Accept json
// @Produce json
// @Param body body dto.InitializePaymentRequest true "Payment initialization"
// @Success 201 {object} dto.PaymentResponse
// @Failure 400 {object} object{error=string}
// @Router /payments/initialize [post]
func (h *PaymentHandler) Initialize(c *gin.Context) {
	ctx := c.Request.Context()
	 _, exists := c.Get("userID")
	 if !exists {
	 	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	 	return
	 }
	var req dto.InitializePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.service.InitializeCheckout(ctx, req)
	h.logger.Info("Initializing payment", zap.Uint("order_id", req.OrderID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// Verify payment
// @Summary Verify payment
// @Description Verifies Paystack transaction by reference
// @Tags Payments
// @Produce json
// @Param reference path string true "Transaction reference"
// @Success 200 {object} dto.PaymentResponse
// @Failure 400 {object} object{error=string}
// @Router /payments/verify/{reference} [get]
func (h *PaymentHandler) Verify(c *gin.Context) {
	ctx := c.Request.Context()
	 _, exists := c.Get("userID")
	 if !exists {
	 	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	 	return
	 }
	reference := c.Param("reference")
	resp, err := h.service.VerifyPayment(ctx, reference)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}





// Webhook handles POST /payments/webhook for Paystack events
// @Summary Paystack webhook
// @Description Receives and processes forwarded Paystack events
// @Tags Payments
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string}
// @Failure 400 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /payments/webhook [post]
func (h *PaymentHandler) Webhook(c *gin.Context) {
	// Verify Paystack signature
	if !h.verifyPaystackSignature(c) {
		h.logger.Warn("Invalid Paystack webhook signature")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature"})
		return
	}

	var event map[string]interface{}
	if err := c.ShouldBindJSON(&event); err != nil {
		h.logger.Error("Failed to bind webhook event", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	// Process the event
	if err := h.service.HandleWebhook(c.Request.Context(), event); err != nil {
		h.logger.Error("Failed to handle webhook", zap.Error(err))
		// Return 200 to acknowledge, as per Paystack recommendation
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// Always acknowledge with 200
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}


// verifyPaystackSignature verifies the x-paystack-signature header
func (h *PaymentHandler) verifyPaystackSignature(c *gin.Context) bool {
	signature := c.GetHeader("x-paystack-signature")
	if signature == "" {
		return false
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Error("Failed to read webhook body", zap.Error(err))
		return false
	}
	// Reset body for ShouldBindJSON
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	key := []byte(h.config.PaystackSecretKey)
	hash := hmac.New(sha512.New, key)
	hash.Write(body)
	expected := hex.EncodeToString(hash.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expected))
}