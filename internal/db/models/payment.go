package models

import (
	"fmt"
	"gorm.io/gorm"
)

// PaymentStatus defines possible payment status values
type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "Pending"
	PaymentStatusCompleted PaymentStatus = "Completed"
	PaymentStatusFailed   PaymentStatus = "Failed"
)

// Valid checks if the status is one of the allowed values
func (s PaymentStatus) Valid() error {
	switch s {
	case PaymentStatusPending, PaymentStatusCompleted, PaymentStatusFailed:
		return nil
	default:
		return fmt.Errorf("invalid payment status: %s", s)
	}
}

type Payment struct {
	gorm.Model
	OrderID uint          `gorm:"not null" json:"order_id"`
	Amount  float64       `gorm:"not null" json:"amount"`
	Status  PaymentStatus `gorm:"type:varchar(20);not null;default:'Pending'" json:"status"`
	Order   Order         `gorm:"foreignKey:OrderID"`
}

// BeforeCreate validates the Status field
func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if err := p.Status.Valid(); err != nil {
		return err
	}
	return nil
}

// BeforeUpdate validates the Status field
func (p *Payment) BeforeUpdate(tx *gorm.DB) error {
	if err := p.Status.Valid(); err != nil {
		return err
	}
	return nil
}