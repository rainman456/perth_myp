package repositories

import (
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	"context"
	"errors"

	"gorm.io/gorm"
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

func (r *OrderRepository) Create(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

// FindByID retrieves an order by ID with associated User and OrderItems
func (r *OrderRepository) FindByID(ctx context.Context, id uint) (*models.Order, error) {
	var order models.Order
	//err := r.db.Preload("User").Preload("OrderItems.Product.Merchant").First(&order, id).Error
	err := r.db.WithContext(ctx).Preload("User").Preload("OrderItems").Preload("OrderItems.Product").Preload("OrderItems.Merchant").First(&order, id).Error
	return &order, err
}

// FindByUserID retrieves all orders for a user
func (r *OrderRepository) FindByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("OrderItems.Product.Merchant").Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}

// FindByMerchantID retrieves all orders containing items from a merchant
func (r *OrderRepository) FindByMerchantID(merchantID uint) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("OrderItems.Product").Joins("JOIN order_items oi ON oi.order_id = orders.id").
		Where("oi.merchant_id = ?", merchantID).Find(&orders).Error
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


func (r *OrderRepository) CancelOrder(ctx context.Context, orderID uint, userID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var ord models.Order
		if err := tx.Preload("OrderItems").First(&ord, "id = ? AND user_id = ?", orderID, userID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("order not found")
			}
			return err
		}
		if ord.Status != models.OrderStatusPending {
			return fmt.Errorf("only pending orders can be cancelled")
		}

		// Restore inventory for each order item
		for _, oi := range ord.OrderItems {
			// find corresponding inventory row
			var inv models.Inventory
			if oi.VariantID != nil && *oi.VariantID != "" {
				if err := tx.First(&inv, "variant_id = ? AND merchant_id = ?", *oi.VariantID, oi.MerchantID).Error; err != nil {
					// continue on errors? better return; but conservative approach: return error
					return fmt.Errorf("inventory lookup failed: %w", err)
				}
			} else {
				if err := tx.First(&inv, "product_id = ? AND merchant_id = ?", oi.ProductID, oi.MerchantID).Error; err != nil {
					return fmt.Errorf("inventory lookup failed: %w", err)
				}
			}

			inv.Quantity = inv.Quantity + oi.Quantity
			if inv.ReservedQuantity >= oi.Quantity {
				inv.ReservedQuantity = inv.ReservedQuantity - oi.Quantity
			} else {
				inv.ReservedQuantity = 0
			}
			if err := tx.Save(&inv).Error; err != nil {
				return fmt.Errorf("failed to update inventory: %w", err)
			}
		}

		ord.Status = models.OrderStatusCancelled
		if err := tx.Save(&ord).Error; err != nil {
			return fmt.Errorf("failed to update order status: %w", err)
		}
		return nil
	})
}