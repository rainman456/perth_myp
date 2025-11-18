package models

import (
	"time"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Settings model represents global marketplace settings
type Settings struct {
	ID              string         `gorm:"type:text;primaryKey;default:'global'" json:"id"`
	Fees            float64        `gorm:"type:decimal(10,2);not null;default:5.00" json:"fees"`
	TaxRate         float64        `gorm:"type:decimal(10,2);not null;default:0.00" json:"tax_rate"`
	ShippingOptions datatypes.JSON `gorm:"type:jsonb;not null" json:"shipping_options"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Settings) TableName() string {
	return "settings"
}

// ShippingOption represents a single shipping option
type ShippingOption struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Enabled     bool    `json:"enabled"`
}