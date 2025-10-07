// db/repositories/wishlist_repository.go (fixed to use UserWishlist explicitly)
package repositories

import (
	"context"
	"errors"
	"time"

	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"

	"gorm.io/gorm"
)

type WishlistRepository struct {
	db *gorm.DB
}

func NewWishlistRepository() *WishlistRepository {
	return &WishlistRepository{db: db.DB}
}

func (r *WishlistRepository) AddToWishlist(ctx context.Context, userID uint, productID string) error {
	// Check if user exists
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Check if product exists
	var product models.Product
	if err := r.db.WithContext(ctx).First(&product, "id = ?", productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return err
	}

	// Check if already in wishlist
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.UserWishlist{}).
		Where("user_id = ? AND product_id = ?", userID, productID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("product already in wishlist")
	}

	// Create UserWishlist entry
	wishlist := &models.UserWishlist{
		UserID:    userID,
		ProductID: productID,
		AddedAt:   time.Now(),
	}
	return r.db.WithContext(ctx).Create(wishlist).Error
}

func (r *WishlistRepository) RemoveFromWishlist(ctx context.Context, userID uint, productID string) error {
	// Delete specific entry
	result := r.db.WithContext(ctx).Where("user_id = ? AND product_id = ?", userID, productID).
		Delete(&models.UserWishlist{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("wishlist entry not found")
	}
	return nil
}


func (r *WishlistRepository) GetWishlist(ctx context.Context, userID uint) ([]models.UserWishlist, error) {
	var wishlists []models.UserWishlist
	err := r.db.WithContext(ctx).
		Preload("Product").
		Where("user_id = ?", userID).
		Find(&wishlists).Error
	if err != nil {
		return nil, err
	}

	products := make([]models.Product, len(wishlists))
	for i, w := range wishlists {
		products[i] = w.Product
	}
	return wishlists, nil
}


func (r *WishlistRepository) IsInWishlist(ctx context.Context, userID uint, productID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.UserWishlist{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *WishlistRepository) ClearWishlist(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&models.UserWishlist{}).Error
}