package models


import (
	"time"
	"gorm.io/datatypes"
)
/*
// MerchantStatus defines the possible status values for a merchant
type MerchantBasicInfo struct {
	StoreName     string `gorm:"column:store_name;size:255;not null" json:"store_name" validate:"required"`
	Name          string `gorm:"column:name;size:255;not null" json:"name" validate:"required"`
	PersonalEmail string `gorm:"column:personal_email;size:255;not null;unique" json:"personal_email" validate:"required,email"`
	WorkEmail     string `gorm:"column:work_email;size:255;not null;unique" json:"work_email" validate:"required,email"`
	PhoneNumber   string `gorm:"column:phone_number;size:50" json:"phone_number"`
}

// MerchantAddress holds address information as JSONB
type MerchantAddress struct {
	PersonalAddress datatypes.JSON `gorm:"column:personal_address;type:jsonb;not null" json:"personal_address" validate:"required"`
	WorkAddress     datatypes.JSON `gorm:"column:work_address;type:jsonb;not null" json:"work_address" validate:"required"`
}

// MerchantBusinessInfo holds business-related information
type MerchantBusinessInfo struct {
	BusinessType               string `gorm:"column:business_type;size:100" json:"business_type"`
	Website                    string `gorm:"column:website;size:255" json:"website"`
	BusinessDescription        string `gorm:"column:business_description;type:text" json:"business_description"`
	BusinessRegistrationNumber string `gorm:"column:business_registration_number;size:255;not null;unique" json:"business_registration_number" validate:"required"`
}

// MerchantDocuments holds document-related information
type MerchantDocuments struct {
	StoreLogoURL                   string `gorm:"column:store_logo_url;size:255" json:"store_logo_url"`
	BusinessRegistrationCertificate string `gorm:"column:business_registration_certificate;size:255" json:"business_registration_certificate"`
}

// MerchantApplication holds the information required for a merchant onboarding application
type MerchantApplication struct {
	ID                string            `gorm:"primaryKey;column:id;type:uuid;default:uuid_generate_v4()" json:"id,omitempty"`
	MerchantBasicInfo                   `gorm:"embedded"`
	MerchantAddress                     `gorm:"embedded"`
	MerchantBusinessInfo                `gorm:"embedded"`
	MerchantDocuments                   `gorm:"embedded"`
	Status            string            `gorm:"column:status;type:varchar(20);default:pending;not null" json:"status"`
	CreatedAt         time.Time         `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time         `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (MerchantApplication) TableName() string {
	return "merchant_application"
}

// MerchantStatus defines the possible statuses for a merchant
type MerchantStatus string

const (
	MerchantStatusActive    MerchantStatus = "active"
	MerchantStatusSuspended MerchantStatus = "suspended"
)

// Merchant holds the active merchant account details after approval
type Merchant struct {
	ID                string            `gorm:"primaryKey;column:id;type:uuid;default:uuid_generate_v4()" json:"id,omitempty"`
	ApplicationID     string            `gorm:"column:application_id;type:uuid;not null;unique" json:"application_id"`
	MerchantID            string            `gorm:"column:merchant_id;type:uuid;not null;unique" json:"user_id"`
	MerchantBasicInfo                   `gorm:"embedded"`
	MerchantAddress                     `gorm:"embedded"`
	MerchantBusinessInfo                `gorm:"embedded"`
	MerchantDocuments                   `gorm:"embedded"`
	Password          string            `gorm:"column:password;size:255;not null" json:"password" validate:"required"`
	Status            MerchantStatus    `gorm:"column:status;type:varchar(20);default:active;index" json:"status"`
	CommissionTier    string            `gorm:"column:commission_tier;default:standard" json:"commission_tier"`
	CommissionRate    float64           `gorm:"column:commission_rate;default:5.00" json:"commission_rate"`
	AccountBalance    float64           `gorm:"column:account_balance;default:0.00" json:"account_balance"`
	TotalSales        float64           `gorm:"column:total_sales;default:0.00" json:"total_sales"`
	TotalPayouts      float64           `gorm:"column:total_payouts;default:0.00" json:"total_payouts"`
	StripeAccountID   string            `gorm:"column:stripe_account_id" json:"stripe_account_id"`
	PayoutSchedule    string            `gorm:"column:payout_schedule;default:weekly" json:"payout_schedule"`
	LastPayoutDate    *time.Time        `gorm:"column:last_payout_date" json:"last_payout_date"`
	Banner            string            `gorm:"column:banner;size:255" json:"banner"`
	Policies          datatypes.JSON    `gorm:"column:policies;type:jsonb" json:"policies"`
	CreatedAt         time.Time         `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time         `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	//Products          []Product         `gorm:"foreignKey:MerchantID" json:"products,omitempty"`
	//CartItems         []CartItem        `gorm:"foreignKey:MerchantID" json:"cart_items,omitempty"`
	//OrderItems        []OrderItem       `gorm:"foreignKey:MerchantID" json:"order_items,omitempty"`
	//Payouts           []Payout          `gorm:"foreignKey:MerchantID" json:"payouts,omitempty"`
}

func (Merchant) TableName() string {
	return "merchant"
}
*/

type MerchantBasicInfo struct {
	StoreName     string `gorm:"column:store_name;size:255;not null" json:"store_name" validate:"required"`
	Name          string `gorm:"column:name;size:255;not null" json:"name" validate:"required"`
	PersonalEmail string `gorm:"column:personal_email;size:255;not null;unique" json:"personal_email" validate:"required,email"`
	WorkEmail     string `gorm:"column:work_email;size:255;not null;unique" json:"work_email" validate:"required,email"`
	PhoneNumber   string `gorm:"column:phone_number;size:50" json:"phone_number"`
}

// MerchantAddress holds address information as JSONB
type MerchantAddress struct {
	PersonalAddress datatypes.JSON `gorm:"column:personal_address;type:jsonb;not null" json:"personal_address" validate:"required"`
	WorkAddress     datatypes.JSON `gorm:"column:work_address;type:jsonb;not null" json:"work_address" validate:"required"`
}

// MerchantBusinessInfo holds business-related information
type MerchantBusinessInfo struct {
	BusinessType               string `gorm:"column:business_type;size:100" json:"business_type"`
	Website                    string `gorm:"column:website;size:255" json:"website"`
	BusinessDescription        string `gorm:"column:business_description;type:text" json:"business_description"`
	BusinessRegistrationNumber string `gorm:"column:business_registration_number;size:255;not null;unique" json:"business_registration_number" validate:"required"`
}

// MerchantDocuments holds document-related information
type MerchantDocuments struct {
	StoreLogoURL                   string `gorm:"column:store_logo_url;size:255" json:"store_logo_url"`
	BusinessRegistrationCertificate string `gorm:"column:business_registration_certificate;size:255" json:"business_registration_certificate"`
}

// MerchantApplication holds the information required for a merchant onboarding application
type MerchantApplication struct {
	ID                string            `gorm:"primaryKey;column:id;type:uuid;default:uuid_generate_v4()" json:"id,omitempty"`
	MerchantBasicInfo                   `gorm:"embedded"`
	MerchantAddress                     `gorm:"embedded"`
	MerchantBusinessInfo                `gorm:"embedded"`
	MerchantDocuments                   `gorm:"embedded"`
	Status            string            `gorm:"column:status;type:varchar(20);default:pending;not null" json:"status"`
	CreatedAt         time.Time         `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time         `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (MerchantApplication) TableName() string {
	return "merchant_application"
}

// MerchantStatus defines the possible statuses for a merchant
type MerchantStatus string

const (
	MerchantStatusActive    MerchantStatus = "active"
	MerchantStatusSuspended MerchantStatus = "suspended"
)

// Merchant holds the active merchant account details after approval
type Merchant struct {
	ID                string            `gorm:"primaryKey;column:id;type:uuid;default:uuid_generate_v4()" json:"id,omitempty"`
	ApplicationID     string            `gorm:"column:application_id;type:uuid;not null;unique" json:"application_id"`
	MerchantID            string            `gorm:"column:merchant_id;type:uuid;not null;unique" json:"user_id"`
	MerchantBasicInfo                   `gorm:"embedded"`
	MerchantAddress                     `gorm:"embedded"`
	MerchantBusinessInfo                `gorm:"embedded"`
	MerchantDocuments                   `gorm:"embedded"`
	Password          string            `gorm:"column:password;size:255;not null" json:"password" validate:"required"`
	Status            MerchantStatus    `gorm:"column:status;type:varchar(20);default:active;index" json:"status"`
	CommissionTier    string            `gorm:"column:commission_tier;default:standard" json:"commission_tier"`
	CommissionRate    float64           `gorm:"column:commission_rate;default:5.00" json:"commission_rate"`
	AccountBalance    float64           `gorm:"column:account_balance;default:0.00" json:"account_balance"`
	TotalSales        float64           `gorm:"column:total_sales;default:0.00" json:"total_sales"`
	TotalPayouts      float64           `gorm:"column:total_payouts;default:0.00" json:"total_payouts"`
	StripeAccountID   string            `gorm:"column:stripe_account_id" json:"stripe_account_id"`
	PayoutSchedule    string            `gorm:"column:payout_schedule;default:weekly" json:"payout_schedule"`
	LastPayoutDate    *time.Time        `gorm:"column:last_payout_date" json:"last_payout_date"`
	Banner            string            `gorm:"column:banner;size:255" json:"banner"`
	Policies          datatypes.JSON    `gorm:"column:policies;type:jsonb" json:"policies"`
	CreatedAt         time.Time         `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time         `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	//Products          []Product         `gorm:"foreignKey:MerchantID" json:"products,omitempty"`
	//CartItems         []CartItem        `gorm:"foreignKey:MerchantID" json:"cart_items,omitempty"`
	//OrderItems        []OrderItem       `gorm:"foreignKey:MerchantID" json:"order_items,omitempty"`
	//Payouts           []Payout          `gorm:"foreignKey:MerchantID" json:"payouts,omitempty"`
}

func (Merchant) TableName() string {
	return "merchant"
}