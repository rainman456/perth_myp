package models

import (
	"fmt"
	"gorm.io/gorm"
)

// FulfillmentStatus defines possible fulfillment status values
type FulfillmentStatus string

const (
	FulfillmentStatusNew     FulfillmentStatus = "New"
	FulfillmentStatusShipped FulfillmentStatus = "Shipped"
)

// Valid checks if the status is one of the allowed values
func (s FulfillmentStatus) Valid() error {
	switch s {
	case FulfillmentStatusNew, FulfillmentStatusShipped:
		return nil
	default:
		return fmt.Errorf("invalid fulfillment status: %s", s)
	}
}

// type OrderItem struct {
// 	gorm.Model
// 	OrderID           uint              `gorm:"not null" json:"order_id"`
// 	ProductID         string              `gorm:"not null" json:"product_id"`
// 	MerchantID        uint              `gorm:"not null" json:"merchant_id"`
// 	Quantity          int               `gorm:"not null" json:"quantity"`
// 	Price             float64           `gorm:"type:decimal(10,2);not null" json:"price"`
// 	FulfillmentStatus FulfillmentStatus `gorm:"type:varchar(20);not null;default:'New'" json:"fulfillment_status"`
// 	Order             Order             `gorm:"foreignKey:OrderID"`
// 	Product           Product           `gorm:"foreignKey:ProductID"`
// 	Merchant          Merchant          `gorm:"foreignKey:MerchantID"`
// }

type OrderItem struct {
	gorm.Model
	OrderID   uint   `gorm:"not null;index" json:"order_id"`
	ProductID string `gorm:"not null;index" json:"product_id"`
	//ProductID         uint              `gorm:"not null;index" json:"product_id"`
	MerchantID        string            `gorm:"not null;index" json:"merchant_id"`
	Quantity          int               `gorm:"not null" json:"quantity"`
	Price             float64           `gorm:"type:decimal(10,2);not null" json:"price"`
	FulfillmentStatus FulfillmentStatus `gorm:"type:varchar(20);not null;default:'New'" json:"fulfillment_status"`
	Order             Order             `gorm:"foreignKey:OrderID"`
	Product           Product           `gorm:"foreignKey:ProductID;references:ID"`
	Merchant          Merchant         `gorm:"foreignKey:MerchantID;references:MerchantID"`
}

// BeforeCreate validates the FulfillmentStatus field
func (oi *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if err := oi.FulfillmentStatus.Valid(); err != nil {
		return err
	}
	return nil
}

// BeforeUpdate validates the FulfillmentStatus field
func (oi *OrderItem) BeforeUpdate(tx *gorm.DB) error {
	if err := oi.FulfillmentStatus.Valid(); err != nil {
		return err
	}
	return nil
}
