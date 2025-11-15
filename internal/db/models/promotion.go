package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// PromotionType defines possible promotion types
type PromotionType string

const (
	PromotionTypePercentage PromotionType = "percentage"
	PromotionTypeFixed      PromotionType = "fixed"
)

// Valid checks if the promotion type is one of the allowed values
func (p PromotionType) Valid() error {
	switch p {
	case PromotionTypePercentage, PromotionTypeFixed:
		return nil
	default:
		return fmt.Errorf("invalid promotion type: %s", p)
	}
}

// PromotionStatus defines possible promotion status values
type PromotionStatus string

const (
	PromotionStatusActive   PromotionStatus = "active"
	PromotionStatusInactive PromotionStatus = "inactive"
	PromotionStatusExpired  PromotionStatus = "expired"
)

// Valid checks if the status is one of the allowed values
func (s PromotionStatus) Valid() error {
	switch s {
	case PromotionStatusActive, PromotionStatusInactive, PromotionStatusExpired:
		return nil
	default:
		return fmt.Errorf("invalid promotion status: %s", s)
	}
}

// Promotion represents a promotional campaign for products
type Promotion struct {
	ID          string         `gorm:"primaryKey;column:id;type:uuid;default:uuid_generate_v4()" json:"id,omitempty"`
	Name        string         `gorm:"column:name;size:255;not null" json:"name" validate:"required"`
	Description string         `gorm:"column:description;type:text" json:"description"`
	Type        PromotionType  `gorm:"column:type;type:varchar(20);default:'percentage'" json:"type"`
	Discount    float64        `gorm:"column:discount;type:decimal(5,2);not null" json:"discount" validate:"required"`
	StartDate   time.Time      `gorm:"column:start_date;not null" json:"start_date" validate:"required"`
	EndDate     time.Time      `gorm:"column:end_date;not null" json:"end_date" validate:"required"`
	Status      PromotionStatus `gorm:"column:status;type:varchar(20);default:'active'" json:"status"`
	MerchantID  string         `gorm:"column:merchant_id;type:uuid;not null;index" json:"merchant_id"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// Relations
	Merchant  Merchant         `gorm:"foreignKey:MerchantID;references:MerchantID"`
	Products  []Product        `gorm:"many2many:promotion_products;" json:"products,omitempty"`
}

// BeforeCreate validates the Type and Status fields
func (p *Promotion) BeforeCreate(tx *gorm.DB) error {
	if err := p.Type.Valid(); err != nil {
		return err
	}
	if err := p.Status.Valid(); err != nil {
		return err
	}
	return nil
}

// BeforeUpdate validates the Type and Status fields
func (p *Promotion) BeforeUpdate(tx *gorm.DB) error {
	if err := p.Type.Valid(); err != nil {
		return err
	}
	if err := p.Status.Valid(); err != nil {
		return err
	}
	return nil
}

// TableName specifies the table name for Promotion
func (Promotion) TableName() string {
	return "promotions"
}