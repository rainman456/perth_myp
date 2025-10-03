package models

import (
	"fmt"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// OrderStatus defines possible order status values
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "Pending"
	OrderStatusCompleted OrderStatus = "Completed"
	OrderStatusCancelled OrderStatus = "Cancelled"
)

// Valid checks if the status is one of the allowed values
func (s OrderStatus) Valid() error {
	switch s {
	case OrderStatusPending, OrderStatusCompleted, OrderStatusCancelled:
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
     UserID         uint              `gorm:"not null"`
    SubTotal       decimal.Decimal   `gorm:"type:decimal(10,2)" json:"sub_total"`
     TotalAmount    decimal.Decimal   `gorm:"type:decimal(10,2)" json:"total_amount"`
     Status         OrderStatus      `gorm:"type:varchar(20);not null;default:'Pending'" json:"status"`
     ShippingMethod string            `gorm:"type:varchar(50)" json:"shipping_method"`
     CouponCode     *string           `gorm:"type:varchar(50)" json:"coupon_code"`
    Currency       string            `gorm:"type:varchar(3);default:'NGN'" json:"currency"`
     User           User              `gorm:"foreignKey:UserID"`
     OrderItems     []OrderItem       `gorm:"foreignKey:OrderID"`
    Payments       []Payment         `gorm:"foreignKey:OrderID"`
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
