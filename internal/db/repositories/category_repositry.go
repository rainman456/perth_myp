package repositories

import (
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{db: db.DB}
}

// Create adds a new category
func (r *CategoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

// FindByID retrieves a category by ID with parent category
func (r *CategoryRepository) FindByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.Preload("Parent").First(&category, id).Error
	return &category, err
}

// FindAll retrieves all categories
func (r *CategoryRepository) FindAll() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Preload("Parent").Find(&categories).Error
	return categories, err
}

// Update modifies an existing category
func (r *CategoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

// Delete removes a category by ID
func (r *CategoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}


func (r *CategoryRepository) GetAllProductsWithCategorySlug(ctx context.Context, limit, offset int, categorySlug string, preloads ...string) ([]models.Product, int64, error) {
	if categorySlug == "" {
		return nil, 0, fmt.Errorf("category slug is required")
	}

	// Step 1: Fetch the Category by slug to get its ID
	var category models.Category
	err := r.db.WithContext(ctx).
		Where("category_slug = ? AND deleted_at IS NULL", categorySlug).
		First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, fmt.Errorf("category not found for slug: %s", categorySlug)
		}
		return nil, 0, fmt.Errorf("failed to fetch category: %w", err)
	}

	// Step 2: Count total products in this category
	var total int64
	err = r.db.WithContext(ctx).
		Model(&models.Product{}).
		Where("category_id = ? AND deleted_at IS NULL", category.ID).
		Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	// Step 3: Fetch products for this category with preloads, limit/offset
	var products []models.Product
	query := r.db.WithContext(ctx).
		Where("category_id = ? AND deleted_at IS NULL", category.ID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset)

	// Apply preloads (e.g., "Media", "Variants", etc.)
	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	err = query.Find(&products).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch products: %w", err)
	}

	return products, total, nil
}
