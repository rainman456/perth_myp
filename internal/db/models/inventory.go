package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// type Inventory struct {
// 	gorm.Model
// 	ProductID         string   `gorm:"not null" json:"product_id"`
// 	StockQuantity     int    `gorm:"not null" json:"stock_quantity"`
// 	LowStockThreshold int    `gorm:"not null;default:10" json:"low_stock_threshold"`
// 	Product           Product `gorm:"foreignKey:ProductID"`
// }

// type VendorInventory struct {
// 	ID                string    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
// 	VariantID         string    `gorm:"type:uuid;not null;unique;index" json:"variant_id"`
// 	MerchantID        string    `gorm:"type:uuid;not null;index" json:"merchant_id"`
// 	ProductID         *string   `gorm:"type:uuid;index"` // Nullable: For simple products
// 	Quantity          int       `gorm:"default:0;not null;check:quantity >= 0" json:"quantity"`
// 	ReservedQuantity  int       `gorm:"default:0;check:reserved_quantity >= 0" json:"reserved_quantity"`
// 	LowStockThreshold int       `gorm:"default:10" json:"low_stock_threshold"`
// 	BackorderAllowed  bool      `gorm:"default:false" json:"backorder_allowed"`
// 	CreatedAt         time.Time `json:"created_at"`
// 	UpdatedAt         time.Time `json:"updated_at"`

// 	Variant  *Variant `gorm:"foreignKey:VariantID"`
// 	Product  *Product `gorm:"foreignKey:ProductID"`
// 	Merchant Merchant `gorm:"foreignKey:MerchantID"`
// }



type Inventory struct {
	ID                string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ProductID         *string    `gorm:"type:uuid;index" json:"product_id,omitempty"` // Optional: For simple products
	VariantID         *string    `gorm:"type:uuid;index" json:"variant_id,omitempty"` // Optional: For variants
	MerchantID        string    `gorm:"type:uuid;not null;index" json:"merchant_id"` // Required: Vendor-specific
	Quantity          int       `gorm:"default:0;not null;check:quantity >= 0" json:"quantity"`
	ReservedQuantity  int       `gorm:"default:0;not null;check:reserved_quantity >= 0" json:"reserved_quantity"`
	LowStockThreshold int       `gorm:"default:5" json:"low_stock_threshold"`
	BackorderAllowed  bool      `gorm:"default:false" json:"backorder_allowed"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	Product  *Product  `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"` // Optional belongs to Product
	Variant  *Variant  `gorm:"foreignKey:VariantID;constraint:OnDelete:CASCADE"` // Optional belongs to Variant
	Merchant Merchant  `gorm:"foreignKey:MerchantID;constraint:OnDelete:RESTRICT"` // Belongs to Merchant, no cascade
}


func (vi *Inventory) BeforeCreate(tx *gorm.DB) error {
	if vi.ID == "" {
		vi.ID = uuid.New().String()
	}
	if (vi.VariantID != nil && vi.ProductID != nil) || (vi.VariantID == nil && vi.ProductID == nil) {
		return errors.New("exactly one of VariantID or ProductID must be set")
	}
	return nil
}
