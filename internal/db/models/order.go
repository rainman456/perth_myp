package models

import (
	"fmt"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// OrderStatus defines possible order status values
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

// order status: paid - confirmed - processing - completed - cancelled
// order item status: processing - confirmed - declined - sent to aronova hub - out for delivery - delivered

// Valid checks if the status is one of the allowed values
func (s OrderStatus) Valid() error {
	switch s {
	case OrderStatusPending, OrderStatusConfirmed, OrderStatusPaid, OrderStatusProcessing, 
		OrderStatusShipped, OrderStatusCompleted, OrderStatusCancelled, 
		OrderStatusOutForDelivery, OrderStatusDelivered:
		return nil
	default:
		return fmt.Errorf("invalid order status: %s", s)
	}
}


// type Order struct {
// 	gorm.Model
// 	UserID      uint        `gorm:"not null" json:"user_id"`
// 	TotalAmount float64     `gorm:"type:decimal(10,2);not null" json:"total_amount"`
// 	Status      OrderStatus `gorm:"type:varchar(20);not null;default:'Pending'" json:"status"`
// 	User        User        `gorm:"foreignKey:UserID"`
// 	OrderItems  []OrderItem `gorm:"foreignKey:OrderID"`
// }

type Order struct {
	gorm.Model
	UserID         uint            `gorm:"not null"`
	SubTotal       decimal.Decimal `gorm:"type:decimal(10,2)" json:"sub_total"`
	TotalAmount    decimal.Decimal `gorm:"type:decimal(10,2)" json:"total_amount"`
	Status         OrderStatus     `gorm:"type:varchar(20);not null;default:'Pending'" json:"status"`
	ShippingMethod string          `gorm:"type:varchar(50)" json:"shipping_method"`
	CouponCode     *string         `gorm:"type:varchar(50)" json:"coupon_code"`
	Currency       string          `gorm:"type:varchar(3);default:'NGN'" json:"currency"`
	User           User            `gorm:"foreignKey:UserID"`
	OrderItems     []OrderItem     `gorm:"foreignKey:OrderID"`
	Payments       []Payment       `gorm:"foreignKey:OrderID"`
}

// BeforeCreate validates the Status field
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if err := o.Status.Valid(); err != nil {
		return err
	}
	return nil
}

// BeforeUpdate validates the Status field
func (o *Order) BeforeUpdate(tx *gorm.DB) error {
	if err := o.Status.Valid(); err != nil {
		return err
	}
	return nil
}




func (o *Order) UpdateStatusBasedOnItems() {
	if len(o.OrderItems) == 0 {
		return
	}

	allDelivered := true
	allDeclined := true
	allSentOrOutForDelivery := true
	anyProcessed := false

	for _, item := range o.OrderItems {
		if item.FulfillmentStatus != FulfillmentStatusDelivered {
			allDelivered = false
		}
		if item.FulfillmentStatus != FulfillmentStatusDeclined {
			allDeclined = false
		}
		if item.FulfillmentStatus != FulfillmentStatusSentToAronovaHub && 
		   item.FulfillmentStatus != FulfillmentStatusOutForDelivery {
			allSentOrOutForDelivery = false
		}
		if item.FulfillmentStatus != FulfillmentStatusProcessing {
			anyProcessed = true
		}
	}

	// Apply business logic
	if allDelivered {
		o.Status = OrderStatusCompleted
	} else if allDeclined {
		o.Status = OrderStatusCancelled
	} else if allSentOrOutForDelivery {
		o.Status = OrderStatusShipped
	} else if anyProcessed {
		o.Status = OrderStatusProcessing
	}
}