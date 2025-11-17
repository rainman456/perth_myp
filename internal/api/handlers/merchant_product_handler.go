package handlers

import (
	//"encoding/csv"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
	//"time"

	//"fmt"
	//"io"
	"net/http"
	"strconv"

	"api-customer-merchant/internal/api/dto" // Assuming this exists for VariantInput
	//"api-customer-merchant/internal/db/repositories"
	//"api-customer-merchant/internal/utils"

	//"api-customer-merchant/internal/db/models"

	"github.com/gin-gonic/gin"
	//"github.com/go-playground/validator/v10"
	//"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"api-customer-merchant/internal/services/product"
)


// // CreateProduct handles product creation for a merchant
// // CreateProduct godoc
// // @Summary Create a new product
// // @Description Creates a product with variants and media for authenticated merchant
// // @Tags Merchant
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param body body dto.ProductInput true "Product details"
// // @Success 201 {object} dto.MerchantProductResponse
// // @Failure 400 {object} object{error=string}
// // @Failure 401 {object} object{error=string}
// // @Router /merchant/products [post]
// func (h *ProductHandler) CreateProduct(c *gin.Context) {
// 	logger := h.logger.With(zap.String("operation", "CreateProduct"))

// 	// Check merchant authorization
// 	merchantID, exists := c.Get("merchantID")
// 	if !exists {
// 		logger.Warn("Unauthorized access attempt")
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 		return
// 	}
// 	merchantIDStr, ok := merchantID.(string)
// 	if !ok || merchantIDStr == "" {
// 		logger.Warn("Invalid merchant ID in context")
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid merchant ID"})
// 		return
// 	}

// 	// Bind and validate input
// 	var input dto.ProductInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		logger.Error("Failed to bind JSON", zap.Error(err))
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
// 		return
// 	}
// 	if err := h.validator.Struct(&input); err != nil {
// 		logger.Error("Input validation failed", zap.Error(err))
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Set merchant ID from context
// 	//input.MerchantID = merchantIDStr
// 	//merchantIDStr = input.MerchantID

// 	// Call service
// 	response, err := h.productService.CreateProductWithVariants(c.Request.Context(), merchantIDStr, &input)
// 	if err != nil {
// 		logger.Error("Failed to create product", zap.Error(err))
// 		if errors.Is(err, product.ErrInvalidProduct) || errors.Is(err, product.ErrInvalidMediaURL) || errors.Is(err, product.ErrInvalidAttributes) ||
// 			errors.Is(err, product.ErrInvalidProduct) {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create product"})
// 		return
// 	}

// 	logger.Info("Product created successfully", zap.String("product_id", response.ID))
// 	c.JSON(http.StatusCreated, response)
// }










































// CreateProduct handles product creation for a merchant with image uploads
// CreateProduct godoc
// @Summary Create a new product with images
// @Description Creates a product with variants and uploads images in a single request
// @Tags Merchant
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param name formData string true "Product name"
// @Param description formData string false "Product description"
// @Param base_price formData number true "Base price"
// @Param category_id formData int true "Category ID"
// @Param category_name formData string true "Category name"
// @Param initial_stock formData int false "Initial stock (for simple products)"
// @Param discount formData number false "Discount amount"
// @Param discount_type formData string false "Discount type (fixed/percentage)"
// @Param variants formData string false "JSON array of variants"
// @Param images formData file false "Product images (multiple files allowed)"
// @Success 201 {object} dto.MerchantProductResponse
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Router /merchant/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	logger := h.logger.With(zap.String("operation", "CreateProduct"))

	// Check merchant authorization
	merchantID, exists := c.Get("merchantID")
	if !exists {
		logger.Warn("Unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	merchantIDStr, ok := merchantID.(string)
	if !ok || merchantIDStr == "" {
		logger.Warn("Invalid merchant ID in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid merchant ID"})
		return
	}

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 32 MB max
		logger.Error("Failed to parse multipart form", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form data"})
		return
	}

	// Build ProductInput from form data
	var input dto.ProductInput
	
	// Parse basic fields
	input.Name = strings.TrimSpace(c.PostForm("name"))
	input.Description = strings.TrimSpace(c.PostForm("description"))
	input.CategoryName = strings.TrimSpace(c.PostForm("category_name"))
	input.DiscountType = strings.TrimSpace(c.PostForm("discount_type"))

	// Parse numeric fields
	basePrice, err := strconv.ParseFloat(c.PostForm("base_price"), 64)
	if err != nil {
		logger.Error("Invalid base_price", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid base_price"})
		return
	}
	input.BasePrice = basePrice

	categoryID, err := strconv.ParseUint(c.PostForm("category_id"), 10, 32)
	if err != nil {
		logger.Error("Invalid category_id", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category_id"})
		return
	}
	input.CategoryID = uint(categoryID)

	// Parse optional discount
	if discountStr := c.PostForm("discount"); discountStr != "" {
		discount, err := strconv.ParseFloat(discountStr, 64)
		if err != nil {
			logger.Error("Invalid discount", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid discount"})
			return
		}
		input.Discount = discount
	}

	// Parse optional initial_stock (for simple products)
	if stockStr := c.PostForm("initial_stock"); stockStr != "" {
		stock, err := strconv.Atoi(stockStr)
		if err != nil {
			logger.Error("Invalid initial_stock", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid initial_stock"})
			return
		}
		input.InitialStock = &stock
	}

	// Parse variants JSON if provided
	if variantsJSON := c.PostForm("variants"); variantsJSON != "" {
		var variants []dto.VariantInput
		if err := json.Unmarshal([]byte(variantsJSON), &variants); err != nil {
			logger.Error("Invalid variants JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid variants format"})
			return
		}
		input.Variants = variants
	}

	// Validate input
	if err := h.validator.Struct(&input); err != nil {
		logger.Error("Input validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Handle file uploads
	uploadedMediaURLs := []dto.MediaInput{}
	uploadedPublicIDs := []string{} // Track for cleanup on failure
	
	form, err := c.MultipartForm()
	if err == nil && form.File["images"] != nil {
		files := form.File["images"]
		
		for i, fileHeader := range files {
			logger.Info("Processing image upload", zap.Int("index", i), zap.String("filename", fileHeader.Filename))
			
			// Validate file type
			contentType := fileHeader.Header.Get("Content-Type")
			if !strings.HasPrefix(contentType, "image/") {
				logger.Error("Invalid file type", zap.String("content_type", contentType))
				h.cleanupCloudinaryUploads(c.Request.Context(), uploadedPublicIDs)
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("file %d is not an image", i+1)})
				return
			}
			
			// Create temp file
			tmpFile, err := os.CreateTemp(os.TempDir(), "product-*.tmp")
			if err != nil {
				logger.Error("Failed to create temp file", zap.Error(err))
				h.cleanupCloudinaryUploads(c.Request.Context(), uploadedPublicIDs)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process image upload"})
				return
			}
			tmpPath := tmpFile.Name()
			tmpFile.Close()
			
			// Save uploaded file to temp location
			if err := c.SaveUploadedFile(fileHeader, tmpPath); err != nil {
				logger.Error("Failed to save uploaded file", zap.Error(err))
				os.Remove(tmpPath)
				h.cleanupCloudinaryUploads(c.Request.Context(), uploadedPublicIDs)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save uploaded file"})
				return
			}
			
			// Upload to Cloudinary
			cloudinaryURL, publicID, err := h.productService.UploadToCloudinary(c.Request.Context(), tmpPath, "image")
			os.Remove(tmpPath) // Clean up temp file immediately
			
			if err != nil {
				logger.Error("Cloudinary upload failed", zap.Error(err), zap.String("filename", fileHeader.Filename))
				h.cleanupCloudinaryUploads(c.Request.Context(), uploadedPublicIDs)
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to upload image %d", i+1)})
				return
			}
			
			uploadedPublicIDs = append(uploadedPublicIDs, publicID)
			uploadedMediaURLs = append(uploadedMediaURLs, dto.MediaInput{
				URL:  cloudinaryURL,
				Type: "image",
			})
			
			logger.Info("Image uploaded successfully", zap.String("public_id", publicID))
		}
	}

	// Add uploaded images to input
	input.Images = uploadedMediaURLs

	// Call service to create product
	response, err := h.productService.CreateProductWithVariants(c.Request.Context(), merchantIDStr, &input)
	if err != nil {
		logger.Error("Failed to create product", zap.Error(err))
		// Clean up uploaded images if product creation fails
		h.cleanupCloudinaryUploads(c.Request.Context(), uploadedPublicIDs)
		
		if errors.Is(err, product.ErrInvalidProduct) || errors.Is(err, product.ErrInvalidMediaURL) || errors.Is(err, product.ErrInvalidAttributes) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create product"})
		return
	}

	logger.Info("Product created successfully with images", zap.String("product_id", response.ID), zap.Int("image_count", len(uploadedMediaURLs)))
	c.JSON(http.StatusCreated, response)
}

// cleanupCloudinaryUploads removes uploaded images from Cloudinary if operation fails
func (h *ProductHandler) cleanupCloudinaryUploads(ctx context.Context, publicIDs []string) {
	if len(publicIDs) == 0 {
		return
	}
	
	h.logger.Info("Cleaning up Cloudinary uploads", zap.Int("count", len(publicIDs)))
	for _, publicID := range publicIDs {
		if err := h.productService.DeleteFromCloudinary(ctx, publicID); err != nil {
			h.logger.Error("Failed to cleanup Cloudinary upload", zap.String("public_id", publicID), zap.Error(err))
		}
	}
}


































// BulkUpdateInventory godoc
// @Summary Bulk update inventory
// @Description Update multiple inventory items at once for authenticated merchant
// @Tags Merchant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body []dto.BulkInventoryUpdateInput true "Array of inventory updates"
// @Success 200 {object} object{message=string,updated_count=int,errors=[]string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /merchant/products/bulk-inventory-update [put]
func (h *ProductHandler) BulkUpdateInventory(c *gin.Context) {
	logger := h.logger.With(zap.String("operation", "BulkUpdateInventory"))

	// Check merchant authorization
	merchantID, exists := c.Get("merchantID")
	if !exists {
		logger.Warn("Unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	merchantIDStr, ok := merchantID.(string)
	if !ok || merchantIDStr == "" {
		logger.Warn("Invalid merchant ID in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid merchant ID"})
		return
	}

	// Bind and validate input
	var inputs []dto.BulkInventoryUpdateInput
	if err := c.ShouldBindJSON(&inputs); err != nil {
		logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	if len(inputs) == 0 {
		logger.Error("No inventory items provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "no inventory items provided"})
		return
	}

	if len(inputs) > 1000 {
		logger.Error("Too many inventory items in bulk update", zap.Int("count", len(inputs)))
		c.JSON(http.StatusBadRequest, gin.H{"error": "maximum 1000 inventory items allowed per bulk update"})
		return
	}

	// Process each inventory update
	updatedCount, errorMessages, err := h.productService.BulkUpdateInventory(c.Request.Context(), merchantIDStr, inputs)
	if err != nil {
		logger.Error("Failed to bulk update inventory", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to bulk update inventory"})
		return
	}

	logger.Info("Bulk inventory update completed", zap.Int("updated_count", updatedCount), zap.Int("total", len(inputs)))

	if updatedCount == 0 && len(errorMessages) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "all inventory items failed to update",
			"updated_count": updatedCount,
			"errors":        errorMessages,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "bulk inventory update completed",
		"updated_count": updatedCount,
		"errors":        errorMessages,
	})
}











// ListProductsByMerchant lists a merchant's products with pagination
// ListProductsByMerchant godoc
// @Summary List merchant's products
// @Description Fetches paginated list of products for authenticated merchant
// @Tags Merchant
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit (default 20)"
// @Param offset query int false "Offset (default 0)"
// @Param active_only query boolean false "Active only (default false)"
// @Success 200 {object} object{products=[]dto.MerchantProductResponse,total=int,limit=int,offset=int}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /merchant/products [get]
func (h *ProductHandler) ListProductsByMerchant(c *gin.Context) {
	logger := h.logger.With(zap.String("operation", "ListProductsByMerchant"))

	// Check merchant authorization
	merchantID, exists := c.Get("merchantID")
	if !exists {
		logger.Warn("Unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	merchantIDStr, ok := merchantID.(string)
	if !ok || merchantIDStr == "" {
		logger.Warn("Invalid merchant ID in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid merchant ID"})
		return
	}

	// Parse query parameters
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit <= 0 || limit > 100 {
		limit = 20
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}
	activeOnly := c.Query("active_only") == "true"

	// Call service
	products, err := h.productService.ListProductsByMerchant(c.Request.Context(), merchantIDStr, limit, offset, activeOnly)
	if err != nil {
		logger.Error("Failed to list products", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list products"})
		return
	}

	logger.Info("Products listed successfully", zap.Int("count", len(products)))
	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"total":    len(products), // Note: Repository doesn't return total; add if needed
		"limit":    limit,
		"offset":   offset,
	})
}










func (h *ProductHandler) processImageUploads(ctx context.Context, files []*multipart.FileHeader, logger *zap.Logger) ([]dto.MediaInput, []string, error) {
	uploadedMediaURLs := []dto.MediaInput{}
	uploadedPublicIDs := []string{}

	for i, fileHeader := range files {
		// Validate file type
		contentType := fileHeader.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "image/") {
			h.cleanupCloudinaryUploads(ctx, uploadedPublicIDs)
			return nil, nil, fmt.Errorf("file %d is not an image", i+1)
		}

		// Create temp file
		tmpFile, err := os.CreateTemp(os.TempDir(), "product-*.tmp")
		if err != nil {
			h.cleanupCloudinaryUploads(ctx, uploadedPublicIDs)
			return nil, nil, fmt.Errorf("failed to create temp file: %w", err)
		}
		tmpPath := tmpFile.Name()
		tmpFile.Close()

		// Save uploaded file
		file, err := fileHeader.Open()
		if err != nil {
			os.Remove(tmpPath)
			h.cleanupCloudinaryUploads(ctx, uploadedPublicIDs)
			return nil, nil, fmt.Errorf("failed to open uploaded file: %w", err)
		}

		tmpFile, err = os.OpenFile(tmpPath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			file.Close()
			os.Remove(tmpPath)
			h.cleanupCloudinaryUploads(ctx, uploadedPublicIDs)
			return nil, nil, fmt.Errorf("failed to open temp file: %w", err)
		}

		_, err = io.Copy(tmpFile, file)
		file.Close()
		tmpFile.Close()

		if err != nil {
			os.Remove(tmpPath)
			h.cleanupCloudinaryUploads(ctx, uploadedPublicIDs)
			return nil, nil, fmt.Errorf("failed to save uploaded file: %w", err)
		}

		// Upload to Cloudinary
		cloudinaryURL, publicID, err := h.productService.UploadToCloudinary(ctx, tmpPath, "image")
		os.Remove(tmpPath) // Clean up temp file

		if err != nil {
			h.cleanupCloudinaryUploads(ctx, uploadedPublicIDs)
			return nil, nil, fmt.Errorf("cloudinary upload failed for image %d: %w", i+1, err)
		}

		uploadedPublicIDs = append(uploadedPublicIDs, publicID)
		uploadedMediaURLs = append(uploadedMediaURLs, dto.MediaInput{
			URL:  cloudinaryURL,
			Type: "image",
		})

		logger.Info("Image uploaded successfully", zap.String("public_id", publicID))
	}

	return uploadedMediaURLs, uploadedPublicIDs, nil
}





// // BulkUpdateProducts godoc
// // @Summary Bulk update products
// // @Description Update multiple products and their variants at once for authenticated merchant
// // @Tags Merchant
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param body body []dto.BulkUpdateProductInput true "Array of product updates"
// // @Success 200 {object} object{message=string,updated_count=int,errors=[]string}
// // @Failure 400 {object} object{error=string}
// // @Failure 401 {object} object{error=string}
// // @Failure 500 {object} object{error=string}
// // @Router /merchant/products/bulk-update [put]
// func (h *ProductHandler) BulkUpdateProducts(c *gin.Context) {
// 	logger := h.logger.With(zap.String("operation", "BulkUpdateProducts"))

// 	// Check merchant authorization
// 	merchantID, exists := c.Get("merchantID")
// 	if !exists {
// 		logger.Warn("Unauthorized access attempt")
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 		return
// 	}
// 	merchantIDStr, ok := merchantID.(string)
// 	if !ok || merchantIDStr == "" {
// 		logger.Warn("Invalid merchant ID in context")
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid merchant ID"})
// 		return
// 	}

// 	// Bind and validate input
// 	var inputs []dto.BulkUpdateProductInput
// 	if err := c.ShouldBindJSON(&inputs); err != nil {
// 		logger.Error("Failed to bind JSON", zap.Error(err))
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
// 		return
// 	}

// 	if len(inputs) == 0 {
// 		logger.Error("No products provided")
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "no products provided"})
// 		return
// 	}

// 	if len(inputs) > 100 {
// 		logger.Error("Too many products in bulk update", zap.Int("count", len(inputs)))
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "maximum 100 products allowed per bulk update"})
// 		return
// 	}

// 	// Process each product update
// 	updatedCount, errorMessages, err := h.productService.BulkUpdateProducts(c.Request.Context(), merchantIDStr, inputs)
// 	if err != nil {
// 		logger.Error("Failed to bulk update products", zap.Error(err))
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to bulk update products"})
// 		return
// 	}

// 	logger.Info("Bulk update completed", zap.Int("updated_count", updatedCount), zap.Int("total", len(inputs)))

// 	if updatedCount == 0 && len(errorMessages) > 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error":        "all products failed to update",
// 			"updated_count": updatedCount,
// 			"errors":       errorMessages,
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message":      "bulk update completed",
// 		"updated_count": updatedCount,
// 		"errors":       errorMessages,
// 	})
// }







// BulkUpdateProducts godoc
// @Summary Bulk update products with images
// @Description Update multiple products and their variants at once with optional new images
// @Tags Merchant
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param updates formData string true "JSON array of product updates"
// @Param images_0 formData file false "New images for product update 0"
// @Param images_1 formData file false "New images for product update 1"
// @Param images_N formData file false "New images for product update N (pattern: images_{index})"
// @Success 200 {object} object{message=string,updated_count=int,errors=[]string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /merchant/products/bulk-update [put]
func (h *ProductHandler) BulkUpdateProducts(c *gin.Context) {
	logger := h.logger.With(zap.String("operation", "BulkUpdateProducts"))

	// Check merchant authorization
	merchantID, exists := c.Get("merchantID")
	if !exists {
		logger.Warn("Unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	merchantIDStr, ok := merchantID.(string)
	if !ok || merchantIDStr == "" {
		logger.Warn("Invalid merchant ID in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid merchant ID"})
		return
	}

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(100 << 20); err != nil {
		logger.Error("Failed to parse multipart form", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form data"})
		return
	}

	// Parse updates JSON array
	updatesJSON := c.PostForm("updates")
	if updatesJSON == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "updates data required"})
		return
	}

	var inputs []dto.BulkUpdateProductInput
	if err := json.Unmarshal([]byte(updatesJSON), &inputs); err != nil {
		logger.Error("Failed to parse updates JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid updates format"})
		return
	}

	if len(inputs) == 0 {
		logger.Error("No products provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "no products provided"})
		return
	}

	if len(inputs) > 100 {
		logger.Error("Too many products in bulk update", zap.Int("count", len(inputs)))
		c.JSON(http.StatusBadRequest, gin.H{"error": "maximum 100 products allowed per bulk update"})
		return
	}

	// Get form with files
	form, _ := c.MultipartForm()

	// Process each product update
	updatedCount := 0
	errorMessages := []string{}
	allUploadedPublicIDs := []string{}

	for i, input := range inputs {
		// Validate product ID
		if input.ProductID == "" {
			errorMessages = append(errorMessages, fmt.Sprintf("Product %d: product_id is required", i+1))
			continue
		}

		// Process images for this product (field name: images_0, images_1, etc.)
		imageFieldName := fmt.Sprintf("images_%d", i)
		uploadedPublicIDs := []string{}

		if files, exists := form.File[imageFieldName]; exists && len(files) > 0 {
			newImages, publicIDs, err := h.processImageUploads(c.Request.Context(), files, logger)
			if err != nil {
				errorMessages = append(errorMessages, fmt.Sprintf("Product %d (%s) image upload failed: %v", i+1, input.ProductID, err.Error()))
				continue
			}
			uploadedPublicIDs = publicIDs
			allUploadedPublicIDs = append(allUploadedPublicIDs, uploadedPublicIDs...)

			// Add new images to product
			if err := h.productService.AddProductMedia(c.Request.Context(), input.ProductID, merchantIDStr, newImages); err != nil {
				h.cleanupCloudinaryUploads(c.Request.Context(), uploadedPublicIDs)
				errorMessages = append(errorMessages, fmt.Sprintf("Product %d (%s) failed to add images: %v", i+1, input.ProductID, err.Error()))
				continue
			}
		}

		// Update product if provided
		if input.Product != nil {
			_, err := h.productService.UpdateProduct(c.Request.Context(), input.ProductID, merchantIDStr, input.Product)
			if err != nil {
				logger.Error("Failed to update product in bulk update", zap.Error(err), zap.String("product_id", input.ProductID))
				h.cleanupCloudinaryUploads(c.Request.Context(), uploadedPublicIDs)
				errorMessages = append(errorMessages, fmt.Sprintf("Product %d (%s): %v", i+1, input.ProductID, err.Error()))
				continue
			}
		}

		// Update variants if provided
		for j, variantUpdate := range input.Variants {
			if variantUpdate.VariantID == "" {
				errorMessages = append(errorMessages, fmt.Sprintf("Product %d, Variant %d: variant_id is required", i+1, j+1))
				continue
			}

			if variantUpdate.Variant != nil {
				err := h.productService.UpdateVariant(c.Request.Context(), variantUpdate.VariantID, merchantIDStr, variantUpdate.Variant)
				if err != nil {
					logger.Error("Failed to update variant in bulk update", zap.Error(err), zap.String("variant_id", variantUpdate.VariantID))
					errorMessages = append(errorMessages, fmt.Sprintf("Product %d (%s), Variant %d (%s): %v", i+1, input.ProductID, j+1, variantUpdate.VariantID, err.Error()))
					continue
				}
			}
		}

		updatedCount++
	}

	logger.Info("Bulk update completed", zap.Int("updated_count", updatedCount), zap.Int("total", len(inputs)))

	if updatedCount == 0 && len(errorMessages) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "all products failed to update",
			"updated_count": updatedCount,
			"errors":        errorMessages,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "bulk update completed",
		"updated_count": updatedCount,
		"errors":        errorMessages,
	})
}






// DeleteProduct soft-deletes a product
// DeleteProduct godoc
// @Summary Delete a product
// @Description Soft-deletes a product by ID for authenticated merchant
// @Tags Merchant
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 200 {object} object{message=string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /merchant/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	logger := h.logger.With(zap.String("operation", "DeleteProduct"))
	productID := c.Param("id")
	if productID == "" {
		logger.Error("Missing product ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "product ID required"})
		return
	}

	// Check merchant authorization
	merchantID, exists := c.Get("merchantID")
	if !exists {
		logger.Warn("Unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	merchantIDStr, ok := merchantID.(string)
	if !ok || merchantIDStr == "" {
		logger.Warn("Invalid merchant ID in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid merchant ID"})
		return
	}

	// Call service
	err := h.productService.DeleteProduct(c.Request.Context(), productID)
	if err != nil {
		logger.Error("Failed to delete product", zap.Error(err), zap.String("product_id", productID))
		if errors.Is(err, product.ErrInvalidProduct) {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete product"})
		return
	}

	logger.Info("Product deleted successfully", zap.String("product_id", productID))
	c.JSON(http.StatusOK, gin.H{"message": "product deleted"})
}




























// UpdateInventory adjusts stock for a given inventory ID
// UpdateInventory godoc
// @Summary Update product inventory
// @Description Adjusts stock delta for a given inventory ID
// @Tags Merchant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Inventory ID"
// @Param body body object{delta=int} true "Stock adjustment"
// @Success 200 {object} object{message=string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /merchant/products/inventory/{id} [put]
func (h *ProductHandler) UpdateInventory(c *gin.Context) {
	logger := h.logger.With(zap.String("operation", "UpdateInventory"))
	inventoryID := c.Param("id")
	if inventoryID == "" {
		logger.Error("Missing inventory ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "inventory ID required"})
		return
	}

	// Check merchant authorization
	merchantID, exists := c.Get("merchantID")
	if !exists {
		logger.Warn("Unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	merchantIDStr, ok := merchantID.(string)
	if !ok || merchantIDStr == "" {
		logger.Warn("Invalid merchant ID in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid merchant ID"})
		return
	}

	// Bind input
	var input struct {
		Delta int `json:"delta" validate:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}
	if err := h.validator.Struct(&input); err != nil {
		logger.Error("Input validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service
	err := h.productService.UpdateInventory(c.Request.Context(), inventoryID, input.Delta)
	if err != nil {
		logger.Error("Failed to update inventory", zap.Error(err), zap.String("inventory_id", inventoryID))
		if errors.Is(err, product.ErrInvalidProduct) {
			c.JSON(http.StatusNotFound, gin.H{"error": "inventory not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update inventory"})
		return
	}

	logger.Info("Inventory updated successfully", zap.String("inventory_id", inventoryID), zap.Int("delta", input.Delta))
	c.JSON(http.StatusOK, gin.H{"message": "inventory updated"})
}














// UpdateVariant handles variant update for a merchant
// UpdateVariant godoc
// @Summary Update a variant
// @Description Updates a variant for authenticated merchant
// @Tags Merchant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Variant ID"
// @Param body body dto.UpdateVariantInput true "Variant update details"
// @Success 200 {object} object{message=string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /merchant/products/variants/{id} [put]
func (h *ProductHandler) UpdateVariant(c *gin.Context) {
	logger := h.logger.With(zap.String("operation", "UpdateVariant"))
	variantID := c.Param("id")
	if variantID == "" {
		logger.Error("Missing variant ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "variant ID required"})
		return
	}

	// Check merchant authorization
	merchantID, exists := c.Get("merchantID")
	if !exists {
		logger.Warn("Unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	merchantIDStr, ok := merchantID.(string)
	if !ok || merchantIDStr == "" {
		logger.Warn("Invalid merchant ID in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid merchant ID"})
		return
	}

	// Bind and validate input
	var input dto.UpdateVariantInput
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}
	if err := h.validator.Struct(&input); err != nil {
		logger.Error("Input validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service
	err := h.productService.UpdateVariant(c.Request.Context(), variantID, merchantIDStr, &input)
	if err != nil {
		logger.Error("Failed to update variant", zap.Error(err))
		if errors.Is(err, product.ErrInvalidProduct) || errors.Is(err, product.ErrUnauthorized) {
			c.JSON(http.StatusNotFound, gin.H{"error": "variant not found or unauthorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update variant"})
		return
	}

	logger.Info("Variant updated successfully", zap.String("variant_id", variantID))
	c.JSON(http.StatusOK, gin.H{"message": "variant updated successfully"})
}

















































































// // UpdateProduct handles product update for a merchant
// // UpdateProduct godoc
// // @Summary Update a product
// // @Description Updates a product with variants and media for authenticated merchant
// // @Tags Merchant
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path string true "Product ID"
// // @Param body body dto.UpdateProductInput true "Product update details"
// // @Success 200 {object} dto.MerchantProductResponse
// // @Failure 400 {object} object{error=string}
// // @Failure 401 {object} object{error=string}
// // @Failure 404 {object} object{error=string}
// // @Router /merchant/products/{id} [put]
// func (h *ProductHandler) UpdateProduct(c *gin.Context) {
// 	logger := h.logger.With(zap.String("operation", "UpdateProduct"))
// 	productID := c.Param("id")
// 	if productID == "" {
// 		logger.Error("Missing product ID")
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "product ID required"})
// 		return
// 	}

// 	// Check merchant authorization
// 	merchantID, exists := c.Get("merchantID")
// 	if !exists {
// 		logger.Warn("Unauthorized access attempt")
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 		return
// 	}
// 	merchantIDStr, ok := merchantID.(string)
// 	if !ok || merchantIDStr == "" {
// 		logger.Warn("Invalid merchant ID in context")
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid merchant ID"})
// 		return
// 	}

// 	// Bind and validate input
// 	var input dto.UpdateProductInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		logger.Error("Failed to bind JSON", zap.Error(err))
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
// 		return
// 	}
// 	if err := h.validator.Struct(&input); err != nil {
// 		logger.Error("Input validation failed", zap.Error(err))
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Call service
// 	response, err := h.productService.UpdateProduct(c.Request.Context(), productID, merchantIDStr, &input)
// 	if err != nil {
// 		logger.Error("Failed to update product", zap.Error(err))
// 		if errors.Is(err, product.ErrInvalidProduct) || errors.Is(err, product.ErrUnauthorized) {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "product not found or unauthorized"})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update product"})
// 		return
// 	}

// 	logger.Info("Product updated successfully", zap.String("product_id", productID))
// 	c.JSON(http.StatusOK, response)
// }































// UpdateProduct handles product update for a merchant with optional new images
// UpdateProduct godoc
// @Summary Update a product
// @Description Updates a product with optional new images for authenticated merchant
// @Tags Merchant
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Param name formData string false "Product name"
// @Param description formData string false "Product description"
// @Param base_price formData number false "Base price"
// @Param category_id formData int false "Category ID"
// @Param category_name formData string false "Category name"
// @Param discount formData number false "Discount amount"
// @Param discount_type formData string false "Discount type (fixed/percentage)"
// @Param images formData file false "New product images (multiple files allowed)"
// @Success 200 {object} dto.MerchantProductResponse
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /merchant/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	logger := h.logger.With(zap.String("operation", "UpdateProduct"))
	productID := c.Param("id")
	if productID == "" {
		logger.Error("Missing product ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "product ID required"})
		return
	}

	// Check merchant authorization
	merchantID, exists := c.Get("merchantID")
	if !exists {
		logger.Warn("Unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	merchantIDStr, ok := merchantID.(string)
	if !ok || merchantIDStr == "" {
		logger.Warn("Invalid merchant ID in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid merchant ID"})
		return
	}

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		logger.Error("Failed to parse multipart form", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form data"})
		return
	}

	// Build UpdateProductInput from form data
	var input dto.UpdateProductInput

	// Parse optional fields
	if name := c.PostForm("name"); name != "" {
		trimmedName := strings.TrimSpace(name)
		input.Name = &trimmedName
	}

	if description := c.PostForm("description"); description != "" {
		trimmedDesc := strings.TrimSpace(description)
		input.Description = &trimmedDesc
	}

	if categoryName := c.PostForm("category_name"); categoryName != "" {
		trimmedCatName := strings.TrimSpace(categoryName)
		input.CategoryName = &trimmedCatName
	}

	if discountType := c.PostForm("discount_type"); discountType != "" {
		trimmedDiscType := strings.TrimSpace(discountType)
		input.DiscountType = &trimmedDiscType
	}

	if basePriceStr := c.PostForm("base_price"); basePriceStr != "" {
		basePrice, err := strconv.ParseFloat(basePriceStr, 64)
		if err != nil {
			logger.Error("Invalid base_price", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid base_price"})
			return
		}
		input.BasePrice = &basePrice
	}

	if categoryIDStr := c.PostForm("category_id"); categoryIDStr != "" {
		categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
		if err != nil {
			logger.Error("Invalid category_id", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category_id"})
			return
		}
		catID := uint(categoryID)
		input.CategoryID = &catID
	}

	if discountStr := c.PostForm("discount"); discountStr != "" {
		discount, err := strconv.ParseFloat(discountStr, 64)
		if err != nil {
			logger.Error("Invalid discount", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid discount"})
			return
		}
		input.Discount = &discount
	}

	// Validate input
	if err := h.validator.Struct(&input); err != nil {
		logger.Error("Input validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Handle optional new image uploads
	uploadedPublicIDs := []string{}
	form, err := c.MultipartForm()
	if err == nil && form.File["images"] != nil {
		files := form.File["images"]

		newImages, publicIDs, err := h.processImageUploads(c.Request.Context(), files, logger)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		uploadedPublicIDs = publicIDs

		// Add new images to product via service
		if len(newImages) > 0 {
			if err := h.productService.AddProductMedia(c.Request.Context(), productID, merchantIDStr, newImages); err != nil {
				h.cleanupCloudinaryUploads(c.Request.Context(), uploadedPublicIDs)
				logger.Error("Failed to add new images", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add new images"})
				return
			}
		}
	}

	// Call service to update product
	response, err := h.productService.UpdateProduct(c.Request.Context(), productID, merchantIDStr, &input)
	if err != nil {
		logger.Error("Failed to update product", zap.Error(err))
		h.cleanupCloudinaryUploads(c.Request.Context(), uploadedPublicIDs)
		if errors.Is(err, product.ErrInvalidProduct) || errors.Is(err, product.ErrUnauthorized) {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found or unauthorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update product"})
		return
	}

	logger.Info("Product updated successfully", zap.String("product_id", productID))
	c.JSON(http.StatusOK, response)
}
