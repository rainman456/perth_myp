package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"

	//"strings"
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
	ErrCategoryNotFound = errors.New("category not found")
)




type ProductFilter struct {
	CategoryID   *uint
	CategoryName *string
	CategorySlug *string
	MerchantID   *string
	MerchantName *string
	MinPrice     *decimal.Decimal
	MaxPrice     *decimal.Decimal
	InStock      *bool
	SearchTerm   *string
	OnSale       *bool
	
	// Variant attributes
	Color    *string
	Size     *string
	Material *string
	Pattern  *string
	
	// Sorting
	SortBy string // "price", "price_desc", "name", "name_desc", "created", "newest", "oldest", "rating"
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


// AutocompleteProducts fetches products matching a name prefix.
type AutocompleteResult struct {
    ID          string
    Name        string
    SKU         string
    Description string
    FinalPrice  float64
}
func (r *ProductRepository) AutocompleteProducts(ctx context.Context, prefix string, limit int) ([]AutocompleteResult, error) {
    if limit <= 0 || limit > 20 { // Sanity check
        limit = 10
    }

    var suggestions []AutocompleteResult
     err := r.db.WithContext(ctx).
        Table("products").
        Select("id, name, sku, description, final_price").
        Where("deleted_at IS NULL").
        Where("name ILIKE ?", prefix+"%").
        Order(gorm.Expr("similarity(name, ?) DESC, name ASC", prefix)).  // Requires pg_trgm extension
        Limit(limit).
        Scan(&suggestions).Error
    
    return suggestions, err

    // Map to DTO
    // suggestions := make([]dto.ProductAutocompleteResponse, len(products))
    // for i, p := range products {
    //     suggestions[i] = dto.ProductAutocompleteResponse{
    //         ID:          p.ID,
    //         Name:        p.Name,
    //         SKU:         p.SKU,
    //         Description: p.Description,
    //     }
    // }

    //return suggestions, nil
}

func (r *ProductRepository) FindByName(ctx context.Context, name string) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).Where("name = ? AND deleted_at IS NULL", name).First(&product).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrProductNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to find product by name: %w", err)
	}
	return &product, nil
}






func (r *ProductRepository) FindByID(ctx context.Context, id string, preloads ...string) (*models.Product, error) {
	var product models.Product
	query := r.db.WithContext(ctx).
	Scopes(product.Scope()).
	Where("id = ? AND deleted_at IS NULL", id)
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







func (r *ProductRepository) GetAllProductsWithCategorySlug(ctx context.Context, limit, offset int, categorySlug string, preloads ...string) ([]models.Product, int64, error) {
	var products []models.Product
	query := r.db.WithContext(ctx).Model(&models.Product{}).Where("deleted_at IS NULL")
	if categorySlug != "" {
		query = query.Where("category_slug = ?", categorySlug)
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








func (r *ProductRepository) FilterProducts(ctx context.Context, filter ProductFilter, limit, offset int) ([]models.Product, int64, error) {
    // qualify deleted_at with products to avoid ambiguity with joined tables
    query := r.db.WithContext(ctx).
        Model(&models.Product{}).
        Where("products.deleted_at IS NULL")

    // Category filters
    if filter.CategoryID != nil {
        query = query.Where("products.category_id = ?", *filter.CategoryID)
    }

    if filter.CategorySlug != nil && *filter.CategorySlug != "" {
        query = query.Joins("JOIN categories ON categories.id = products.category_id").
            Where("categories.category_slug = ?", *filter.CategorySlug)
    }

    if filter.CategoryName != nil && *filter.CategoryName != "" {
        query = query.Joins("JOIN categories ON categories.id = products.category_id").
            Where("categories.name ILIKE ?", "%"+*filter.CategoryName+"%")
    }

    // Merchant filters
    if filter.MerchantID != nil {
        query = query.Where("products.merchant_id = ?", *filter.MerchantID)
    }

    if filter.MerchantName != nil && *filter.MerchantName != "" {
        query = query.Joins("JOIN merchants ON merchants.merchant_id = products.merchant_id").
            Where("merchants.store_name ILIKE ?", "%"+*filter.MerchantName+"%")
    }

    // Price range (qualify final_price)
    if filter.MinPrice != nil {
        query = query.Where("products.final_price >= ?", *filter.MinPrice)
    }

    if filter.MaxPrice != nil {
        query = query.Where("products.final_price <= ?", *filter.MaxPrice)
    }

    // On sale filter
    if filter.OnSale != nil && *filter.OnSale {
        query = query.Where("products.discount > 0")
    }

    // Search term
    if filter.SearchTerm != nil && *filter.SearchTerm != "" {
        searchPattern := "%" + *filter.SearchTerm + "%"
        query = query.Where("products.name ILIKE ? OR products.description ILIKE ?", searchPattern, searchPattern)
    }

    // In stock filter - qualify columns inside the EXISTS subquery
    if filter.InStock != nil && *filter.InStock {
        query = query.Where(`EXISTS (
            SELECT 1 FROM inventories
            WHERE (inventories.product_id = products.id
               OR inventories.variant_id IN (
                   SELECT id FROM variants
                   WHERE variants.product_id = products.id
                   AND variants.deleted_at IS NULL
               ))
            AND (inventories.quantity - inventories.reserved_quantity) > 0
        )`)
    }

    // Variant attribute filters 
    // Sorting block: qualify order columns where needed (products.final_price, products.name, products.created_at)
	hasVariantFilter := (filter.Color != nil && *filter.Color != "") ||
		(filter.Size != nil && *filter.Size != "") ||
		(filter.Material != nil && *filter.Material != "") ||
		(filter.Pattern != nil && *filter.Pattern != "")

	if hasVariantFilter {
		variantQuery := `EXISTS (
			SELECT 1 FROM variants 
			WHERE variants.product_id = products.id 
			AND variants.deleted_at IS NULL
			AND variants.is_active = true`

		conditions := []string{}
		args := []interface{}{}

		if filter.Color != nil && *filter.Color != "" {
			conditions = append(conditions, "variants.attributes->>'color' ILIKE ?")
			args = append(args, "%"+*filter.Color+"%")
		}

		if filter.Size != nil && *filter.Size != "" {
			conditions = append(conditions, "variants.attributes->>'size' ILIKE ?")
			args = append(args, "%"+*filter.Size+"%")
		}

		if filter.Material != nil && *filter.Material != "" {
			conditions = append(conditions, "variants.attributes->>'material' ILIKE ?")
			args = append(args, "%"+*filter.Material+"%")
		}

		if filter.Pattern != nil && *filter.Pattern != "" {
			conditions = append(conditions, "variants.attributes->>'pattern' ILIKE ?")
			args = append(args, "%"+*filter.Pattern+"%")
		}

		if len(conditions) > 0 {
			variantQuery += " AND (" + strings.Join(conditions, " OR ") + ")"
			variantQuery += ")"
			query = query.Where(variantQuery, args...)
		}
	}


    // Count total before applying limit/offset
    var total int64
    countQuery := query.Session(&gorm.Session{}) // clone
    if err := countQuery.Count(&total).Error; err != nil {
        return nil, 0, fmt.Errorf("failed to count products: %w", err)
    }

    // Fetch products with selective preloads (qualify any raw SQL if used)
    var products []models.Product
    err := query.
        Preload("Category", func(db *gorm.DB) *gorm.DB {
            return db.Select("id, name, category_slug")
        }).
        Preload("Media", func(db *gorm.DB) *gorm.DB {
            return db.Select("id, product_id, url, type").
                Where("type = ?", "image").
                Order("created_at ASC").
                Limit(3)
        }).
		//"Variants.Inventory",
		//	"SimpleInventory",
        Preload("Merchant", func(db *gorm.DB) *gorm.DB {
            return db.Select("id, merchant_id, store_name, name")
        }).
        Preload("Variants", func(db *gorm.DB) *gorm.DB {
            return db.Select("id, product_id, attributes, final_price, is_active").
                Where("is_active = ?", true).
                Limit(5)
        }).
		Preload("Variants.Inventory", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, variant_id, merchant_id, quantity, reserved_quantity, backorder_allowed, updated_at")
		}).
		Preload("SimpleInventory", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, product_id, merchant_id, quantity, reserved_quantity, backorder_allowed, updated_at")
		}).
        Limit(limit).
        Offset(offset).
        Find(&products).Error

    if err != nil {
        return nil, 0, fmt.Errorf("failed to fetch products: %w", err)
    }

    return products, total, nil
}






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
		// Validate and set CategoryName if needed
		if product.CategoryID > 0 {
			var cat models.Category
			if err := tx.First(&cat, product.CategoryID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return ErrCategoryNotFound  // Define this constant, e.g., var ErrCategoryNotFound = errors.New("category not found")
				}
				return fmt.Errorf("failed to fetch category: %w", err)
			}

			if product.CategoryName == "" {
				// Automatically set from category name
				product.CategoryName = cat.Name
			} else if product.CategoryName != cat.Name {
				// Validate mismatch
				return fmt.Errorf("category name mismatch: expected '%s', got '%s'", cat.Name, product.CategoryName)
			}
		} else if product.CategoryName != "" {
			// If CategoryID is 0 but CategoryName is set, that's fine (no validation needed)
			// Or error if you want to enforce CategoryID presence
		}

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

func (r *ProductRepository) GetReviewStats(ctx context.Context, productID string) (float64, int, error) {
    var result struct {
        AvgRating   float64
        ReviewCount int
    }
    err := r.db.WithContext(ctx).
        Table("reviews").
        Select("COALESCE(AVG(rating), 0) as avg_rating, COUNT(*) as review_count").
        Where("product_id = ?", productID).
        Scan(&result).Error
    return result.AvgRating, result.ReviewCount, err
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

// UpdateProduct updates a product with the provided fields
func (r *ProductRepository) UpdateProduct(ctx context.Context, productID string, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&models.Product{}).Where("id = ?", productID).Updates(updates).Error
}

// UpdateVariant updates a variant with the provided fields
func (r *ProductRepository) UpdateVariant(ctx context.Context, variantID string, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&models.Variant{}).Where("id = ?", variantID).Updates(updates).Error
}

// FindVariantByID finds a variant by its ID
func (r *ProductRepository) FindVariantByID(ctx context.Context, id string) (*models.Variant, error) {
	var variant models.Variant
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&variant).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrProductNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to find variant by ID: %w", err)
	}
	return &variant, nil
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