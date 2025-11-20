package dto

import (
	//"api-customer-merchant/internal/db/models"
	"time"
)

//import "api-customer-merchant/internal/db/models"

//CreateOrderRequest defines the request body for creating an order.
// type CreateOrderRequest struct {
// 	ShippingAddress string             `json:"shipping_address"`
// }

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
	OrderStatusPending       OrderStatus = "Pending"
	OrderStatusConfirmed     OrderStatus = "Confirmed"
	OrderStatusPaid          OrderStatus = "Paid"
	OrderStatusProcessing    OrderStatus = "Processing"
	OrderStatusShipped       OrderStatus = "Shipped"
	OrderStatusCompleted     OrderStatus = "Completed"
	OrderStatusCancelled     OrderStatus = "Cancelled"
	OrderStatusOutForDelivery OrderStatus = "OutForDelivery" // Keep for backward compatibility if needed
	OrderStatusDelivered     OrderStatus = "Delivered"       
)

// Add this new DTO
type CreateOrderRequest struct {
	ShippingMethod string `json:"shipping_method" validate:"required"`
}


// OrderResponse defines the structure for an order response.
type OrderResponse struct {
	ID           uint                `json:"id"`
	UserID       uint                `json:"user_id"`
	Status       OrderStatus         `json:"status"`
	OrderItems   []OrderItemResponse `json:"order_items"`
	TotalAmount  float64             `json:"total_amount"`
	DeliveryAddress string             `json:"delivery_address"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	PaymentAuthorizationURL string              `json:"payment_authorization_url,omitempty"` 
    PaymentReference        string              `json:"payment_reference,omitempty"`         
}

// OrderItemResponse defines the structure for individual items in an order.
type OrderItemResponse struct {
	ProductID string  `json:"product_id"`
	Name      string    `json:"name"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Image   string    `json:"image_url"`
	CategorySlug string  `json:"category_slug"`
	

}

type CancelOrderRequest struct {
	Reason string `json:"reason" validate:"omitempty,max=500"` // Optional cancellation reason
}





type OrdersResponse struct {
	ID        uint             `json:"id"`
	UserID    uint             `json:"user_id"`
	OrderItems []OrdersItemResponse `json:"order_items"`
	Status    OrderStatus `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type OrdersItemResponse struct {
	ID        uint   `json:"id"`
	OrderID   uint   `json:"order_id"`
	Product   OrderProductResponse `json:"product"`
	Merchant  OrderMerchantResponse `json:"merchant"`
	Quantity  uint   `json:"quantity"`
}

type OrderProductResponse struct {
	ID          string   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       float64 `json:"price"`
	Image   string    `json:"image_url"`
	CategorySlug string  `json:"category_slug"`
}

type OrderMerchantResponse struct {
	ID        string `json:"id"`
	StoreName string `json:"store_name"`
}