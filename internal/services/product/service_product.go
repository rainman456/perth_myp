package product
/*
import (
	"context"
	"errors"
	"fmt"

	//"net/url"
	"regexp"
	"strings"

	//"github.com/shopspring/decimal"
	"api-customer-merchant/internal/api/dto" // Assuming this exists for VariantInput
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

var (
	ErrInvalidProduct    = errors.New("invalid product data")
	ErrInvalidSKU        = errors.New("invalid SKU format")
	ErrInvalidMediaURL   = errors.New("invalid media URL")
	ErrInvalidAttributes = errors.New("invalid variant attributes")
	ErrUnauthorized      = errors.New("unauthorized operation")
)

// SKU validation regex: alphanumeric, hyphens, underscores, max 100 chars
var skuRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,100}$`)

type ProductService struct {
	productRepo *repositories.ProductRepository
	logger      *zap.Logger
	validator   *validator.Validate
}

func NewProductService(productRepo *repositories.ProductRepository, logger *zap.Logger) *ProductService {
	return &ProductService{
		productRepo: productRepo,
		logger:      logger,
		validator:   validator.New(),
	}
}

// CreateProductWithVariants creates a product from input DTO
func (s *ProductService) CreateProductWithVariants(ctx context.Context, input *dto.ProductInput) (*dto.ProductResponse, error) {
	logger := s.logger.With(zap.String("operation", "CreateProductWithVariants"))

	// Validate input
	if err := s.validator.Struct(input); err != nil {
		logger.Error("Input validation failed", zap.Error(err))
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	// Additional validation
	if !skuRegex.MatchString(input.SKU) {
		logger.Error("Invalid SKU format", zap.String("sku", input.SKU))
		return nil, ErrInvalidSKU
	}
	isSimple := len(input.Variants) == 0
	if isSimple && input.InitialStock == nil {
		logger.Error("Initial stock required for simple product")
		return nil, ErrInvalidProduct
	}
	for _, v := range input.Variants {
		if !skuRegex.MatchString(v.SKU) {
			logger.Error("Invalid variant SKU format", zap.String("sku", v.SKU))
			return nil, ErrInvalidSKU
		}
	}

	// Check SKU uniqueness
	if _, err := s.productRepo.FindBySKU(input.SKU); err == nil {
		logger.Warn("Duplicate product SKU", zap.String("sku", input.SKU))
		return nil, fmt.Errorf("product SKU %s already exists", input.SKU)
	} else if !errors.Is(err, repositories.ErrProductNotFound) {
		logger.Error("Failed to check product SKU", zap.Error(err))
		return nil, fmt.Errorf("failed to check SKU: %w", err)
	}
	for _, v := range input.Variants {
		if _, err := s.productRepo.FindBySKU(v.SKU); err == nil {
			logger.Warn("Duplicate variant SKU", zap.String("sku", v.SKU))
			return nil, fmt.Errorf("variant SKU %s already exists", v.SKU)
		} else if !errors.Is(err, repositories.ErrProductNotFound) {
			logger.Error("Failed to check variant SKU", zap.Error(err))
			return nil, fmt.Errorf("failed to check variant SKU: %w", err)
		}
	}

	// Map DTO to models
	product := &models.Product{
		Name:        strings.TrimSpace(input.Name),
		MerchantID:  strings.TrimSpace(input.MerchantID),
		Description: strings.TrimSpace(input.Description),
		SKU:         strings.TrimSpace(input.SKU),
		BasePrice:   decimal.NewFromFloat(input.BasePrice),
		CategoryID:  input.CategoryID,
	}
	variants := make([]models.Variant, len(input.Variants))
	for i, v := range input.Variants {
		variants[i] = models.Variant{
			SKU:             strings.TrimSpace(v.SKU),
			PriceAdjustment: decimal.NewFromFloat(v.PriceAdjustment),
			Attributes:      v.Attributes,
			IsActive:        true,
		}
	}
	media := make([]models.Media, len(input.Media))
	for i, m := range input.Media {
		media[i] = models.Media{
			URL:  strings.TrimSpace(m.URL),
			Type: models.MediaType(m.Type),
		}
	}

	// Delegate to repo
	err := s.productRepo.CreateProductWithVariantsAndInventory(ctx, product, variants, input.Variants, media, nil, isSimple)
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
	response := &dto.ProductResponse{
		ID:          product.ID,
		MerchantID:  product.MerchantID,
		Name:        product.Name,
		Description: product.Description,
		SKU:         product.SKU,
		BasePrice:   (product.BasePrice).InexactFloat64(),
		CategoryID:  product.CategoryID,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
		Media:       make([]dto.MediaResponse, len(product.Media)),
		Variants:    make([]dto.VariantResponse, len(product.Variants)),
	}
	for i, v := range product.Variants {
		response.Variants[i] = dto.VariantResponse{
			ID:              v.ID,
			ProductID:       v.ProductID,
			SKU:             v.SKU,
			PriceAdjustment: v.PriceAdjustment.InexactFloat64(),
			TotalPrice:      v.TotalPrice.InexactFloat64(),
			Attributes:      v.Attributes,
			IsActive:        v.IsActive,
			CreatedAt:       v.CreatedAt,
			UpdatedAt:       v.UpdatedAt,
			Inventory: dto.InventoryResponse{
				ID:                v.Inventory.ID,
				Quantity:          v.Inventory.Quantity,
				ReservedQuantity:  v.Inventory.ReservedQuantity,
				LowStockThreshold: v.Inventory.LowStockThreshold,
				BackorderAllowed:  v.Inventory.BackorderAllowed,
			},
		}
	}

	// Map media
	for i, m := range product.Media {
		response.Media[i] = dto.MediaResponse{
			ID:        m.ID,
			ProductID: m.ProductID,
			URL:       m.URL,
			Type:      string(m.Type),
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		}
	}

	// SimpleInventory is always nil for simple products
	response.SimpleInventory = nil

	logger.Info("Product created successfully", zap.String("product_id", product.ID))
	return response, nil
}

// GetProductByID fetches a product with optional preloads
func (s *ProductService) GetProductByID(ctx context.Context, id string, preloads ...string) (*dto.ProductResponse, error) {
	logger := s.logger.With(zap.String("operation", "GetProductByID"), zap.String("product_id", id))
	product, err := s.productRepo.FindByID(id, preloads...)
	if err != nil {
		if errors.Is(err, repositories.ErrProductNotFound) {
			return nil, err
		}
		logger.Error("Failed to fetch product", zap.Error(err))
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}

	response := &dto.ProductResponse{
		ID:          product.ID,
		MerchantID:  product.MerchantID,
		Name:        product.Name,
		Description: product.Description,
		SKU:         product.SKU,
		BasePrice:   (product.BasePrice).InexactFloat64(),
		CategoryID:  product.CategoryID,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
		Variants:    make([]dto.VariantResponse, len(product.Variants)),
		Media:       make([]dto.MediaResponse, len(product.Media)),
	}
	for i, v := range product.Variants {
		response.Variants[i] = dto.VariantResponse{
			ID:              v.ID,
			ProductID:       v.ProductID,
			SKU:             v.SKU,
			PriceAdjustment: (v.PriceAdjustment).InexactFloat64(),
			TotalPrice:      (v.TotalPrice).InexactFloat64(),
			Attributes:      v.Attributes,
			IsActive:        v.IsActive,
			CreatedAt:       v.CreatedAt,
			UpdatedAt:       v.UpdatedAt,
			Inventory: dto.InventoryResponse{
				ID:                v.Inventory.ID,
				Quantity:          v.Inventory.Quantity,
				ReservedQuantity:  v.Inventory.ReservedQuantity,
				LowStockThreshold: v.Inventory.LowStockThreshold,
				BackorderAllowed:  v.Inventory.BackorderAllowed,
			},
		}
	}
	for i, m := range product.Media {
		response.Media[i] = dto.MediaResponse{
			ID:        m.ID,
			ProductID: m.ProductID,
			URL:       m.URL,
			Type:      string(m.Type),
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		}
	}
	if product.SimpleInventory != nil {
		response.SimpleInventory = &dto.InventoryResponse{
			ID:                product.SimpleInventory.ID,
			Quantity:          product.SimpleInventory.Quantity,
			ReservedQuantity:  product.SimpleInventory.ReservedQuantity,
			LowStockThreshold: product.SimpleInventory.LowStockThreshold,
			BackorderAllowed:  product.SimpleInventory.BackorderAllowed,
		}
	}

	logger.Info("Product fetched successfully")
	return response, nil
}

// ListProductsByMerchant lists products for a merchant
func (s *ProductService) ListProductsByMerchant(ctx context.Context, merchantID string, limit, offset int, activeOnly bool) ([]dto.ProductResponse, error) {
	logger := s.logger.With(zap.String("operation", "ListProductsByMerchant"), zap.String("merchant_id", merchantID))
	products, err := s.productRepo.ListByMerchant(merchantID, limit, offset, activeOnly)
	if err != nil {
		logger.Error("Failed to list products", zap.Error(err))
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	responses := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		responses[i] = dto.ProductResponse{
			ID:          p.ID,
			MerchantID:  p.MerchantID,
			Name:        p.Name,
			Description: p.Description,
			SKU:         p.SKU,
			BasePrice:   (p.BasePrice).InexactFloat64(),
			CategoryID:  p.CategoryID,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
			Variants:    make([]dto.VariantResponse, len(p.Variants)),
			Media:       make([]dto.MediaResponse, len(p.Media)),
		}
		for j, v := range p.Variants {
			responses[i].Variants[j] = dto.VariantResponse{
				ID:              v.ID,
				ProductID:       v.ProductID,
				SKU:             v.SKU,
				PriceAdjustment: (v.PriceAdjustment).InexactFloat64(),
				TotalPrice:      (v.TotalPrice).InexactFloat64(),
				Attributes:      v.Attributes,
				IsActive:        v.IsActive,
				CreatedAt:       v.CreatedAt,
				UpdatedAt:       v.UpdatedAt,
				Inventory: dto.InventoryResponse{
					ID:                v.Inventory.ID,
					Quantity:          v.Inventory.Quantity,
					ReservedQuantity:  v.Inventory.ReservedQuantity,
					LowStockThreshold: v.Inventory.LowStockThreshold,
					BackorderAllowed:  v.Inventory.BackorderAllowed,
				},
			}
		}
		for j, m := range p.Media {
			responses[i].Media[j] = dto.MediaResponse{
				ID:        m.ID,
				ProductID: m.ProductID,
				URL:       m.URL,
				Type:      string(m.Type),
				CreatedAt: m.CreatedAt,
				UpdatedAt: m.UpdatedAt,
			}
		}
		if p.SimpleInventory != nil {
			responses[i].SimpleInventory = &dto.InventoryResponse{
				ID:                p.SimpleInventory.ID,
				Quantity:          p.SimpleInventory.Quantity,
				ReservedQuantity:  p.SimpleInventory.ReservedQuantity,
				LowStockThreshold: p.SimpleInventory.LowStockThreshold,
				BackorderAllowed:  p.SimpleInventory.BackorderAllowed,
			}
		}
	}

	logger.Info("Products listed successfully", zap.Int("count", len(responses)))
	return responses, nil
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

	products, total, err := s.productRepo.GetAllProducts(limit, offset, categoryID, "Media", "Variants", "Variants.Inventory", "SimpleInventory")
	if err != nil {
		logger.Error("Failed to fetch all products", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to fetch products: %w", err)
	}

	responses := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		responses[i] = dto.ProductResponse{
			ID:          p.ID,
			MerchantID:  "", // Exclude for customer-facing API
			Name:        p.Name,
			Description: p.Description,
			SKU:         p.SKU,
			BasePrice:   (p.BasePrice).InexactFloat64(),
			CategoryID:  p.CategoryID,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
			Variants:    make([]dto.VariantResponse, len(p.Variants)),
			Media:       make([]dto.MediaResponse, len(p.Media)),
		}
		for j, v := range p.Variants {
			responses[i].Variants[j] = dto.VariantResponse{
				ID:              v.ID,
				ProductID:       v.ProductID,
				SKU:             v.SKU,
				PriceAdjustment: (v.PriceAdjustment).InexactFloat64(),
				TotalPrice:      (v.TotalPrice).InexactFloat64(),
				Attributes:      v.Attributes,
				IsActive:        v.IsActive,
				CreatedAt:       v.CreatedAt,
				UpdatedAt:       v.UpdatedAt,
				Inventory: dto.InventoryResponse{
					ID:                v.Inventory.ID,
					Quantity:          v.Inventory.Quantity,
					ReservedQuantity:  v.Inventory.ReservedQuantity,
					LowStockThreshold: v.Inventory.LowStockThreshold,
					BackorderAllowed:  v.Inventory.BackorderAllowed,
				},
			}
		}
		for j, m := range p.Media {
			responses[i].Media[j] = dto.MediaResponse{
				ID:        m.ID,
				ProductID: m.ProductID,
				URL:       m.URL,
				Type:      string(m.Type),
				CreatedAt: m.CreatedAt,
				UpdatedAt: m.UpdatedAt,
			}
		}
		if p.SimpleInventory != nil {
			responses[i].SimpleInventory = &dto.InventoryResponse{
				ID:                p.SimpleInventory.ID,
				Quantity:          p.SimpleInventory.Quantity,
				ReservedQuantity:  p.SimpleInventory.ReservedQuantity,
				LowStockThreshold: p.SimpleInventory.LowStockThreshold,
				BackorderAllowed:  p.SimpleInventory.BackorderAllowed,
			}
		}
	}

	logger.Info("Products fetched for landing page", zap.Int("count", len(responses)), zap.Int64("total", total))
	return responses, total, nil
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
*/