package models

import (
	"gorm.io/gorm"
)

type CartItem struct {
	gorm.Model
	CartID     uint     `gorm:"not null;index:idx_cartitem_cart_product" json:"cart_id"` // Composite for common WHERE
    VariantID  *string  `gorm:"type:uuid;index:idx_cartitem_cart_product" json:"variant_id"`
    ProductID  string   `gorm:"not null;index:idx_cartitem_cart_product" json:"product_id"`
    Quantity   int      `gorm:"not null" json:"quantity"`
    MerchantID string   `gorm:"not null;index" json:"merchant_id"`
	Cart       Cart     `gorm:"foreignKey:CartID"`
	Product    Product  `gorm:"foreignKey:ProductID"`
	Merchant   Merchant `gorm:"foreignKey:MerchantID;references:MerchantID"`
	Variant    *Variant `gorm:"foreignKey:VariantID"`
}
