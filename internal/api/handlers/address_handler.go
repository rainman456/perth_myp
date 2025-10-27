package handlers

import (
	"net/http"
	"strconv"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/services/user"

	"github.com/gin-gonic/gin"
)

// AddressHandler handles address CRUD
type AddressHandler struct {
	service *user.AddressService
}

func NewAddressHandler(s *user.AddressService) *AddressHandler {
	return &AddressHandler{service: s}
}

// CreateAddress godoc
// @Summary Create a new address
// @Description Create a user address (user must be authenticated)
// @Tags Customer
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body dto.CreateAddressRequest true "Address payload"
// @Success 201 {object} dto.AddressResponse
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /customer/addresses [post]
func (h *AddressHandler) CreateAddress(c *gin.Context) {
	ctx := c.Request.Context()

	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not logged in"})
		return
	}
	userIDStr, ok := userIDVal.(string)
	if !ok || userIDStr == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user session"})
		return
	}
	uid64, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id"})
		return
	}
	userID := uint(uid64)

	var req dto.CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	addr, err := h.service.CreateAddress(ctx, userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := dto.AddressResponse{
		ID:                    addr.ID,
		//Address:               addr.Address,
		PhoneNumber:           addr.PhoneNumber,
		AdditionalPhoneNumber: addr.AdditionalPhoneNumber,
		DeliveryAddress:       addr.DeliveryAddress,
		ShippingAddress:       addr.ShippingAddress,
		State:                 addr.State,
		LGA:                   addr.LGA,
		
	}

	c.JSON(http.StatusCreated, gin.H{"address": resp})
}

// ListAddresses godoc
// @Summary List user's addresses
// @Description Returns all addresses for the authenticated user
// @Tags Customer
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.ListAddressesResponse
// @Failure 401 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /customer/addresses [get]
func (h *AddressHandler) ListAddresses(c *gin.Context) {
	ctx := c.Request.Context()

	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not logged in"})
		return
	}
	userIDStr, ok := userIDVal.(string)
	if !ok || userIDStr == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user session"})
		return
	}
	uid64, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id"})
		return
	}
	userID := uint(uid64)

	addrs, err := h.service.ListAddresses(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(addrs) == 0 {
		empty := make([]dto.AddressResponse, 0)
		c.JSON(http.StatusOK, gin.H{"addresses": dto.ListAddressesResponse{Items: empty, Count: 0}})
		return
	}

	// if addrs.UserID != userID {
	// 	c.JSON(http.StatusForbidden, gin.H{"error": "address does not belong to user"})
	// 	return
	// }

	out := make([]dto.AddressResponse, 0, len(addrs))
	for _, a := range addrs {
		out = append(out, dto.AddressResponse{
			ID:                    a.ID,
			//Address:               a.Address,
			PhoneNumber:           a.PhoneNumber,
			AdditionalPhoneNumber: a.AdditionalPhoneNumber,
			DeliveryAddress:       a.DeliveryAddress,
			ShippingAddress:       a.ShippingAddress,
			State:                 a.State,
			LGA:                   a.LGA,
		})
	}

	c.JSON(http.StatusOK, gin.H{"addresses": dto.ListAddressesResponse{Items: out, Count: len(out)}})
}

// GetAddress godoc
// @Summary Get single address
// @Description Get a single address by id (must belong to authenticated user)
// @Tags Customer
// @Security BearerAuth
// @Produce json
// @Param id path int true "address id"
// @Success 200 {object} dto.AddressResponse
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /customer/addresses/{id} [get]
func (h *AddressHandler) GetAddress(c *gin.Context) {
	ctx := c.Request.Context()
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	id := uint(id64)

	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not logged in"})
		return
	}
	userIDStr, ok := userIDVal.(string)
	if !ok || userIDStr == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user session"})
		return
	}
	uid64, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id"})
		return
	}
	userID := uint(uid64)

	addr, err := h.service.GetAddress(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if addr == nil {
		c.JSON(http.StatusOK, gin.H{"empty": []any{}})
		return
	}
	if addr.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "address does not belong to user"})
		return
	}

	resp := dto.AddressResponse{
		ID:                    addr.ID,
		//Address:               addr.Address,
		PhoneNumber:           addr.PhoneNumber,
		AdditionalPhoneNumber: addr.AdditionalPhoneNumber,
		DeliveryAddress:       addr.DeliveryAddress,
		ShippingAddress:       addr.ShippingAddress,
		State:                 addr.State,
		LGA:                   addr.LGA,
		
	}
	c.JSON(http.StatusOK, gin.H{"address": resp})
}

// UpdateAddress godoc
// @Summary Update an address
// @Description Update fields of an address (must belong to authenticated user)
// @Tags Customer
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "address id"
// @Param body body dto.UpdateAddressRequest true "Update payload"
// @Success 200 {object} dto.AddressResponse
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /customer/addresses/{id} [patch]
func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	ctx := c.Request.Context()
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	id := uint(id64)

	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not logged in"})
		return
	}
	userIDStr, ok := userIDVal.(string)
	if !ok || userIDStr == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user session"})
		return
	}
	uid64, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id"})
		return
	}
	userID := uint(uid64)

	var req dto.UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	addr, err := h.service.UpdateAddress(ctx, userID, id, req)
	if err != nil {
		if err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, gin.H{"error": "not allowed to update this address"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if addr == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "address not found"})
		return
	}

	resp := dto.AddressResponse{
		ID:                    addr.ID,
		//Address:               addr.Address,
		PhoneNumber:           addr.PhoneNumber,
		AdditionalPhoneNumber: addr.AdditionalPhoneNumber,
		DeliveryAddress:       addr.DeliveryAddress,
		ShippingAddress:       addr.ShippingAddress,
		State:                 addr.State,
		LGA:                   addr.LGA,
		
	}
	c.JSON(http.StatusOK, gin.H{"address": resp})
}

// DeleteAddress godoc
// @Summary Delete an address
// @Description Delete an address (must belong to authenticated user)
// @Tags Customer
// @Security BearerAuth
// @Produce json
// @Param id path int true "address id"
// @Success 204
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /customer/addresses/{id} [delete]
func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	ctx := c.Request.Context()
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	id := uint(id64)

	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not logged in"})
		return
	}
	userIDStr, ok := userIDVal.(string)
	if !ok || userIDStr == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user session"})
		return
	}
	uid64, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id"})
		return
	}
	userID := uint(uid64)

	if err := h.service.DeleteAddress(ctx, userID, id); err != nil {
		if err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, gin.H{"error": "not allowed to delete this address"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
