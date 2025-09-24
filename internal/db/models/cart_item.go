package models

import (
	"gorm.io/gorm"
)

type CartItem struct {
	gorm.Model
	CartID     uint     `gorm:"not null" json:"cart_id"`
	VariantID  *string  `gorm:"type:uuid;index"`
	ProductID  string   `gorm:"not null" json:"product_id"`
	Quantity   int      `gorm:"not null" json:"quantity"`
	MerchantID string   `gorm:"not null" json:"merchant_id"`
	Cart       Cart     `gorm:"foreignKey:CartID"`
	Product    Product  `gorm:"foreignKey:ProductID"`
	Merchant   Merchant `gorm:"foreignKey:MerchantID;references:MerchantID"`
	Variant    *Variant `gorm:"foreignKey:VariantID"`
}
