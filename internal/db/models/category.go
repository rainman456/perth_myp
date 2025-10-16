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

type Category struct {
	gorm.Model
	Name       string                 `gorm:"size:255;not null" json:"name"`
	ParentID   *uint                  `json:"parent_id"`
	Attributes Attributes             `gorm:"type:jsonb" json:"attributes"`
	Parent     *Category              `gorm:"foreignKey:ParentID"`
}


func (c *Category) Slug() string {
    if c.Name == "" {
        return ""
    }

    // Simple slug generation
    slug := strings.ToLower(strings.TrimSpace(c.Name))
    // Replace spaces and common special chars with hyphens
    re := regexp.MustCompile(`[^a-z0-9]+`)
    slug = re.ReplaceAllString(slug, "-")
    // Trim leading/trailing hyphens and replace multiple hyphens with single
    reMulti := regexp.MustCompile(`-+`)
    slug = reMulti.ReplaceAllString(strings.Trim(slug, "-"), "-")

    return slug
}