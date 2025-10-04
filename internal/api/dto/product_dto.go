package dto

import (
	"time"
	//"github.com/shopspring/decimal"
)

// ProductInput represents the request body for creating a product
type ProductInput struct {
	Name         string         `json:"name" validate:"required,max=255"`
	MerchantID   string         `json:"merchant_id" validate:"required,max=255"`
	Description  string         `json:"description" validate:"max=1000"`
	//SKU          string         `json:"sku" validate:"required,max=100"`
	BasePrice    float64        `json:"base_price" validate:"required,gt=0"`
	CategoryID   uint           `json:"category_id" validate:"required"`
	InitialStock *int           `json:"initial_stock" validate:"omitempty,gte=0"` // For simple products
	Variants     []VariantInput `json:"variants,omitempty" validate:"dive,omitempty"`
	Media        []MediaInput   `json:"media,omitempty" validate:"dive,omitempty"`
}

type VariantInput struct {
	//SKU             string            `json:"sku" validate:"required,max=100"`
	PriceAdjustment float64           `json:"price_adjustment" validate:"gte=0"`
	Attributes      map[string]string `json:"attributes" validate:"required,dive,required"`
	InitialStock    int               `json:"initial_stock" validate:"gte=0"`
}

type MediaInput struct {
	URL  string `json:"url" validate:"required,url,max=500"`
	Type string `json:"type" validate:"required,oneof=image video"`
}

// ProductResponse for API output
type ProductResponse struct {
	ID              string             `json:"id"`
	MerchantID      string             `json:"merchant_id"`
	Name            string             `json:"name"`
	Description     string             `json:"description"`
	SKU             string             `json:"sku"`
	BasePrice       float64            `json:"base_price"`
	CategoryID      uint               `json:"category_id"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
	Variants        []VariantResponse  `json:"variants,omitempty"`
	Media           []MediaResponse    `json:"media,omitempty"`
	SimpleInventory *InventoryResponse `json:"simple_inventory,omitempty"` // For simple products
}

type VariantResponse struct {
	ID              string            `json:"id"`
	ProductID       string            `json:"product_id"`
	SKU             string            `json:"sku"`
	PriceAdjustment float64           `json:"price_adjustment"`
	TotalPrice      float64           `json:"total_price"`
	Attributes      map[string]string `json:"attributes"`
	IsActive        bool              `json:"is_active"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	Inventory       InventoryResponse `json:"inventory"`
}

type MediaResponse struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	URL       string    `json:"url"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type InventoryResponse struct {
	ID                string `json:"id"`
	Quantity          int    `json:"quantity"`
	ReservedQuantity  int    `json:"reserved_quantity"`
	LowStockThreshold int    `json:"low_stock_threshold"`
	BackorderAllowed  bool   `json:"backorder_allowed"`
}


type MediaUploadRequest struct {
	File  string `form:"file" validate:"required"` // Multipart file key
	Type  string `form:"type" validate:"required,oneof=image video"`
	// Optional: Position or other metadata
}

// MediaUploadResponse
type MediaUploadResponse struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	URL       string    `json:"url"` // Secure Cloudinary URL
	PublicID  string    `json:"public_id"` // For delete/update
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// MediaUpdateRequest for PUT /merchant/products/:media_id
type MediaUpdateRequest struct {
	File *string `form:"file" validate:"omitempty"` // Optional new file
	URL  *string `form:"url" validate:"omitempty,url"` // Or new URL (e.g., external)
	Type *string `form:"type" validate:"omitempty,oneof=image video"`
}

// MediaDeleteRequest (empty body, if needed for reason)
type MediaDeleteRequest struct {
	Reason string `json:"reason" validate:"omitempty,max=500"`
}

type CategoryResponse struct {
	ID         uint                  `json:"id"`
	Name       string                `json:"name"`
	ParentID   *uint                 `json:"parent_id"`
	Attributes map[string]interface{}`json:"attributes"`
	Parent     *CategoryResponse     `json:"parent"`
}