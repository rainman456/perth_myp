// services/wishlist/wishlist_service.go (minor fix: update return in GetWishlist to match []models.Product)
package wishlist

import (
	"context"

	"api-customer-merchant/internal/api/dto"
	//"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"

	"go.uber.org/zap"
)

type WishlistService struct {
	repo   *repositories.WishlistRepository
	logger *zap.Logger
}

func NewWishlistService(repo *repositories.WishlistRepository, logger *zap.Logger) *WishlistService {
	return &WishlistService{
		repo:   repo,
		logger: logger,
	}
}

func (s *WishlistService) AddToWishlist(ctx context.Context, userID uint, productID string) error {
	if err := s.repo.AddToWishlist(ctx, userID, productID); err != nil {
		s.logger.Error("Failed to add to wishlist", zap.Uint("user_id", userID), zap.String("product_id", productID), zap.Error(err))
		return err
	}
	return nil
}

func (s *WishlistService) RemoveFromWishlist(ctx context.Context, userID uint, productID string) error {
	if err := s.repo.RemoveFromWishlist(ctx, userID, productID); err != nil {
		s.logger.Error("Failed to remove from wishlist", zap.Uint("user_id", userID), zap.String("product_id", productID), zap.Error(err))
		return err
	}
	return nil
}

// func (s *WishlistService) GetWishlist(ctx context.Context, userID uint) ([]models.Product, error) {
// 	products, err := s.repo.GetWishlist(ctx, userID)
// 	if err != nil {
// 		s.logger.Error("Failed to get wishlist", zap.Uint("user_id", userID), zap.Error(err))
// 		return nil, err
// 	}
// 	return products, nil
// }


func (s *WishlistService) GetWishlist(ctx context.Context, userID uint) (*dto.WishlistResponseDTO, error) {
	wishlists, err := s.repo.GetWishlist(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get wishlist", zap.Uint("user_id", userID), zap.Error(err))
		return nil, err
	}

	items := make([]dto.WishlistItemResponseDTO, len(wishlists))
	
	for i, w := range wishlists {
		primaryImage := ""
	for _, media := range w.Product.Media {
		if media.Type == models.MediaTypeImage {
			primaryImage = media.URL
			break // Only first image
		}
	}
		items[i] = dto.WishlistItemResponseDTO{
			
			ProductID:  w.ProductID,
			Name:       w.Product.Name,
			//SKU: w.Product.SKU,
			CategorySlug: w.Product.Category.CategorySlug,
			PrimaryImage: primaryImage,
			Discount: w.Product.Discount.InexactFloat64(),
			FinalPrice: w.Product.FinalPrice.InexactFloat64(),

			//DiscountType: string(w.Product.DiscountType),
		}
	}

	return &dto.WishlistResponseDTO{
		UserID: userID,
		Items:  items,
	}, nil
}


func (s *WishlistService) IsInWishlist(ctx context.Context, userID uint, productID string) (bool, error) {
	isIn, err := s.repo.IsInWishlist(ctx, userID, productID)
	if err != nil {
		s.logger.Error("Failed to check if in wishlist", zap.Uint("user_id", userID), zap.String("product_id", productID), zap.Error(err))
		return false, err
	}
	return isIn, nil
}

func (s *WishlistService) ClearWishlist(ctx context.Context, userID uint) error {
	if err := s.repo.ClearWishlist(ctx, userID); err != nil {
		s.logger.Error("Failed to clear wishlist", zap.Uint("user_id", userID), zap.Error(err))
		return err
	}
	return nil
}