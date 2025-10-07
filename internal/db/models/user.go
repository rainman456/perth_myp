package models

import (
	

	"gorm.io/gorm"
)



type UserAddress struct {
	gorm.Model
	UserID  uint   `gorm:"not null;index"`  // Foreign key to User
	Address string `gorm:"type:varchar(255);not null"`  // Increased length for real addresses
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
