package handlers

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/services/product"
	"api-customer-merchant/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type ProductMediaHandler struct {
	mediaService *product.ProductMediaService
	logger       *zap.Logger
	validate     *validator.Validate
}

func NewProductMediaHandler(mediaService *product.ProductMediaService, logger *zap.Logger) *ProductMediaHandler {
	return &ProductMediaHandler{
		mediaService: mediaService,
		logger:       logger,
		validate:     validator.New(),
	}
}

// UploadMedia handles POST /merchant/products/:product_id/media
func (h *ProductMediaHandler) UploadMedia(c *gin.Context) {
	ctx := c.Request.Context()
	merchantIDStr, exists := c.Get("merchantID")
	if !exists {
		h.logger.Warn("Unauthorized merchant access")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	merchantID := merchantIDStr.(string)

	productID := strings.TrimSpace(c.Param("product_id"))
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	// Multipart form
	file, err := c.FormFile("file")
	if err != nil {
		h.logger.Error("No file in request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}

	// Temp file
	tmpFile, err := ioutil.TempFile(os.TempDir(), "upload-*.tmp")
	if err != nil {
		h.logger.Error("Temp file creation failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "upload failed"})
		return
	}
	defer os.Remove(tmpFile.Name()) // Cleanup

	if err := c.SaveUploadedFile(file, tmpFile.Name()); err != nil {
		h.logger.Error("File save failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "upload failed"})
		return
	}

	var req dto.MediaUploadRequest
	req.Type = c.PostForm("type")
	if err := h.validate.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Service call
	media, err := h.mediaService.UploadMedia(ctx, productID, merchantID, tmpFile.Name(), req.Type)
	if err != nil {
		h.logger.Error("Upload service failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := &dto.MediaUploadResponse{}
	if err := utils.RespMap(media, resp); err != nil {
		h.logger.Error("Mapping error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	h.logger.Info("Media uploaded", zap.String("product_id", productID), zap.String("media_id", media.ID))
	c.JSON(http.StatusCreated, resp)
}

// UpdateMedia handles PUT /merchant/products/:product_id/media/:media_id
func (h *ProductMediaHandler) UpdateMedia(c *gin.Context) {
	ctx := c.Request.Context()
	merchantIDStr, exists := c.Get("merchantID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	merchantID := merchantIDStr.(string)

	productID := strings.TrimSpace(c.Param("product_id"))
	mediaID := strings.TrimSpace(c.Param("media_id"))
	if productID == "" || mediaID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid IDs"})
		return
	}

	var req dto.MediaUpdateRequest
	if file, err := c.FormFile("file"); err == nil && file != nil {
		// Temp file for new upload
		tmpFile, err := ioutil.TempFile(os.TempDir(), "update-*.tmp")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "upload failed"})
			return
		}
		defer os.Remove(tmpFile.Name())
		if err := c.SaveUploadedFile(file, tmpFile.Name()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "upload failed"})
			return
		}
		req.File = &tmpFile.Name() // Pass temp path
	}
	req.URL = parseOptionalString(c.PostForm("url"))
	req.Type = parseOptionalString(c.PostForm("type"))

	if err := h.validate.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedMedia, err := h.mediaService.UpdateMedia(ctx, mediaID, productID, merchantID, &req)
	if err != nil {
		h.logger.Error("Update service failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := &dto.MediaUploadResponse{} // Reuse
	if err := utils.RespMap(updatedMedia, resp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteMedia handles DELETE /merchant/products/:product_id/media/:media_id
func (h *ProductMediaHandler) DeleteMedia(c *gin.Context) {
	ctx := c.Request.Context()
	merchantIDStr, exists := c.Get("merchantID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	merchantID := merchantIDStr.(string)

	productID := strings.TrimSpace(c.Param("product_id"))
	mediaID := strings.TrimSpace(c.Param("media_id"))
	if productID == "" || mediaID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid IDs"})
		return
	}

	var req dto.MediaDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.mediaService.DeleteMedia(ctx, mediaID, productID, merchantID, req.Reason)
	if err != nil {
		h.logger.Error("Delete service failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Media deleted", zap.String("media_id", mediaID))
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Helper
func parseOptionalString(s string) *string {
	if s != "" {
		return &s
	}
	return nil
}