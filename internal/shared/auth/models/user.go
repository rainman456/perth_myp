package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"unique;not null"`
	Name     string `gorm:"type:varchar(100);not null"`
	Password string // Empty for OAuth users
	Role     string `gorm:"not null"` // "customer" (default) or "merchant" (upgraded by admin)
	GoogleID string    // Google ID for OAuth
	Country  string `gorm:"type:varchar(100)"` // Optional country field
}