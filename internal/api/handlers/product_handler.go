package handlers

import (
	//"encoding/csv"
	"errors"
	"fmt"
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

// CreateProduct handles product creation for a merchant
// CreateProduct godoc
// @Summary Create a new product
// @Description Creates a product with variants and media for authenticated merchant
// @Tags Merchant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.ProductInput true "Product details"
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
	//input.MerchantID = merchantIDStr
	//merchantIDStr = input.MerchantID

	// Call service
	response, err := h.productService.CreateProductWithVariants(c.Request.Context(),merchantIDStr, &input)
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