package models

import (
	"time"

	"gorm.io/gorm"
)

// Announcement model (matching TS announcements)
type Announcement struct {
	gorm.Model
	ID        string    `gorm:"type:varchar;primaryKey" json:"id"`
	Title     string    `gorm:"type:text;not null" json:"title"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// Dispute model (matching TS disputes)
type Dispute struct {
	gorm.Model
	ID          string    `gorm:"type:varchar;primaryKey" json:"id"`
	OrderID     string    `gorm:"type:varchar;not null" json:"order_id"`
	CustomerID  uint    `gorm:"not null" json:"customer_id"`
	MerchantID  string    `gorm:"type:varchar;not null" json:"merchant_id"`
	Reason      string    `gorm:"type:text;not null" json:"reason"`
	Description string    `gorm:"type:text;not null" json:"description"`
	Status      string    `gorm:"type:text;not null;default:'open'" json:"status"`
	Resolution  string    `gorm:"type:text" json:"resolution"`
	Customer           User              `gorm:"foreignKey:CustomerID"`
	Order         Order                 `gorm:"foreignKey:OrderID"`
	Merchant          Merchant         `gorm:"foreignKey:MerchantID;references:MerchantID"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	ResolvedAt  time.Time `json:"resolved_at"`
}

// ReturnRequest model (matching TS return_requests)
type ReturnRequest struct {
	gorm.Model
	ID               string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	OrderItemID      uint    `gorm:"not null" json:"order_item_id"`
	CustomerID        uint    `gorm:"not null" json:"customer_id"`
	Reason           string    `gorm:"type:text" json:"reason"`
	Status           string    `gorm:"type:varchar(255);default:'Pending'" json:"status"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	OrderItem        OrderItem `gorm:"foreignKey:OrderItemID"`
	Customer           User              `gorm:"foreignKey:CustomerID"`

}

// Settings model (matching TS settings)
// type Settings struct {
// 	gorm.Model
// 	ID              string                 `gorm:"type:text;primaryKey;default:'global'" json:"id"`
// 	Fees            float64                `gorm:"type:decimal(10,2);not null;default:5.00" json:"fees"`
// 	TaxRate         float64                `gorm:"type:decimal(10,2);not null;default:0.00" json:"tax_rate"`
// 	ShippingOptions map[string]interface{} `gorm:"type:jsonb;not null" json:"shipping_options"`
// }
