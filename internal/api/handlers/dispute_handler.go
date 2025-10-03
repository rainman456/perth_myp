package handlers

import (
	"net/http"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/services/dispute"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type DisputeHandler struct {
	service  *dispute.DisputeService
	validate *validator.Validate
}

func NewDisputeHandler(service *dispute.DisputeService) *DisputeHandler {
	return &DisputeHandler{
		service:  service,
		validate: validator.New(),
	}
}

// CreateDispute handles POST /disputes
// @Summary Create a dispute
// @Description Customer creates a dispute for an order
// @Tags Disputes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.CreateDisputeDTO true "Dispute details"
// @Success 201 {object} dto.DisputeResponseDTO
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Router /disputes [post]
func (h *DisputeHandler) CreateDispute(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.CreateDisputeDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.CreateDispute(c.Request.Context(), userID.(uint), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetDispute handles GET /disputes/:id
// @Summary Get dispute details
// @Description Retrieve a specific dispute by ID
// @Tags Disputes
// @Produce json
// @Security BearerAuth
// @Param id path string true "Dispute ID"
// @Success 200 {object} dto.DisputeResponseDTO
// @Failure 404 {object} object{error=string}
// @Router /disputes/{id} [get]
func (h *DisputeHandler) GetDispute(c *gin.Context) {
	userID, _ := c.Get("userID")
	disputeID := c.Param("id")

	resp, err := h.service.GetDispute(c.Request.Context(), disputeID, userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}