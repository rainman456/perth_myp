package helpers

import (
	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/db/models"
	"fmt"
	//"github.com/shopspring/decimal"
)

// ToMerchantProductResponse converts model to merchant-specific DTO
func ToOrderResponse(p *models.Order) *dto.OrderResponse {
	resp := &dto.OrderResponse{
		ID:         p.ID,
        UserID:     p.UserID,
        Status:     dto.OrderStatus(p.Status),
		TotalAmount:   p.TotalAmount.InexactFloat64(),
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}

	
	resp.OrderItems= make([]dto.OrderItemResponse, len(p.OrderItems))
	for i,v:= range p.OrderItems{
		resp.OrderItems[i]=*ToOrderItemResponse(&v)
	}

	return resp
}

// ToProductVariantResponse converts variant model to merchant-specific DTO (no flattened attributes, direct pricing fields)
func ToOrderItemResponse(v *models.OrderItem) *dto.OrderItemResponse {
	resp := &dto.OrderItemResponse{
		ProductID: fmt.Sprint(v.ProductID),
		Name: v.Product.Name,
		Quantity:  v.Quantity,
		Price:     v.Price,
		Image:  "",
	}
	if len(v.Product.Media) > 0 {
		resp.Image = v.Product.Media[0].URL // Assume Media has URL field
	}

	// Inventory (check if set)
	

	return resp
}