package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	MerchantID  uint    `gorm:"not null" json:"merchant_id"`
	Name        string  `gorm:"size:255;not null" json:"name"`
	Description string  `gorm:"type:text" json:"description"`
	SKU         string  `gorm:"size:100;unique;not null" json:"sku"`
	Price       float64 `gorm:"type:decimal(10,2);not null" json:"price"`
	CategoryID  uint    `gorm:"not null" json:"category_id"`
	Merchant    Merchant `gorm:"foreignKey:MerchantID"`
	Category    Category `gorm:"foreignKey:CategoryID"`
}