package helpers

import (
	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/db/models"
	"fmt"

	//"fmt"
	"math"

	"github.com/shopspring/decimal"
)

// ToCartResponse converts a Cart model to CartResponse DTO
// func ToCartResponse(cart *models.Cart) *dto.CartResponse {
// 	// Prepare items
// 	items := make([]dto.CartItemResponse, len(cart.CartItems))
// 	total := 0.0
// 	for i, item := range cart.CartItems {
// 		// Compute subtotal using Product and Variant prices
// 		var price decimal.Decimal
// 		if item.Product.ID != "" { // Ensure Product is preloaded
// 			price = item.Product.FinalPrice
// 			if item.VariantID != nil && item.Variant != nil && item.Variant.ID != "" {
// 				price = item.Variant.FinalPrice // Use variant's final price if applicable
// 			}
// 		}
// 		subtotal := float64(item.Quantity) * price.InexactFloat64()

// 		// Create item response with embedded ProductResponse
// 		items[i] = dto.CartItemResponse{
// 			ID:         item.ID,
// 			ProductID:  item.ProductID,
// 			Name:       item.Product.Name, // From preloaded Product
// 			VariantID:  item.VariantID,
// 			Quantity:   item.Quantity,
// 			Subtotal:   subtotal,
// 			// Embed full ProductResponse (includes category_slug, etc.)
// 			Product: ToProductResponse(
// 				&item.Product,
// 				nil, // Variants not needed here (already in VariantID/Variant)
// 				nil, // Reviews optional (add if preloaded)
// 				&item.Merchant, // Merchant from preloaded CartItem
// 			),
// 			// Embed VariantResponse if applicable
// 			Variant: func() *dto.VariantResponse {
// 				if item.Variant != nil && item.Variant.ID != "" {
// 					return ToVariantResponse(item.Variant, item.Product.BasePrice)
// 				}
// 				return nil
// 			}(),
// 		}
// 		total += subtotal
// 	}

// 	return &dto.CartResponse{
// 		ID:        cart.ID,
// 		UserID:    cart.UserID,
// 		Status:    cart.Status,
// 		Items:     items,
// 		Total:     math.Round(total*100) / 100, // Round to 2 decimals
// 		CreatedAt: cart.CreatedAt,
// 		UpdatedAt: cart.UpdatedAt,
// 	}
// }



func ToCartResponse(cart *models.Cart) *dto.CartResponse {
	items := make([]dto.CartItemResponse, len(cart.CartItems))
	subtotal := 0.0
	// Stubbed tax/shipping (replace with real logic if available)
	//taxTotal := 0.0     // e.g., 0.1 * subtotal
	//shippingTotal := 0.0 // e.g., flat rate or API call

	for i, item := range cart.CartItems {
		// Compute price (variant if present, else product)
		var price decimal.Decimal
		if item.Product.ID != "" {
			price = item.Product.FinalPrice
			if item.VariantID != nil && item.Variant != nil && item.Variant.ID != "" {
				price = item.Variant.FinalPrice
			}
		}
		itemSubtotal := float64(item.Quantity) * price.InexactFloat64()

		// Build item response
		items[i] = dto.CartItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Name:      item.Product.Name,
			VariantID: item.VariantID,
			Quantity:  item.Quantity,
			Subtotal:  math.Round(itemSubtotal*100) / 100,
			Product: ToCartProductResponse(
				&item.Product,
				nil, // No variants needed
				nil, // No reviews needed
				&item.Merchant,
			),
			Variant: func() *dto.CartVariantResponse {
				if item.Variant != nil && item.Variant.ID != "" {
					available := 0
					if item.Variant.Inventory.ID != "" {
						available = item.Variant.Inventory.Quantity - item.Variant.Inventory.ReservedQuantity
						if available < 0 {
							available = 0
						}
					}
					var color, size *string
					if c, ok := item.Variant.Attributes["color"]; ok {
						color = &c
					}
					if s, ok := item.Variant.Attributes["size"]; ok {
						size = &s
					}
					return &dto.CartVariantResponse{
						ID:              item.Variant.ID,
						ProductID:       item.Variant.ProductID,
						Color:           color,
						Size:            size,
						Pricing: dto.VariantPricingResponse{ // Added full pricing
							BasePrice:       item.Product.BasePrice.InexactFloat64(),
							PriceAdjustment: item.Variant.PriceAdjustment.InexactFloat64(),
							TotalPrice:      item.Variant.TotalPrice.InexactFloat64(),
							Discount:        item.Variant.Discount.InexactFloat64(),
							FinalPrice:      item.Variant.FinalPrice.InexactFloat64(),
						},
						FinalPrice:      item.Variant.FinalPrice.InexactFloat64(),
						Available:       available,
						BackorderAllowed: item.Variant.Inventory.ID != "" && item.Variant.Inventory.BackorderAllowed,
					}
				}
				return nil
			}(),
		}
		subtotal += itemSubtotal
	}

	return &dto.CartResponse{
		ID:            cart.ID,
		UserID:        cart.UserID,
		Status:        cart.Status,
		Items:         items,
		// Subtotal:      math.Round(subtotal*100) / 100,
		// TaxTotal:      taxTotal,
		// ShippingTotal: shippingTotal,
		// GrandTotal:    math.Round((subtotal+taxTotal+shippingTotal)*100) / 100,
		CreatedAt:     cart.CreatedAt,
		UpdatedAt:     cart.UpdatedAt,
	}
}

// ToProductResponse - Slim version for cart
func ToCartProductResponse(
	p *models.Product,
	variants []dto.CartVariantResponse,
	reviews []dto.ReviewResponseDTO,
	merchant *models.Merchant,
) *dto.CartProductResponse {
	// First image as primary
	primaryImage := ""
	for _, media := range p.Media {
		if media.Type == models.MediaTypeImage {
			primaryImage = media.URL
			break // Only first image
		}
	}

	resp := &dto.CartProductResponse{
		ID:           p.ID,
		Name:         p.Name,
		Slug:         p.Slug,
		Pricing: dto.ProductPricingResponse{ // Added full pricing
			BasePrice:    p.BasePrice.InexactFloat64(),
			Discount:     p.Discount.InexactFloat64(),
			FinalPrice:   p.FinalPrice.InexactFloat64(),
		},
		FinalPrice:   p.FinalPrice.InexactFloat64(),
		PrimaryImage: primaryImage,
	}

	// Category (use DB-stored slug)
	if p.CategoryID != 0 && p.Category.ID != 0 {
		resp.CategoryName = p.Category.Name
		resp.CategorySlug = p.Category.CategorySlug
		fmt.Printf("Category set for product %s: Name='%s', Slug='%s'\n", p.ID, resp.CategoryName, resp.CategorySlug)
	} else {
		fmt.Printf("Category missing for product %s (CategoryID=%d, Loaded=%v)\n", p.ID, p.CategoryID, p.Category.ID != 0)
	}

	// Merchant
	if merchant != nil && merchant.ID != "" {
		resp.MerchantName = merchant.Name
		resp.MerchantStoreName = merchant.StoreName
	}

	return resp
}
