package dto

import (
	"api-customer-merchant/internal/db/models" // For enums if needed
	"time"
)

// AddItemRequest: For POST /cart/add (add new item or increment existing)
type AddItemRequest struct {
	//UserID    uint    `json:"user_id,omitempty"`
	ProductID string    `json:"product_id" validate:"required"`
	VariantID *string `json:"variant_id,omitempty" validate:"omitempty,uuid"`
	Quantity  int     `json:"quantity" validate:"required,gt=0"`
}

// UpdateItemRequest: For PUT /cart/items/:id (full replace quantity) or PATCH /cart/items/:id (partial update)
type UpdateItemRequest struct {
	UserID   uint `json:"user_id,omitempty"`
	Quantity int  `json:"quantity" validate:"omitempty,gt=0"` // Omitempty for PATCH
}

// BulkUpdateRequest: For POST /cart/bulk (add/update multiple—extension for prod)
type BulkUpdateRequest struct {
	UserID uint
	Items  []struct {
		ProductID string    `json:"product_id" validate:"required"`
		VariantID *string `json:"variant_id,omitempty"`
		Quantity  int     `json:"quantity" validate:"required,gt=0"`
	} `json:"items" validate:"dive"`
}

// CartResponse: For all responses (shared output DTO)
type CartResponse struct {
	ID        uint             `json:"id"`
	UserID    uint               `json:"user_id,omitempty"`
	Status    models.CartStatus  `json:"status"`
	Items     []CartItemResponse `json:"items"`
	Total     float64            `json:"total,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty"` // Added
	UpdatedAt time.Time          `json:"updated_at,omitempty"` // Added
}

type CartItemResponse struct {
	ID         uint            `json:"id"`
	ProductID  string            `json:"product_id"`
	VariantID  *string           `json:"variant_id,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
	Product    *ProductResponse  `json:"product,omitempty"`
	Variant    *VariantResponse  `json:"variant,omitempty"` // Embed for display
	Quantity   int               `json:"quantity"`
	Subtotal   float64           `json:"subtotal"`
}
