package dto

import "time"

// CreateDisputeDTO for creating a dispute
type CreateDisputeDTO struct {
	OrderID     string `json:"order_id" binding:"required"`
	Reason      string `json:"reason" binding:"required,max=100"`
	Description string `json:"description" binding:"required,max=1000"`
}

// DisputeResponseDTO for dispute responses
type DisputeResponseDTO struct {
	ID          string    `json:"id"`
	OrderID     string    `json:"order_id"`
	CustomerID  uint      `json:"customer_id"`
	MerchantID  string    `json:"merchant_id"`
	Reason      string    `json:"reason"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Resolution  string    `json:"resolution,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	ResolvedAt  time.Time `json:"resolved_at,omitempty"`
}

// CreateReturnRequestDTO for return requests
type CreateReturnRequestDTO struct {
	OrderItemID string `json:"order_item_id" binding:"required,uuid"`
	Reason      string `json:"reason" binding:"required,max=500"`
}

// ReturnRequestResponseDTO for return request responses
type ReturnRequestResponseDTO struct {
	ID          string    `json:"id"`
	OrderItemID string    `json:"order_item_id"`
	CustomerID  uint      `json:"customer_id"`
	Reason      string    `json:"reason"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}