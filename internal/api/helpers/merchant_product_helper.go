package helpers

import (
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/api/dto"

	//"github.com/shopspring/decimal"
)

// ToMerchantProductResponse converts model to merchant-specific DTO
func ToMerchantProductResponse(p *models.Product) *dto.MerchantProductResponse {
	resp := &dto.MerchantProductResponse{
		ID:          p.ID,
		MerchantID:  p.MerchantID,
		Name:        p.Name,
		Description: p.Description,
		BasePrice:   p.BasePrice.InexactFloat64(),
		Discount:    p.Discount.InexactFloat64(),
		DiscountType: string(p.DiscountType),
		FinalPrice:  p.FinalPrice.InexactFloat64(),
		CategoryID:  p.CategoryID,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}

	// Variants
	resp.Variants = make([]dto.ProductVariantResponse, len(p.Variants))
	for i, v := range p.Variants {
		resp.Variants[i] = *ToProductVariantResponse(&v)
	}

	// Images (full MediaResponse)
	resp.Images = make([]dto.MediaResponse, len(p.Media))
	for i, m := range p.Media {
		resp.Images[i] = *ToMediaResponse(&m)
	}

	// Reviews
	resp.Reviews = make([]dto.ReviewResponseDTO, len(p.Reviews))
	for i, r := range p.Reviews {
		resp.Reviews[i] = *ToMerchantReviewResponse(&r, p.Name)
	}

	// SimpleInventory
	if p.SimpleInventory != nil {
		resp.SimpleInventory = ToInventoryResponse(p.SimpleInventory)
	}

	return resp
}

// ToProductVariantResponse converts variant model to merchant-specific DTO (no flattened attributes, direct pricing fields)
func ToProductVariantResponse(v *models.Variant) *dto.ProductVariantResponse {
	resp := &dto.ProductVariantResponse{
		ID:              v.ID,
		ProductID:       v.ProductID,
		SKU:             v.SKU,
		PriceAdjustment: v.PriceAdjustment.InexactFloat64(),
		TotalPrice:      v.TotalPrice.InexactFloat64(),
		Discount:        v.Discount.InexactFloat64(),
		DiscountType:    string(v.DiscountType),
		FinalPrice:      v.FinalPrice.InexactFloat64(),
		Attributes:      v.Attributes,
		IsActive:        v.IsActive,
		CreatedAt:       v.CreatedAt,
		UpdatedAt:       v.UpdatedAt,
	}

	// Inventory (check if set)
	if v.Inventory.ID != "" {
		resp.Inventory = *ToInventoryResponse(&v.Inventory)
	}

	return resp
}

// ToMediaResponse converts media model to DTO
func ToMediaResponse(m *models.Media) *dto.MediaResponse {
	return &dto.MediaResponse{
		ID:        m.ID,
		ProductID: m.ProductID,
		URL:       m.URL,
		Type:      string(m.Type),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// ToMerchantReviewResponse converts review model to merchant-specific DTO (includes product details)
func ToMerchantReviewResponse(r *models.Review, productName string) *dto.ReviewResponseDTO {
	return &dto.ReviewResponseDTO{
		ProductName: productName,
		Rating:      r.Rating,
		Comment:     r.Comment,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
		UserName:    r.User.Name,
	}
}