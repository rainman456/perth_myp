package product

import (
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"
	"api-customer-merchant/internal/api/dto"
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