package models

import (
	"fmt"
	"gorm.io/gorm"
)

// PayoutStatus defines possible payout status values
type PayoutStatus string

const (
	PayoutStatusPending   PayoutStatus = "Pending"
	PayoutStatusCompleted PayoutStatus = "Completed"
)

// Valid checks if the status is one of the allowed values
func (s PayoutStatus) Valid() error {
	switch s {
	case PayoutStatusPending, PayoutStatusCompleted:
		return nil
	default:
		return fmt.Errorf("invalid payout status: %s", s)
	}

	
}

type Payout struct {
	gorm.Model
	MerchantID      string         `gorm:"not null" json:"merchant_id"`
	Amount          float64      `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status          PayoutStatus `gorm:"type:varchar(20);not null;default:'Pending'" json:"status"`
	PayoutAccountID     string `gorm:"size:255" json:"payout_account_id"`        // ← Remove not null
    PayStackTransferID  string `gorm:"size:255" json:"paystack_transfer_id"`     // ← R
	Merchant        Merchant     `gorm:"foreignKey:MerchantID"`
}

// BeforeCreate validates the Status field
func (p *Payout) BeforeCreate(tx *gorm.DB) error {
	if err := p.Status.Valid(); err != nil {
		return err
	}
	return nil
}

// BeforeUpdate validates the Status field
func (p *Payout) BeforeUpdate(tx *gorm.DB) error {
	if err := p.Status.Valid(); err != nil {
		return err
	}
	return nil
}
