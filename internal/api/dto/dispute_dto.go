package dto

import "time"

// CreateDisputeDTO for creating a dispute
type CreateDisputeDTO struct {
	OrderID     string `json:"order_id" binding:"required"`
	Reason      string `json:"reason" binding:"required,max=100"`
	Description string `json:"description" binding:"required,max=1000"`
}


type PayoutStatus string

const (
	PayoutStatusPending   PayoutStatus = "Pending"
	PayoutStatusCompleted PayoutStatus = "Completed"
	PayoutStatusOpen PayoutStatus = "Open"
)
// DisputeResponseDTO for dispute responses
type CreateDisputeResponseDTO struct {
	ID          string    `json:"id"`
	OrderID     string    `json:"order_id"`
	CustomerID  uint      `json:"customer_id"`
	MerchantID  string    `json:"merchant_id"`
	Reason      string    `json:"reason"`
	Description string    `json:"description"`
	Status      PayoutStatus    `json:"status"`
	//Resolution  string    `json:"resolution,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	ResolvedAt  time.Time `json:"resolved_at,omitempty"`
}



type DisputeItemDTO struct {
	ProductID       string    `json:"product_id"`
	ProductName     string    `json:"product_name"`
	ProductImageURL string    `json:"product_image_url"`
	CategorySlug    string    `json:"category_slug"`
	Reason          string    `json:"reason"`
	Description     string    `json:"description"`
	Resolution      string    `json:"resolution"`
	ResolvedAt      *time.Time `json:"resolved_at"` // Use pointer for nullable field
	CreatedAt       time.Time `json:"created_at"`
}

// DisputeResponseDTO represents the grouped dispute response
type DisputeResponseDTO struct {
	OrderID        uint            `json:"order_id"`
	OrderCreatedAt time.Time       `json:"order_created_at"`
	Status         PayoutStatus          `json:"status"`
	CustomerID     uint            `json:"customer_id"`
	MerchantID     string          `json:"merchant_id"`
	Disputes       []DisputeItemDTO `json:"disputes"`
}




// CreateReturnRequestDTO for return requests
type CreateReturnRequestDTO struct {
	OrderItemID uint `json:"order_item_id" binding:"required,uuid"`
	Reason      string `json:"reason" binding:"required,max=500"`
}

// ReturnRequestResponseDTO for return request responses
type ReturnRequestResponseDTO struct {
	ID          string    `json:"id"`
	OrderItemID uint    `json:"order_item_id"`
	CustomerID  uint      `json:"customer_id"`
	Reason      string    `json:"reason"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}



// UpdateReturnRequestDTO - For future use
type UpdateReturnRequestDTO struct {
	Reason *string `json:"reason" validate:"omitempty"`
	Status *string `json:"status" validate:"omitempty,oneof=Pending Approved Rejected"`
}

// ReturnRequestResponseDTO - Unchanged (for single return request)
type CreateReturnRequestResponseDTO struct {
	ID          string    `json:"id"`
	OrderItemID  uint    `json:"order_item_id"`
	CustomerID  uint      `json:"customer_id"`
	Reason      string    `json:"reason"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ReturnItemDTO struct {
	ProductID       string    `json:"product_id"`
	ProductName     string    `json:"product_name"`
	ProductImageURL string    `json:"product_image_url"`
	CategorySlug    string    `json:"category_slug"`
	Reason          string    `json:"reason"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

type ReturnResponseDTO struct {
	OrderID        uint            `json:"order_id"`
	OrderCreatedAt time.Time       `json:"order_created_at"`
	CustomerID     uint            `json:"customer_id"`
	Returns        []ReturnItemDTO `json:"returns"`
}