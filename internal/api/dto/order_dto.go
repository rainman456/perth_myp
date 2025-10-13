package dto

import (
	"api-customer-merchant/internal/db/models"
	"time"
)

//import "api-customer-merchant/internal/db/models"

// CreateOrderRequest defines the request body for creating an order.
type CreateOrderRequest struct {
	UserID uint `json:"user_id"` // The ID of the user placing the order.
	//UserID uint `json:"user_id"` // The ID of the user placing the order
}

// OrderResponse defines the structure for order-related responses.
// OrderStatus describes the status of an order.
//
// swagger:enum
//
// Pending
// Completed
// Cancelled
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "Pending"
	OrderStatusCompleted OrderStatus = "Completed"
	OrderStatusCancelled OrderStatus = "Cancelled"
)

// OrderResponse defines the structure for an order response.
type OrderResponse struct {
	ID           uint                `json:"id"`
	UserID       uint                `json:"user_id"`
	Status       OrderStatus         `json:"status"`
	OrderItems   []OrderItemResponse `json:"order_items"`
	TotalAmount  float64             `json:"total_amount"`
	PaymentStatus string             `json:"payment_status"`
}

// OrderItemResponse defines the structure for individual items in an order.
type OrderItemResponse struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type CancelOrderRequest struct {
	Reason string `json:"reason" validate:"omitempty,max=500"` // Optional cancellation reason
}





type OrdersResponse struct {
	ID        uint             `json:"id"`
	UserID    uint             `json:"user_id"`
	OrderItems []OrdersItemResponse `json:"order_items"`
	Status    models.OrderStatus `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type OrdersItemResponse struct {
	ID        uint   `json:"id"`
	OrderID   uint   `json:"order_id"`
	ProductID string   `json:"product_id"`
	Product   OrderProductResponse `json:"product"`
	MerchantID string  `json:"merchant_id"`
	Merchant  OrderMerchantResponse `json:"merchant"`
	Quantity  uint   `json:"quantity"`
}

type OrderProductResponse struct {
	ID          string   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       float64 `json:"price"`
}

type OrderMerchantResponse struct {
	ID        string `json:"id"`
	StoreName string `json:"store_name"`
}