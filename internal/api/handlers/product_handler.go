package handlers

import (
	//"encoding/csv"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	//"fmt"
	//"io"
	"net/http"
	"strconv"

	"api-customer-merchant/internal/api/dto" // Assuming this exists for VariantInput
	"api-customer-merchant/internal/db/repositories"
	"api-customer-merchant/internal/utils"

	//"api-customer-merchant/internal/db/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"api-customer-merchant/internal/services/product"
)

type CategoryHandler struct {
	service *product.CategoryService
}

func NewCategoryHandler(service *product.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

type ProductHandler struct {
	productService *product.ProductService
	logger         *zap.Logger
	validator      *validator.Validate
}

func NewProductHandlers(productService *product.ProductService, logger *zap.Logger) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		logger:         logger,
		validator:      validator.New(),
	}
}

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
























// GetAllProducts handles fetching paginated products for the landing page
// GetAllProducts godoc
// @Summary Get all products
// @Description Fetches paginated list of products, optionally filtered by category
// @Tags Products
// @Produce json
// @Param limit query int false "Limit (default 20)"
// @Param offset query int false "Offset (default 0)"
// @Param category_id query int false "Category ID"
// @Success 200 {object} object{products=[]dto.ProductResponse,total=int64,limit=int,offset=int}
// @Failure 400 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /products [get]
func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	logger := h.logger.With(zap.String("operation", "GetAllProducts"))

	// Parse query parameters
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit <= 0 || limit > 100 {
		limit = 20
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}
	var categoryID *uint
	if catIDStr := c.Query("category_id"); catIDStr != "" {
		catID, err := strconv.ParseUint(catIDStr, 10, 32)
		if err == nil {
			id := uint(catID)
			categoryID = &id
		} else {
			logger.Error("Invalid category ID", zap.String("category_id", catIDStr))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
			return
		}
	}

	// Call service
	products, total, err := h.productService.GetAllProducts(c.Request.Context(), limit, offset, categoryID)
	if err != nil {
		logger.Error("Failed to fetch products", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch products"})
		return
	}

	logger.Info("Products fetched successfully", zap.Int("count", len(products)), zap.Int64("total", total))
	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
	})
}

// GetProductByID fetches a single product by ID
// GetProductByID godoc
// @Summary Get product by ID
// @Description Fetches a single product with media and variants
// @Tags Products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /products/{id} [get]
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	logger := h.logger.With(zap.String("operation", "GetProductByID"))
	productID := c.Param("id")
	if productID == "" {
		logger.Error("Missing product ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "product ID required"})
		return
	}

	// Call service with preloads
	response, err := h.productService.GetProductByID(c.Request.Context(), productID)
	if err != nil {
		logger.Error("Failed to fetch product", zap.Error(err), zap.String("product_id", productID))
		if errors.Is(err, product.ErrInvalidProduct) {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch product"})
		return
	}

	logger.Info("Product fetched successfully", zap.String("product_id", productID))
	c.JSON(http.StatusOK, response)
}

// GetProductByName fetches a single product by name
// GetProductByID godoc
// @Summary Get product by Name
// @Description Fetches a single product with media and variants
// @Tags Products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /products/by-name/{name} [get]
func (h *ProductHandler) GetProductByName(c *gin.Context) {
	logger := h.logger.With(zap.String("operation", "GetProductByName"))
	productName := c.Param("name")
	if productName == "" {
		logger.Error("Missing product ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "product ID required"})
		return
	}

	// Call service with preloads
	response, err := h.productService.GetProductByName(c.Request.Context(), productName)
	if err != nil {
		logger.Error("Failed to fetch product", zap.Error(err), zap.String("product_name", productName))
		if errors.Is(err, product.ErrInvalidProduct) {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch product"})
		return
	}

	logger.Info("Product fetched successfully", zap.String("product_name", productName))
	c.JSON(http.StatusOK, response)
}

// BulkUploadProducts godoc
// @Summary Bulk upload products
// @Description Upload multiple products at once for authenticated merchant
// @Tags Merchant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body []dto.ProductInput true "Array of product details"
// @Success 201 {object} object{message=string,created_count=int,errors=[]string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /merchant/products/bulk-upload [post]
func (h *ProductHandler) BulkUploadProducts(c *gin.Context) {
	logger := h.logger.With(zap.String("operation", "BulkUploadProducts"))

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
	var inputs []dto.ProductInput
	if err := c.ShouldBindJSON(&inputs); err != nil {
		logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	if len(inputs) == 0 {
		logger.Error("No products provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "no products provided"})
		return
	}

	if len(inputs) > 100 {
		logger.Error("Too many products in bulk upload", zap.Int("count", len(inputs)))
		c.JSON(http.StatusBadRequest, gin.H{"error": "maximum 100 products allowed per bulk upload"})
		return
	}

	// Process each product
	createdCount := 0
	errorMessages := []string{}

	for i, input := range inputs {
		// Validate each input
		if err := h.validator.Struct(&input); err != nil {
			errorMessages = append(errorMessages, fmt.Sprintf("Product %d validation error: %v", i+1, err.Error()))
			continue
		}

		// Call service to create product
		_, err := h.productService.CreateProductWithVariants(c.Request.Context(), merchantIDStr, &input)
		if err != nil {
			logger.Error("Failed to create product in bulk upload", zap.Error(err), zap.Int("index", i))
			errorMessages = append(errorMessages, fmt.Sprintf("Product %d creation failed: %v", i+1, err.Error()))
			continue
		}

		createdCount++
	}

	logger.Info("Bulk upload completed", zap.Int("created_count", createdCount), zap.Int("total", len(inputs)))

	if createdCount == 0 && len(errorMessages) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "all products failed to upload",
			"created_count": createdCount,
			"errors":        errorMessages,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":       "bulk upload completed",
		"created_count": createdCount,
		"errors":        errorMessages,
	})
}

// AutocompleteHandler godoc
// @Summary      Product Autocomplete
// @Description  Get product suggestions based on a name prefix for search autocomplete.
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        query     query     string  true  "Search prefix (e.g., 'a' for products starting with 'a')"
// @Param        limit     query     int     false "Number of results (default 10, max 20)"
// @Success      200  {object}  dto.AutocompleteResponse
// @Failure      400  {object}  map[string]string  "Invalid query parameter"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Router       /products/autocomplete [get]
func (h *ProductHandler) AutocompleteHandler(c *gin.Context) {
	prefix := c.Query("query")
	if prefix == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'query' is required"})
		return
	}

	limitStr := c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10 // Default
	}

	response, err := h.productService.Autocomplete(c.Request.Context(), prefix, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
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

// UpdateProduct handles product update for a merchant
// UpdateProduct godoc
// @Summary Update a product
// @Description Updates a product with variants and media for authenticated merchant
// @Tags Merchant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Param body body dto.UpdateProductInput true "Product update details"
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

	// Bind and validate input
	var input dto.UpdateProductInput
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
	response, err := h.productService.UpdateProduct(c.Request.Context(), productID, merchantIDStr, &input)
	if err != nil {
		logger.Error("Failed to update product", zap.Error(err))
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

// GetCategories godoc
// @Summary Get all categories
// @Description Retrieves all categories with parent information
// @Tags Categories
// @Produce json
// @Success 200 {array} dto.CategoryResponse "List of categories"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /categories [get]
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}

// BulkUpdateProducts godoc
// @Summary Bulk update products
// @Description Update multiple products and their variants at once for authenticated merchant
// @Tags Merchant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body []dto.BulkUpdateProductInput true "Array of product updates"
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

	// Bind and validate input
	var inputs []dto.BulkUpdateProductInput
	if err := c.ShouldBindJSON(&inputs); err != nil {
		logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
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

	// Process each product update
	updatedCount, errorMessages, err := h.productService.BulkUpdateProducts(c.Request.Context(), merchantIDStr, inputs)
	if err != nil {
		logger.Error("Failed to bulk update products", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to bulk update products"})
		return
	}

	logger.Info("Bulk update completed", zap.Int("updated_count", updatedCount), zap.Int("total", len(inputs)))

	if updatedCount == 0 && len(errorMessages) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":        "all products failed to update",
			"updated_count": updatedCount,
			"errors":       errorMessages,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "bulk update completed",
		"updated_count": updatedCount,
		"errors":       errorMessages,
	})
}

// GetAllProductsWithCategorySlug handles fetching paginated products for the landing page
// GetAllProductsWithCategorySlug godoc
// @Summary Get all products using category slug
// @Description Fetches paginated list of products,  filtered by category slug
// @Tags Categories
// @Produce json
// @Param limit query int false "Limit (default 20)"
// @Param offset query int false "Offset (default 0)"
// @Param slug path string true "Category Slug"
// @Success 200 {object} object{products=[]dto.ProductResponse,total=int64,limit=int,offset=int}
// @Failure 400 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /categories/{slug} [get]
func (h *CategoryHandler) GetAllProductsWithCategorySlug(c *gin.Context) {
	//logger := h.logger.With(zap.String("operation", "GetAllProductsWithCategorySlug"))

	// Parse query parameters
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit <= 0 || limit > 100 {
		limit = 20
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}
	//var categorySlug *string
	categorySlug := c.Param("slug")
	if categorySlug == "" {
		//logger.Error("Missing category slug")
		c.JSON(http.StatusBadRequest, gin.H{"error": "category slug required"})
		return
	}

	// Call service
	products, total, err := h.service.GetAllProductsWithCategorySlug(c.Request.Context(), limit, offset, categorySlug)
	if err != nil {
		//logger.Error("Failed to fetch products", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch products"})
		return
	}

	//logger.Info("Products fetched successfully", zap.Int("count", len(products)), zap.Int64("total", total))
	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
	})
}

// FilterProducts godoc
// @Summary Filter products with advanced options
// @Description Filter and search products by multiple criteria including price, category, attributes, etc.
// @Tags Products
// @Produce json
// @Param category_id query int false "Category ID"
// @Param category_name query string false "Category Name"
// @Param category_slug query string false "Category Slug"
// @Param min_price query number false "Minimum Price"
// @Param max_price query number false "Maximum Price"
// @Param in_stock query bool false "In Stock Only"
// @Param search query string false "Search Term"
// @Param color query string false "Color Filter"
// @Param size query string false "Size Filter"
// @Param material query string false "Material Filter"
// @Param pattern query string false "Pattern Filter"
// @Param on_sale query bool false "On Sale Only"
// @Param sort_by query string false "Sort By" Enums(price, price_desc, name, name_desc, newest, oldest, rating)
// @Param page query int false "Page Number" default(1)
// @Param limit query int false "Items Per Page" default(20)
// @Success 200 {object} object{products=[]dto.ProductResponse,total=int64,page=int,limit=int}
// @Failure 400 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /products/filter [get]
func (h *ProductHandler) FilterProducts(c *gin.Context) {
	logger := h.logger.With(zap.String("operation", "FilterProducts"))

	// Bind query parameters
	var req dto.ProductFilterRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error("Failed to bind query params", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate cache key
	cacheKey := fmt.Sprintf("product:filter:%s:p%d:l%d", req.Hash(), req.Page, req.Limit)

	type CachedResult struct {
		Products []dto.ProductResponse `json:"products"`
		Total    int64                 `json:"total"`
	}

	result, err := utils.GetOrSetCacheJSON(c.Request.Context(), cacheKey, 2*time.Minute, func() (*CachedResult, error) {
		logger.Debug("Cache miss - filtering products from DB")

		// Build filter
		var sortBy string
		if req.SortBy != nil && *req.SortBy != "" {
			sortBy = *req.SortBy
		} else {
			// choose a default that matches your repo behavior;
			// "newest" keeps the existing default ordering of created_at DESC
			sortBy = "newest"
		}

		filter := repositories.ProductFilter{
			CategoryID:   req.CategoryID,
			CategoryName: req.CategoryName,
			CategorySlug: req.CategorySlug,
			MerchantName: req.MerchantName,
			SearchTerm:   req.SearchTerm,
			InStock:      req.InStock,
			OnSale:       req.OnSale,
			Color:        req.Color,
			Size:         req.Size,
			Material:     req.Material,
			Pattern:      req.Pattern,
			SortBy:       sortBy,
		}
		if req.MinPrice != nil {
			minPrice := decimal.NewFromFloat(*req.MinPrice)
			filter.MinPrice = &minPrice
		}

		if req.MaxPrice != nil {
			maxPrice := decimal.NewFromFloat(*req.MaxPrice)
			filter.MaxPrice = &maxPrice
		}

		// Fetch products
		products, total, err := h.productService.FilterProducts(
			c.Request.Context(),
			filter,
			req.GetLimit(),
			req.GetOffset(),
		)

		if err != nil {
			return nil, err
		}

		return &CachedResult{
			Products: products,
			Total:    total,
		}, nil
	})

	if err != nil {
		logger.Error("Failed to filter products", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to filter products"})
		return
	}

	logger.Info("Products filtered successfully",
		zap.Int("count", len(result.Products)),
		zap.Int64("total", result.Total))

	c.JSON(http.StatusOK, gin.H{
		"products": result.Products,
		"total":    result.Total,
		"page":     req.Page,
		"limit":    req.Limit,
	})
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
