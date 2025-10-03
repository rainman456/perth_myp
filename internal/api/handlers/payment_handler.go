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