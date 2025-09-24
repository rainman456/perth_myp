package product

/*
type ProductService struct {
	productRepo  *repositories.ProductRepository
	inventoryRepo *repositories.InventoryRepository
}

func NewProductService(productRepo *repositories.ProductRepository, inventoryRepo *repositories.InventoryRepository) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		inventoryRepo: inventoryRepo,
	}
}

// GetProductByID retrieves a product by its ID
func (s *ProductService) GetProductByID(id uint) (*models.Product, error) {
	if id == 0 {
		return nil, errors.New("invalid product ID")
	}
	return s.productRepo.FindByID(id)
}

// GetProductBySKU retrieves a product by its SKU
func (s *ProductService) GetProductBySKU(sku string) (*models.Product, error) {
	if strings.TrimSpace(sku) == "" {
		return nil, errors.New("SKU cannot be empty")
	}
	return s.productRepo.FindBySKU(sku)
}

// SearchProductsByName searches for products by name (case-insensitive)
func (s *ProductService) SearchProductsByName(name string) ([]models.Product, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("search name cannot be empty")
	}
	return s.productRepo.SearchByName(name)
}

// GetProductsByCategory retrieves products in a category
// In ProductService
func (s *ProductService) GetProductsByCategory(categoryID uint, limit, offset int) ([]models.Product, error) {
    if categoryID == 0 {
        return nil, errors.New("invalid category ID")
    }
    return s.productRepo.FindByCategoryWithPagination(categoryID, limit, offset)
}

// CreateProduct creates a new product for a merchant
func (s *ProductService) CreateProduct(product *models.Product, merchantID uint) error {
	if product == nil {
		return errors.New("product cannot be nil")
	}
	if merchantID == 0 {
		return errors.New("invalid merchant ID")
	}
	if strings.TrimSpace(product.Name) == "" {
		return errors.New("product name cannot be empty")
	}
	if strings.TrimSpace(product.SKU) == "" {
		return errors.New("SKU cannot be empty")
	}
	if product.Price <= 0 {
		return errors.New("price must be positive")
	}
	if product.CategoryID == 0 {
		return errors.New("category ID must be set")
	}

	// Check if SKU is unique
	if _, err := s.productRepo.FindBySKU(product.SKU); err == nil {
		return errors.New("SKU already exists")
	}

	product.MerchantID = merchantID
	return s.productRepo.Create(product)
}

// UpdateProduct updates an existing product (merchant only)
func (s *ProductService) UpdateProduct(product *models.Product, merchantID uint) error {
	if product == nil || product.ID == 0 {
		return errors.New("invalid product or product ID")
	}
	if merchantID == 0 {
		return errors.New("invalid merchant ID")
	}
	if strings.TrimSpace(product.Name) == "" {
		return errors.New("product name cannot be empty")
	}
	if strings.TrimSpace(product.SKU) == "" {
		return errors.New("SKU cannot be empty")
	}
	if product.Price <= 0 {
		return errors.New("price must be positive")
	}
	if product.CategoryID == 0 {
		return errors.New("category ID must be set")
	}

	// Verify product belongs to merchant
	existing, err := s.productRepo.FindByID(product.ID)
	if err != nil {
		return err
	}
	if existing.MerchantID != merchantID {
		return errors.New("product does not belong to merchant")
	}

	// Check if SKU is unique (excluding current product)
	if p, err := s.productRepo.FindBySKU(product.SKU); err == nil && p.ID != product.ID {
		return errors.New("SKU already exists")
	}

	return s.productRepo.Update(product)
}

// DeleteProduct deletes a product (merchant only)
func (s *ProductService) DeleteProduct(id uint, merchantID uint) error {
	if id == 0 {
		return errors.New("invalid product ID")
	}
	if merchantID == 0 {
		return errors.New("invalid merchant ID")
	}

	// Verify product belongs to merchant
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return err
	}
	if product.MerchantID != merchantID {
		return errors.New("product does not belong to merchant")
	}

	return s.productRepo.Delete(id)
}
*/

/*

import (
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"
	"api-customer-merchant/internal/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type ProductService struct {
	productRepo  *repositories.ProductRepository
	inventoryRepo *repositories.InventoryRepository
}

func NewProductService(productRepo *repositories.ProductRepository, inventoryRepo *repositories.InventoryRepository) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		inventoryRepo: inventoryRepo,
	}
}

// GetProductByID retrieves a product by its ID
func (s *ProductService) GetProductByID(id string) (*models.Product, error) {
    if id == "" {
        return nil, errors.New("invalid product ID")
    }

    ctx := context.Background()
    key := fmt.Sprintf("product:%s", id)
    data, err := utils.GetOrSetCache(ctx, key, 10*time.Minute, func() (any, error) {
        product, err := s.productRepo.FindByID(id)
        if err != nil {
            return nil, err
        }
        return json.Marshal(product) // Serialize
    })
    if err != nil {
        return nil, err
    }

    var product models.Product
    if err := json.Unmarshal([]byte(data.(string)), &product); err != nil {
        return nil, err
    }
    return &product, nil
}

// GetProductBySKU retrieves a product by its SKU
func (s *ProductService) GetProductBySKU(sku string) (*models.Product, error) {
	if strings.TrimSpace(sku) == "" {
		return nil, errors.New("SKU cannot be empty")
	}
	return s.productRepo.FindBySKU(sku)
}

// SearchProductsByName searches for products by name (case-insensitive)
func (s *ProductService) SearchProductsByName(name string) ([]models.Product, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("search name cannot be empty")
	}
	return s.productRepo.SearchByName(name)
}

// GetProductsByCategory retrieves products in a category
// In ProductService
func (s *ProductService) GetProductsByCategory(categoryID string, limit, offset int) ([]models.Product, error) {
    if categoryID == "" {
        return nil, errors.New("invalid category ID")
    }
    return s.productRepo.FindByCategoryWithPagination(categoryID, limit, offset)
}

// CreateProduct creates a new product for a merchant
func (s *ProductService) CreateProduct(product *models.Product, merchantID string) error {
	if product == nil {
		return errors.New("product cannot be nil")
	}
	if merchantID == "" {
		return errors.New("invalid merchant ID")
	}
	if strings.TrimSpace(product.Name) == "" {
		return errors.New("product name cannot be empty")
	}
	if strings.TrimSpace(product.SKU) == "" {
		return errors.New("SKU cannot be empty")
	}
	if product.Price <= 0 {
		return errors.New("price must be positive")
	}
	if product.CategoryID == 0 {
		return errors.New("category ID must be set")
	}

	// Check if SKU is unique
	if _, err := s.productRepo.FindBySKU(product.SKU); err == nil {
		return errors.New("SKU already exists")
	}

	product.MerchantID = merchantID
	return s.productRepo.Create(product)
}

// UpdateProduct updates an existing product (merchant only)
func (s *ProductService) UpdateProduct(product *models.Product, merchantID string) error {
	if product == nil || product.ID == "" {
		return errors.New("invalid product or product ID")
	}
	if merchantID == "" {
		return errors.New("invalid merchant ID")
	}
	if strings.TrimSpace(product.Name) == "" {
		return errors.New("product name cannot be empty")
	}
	if strings.TrimSpace(product.SKU) == "" {
		return errors.New("SKU cannot be empty")
	}
	if product.Price <= 0 {
		return errors.New("price must be positive")
	}
	if product.CategoryID == 0 {
		return errors.New("category ID must be set")
	}

	// Verify product belongs to merchant
	existing, err := s.productRepo.FindByID(product.ID)
	if err != nil {
		return err
	}
	if existing.MerchantID != merchantID {
		return errors.New("product does not belong to merchant")
	}

	// Check if SKU is unique (excluding current product)
	if p, err := s.productRepo.FindBySKU(product.SKU); err == nil && p.ID != product.ID {
		return errors.New("SKU already exists")
	}

	return s.productRepo.Update(product)
}

// DeleteProduct deletes a product (merchant only)
func (s *ProductService) DeleteProduct(id string, merchantID string) error {
	if id == "" {
		return errors.New("invalid product ID")
	}
	if merchantID == "" {
		return errors.New("invalid merchant ID")
	}

	// Verify product belongs to merchant
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return err
	}
	if product.MerchantID != merchantID {
		return errors.New("product does not belong to merchant")
	}

	return s.productRepo.Delete(id)
}

func (s *ProductService) GetProductsByMerchantID(merchantID string) ([]models.Product, error) {
	return s.productRepo.FindByMerchantID(merchantID)
}


*/
