package models

import (
	"fmt"
	"gorm.io/gorm"
)

// MerchantStatus defines the possible status values for a merchant
type MerchantStatus string

const (
	MerchantStatusApproved  MerchantStatus = "Approved"
	MerchantStatusPending   MerchantStatus = "Pending"
	MerchantStatusInReview  MerchantStatus = "InReview"
	MerchantStatusRejected  MerchantStatus = "Rejected"
)

// Valid checks if the status is one of the allowed values
func (s MerchantStatus) Valid() error {
	switch s {
	case MerchantStatusApproved, MerchantStatusPending, MerchantStatusInReview, MerchantStatusRejected:
		return nil
	default:
		return fmt.Errorf("invalid merchant status: %s", s)
	}
}

// Basic merchant info
type MerchantBasicInfo struct {
	Name          string `gorm:"size:255;not null" json:"name"`
	StoreName     string `gorm:"size:255;not null" json:"store_name"`
	PersonalEmail string `gorm:"size:255;not null;unique" json:"personal_email"`
	WorkEmail     string `gorm:"size:255;not null;unique" json:"work_email"`
	Password      string `gorm:"type:varchar(255);not null;default:'XXXX'"`
	PhoneNumber   string `gorm:"size:50" json:"phone_number"`
}

// Address info
type MerchantAddress struct {
	StreetAddress string `gorm:"size:255" json:"street_address"`
	City          string `gorm:"size:100" json:"city"`
	Country       string `gorm:"size:100" json:"country"`
	ZipCode       string `gorm:"size:20" json:"zip_code"`
	WorkAddress   string `gorm:"size:255" json:"work_address"`
}

// Business info
type MerchantBusinessInfo struct {
	BusinessType        string `gorm:"size:100" json:"business_type"`
	Website             string `gorm:"size:255" json:"website"`
	BusinessDescription string `gorm:"type:text" json:"business_description"`
}

// Document info (only store logo and business registration certificate)
type MerchantDocuments struct {
	StoreLogoURL                   string `gorm:"size:255" json:"store_logo_url"`                    // image file (jpg, png)
	BusinessRegistrationCertificate string `gorm:"size:255" json:"business_registration_certificate"` // file (pdf, doc, docx, jpg, png)
}

// Main merchant model
type Merchant struct {
	gorm.Model
	MerchantBasicInfo
	MerchantAddress
	MerchantBusinessInfo
	MerchantDocuments
	Status MerchantStatus `gorm:"type:varchar(20);not null;default:'Pending'" json:"status"`
}

// BeforeCreate validates the Status field before saving to the database
func (m *Merchant) BeforeCreate(tx *gorm.DB) error {
	if err := m.Status.Valid(); err != nil {
		return err
	}
	return nil
}

// BeforeUpdate validates the Status field before updating the database
func (m *Merchant) BeforeUpdate(tx *gorm.DB) error {
	if err := m.Status.Valid(); err != nil {
		return err
	}
	return nil
}