package repositories

import (
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"

	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository() *CartRepository {
	return &CartRepository{db: db.DB}
}

// Create adds a new cart
func (r *CartRepository) Create(cart *models.Cart) error {
	return r.db.Create(cart).Error
}

// FindByID retrieves a cart by ID with associated User and CartItems
func (r *CartRepository) FindByID(id uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.Preload("User").Preload("CartItems.Product.Merchant").First(&cart, id).Error
	return &cart, err
}

// FindActiveCart retrieves the user's most recent active cart
func (r *CartRepository) FindActiveCart(userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.Preload("CartItems.Product.Merchant").Where("user_id = ? AND status = ?", userID, models.CartStatusActive).Order("created_at DESC").First(&cart).Error
	return &cart, err
}

// FindByUserIDAndStatus retrieves carts for a user by status
func (r *CartRepository) FindByUserIDAndStatus(userID uint, status models.CartStatus) ([]models.Cart, error) {
	var carts []models.Cart
	err := r.db.Preload("CartItems.Product.Merchant").Where("user_id = ? AND status = ?", userID, status).Find(&carts).Error
	return carts, err
}

// FindByUserID retrieves all carts for a user
func (r *CartRepository) FindByUserID(userID uint) ([]models.Cart, error) {
	var carts []models.Cart
	err := r.db.Preload("CartItems.Product.Merchant").Where("user_id = ?", userID).Find(&carts).Error
	return carts, err
}

// Update modifies an existing cart
func (r *CartRepository) Update(cart *models.Cart) error {
	return r.db.Save(cart).Error
}

// Delete removes a cart by ID
func (r *CartRepository) Delete(id uint) error {
	return r.db.Delete(&models.Cart{}, id).Error
}