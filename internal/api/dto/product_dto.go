package dto

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"time"
	//"github.com/shopspring/decimal"
)

// ProductInput represents the request body for creating a product
type ProductInput struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description" validate:"max=1000"`
	//SKU          string         `json:"sku" validate:"required,max=100"`
	BasePrice    float64        `json:"base_price" validate:"required,gt=0"`
	CategoryID   uint           `json:"category_id" validate:"required"`
	CategoryName string         `json:"category_name" validate:"required"`
	InitialStock *int           `json:"initial_stock" validate:"omitempty,gte=0"` // For simple products
	Discount     float64        `json:"discount" validate:"gte=0"`
	DiscountType string         `json:"discount_type" validate:"oneof=fixed percentage ''"`
	Variants     []VariantInput `json:"variants,omitempty" validate:"dive,omitempty"`
	Images       []MediaInput   `json:"media,omitempty" validate:"dive,omitempty"`
}

// UpdateProductInput represents the request body for updating a product
type UpdateProductInput struct {
	Name         *string  `json:"name" validate:"omitempty,max=255"`
	Description  *string  `json:"description" validate:"omitempty,max=1000"`
	BasePrice    *float64 `json:"base_price" validate:"omitempty,gt=0"`
	CategoryID   *uint    `json:"category_id" validate:"omitempty"`
	CategoryName *string  `json:"category_name" validate:"omitempty"`
	Discount     *float64 `json:"discount" validate:"omitempty,gte=0"`
	DiscountType *string  `json:"discount_type" validate:"omitempty,oneof=fixed percentage ''"`
}

// BulkUpdateProductInput represents the request body for bulk updating products
type BulkUpdateProductInput struct {
	ProductID string              `json:"product_id" validate:"required,uuid"`
	Product   *UpdateProductInput `json:"product" validate:"omitempty"`
	Variants  []BulkUpdateVariant `json:"variants" validate:"omitempty,dive"`
}

// BulkUpdateVariant represents a variant update in bulk update
type BulkUpdateVariant struct {
	VariantID string              `json:"variant_id" validate:"required,uuid"`
	Variant   *UpdateVariantInput `json:"variant" validate:"omitempty"`
}

// VariantInput represents the request body for creating a variant
type VariantInput struct {
	//SKU             string            `json:"sku" validate:"required,max=100"`
	PriceAdjustment float64           `json:"price_adjustment" validate:"gte=0"`
	Discount        float64           `json:"discount" validate:"gte=0"`
	DiscountType    string            `json:"discount_type" validate:"oneof=fixed percentage ''"`
	Attributes      map[string]string `json:"attributes" validate:"required,dive,required"`
	InitialStock    int               `json:"initial_stock" validate:"gte=0"`
}

// UpdateVariantInput represents the request body for updating a variant
type UpdateVariantInput struct {
	PriceAdjustment *float64          `json:"price_adjustment" validate:"omitempty,gte=0"`
	Discount        *float64          `json:"discount" validate:"omitempty,gte=0"`
	DiscountType    *string           `json:"discount_type" validate:"omitempty,oneof=fixed percentage ''"`
	Attributes      map[string]string `json:"attributes" validate:"omitempty,dive,required"`
	IsActive        *bool             `json:"is_active" validate:"omitempty"`
}

type MediaInput struct {
	URL  string `json:"url" validate:"required,url,max=500"`
	Type string `json:"type" validate:"required,oneof=image video"`
}

// ProductResponse for API output
type MerchantProductResponse struct {
	ID          string `json:"id"`
	MerchantID  string `json:"merchant_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	//SKU             string             `json:"sku"`
	BasePrice       float64                  `json:"base_price"`
	Discount        float64                  `json:"discount" validate:"gte=0"`
	DiscountType    string                   `json:"discount_type" validate:"oneof=fixed percentage ''"`
	FinalPrice      float64                  `json:"final_price"`
	CategoryID      uint                     `json:"category_id"`
	CategoryName    string                   `json:"category_name" validate:"required"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
	Variants        []ProductVariantResponse `json:"variants,omitempty"`
	Images          []MediaResponse          `json:"media,omitempty"`
	Reviews         []ReviewResponseDTO      `json:"reviews,omitempty"`
	SimpleInventory *InventoryResponse       `json:"simple_inventory,omitempty"` // For simple products
}

type ProductVariantResponse struct {
	ID              string            `json:"id"`
	ProductID       string            `json:"product_id"`
	SKU             string            `json:"sku"`
	PriceAdjustment float64           `json:"price_adjustment"`
	TotalPrice      float64           `json:"total_price"`
	Discount        float64           `json:"discount" validate:"gte=0"`
	DiscountType    string            `json:"discount_type" validate:"oneof=fixed percentage ''"`
	FinalPrice      float64           `json:"final_price"`
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

// type InventoryResponse struct {
// 	ID                string `json:"id"`
// 	Quantity          int    `json:"quantity"`
// 	ReservedQuantity  int    `json:"reserved_quantity"`
// 	LowStockThreshold int    `json:"low_stock_threshold"`
// 	BackorderAllowed  bool   `json:"backorder_allowed"`
// }

type VariantResponse struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	//SKU       string `json:"sku"`

	// Flattened attributes for convenience
	Color *string `json:"color,omitempty"`
	Size  *string `json:"size,omitempty"`
	//Material *string `json:"material,omitempty"`
	//Pattern  *string `json:"pattern,omitempty"`

	Pricing   VariantPricingResponse `json:"pricing"`
	Inventory InventoryResponse      `json:"inventory"`
	IsActive  bool                   `json:"is_active"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// VariantPricingResponse - Variant pricing
type VariantPricingResponse struct {
	BasePrice       float64 `json:"base_price"`       // Product base price
	PriceAdjustment float64 `json:"price_adjustment"` // Variant markup/markdown
	TotalPrice      float64 `json:"total_price"`      // BasePrice + Adjustment
	Discount        float64 `json:"discount"`         // Discount amount or percentage
	//DiscountType    string  `json:"discount_type"`    // "fixed", "percentage", or ""
	FinalPrice float64 `json:"final_price"` // Pre-calculated final price
}

// ProductResponse - Main product response
type ProductResponse struct {
	ID string `json:"id"`
	//SKU             string            `json:"sku"`
	MerchantID        string `json:"merchant_id"`
	MerchantName      string `json:"merchant_name"`
	MerchantStoreName string `json:"merchant_store_name"`

	Name        string `json:"name"`
	Description string `json:"description"`
	Slug        string `json:"slug"`
	//	CategoryID  uint   `json:"category_id"`
	CategorySlug string `json:"category_slug"`
	CategoryName string `json:"category_name"`

	Pricing   ProductPricingResponse `json:"pricing"`
	Inventory *InventoryResponse     `json:"inventory,omitempty"` // nil for variant products

	Reviews  []ReviewResponseDTO `json:"reviews,omitempty"`
	Images   []string            `json:"images"`
	Variants []VariantResponse   `json:"variants,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	AvgRating   float64 `json:"average_rating"`
	ReviewCount int     `json:"review_count"`
}

// ProductPricingResponse - Product pricing
type ProductPricingResponse struct {
	BasePrice float64 `json:"base_price"`
	Discount  float64 `json:"discount"`
	//DiscountType string  `json:"discount_type"` // "fixed", "percentage", or ""
	FinalPrice float64 `json:"final_price"` // Pre-calculated
}

// InventoryResponse - Inventory/stock info
type InventoryResponse struct {
	ID                string `json:"id"`
	Quantity          int    `json:"quantity"`
	Reserved          int    `json:"reserved"`
	Available         int    `json:"available"` // quantity - reserved
	Status            string `json:"status"`    // "in_stock", "low_stock", "out_of_stock", "backorder"
	BackorderAllowed  bool   `json:"backorder_allowed"`
	LowStockThreshold int    `json:"low_stock_threshold"`
}

// type MediaUploadRequest struct {
// 	File string `form:"file" validate:"required"` // Multipart file key
// 	Type string `form:"type" validate:"required,oneof=image video"`
// 	// Optional: Position or other metadata
// }

type MediaUploadRequest struct {
	Type string `form:"type" validate:"required,oneof=image video"`
	// Optional: Position or other metadata
}

// MediaUploadResponse
type MediaUploadResponse struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	URL       string    `json:"url"`       // Secure Cloudinary URL
	PublicID  string    `json:"public_id"` // For delete/update
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// MediaUpdateRequest for PUT /merchant/products/:media_id
type MediaUpdateRequest struct {
	File *string `form:"file" validate:"omitempty"`    // Optional new file
	URL  *string `form:"url" validate:"omitempty,url"` // Or new URL (e.g., external)
	Type *string `form:"type" validate:"omitempty,oneof=image video"`
}

// MediaDeleteRequest (empty body, if needed for reason)
type MediaDeleteRequest struct {
	Reason string `json:"reason" validate:"omitempty,max=500"`
}

type CategoryResponse struct {
	ID           uint                   `json:"id"`
	Name         string                 `json:"name"`
	ParentID     *uint                  `json:"parent_id"`
	CategorySlug string                 `json:"category_slug"`
	Attributes   map[string]interface{} `json:"attributes"`
	Parent       *CategoryResponse      `json:"parent"`
}

type CreateReviewDTO struct {
	ProductID string `json:"product_id" validate:"required"`
	Rating    int    `json:"rating" validate:"required,min=1,max=5"`
	Comment   string `json:"comment" validate:"omitempty"`
}

type UpdateReviewDTO struct {
	Rating  *int    `json:"rating" validate:"omitempty,min=1,max=5"`
	Comment *string `json:"comment" validate:"omitempty"`
}

type ReviewResponseDTO struct {
	//ID        uint      `json:"id"`
	ProductID string    `json:"product_id"`
	ProductName string `json:"product_name"`
	//UserID    uint      `json:"user_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserName  string    `json:"user_name"`
	Image   string    `json:"image_url"`
	//AverageRating uint  `json:"average_rating"`
}

type AddWishlistItemDTO struct {
	ProductID string `json:"product_id" validate:"required,uuid"` // UUID validation for product_id
}

// WishlistItemResponseDTO represents a single wishlist item in the response
type WishlistItemResponseDTO struct {
	ProductID  string  `json:"product_id"`
	Name       string  `json:"name"`
	FinalPrice float64 `json:"total_price"`
	Discount   float64 `json:"discount" validate:"gte=0"`
	//DiscountType string    `json:"discount_type" validate:"oneof=fixed percentage ''"`
	CategorySlug string `json:"category_slug"`
	PrimaryImage string `json:"primary_image,omitempty"`

	//SKU          string    `json:"sku"`
	//MerchantID   string    `json:"merchant_id"`
	AddedAt time.Time `json:"added_at"`
}

// WishlistResponseDTO represents the entire wishlist in the response
type WishlistResponseDTO struct {
	UserID uint                      `json:"user_id"`
	Items  []WishlistItemResponseDTO `json:"items"`
}

// ProductAutocompleteResponse represents a single product suggestion for autocomplete.
type ProductAutocompleteResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	SKU         string `json:"sku"`
	Description string `json:"description,omitempty"`
}

// AutocompleteResponse is the full response with a list of suggestions.
type AutocompleteResponse struct {
	Suggestions []ProductAutocompleteResponse `json:"suggestions"`
}

type ProductFilterRequest struct {
	CategoryID   *uint    `form:"category_id"`
	CategoryName *string  `form:"category_name"`
	CategorySlug *string  `form:"category_slug"`
	MinPrice     *float64 `form:"min_price"`
	MaxPrice     *float64 `form:"max_price"`
	InStock      *bool    `form:"in_stock"`
	SearchTerm   *string  `form:"search" binding:"omitempty,max=100"`
	MerchantName *string  `form:"merchant_name"`

	// NEW: Variant attribute filters
	Color    *string `form:"color"`
	Size     *string `form:"size"`
	Material *string `form:"material"`
	Pattern  *string `form:"pattern"`

	// NEW: Sorting
	SortBy   *string `form:"sort_by" binding:"omitempty,oneof=price price_desc name name_desc created newest oldest rating"`
	OnSale   *bool   `form:"on_sale"` // Products with discounts
	Featured *bool   `form:"featured"`

	// Pagination
	Page  int `form:"page" binding:"omitempty,min=1"`
	Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
}

func (f *ProductFilterRequest) GetOffset() int {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Limit <= 0 {
		f.Limit = 20
	}
	return (f.Page - 1) * f.Limit
}

func (f *ProductFilterRequest) GetLimit() int {
	if f.Limit <= 0 || f.Limit > 100 {
		return 20
	}
	return f.Limit
}

// Generate hash for cache key
func (f *ProductFilterRequest) Hash() string {
	data, _ := json.Marshal(f)
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}

// BulkInventoryUpdateInput represents the request body for bulk inventory updates
type BulkInventoryUpdateInput struct {
	InventoryID string `json:"inventory_id" validate:"required,uuid"`
	Delta       int    `json:"delta" validate:"required"`
}
