// services/review/review_service.go (edited)
package review

import (
	"context"
	"errors"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrInvalidReview    = errors.New("invalid review data")
	ErrUnauthorized     = errors.New("unauthorized operation")
	ErrReviewNotFound   = errors.New("review not found")
	ErrNotPurchased     = errors.New("you must purchase the product before reviewing it")
)

type ReviewService struct {
	repo      *repositories.ReviewRepository
	orderRepo *repositories.OrderRepository // Added for purchase check
	logger    *zap.Logger
	validator *validator.Validate
}

func NewReviewService(repo *repositories.ReviewRepository, orderRepo *repositories.OrderRepository, logger *zap.Logger) *ReviewService {
	return &ReviewService{
		repo:      repo,
		orderRepo: orderRepo, // Injected
		logger:    logger,
		validator: validator.New(),
	}
}

func (s *ReviewService) CreateReview(ctx context.Context, userID uint, input dto.CreateReviewDTO) (*dto.ReviewResponseDTO, error) {
	if err := s.validator.Struct(input); err != nil {
		s.logger.Error("Validation failed", zap.Error(err))
		return nil, ErrInvalidReview
	}

	// Added purchase check
	hasPurchased, err := s.orderRepo.HasUserPurchasedProduct(ctx, userID, input.ProductID)
	if err != nil {
		s.logger.Error("Failed to check purchase", zap.Error(err))
		return nil, err
	}
	if !hasPurchased {
		return nil, ErrNotPurchased
	}

	review := &models.Review{
		ProductID: input.ProductID,
		UserID:    userID,
		Rating:    input.Rating,
		Comment:   input.Comment,
	}

	if err := s.repo.Create(ctx, review); err != nil {
		s.logger.Error("Failed to create review", zap.Error(err))
		return nil, err
	}

	return s.mapToDTO(review), nil
}

func (s *ReviewService) GetReview(ctx context.Context, id uint) (*dto.ReviewResponseDTO, error) {
	review, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReviewNotFound
		}
		s.logger.Error("Failed to get review", zap.Error(err))
		return nil, err
	}
	return s.mapToDTO(review), nil
}

func (s *ReviewService) GetReviewsByProduct(ctx context.Context, productID string, limit, offset int) ([]dto.ReviewResponseDTO, error) {
	reviews, err := s.repo.FindByProductID(ctx, productID)
	if err != nil {
		s.logger.Error("Failed to get reviews by product", zap.Error(err))
		return nil, err
	}
	dtos := make([]dto.ReviewResponseDTO, len(reviews))
	for i, r := range reviews {
		dtos[i] = *s.mapToDTO(&r)
	}
	return dtos, nil
}

func (s *ReviewService) GetReviewsByUser(ctx context.Context, userID uint) ([]dto.ReviewResponseDTO, error) {
	reviews, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get reviews by user", zap.Error(err))
		return nil, err
	}
	dtos := make([]dto.ReviewResponseDTO, len(reviews))
	for i, r := range reviews {
		dtos[i] = *s.mapToDTO(&r)
	}
	return dtos, nil
}

func (s *ReviewService) UpdateReview(ctx context.Context, id uint, userID uint, input dto.UpdateReviewDTO) (*dto.ReviewResponseDTO, error) {
	review, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReviewNotFound
		}
		s.logger.Error("Failed to get review for update", zap.Error(err))
		return nil, err
	}
	if review.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Optional: Re-check purchase if needed for updates, but typically not required if create already enforced
	// hasPurchased, err := s.orderRepo.HasUserPurchasedProduct(ctx, userID, review.ProductID)
	// if err != nil || !hasPurchased { ... }

	if input.Rating != nil {
		review.Rating = *input.Rating
	}
	if input.Comment != nil {
		review.Comment = *input.Comment
	}
	if err := s.repo.Update(ctx, review); err != nil {
		s.logger.Error("Failed to update review", zap.Error(err))
		return nil, err
	}
	return s.mapToDTO(review), nil
}

func (s *ReviewService) DeleteReview(ctx context.Context, id uint, userID uint) error {
	review, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrReviewNotFound
		}
		s.logger.Error("Failed to get review for delete", zap.Error(err))
		return err
	}
	if review.UserID != userID {
		return ErrUnauthorized
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete review", zap.Error(err))
		return err
	}
	return nil
}

func (s *ReviewService) mapToDTO(r *models.Review) *dto.ReviewResponseDTO {
	primaryImage:=""
	for _, media := range r.Product.Media {
		if media.Type == models.MediaTypeImage {
			primaryImage = media.URL
			break // Only first image
		}
	}
	return &dto.ReviewResponseDTO{

		

		ProductName:  r.Product.Name,
		ProductID: r.ProductID,
		Rating:      r.Rating,
		Comment:     r.Comment,
		Image: primaryImage,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
		UserName:    r.User.Name,
	}

	//primaryImage := ""
	
	
}