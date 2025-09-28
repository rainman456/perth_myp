package handlers

import (
	"net/http"
	"strconv"

	"api-customer-merchant/internal/services/order"
	"github.com/gin-gonic/gin"
)

// OrderHandler manages order-related API requests.
type OrderHandler struct {
	orderService *order.OrderService
}

// NewOrderHandler creates a new OrderHandler instance.
func NewOrderHandler(orderService *order.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// CreateOrder handles the creation of a new order.
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
