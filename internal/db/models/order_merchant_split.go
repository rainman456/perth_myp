package models

import (
    "time"
    "github.com/shopspring/decimal"
    "gorm.io/gorm"
)

type OrderMerchantSplit struct {
    gorm.Model
    OrderID    uint    `gorm:"index"`
    MerchantID string  `gorm:"type:uuid;index"`
    
    // Use gorm:"type:numeric(12,2)" to force PostgreSQL numeric
    AmountDue  decimal.Decimal `gorm:"type:numeric(12,2)"`
    Fee        decimal.Decimal `gorm:"type:numeric(12,2)"`
    
    Status     string  `gorm:"default:'pending'"`
    HoldUntil  time.Time
    
    Merchant   Merchant `gorm:"foreignKey:MerchantID;references:MerchantID"`
    Order      Order    `gorm:"foreignKey:OrderID"`
}