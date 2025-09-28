package models

import (
	"fmt"
	"gorm.io/gorm"
)

// CartStatus defines possible cart status values
type CartStatus string

const (
	CartStatusActive    CartStatus = "Active"
	CartStatusAbandoned CartStatus = "Abandoned"
	CartStatusConverted CartStatus = "Converted"
)

// Valid checks if the status is one of the allowed values
func (s CartStatus) Valid() error {
	switch s {
	case CartStatusActive, CartStatusAbandoned, CartStatusConverted:
		return nil
	default:
		return fmt.Errorf("invalid cart status: %s", s)
	}
}

type Cart struct {
	gorm.Model
	UserID     uint       `gorm:"not null" json:"user_id"`
	Status     CartStatus `gorm:"type:varchar(20);not null;default:'Active'" json:"status"`
	SubTotal   float64    `gorm:"-" json:"subtotal"` // Computed
	TaxTotal   float64    `gorm:"-" json:"tax_total"`
	ShipTotal  float64    `gorm:"-" json:"shipping_total"`
	GrandTotal float64    `gorm:"-" json:"grand_total"`
	User       User       `gorm:"foreignKey:UserID"`
	CartItems  []CartItem `gorm:"foreignKey:CartID"`
}

// BeforeCreate validates the Status field
func (c *Cart) BeforeCreate(tx *gorm.DB) error {
	if err := c.Status.Valid(); err != nil {
		return err
	}
	return nil
}

// BeforeUpdate validates the Status field
func (c *Cart) BeforeUpdate(tx *gorm.DB) error {
	if err := c.Status.Valid(); err != nil {
		return err
	}
	return nil
}

func (c *Cart) AfterFind(tx *gorm.DB) error {
	c.ComputeTotals()
	return nil
}

func (c *Cart) ComputeTotals() {
	c.SubTotal = 0
	for _, item := range c.CartItems {
		c.SubTotal += float64(item.Quantity) * (item.Product.BasePrice).InexactFloat64() // Assume BasePrice in Product
	}
	// Stub: c.TaxTotal = 0.1 * c.SubTotal // Or call pricing
	// c.ShipTotal = 10.00
	c.GrandTotal = c.SubTotal // + c.TaxTotal + c.ShipTotal
}
