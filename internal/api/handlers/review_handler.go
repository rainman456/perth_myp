// api/handlers/review_handler.go (boilerplate assuming Gin framework)
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/services/review"
)

type ReviewHandler struct {
	service *review.ReviewService
}

func NewReviewHandler(service *review.ReviewService) *ReviewHandler {
	return &ReviewHandler{service: service}
}

// CreateReview godoc
// @Summary Create a review
// @Description Create a new review by user
// @Tags reviews
// @Accept json
// @Produce json
// @Param input body dto.CreateReviewDTO true "Review input"
// @Success 201 {object} dto.ReviewResponseDTO
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Router /review [post]
func (h *ReviewHandler) CreateReview(c *gin.Context) {
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
	var input dto.CreateReviewDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.service.CreateReview(c.Request.Context(), uint(userIDUint), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// GetReview godoc
// @Summary Get a review by ID
// @Description Retrieve a review by its ID
// @Tags reviews
// @Produce json
// @Param id path int true "Review ID"
// @Success 200 {object} dto.ReviewResponseDTO
// @Failure 404 {object}  object{error=string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Router /reviews/{id} [get]
func (h *ReviewHandler) GetReview(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)
	resp, err := h.service.GetReview(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetReviewsByProduct godoc
// @Summary Get reviews by product ID
// @Description Retrieve reviews for a specific product
// @Tags reviews
// @Produce json
// @Param productID path string true "Product ID"
// @Param limit query int false "Limit of reviews" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {array} dto.ReviewResponseDTO
// @Failure 500 {object} object{error=string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Router /products/{productID}/reviews [get]
func (h *ReviewHandler) GetReviewsByProduct(c *gin.Context) {
	productID := c.Param("productID")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	reviews, err := h.service.GetReviewsByProduct(c.Request.Context(), productID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reviews)
}

// UpdateReview godoc
// @Summary Update a review
// @Description Update an existing review by ID
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "Review ID"
// @Param input body dto.UpdateReviewDTO true "Updated review input"
// @Success 200 {object} dto.ReviewResponseDTO
// @Failure 400 {object}  object{error=string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Router /reviews/{id} [put]
func (h *ReviewHandler) UpdateReview(c *gin.Context) {
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
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)
	var input dto.UpdateReviewDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.service.UpdateReview(c.Request.Context(), uint(id), uint(userIDUint), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// DeleteReview godoc
// @Summary Delete a review
// @Description Delete a review by ID
// @Tags reviews
// @Param id path int true "Review ID"
// @Success 204
// @Failure 400 {object}  object{error=string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Router /reviews/{id} [delete]
func (h *ReviewHandler) DeleteReview(c *gin.Context) {
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
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)
	perr := h.service.DeleteReview(c.Request.Context(), uint(id), uint(userIDUint))
	if perr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": perr.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}