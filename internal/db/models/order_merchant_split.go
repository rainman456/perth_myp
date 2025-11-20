// Add to internal/db/models/order_merchant_split.go

package models

import (
    "fmt"
    "time"
    "github.com/shopspring/decimal"
    "gorm.io/gorm"
)

// OrderMerchantSplitStatus defines possible status values for order merchant splits
type OrderMerchantSplitStatus string

const (
    OrderMerchantSplitStatusPending    OrderMerchantSplitStatus = "pending"
    OrderMerchantSplitStatusProcessing OrderMerchantSplitStatus = "processing"
    OrderMerchantSplitStatusCompleted  OrderMerchantSplitStatus = "completed"
)

// Valid checks if the status is one of the allowed values
func (s OrderMerchantSplitStatus) Valid() error {
    switch s {
    case OrderMerchantSplitStatusPending, OrderMerchantSplitStatusProcessing, OrderMerchantSplitStatusCompleted:
        return nil
    default:
        return fmt.Errorf("invalid order merchant split status: %s", s)
    }
}

type OrderMerchantSplit struct {
    gorm.Model
    OrderID    uint    `gorm:"index"`
    MerchantID string  `gorm:"type:uuid;index"`
    
    // Use gorm:"type:numeric(12,2)" to force PostgreSQL numeric
    AmountDue  decimal.Decimal `gorm:"type:numeric(12,2)"`
    Fee        decimal.Decimal `gorm:"type:numeric(12,2)"`
    
    Status     OrderMerchantSplitStatus `gorm:"type:varchar(20);default:'pending'"`
    HoldUntil  time.Time
    
    Merchant   Merchant `gorm:"foreignKey:MerchantID;references:MerchantID"`
    Order      Order    `gorm:"foreignKey:OrderID"`
}

// BeforeCreate validates the Status field
func (oms *OrderMerchantSplit) BeforeCreate(tx *gorm.DB) error {
    if err := oms.Status.Valid(); err != nil {
        return err
    }
    return nil
}

// BeforeUpdate validates the Status field
func (oms *OrderMerchantSplit) BeforeUpdate(tx *gorm.DB) error {
    if err := oms.Status.Valid(); err != nil {
        return err
    }
    return nil
}