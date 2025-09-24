package repositories

import (
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	"fmt"

	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause" // Added for Locking
)

var (
	ErrCartItemNotFound  = errors.New("cart item not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrReservationFailed = errors.New("failed to reserve stock")
)

type CartItemRepository struct {
	db *gorm.DB
}

func NewCartItemRepository() *CartItemRepository {
	return &CartItemRepository{db: db.DB}
}

func (r *CartItemRepository) Create(ctx context.Context, cartItem *models.CartItem) error {
	return r.db.WithContext(ctx).Create(cartItem).Error
}

func (r *CartItemRepository) FindByID(ctx context.Context, id uint) (*models.CartItem, error) {
	var cartItem models.CartItem
	err := r.db.WithContext(ctx).
		Preload("Cart.User").
		Preload("Product.Merchant").
		Preload("Merchant").
		First(&cartItem, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrCartItemNotFound
	}
	return &cartItem, err
}

func (r *CartItemRepository) FindByCartID(ctx context.Context, cartID uint) ([]models.CartItem, error) {
	var cartItems []models.CartItem
	err := r.db.WithContext(ctx).
		Preload("Product.Merchant").
		Preload("Merchant").
		Where("cart_id = ?", cartID).Find(&cartItems).Error
	return cartItems, err
}



func (r *CartItemRepository) FindByProductIDAndCartID(ctx context.Context, productID string, variantID *string, cartID uint) (*models.CartItem, error) {
	var item models.CartItem
	query := r.db.WithContext(ctx).Where("product_id = ? AND cart_id = ?", productID, cartID)
	if variantID != nil {
		query = query.Where("variant_id = ?", *variantID)
	}
	err := query.First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &item, err
}



func (r *CartItemRepository) UpdateQuantityWithReservation(ctx context.Context, itemID uint, newQuantity int, vendorInvID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var item models.CartItem
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&item, itemID).Error; err != nil {
			return fmt.Errorf("failed to lock item: %w", err)
		}

		delta := newQuantity - item.Quantity
		if delta > 0 {
			// Reserve extra stock
			if err := tx.Model(&models.Inventory{}).Where("id = ?", vendorInvID).
				Update("reserved_quantity", gorm.Expr("reserved_quantity + ?", delta)).Error; err != nil {
				return fmt.Errorf("stock reservation failed: %w", err)
			}
		} else if delta < 0 {
			// Unreserve stock
			if err := tx.Model(&models.Inventory{}).Where("id = ?", vendorInvID).
				Update("reserved_quantity", gorm.Expr("reserved_quantity - ?", -delta)).Error; err != nil {
				return fmt.Errorf("stock unreservation failed: %w", err)
			}
		}

		return tx.Model(&models.CartItem{}).Where("id = ?", itemID).Update("quantity", newQuantity).Error
	})
}

func (r *CartItemRepository) Update(ctx context.Context, cartItem *models.CartItem) error {
	return r.db.WithContext(ctx).Save(cartItem).Error
}



func (r *CartItemRepository) DeleteWithUnreserve(ctx context.Context, id uint, vendorInvID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var item models.CartItem
		if err := tx.First(&item, id).Error; err != nil {
			return ErrCartItemNotFound
		}
		// Release reserved stock
		if err := tx.Model(&models.Inventory{}).Where("id = ?", vendorInvID).
			Update("reserved_quantity", gorm.Expr("reserved_quantity - ?", item.Quantity)).Error; err != nil {
			return fmt.Errorf("stock unreservation failed: %w", err)
		}
		return tx.Delete(&models.CartItem{}, id).Error
	})
}

func (r *CartItemRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.CartItem{}, id).Error
}
