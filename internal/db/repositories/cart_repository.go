package repositories

import (
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

var ErrCartNotFound = errors.New("cart not found")

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository() *CartRepository {
	return &CartRepository{db: db.DB}
}

// Create adds a new cart
func (r *CartRepository) Create(ctx context.Context, cart *models.Cart) error {
	return r.db.WithContext(ctx).Create(cart).Error
}

// FindByID retrieves a cart by ID with associated User and CartItems
// func (r *CartRepository) FindByID(id uint) (*models.Cart, error) {
// 	var cart models.Cart
// 	err := r.db.Preload("User").Preload("CartItems.Product.Merchant").First(&cart, id).Error
// 	return &cart, err
// }

func (r *CartRepository) FindByID(ctx context.Context, id uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.WithContext(ctx).
		//Preload("User").
		Preload("CartItems.Product.Category").  // For category_name/slug
		Preload("CartItems.Product.Media").      // For images
		//Preload("CartItems.Product.Merchant").   // For merchant_name/store_name
		//Preload("CartItems.Product.Reviews").    // For avg_rating/review_count (or limit: Preload("CartItems.Product.Reviews", func(db *gorm.DB) *gorm.DB { return db.Limit(10).Order("created_at DESC") }))
		Preload("CartItems.Variant.Inventory"). // Existing + variants for full response
		//Preload("CartItems.Variant").            // For variant details in item
		//Preload("CartItems.Merchant").           // Extra if needed, but Product.Merchant covers
		First(&cart, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrCartNotFound
	}
	return &cart, err
}
// FindActiveCart retrieves the user's most recent active cart
// func (r *CartRepository) FindActiveCart(userID uint) (*models.Cart, error) {
// 	var cart models.Cart
// 	err := r.db.Preload("CartItems.Product.Merchant").Where("user_id = ? AND status = ?", userID, models.CartStatusActive).Order("created_at DESC").First(&cart).Error
// 	return &cart, err
// }

func (r *CartRepository) FindActiveCart(ctx context.Context, userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.WithContext(ctx).
		Preload("CartItems.Product.Category").  // For category_name/slug
		Preload("CartItems.Product.Media").      // For images
		//Preload("CartItems.Product.Merchant").   // For merchant_name/store_name
		//Preload("CartItems.Product.Reviews").    // For avg_rating/review_count (or limit as above)
		Preload("CartItems.Variant.Inventory"). // For variants
		//Preload("CartItems.Variant").            // For variant details
		//Preload("CartItems.Merchant").           // Existing + extra
		Where("user_id = ? AND status = ?", userID, models.CartStatusActive).
		Order("created_at DESC").First(&cart).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound //ErrCartNotFound
	}
	return &cart, err
}

// FindByUserIDAndStatus retrieves carts for a user by status
func (r *CartRepository) FindByUserIDAndStatus(ctx context.Context, userID uint, status models.CartStatus) ([]models.Cart, error) {
	var carts []models.Cart
	err := r.db.WithContext(ctx).
		Preload("CartItems.Product.Merchant").Where("user_id = ? AND status = ?", userID, status).Find(&carts).Error
	return carts, err
}

// FindByUserID retrieves all carts for a user
func (r *CartRepository) FindByUserID(ctx context.Context, userID uint) ([]models.Cart, error) {
	var carts []models.Cart
	err := r.db.WithContext(ctx).Preload("CartItems.Product.Merchant").Where("user_id = ?", userID).Find(&carts).Error
	return carts, err
}

// Update modifies an existing cart
func (r *CartRepository) Update(ctx context.Context, cart *models.Cart) error {
	return r.db.WithContext(ctx).Save(cart).Error
}

// Delete removes a cart by ID
func (r *CartRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Cart{}, id).Error
}

func (r *CartRepository) FindActiveCartLight(ctx context.Context, userID uint) (*models.Cart, error) {
    var cart models.Cart
    err := r.db.WithContext(ctx).
        Where("user_id = ? AND status = ?", userID, models.CartStatusActive).
        Order("created_at DESC").First(&cart).Error
    return &cart, err
}