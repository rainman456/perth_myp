package product

import (
	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"
	"api-customer-merchant/internal/api/helpers"
	"context"
	"fmt"
)

type CategoryService struct {
	categoryRepo *repositories.CategoryRepository
}

func NewCategoryService(categoryRepo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

func mapToCategoryDTO(cat *models.Category) *dto.CategoryResponse {
	if cat == nil {
		return nil
	}
	return &dto.CategoryResponse{
		ID:         cat.ID,
		Name:       cat.Name,
		ParentID:   cat.ParentID,
		Attributes: cat.Attributes,
		CategorySlug: cat.CategorySlug,
		Parent:     mapToCategoryDTO(cat.Parent),
	}
}

func (s *CategoryService) GetAllCategories() ([]dto.CategoryResponse, error) {
	cats, err := s.categoryRepo.FindAll()
	if err != nil {
		return nil, err
	}
	var dtos []dto.CategoryResponse
	for _, cat := range cats {
		dtos = append(dtos, *mapToCategoryDTO(&cat))
	}
	return dtos, nil
}










func (s *CategoryService) GetAllProductsWithCategorySlug(ctx context.Context, limit, offset int, categorySlug string) ([]dto.ProductResponse, int64, error) {
	//logger := s.logger.With(zap.String("operation", "GetAllProductsWithCategorySlug"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	preloads := []string{"Media", "Merchant" ,"Variants", "Variants.Inventory", "SimpleInventory", "Category", "Reviews", }

	products, total, err := s.categoryRepo.GetAllProductsWithCategorySlug(ctx, limit, offset, categorySlug, preloads... )  // Fixed: Added ctx (resolves type shifts)
	if err != nil {
		//logger.Error("Failed to fetch all products", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to fetch products: %w", err)
	}

	responses := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		// Prepare variants DTOs
		variantDTOs := make([]dto.VariantResponse, len(p.Variants))
		for j, v := range p.Variants {
			variantDTOs[j] = *helpers.ToVariantResponse(&v, p.BasePrice)
		}

		// Prepare reviews DTOs
		reviewDTOs := make([]dto.ReviewResponseDTO, len(p.Reviews))
		for j, r := range p.Reviews {
			reviewDTOs[j] = *helpers.ToReviewResponse(&r)
		}

		// Use helper (nil merchant for customer-facing, and set MerchantID = "")
		resp := helpers.ToProductResponse(&p, variantDTOs, reviewDTOs, &p.Merchant)
		resp.MerchantID = ""
		responses[i] = *resp
	}

	//logger.Info("Products fetched for landing page", zap.Int("count", len(responses)), zap.Int64("total", total))
	return responses, total, nil
}

