package repositories

import (
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	"context"

	"gorm.io/gorm"
)

type OrderItemRepository struct {
	db *gorm.DB
}

func NewOrderItemRepository() *OrderItemRepository {
	return &OrderItemRepository{db: db.DB}
}

// Create adds a new order item
func (r *OrderItemRepository) Create(orderItem *models.OrderItem) error {
	return r.db.Create(orderItem).Error
}

// FindByID retrieves an order item by ID with associated Order, Product, and Merchant
func (r *OrderItemRepository) FindByID(id uint) (*models.OrderItem, error) {
	var orderItem models.OrderItem
	err := r.db.Preload("Order.User").
	Preload("Product.Merchant").
	Preload("Product.Media").
	Preload("Merchant").First(&orderItem, id).Error
	return &orderItem, err
}

// FindByOrderID retrieves all order items for an order
func (r *OrderItemRepository) FindByOrderID(orderID uint) ([]models.OrderItem, error) {
	var orderItems []models.OrderItem
	err := r.db.Preload("Product.Merchant").Preload("Merchant").Where("order_id = ?", orderID).Find(&orderItems).Error
	return orderItems, err
}

// Update modifies an existing order item
func (r *OrderItemRepository) Update(orderItem *models.OrderItem) error {
	return r.db.Save(orderItem).Error
}

// Delete removes an order item by ID
func (r *OrderItemRepository) Delete(id uint) error {
	return r.db.Delete(&models.OrderItem{}, id).Error
}



// In orderItemRepository
func (r *OrderItemRepository) FindOrderItemsByOrderID(ctx context.Context, orderID uint) ([]models.OrderItem, error) {
	var items []models.OrderItem
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).Find(&items).Error
	return items, err
}