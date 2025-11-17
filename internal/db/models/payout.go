package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PayoutStatus defines possible payout status values
type PayoutStatus string

const (
	PayoutStatusPending   PayoutStatus = "Pending"
	PayoutStatusCompleted PayoutStatus = "Completed"
	PayoutStatusOpen PayoutStatus = "Open"
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
	ID                 string       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CreatedAt          time.Time    `json:"created_at"`
	UpdatedAt          time.Time    `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	MerchantID         string       `gorm:"type:uuid;not null;index" json:"merchant_id"`
	Amount             float64      `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status             PayoutStatus `gorm:"type:varchar(20);not null;default:'Pending'" json:"status"`
	PayoutAccountID    string       `gorm:"size:255" json:"payout_account_id"`
	PayStackTransferID string       `gorm:"size:255" json:"paystack_transfer_id"`
	
	Merchant           Merchant     `gorm:"foreignKey:MerchantID;references:MerchantID"`
}
// BeforeCreate validates the Status field
func (p *Payout) BeforeCreate(tx *gorm.DB) error {
	if err := p.Status.Valid(); err != nil {
		return err
	}

	if p.ID == "" {
		p.ID = uuid.New().String()
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
