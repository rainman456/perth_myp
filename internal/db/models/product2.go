package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// AttributesMap for JSONB
type AttributesMap map[string]string

func (a *AttributesMap) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, a)
}

func (a AttributesMap) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// MediaType enum-like type
type MediaType string

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
)

func (mt *MediaType) Scan(value interface{}) error {
	*mt = MediaType(value.(string))
	return nil
}

func (mt MediaType) Value() (driver.Value, error) {
	return string(mt), nil
}


// type Product struct {
// 	ID          string          `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
// 	MerchantID  string          `gorm:"type:uuid;not null;index" json:"merchant_id"`
// 	Name        string          `gorm:"size:255;not null" json:"name"`
// 	Description string          `gorm:"type:text" json:"description"`
// 	SKU         string          `gorm:"size:100;unique;not null;index" json:"sku"`
// 	BasePrice   decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"base_price"`
// 	CategoryID  uint            `gorm:"type:int;index" json:"category_id"`
// 	CreatedAt   time.Time       `json:"created_at"`
// 	UpdatedAt   time.Time       `json:"updated_at"`
// 	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"deleted_at"`

// 	Merchant        Merchant         `gorm:"foreignKey:MerchantID;references:id"`
// 	Category        Category         `gorm:"foreignKey:CategoryID"`
// 	Variants        []Variant        `gorm:"foreignKey:ProductID" json:"variants,omitempty"`
// 	Media           []Media          `gorm:"foreignKey:ProductID" json:"media,omitempty"`
// 	SimpleInventory *VendorInventory `gorm:"foreignKey:ProductID"` // Only for simple products (no variants)
// }



type Product struct {
	ID          string          `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	MerchantID  string          `gorm:"type:uuid;not null;index" json:"merchant_id"`
	Name        string          `gorm:"size:255;not null" json:"name"`
	Description string          `gorm:"type:text" json:"description"`
	SKU         string          `gorm:"size:100;unique;not null;index" json:"sku"`
	BasePrice   decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"base_price"`
	CategoryID  uint            `gorm:"type:int;index" json:"category_id"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"` // Soft deletes for recovery

	Merchant        Merchant       `gorm:"foreignKey:MerchantID;references:MerchantID;constraint:OnDelete:RESTRICT"` // Belongs to Merchant, no cascade to protect merchants
	Category        Category       `gorm:"foreignKey:CategoryID;constraint:OnDelete:RESTRICT"` // Belongs to Category
	Variants        []Variant      `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"variants,omitempty"` // Has many Variants, cascade delete
	Media           []Media        `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"media,omitempty"` // Has many Media, cascade delete
	SimpleInventory *Inventory     `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"` // Has one optional SimpleInventory for non-variant products
}


func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

// type Variant struct {
// 	ID              string          `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
// 	ProductID       string          `gorm:"type:uuid;not null;index" json:"product_id"`
// 	SKU             string          `gorm:"size:100;unique;not null;index" json:"sku"`
// 	PriceAdjustment decimal.Decimal `gorm:"type:decimal(10,2);not null;default:0.00" json:"price_adjustment"`
// 	TotalPrice      decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"total_price"` // Computed: BasePrice + PriceAdjustment
// 	Attributes      AttributesMap   `gorm:"type:jsonb;default:'{}'" json:"attributes"`
// 	IsActive        bool            `gorm:"default:true" json:"is_active"`
// 	CreatedAt       time.Time       `json:"created_at"`
// 	UpdatedAt       time.Time       `json:"updated_at"`

// 	Product   Product         `gorm:"foreignKey:ProductID"`
// 	Inventory VendorInventory `gorm:"foreignKey:VariantID"`
// }




type Variant struct {
	ID              string          `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ProductID       string          `gorm:"type:uuid;not null;index" json:"product_id"`
	SKU             string          `gorm:"size:100;unique;not null;index" json:"sku"`
	PriceAdjustment decimal.Decimal `gorm:"type:decimal(10,2);not null;default:0.00" json:"price_adjustment"`
	TotalPrice      decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"total_price"` // Computed: BasePrice + PriceAdjustment
	Attributes      AttributesMap `gorm:"type:jsonb;default:'{}'" json:"attributes"` // Use map for simplicity; can change to custom AttributesMap if needed
	IsActive        bool            `gorm:"default:true" json:"is_active"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"` // Soft deletes for recovery

	Product   Product   `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"` // Belongs to Product, cascade from parent
	Inventory Inventory `gorm:"foreignKey:VariantID;constraint:OnDelete:CASCADE"` // Has one Inventory, cascade delete
}

func (v *Variant) BeforeCreate(tx *gorm.DB) error {
	if v.ID == "" {
		v.ID = uuid.New().String()
	}
	// Fetch Product to compute TotalPrice
	var product Product
	if err := tx.Where("id = ?", v.ProductID).First(&product).Error; err != nil {
		return err
	}
	v.TotalPrice = product.BasePrice.Add(v.PriceAdjustment)
	return nil
}

func (v *Variant) BeforeUpdate(tx *gorm.DB) error {
	// Recompute TotalPrice on update
	var product Product
	if err := tx.Where("id = ?", v.ProductID).First(&product).Error; err != nil {
		return err
	}
	v.TotalPrice = product.BasePrice.Add(v.PriceAdjustment)
	return nil
}

// type VendorInventory struct {
// 	ID               string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
// 	VariantID        string         `gorm:"type:uuid;not null;unique;index" json:"variant_id"`
// 	MerchantID       string         `gorm:"type:uuid;not null;index" json:"merchant_id"`
// 	ProductID        *string        `gorm:"type:uuid;index"` // Nullable: For simple products
// 	Quantity         int            `gorm:"default:0;not null;check:quantity >= 0" json:"quantity"`
// 	ReservedQuantity int            `gorm:"default:0;check:reserved_quantity >= 0" json:"reserved_quantity"`
// 	LowStockThreshold int           `gorm:"default:10" json:"low_stock_threshold"`
// 	BackorderAllowed bool           `gorm:"default:false" json:"backorder_allowed"`
// 	CreatedAt        time.Time      `json:"created_at"`
// 	UpdatedAt        time.Time      `json:"updated_at"`

// 	Variant  *Variant `gorm:"foreignKey:VariantID"`
// 	Product  *Product `gorm:"foreignKey:ProductID"`
// 	Merchant Merchant `gorm:"foreignKey:MerchantID"`
// }

// func (vi *VendorInventory) BeforeCreate(tx *gorm.DB) error {
// 	if vi.ID == "" {
// 		vi.ID = uuid.New().String()
// 	}
// 	if (vi.VariantID != "" && vi.ProductID != nil) || (vi.VariantID == "" && vi.ProductID == nil) {
// 		return errors.New("exactly one of VariantID or ProductID must be set")
// 	}
// 	return nil
// }

// type Media struct {
// 	ID        string    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
// 	ProductID string    `gorm:"type:uuid;not null;index" json:"product_id"`
// 	URL       string    `gorm:"size:500;not null" json:"url"`
// 	Type      MediaType `gorm:"type:varchar(20);default:image;not null" json:"type"`
// 	CreatedAt time.Time `json:"created_at"`
// 	UpdatedAt time.Time `json:"updated_at"`

// 	Product Product `gorm:"foreignKey:ProductID"`
// }


type Media struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ProductID string    `gorm:"type:uuid;index" json:"product_id"`
	URL       string    `gorm:"not null" json:"url"`
	Type      MediaType `gorm:"not null" json:"type"` // enum: image, video
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Product Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"` // Belongs to Product (bidirectional for easier queries)
}


func (m *Media) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}




func (p *Product) GenerateSKU(merchantID string) {
	base := strings.ToUpper(strings.ReplaceAll(p.Name, " ", "-"))
	unique := strings.ToUpper(uuid.NewString()[:8])
	p.SKU = fmt.Sprintf("%s-%s-%s", merchantID[:4], base, unique) // Prefix min(4, len(merchantID))
}

// GenerateSKU auto-generates SKU for the variant based on product SKU and attributes
func (v *Variant) GenerateSKU(productSKU string) {
	if len(v.Attributes) == 0 {
		v.SKU = productSKU + "-DEFAULT"
		return
	}

	// Sort keys for consistent order
	keys := make([]string, 0, len(v.Attributes))
	for k := range v.Attributes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	attrStr := ""
	for _, k := range keys {
		val := v.Attributes[k]
		attrStr += fmt.Sprintf("-%s-%s", strings.ToUpper(k), strings.ToUpper(strings.ReplaceAll(val, " ", "-")))
	}
	v.SKU = productSKU + attrStr
}