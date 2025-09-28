package models

import (
    "time"
    "github.com/shopspring/decimal"
    "gorm.io/gorm"
)

type OrderMerchantSplit struct {
    gorm.Model
    OrderID    uint    `gorm:"index"`
    MerchantID string  `gorm:"type:uuid;index"`  // Match Merchant.ID type
    AmountDue  decimal.Decimal
    Fee        decimal.Decimal  // Platform cut
    Status     string  `gorm:"default:'pending'"`  // pending, payout_requested, paid, reversed
    HoldUntil  time.Time
	Merchant   Merchant `gorm:"foreignKey:MerchantID;references:MerchantID"`
	Order             Order             `gorm:"foreignKey:OrderID"`
}