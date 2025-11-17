package dto

type MerchantApplyDTO struct {
	StoreName                       string         `json:"store_name" binding:"required"`
	Name                            string         `json:"name" binding:"required"`
	PersonalEmail                   string         `json:"personal_email" binding:"required,email"`
	WorkEmail                       string         `json:"work_email" binding:"required,email"`
	PhoneNumber                     string         `json:"phone_number"`
	PersonalAddress                 map[string]any `json:"personal_address" binding:"required"`
	WorkAddress                     map[string]any `json:"work_address" binding:"required"`
	BusinessType                    string         `json:"business_type"`
	Website                         string         `json:"website"`
	BusinessDescription             string         `json:"business_description"`
	BusinessRegistrationNumber      string         `json:"business_registration_number" binding:"required"`
	StoreLogoURL                    string         `json:"store_logo_url"`
	BusinessRegistrationCertificate string         `json:"business_registration_certificate"`
}

type MerchantApplyResponse struct {
	ID                              string         `json:"id"`
	StoreName                       string         `json:"store_name"`
	Name                            string         `json:"name"`
	PersonalEmail                   string         `json:"personal_email"`
	WorkEmail                       string         `json:"work_email"`
	PhoneNumber                     string         `json:"phone_number"`
	PersonalAddress                 map[string]any `json:"personal_address"`
	WorkAddress                     map[string]any `json:"work_address"`
	BusinessType                    string         `json:"business_type"`
	Website                         string         `json:"website"`
	BusinessDescription             string         `json:"business_description"`
	BusinessRegistrationNumber      string         `json:"business_registration_number"`
	StoreLogoURL                    string         `json:"store_logo_url"`
	BusinessRegistrationCertificate string         `json:"business_registration_certificate"`
	Status                          string         `json:"status"`
	CreatedAt                       string         `json:"created_at"`
	UpdatedAt                       string         `json:"updated_at"`
}




type MerchRequestPasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest is used to reset password with a token
type MerchResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type MerchantLogin struct {
	Work_Email string `json:"work_email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
}

type BankDetailsRequest struct {
	BankName      string `json:"bank_name,omitempty"`
	Currency      string `json:"currency,omitempty"`
	AccountName   string `json:"account_name"`
	AccountNumber string `json:"account_number"`
	BankCode      string `json:"bank_code"`
}

// BankDetailsResponse represents the response for bank details operations
type BankDetailsResponse struct {
	ID            uint   `json:"id"`
	MerchantID    string `json:"merchant_id"`
	BankName      string `json:"bank_name"`
	BankCode      string `json:"bank_code"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
	RecipientCode string `json:"recipient_code,omitempty"`
	Currency      string `json:"currency"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// Merchant Order Management DTOs
type MerchantOrderResponse struct {
	ID              uint                        `json:"id"`
	UserID          uint                        `json:"user_id"`
	Status          string                      `json:"status"`
	OrderItems      []MerchantOrderItemResponse `json:"order_items"`
	TotalAmount     float64                     `json:"total_amount"`
	DeliveryAddress string                      `json:"delivery_address"`
	CreatedAt       string                      `json:"created_at"`
	UpdatedAt       string                      `json:"updated_at"`
}

// OrderItemActionResponse represents the response for order item actions
type OrderItemActionResponse struct {
	Message string `json:"message"`
}

// MerchantOrderItemResponse represents an order item in a merchant's order
type MerchantOrderItemResponse struct {
	ID                uint    `json:"id"`
	ProductID         string  `json:"product_id"`
	Name              string  `json:"name"`
	Quantity          int     `json:"quantity"`
	Price             float64 `json:"price"`
	Image             string  `json:"image_url"`
	FulfillmentStatus string  `json:"fulfillment_status"`
}

// Payout DTOs
type PayoutRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
}

type PayoutResponse struct {
	ID                 string  `json:"id"`
	MerchantID      string  `json:"merchant_id"`
	Amount          float64 `json:"amount"`
	Status          string  `json:"status"`
	PayoutAccountID    string  `json:"payout_account_id,omitempty"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

type PayoutHistoryResponse struct {
	Payouts []PayoutResponse `json:"payouts"`
}

// Promotion DTOs
type CreatePromotionRequest struct {
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description"`
	Discount    float64  `json:"discount" validate:"required,gt=0,lte=100"`
	StartDate   string   `json:"start_date" validate:"required"`
	EndDate     string   `json:"end_date" validate:"required"`
	ProductIDs  []string `json:"product_ids" validate:"required"`
}

type PromotionResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Discount    float64  `json:"discount"`
	StartDate   string   `json:"start_date"`
	EndDate     string   `json:"end_date"`
	ProductIDs  []string `json:"product_ids"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

// Store Profile DTOs
type UpdateStoreProfileRequest struct {
	StoreName           string `json:"store_name"`
	Name                string `json:"name"`
	WorkEmail           string `json:"work_email" validate:"email"`
	PhoneNumber         string `json:"phone_number"`
	BusinessType        string `json:"business_type"`
	Website             string `json:"website"`
	BusinessDescription string `json:"business_description"`
	Banner              string `json:"banner"`
}

type StoreProfileResponse struct {
	ID                         string         `json:"id"`
	StoreName                  string         `json:"store_name"`
	Name                       string         `json:"name"`
	PersonalEmail              string         `json:"personal_email"`
	WorkEmail                  string         `json:"work_email"`
	PhoneNumber                string         `json:"phone_number"`
	PersonalAddress            map[string]any `json:"personal_address"`
	WorkAddress                map[string]any `json:"work_address"`
	BusinessType               string         `json:"business_type"`
	Website                    string         `json:"website"`
	BusinessDescription        string         `json:"business_description"`
	BusinessRegistrationNumber string         `json:"business_registration_number"`
	StoreLogoURL               string         `json:"store_logo_url"`
	Banner                     string         `json:"banner"`
	CommissionTier             string         `json:"commission_tier"`
	CommissionRate             float64        `json:"commission_rate"`
	AccountBalance             float64        `json:"account_balance"`
	TotalSales                 float64        `json:"total_sales"`
	TotalPayouts               float64        `json:"total_payouts"`
	PayoutSchedule             string         `json:"payout_schedule"`
	LastPayoutDate             *string        `json:"last_payout_date"`
	CreatedAt                  string         `json:"created_at"`
	UpdatedAt                  string         `json:"updated_at"`
}

// UpdateMerchantProfileInput represents the input for updating a merchant's profile
type UpdateMerchantProfileInput struct {
	StoreName           *string         `json:"store_name,omitempty"`
	Name                *string         `json:"name,omitempty"`
	PersonalEmail       *string         `json:"personal_email,omitempty" validate:"omitempty,email"`
	WorkEmail           *string         `json:"work_email,omitempty" validate:"omitempty,email"`
	PhoneNumber         *string         `json:"phone_number,omitempty"`
	PersonalAddress     *map[string]any `json:"personal_address,omitempty"`
	WorkAddress         *map[string]any `json:"work_address,omitempty"`
	BusinessType        *string         `json:"business_type,omitempty"`
	Website             *string         `json:"website,omitempty"`
	BusinessDescription *string         `json:"business_description,omitempty"`
	StoreLogoURL        *string         `json:"store_logo_url,omitempty"`
	Banner              *string         `json:"banner,omitempty"`
}
