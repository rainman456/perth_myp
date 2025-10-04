package handlers

import (
	"net/http"
	//"strconv"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/services/payment"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	service *payment.PaymentService
}

func NewPaymentHandler(s *payment.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: s}
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
	var req dto.InitializePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.service.InitializeCheckout(c.Request.Context(), req)
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
	reference := c.Param("reference")
	resp, err := h.service.VerifyPayment(c.Request.Context(), reference)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}




func (s *PaymentHandler) HandleWebhook(c *gin.Context) {
	

	// var event paystack.Event
	// if err := c.ShouldBindJSON(&event); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
	// 	return
	// }

	// // Verify signature
	// sig := c.GetHeader("x-paystack-signature")
	// if !paystack.VerifySignature([]byte(event.Raw), sig, s.conf.PaystackSecretKey) {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
	// 	return
	// }

	// // Handle charge.success
	// if event.Event == "charge.success" {
	// 	var data paystack.TransactionData
	// 	if err := json.Unmarshal(event.Data, &data); err != nil {
	// 		logger.Error("Unmarshal webhook data failed", zap.Error(err))
	// 		c.Status(http.StatusBadRequest)
	// 		return
	// 	}
	// 	paymentResp, err := s.VerifyPayment(c.Request.Context(), data.Reference)
	// 	if err != nil {
	// 		logger.Error("Webhook verification failed", zap.Error(err))
	// 		c.Status(http.StatusInternalServerError)
	// 		return
	// 	}
	// 	// Trigger notifications, payouts, etc.
	// }

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}