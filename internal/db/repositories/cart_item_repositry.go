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

// FindByID retrieves a cart item by ID with associated Cart, Product, and Merchant
func (r *CartItemRepository) FindByID(id uint) (*models.CartItem, error) {
	var cartItem models.CartItem
	err := r.db.Preload("Cart.User").Preload("Product.Merchant").Preload("Merchant").First(&cartItem, id).Error
	return &cartItem, err
}

// FindByCartID retrieves all cart items for a cart
func (r *CartItemRepository) FindByCartID(cartID uint) ([]models.CartItem, error) {
	var cartItems []models.CartItem
	err := r.db.Preload("Product.Merchant").Preload("Merchant").Where("cart_id = ?", cartID).Find(&cartItems).Error
	return cartItems, err
}

// FindByProductIDAndCartID retrieves a cart item by product ID and cart ID
func (r *CartItemRepository) FindByProductIDAndCartID(productID, cartID uint) (*models.CartItem, error) {
	var cartItem models.CartItem
	err := r.db.Preload("Product.Merchant").Preload("Merchant").Where("product_id = ? AND cart_id = ?", productID, cartID).First(&cartItem).Error
	return &cartItem, err
}

// UpdateQuantity updates the quantity of a cart item
func (r *CartItemRepository) UpdateQuantity(id uint, quantity int) error {
	return r.db.Model(&models.CartItem{}).Where("id = ?", id).Update("quantity", quantity).Error
}

// Update modifies an existing cart item
func (r *CartItemRepository) Update(cartItem *models.CartItem) error {
	return r.db.Save(cartItem).Error
}

// Delete removes a cart item by ID
func (r *CartItemRepository) Delete(id uint) error {
	return r.db.Delete(&models.CartItem{}, id).Error
}