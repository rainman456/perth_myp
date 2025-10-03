package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/services/order"
	"api-customer-merchant/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// OrderHandler manages order-related API requests.
type OrderHandler struct {
	orderService *order.OrderService
	logger      *zap.Logger
	validate    *validator.Validate
}

// NewOrderHandler creates a new OrderHandler instance.
func NewOrderHandler(orderService *order.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// CreateOrder handles the creation of a new order.
// CreateOrder handles the creation of a new order
// @Summary Create order from cart
// @Description Converts user's active cart into an order
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id query string true "User ID (for testing)"
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} object{error=string}
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userIDStr := c.Query("user_id") // For testing, get from query/body
	userID, _ := strconv.ParseUint(userIDStr, 10, 32)
	ctx := c.Request.Context()
	newOrder, err := h.orderService.CreateOrder(ctx, uint(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, newOrder)
	// Get authURL from payment init (modify service to return it)
//c.JSON(http.StatusOK, gin.H{"new_order": newOrder, "paystack_url": authURL})
}

// GetOrder handles the request to retrieve a specific order by ID.
// GetOrder godoc
// @Summary Get order by ID
// @Description Retrieves a specific order with items
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := h.orderService.GetOrder(c.Request.Context(), uint(orderID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Adjust status code as needed
		return
	}

	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}




// CancelOrder handles POST /orders/:id/cancel (user-authenticated)
// CancelOrder handles POST /orders/:id/cancel
// @Summary Cancel an order
// @Description Customer cancels their pending order
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Param body body dto.CancelOrderRequest true "Cancellation reason"
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Router /orders/{id}/cancel [post]
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	ctx := c.Request.Context()
	userIDStr, exists := c.Get("userID")
	if !exists {
		h.logger.Warn("Unauthorized access to CancelOrder")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID, err := strconv.ParseUint(userIDStr.(string), 10, 32)
	if err != nil {
		h.logger.Error("Invalid user ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	orderIDStr := strings.TrimSpace(c.Param("id"))
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid order ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	var req dto.CancelOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Bind error", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.validate.Struct(&req); err != nil {
		h.logger.Error("Validation error", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.orderService.CancelOrder(ctx, uint(orderID), uint(userID), req.Reason)
	if err != nil {
		h.logger.Error("CancelOrder failed", zap.Uint("order_id", uint(orderID)), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch updated order for response
	updatedOrder, err := h.orderService.GetOrder(ctx, uint(orderID))
	if err != nil {
		h.logger.Error("Failed to fetch updated order", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch updated order"})
		return
	}

	resp := &dto.OrderResponse{}
	if err := utils.RespMap(updatedOrder, resp); err != nil {
		h.logger.Error("Response mapping error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	h.logger.Info("Order cancelled successfully", zap.Uint("order_id", uint(orderID)), zap.Uint("user_id", uint(userID)))
	c.JSON(http.StatusOK, resp)
}