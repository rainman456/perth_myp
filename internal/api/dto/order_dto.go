package dto

//import "api-customer-merchant/internal/db/models"

// CreateOrderRequest defines the request body for creating an order.
type CreateOrderRequest struct {
	UserID uint `json:"user_id"` // The ID of the user placing the order.
	//UserID uint `json:"user_id"` // The ID of the user placing the order
}

// OrderResponse defines the structure for order-related responses.
type OrderResponse struct {
	ID         uint                `json:"id"`
	UserID     uint                `json:"user_id"`
	Status     string              `json:"status"` // Consider using the models.OrderStatus type directly
	OrderItems []OrderItemResponse `json:"order_items"`
	TotalAmount   float64 `json:"total_amount"`
	PaymentStatus string  `json:"payment_status"`
}

// OrderItemResponse defines the structure for individual items in an order.
type OrderItemResponse struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
