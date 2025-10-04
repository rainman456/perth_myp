// Updated category model with custom Attributes type
package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

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