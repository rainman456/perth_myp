// db/repositories/review_repository.go
package repositories

import (
	"context"

	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"

	"gorm.io/gorm"
)

type ReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository() *ReviewRepository {
	return &ReviewRepository{db: db.DB}
}

func (r *ReviewRepository) Create(ctx context.Context, review *models.Review) error {
	return r.db.WithContext(ctx).Create(review).Error
}

func (r *ReviewRepository) FindByID(ctx context.Context, id uint) (*models.Review, error) {
	var review models.Review
	err := r.db.WithContext(ctx).
	Preload("User").
	Preload("Product").
	Preload("Product.Media").
	First(&review, id).Error
	return &review, err
}

func (r *ReviewRepository) FindByProductID(ctx context.Context, productID string) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.WithContext(ctx).Preload("User").
	Preload("Product").
	Preload("Product.Media").
	Where("product_id = ?", productID).Find(&reviews).Error
	return reviews, err
}

func (r *ReviewRepository) FindByUserID(ctx context.Context, userID uint) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.WithContext(ctx).
	Preload("Product").
	Preload("Product").
	Preload("Product.Media").
	Where("user_id = ?", userID).Find(&reviews).Error
	return reviews, err
}




// Update this method in ReviewRepository (db/repositories/review_repository.go)
// Now returns avg and count directly with two efficient queries (avoids extra struct)
func (r *ReviewRepository) GetAverageRatingByProductID(ctx context.Context, productID string) (float64, int, error) {
    var avg float64
    err := r.db.WithContext(ctx).
        Model(&models.Review{}).
        Where("product_id = ?", productID).
        Select("COALESCE(AVG(rating), 0)").
        Scan(&avg).Error
    if err != nil {
        return 0, 0, err
    }

    var count int64
    err = r.db.WithContext(ctx).
        Model(&models.Review{}).
        Where("product_id = ?", productID).
        Count(&count).Error
    if err != nil {
        return avg, 0, err
    }

    return avg, int(count), nil
}

func (r *ReviewRepository) Update(ctx context.Context, review *models.Review) error {
	return r.db.WithContext(ctx).Save(review).Error
}

func (r *ReviewRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Review{}, id).Error
}