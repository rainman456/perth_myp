package repositories

import (
	//"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	"context"
	"errors"

	//"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)


type ProductCacheRepository struct {
    // db    *gorm.DB
    // cache *redis.Client
}

// FindByID with selective loading
func (r *ProductRepository) FindByIDOptimized(ctx context.Context, id string, fields []string) (*models.Product, error) {
    var product models.Product
    
    query := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id)
    
    // Select only needed fields if specified
    if len(fields) > 0 {
        query = query.Select(fields)
    }
    
    err := query.First(&product).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, ErrProductNotFound
    }
    
    return &product, err
}

// Batch fetch products (for cart/order processing)
func (r *ProductRepository) FindByIDs(ctx context.Context, ids []string, preloads ...string) (map[string]*models.Product, error) {
    var products []models.Product
    
    query := r.db.WithContext(ctx).Where("id IN ? AND deleted_at IS NULL", ids)
    
    for _, preload := range preloads {
        query = query.Preload(preload)
    }
    
    if err := query.Find(&products).Error; err != nil {
        return nil, err
    }
    
    // Return as map for O(1) lookup
    productMap := make(map[string]*models.Product, len(products))
    for i := range products {
        productMap[products[i].ID] = &products[i]
    }
    
    return productMap, nil
}

// Efficient count query
func (r *ProductRepository) CountByMerchant(ctx context.Context, merchantID string) (int64, error) {
    var count int64
    err := r.db.WithContext(ctx).
        Model(&models.Product{}).
        Where("merchant_id = ? AND deleted_at IS NULL", merchantID).
        Count(&count).Error
    
    return count, err
}

// Get products with inventory status (optimized)
type ProductWithInventoryStatus struct {
	models.Product
	InStock      bool
	TotalStock   int
	ReservedStock int
}
func (r *ProductRepository) GetProductsWithInventoryStatus(ctx context.Context, merchantID string, limit, offset int) ([]ProductWithInventoryStatus, error) {
	
    
    var results []ProductWithInventoryStatus
    
    query := `
        SELECT p.*,
               CASE WHEN COALESCE(inv.available, 0) > 0 THEN true ELSE false END as in_stock,
               COALESCE(inv.total, 0) as total_stock,
               COALESCE(inv.reserved, 0) as reserved_stock
        FROM products p
        LEFT JOIN (
            SELECT 
                COALESCE(product_id, (SELECT product_id FROM variants WHERE id = variant_id)) as product_id,
                SUM(quantity - reserved_quantity) as available,
                SUM(quantity) as total,
                SUM(reserved_quantity) as reserved
            FROM inventories
            WHERE merchant_id = ?
            GROUP BY COALESCE(product_id, (SELECT product_id FROM variants WHERE id = variant_id))
        ) inv ON inv.product_id = p.id
        WHERE p.merchant_id = ? AND p.deleted_at IS NULL
        ORDER BY p.created_at DESC
        LIMIT ? OFFSET ?
    `
    
    err := r.db.WithContext(ctx).Raw(query, merchantID, merchantID, limit, offset).Scan(&results).Error
    
    return results, err
}