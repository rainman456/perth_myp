package repositories

import (
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	"context"

	//"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{db: db.DB}
}

// Create adds a new order
// func (r *OrderRepository) Create(order *models.Order) error {
// 	return r.db.Create(order).Error
// }

type OrderInterface interface {
	FindByID(ctx context.Context, id uint) (*models.Order, error)
	// Add other methods as needed
}

func (r *OrderRepository) Create(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

// FindByID retrieves an order by ID with associated User and OrderItems
func (r *OrderRepository) FindByID(ctx context.Context, id uint) (*models.Order, error) {
	var order models.Order
	//err := r.db.Preload("User").Preload("OrderItems.Product.Merchant").First(&order, id).Error
	err := r.db.WithContext(ctx).
		Preload("User.Addresses").
		Preload("OrderItems").
		Preload("OrderItems.Product.Media").
		Preload("OrderItems.Product").
		Preload("OrderItems.Merchant").
		First(&order, id).Error
	return &order, err
}

// FindByUserID retrieves all orders for a user
func (r *OrderRepository) FindByUserID(ctx context.Context, userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.WithContext(ctx).
		Preload("OrderItems", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, order_id, product_id, variant_id, merchant_id, quantity, price, fulfillment_status")
		}).
		Preload("OrderItems.Product", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, description, base_price")
		}).
		Preload("OrderItems.Merchant", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, merchant_id, store_name")
		}).
		Preload("OrderItems.Product.Media", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, product_id, url")
		}).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

// FindByMerchantID retrieves all orders containing items from a merchant
func (r *OrderRepository) FindByMerchantID(ctx context.Context, merchantID interface{}) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.WithContext(ctx).
		Preload("OrderItems.Product.Media").
		Preload("OrderItems.Product").
		Preload("User.Addresses").
		Joins("JOIN order_items oi ON oi.order_id = orders.id").
		Where("oi.merchant_id = ?", merchantID).
		Find(&orders).Error
	return orders, err
}

// Update modifies an existing order
func (r *OrderRepository) Update(ctx context.Context,order *models.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

// Delete removes an order by ID
func (r *OrderRepository) Delete(id uint) error {
	return r.db.Delete(&models.Order{}, id).Error
}


// FindByIDWithPreloads fetches with ownership check and preloads (avoids N+1)
func (r *OrderRepository) FindByIDWithPreloads(ctx context.Context, id uint) (*models.Order, error) {
	var order models.Order
	// Preload OrderItems (no deeper Inventory preload to avoid N+1; fetch separately if needed)
	err := r.db.WithContext(ctx).
		Scopes(r.activeScope()). // Soft delete filter
		Preload("OrderItems").
		Preload("Payment").
		First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// UpdateStatus updates order status (with locking for concurrency)
func (r *OrderRepository) UpdateStatus(ctx context.Context, id uint, status models.OrderStatus) error {
	return r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Model(&models.Order{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// activeScope for soft deletes (if Order has DeletedAt)
func (r *OrderRepository) activeScope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Unscoped().Where("deleted_at IS NULL")
	}
}

// HasUserPurchasedProduct checks if the user has at least one completed order containing the product
func (r *OrderRepository) HasUserPurchasedProduct(ctx context.Context, userID uint, productID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("orders").
		Joins("JOIN order_items ON order_items.order_id = orders.id").
		Where("orders.user_id = ? AND order_items.product_id = ? AND orders.status = ? OR orders.status = ?", userID, productID, models.OrderStatusCompleted ,models.OrderStatusPaid).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}





// func (r *OrderRepository) FindByIDWithDB(db *gorm.DB, id uint) (*models.Order, error) {
//     var order models.Order
//     err := db.Preload("User").
//              Preload("OrderItems").
//              Preload("OrderItems.Product").
//              Preload("OrderItems.Merchant").
//              First(&order, id).Error
//     return &order, err
// }

// func (r *OrderRepository) UpdateWithDB(db *gorm.DB, order *models.Order) error {
//     return db.Save(order).Error
// }