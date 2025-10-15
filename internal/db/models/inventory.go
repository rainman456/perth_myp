package models

import (
	"errors"
	"fmt"
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

// type Inventory struct {
// 	ID                string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
// 	ProductID         *string    `gorm:"type:uuid;index" json:"product_id,omitempty"` // Optional: For simple products
// 	VariantID         *string    `gorm:"type:uuid;index" json:"variant_id,omitempty"` // Optional: For variants
// 	MerchantID        string    `gorm:"type:uuid;not null;index" json:"merchant_id"` // Required: Vendor-specific
// 	Quantity          int       `gorm:"default:0;not null;check:quantity >= 0" json:"quantity"`
// 	ReservedQuantity  int       `gorm:"default:0;not null;check:reserved_quantity >= 0" json:"reserved_quantity"`
// 	LowStockThreshold int       `gorm:"default:5" json:"low_stock_threshold"`
// 	BackorderAllowed  bool      `gorm:"default:false" json:"backorder_allowed"`
// 	CreatedAt         time.Time `json:"created_at"`
// 	UpdatedAt         time.Time `json:"updated_at"`

// 	Product  *Product  `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"` // Optional belongs to Product
// 	Variant  *Variant  `gorm:"foreignKey:VariantID;constraint:OnDelete:CASCADE"` // Optional belongs to Variant
// 	Merchant Merchant  `gorm:"foreignKey:MerchantID;constraint:OnDelete:RESTRICT"` // Belongs to Merchant, no cascade
// }

// func (vi *Inventory) BeforeCreate(tx *gorm.DB) error {
// 	if vi.ID == "" {
// 		vi.ID = uuid.New().String()
// 	}
// 	if (vi.VariantID != nil && vi.ProductID != nil) || (vi.VariantID == nil && vi.ProductID == nil) {
// 		return errors.New("exactly one of VariantID or ProductID must be set")
// 	}
// 	return nil
// }


type InventoryStatus string

const (
	InventoryStatusInStock    InventoryStatus = "in_stock"
	InventoryStatusLowStock   InventoryStatus = "low_stock"
	InventoryStatusOutOfStock InventoryStatus = "out_of_stock"
	InventoryStatusBackorder  InventoryStatus = "backorder"
)

// Inventory - Stock management (merchant/vendor specific)
type Inventory struct {
	// Identity
	ID         string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ProductID  *string `gorm:"type:uuid;index" json:"product_id,omitempty"`   // NULL for variant inventory
	VariantID  *string `gorm:"type:uuid;index" json:"variant_id,omitempty"`   // NULL for simple product inventory
	MerchantID string  `gorm:"type:uuid;not null;index" json:"merchant_id"`  // Vendor owns inventory

	// Stock levels
	Quantity         int `gorm:"not null;default:0;check:quantity >= 0" json:"quantity"`
	ReservedQuantity int `gorm:"not null;default:0;check:reserved_quantity >= 0" json:"reserved_quantity"`
    
	// Configuration
	LowStockThreshold int  `gorm:"not null;default:5" json:"low_stock_threshold"`
	BackorderAllowed  bool `gorm:"default:false" json:"backorder_allowed"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Product  *Product  `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	Variant  *Variant  `gorm:"foreignKey:VariantID;constraint:OnDelete:CASCADE"`
	Merchant *Merchant `gorm:"foreignKey:MerchantID;constraint:OnDelete:RESTRICT"`
}

// ============ VALIDATION ============

// BeforeCreate hook - validate inventory setup
func (inv *Inventory) BeforeCreate(tx *gorm.DB) error {
	if inv.ID == "" {
		inv.ID = uuid.New().String()
	}

	// Exactly one of ProductID or VariantID must be set
	hasProduct := inv.ProductID != nil
	hasVariant := inv.VariantID != nil

	if (hasProduct && hasVariant) || (!hasProduct && !hasVariant) {
		return fmt.Errorf("exactly one of ProductID or VariantID must be set")
	}

	// Validate reserved quantity doesn't exceed total
	if inv.ReservedQuantity > inv.Quantity {
		return fmt.Errorf("reserved quantity (%d) cannot exceed total quantity (%d)", inv.ReservedQuantity, inv.Quantity)
	}

	// Validate low stock threshold
	if inv.LowStockThreshold < 0 {
		return errors.New("low stock threshold cannot be negative")
	}

	return nil
}

// BeforeUpdate hook - validate updates
func (inv *Inventory) BeforeUpdate(tx *gorm.DB) error {
	if inv.ReservedQuantity > inv.Quantity {
		return fmt.Errorf("reserved quantity (%d) cannot exceed total quantity (%d)", inv.ReservedQuantity, inv.Quantity)
	}
	return nil
}

// ============ BUSINESS LOGIC METHODS ============

// GetAvailableQuantity returns quantity available for purchase
func (inv *Inventory) GetAvailableQuantity() int {
	available := inv.Quantity - inv.ReservedQuantity
	if available < 0 {
		return 0
	}
	return available
}

// GetStatus computes current inventory status
func (inv *Inventory) GetStatus() InventoryStatus {
	available := inv.GetAvailableQuantity()

	// Out of stock
	if available <= 0 {
		if inv.BackorderAllowed {
			return InventoryStatusBackorder
		}
		return InventoryStatusOutOfStock
	}

	// Low stock
	if available <= inv.LowStockThreshold {
		return InventoryStatusLowStock
	}

	// In stock
	return InventoryStatusInStock
}

// CanFulfill checks if order can be fulfilled
func (inv *Inventory) CanFulfill(quantity int) bool {
	available := inv.GetAvailableQuantity()
	
	if available >= quantity {
		return true
	}

	// Can fulfill with backorder if allowed
	if inv.BackorderAllowed {
		return true
	}

	return false
}

// Reserve reserves stock for an order
func (inv *Inventory) Reserve(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	if !inv.CanFulfill(quantity) {
		return fmt.Errorf(
			"cannot reserve %d units. Available: %d, Backorder allowed: %v",
			quantity,
			inv.GetAvailableQuantity(),
			inv.BackorderAllowed,
		)
	}

	inv.ReservedQuantity += quantity
	return nil
}

// Release cancels a reservation
func (inv *Inventory) Release(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	if inv.ReservedQuantity < quantity {
		return fmt.Errorf("cannot release %d units. Reserved: %d", quantity, inv.ReservedQuantity)
	}

	inv.ReservedQuantity -= quantity
	return nil
}

// Commit converts reservation to actual sale (reduces total quantity)
func (inv *Inventory) Commit(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	if inv.ReservedQuantity < quantity {
		return fmt.Errorf("cannot commit %d units. Reserved: %d", quantity, inv.ReservedQuantity)
	}

	inv.Quantity -= quantity
	inv.ReservedQuantity -= quantity
	return nil
}

// Refund adds inventory back (e.g., order cancellation or return)
func (inv *Inventory) Refund(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	inv.Quantity += quantity
	if inv.ReservedQuantity >= quantity {
		inv.ReservedQuantity -= quantity
	} else {
		inv.ReservedQuantity = 0
	}

	return nil
}

// AdjustStock manually adjusts inventory (for admin/merchant operations)
func (inv *Inventory) AdjustStock(delta int) error {
	newQuantity := inv.Quantity + delta

	if newQuantity < 0 {
		return fmt.Errorf("cannot set quantity below 0. Current: %d, Delta: %d", inv.Quantity, delta)
	}

	if newQuantity < inv.ReservedQuantity {
		return fmt.Errorf(
			"cannot reduce quantity below reserved amount. New: %d, Reserved: %d",
			newQuantity,
			inv.ReservedQuantity,
		)
	}

	inv.Quantity = newQuantity
	return nil
}

// IsLowStock checks if inventory is below threshold
func (inv *Inventory) IsLowStock() bool {
	return inv.GetAvailableQuantity() <= inv.LowStockThreshold
}
