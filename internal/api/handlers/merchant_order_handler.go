package handlers

import (
	"net/http"
	"strconv"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/services/order"
	//"api-customer-merchant/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// MerchantOrderHandler manages merchant order-related API requests.
type MerchantOrderHandler struct {
	orderService *order.OrderService
	logger       *zap.Logger
	validate     *validator.Validate
}

// NewMerchantOrderHandler creates a new MerchantOrderHandler instance.
func NewMerchantOrderHandler(orderService *order.OrderService, logger *zap.Logger) *MerchantOrderHandler {
	return &MerchantOrderHandler{
		orderService: orderService,
		logger:       logger,
		validate:     validator.New(),
	}
}

// GetMerchantOrders retrieves all orders for a merchant
// @Summary Get merchant orders
// @Description Retrieves all orders containing items from the authenticated merchant
// @Tags Merchant Orders
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.MerchantOrderResponse
// @Failure 401 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /merchant/orders [get]
func (h *MerchantOrderHandler) GetMerchantOrders(c *gin.Context) {
	ctx := c.Request.Context()

	// Get merchant ID from context
	merchantID, exists := c.Get("merchantID")
	if !exists {
		h.logger.Warn("Unauthorized access to GetMerchantOrders")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	merchantIDStr, ok := merchantID.(string)
	if !ok || merchantIDStr == "" {
		h.logger.Warn("Invalid merchant ID in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid merchant ID"})
		return
	}

	orders, err := h.orderService.GetMerchantOrders(ctx, merchantIDStr)
	if err != nil {
		h.logger.Error("Failed to get merchant orders", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve orders"})
		return
	}

	// Convert to merchant order response DTOs
	var responses []dto.MerchantOrderResponse
	for _, order := range orders {
		var items []dto.MerchantOrderItemResponse
		for _, item := range order.OrderItems {
			// Only include items belonging to this merchant
			if item.MerchantID == merchantIDStr {
				imageURL := ""
				if len(item.Product.Media) > 0 {
					imageURL = item.Product.Media[0].URL
				}

				items = append(items, dto.MerchantOrderItemResponse{
					ID:                item.ID,
					ProductID:         item.ProductID,
					Name:              item.Product.Name,
					Quantity:          item.Quantity,
					Price:             item.Price,
					Image:             imageURL,
					FulfillmentStatus: string(item.FulfillmentStatus),
				})
			}
		}

		// Get delivery address from user's default address
		deliveryAddress := ""
		for _, addr := range order.User.Addresses {
			if addr.IsDefault {
				deliveryAddress = addr.DeliveryAddress
				break
			}
		}

		responses = append(responses, dto.MerchantOrderResponse{
			ID:              order.ID,
			UserID:          order.UserID,
			Status:          string(order.Status),
			OrderItems:      items,
			TotalAmount:     order.TotalAmount.InexactFloat64(),
			DeliveryAddress: deliveryAddress,
			CreatedAt:       order.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:       order.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, responses)
}

// GetMerchantOrder retrieves a specific order for a merchant
// @Summary Get merchant order by ID
// @Description Retrieves a specific order containing items from the authenticated merchant
// @Tags Merchant Orders
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Success 200 {object} dto.MerchantOrderResponse
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /merchant/orders/{id} [get]
func (h *MerchantOrderHandler) GetMerchantOrder(c *gin.Context) {
	ctx := c.Request.Context()

	// Get merchant ID from context
	merchantID, exists := c.Get("merchantID")
	if !exists {
		h.logger.Warn("Unauthorized access to GetMerchantOrder")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	merchantIDStr, ok := merchantID.(string)
	if !ok || merchantIDStr == "" {
		h.logger.Warn("Invalid merchant ID in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid merchant ID"})
		return
	}

	// Parse order ID
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid order ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	order, err := h.orderService.GetOrder(ctx, uint(orderID))
	if err != nil {
		h.logger.Error("Failed to get order", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve order"})
		return
	}

	// Check if merchant has items in this order
	hasMerchantItems := false
	for _, item := range order.OrderItems {
		if item.MerchantID == merchantIDStr {
			hasMerchantItems = true
			break
		}
	}

	if !hasMerchantItems {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	// Convert to merchant order response DTO
	var items []dto.MerchantOrderItemResponse
	for _, item := range order.OrderItems {
		if item.MerchantID == merchantIDStr {
			imageURL := ""
			if len(item.Product.Media) > 0 {
				imageURL = item.Product.Media[0].URL
			}

			items = append(items, dto.MerchantOrderItemResponse{
				ID:                item.ID,
				ProductID:         item.ProductID,
				Name:              item.Product.Name,
				Quantity:          item.Quantity,
				Price:             item.Price,
				Image:             imageURL,
				FulfillmentStatus: string(item.FulfillmentStatus),
			})
		}
	}

	// Get delivery address from user's default address
	deliveryAddress := ""
	for _, addr := range order.User.Addresses {
		if addr.IsDefault {
			deliveryAddress = addr.DeliveryAddress
			break
		}
	}

	response := dto.MerchantOrderResponse{
		ID:              order.ID,
		UserID:          order.UserID,
		Status:          string(order.Status),
		OrderItems:      items,
		TotalAmount:     order.TotalAmount.InexactFloat64(),
		DeliveryAddress: deliveryAddress,
		CreatedAt:       order.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       order.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, response)
}

// AcceptOrderItem allows a merchant to accept an order item
// @Summary Accept an order item
// @Description Allows a merchant to accept an order item
// @Tags Merchant Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order Item ID"
// @Success 200 {object} object{message=string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /merchant/orders/items/{id}/accept [post]
func (h *MerchantOrderHandler) AcceptOrderItem(c *gin.Context) {
	ctx := c.Request.Context()

	// Get merchant ID from context
	merchantID, exists := c.Get("merchantID")
	if !exists {
		h.logger.Warn("Unauthorized access to AcceptOrderItem")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	merchantIDStr, ok := merchantID.(string)
	if !ok || merchantIDStr == "" {
		h.logger.Warn("Invalid merchant ID in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid merchant ID"})
		return
	}

	// Parse order item ID
	orderItemIDStr := c.Param("id")
	orderItemID, err := strconv.ParseUint(orderItemIDStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid order item ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order item ID"})
		return
	}

	// Call service to accept the order item
	if err := h.orderService.AcceptOrderItem(ctx, uint(orderItemID), merchantIDStr); err != nil {
		h.logger.Error("Failed to accept order item", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order item accepted successfully"})
}

// DeclineOrderItem allows a merchant to decline an order item
// @Summary Decline an order item
// @Description Allows a merchant to decline an order item
// @Tags Merchant Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order Item ID"
// @Success 200 {object} object{message=string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /merchant/orders/items/{id}/decline [post]
func (h *MerchantOrderHandler) DeclineOrderItem(c *gin.Context) {
	ctx := c.Request.Context()

	// Get merchant ID from context
	merchantID, exists := c.Get("merchantID")
	if !exists {
		h.logger.Warn("Unauthorized access to DeclineOrderItem")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	merchantIDStr, ok := merchantID.(string)
	if !ok || merchantIDStr == "" {
		h.logger.Warn("Invalid merchant ID in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid merchant ID"})
		return
	}

	// Parse order item ID
	orderItemIDStr := c.Param("id")
	orderItemID, err := strconv.ParseUint(orderItemIDStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid order item ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order item ID"})
		return
	}

	// Call service to decline the order item
	if err := h.orderService.DeclineOrderItem(ctx, uint(orderItemID), merchantIDStr); err != nil {
		h.logger.Error("Failed to decline order item", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order item declined successfully"})
}

// UpdateOrderItemToSentToAronovaHub allows a merchant to update an order item to "SentToAronovaHub" status
// @Summary Update order item to SentToAronovaHub
// @Description Allows a merchant to update an order item to "SentToAronovaHub" status (only allowed after acceptance)
// @Tags Merchant Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order Item ID"
// @Success 200 {object} object{message=string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /merchant/orders/items/{id}/sent-to-aronova-hub [post]
func (h *MerchantOrderHandler) UpdateOrderItemToSentToAronovaHub(c *gin.Context) {
	ctx := c.Request.Context()

	// Get merchant ID from context
	merchantID, exists := c.Get("merchantID")
	if !exists {
		h.logger.Warn("Unauthorized access to UpdateOrderItemToSentToAronovaHub")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	merchantIDStr, ok := merchantID.(string)
	if !ok || merchantIDStr == "" {
		h.logger.Warn("Invalid merchant ID in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid merchant ID"})
		return
	}

	// Parse order item ID
	orderItemIDStr := c.Param("id")
	orderItemID, err := strconv.ParseUint(orderItemIDStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid order item ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order item ID"})
		return
	}

	// Call service to update the order item status
	if err := h.orderService.UpdateOrderItemToSentToAronovaHub(ctx, uint(orderItemID), merchantIDStr); err != nil {
		h.logger.Error("Failed to update order item to SentToAronovaHub", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order item updated to SentToAronovaHub successfully"})
}
