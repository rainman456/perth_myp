package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

var (
	ErrProductNotFound  = errors.New("product not found")
	ErrDuplicateSKU     = errors.New("duplicate SKU")
	ErrInvalidInventory = errors.New("invalid inventory setup")
	ErrMerchantNotFound = errors.New("merchant not found")
)




type ProductFilter struct {
    CategoryName   *string
    CategoryID     *uint
    MinPrice       *decimal.Decimal
    MaxPrice       *decimal.Decimal
    InStock        *bool
    VariantAttrs   map[string]interface{}
    MerchantName   *string
}


type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{db: db.DB}
}

// func (r *ProductRepository) FindBySKU(sku string) (*models.Product, error) {
// 	var product models.Product
// 	err := r.db.Where("sku = ?", sku).First(&product).Error
// 	if errors.Is(err, gorm.ErrRecordNotFound) {
// 		return nil, ErrProductNotFound
// 	} else if err != nil {
// 		return nil, fmt.Errorf("failed to find product by SKU: %w", err)
// 	}
// 	return &product, nil
// }



func (r *ProductRepository) FindBySKU(ctx context.Context, sku string) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).Where("sku = ? AND deleted_at IS NULL", sku).First(&product).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrProductNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to find product by SKU: %w", err)
	}
	return &product, nil
}



// func (r *ProductRepository) FindByID(id string, preloads ...string) (*models.Product, error) {
// 	var product models.Product
// 	query := r.db.Where("id = ?", id)
// 	for _, preload := range preloads {
// 		query = query.Preload(preload)
// 	}
// 	err := query.First(&product).Error
// 	if errors.Is(err, gorm.ErrRecordNotFound) {
// 		return nil, ErrProductNotFound
// 	} else if err != nil {
// 		return nil, fmt.Errorf("failed to find product by ID: %w", err)
// 	}
// 	return &product, nil
// }






func (r *ProductRepository) FindByID(ctx context.Context, id string, preloads ...string) (*models.Product, error) {
	var product models.Product
	query := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id)
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	err := query.First(&product).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrProductNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to find product by ID: %w", err)
	}
	return &product, nil
}





// func (r *ProductRepository) ListByMerchant(merchantID string, limit, offset int, filterActive bool) ([]models.Product, error) {
// 	var products []models.Product
// 	query := r.db.Where("merchant_id = ?", merchantID).Limit(limit).Offset(offset)
// 	if filterActive {
// 		query = query.Where("deleted_at IS NULL")
// 	}
// 	err := query.Find(&products).Error
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to list products: %w", err)
// 	}
// 	return products, nil
// }






func (r *ProductRepository) ListByMerchant(ctx context.Context, merchantID string, limit, offset int, filterActive bool) ([]models.Product, error) {
	var products []models.Product
	query := r.db.WithContext(ctx).Where("merchant_id = ?", merchantID).Limit(limit).Offset(offset)
	if filterActive {
		query = query.Where("deleted_at IS NULL")
	}
	err := query.Find(&products).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	return products, nil
}







// func (r *ProductRepository) GetAllProducts(limit, offset int, categoryID *uint, preloads ...string) ([]models.Product, int64, error) {
// 	var products []models.Product
// 	query := r.db.Model(&models.Product{}).Where("deleted_at IS NULL")
// 	if categoryID != nil {
// 		query = query.Where("category_id = ?", *categoryID)
// 	}
// 	var total int64
// 	if err := query.Count(&total).Error; err != nil {
// 		return nil, 0, fmt.Errorf("failed to count products: %w", err)
// 	}
// 	for _, preload := range preloads {
// 		query = query.Preload(preload)
// 	}
// 	query = query.Limit(limit).Offset(offset).Order("created_at DESC")
// 	err := query.Find(&products).Error
// 	if err != nil {
// 		return nil, 0, fmt.Errorf("failed to fetch products: %w", err)
// 	}
// 	return products, total, nil
// }






func (r *ProductRepository) GetAllProducts(ctx context.Context, limit, offset int, categoryID *uint, preloads ...string) ([]models.Product, int64, error) {
	var products []models.Product
	query := r.db.WithContext(ctx).Model(&models.Product{}).Where("deleted_at IS NULL")
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	query = query.Limit(limit).Offset(offset).Order("created_at DESC")
	err := query.Find(&products).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch products: %w", err)
	}
	return products, total, nil
}






func (r *ProductRepository) ProductsFilter(
    ctx context.Context,
    filter ProductFilter,
    limit, offset int,
    preloads ...string,
) ([]models.Product, int64, error) {
    var products []models.Product

    query := r.db.WithContext(ctx).
        Model(&models.Product{}).
        Joins("LEFT JOIN categories ON categories.id = products.category_id").
        Joins("LEFT JOIN merchants ON merchants.id = products.merchant_id").
        Joins("LEFT JOIN variants ON variants.product_id = products.id").
        Joins("LEFT JOIN inventories ON inventories.product_id = products.id OR inventories.variant_id = variants.id").
        Where("products.deleted_at IS NULL")

    // --- Apply filters ---
    if filter.CategoryID != nil {
        query = query.Where("products.category_id = ?", *filter.CategoryID)
    }
    if filter.CategoryName != nil {
        query = query.Where("categories.name ILIKE ?", "%"+*filter.CategoryName+"%")
    }
    if filter.MinPrice != nil {
        query = query.Where("products.base_price >= ?", *filter.MinPrice)
    }
    if filter.MaxPrice != nil {
        query = query.Where("products.base_price <= ?", *filter.MaxPrice)
    }
    if filter.InStock != nil {
        if *filter.InStock {
            query = query.Where("(inventories.quantity - inventories.reserved_quantity) > 0")
        } else {
            query = query.Where("(inventories.quantity - inventories.reserved_quantity) <= 0")
        }
    }
    if filter.MerchantName != nil {
        query = query.Where("merchant.store_name ILIKE ?", "%"+*filter.MerchantName+"%")
    }
    if len(filter.VariantAttrs) > 0 {
        for key, val := range filter.VariantAttrs {
            // Postgres JSONB query on variant.attributes
            query = query.Where("variants.attributes ->> ? = ?", key, fmt.Sprintf("%v", val))
        }
    }

    // --- Count total ---
    var total int64
    if err := query.Distinct("products.id").Count(&total).Error; err != nil {
        return nil, 0, fmt.Errorf("failed to count products: %w", err)
    }

    // --- Preloads ---
    for _, preload := range preloads {
        query = query.Preload(preload)
    }

    // --- Fetch results ---
    err := query.Distinct("products.id").
        Limit(limit).
        Offset(offset).
        Order("products.created_at DESC").
        Find(&products).Error

    if err != nil {
        return nil, 0, fmt.Errorf("failed to fetch products: %w", err)
    }

    return products, total, nil
}





// func (r *ProductRepository) CreateProductWithVariantsAndInventory(ctx context.Context, product *models.Product, variants []models.Variant, variantInputs []dto.VariantInput, media []models.Media, simpleInitialStock *int, isSimple bool) error {
// 	if isSimple && len(variants) > 0 {
// 		return ErrInvalidInventory // Cannot have variants for simple products
// 	}
// 	if !isSimple && (len(variants) == 0 || len(variants) != len(variantInputs)) {
// 		return ErrInvalidInventory // Must provide matching variants and inputs
// 	}

// 	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
// 		// Create product
// 		if err := tx.Create(product).Error; err != nil {
// 			if errors.Is(err, gorm.ErrDuplicatedKey) {
// 				return ErrDuplicateSKU
// 			}
// 			return fmt.Errorf("failed to create product: %w", err)
// 		}

// 		if !isSimple {
// 			// Variant-based product
// 			for i := range variants {
// 				variants[i].ProductID = product.ID
// 				if err := tx.Create(&variants[i]).Error; err != nil {
// 					return fmt.Errorf("failed to create variant: %w", err)
// 				}
// 				variantIDPtr := variants[i].ID
// 				inventory := models.Inventory{
// 					VariantID:         variantIDPtr,
// 					ProductID:         nil,
// 					MerchantID:        product.MerchantID,
// 					Quantity:          variantInputs[i].InitialStock,
// 					ReservedQuantity:  0,
// 					LowStockThreshold: 10,
// 					BackorderAllowed:  false,
// 				}
// 				if err := tx.Create(&inventory).Error; err != nil {
// 					return fmt.Errorf("failed to create variant inventory: %w", err)
// 				}
// 				variants[i].Inventory = inventory
// 			}
// 		}
// 		// Note: Skip VendorInventory creation for simple products

// 		// Create media
// 		for i := range media {
// 			media[i].ProductID = product.ID
// 			if err := tx.Create(&media[i]).Error; err != nil {
// 				return fmt.Errorf("failed to create media: %w", err)
// 			}
// 		}

// 		// Reload with preloads
// 		preloadQuery := tx.Where("id = ?", product.ID)
// 		if !isSimple {
// 			preloadQuery = preloadQuery.Preload("Variants.Inventory")
// 		}
// 		preloadQuery = preloadQuery.Preload("Media")
// 		if err := preloadQuery.First(product).Error; err != nil {
// 			return fmt.Errorf("failed to preload associations: %w", err)
// 		}

// 		return nil
// 	})
// }




func (r *ProductRepository) CreateProductWithVariantsAndInventory(ctx context.Context, product *models.Product, variants []models.Variant, variantInputs []dto.VariantInput, media []models.Media, simpleInitialStock *int, isSimple bool) error {
	// Validate Merchant exists
	var merchant models.Merchant
	if err := r.db.WithContext(ctx).Where("merchant_id = ?", product.MerchantID).First(&merchant).Error; err != nil {
		return ErrMerchantNotFound
	}

	// Validate SKU uniqueness
	if p, _ := r.FindBySKU(ctx, product.SKU); p != nil {
		return ErrDuplicateSKU
	}
	// for _, v := range variants {
	// 	if v2, _ := r.db.WithContext(ctx).Where("sku = ? AND deleted_at IS NULL", v.SKU).First(&models.Variant{}).Error; v2 == nil {
	// 		return ErrDuplicateSKU
	// 	}
	// }

	for _, v := range variants {
var temp models.Variant
if err := r.db.WithContext(ctx).Where("sku = ? AND deleted_at IS NULL", v.SKU).First(&temp).Error; err == nil {
return ErrDuplicateSKU
} else if !errors.Is(err, gorm.ErrRecordNotFound) {
return err  // Propagate unexpected errors
}
}

	// Validate inputs
	if isSimple && len(variants) > 0 {
		return ErrInvalidInventory // Cannot have variants for simple products
	}
	if !isSimple && (len(variants) == 0 || len(variants) != len(variantInputs)) {
		return ErrInvalidInventory // Must provide matching variants and inputs
	}
	for _, vi := range variantInputs {
		if vi.InitialStock < 0 {
			return errors.New("initial stock cannot be negative")
		}
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create product
		if err := tx.Create(product).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return ErrDuplicateSKU
			}
			return fmt.Errorf("failed to create product: %w", err)
		}

		if isSimple {
			// Create inventory for simple product
			if simpleInitialStock == nil {
				return errors.New("simpleInitialStock required for simple products")
			}
			inventory := models.Inventory{
				ProductID:         &product.ID,
				MerchantID:        product.MerchantID,
				Quantity:          *simpleInitialStock,
				ReservedQuantity:  0,
				LowStockThreshold: 5, // From merged model
				BackorderAllowed:  false,
			}
			if err := tx.Create(&inventory).Error; err != nil {
				return fmt.Errorf("failed to create simple inventory: %w", err)
			}
			product.SimpleInventory = &inventory
		} else {
			// Create variants and their inventories
			for i := range variants {
				variants[i].ProductID = product.ID
				if err := tx.Create(&variants[i]).Error; err != nil {
					return fmt.Errorf("failed to create variant: %w", err)
				}
				inventory := models.Inventory{
					//ProductID:         &product.ID, // Explicit link to product (per requirement)
					VariantID:         &variants[i].ID,
					MerchantID:        product.MerchantID,
					Quantity:          variantInputs[i].InitialStock,
					ReservedQuantity:  0,
					LowStockThreshold: 5,
					BackorderAllowed:  false,
				}
				if err := tx.Create(&inventory).Error; err != nil {
					return fmt.Errorf("failed to create variant inventory: %w", err)
				}
				variants[i].Inventory = inventory
			}
		}

		// Create media
		for i := range media {
			media[i].ProductID = product.ID
			if err := tx.Create(&media[i]).Error; err != nil {
				return fmt.Errorf("failed to create media: %w", err)
			}
		}

		// Reload with preloads
		preloadQuery := tx.Where("id = ? AND deleted_at IS NULL", product.ID)
		if !isSimple {
			preloadQuery = preloadQuery.Preload("Variants.Inventory")
		}
		preloadQuery = preloadQuery.Preload("SimpleInventory").Preload("Media")
		if err := preloadQuery.First(product).Error; err != nil {
			return fmt.Errorf("failed to preload associations: %w", err)
		}

		return nil
	})
}




func (r *ProductRepository) UpdateInventoryQuantity(inventoryID string, delta int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var inventory models.Inventory
		if err := tx.First(&inventory, "id = ?", inventoryID).Error; err != nil {
			return fmt.Errorf("failed to find inventory: %w", err)
		}
		newQuantity := inventory.Quantity + delta
		if newQuantity < 0 && !inventory.BackorderAllowed {
			return errors.New("insufficient stock and backorders not allowed")
		}
		inventory.Quantity = newQuantity
		return tx.Save(&inventory).Error
	})
}

func (r *ProductRepository) SoftDeleteProduct(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Product{}).Error
}





//For media uploads
func (r *ProductRepository) CreateMedia(ctx context.Context, media *models.Media) error {
	return r.db.WithContext(ctx).Create(media).Error
}

// FindMediaByID fetches media
func (r *ProductRepository) FindMediaByID(ctx context.Context, id string) (*models.Media, error) {
	var media models.Media
	err := r.db.WithContext(ctx).Scopes(r.activeScope()).First(&media, "id = ?", id).Error
	return &media, err
}

// UpdateMedia updates fields
func (r *ProductRepository) UpdateMedia(ctx context.Context, id string, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&models.Media{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteMedia soft-deletes
func (r *ProductRepository) DeleteMedia(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&models.Media{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

// activeScope (if soft delete)
func (r *ProductRepository) activeScope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB { return db.Where("deleted_at IS NULL") }
}