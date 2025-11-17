package dto

import "time"


type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "Pending"
	PaymentStatusCompleted PaymentStatus = "Completed"
	PaymentStatusFailed    PaymentStatus = "Failed"
	PaymentStatusRefunded    PaymentStatus = "Refunded"
)

type InitializePaymentRequest struct {
	OrderID   uint    `json:"order_id" validate:"required"`
	Amount    float64 `json:"amount" validate:"required,gt=0"` // In kobo for Paystack
	Email     string  `json:"email" validate:"required,email"`
	Currency  string  `json:"currency" validate:"required"` // e.g., "NGN"
}

type PaymentResponse struct {
	ID             uint      `json:"id"`
	OrderID        uint      `json:"order_id"`
	Amount         float64   `json:"amount"`
	Currency       string    `json:"currency"`
	Status         PaymentStatus    `json:"status"` // e.g., "success", "pending"
	TransactionID  string    `json:"transaction_id"` // Paystack ref
	AuthorizationURL string  `json:"authorization_url,omitempty"` // For checkout redirect
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}