package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// type Product struct {
// 	gorm.Model
// 	MerchantID  uint    `gorm:"not null" json:"merchant_id"`
// 	Name        string  `gorm:"size:255;not null" json:"name"`
// 	Description string  `gorm:"type:text" json:"description"`
// 	SKU         string  `gorm:"size:100;unique;not null" json:"sku"`
// 	Price       float64 `gorm:"type:decimal(10,2);not null" json:"price"`
// 	CategoryID  uint    `gorm:"not null" json:"category_id"`
// 	Merchant    Merchant `gorm:"foreignKey:MerchantID"`
// 	Category    Category `gorm:"foreignKey:CategoryID"`
// }

type AttributesMap map[string]string

// Value implements driver.Valuer
func (a AttributesMap) Value() (driver.Value, error) {
    return json.Marshal(a)
}


// Scan implements sql.Scanner
func (a *AttributesMap) Scan(value interface{}) error {
    b, ok := value.([]byte)
    if !ok {
        return errors.New("type assertion to []byte failed")
    }
    return json.Unmarshal(b, a)
}

type Product struct {
	ID          string                 `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"` // Override gorm.Model ID
	MerchantID  string                 `gorm:"type:uuid;not null;index"` // FK to merchants.id (UUID string)
	Name        string                 `gorm:"size:255;not null" json:"name"`
	Description string                 `gorm:"type:text" json:"description"`
	SKU         string                 `gorm:"size:100;unique;not null;index" json:"sku"`
	Price       float64                `gorm:"type:decimal(10,2);not null" json:"price"`
	CategoryID  uint                 `gorm:"type:int;index" json:"category_id"` // Changed to string for UUID; revert to uint if numeric
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"` 

	Merchant    Merchant    `gorm:"foreignKey:MerchantID;references:id"` // Ensure Merchant.ID is string
	Category    Category    `gorm:"foreignKey:CategoryID"` // Adjust if Category.ID is uint
	Variants    []Variant   `gorm:"foreignKey:ProductID"`  // Keep relational
	Media       []Media     `gorm:"foreignKey:ProductID"`  // Keep relational
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}


// type Variant struct {
//     gorm.Model
//     ProductID  uint
//     Attributes map[string]string `gorm:"type:jsonb"`
//     Price      float64
//     SKU        string
// }


// //define merchant variant

// type Media struct {
//     gorm.Model
//     ProductID uint
//     URL       string
//     Type      string // image/video
//}

type Variant struct {
	gorm.Model
	ID        string             `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"` // Add UUID PK
	ProductID string             `gorm:"type:uuid;not null;index"` // Fixed: string for UUID (references products.id)
	//Attributes map[string]string `gorm:"type:jsonb;default:'{}'"`
	Attributes AttributesMap `gorm:"type:jsonb;default:'{}'"`
	Price     float64            `gorm:"type:decimal(10,2);not null"`
	SKU       string             `gorm:"size:100;unique;not null;index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (v *Variant) BeforeCreate(tx *gorm.DB) error {
	if v.ID == "" {
		v.ID = uuid.New().String()
	}
	return nil
}

type Media struct {
	gorm.Model
	ID        string  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"` // Add UUID PK
	ProductID string  `gorm:"type:uuid;not null;index"` // Fixed: string for UUID (references products.id)
	URL       string  `gorm:"size:500;not null"`
	Type      string  `gorm:"size:20;default:'image';not null"` // e.g., "image", "video"
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m *Media) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}