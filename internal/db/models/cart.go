package models

import (
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	UserID uint `gorm:"not null" json:"user_id"`
	User   User `gorm:"foreignKey:UserID"`
}