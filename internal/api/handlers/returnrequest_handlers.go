package handlers

import (
    "net/http"

    "api-customer-merchant/internal/api/dto"
    "api-customer-merchant/internal/services/return_request"

    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
)

type ReturnRequestHandler struct {
    service  *return_request.ReturnRequestService
    validate *validator.Validate
}

func NewReturnRequestHandler(service *return_request.ReturnRequestService) *ReturnRequestHandler {
    return &ReturnRequestHandler{
        service:  service,
        validate: validator.New(),
    }
}

// CreateReturnRequest handles POST /return-requests
// @Summary Create a return request
// @Description Customer creates a return request for an order item
// @Tags Return Requests
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.CreateReturnRequestDTO true "Return request details"
// @Success 201 {object} dto.ReturnRequestResponseDTO
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Router /return-requests [post]
func (h *ReturnRequestHandler) CreateReturnRequest(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    var req dto.CreateReturnRequestDTO
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.validate.Struct(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    resp, err := h.service.CreateReturnRequest(c.Request.Context(), userID.(uint), req)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, resp)
}

// GetReturnRequest handles GET /return-requests/:id
// @Summary Get return request details
// @Description Retrieve a specific return request by ID
// @Tags Return Requests
// @Produce json
// @Security BearerAuth
// @Param id path string true "Return Request ID"
// @Success 200 {object} dto.ReturnRequestResponseDTO
// @Failure 404 {object} object{error=string}
// @Router /return-requests/{id} [get]
func (h *ReturnRequestHandler) GetReturnRequest(c *gin.Context) {
    userID, _ := c.Get("userID")
    requestID := c.Param("id")

    resp, err := h.service.GetReturnRequest(c.Request.Context(), requestID, userID.(uint))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, resp)
}

// ListCustomerReturnRequests handles GET /return-requests
// @Summary List customer's return requests
// @Description Retrieve all return requests for the authenticated customer
// @Tags Return Requests
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.ReturnRequestResponseDTO
// @Failure 401 {object} object{error=string}
// @Router /return-requests [get]
func (h *ReturnRequestHandler) ListCustomerReturnRequests(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    returnRequests, err := h.service.GetCustomerReturnRequests(c.Request.Context(), userID.(uint))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, returnRequests)
}