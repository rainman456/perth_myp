package helpers

import (
	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/db/models"
	"fmt"

	"github.com/shopspring/decimal"
)

// ToVariantResponse converts model to DTO with flattened attributes
func ToVariantResponse(v *models.Variant, productBasePrice decimal.Decimal) *dto.VariantResponse {
	resp := &dto.VariantResponse{
		ID:        v.ID,
		ProductID: v.ProductID,
		//SKU:       v.SKU,
		//Attributes: v.Attributes,
		Pricing: dto.VariantPricingResponse{
			BasePrice:       productBasePrice.InexactFloat64(),
			PriceAdjustment: v.PriceAdjustment.InexactFloat64(),
			TotalPrice:      v.TotalPrice.InexactFloat64(),
			Discount:        v.Discount.InexactFloat64(),
			//DiscountType:    string(v.DiscountType),
			FinalPrice:      v.FinalPrice.InexactFloat64(),
		},
		IsActive:  v.IsActive,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}

	// Flatten common attributes
	if v.Attributes != nil {
		if color, ok := v.Attributes["color"]; ok {
			resp.Color = &color
		}
		if size, ok := v.Attributes["size"]; ok {
			resp.Size = &size
		}
		// if material, ok := v.Attributes["material"]; ok {
		// 	resp.Material = &material
		// }
		// if pattern, ok := v.Attributes["pattern"]; ok {
		// 	resp.Pattern = &pattern
		// }
	}

	// Inventory
	if v.Inventory.ID != "" {
		resp.Inventory = *ToInventoryResponse(&v.Inventory)
	}

	return resp
}

func ToProductResponse(
    p *models.Product,
    variants []dto.VariantResponse,
    reviews []dto.ReviewResponseDTO,
    merchant *models.Merchant,
) *dto.ProductResponse {
    imageURLs := []string{}
    for _, media := range p.Media {
        if media.Type == models.MediaTypeImage {
            imageURLs = append(imageURLs, media.URL)
        }
    }

    // Compute average rating from reviews
    var totalRating int
    reviewCount := len(reviews)
    for _, rev := range reviews {
        totalRating += rev.Rating
    }
    var avgRating float64
    if reviewCount > 0 {
        avgRating = float64(totalRating) / float64(reviewCount)
    }

    resp := &dto.ProductResponse{
        ID:               p.ID,
        //SKU:              p.SKU,
        MerchantID:       p.MerchantID,
        Name:             p.Name,
        Description:      p.Description,
    
        Pricing: dto.ProductPricingResponse{
            BasePrice:    p.BasePrice.InexactFloat64(),
            Discount:     p.Discount.InexactFloat64(),
            //DiscountType: string(p.DiscountType),
            FinalPrice:   p.FinalPrice.InexactFloat64(),
        },
        Images:     imageURLs,
        Variants:   variants,
        Slug:       p.Slug,
		CategoryName: p.Category.Name,
        Reviews:    reviews,
        AvgRating:  avgRating,  // Add this field to dto.ProductResponse if not already present
		ReviewCount: reviewCount,
        CreatedAt:  p.CreatedAt,
        UpdatedAt:  p.UpdatedAt,
    }

    // Handle optional merchant
	if merchant != nil {
		resp.MerchantName = merchant.Name
		resp.MerchantStoreName = merchant.StoreName
	} else {
		fmt.Println("Merchant not provided")  // Debug if empty
	}

    // Simple product inventory
    if p.SimpleInventory != nil {
        resp.Inventory = ToInventoryResponse(p.SimpleInventory)
    }
	if p.CategoryID != 0 && p.Category.ID != 0 {
		resp.CategoryName = p.Category.Name
		resp.CategorySlug = p.Category.Slug()  // Calls your Slug() method
		fmt.Printf("Category loaded: Name=%s, Slug=%s\n", resp.CategoryName, resp.CategorySlug)  // Debug print
	} else {
		fmt.Println("Category not loaded or ID=0")  // This will show if issue persists
		resp.CategoryName = ""
		resp.CategorySlug = ""
	}

    return resp
}

func ToInventoryResponse(inv *models.Inventory) *dto.InventoryResponse {
	if inv == nil {
		return nil
	}

	available := inv.Quantity - inv.ReservedQuantity
	if available < 0 {
		available = 0
	}

	return &dto.InventoryResponse{
		ID:                inv.ID,
		Quantity:          inv.Quantity,
		Reserved:          inv.ReservedQuantity,
		Available:         available,
		Status:            string(inv.GetStatus()),
		BackorderAllowed:  inv.BackorderAllowed,
		LowStockThreshold: inv.LowStockThreshold,
	}
}

func ToReviewResponse(r *models.Review) *dto.ReviewResponseDTO {
	resp := &dto.ReviewResponseDTO{
		//ID:        r.ID,
		//UserID:    r.UserID,
		Rating:    r.Rating,
		Comment:   r.Comment,
		CreatedAt: r.CreatedAt,
		
	}

	// Safe access to related data (assume preloaded; add checks if needed)
	if r.Product.Name != "" {
		resp.ProductName = r.Product.Name
	}
	if r.User.Name != "" {
		resp.UserName = r.User.Name
	}

	return resp
}