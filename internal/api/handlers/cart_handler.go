//package handlers
/*
import (
	//"net/http"
	//"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/services/cart"
	"net/http"
	"strconv"

	//"api-customer-merchant/internal/domain/order"
	//"api-customer-merchant/internal/domain/payment"
	"api-customer-merchant/internal/services/product"
	//"strconv"




	"github.com/gin-gonic/gin"
)

type CartHandlers struct {
	productService *product.ProductService
	cartService    *cart.CartService
	// orderService   *order.OrderService
	// paymentService *payment.PaymentService
}

func NewCartHandlers(productService *product.ProductService, cartService *cart.CartService) *CartHandlers {
	return &CartHandlers{
		productService: productService,
		cartService:    cartService,
		// orderService:   orderService,
		// paymentService: paymentService,
	}
}





func (h *CartHandlers) AddItemToCart(c *gin.Context) {
	// userID, exists := c.Get("userID")
	// if !exists {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	// 	return
	// }

	var req dto.AddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cart, err := h.cartService.AddItemToCart(req.UserID, req.ProductID, req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}


// UpdateCartItemQuantity handles PUT /customer/cart/update/:cartItemID
func (h *CartHandlers) UpdateCartItemQuantity(c *gin.Context) {
	// userID, exists := c.Get("userID")
	// if !exists {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	// 	return
	// }

	cartItemIDStr := c.Param("cartItemID")
	cartItemID, err := strconv.ParseUint(cartItemIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cart item ID"})
		return
	}
	// Verify cart item belongs to user's active cart
	cartItem, err := h.cartService.GetCartItemByID(uint(cartItemID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cart item not found"})
		return
	}

	var req dto.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cart, err := h.cartService.GetActiveCart(req.UserID)
	if err != nil || cart.ID != cartItem.CartID {
		c.JSON(http.StatusForbidden, gin.H{"error": "cart item does not belong to user"})
		return
	}

	updatedCart, err := h.cartService.UpdateCartItemQuantity(uint(cartItemID), req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedCart)
}



// RemoveCartItem handles DELETE /customer/cart/remove/:cartItemID
func (h *CartHandlers) RemoveCartItem(c *gin.Context) {
	 userID, exists := c.Get("userID")
	 if !exists {
	 	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	 	return
	 }

	cartItemIDStr := c.Param("cartItemID")
	cartItemID, err := strconv.ParseUint(cartItemIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cart item ID"})
		return
	}
	// Verify cart item belongs to user's active cart
	cartItem, err := h.cartService.GetCartItemByID(uint(cartItemID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cart item not found"})
		return
	}
	 cart, err := h.cartService.GetActiveCart(userID.(uint))
	 if err != nil || cart.ID != cartItem.CartID {
	 	c.JSON(http.StatusForbidden, gin.H{"error": "cart item does not belong to user"})
	 	return
	 }


	updatedCart, err := h.cartService.RemoveCartItem(uint(cartItemID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedCart)
}
*/

package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/utils"
	"api-customer-merchant/internal/services/cart" // Assuming service import
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// type CartHandler struct {
// 	cartService *cart.CartService
// }

// func NewCartHandler(cartService *cart.CartService) *CartHandler {
// 	return &CartHandler{cartService: cartService}
// }


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
// func (h *CartHandler) AddToCart(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	userIDStr := c.Query("user_id") // For testing, get from query/body
// 	userID, _ := strconv.ParseUint(userIDStr, 10, 32)
// 	productID := c.Query("product_id")
// 	quantityStr := c.Query("quantity")
// 	quantity, _ := strconv.ParseUint(quantityStr, 10, 32)

// 	updatedCart, err := h.cartService.AddItemToCart(ctx, uint(userID), uint(quantity), productID)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, updatedCart)
// }




// AddToCart godoc
// @Summary Add item to cart
// @Description Adds a product variant to the user's cart or increments quantity if exists
// @Tags Cart
// @Accept json
// @Produce json
// @Param body body dto.AddItemRequest true "Item to add"
// @Param user_id query uint false "User ID (for testing)"
// @Security BearerAuth
// @Success 200 {object} dto.CartResponse "Item added successfully"
// @Failure 400 {object} object{error=string} "Invalid request or validation failed"
// @Failure 500 {object} object{error=string} "Failed to add item"
// @Router /cart/items [post]
func (h *CartHandler) AddToCart(c *gin.Context) {
	ctx := c.Request.Context()
	userIDStr := c.Query("user_id") // For testing
	userID, _ := strconv.ParseUint(userIDStr, 10, 32)  // Helper to parse, assume implemented

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

	updatedCart, err := h.cartService.AddItemToCart(ctx, uint(userID), req.Quantity, req.ProductID, req.VariantID)
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
// @Param id path uint true "Cart item ID"
// @Success 200 {object} dto.CartItemResponse "Cart item retrieved successfully"
// @Failure 400 {object} object{error=string} "Invalid cart item ID"
// @Failure 404 {object} object{error=string} "Cart item not found"
// @Failure 500 {object} object{error=string} "Failed to fetch cart item"
// @Router /cart/items/{id} [get]
func (h *CartHandler) GetCartItem(c *gin.Context) {
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
// @Description Retrieves the user's active cart
// @Tags Cart
// @Produce json
// @Param user_id query uint false "User ID (for testing)"
// @Security BearerAuth
// @Success 200 {object} dto.CartResponse "Cart retrieved successfully"
// @Failure 400 {object} object{error=string} "Invalid user ID"
// @Failure 404 {object} object{error=string} "Cart not found"
// @Failure 500 {object} object{error=string} "Failed to fetch cart"
// @Router /cart [get]
func (h *CartHandler) GetCart(c *gin.Context) {
	ctx := c.Request.Context()
	userIDStr := c.Query("user_id")
	userID, _ := strconv.ParseUint(userIDStr, 10, 32)

	cart, err := h.cartService.GetActiveCart(ctx, uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cart)
}

// UpdateCartItemQuantity handles updating quantity
// func (h *CartHandler) UpdateCartItemQuantity(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	itemIDStr := c.Param("id")
// 	itemID, _ := strconv.ParseUint(itemIDStr, 10, 32)
// 	quantityStr := c.Query("quantity")
// 	quantity, _ := strconv.Atoi(quantityStr)

// 	updatedCart, err := h.cartService.UpdateCartItemQuantity(ctx, uint(itemID), quantity)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, updatedCart)
// }


// UpdateCartItemQuantity godoc
// @Summary Update cart item quantity
// @Description Updates the quantity of a specific cart item
// @Tags Cart
// @Accept json
// @Produce json
// @Param id path uint true "Cart item ID"
// @Param body body dto.UpdateItemRequest true "Updated quantity"
// @Security BearerAuth
// @Success 200 {object} dto.CartResponse "Cart item updated successfully"
// @Failure 400 {object} object{error=string} "Invalid request or validation failed"
// @Failure 404 {object} object{error=string} "Cart item not found"
// @Failure 500 {object} object{error=string} "Failed to update cart item"
// @Router /cart/items/{id} [put]
func (h *CartHandler) UpdateCartItemQuantity(c *gin.Context) {
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
// @Description Removes a specific item from the cart
// @Tags Cart
// @Produce json
// @Param id path uint true "Cart item ID"
// @Security BearerAuth
// @Success 200 {object} dto.CartResponse "Cart item removed successfully"
// @Failure 400 {object} object{error=string} "Invalid cart item ID"
// @Failure 404 {object} object{error=string} "Cart item not found"
// @Failure 500 {object} object{error=string} "Failed to remove cart item"
// @Router /cart/items/{id} [delete]
func (h *CartHandler) RemoveCartItem(c *gin.Context) {
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
