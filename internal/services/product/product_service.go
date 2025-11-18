package product

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	//"os"
	"path/filepath"

	//"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/api/helpers"
	"api-customer-merchant/internal/config"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"
	"api-customer-merchant/internal/utils"

	//"api-customer-merchant/internal/services/review"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var (
	ErrInvalidProduct = errors.New("invalid product data")
	//ErrInvalidSKU        = errors.New("invalid SKU format")
	ErrInvalidMediaURL   = errors.New("invalid media URL")
	ErrInvalidAttributes = errors.New("invalid variant attributes")
	ErrUnauthorized      = errors.New("unauthorized operation")
	ErrUploadFailed      = errors.New("upload to Cloudinary failed")
	ErrUpdateFailed      = errors.New("update failed")
	ErrDeleteFailed      = errors.New("delete failed")
	ErrUnauthorizedMedia = errors.New("unauthorized for this media")
)

// SKU validation regex: alphanumeric, hyphens, underscores, max 100 chars
//var skuRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,100}$`)

type ProductService struct {
	productRepo *repositories.ProductRepository
	//reviewRepo  *repositories.ReviewRepository
	logger    *zap.Logger
	validator *validator.Validate
	cld       *cloudinary.Cloudinary
	//config  *config.Config

}

func NewProductService(productRepo *repositories.ProductRepository, cfg *config.Config, logger *zap.Logger) *ProductService {
	cld, err := cloudinary.NewFromParams(cfg.CloudinaryCloudName, cfg.CloudinaryAPIKey, cfg.CloudinaryAPISecret)
	if err != nil {
		logger.Fatal("Cloudinary init failed", zap.Error(err))
	}

	return &ProductService{
		productRepo: productRepo,
		//reviewRepo: reviewRepo,
		logger:    logger,
		validator: validator.New(),
		cld:       cld,
	}
}

// CreateProductWithVariants creates a product from input DTO
func (s *ProductService) CreateProductWithVariants(ctx context.Context, merchant_id string, input *dto.ProductInput) (*dto.MerchantProductResponse, error) {
	logger := s.logger.With(zap.String("operation", "CreateProductWithVariants"))

	// Validate input
	if err := s.validator.Struct(input); err != nil {
		logger.Error("Input validation failed", zap.Error(err))
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	// Additional validation
	// if !skuRegex.MatchString(input.SKU) {
	// 	logger.Error("Invalid SKU format", zap.String("sku", input.SKU))
	// 	return nil, ErrInvalidSKU
	// }

	isSimple := len(input.Variants) == 0
	if isSimple && input.InitialStock == nil {
		logger.Error("Initial stock required for simple product")
		return nil, ErrInvalidProduct
	}

	// Map DTO to models
	product := &models.Product{
		Name:        strings.TrimSpace(input.Name),
		MerchantID:  merchant_id,
		Description: strings.TrimSpace(input.Description),
		//SKU:         strings.TrimSpace(input.SKU),
		BasePrice:    decimal.NewFromFloat(input.BasePrice),
		Discount:     decimal.NewFromFloat(input.Discount),
		DiscountType: models.DiscountType(input.DiscountType),
		CategoryID:   input.CategoryID,
		CategoryName: input.CategoryName,
	}
	variants := make([]models.Variant, len(input.Variants))
	for i, v := range input.Variants {
		variants[i] = models.Variant{
			//SKU:             strings.TrimSpace(v.SKU),
			PriceAdjustment: decimal.NewFromFloat(v.PriceAdjustment),
			Discount:        decimal.NewFromFloat(input.Discount),
			DiscountType:    models.DiscountType(input.DiscountType),
			Attributes:      v.Attributes,
			IsActive:        true,
		}
	}
	media := make([]models.Media, len(input.Images))
	for i, m := range input.Images {
		media[i] = models.Media{
			URL:  strings.TrimSpace(m.URL),
			Type: models.MediaType(m.Type),
		}
	}

	product.GenerateSKU(merchant_id)
	for i := range variants {
		variants[i].GenerateSKU(product.SKU)
	}

	// Delegate to repo
	var simpleStock *int
	if isSimple {
		simpleStock = input.InitialStock
	}
	err := s.productRepo.CreateProductWithVariantsAndInventory(ctx, product, variants, input.Variants, media, simpleStock, isSimple)
	if err != nil {
		if errors.Is(err, repositories.ErrDuplicateSKU) {
			return nil, fmt.Errorf("duplicate SKU: %w", err)
		}
		if errors.Is(err, repositories.ErrInvalidInventory) {
			return nil, fmt.Errorf("invalid inventory setup: %w", err)
		}
		logger.Error("Failed to create product", zap.Error(err))
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// Map to response DTO
	go utils.InvalidateCachePattern(context.Background(), "product:list:*")
	response := helpers.ToMerchantProductResponse(product)

	logger.Info("Product created successfully", zap.String("product_id", product.ID))
	return response, nil
}

// GetProductByID fetches a product with optional preloads
func (s *ProductService) GetProductByID(ctx context.Context, id string, preloads ...string) (*dto.ProductResponse, error) {
	logger := s.logger.With(zap.String("operation", "GetProductByID"), zap.String("product_id", id))

	// Try cache first
	cacheKey := utils.ProductCacheKey(id)

	response, err := utils.GetOrSetCacheJSON(ctx, cacheKey, 5*time.Minute, func() (*dto.ProductResponse, error) {
		logger.Debug("Cache miss - fetching from DB")

		// Set default preloads if none provided (exclude Reviews for list view)
		if len(preloads) == 0 {
			preloads = []string{
				"Category",
				"Merchant",
				"Variants",
				"Media",
				"Variants.Inventory",
				"SimpleInventory",
			}
		}

		product, err := s.productRepo.FindByID(ctx, id, preloads...)
		if err != nil {
			return nil, err
		}

		// Build variant DTOs
		variantDTOs := make([]dto.VariantResponse, len(product.Variants))
		for i, v := range product.Variants {
			variantDTOs[i] = *helpers.ToVariantResponse(&v, product.BasePrice)
		}

		// Fetch review stats separately (aggregated, not all reviews)
		avgRating, reviewCount, _ := s.productRepo.GetReviewStats(ctx, id)

		// Build response
		response := helpers.ToProductResponse(product, variantDTOs, nil, &product.Merchant)
		response.AvgRating = avgRating
		response.ReviewCount = reviewCount

		return response, nil
	})

	if err != nil {
		logger.Error("Failed to fetch product", zap.Error(err))
		return nil, err
	}

	logger.Info("Product fetched successfully", zap.Bool("from_cache", err == nil))
	return response, nil
}

// ListProductsByMerchant lists products for a merchant
func (s *ProductService) ListProductsByMerchant(ctx context.Context, merchantID string, limit, offset int, activeOnly bool) ([]dto.MerchantProductResponse, error) {
	logger := s.logger.With(zap.String("operation", "ListProductsByMerchant"), zap.String("merchant_id", merchantID))
	products, err := s.productRepo.ListByMerchant(ctx, merchantID, limit, offset, activeOnly) // Fixed: Added ctx
	if err != nil {
		logger.Error("Failed to list products", zap.Error(err))
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	responses := make([]dto.MerchantProductResponse, len(products))
	for i, p := range products {
		responses[i] = *helpers.ToMerchantProductResponse(&p)
	}

	logger.Info("Products listed successfully", zap.Int("count", len(responses)))
	return responses, nil
}

// Autocomplete fetches product suggestions for search autocomplete.
func (s *ProductService) Autocomplete(ctx context.Context, prefix string, limit int) (*dto.AutocompleteResponse, error) {
	s.logger.Info("Fetching autocomplete suggestions", zap.String("prefix", prefix), zap.Int("limit", limit))

	suggestions, err := s.productRepo.AutocompleteProducts(ctx, prefix, limit)
	if err != nil {
		s.logger.Error("Failed to fetch autocomplete products", zap.Error(err))
		return nil, fmt.Errorf("autocomplete failed: %w", err)
	}

	suggest := make([]dto.ProductAutocompleteResponse, len(suggestions))
	for i, p := range suggestions {
		suggest[i] = dto.ProductAutocompleteResponse{
			ID:          p.ID,
			Name:        p.Name,
			SKU:         p.SKU,
			Description: p.Description,
		}
	}

	response := &dto.AutocompleteResponse{
		Suggestions: suggest,
	}

	s.logger.Info("Autocomplete suggestions returned", zap.Int("count", len(suggestions)))
	return response, nil
}

// GetAllProducts fetches all active products for the landing page
func (s *ProductService) GetAllProducts(ctx context.Context, limit, offset int, categoryID *uint) ([]dto.ProductResponse, int64, error) {
	logger := s.logger.With(zap.String("operation", "GetAllProducts"))

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	// Generate cache key
	filterHash := fmt.Sprintf("cat%v", categoryID)
	cacheKey := utils.ProductListCacheKey(offset/limit+1, limit, filterHash)

	type CachedResult struct {
		Products []dto.ProductResponse `json:"products"`
		Total    int64                 `json:"total"`
	}

	result, err := utils.GetOrSetCacheJSON(ctx, cacheKey, 2*time.Minute, func() (*CachedResult, error) {
		logger.Debug("Cache miss - fetching products from DB")

		// Selective preloads - no Reviews in list view
		preloads := []string{
			"Media",
			"Merchant",
			"Variants",
			"Variants.Inventory",
			"SimpleInventory",
			"Category",
		}

		products, total, err := s.productRepo.GetAllProducts(ctx, limit, offset, categoryID, preloads...)
		if err != nil {
			return nil, err
		}

		responses := make([]dto.ProductResponse, len(products))
		for i, p := range products {
			// Prepare variants DTOs
			variantDTOs := make([]dto.VariantResponse, len(p.Variants))
			for j, v := range p.Variants {
				variantDTOs[j] = *helpers.ToVariantResponse(&v, p.BasePrice)
			}

			// Get review stats (cached separately)
			avgRating, reviewCount, _ := s.productRepo.GetReviewStats(ctx, p.ID)

			resp := helpers.ToProductResponse(&p, variantDTOs, nil, &p.Merchant)
			resp.MerchantID = "" // Hide for customer view
			resp.AvgRating = avgRating
			resp.ReviewCount = reviewCount

			responses[i] = *resp
		}

		return &CachedResult{
			Products: responses,
			Total:    total,
		}, nil
	})

	if err != nil {
		logger.Error("Failed to fetch products", zap.Error(err))
		return nil, 0, err
	}

	logger.Info("Products fetched successfully", zap.Int("count", len(result.Products)), zap.Int64("total", result.Total))
	return result.Products, result.Total, nil
}

// Invalidate cache when product is updated
func (s *ProductService) InvalidateProductCache(ctx context.Context, productID string) {
	// Delete product detail cache
	err := utils.InvalidateCache(ctx, utils.ProductCacheKey(productID))
	if err != nil {
		log.Fatal(err)
	}

	// Delete list caches (all pages might contain this product)
	// err = utils.InvalidateCachePattern(ctx, "product:list:*")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	s.logger.Info("Product cache invalidated", zap.String("product_id", productID))
}

// GetAllProducts fetches all active products for the landing page
// Assumes ProductFilter is defined in the same package or imported.
func (s *ProductService) FilterProducts(ctx context.Context, filter repositories.ProductFilter, limit, offset int) ([]dto.ProductResponse, int64, error) {
	logger := s.logger.With(zap.String("operation", "FilterProducts"))

	products, total, err := s.productRepo.FilterProducts(ctx, filter, limit, offset)
	if err != nil {
		logger.Error("Failed to filter products", zap.Error(err))
		return nil, 0, err
	}

	responses := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		// Prepare variants DTOs
		variantDTOs := make([]dto.VariantResponse, len(p.Variants))
		for j, v := range p.Variants {
			variantDTOs[j] = *helpers.ToVariantResponse(&v, p.BasePrice)
		}

		// Get review stats
		avgRating, reviewCount, _ := s.productRepo.GetReviewStats(ctx, p.ID)

		resp := helpers.ToProductResponse(&p, variantDTOs, nil, &p.Merchant)
		resp.MerchantID = "" // Hide for customer view
		resp.AvgRating = avgRating
		resp.ReviewCount = reviewCount

		responses[i] = *resp
	}

	logger.Info("Products filtered successfully", zap.Int("count", len(responses)), zap.Int64("total", total))
	return responses, total, nil
}

// GetProductByID fetches a product by name
func (s *ProductService) GetProductByName(ctx context.Context, name string) (*dto.ProductResponse, error) {
	logger := s.logger.With(zap.String("operation", "GetProductByName"), zap.String("product_id", name))
	product, err := s.productRepo.FindByName(ctx, name) // Fixed: Added ctx
	if err != nil {
		if errors.Is(err, repositories.ErrProductNotFound) {
			return nil, err
		}
		logger.Error("Failed to fetch product", zap.Error(err))
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}

	variantDTOs := make([]dto.VariantResponse, len(product.Variants))
	for i, v := range product.Variants {
		variantDTOs[i] = *helpers.ToVariantResponse(&v, product.BasePrice)
	}

	// Prepare reviews DTOs
	reviewDTOs := make([]dto.ReviewResponseDTO, len(product.Reviews))
	for i, r := range product.Reviews {
		reviewDTOs[i] = *helpers.ToReviewResponse(&r)
	}

	// Use helper with loaded merchant
	response := helpers.ToProductResponse(product, variantDTOs, reviewDTOs, &product.Merchant)

	logger.Info("Product fetched successfully")
	return response, nil
}

// UpdateProduct updates a product with the provided fields
func (s *ProductService) UpdateProduct(ctx context.Context, productID string, merchantID string, input *dto.UpdateProductInput) (*dto.MerchantProductResponse, error) {
	logger := s.logger.With(zap.String("operation", "UpdateProduct"), zap.String("product_id", productID))

	// Validate input
	if err := s.validator.Struct(input); err != nil {
		logger.Error("Input validation failed", zap.Error(err))
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	// Check if product exists and belongs to merchant
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		logger.Error("Failed to find product", zap.Error(err))
		return nil, ErrInvalidProduct
	}

	if product.MerchantID != merchantID {
		logger.Error("Unauthorized access to product", zap.String("merchant_id", merchantID))
		return nil, ErrUnauthorized
	}

	// Prepare updates
	updates := make(map[string]interface{})

	if input.Name != nil {
		updates["name"] = strings.TrimSpace(*input.Name)
	}
	if input.Description != nil {
		updates["description"] = strings.TrimSpace(*input.Description)
	}
	if input.BasePrice != nil {
		updates["base_price"] = decimal.NewFromFloat(*input.BasePrice)
	}
	if input.CategoryID != nil {
		updates["category_id"] = *input.CategoryID
	}
	if input.CategoryName != nil {
		updates["category_name"] = strings.TrimSpace(*input.CategoryName)
	}
	if input.Discount != nil {
		updates["discount"] = decimal.NewFromFloat(*input.Discount)
	}
	if input.DiscountType != nil {
		updates["discount_type"] = models.DiscountType(*input.DiscountType)
	}

	// Update product
	if err := s.productRepo.UpdateProduct(ctx, productID, updates); err != nil {
		logger.Error("Failed to update product", zap.Error(err))
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	// Invalidate cache
	go s.InvalidateProductCache(context.Background(), productID)

	// Fetch updated product
	updatedProduct, err := s.productRepo.FindByID(ctx, productID, "Variants", "Media", "SimpleInventory")
	if err != nil {
		logger.Error("Failed to fetch updated product", zap.Error(err))
		return nil, fmt.Errorf("failed to fetch updated product: %w", err)
	}

	response := helpers.ToMerchantProductResponse(updatedProduct)
	logger.Info("Product updated successfully", zap.String("product_id", productID))
	return response, nil
}

// UpdateVariant updates a variant with the provided fields
func (s *ProductService) UpdateVariant(ctx context.Context, variantID string, merchantID string, input *dto.UpdateVariantInput) error {
	logger := s.logger.With(zap.String("operation", "UpdateVariant"), zap.String("variant_id", variantID))

	// Validate input
	if err := s.validator.Struct(input); err != nil {
		logger.Error("Input validation failed", zap.Error(err))
		return fmt.Errorf("invalid input: %w", err)
	}

	// Check if variant exists and belongs to merchant's product
	variant, err := s.productRepo.FindVariantByID(ctx, variantID)
	if err != nil {
		logger.Error("Failed to find variant", zap.Error(err))
		return ErrInvalidProduct
	}

	product, err := s.productRepo.FindByID(ctx, variant.ProductID)
	if err != nil {
		logger.Error("Failed to find product", zap.Error(err))
		return ErrInvalidProduct
	}

	if product.MerchantID != merchantID {
		logger.Error("Unauthorized access to variant", zap.String("merchant_id", merchantID))
		return ErrUnauthorized
	}

	// Prepare updates
	updates := make(map[string]interface{})

	if input.PriceAdjustment != nil {
		updates["price_adjustment"] = decimal.NewFromFloat(*input.PriceAdjustment)
	}
	if input.Discount != nil {
		updates["discount"] = decimal.NewFromFloat(*input.Discount)
	}
	if input.DiscountType != nil {
		updates["discount_type"] = models.DiscountType(*input.DiscountType)
	}
	if input.Attributes != nil {
		updates["attributes"] = input.Attributes
	}
	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}

	// Update variant
	if err := s.productRepo.UpdateVariant(ctx, variantID, updates); err != nil {
		logger.Error("Failed to update variant", zap.Error(err))
		return fmt.Errorf("failed to update variant: %w", err)
	}

	// Invalidate cache
	go s.InvalidateProductCache(context.Background(), variant.ProductID)

	logger.Info("Variant updated successfully", zap.String("variant_id", variantID))
	return nil
}

// BulkUpdateProducts updates multiple products and their variants
func (s *ProductService) BulkUpdateProducts(ctx context.Context, merchantID string, inputs []dto.BulkUpdateProductInput) (int, []string, error) {
	logger := s.logger.With(zap.String("operation", "BulkUpdateProducts"))

	updatedCount := 0
	errorMessages := []string{}

	for i, input := range inputs {
		// Validate product ID
		if input.ProductID == "" {
			errorMessages = append(errorMessages, fmt.Sprintf("Product %d: product_id is required", i+1))
			continue
		}

		// Update product if provided
		if input.Product != nil {
			_, err := s.UpdateProduct(ctx, input.ProductID, merchantID, input.Product)
			if err != nil {
				logger.Error("Failed to update product in bulk update", zap.Error(err), zap.String("product_id", input.ProductID))
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
				err := s.UpdateVariant(ctx, variantUpdate.VariantID, merchantID, variantUpdate.Variant)
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
	return updatedCount, errorMessages, nil
}

// BulkUpdateInventory updates multiple inventory items
func (s *ProductService) BulkUpdateInventory(ctx context.Context, merchantID string, inputs []dto.BulkInventoryUpdateInput) (int, []string, error) {
	logger := s.logger.With(zap.String("operation", "BulkUpdateInventory"))

	updatedCount := 0
	errorMessages := []string{}

	for i, input := range inputs {
		// Validate inventory ID
		if input.InventoryID == "" {
			errorMessages = append(errorMessages, fmt.Sprintf("Item %d: inventory_id is required", i+1))
			continue
		}

		// Update inventory
		err := s.UpdateInventory(ctx, input.InventoryID, input.Delta)
		if err != nil {
			logger.Error("Failed to update inventory in bulk update", zap.Error(err), zap.String("inventory_id", input.InventoryID))
			errorMessages = append(errorMessages, fmt.Sprintf("Item %d (%s): %v", i+1, input.InventoryID, err.Error()))
			continue
		}

		updatedCount++
	}

	logger.Info("Bulk inventory update completed", zap.Int("updated_count", updatedCount), zap.Int("total", len(inputs)))
	return updatedCount, errorMessages, nil
}

// UpdateInventory adjusts stock for a given inventory ID
func (s *ProductService) UpdateInventory(ctx context.Context, inventoryID string, delta int) error {
	logger := s.logger.With(zap.String("operation", "UpdateInventory"), zap.String("inventory_id", inventoryID))
	err := s.productRepo.UpdateInventoryQuantity(inventoryID, delta)
	if err != nil {
		logger.Error("Failed to update inventory", zap.Error(err))
		return fmt.Errorf("failed to update inventory: %w", err)
	}
	logger.Info("Inventory updated successfully", zap.Int("delta", delta))
	return nil
}

// DeleteProduct soft-deletes a product
func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
	logger := s.logger.With(zap.String("operation", "DeleteProduct"), zap.String("product_id", id))
	err := s.productRepo.SoftDeleteProduct(id)
	if err != nil {
		logger.Error("Failed to delete product", zap.Error(err))
		return fmt.Errorf("failed to delete product: %w", err)
	}
	logger.Info("Product deleted successfully")
	go func() {
		s.InvalidateProductCache(context.Background(), id)
	}()
	return nil
}

//Media service

// UploadMedia uploads file to Cloudinary, saves to DB
func (s *ProductService) UploadMedia(ctx context.Context, productID, merchantID, filePath, mediaType string) (*models.Media, error) {
	logger := s.logger.With(zap.String("operation", "UploadMedia"), zap.String("product_id", productID))

	// Validate merchant owns product
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil || product.MerchantID != merchantID {
		return nil, ErrUnauthorizedMedia
	}

	// Upload to Cloudinary
	params := uploader.UploadParams{
		Folder:       "merchant_media",                                         // Organized folder
		ResourceType: mediaType,                                                // image/video
		PublicID:     fmt.Sprintf("%s_%s", productID, filepath.Base(filePath)), // Unique ID
	}
	resp, err := s.cld.Upload.Upload(ctx, filePath, params)
	if err != nil {
		logger.Error("Cloudinary upload failed", zap.Error(err))
		return nil, ErrUploadFailed
	}

	// Save to DB
	//mediaType models.Media
	media := &models.Media{
		ProductID: productID,
		URL:       resp.SecureURL,
		Type:      models.MediaType(mediaType),
		PublicID:  resp.PublicID, // Store for delete/update (add to model if missing)
	}
	if err := s.productRepo.CreateMedia(ctx, media); err != nil {
		// Cleanup on failure
		_, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: media.PublicID})
		return nil, err
	}

	logger.Info("Media uploaded", zap.String("public_id", resp.PublicID))
	return media, nil
}

// UpdateMedia re-uploads or updates URL
func (s *ProductService) UpdateMedia(ctx context.Context, mediaID, productID, merchantID string, req *dto.MediaUpdateRequest) (*models.Media, error) {
	logger := s.logger.With(zap.String("operation", "UpdateMedia"), zap.String("media_id", mediaID))

	media, err := s.productRepo.FindMediaByID(ctx, mediaID)
	if err != nil || media.ProductID != productID || !s.merchantOwnsProduct(ctx, productID, merchantID) {
		return nil, ErrUnauthorizedMedia
	}

	var newURL string
	var newPublicID string
	if req.File != nil {
		// Re-upload
		resp, err := s.cld.Upload.Upload(ctx, *req.File, uploader.UploadParams{
			PublicID:     media.PublicID, // Overwrite existing
			ResourceType: string(media.Type),
		})
		if err != nil {
			logger.Error("Cloudinary re-upload failed", zap.Error(err))
			return nil, ErrUpdateFailed
		}
		newURL = resp.SecureURL
		newPublicID = resp.PublicID
	} else if req.URL != nil {
		newURL = *req.URL
	}

	// Update DB
	updates := map[string]interface{}{"url": newURL}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if newPublicID != "" {
		updates["public_id"] = newPublicID
	}
	if err := s.productRepo.UpdateMedia(ctx, mediaID, updates); err != nil {
		return nil, err
	}

	media.URL = newURL
	if req.Type != nil {
		media.Type = models.MediaType(*req.Type)
	}
	return media, nil
}

// DeleteMedia destroys on Cloudinary, deletes from DB
func (s *ProductService) DeleteMedia(ctx context.Context, mediaID, productID, merchantID, reason string) error {
	logger := s.logger.With(zap.String("operation", "DeleteMedia"), zap.String("media_id", mediaID))

	media, err := s.productRepo.FindMediaByID(ctx, mediaID)
	if err != nil || media.ProductID != productID || !s.merchantOwnsProduct(ctx, productID, merchantID) {
		return ErrUnauthorizedMedia
	}

	// Destroy on Cloudinary
	_, err = s.cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: media.PublicID})
	if err != nil {
		logger.Error("Cloudinary destroy failed", zap.Error(err))
		return ErrDeleteFailed
	}

	// Soft delete from DB
	if err := s.productRepo.DeleteMedia(ctx, mediaID); err != nil {
		return err
	}

	logger.Info("Media deleted", zap.String("public_id", media.PublicID), zap.String("reason", reason))
	return nil
}






func (s *ProductService) UploadToCloudinary(ctx context.Context, filePath, mediaType string) (string, string, error) {
	logger := s.logger.With(zap.String("operation", "UploadToCloudinary"), zap.String("file_path", filePath))

	// Validate file exists
	if filePath == "" {
		return "", "", fmt.Errorf("file path is required")
	}

	// Upload to Cloudinary
	params := uploader.UploadParams{
		Folder:       "merchant_products", // Organized folder for product images
		ResourceType: mediaType,           // "image" or "video"
		PublicID:     filepath.Base(filePath), // Use filename as base for public ID
	}

	resp, err := s.cld.Upload.Upload(ctx, filePath, params)
	if err != nil {
		logger.Error("Cloudinary upload failed", zap.Error(err))
		return "", "", fmt.Errorf("cloudinary upload failed: %w", err)
	}

	logger.Info("File uploaded to Cloudinary successfully", 
		zap.String("public_id", resp.PublicID),
		zap.String("secure_url", resp.SecureURL))

	return resp.SecureURL, resp.PublicID, nil
}

// DeleteFromCloudinary deletes a file from Cloudinary by public ID
func (s *ProductService) DeleteFromCloudinary(ctx context.Context, publicID string) error {
	logger := s.logger.With(zap.String("operation", "DeleteFromCloudinary"), zap.String("public_id", publicID))

	if publicID == "" {
		return fmt.Errorf("public ID is required")
	}

	_, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})

	if err != nil {
		logger.Error("Failed to delete from Cloudinary", zap.Error(err))
		return fmt.Errorf("cloudinary delete failed: %w", err)
	}

	logger.Info("File deleted from Cloudinary successfully")
	return nil
}
















// Helper: Check merchant owns product
func (s *ProductService) merchantOwnsProduct(ctx context.Context, productID, merchantID string) bool {
	product, err := s.productRepo.FindByID(ctx, productID)
	return err == nil && product.MerchantID == merchantID
}





// AddProductMedia adds new media to an existing product
func (s *ProductService) AddProductMedia(ctx context.Context, productID, merchantID string, mediaInputs []dto.MediaInput) error {
	logger := s.logger.With(zap.String("operation", "AddProductMedia"), zap.String("product_id", productID))

	// Verify merchant owns the product
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		logger.Error("Failed to find product", zap.Error(err))
		return ErrInvalidProduct
	}

	if product.MerchantID != merchantID {
		logger.Error("Unauthorized access to product", zap.String("merchant_id", merchantID))
		return ErrUnauthorized
	}

	// Create media records
	for _, mediaInput := range mediaInputs {
		media := &models.Media{
			ProductID: productID,
			URL:       mediaInput.URL,
			Type:      models.MediaType(mediaInput.Type),
			// Note: PublicID should be extracted from the URL or passed separately if needed
		}

		if err := s.productRepo.CreateMedia(ctx, media); err != nil {
			logger.Error("Failed to create media record", zap.Error(err))
			return fmt.Errorf("failed to add media: %w", err)
		}
	}

	// Invalidate cache
	go s.InvalidateProductCache(context.Background(), productID)

	logger.Info("Media added to product successfully", zap.Int("count", len(mediaInputs)))
	return nil
}