package models

import (
	

	"gorm.io/gorm"
)



// type UserAddress struct {
// 	gorm.Model
// 	UserID  uint   `gorm:"not null;index"`  // Foreign key to User
// 	Address string `gorm:"type:varchar(255);not null"`  // Increased length for real addresses
// }

type UserAddress struct {
	gorm.Model

	UserID                uint   `gorm:"not null;index" json:"user_id"`                            // FK to users table
	PhoneNumber           string `gorm:"type:varchar(20)" json:"phone_number,omitempty"`          // primary phone
	AdditionalPhoneNumber string `gorm:"type:varchar(20)" json:"additional_phone_number,omitempty"` // secondary phone
	DeliveryAddress       string `gorm:"type:text" json:"delivery_address,omitempty"`             // delivery address (freeform)
	ShippingAddress       string `gorm:"type:text" json:"shipping_address,omitempty"`             // shipping address (freeform)
	State                 string `gorm:"type:varchar(100)" json:"state,omitempty"`
	LGA                   string `gorm:"type:varchar(100)" json:"lga,omitempty"`
	User             User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	// you can add a 'Label' or 'Type' field (home, work) if you want to categorize addresses
	// Type                  string `gorm:"type:varchar(50)" json:"type,omitempty"`
}

type User struct {
	gorm.Model
	Email    string `gorm:"unique;not null"`
	Name     string `gorm:"type:varchar(100);not null"`
	Password string // Empty for OAuth users
	//Role     string `gorm:"not null"` // "customer" (default) or "merchant" (upgraded by admin)
	GoogleID string // Google ID for OAuth
	Country  string `gorm:"type:varchar(100)"` // Optional country field
	Addresses []UserAddress      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	//Carts    []Cart  `gorm:"foreignKey:UserID" json:"carts,omitempty"`
	//Orders   []Order `gorm:"foreignKey:UserID" json:"orders,omitempty"`
	Wishlists    []UserWishlist `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"` // Has many UserWishlists
	Reviews []Review `gorm:"foreignKey:UserID"`
}
