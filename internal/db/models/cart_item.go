package models

import (
	"gorm.io/gorm"
)

type CartItem struct {
	gorm.Model
	CartID     uint    `gorm:"not null" json:"cart_id"`
	ProductID  uint    `gorm:"not null" json:"product_id"`
	Quantity   int     `gorm:"not null" json:"quantity"`
	MerchantID uint    `gorm:"not null" json:"merchant_id"`
	Cart       Cart    `gorm:"foreignKey:CartID"`
	Product    Product `gorm:"foreignKey:ProductID"`
	Merchant   Merchant `gorm:"foreignKey:MerchantID"`
}