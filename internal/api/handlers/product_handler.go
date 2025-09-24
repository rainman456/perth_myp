package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"api-customer-merchant/internal/api/dto" // Assuming this exists for VariantInput
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"api-customer-merchant/internal/services/product"
)

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

// CreateProduct handles product creation for a merchant
// CreateProduct godoc
// @Summary Create a new product
// @Description Allows a merchant to create a new product with variants and media
// @Tags Product
// @Accept json
// @Produce json
// @Param body body dto.ProductInput true "Product creation details"
// @Security BearerAuth
// @Success 201 {object} dto.ProductResponse "Product created successfully"
// @Failure 400 {object} object{error=string} "Invalid request payload or validation failed"
// @Failure 401 {object} object{error=string} "Unauthorized access"
// @Failure 500 {object} object{error=string} "Failed to create product"
// @Router /product/create [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	logger := h.logger.With(zap.String("operation", "CreateProduct"))

	// Check merchant authorization
	// merchantID, exists := c.Get("merchantID")
	// if !exists {
	// 	logger.Warn("Unauthorized access attempt")
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	// 	return
	// }
	// merchantIDStr, ok := merchantID.(string)
	// if !ok || merchantIDStr == "" {
	// 	logger.Warn("Invalid merchant ID in context")
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid merchant ID"})
	// 	return
	// }

	// Bind and validate input
	var input dto.ProductInput
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

	// Set merchant ID from context
	//input.MerchantID := merchantIDStr
	//merchantIDStr := input.MerchantID

	// Call service
	response, err := h.productService.CreateProductWithVariants(c.Request.Context(), &input)
	if err != nil {
		logger.Error("Failed to create product", zap.Error(err))
		if errors.Is(err, product.ErrInvalidProduct) || errors.Is(err, product.ErrInvalidMediaURL) || errors.Is(err, product.ErrInvalidAttributes) ||
			errors.Is(err, product.ErrInvalidProduct) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create product"})
		return
	}

	logger.Info("Product created successfully", zap.String("product_id", response.ID))
	c.JSON(http.StatusCreated, response)
}

// GetAllProducts handles fetching paginated products for the landing page
// GetAllProducts godoc
// @Summary Fetch all products with pagination
// @Description Retrieves a paginated list of products, optionally filtered by category for the landing page
// @Tags Product
// @Produce json
// @Param limit query int false "Number of products per page (default 20, max 100)" default(20)
// @Param offset query int false "Offset for pagination (default 0)" default(0)
// @Param category_id query uint false "Filter by category ID"
// @Success 200 {object} object{products=[]dto.ProductResponse,total=int64,limit=int,offset=int} "Products fetched successfully"
// @Failure 400 {object} object{error=string} "Invalid category ID or pagination parameters"
// @Failure 500 {object} object{error=string} "Failed to fetch products"
// @Router /product/products [get]
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
// @Summary Fetch a single product by ID
// @Description Retrieves detailed product information including media, variants, and inventory
// @Tags Product
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dto.ProductResponse "Product fetched successfully"
// @Failure 400 {object} object{error=string} "Product ID required"
// @Failure 404 {object} object{error=string} "Product not found"
// @Failure 500 {object} object{error=string} "Failed to fetch product"
// @Router /product/{id} [get]
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	logger := h.logger.With(zap.String("operation", "GetProductByID"))
	productID := c.Param("id")
	if productID == "" {
		logger.Error("Missing product ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "product ID required"})
		return
	}

	// Call service with preloads
	response, err := h.productService.GetProductByID(c.Request.Context(), productID, "Media", "Variants", "Variants.Inventory", "SimpleInventory")
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

// ListProductsByMerchant lists a merchant's products with pagination
// ListProductsByMerchant godoc
// @Summary List products by merchant
// @Description Retrieves a paginated list of products for a specific merchant, optionally active only
// @Tags Product
// @Produce json
// @Param id path string true "Merchant ID"
// @Param limit query int false "Number of products per page (default 20, max 100)" default(20)
// @Param offset query int false "Offset for pagination (default 0)" default(0)
// @Param active_only query boolean false "Filter active products only (default false)" default(false)
// @Security BearerAuth
// @Success 200 {object} object{products=[]dto.ProductResponse,total=int,limit=int,offset=int} "Products listed successfully"
// @Failure 400 {object} object{error=string} "Invalid pagination parameters"
// @Failure 401 {object} object{error=string} "Unauthorized access"
// @Failure 500 {object} object{error=string} "Failed to list products"
// @Router /product/merchant/{id} [get]
func (h *ProductHandler) ListProductsByMerchant(c *gin.Context) {
	logger := h.logger.With(zap.String("operation", "ListProductsByMerchant"))

	// Check merchant authorization
	merchantID := c.Param("id")
	if merchantID == "" {
		logger.Warn("Unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// merchantIDStr, ok := merchantID.(string)
	// if !ok || merchantIDStr == "" {
	// 	logger.Warn("Invalid merchant ID in context")
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid merchant ID"})
	// 	return
	// }

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
	products, err := h.productService.ListProductsByMerchant(c.Request.Context(), merchantID, limit, offset, activeOnly)
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

// UpdateInventory adjusts stock for a given inventory ID
// UpdateInventory godoc
// @Summary Update inventory stock
// @Description Adjusts the stock quantity for a specific inventory by delta (positive for add, negative for subtract)
// @Tags Product
// @Accept json
// @Produce json
// @Param id path string true "Inventory ID"
// @Param body body object{delta=int} true "Delta adjustment (required)"
// @Security BearerAuth
// @Success 200 {object} object{message=string} "Inventory updated successfully"
// @Failure 400 {object} object{error=string} "Inventory ID required or invalid request payload"
// @Failure 401 {object} object{error=string} "Unauthorized access"
// @Failure 404 {object} object{error=string} "Inventory not found"
// @Failure 500 {object} object{error=string} "Failed to update inventory"
// @Router /product/products/inventory [put]
func (h *ProductHandler) UpdateInventory(c *gin.Context) {
	logger := h.logger.With(zap.String("operation", "UpdateInventory"))
	inventoryID := c.Param("id")
	if inventoryID == "" {
		logger.Error("Missing inventory ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "inventory ID required"})
		return
	}

	// Check merchant authorization
	// merchantID, exists := c.Get("merchantID")
	// if !exists {
	// 	logger.Warn("Unauthorized access attempt")
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	// 	return
	// }
	// merchantIDStr, ok := merchantID.(string)
	// if !ok || merchantIDStr == "" {
	// 	logger.Warn("Invalid merchant ID in context")
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid merchant ID"})
	// 	return
	// }

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

// DeleteProduct soft-deletes a 
// DeleteProduct godoc
// @Summary Delete a product
// @Description Soft-deletes a product by ID (requires merchant authorization)
// @Tags Product
// @Produce json
// @Param id path string true "Product ID"
// @Security BearerAuth
// @Success 200 {object} object{message=string} "Product deleted successfully"
// @Failure 400 {object} object{error=string} "Product ID required"
// @Failure 401 {object} object{error=string} "Unauthorized access"
// @Failure 404 {object} object{error=string} "Product not found"
// @Failure 500 {object} object{error=string} "Failed to delete product"
// @Router /product/delete/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	logger := h.logger.With(zap.String("operation", "DeleteProduct"))
	productID := c.Param("id")
	if productID == "" {
		logger.Error("Missing product ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "product ID required"})
		return
	}

	// Check merchant authorization
	// merchantID, exists := c.Get("merchantID")
	// if !exists {
	// 	logger.Warn("Unauthorized access attempt")
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	// 	return
	// }
	// merchantIDStr, ok := merchantID.(string)
	// if !ok || merchantIDStr == "" {
	// 	logger.Warn("Invalid merchant ID in context")
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid merchant ID"})
	// 	return
	// }

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
