package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/services/cart" // Assuming service import
	"api-customer-merchant/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)



type CartHandler struct {
	cartService *cart.CartService
	logger      *zap.Logger
	validate    *validator.Validate
}

func NewCartHandler(cartService *cart.CartService, logger *zap.Logger) *CartHandler {
	return &CartHandler{
		cartService: cartService,
		logger:      logger,
		validate:    validator.New(),
	}
}

// AddToCart handles adding an item to the cart
// @Summary Add item to cart
// @Description Adds a product (with optional variant) to user's active cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.AddItemRequest true "Item details"
// @Success 200 {object} dto.CartResponse
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Router /cart/items [post]
func (h *CartHandler) AddToCart(c *gin.Context) {
	ctx := c.Request.Context()
	//userIDStr := c.Query("user_id") // For testing
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	//userID, _ := strconv.ParseUint(userIDStr, 10, 32)  // Helper to parse, assume implemented

	var req dto.AddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Bind error", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.validate.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		h.logger.Error("Invalid userID type in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user session"})
		return
	}
	userIDUint, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		h.logger.Error("Failed to parse userID", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	updatedCart, err := h.cartService.AddItemToCart(ctx, uint(userIDUint), req.Quantity, req.ProductID, req.VariantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Helper to map model to DTO, assume implemented
	resp := &dto.CartResponse{}
	if err := utils.RespMap(updatedCart, resp); err != nil {
		h.logger.Error(" error", zap.Error(err))
	}
	c.JSON(http.StatusOK, resp)
}

// GetCartItem handles getting a cart item
// GetCartItem godoc
// @Summary Get cart item by ID
// @Description Retrieves a specific cart item
// @Tags Cart
// @Produce json
// @Security BearerAuth
// @Param id path string true "Cart Item ID"
// @Success 200 {object} dto.CartItemResponse
// @Failure 401 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /cart/items/{id} [get]
func (h *CartHandler) GetCartItem(c *gin.Context) {
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	ctx := c.Request.Context()
	itemIDStr := c.Param("id")
	itemID, _ := strconv.ParseUint(itemIDStr, 10, 32)

	item, err := h.cartService.GetCartItemByID(ctx, uint(itemID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// GetCart handles getting the active cart
// GetCart godoc
// @Summary Get active cart
// @Description Retrieves the user's active cart with items
// @Tags Cart
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.CartResponse
// @Failure 401 {object} object{error=string}
// @Router /cart [get]
func (h *CartHandler) GetCart(c *gin.Context) {
	ctx := c.Request.Context()
	//userIDStr := c.Query("user_id")
	//userID, _ := strconv.ParseUint(userIDStr, 10, 32)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		h.logger.Error("Invalid userID type in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user session"})
		return
	}
	userIDUint, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		h.logger.Error("Failed to parse userID", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	cart, err := h.cartService.GetActiveCart(ctx,uint(userIDUint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cart)
}

// UpdateCartItemQuantity godoc
// @Summary Update cart item quantity
// @Description Updates the quantity of a cart item
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Cart Item ID"
// @Param body body dto.UpdateItemRequest true "New quantity"
// @Success 200 {object} dto.CartResponse
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /cart/items/{id} [put]
func (h *CartHandler) UpdateCartItemQuantity(c *gin.Context) {
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	ctx := c.Request.Context()
	itemIDstr := strings.TrimSpace(c.Param("id"))
	itemID, _ := strconv.ParseUint(itemIDstr, 10, 32)

	var req dto.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.validate.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	

	updatedCart, err := h.cartService.UpdateCartItemQuantity(ctx, uint(itemID), req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := &dto.CartResponse{}
	if err := utils.RespMap(updatedCart, resp); err != nil {
		h.logger.Error(" error", zap.Error(err))
	}
	c.JSON(http.StatusOK, resp)

}

// RemoveCartItem handles removing an item
// RemoveCartItem godoc
// @Summary Remove cart item
// @Description Removes a cart item by ID
// @Tags Cart
// @Produce json
// @Security BearerAuth
// @Param id path string true "Cart Item ID"
// @Success 200 {object} dto.CartResponse
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /cart/items/{id} [delete]
func (h *CartHandler) RemoveCartItem(c *gin.Context) {
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	ctx := c.Request.Context()
	itemIDStr := c.Param("id")
	itemID, _ := strconv.ParseUint(itemIDStr, 10, 32)

	updatedCart, err := h.cartService.RemoveCartItem(ctx, uint(itemID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedCart)
}

// ClearCart handles DELETE /cart or POST /cart/clear
// @Summary Clear the cart
// @Description Clears all items from the user's active cart
// @Tags Cart
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{message=string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Router /cart/clear [post]
func (h *CartHandler) ClearCart(c *gin.Context) {
	ctx := c.Request.Context()
	userID, exists := c.Get("userID")
	if !exists {
		h.logger.Warn("Unauthorized access to ClearCart")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		h.logger.Error("Invalid userID type in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user session"})
		return
	}
	userIDUint, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		h.logger.Error("Failed to parse userID", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	perr := h.cartService.ClearCart(ctx, uint(userIDUint))
	if perr != nil {
		h.logger.Error("ClearCart failed", zap.Uint("user_id", userID.(uint)), zap.Error(perr))
		c.JSON(http.StatusBadRequest, gin.H{"error": perr.Error()})
		return
	}

	h.logger.Info("Cart cleared successfully", zap.Uint("user_id", userID.(uint)))
	c.JSON(http.StatusOK, gin.H{"message": "cart cleared"})
}

// BulkAddItems handles POST /cart/bulk
// @Summary Bulk add items to cart
// @Description Adds multiple items to the user's active cart in one request
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.BulkUpdateRequest true "Bulk items details"
// @Success 200 {object} dto.CartResponse
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Router /cart/bulk [post]
func (h *CartHandler) BulkAddItems(c *gin.Context) {
	ctx := c.Request.Context()
	userID, exists := c.Get("userID")
	if !exists {
		h.logger.Warn("Unauthorized access to BulkAddItems")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.BulkUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Bind error", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.validate.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		h.logger.Error("Invalid userID type in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user session"})
		return
	}
	userIDUint, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		h.logger.Error("Failed to parse userID", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	updatedCart, err := h.cartService.BulkAddItems(ctx, uint(userIDUint), req)
	if err != nil {
		h.logger.Error("BulkAddItems failed", zap.Uint("user_id", userID.(uint)), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := &dto.CartResponse{}
	if err := utils.RespMap(updatedCart, resp); err != nil {
		h.logger.Error("Response mapping error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	h.logger.Info("Bulk items added successfully", zap.Uint("user_id", userID.(uint)), zap.Int("item_count", len(req.Items)))
	c.JSON(http.StatusOK, resp)
}
