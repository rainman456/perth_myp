// api/handlers/wishlist_handler.go
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/services/wishlist"
)

type WishlistHandler struct {
	service *wishlist.WishlistService
}

func NewWishlistHandler(service *wishlist.WishlistService) *WishlistHandler {
	return &WishlistHandler{service: service}
}

// AddToWishlist godoc
// @Summary Add a product to the wishlist
// @Description Add a product to the user's wishlist
// @Tags wishlist
// @Accept json
// @Produce json
// @Param body body dto.AddItemRequest true  "Wishlist body"
// @Success 201 {object} object{message=string} "product added to wishlist"
// @Failure 400 {object} object{error=string}
// @Router /wishlist [post]
func (h *WishlistHandler) AddToWishlist(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context - authentication required"})
		return
	}
	userID, ok := userIDInterface.(string)
	if !ok || userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID from authentication"})
		return
	}
	userIDUint, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
	
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}
	var input dto.AddWishlistItemDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productID := input.ProductID
	if err := h.service.AddToWishlist(c.Request.Context(), uint(userIDUint), productID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "product added to wishlist"})
}

// RemoveFromWishlist godoc
// @Summary Remove a product from the wishlist
// @Description Remove a product from the user's wishlist
// @Tags wishlist
// @Accept json
// @Produce json
// @Param productID path string true "Product ID"
// @Success 204 {object} object{message=string}
// @Failure 400 {object}  object{error=string}
// @Router /wishlist/{productID} [delete]
func (h *WishlistHandler) RemoveFromWishlist(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user session"})
		return
	}
	userIDUint, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	productID := c.Param("productID")
	if err := h.service.RemoveFromWishlist(c.Request.Context(), uint(userIDUint), productID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}


// GetWishlist godoc
// @Summary Get the user's wishlist
// @Description Retrieve the user's wishlist
// @Tags wishlist
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.WishlistResponseDTO
// @Failure 500 {object} object{error=string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Router /wishlist [get]  // Fixed path comment to match routes.go
func (h *WishlistHandler) GetWishlist(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user session"})
		return
	}
	userIDUint, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	products, err := h.service.GetWishlist(c.Request.Context(), uint(userIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products) // Or map to DTO if needed
}

// IsInWishlist godoc
// @Summary Check if a product is in the wishlist
// @Description Check if a specific product is in the user's wishlist
// @Tags wishlist
// @Accept json
// @Produce json
// @Param productID path string true "Product ID"
// @Success 200 {object} object{message=string}
// @Failure 500 {object} object{error=string}
// @Router /wishlist/{productID}/check [get]
func (h *WishlistHandler) IsInWishlist(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user session"})
		return
	}
	userIDUint, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	productID := c.Param("productID")
	isIn, err := h.service.IsInWishlist(c.Request.Context(), uint(userIDUint), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"is_in_wishlist": isIn})
}

// ClearWishlist godoc
// @Summary Clear the user's wishlist
// @Description Remove all products from the user's wishlist
// @Tags wishlist
// @Accept json
// @Produce json
// @Success 204
// @Failure 500 {object} object{error=string}
// @Router /wishlist/clear [delete]
func (h *WishlistHandler) ClearWishlist(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user session"})
		return
	}
	userIDUint, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	if err := h.service.ClearWishlist(c.Request.Context(), uint(userIDUint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}