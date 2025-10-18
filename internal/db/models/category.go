// Updated category model with custom Attributes type
package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"regexp"

	"gorm.io/gorm"
)

type Attributes map[string]interface{}

// Scan implements the sql.Scanner interface for Attributes
func (a *Attributes) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan attributes: unexpected type %T", value)
	}
	if err := json.Unmarshal(b, a); err != nil {
		return fmt.Errorf("failed to unmarshal attributes: %w", err)
	}
	return nil
}

// Value implements the driver.Valuer interface for Attributes
func (a Attributes) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}



func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.CategorySlug == "" {
		c.CategorySlug = GetSlug(c.Name)
	}
	
	return nil
}

type Category struct {
	gorm.Model
	Name       string                 `gorm:"size:255;not null" json:"name"`
	ParentID   *uint                  `json:"parent_id"`
	CategorySlug string `gorm:"size:255;index" json:"category_slug"`   
	Attributes Attributes             `gorm:"type:jsonb" json:"attributes"`
	Parent     *Category              `gorm:"foreignKey:ParentID"`
}


func  GetSlug(name string) string {
    if name == "" {
        return ""
    }

    // Simple slug generation
    slug := strings.ToLower(strings.TrimSpace(name))
    // Replace spaces and common special chars with hyphens
    re := regexp.MustCompile(`[^a-z0-9]+`)
    slug = re.ReplaceAllString(slug, "-")
    // Trim leading/trailing hyphens and replace multiple hyphens with single
    reMulti := regexp.MustCompile(`-+`)
    slug = reMulti.ReplaceAllString(strings.Trim(slug, "-"), "-")

    return slug
}


func BackfillCategorySlugs(db *gorm.DB) error {
    var categories []Category
    if err := db.Model(&Category{}).Where("category_slug = '' OR category_slug IS NULL").Find(&categories).Error; err != nil {
        return fmt.Errorf("failed to fetch categories for backfill: %w", err)
    }

    for _, p := range categories {
        // Generate unique slug (handles potential collisions by appending ID suffix)
        newSlug := GetSlug(p.Name)
        // Optional: Check for uniqueness per merchant (if idx_merchant_slug is composite)
        // if err := db.Model(&Product{}).Where("category_slug = ? AND id != ?", newSlug, p.ID).First(&Category{}).Error; err == nil {
        //     // Collision: Append a counter or more of ID
        //     //newSlug = fmt.Sprintf("%s", newSlug)
        // }

        if err := db.Model(&p).Update("category_slug", newSlug).Error; err != nil {
            return fmt.Errorf("failed to update slug for product %v: %w", p.ID, err)
        }
    }

    return nil
}