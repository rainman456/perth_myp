package models

import (
	"fmt"
	"gorm.io/gorm"
)

// FulfillmentStatus defines possible fulfillment status values
type FulfillmentStatus string

const (
    FulfillmentStatusProcessing       FulfillmentStatus = "Processing" // NEW
    FulfillmentStatusConfirmed        FulfillmentStatus = "Confirmed"
    FulfillmentStatusDeclined         FulfillmentStatus = "Declined"
    FulfillmentStatusSentToAronovaHub FulfillmentStatus = "SentToAronovaHub"
    FulfillmentStatusShipped          FulfillmentStatus = "Shipped"
)
// Valid checks if the status is one of the allowed values
func (s FulfillmentStatus) Valid() error {
	switch s {
	case FulfillmentStatusProcessing, FulfillmentStatusConfirmed, FulfillmentStatusDeclined, FulfillmentStatusSentToAronovaHub, FulfillmentStatusShipped:
		return nil
	default:
		return fmt.Errorf("invalid fulfillment status: %s", s)
	}
}

type OrderItem struct {
	gorm.Model
	OrderID   uint    `gorm:"not null;index" json:"order_id"`
	ProductID string  `gorm:"not null;index" json:"product_id"`
	VariantID *string `gorm:"type:uuid;index" json:"variant_id"`
	//ProductID         uint              `gorm:"not null;index" json:"product_id"`
	MerchantID        string            `gorm:"not null;index" json:"merchant_id"`
	Quantity          int               `gorm:"not null" json:"quantity"`
	Price             float64           `gorm:"type:decimal(10,2);not null" json:"price"`
	FulfillmentStatus FulfillmentStatus `gorm:"type:varchar(20);not null;default:'New'" json:"fulfillment_status"`
	Order             Order             `gorm:"foreignKey:OrderID"`
	Product           Product           `gorm:"foreignKey:ProductID;references:ID"`
	Merchant          Merchant          `gorm:"foreignKey:MerchantID;references:MerchantID"`
	Variant           *Variant          `gorm:"foreignKey:VariantID"`
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

func (oi *OrderItem) CanBeModified() bool {
    return oi.FulfillmentStatus != FulfillmentStatusSentToAronovaHub &&
           oi.FulfillmentStatus != FulfillmentStatusShipped
}


// ValidateStatusTransition checks if status change is allowed
func (oi *OrderItem) ValidateStatusTransition(newStatus FulfillmentStatus) error {
    switch oi.FulfillmentStatus {
    case FulfillmentStatusProcessing:
        if newStatus != FulfillmentStatusConfirmed && newStatus != FulfillmentStatusDeclined {
            return fmt.Errorf("can only transition to Confirmed/Declined from Processing")
        }
    case FulfillmentStatusConfirmed:
        if newStatus != FulfillmentStatusSentToAronovaHub {
            return fmt.Errorf("can only transition to SentToAronovaHub from Confirmed")
        }
    case FulfillmentStatusSentToAronovaHub, FulfillmentStatusShipped, FulfillmentStatusDeclined:
        return fmt.Errorf("cannot transition from terminal status: %s", oi.FulfillmentStatus)
    }
    return nil
}