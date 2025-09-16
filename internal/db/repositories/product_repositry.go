package repositories

import (
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	"errors"

	"gorm.io/gorm"
)

/*
type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{db: db.DB}
}

// Create adds a new product to the database
func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

// FindByID retrieves a product by ID with associated Merchant and Category
func (r *ProductRepository) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Merchant").Preload("Category").First(&product, id).Error
	return &product, err
}

// FindBySKU retrieves a product by SKU with associated Merchant and Category
func (r *ProductRepository) FindBySKU(sku string) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Merchant").Preload("Category").Where("sku = ?", sku).First(&product).Error
	return &product, err
}

// SearchByName retrieves products matching a name (partial match)
func (r *ProductRepository) SearchByName(name string) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Merchant").Preload("Category").Where("name ILIKE ?", "%"+name+"%").Find(&products).Error
	return products, err
}

// FindByMerchantID retrieves all products for a merchant
func (r *ProductRepository) FindByMerchantID(merchantID uint) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Category").Where("merchant_id = ?", merchantID).Find(&products).Error
	return products, err
}

// FindByCategoryID retrieves all products in a category
func (r *ProductRepository) FindByCategoryID(categoryID uint) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Merchant").Where("category_id = ?", categoryID).Find(&products).Error
	return products, err
}

// Update modifies an existing product
func (r *ProductRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

// Delete removes a product by ID
func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}
// In ProductRepository
func (r *ProductRepository) FindByCategoryWithPagination(categoryID uint, limit, offset int) ([]models.Product, error) {
    var products []models.Product
    err := r.db.Preload("Merchant").Preload("Category").Where("category_id = ?", categoryID).Limit(limit).Offset(offset).Find(&products).Error
    return products, err
}
*/

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{db: db.DB}
}

// Create adds a new product to the database
func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

// FindByID retrieves a product by ID with associated Merchant and Category
func (r *ProductRepository) FindByID(id string) (*models.Product, error) {
    var product models.Product
    err := r.db.Preload("Merchant").Preload("Category").First(&product, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("product not found")
        }
        return nil, err
    }
    return &product, nil
}

// FindBySKU retrieves a product by SKU with associated Merchant and Category
func (r *ProductRepository) FindBySKU(sku string) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Merchant").Preload("Category").Where("sku = ?", sku).First(&product).Error
	return &product, err
}

// SearchByName retrieves products matching a name (partial match)
func (r *ProductRepository) SearchByName(name string) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Merchant").Preload("Category").Where("name ILIKE ?", "%"+name+"%").Find(&products).Error
	return products, err
}

// FindByMerchantID retrieves all products for a merchant
func (r *ProductRepository) FindByMerchantID(merchantID string) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Category").Where("merchant_id = ?", merchantID).Find(&products).Error
	return products, err
}

// FindByCategoryID retrieves all products in a category
func (r *ProductRepository) FindByCategoryID(categoryID string) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Merchant").Where("category_id = ?", categoryID).Find(&products).Error
	return products, err
}

// Update modifies an existing product
func (r *ProductRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

// Delete removes a product by ID
func (r *ProductRepository) Delete(id string) error {
	return r.db.Delete(&models.Product{}, id).Error
}
// In ProductRepository
func (r *ProductRepository) FindByCategoryWithPagination(categoryID string, limit, offset int) ([]models.Product, error) {
    var products []models.Product
    err := r.db.Preload("Merchant").Preload("Category").Where("category_id = ?", categoryID).Limit(limit).Offset(offset).Find(&products).Error
    return products, err
}

