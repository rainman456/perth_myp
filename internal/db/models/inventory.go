package models

import (
	"gorm.io/gorm"
)

type Inventory struct {
	gorm.Model
	ProductID         uint   `gorm:"not null" json:"product_id"`
	StockQuantity     int    `gorm:"not null" json:"stock_quantity"`
	LowStockThreshold int    `gorm:"not null;default:10" json:"low_stock_threshold"`
	Product           Product `gorm:"foreignKey:ProductID"`
}