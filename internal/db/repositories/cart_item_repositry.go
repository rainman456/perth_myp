package repositories

import (
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"

	"gorm.io/gorm"
)

type CartItemRepository struct {
	db *gorm.DB
}

func NewCartItemRepository() *CartItemRepository {
	return &CartItemRepository{db: db.DB}
}

// Create adds a new cart item
func (r *CartItemRepository) Create(cartItem *models.CartItem) error {
	return r.db.Create(cartItem).Error
}

// FindByCartID retrieves all cart items for a cart
func (r *CartItemRepository) FindByCartID(cartID uint) ([]models.CartItem, error) {
	var cartItems []models.CartItem
	err := r.db.Preload("Product.Merchant").Where("cart_id = ?", cartID).Find(&cartItems).Error
	return cartItems, err
}

// Update modifies an existing cart item
func (r *CartItemRepository) Update(cartItem *models.CartItem) error {
	return r.db.Save(cartItem).Error
}

// Delete removes a cart item by ID
func (r *CartItemRepository) Delete(id uint) error {
	return r.db.Delete(&models.CartItem{}, id).Error
}