package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// AttributesMap for JSONB

type DiscountType string

const (
	DiscountTypeFixed      DiscountType = "fixed"      // e.g., N5 off
	DiscountTypePercentage DiscountType = "percentage" // e.g., 10% off
	DiscountTypeNone       DiscountType = ""           // No discount
)

func (dt *DiscountType) Scan(value interface{}) error {
	*dt = DiscountType(value.(string))
	return nil
}

func (dt DiscountType) Value() (driver.Value, error) {
	return string(dt), nil
}

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

type Product struct {
	ID          string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	MerchantID  string `gorm:"type:uuid;not null;index" json:"merchant_id"`
	Name        string `gorm:"size:255;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	SKU         string `gorm:"size:100;unique;not null;index" json:"sku"`

	BasePrice       decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"base_price"`
	Discount        decimal.Decimal `gorm:"type:decimal(10,2);not null;default:0.00" json:"discount"`  // NEW: Discount amount
	DiscountType    DiscountType    `gorm:"type:varchar(20);not null;default:''" json:"discount_type"` // NEW: fixed/percentage
	FinalPrice      decimal.Decimal `gorm:"type:decimal(10,2);not null;default:0.00" json:"final_price"`
	CategoryID      uint            `gorm:"type:int;index" json:"category_id"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	DeletedAt       gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`                                          // Soft deletes for recovery
	//Slug            string          `gorm:"size:255;not null;uniqueIndex:idx_merchant_slug" json:"slug"`
	Slug string `gorm:"size:255;index" json:"slug"`                // Add this
	Merchant        Merchant        `gorm:"foreignKey:MerchantID;references:MerchantID;constraint:OnDelete:RESTRICT"`   // Belongs to Merchant, no cascade to protect merchants
	Category        Category        `gorm:"foreignKey:CategoryID;constraint:OnDelete:RESTRICT"`                         // Belongs to Category
	Variants        []Variant       `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"variants,omitempty"` // Has many Variants, cascade delete
	Media           []Media         `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"media,omitempty"`    // Has many Media, cascade delete
	SimpleInventory *Inventory      `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`                           // Has one optional SimpleInventory for non-variant products
	Wishlists       []UserWishlist  `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`                           // Has many UserWishlists
	Reviews         []Review        `gorm:"foreignKey:ProductID"`
}

// func (p *Product) BeforeCreate(tx *gorm.DB) error {
// 	if p.ID == "" {
// 		p.ID = uuid.New().String()
// 	}
// 	return nil
// }

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	if p.Slug == "" {
		p.Slug = GenerateSlug(p.Name, p.ID)
	}
	p.ComputeFinalPrice()
	return nil
}

// Helper function to generate slug
func GenerateSlug(name, id string) string {
	// Convert to lowercase, replace spaces with hyphens, remove special chars
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(slug, "")
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	
	// Add UUID suffix to ensure uniqueness if needed
	if slug == "" {
		slug = id[:8]
	} else {
		slug = slug + "-" + id[:8]
	}
	return slug
}


func (p *Product) BeforeUpdate(tx *gorm.DB) error {
	p.ComputeFinalPrice()
	return nil
}

func (p *Product) ComputeFinalPrice() {
	if p.DiscountType == DiscountTypePercentage && !p.Discount.Equal(decimal.Zero) {
		// e.g., 10% off = BasePrice * (1 - Discount/100)
		discountFraction := p.Discount.Div(decimal.NewFromInt(100))
		p.FinalPrice = p.BasePrice.Mul(decimal.NewFromInt(1).Sub(discountFraction))
	} else if p.DiscountType == DiscountTypeFixed && !p.Discount.Equal(decimal.Zero) {
		// e.g., $5 off = BasePrice - Discount
		p.FinalPrice = p.BasePrice.Sub(p.Discount)
	} else {
		p.FinalPrice = p.BasePrice
	}
	// Ensure non-negative
	if p.FinalPrice.LessThan(decimal.Zero) {
		p.FinalPrice = decimal.Zero
	}
}

type Variant struct {
	ID              string          `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ProductID       string          `gorm:"type:uuid;not null;index" json:"product_id"`
	SKU             string          `gorm:"size:100;unique;not null;index" json:"sku"`
	PriceAdjustment decimal.Decimal `gorm:"type:decimal(10,2);not null;default:0.00" json:"price_adjustment"`
	TotalPrice      decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"total_price"`              // Computed: BasePrice + PriceAdjustment
	Discount        decimal.Decimal `gorm:"type:decimal(10,2);not null;default:0.00" json:"discount"`    // NEW
	DiscountType    DiscountType    `gorm:"type:varchar(20);not null;default:''" json:"discount_type"`   // NEW
	FinalPrice      decimal.Decimal `gorm:"type:decimal(10,2);not null;default:0.00" json:"final_price"` // NEW: TotalPrice - Discount
	Attributes      AttributesMap   `gorm:"type:jsonb;default:'{}'" json:"attributes"`                   // Use map for simplicity; can change to custom AttributesMap if needed
	IsActive        bool            `gorm:"default:true" json:"is_active"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	DeletedAt       gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"` // Soft deletes for recovery

	Product   Product   `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"` // Belongs to Product, cascade from parent
	Inventory Inventory `gorm:"foreignKey:VariantID;constraint:OnDelete:CASCADE"` // Has one Inventory, cascade delete

}

type Review struct {
	ID        uint      `gorm:"primaryKey"`
	ProductID string    `gorm:"type:uuid;index;not null"`
	UserID    uint      `gorm:"index;not null"`
	Rating    int       `gorm:"not null;check:rating >= 1 AND rating <= 5"`
	Comment   string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Product   Product   `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type UserWishlist struct {
	UserID    uint   `gorm:"primaryKey" json:"user_id"`
	ProductID string `gorm:"primaryKey;type:uuid" json:"product_id"`

	AddedAt time.Time `json:"added_at"`

	User    User    `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Product Product `gorm:"foreignKey:ProductID;references:ID;constraint:OnDelete:CASCADE"`
}

func (uw *UserWishlist) BeforeCreate(tx *gorm.DB) error {
	uw.AddedAt = time.Now()
	return nil
}

// func (v *Variant) BeforeCreate(tx *gorm.DB) error {
// 	if v.ID == "" {
// 		v.ID = uuid.New().String()
// 	}
// 	// Fetch Product to compute TotalPrice
// 	var product Product
// 	if err := tx.Where("id = ?", v.ProductID).First(&product).Error; err != nil {
// 		return err
// 	}
// 	v.TotalPrice = product.BasePrice.Add(v.PriceAdjustment)
// 	return nil
// }

// func (v *Variant) BeforeUpdate(tx *gorm.DB) error {
// 	// Recompute TotalPrice on update
// 	var product Product
// 	if err := tx.Where("id = ?", v.ProductID).First(&product).Error; err != nil {
// 		return err
// 	}
// 	v.TotalPrice = product.BasePrice.Add(v.PriceAdjustment)
// 	return nil
// }

func (v *Variant) BeforeCreate(tx *gorm.DB) error {
	if v.ID == "" {
		v.ID = uuid.New().String()
	}
	// Fetch Product to compute TotalPrice and FinalPrice
	var product Product
	if err := tx.Where("id = ?", v.ProductID).First(&product).Error; err != nil {
		return err
	}
	v.TotalPrice = product.BasePrice.Add(v.PriceAdjustment)
	v.computeFinalPrice()
	return nil
}

func (v *Variant) BeforeUpdate(tx *gorm.DB) error {
	// Recompute TotalPrice and FinalPrice
	var product Product
	if err := tx.Where("id = ?", v.ProductID).First(&product).Error; err != nil {
		return err
	}
	v.TotalPrice = product.BasePrice.Add(v.PriceAdjustment)
	v.computeFinalPrice()
	return nil
}

func (v *Variant) computeFinalPrice() {
	if v.DiscountType == DiscountTypePercentage && !v.Discount.Equal(decimal.Zero) {
		discountFraction := v.Discount.Div(decimal.NewFromInt(100))
		v.FinalPrice = v.TotalPrice.Mul(decimal.NewFromInt(1).Sub(discountFraction))
	} else if v.DiscountType == DiscountTypeFixed && !v.Discount.Equal(decimal.Zero) {
		v.FinalPrice = v.TotalPrice.Sub(v.Discount)
	} else {
		v.FinalPrice = v.TotalPrice
	}
	if v.FinalPrice.LessThan(decimal.Zero) {
		v.FinalPrice = decimal.Zero
	}
}

type Media struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ProductID string    `gorm:"type:uuid;index" json:"product_id"`
	URL       string    `gorm:"not null" json:"url"`
	Type      MediaType `gorm:"not null" json:"type"` // enum: image, video
	PublicID  string    `gorm:"index" json:"public_id"`
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




func BackfillProductSlugs(db *gorm.DB) error {
    var products []Product
    if err := db.Model(&Product{}).Where("slug = '' OR slug IS NULL").Find(&products).Error; err != nil {
        return fmt.Errorf("failed to fetch products for backfill: %w", err)
    }

    for _, p := range products {
        // Generate unique slug (handles potential collisions by appending ID suffix)
        newSlug := GenerateSlug(p.Name, p.ID)
        // Optional: Check for uniqueness per merchant (if idx_merchant_slug is composite)
        if err := db.Model(&Product{}).Where("merchant_id = ? AND slug = ? AND id != ?", p.MerchantID, newSlug, p.ID).First(&Product{}).Error; err == nil {
            // Collision: Append a counter or more of ID
            newSlug = fmt.Sprintf("%s-%s", newSlug, uuid.NewString()[:4])
        }

        if err := db.Model(&p).Update("slug", newSlug).Error; err != nil {
            return fmt.Errorf("failed to update slug for product %s: %w", p.ID, err)
        }
    }

    return nil
}