package product

import (
	"context"
	"errors"
	"fmt"

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

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var (
	ErrInvalidProduct    = errors.New("invalid product data")
	//ErrInvalidSKU        = errors.New("invalid SKU format")
	ErrInvalidMediaURL   = errors.New("invalid media URL")
	ErrInvalidAttributes = errors.New("invalid variant attributes")
	ErrUnauthorized      = errors.New("unauthorized operation")
	ErrUploadFailed     = errors.New("upload to Cloudinary failed")
	ErrUpdateFailed     = errors.New("update failed")
	ErrDeleteFailed     = errors.New("delete failed")
	ErrUnauthorizedMedia = errors.New("unauthorized for this media")
)

// SKU validation regex: alphanumeric, hyphens, underscores, max 100 chars
//var skuRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,100}$`)

type ProductService struct {
	productRepo *repositories.ProductRepository
	logger      *zap.Logger
	validator   *validator.Validate
	cld         *cloudinary.Cloudinary
	//config  *config.Config
	
}

func NewProductService(productRepo *repositories.ProductRepository,  cfg *config.Config,logger *zap.Logger) *ProductService {
	cld, err := cloudinary.NewFromParams(cfg.CloudinaryCloudName, cfg.CloudinaryAPIKey, cfg.CloudinaryAPISecret)
	if err != nil {
		logger.Fatal("Cloudinary init failed", zap.Error(err))
	}

	return &ProductService{
		productRepo: productRepo,
		logger:      logger,
		validator:   validator.New(),
		cld:         cld,
	}
}

// CreateProductWithVariants creates a product from input DTO
func (s *ProductService) CreateProductWithVariants(ctx context.Context, merchant_id string,input *dto.ProductInput) (*dto.MerchantProductResponse, error) {
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
		BasePrice:   decimal.NewFromFloat(input.BasePrice),
		Discount: decimal.NewFromFloat(input.Discount),
		DiscountType: models.DiscountType(input.DiscountType),
		CategoryID:  input.CategoryID,
	}
	variants := make([]models.Variant, len(input.Variants))
	for i, v := range input.Variants {
		variants[i] = models.Variant{
			//SKU:             strings.TrimSpace(v.SKU),
			PriceAdjustment: decimal.NewFromFloat(v.PriceAdjustment),
			Discount: decimal.NewFromFloat(input.Discount),
			DiscountType: models.DiscountType(input.DiscountType),
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
	response := helpers.ToMerchantProductResponse(product)

	logger.Info("Product created successfully", zap.String("product_id", product.ID))
	return response, nil
}

// GetProductByID fetches a product with optional preloads
func (s *ProductService) GetProductByID(ctx context.Context, id string, preloads ...string) (*dto.ProductResponse, error) {
	logger := s.logger.With(zap.String("operation", "GetProductByID"), zap.String("product_id", id))
	product, err := s.productRepo.FindByID(ctx, id, preloads...)  // Fixed: Added ctx
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


// ListProductsByMerchant lists products for a merchant
func (s *ProductService) ListProductsByMerchant(ctx context.Context, merchantID string, limit, offset int, activeOnly bool) ([]dto.MerchantProductResponse, error) {
	logger := s.logger.With(zap.String("operation", "ListProductsByMerchant"), zap.String("merchant_id", merchantID))
	products, err := s.productRepo.ListByMerchant(ctx, merchantID, limit, offset, activeOnly)  // Fixed: Added ctx
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

	products, total, err := s.productRepo.GetAllProducts(ctx, limit, offset, categoryID, "Media", "Variants", "Variants.Inventory", "SimpleInventory")  // Fixed: Added ctx (resolves type shifts)
	if err != nil {
		logger.Error("Failed to fetch all products", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to fetch products: %w", err)
	}

	responses := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		// Prepare variants DTOs
		variantDTOs := make([]dto.VariantResponse, len(p.Variants))
		for j, v := range p.Variants {
			variantDTOs[j] = *helpers.ToVariantResponse(&v, p.BasePrice)
		}

		// Prepare reviews DTOs
		reviewDTOs := make([]dto.ReviewResponseDTO, len(p.Reviews))
		for j, r := range p.Reviews {
			reviewDTOs[j] = *helpers.ToReviewResponse(&r)
		}

		// Use helper (nil merchant for customer-facing, and set MerchantID = "")
		resp := helpers.ToProductResponse(&p, variantDTOs, reviewDTOs, nil)
		resp.MerchantID = ""
		responses[i] = *resp
	}

	logger.Info("Products fetched for landing page", zap.Int("count", len(responses)), zap.Int64("total", total))
	return responses, total, nil
}





// GetAllProducts fetches all active products for the landing page
// Assumes ProductFilter is defined in the same package or imported.
 type ProductFilter struct {
     CategoryName   *string
     CategoryID     *uint
     MinPrice       *decimal.Decimal
     MaxPrice       *decimal.Decimal
     InStock        *bool
     VariantAttrs   map[string]interface{}
     MerchantName   *string
 }

func (s *ProductService) FilterProducts(ctx context.Context, filter ProductFilter, limit, offset int) ([]dto.ProductResponse, int64, error) {
	logger := s.logger.With(zap.String("operation", "FilterProducts"))

	// --- pagination sanitization ---
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}




	// --- fetch products from repository using the provided filter ---
	repoFilter := repositories.ProductFilter{
    CategoryName: filter.CategoryName,
    CategoryID:   filter.CategoryID,
    MinPrice:     filter.MinPrice,
    MaxPrice:     filter.MaxPrice,
    InStock:      filter.InStock,
    VariantAttrs: filter.VariantAttrs,
    MerchantName: filter.MerchantName,
}

products, total, err := s.productRepo.ProductsFilter(ctx, repoFilter, limit, offset, "Media", "Variants", "Variants.Inventory", "SimpleInventory")

	if err != nil {
		logger.Error("Failed to fetch products", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to fetch products: %w", err)
	}

	// --- map DB models -> DTOs ---
	responses := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		// Prepare variants DTOs
		variantDTOs := make([]dto.VariantResponse, len(p.Variants))
		for j, v := range p.Variants {
			variantDTOs[j] = *helpers.ToVariantResponse(&v, p.BasePrice)
		}

		// Prepare reviews DTOs
		reviewDTOs := make([]dto.ReviewResponseDTO, len(p.Reviews))
		for j, r := range p.Reviews {
			reviewDTOs[j] = *helpers.ToReviewResponse(&r)
		}

		// Use helper (nil merchant for customer-facing, and set MerchantID = "")
		resp := helpers.ToProductResponse(&p, variantDTOs, reviewDTOs, nil)
		resp.MerchantID = ""
		responses[i] = *resp
	}

	logger.Info("Products fetched for filter", zap.Int("count", len(responses)), zap.Int64("total", total))
	return responses, total, nil
}



// GetProductByID fetches a product by name
func (s *ProductService) GetProductByName(ctx context.Context, name string) (*dto.ProductResponse, error) {
	logger := s.logger.With(zap.String("operation", "GetProductByName"), zap.String("product_id", name))
	product, err := s.productRepo.FindByName(ctx,  name)  // Fixed: Added ctx
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
		Folder:     "merchant_media", // Organized folder
		ResourceType: mediaType, // image/video
		PublicID:    fmt.Sprintf("%s_%s", productID, filepath.Base(filePath)), // Unique ID
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
			PublicID:    media.PublicID, // Overwrite existing
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





// Helper: Check merchant owns product
func (s *ProductService) merchantOwnsProduct(ctx context.Context, productID, merchantID string) bool {
	product, err := s.productRepo.FindByID(ctx, productID)
	return err == nil && product.MerchantID == merchantID
}








