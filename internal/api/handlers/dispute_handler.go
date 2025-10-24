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
// @Success 201 {object} dto.CreateDisputeResponseDTO
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Router /disputes [post]
func (h *DisputeHandler) CreateDispute(c *gin.Context) {
	userID := getUserIDFromContext(c)
	// if !exists {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	// 	return
	// }

	var req dto.CreateDisputeDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.CreateDispute(c.Request.Context(), userID, req)
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
// @Success 200 {object} dto.CreateDisputeResponseDTO
// @Failure 404 {object} object{error=string}
// @Router /disputes/{id} [get]
func (h *DisputeHandler) GetDispute(c *gin.Context) {
	userID := getUserIDFromContext(c)
	disputeID := c.Param("id")

	resp, err := h.service.GetDispute(c.Request.Context(), disputeID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetDisputesByOrderID handles GET /dispute/:orderId
// @Summary Get disputes for a specific order
// @Description Retrieve all disputes for a specific order by order ID for the authenticated customer
// @Tags Disputes
// @Produce json
// @Security BearerAuth
// @Param orderId path string true "Order ID"
// @Success 200 {object} dto.DisputeResponseDTO
// @Failure 400 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Router /disputes/order/{id} [get]
func (h *DisputeHandler) GetDisputesByOrderID(c *gin.Context) {
	userID := getUserIDFromContext(c)
	// if !exists {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	// 	return
	// }

	orderID := c.Param("id")

	resp, err := h.service.GetDisputesByOrderID(c.Request.Context(), orderID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListCustomerDisputes handles GET /disputes
// @Summary List customer's disputes
// @Description Retrieve all disputes for the authenticated customer, grouped by order
// @Tags Disputes
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.DisputeResponseDTO
// @Failure 401 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /disputes [get]
func (h *DisputeHandler) ListCustomerDisputes(c *gin.Context) {
	userID := getUserIDFromContext(c)
	// if !exists {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	// 	return
	// }

	disputes, err := h.service.GetCustomerDisputes(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, disputes)
}