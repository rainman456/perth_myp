# Codebase Analysis: internal
Generated: 2025-10-02 15:26:53
---

## 📂 Project Structure
```tree
📁 internal
├── db/
│   ├── models/
│   │   ├── bank_details.go
│   │   ├── cart.go
│   │   ├── cart_item.go
│   │   ├── category.go
│   │   ├── disputes.go
│   │   ├── inventory.go
│   │   ├── merchant.go
│   │   ├── order.go
│   │   ├── order_item.go
│   │   ├── order_merchant_split.go
│   │   ├── payment.go
│   │   ├── payout.go
│   │   ├── product.go
│   │   ├── product2.go
│   │   └── user.go
│   ├── repositories/
│   │   ├── cart_item_repository.go
│   │   ├── cart_repository.go
│   │   ├── category_repositry.go
│   │   ├── disputes_repository.go
│   │   ├── inventory_repository.go
│   │   ├── merchant_repository.go
│   │   ├── order_item_repository.go
│   │   ├── order_repository.go
│   │   ├── payment_repository.go
│   │   ├── payout_repository.go
│   │   ├── product_repo.go
│   │   ├── product_repositry.go
│   │   ├── user_repository.go
│   │   └── variant_repository.go
│   └── db.go
└── services/
    ├── cart/
    │   └── cart_service.go
    ├── dispute/
    │   └── dispute_service.go
    ├── merchant/
    │   └── merchant_service.go
    ├── notifications/
    │   └── notifcation_service.go
    ├── order/
    │   └── order_service.go
    ├── payment/
    │   └── payment_service.go
    ├── payout/
    │   └── payout_service.go
    ├── pricing/
    │   └── pricing_service.go
    ├── product/
    │   └── proudct_service.go
    ├── return_request/
    │   └── return_request_service.go
    └── user/
        └── user_service.go
```
---

## 📄 File Contents
### db/db.go
- Size: 5.76 KB
- Lines: 217
- Last Modified: 2025-09-30 12:28:15

```go
package db

import (
	//"api-customer-merchant/internal/db/models"
	//"api-customer-merchant/internal/db/models"
	//"api-customer-merchant/internal/db/models"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

/*
func Connect() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN environment variable not set")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully")
}

func AutoMigrate() {
	//    err := DB.AutoMigrate(
	//        &models.User{},
	// 	   &models.Merchant{},
	//        // Add other models here when implemented
	//    )
	//    if err != nil {
	//        log.Fatalf("Failed to auto-migrate: %v", err)
	//    }

}
*/

func Connect() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN environment variable not set")
	}

	var err error
	//DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true,})
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{PrepareStmt: false})
	//DB = DB.Debug()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	log.Println("Database connected successfully")
}

func AutoMigrate() {
	// Run AutoMigrate with all models
	if err := DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		log.Fatal("Failed to enable uuid-ossp extension:", err)

	}

	//DB = DB.Session(&gorm.Session{DisableForeignKeyConstraintWhenMigrating: true})
	log.Println("Starting AutoMigrate...")

	// Independent tables (no incoming FKs or self-referential only)
	// log.Println("Migrating Category...")
	// if err := DB.AutoMigrate(&models.Category{}); err != nil {
	//     log.Printf("Failed to migrate Category: %v", err)
	//     return
	// }

	//  log.Println("Migrating MerchantApplication...")
	//  if err := DB.AutoMigrate(&models.MerchantApplication{}); err != nil {
	//     log.Printf("Failed to migrate MerchantApplication: %v", err)
	//      return
	//  }

	// log.Println("Migrating Merchant...")
	// if err := DB.AutoMigrate(&models.Merchant{}); err != nil {
	//     log.Printf("Failed to migrate Merchant: %v", err)
	//     return
	// }

	//  log.Println("Migrating User...")
	//  if err := DB.AutoMigrate(&models.User{}); err != nil {
	//      log.Printf("Failed to migrate User: %v", err)
	//      return
	//  }

	// Tables depending on Merchant/Category/User
	//    log.Println("Migrating Product ecosystem (Product, Variant, Media)...")
	//  if err := DB.AutoMigrate(&models.Product{}, &models.Variant{}, &models.Media{}, &models.VendorInventory{}); err != nil {
	//      log.Printf("Failed to migrate Product/Variant/Media: %v", err)
	//      return
	//  }

	// log.Println("Migrating Inventory...")
	// if err := DB.AutoMigrate(&models.Inventory{}); err != nil {
	//     log.Printf("Failed to migrate Inventory: %v", err)
	//     return
	// }

	// log.Println("Migrating Promotion...")
	// if err := DB.AutoMigrate(&models.Promotion{}); err != nil {
	//     log.Printf("Failed to migrate Promotion: %v", err)
	//     return
	// }

	//  log.Println("Migrating Cart...")
	//  if err := DB.AutoMigrate(&models.Cart{}); err != nil {
	//      log.Printf("Failed to migrate Cart: %v", err)
	//      return
	//  }

	//  log.Println("Migrating CartItem...")
	//  if err := DB.AutoMigrate(&models.CartItem{}); err != nil {
	//      log.Printf("Failed to migrate CartItem: %v", err)
	//      return
	//  }

	//  log.Println("Migrating Order...")
	//  if err := DB.AutoMigrate(&models.Order{}); err != nil {
	//      log.Printf("Failed to migrate Order: %v", err)
	//      return
	//  }

	//  log.Println("Migrating OrderItem...")
	//  if err := DB.AutoMigrate(&models.OrderItem{}); err != nil {
	//      log.Printf("Failed to migrate OrderItem: %v", err)
	//      return
	//  }

	// log.Println("Migrating Payment...")
	// if err := DB.AutoMigrate(&models.Payment{}); err != nil {
	//     log.Printf("Failed to migrate Payment: %v", err)
	//     return
	// }

	// log.Println("Migrating Payout...")
	// if err := DB.AutoMigrate(&models.Payout{}); err != nil {
	//     log.Printf("Failed to migrate Payout: %v", err)
	//     return
	// }

	// log.Println("Migrating ReturnRequest...")
	// if err := DB.AutoMigrate(&models.ReturnRequest{}); err != nil {
	//     log.Printf("Failed to migrate ReturnRequest: %v", err)
	//     return
	// }

	err := DB.AutoMigrate(
	//&models.User{},
	//&models.MerchantApplication{},
	//&models.Product{},
	//&models.Variant{},
	//&models.Media{},
	// &models.Cart{},
	// &models.Order{},
	 //&models.OrderItem{},
	 //&models.CartItem{},
	//&models.Category{},
	 //&models.Inventory{},
	// &models.Promotion{},
	// &models.ReturnRequest{},
	// &models.Payout{},
	)

	if err != nil {
		log.Fatalf("Failed to auto-migrate: %v", err)
	}

	// Get the underlying SQL database connection
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL DB: %v", err)
	}

	// Close all connections to clear cached plans
	if err := sqlDB.Close(); err != nil {
		log.Printf("Failed to close connections: %v", err)
	}

	// Reconnect to ensure fresh connections
	dsn := os.Getenv("DB_DSN")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to reconnect to database: %v", err)
	}

	// Reconfigure connection pool
	sqlDB, err = DB.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL DB after reconnect: %v", err)
	}
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	log.Println("Database migrated and reconnected successfully")
}

```

---
### db/models/bank_details.go
- Size: 0.62 KB
- Lines: 21
- Last Modified: 2025-09-30 12:28:15

```go
package models

   import (
       "time"
       "gorm.io/gorm"
   )

   type MerchantBankDetails struct {
       gorm.Model
       MerchantID     string `gorm:"type:uuid;uniqueIndex"`
       BankName       string
       BankCode       string `gorm:"size:5"`
       AccountNumber  string `gorm:"size:15"`
       AccountName    string
       RecipientCode  string `gorm:"size:50"`
       Currency       string `gorm:"size:3;default:'NGN'"`
       Status         string `gorm:"default:'pending'"`
       CreatedAt      time.Time
       UpdatedAt      time.Time
	   Merchant   Merchant `gorm:"foreignKey:MerchantID;references:MerchantID"`
   }
```

---
### db/models/cart.go
- Size: 1.74 KB
- Lines: 68
- Last Modified: 2025-09-30 12:22:22

```go
package models

import (
	"fmt"
	"gorm.io/gorm"
)

// CartStatus defines possible cart status values
type CartStatus string

const (
	CartStatusActive    CartStatus = "Active"
	CartStatusAbandoned CartStatus = "Abandoned"
	CartStatusConverted CartStatus = "Converted"
)

// Valid checks if the status is one of the allowed values
func (s CartStatus) Valid() error {
	switch s {
	case CartStatusActive, CartStatusAbandoned, CartStatusConverted:
		return nil
	default:
		return fmt.Errorf("invalid cart status: %s", s)
	}
}

type Cart struct {
	gorm.Model
	UserID     uint       `gorm:"not null" json:"user_id"`
	Status     CartStatus `gorm:"type:varchar(20);not null;default:'Active'" json:"status"`
	SubTotal   float64    `gorm:"-" json:"subtotal"` // Computed
	TaxTotal   float64    `gorm:"-" json:"tax_total"`
	ShipTotal  float64    `gorm:"-" json:"shipping_total"`
	GrandTotal float64    `gorm:"-" json:"grand_total"`
	User       User       `gorm:"foreignKey:UserID"`
	CartItems  []CartItem `gorm:"foreignKey:CartID"`
}

// BeforeCreate validates the Status field
func (c *Cart) BeforeCreate(tx *gorm.DB) error {
	if err := c.Status.Valid(); err != nil {
		return err
	}
	return nil
}

// BeforeUpdate validates the Status field
func (c *Cart) BeforeUpdate(tx *gorm.DB) error {
	if err := c.Status.Valid(); err != nil {
		return err
	}
	return nil
}

func (c *Cart) AfterFind(tx *gorm.DB) error {
	c.ComputeTotals()
	return nil
}

func (c *Cart) ComputeTotals() {
	c.SubTotal = 0
	for _, item := range c.CartItems {
		c.SubTotal += float64(item.Quantity) * (item.Product.BasePrice).InexactFloat64() // Assume BasePrice in Product
	}
	// Stub: c.TaxTotal = 0.1 * c.SubTotal // Or call pricing
	// c.ShipTotal = 10.00
	c.GrandTotal = c.SubTotal // + c.TaxTotal + c.ShipTotal
}

```

---
### db/models/cart_item.go
- Size: 0.56 KB
- Lines: 18
- Last Modified: 2025-09-30 12:22:22

```go
package models

import (
	"gorm.io/gorm"
)

type CartItem struct {
	gorm.Model
	CartID     uint     `gorm:"not null" json:"cart_id"`
	VariantID  *string  `gorm:"type:uuid;index"`
	ProductID  string   `gorm:"not null" json:"product_id"`
	Quantity   int      `gorm:"not null" json:"quantity"`
	MerchantID string   `gorm:"not null" json:"merchant_id"`
	Cart       Cart     `gorm:"foreignKey:CartID"`
	Product    Product  `gorm:"foreignKey:ProductID"`
	Merchant   Merchant `gorm:"foreignKey:MerchantID;references:MerchantID"`
	Variant    *Variant `gorm:"foreignKey:VariantID"`
}

```

---
### db/models/category.go
- Size: 0.34 KB
- Lines: 13
- Last Modified: 2025-09-30 12:28:15

```go
package models

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name       string                 `gorm:"size:255;not null" json:"name"`
	ParentID   *uint                  `json:"parent_id"`
	Attributes map[string]interface{} `gorm:"type:jsonb" json:"attributes"`
	Parent     *Category              `gorm:"foreignKey:ParentID"`
}

```

---
### db/models/disputes.go
- Size: 2.51 KB
- Lines: 58
- Last Modified: 2025-10-02 10:55:47

```go
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
	OrderItemID      string    `gorm:"type:uuid;not null" json:"order_item_id"`
	CustomerID        uint    `gorm:"not null" json:"customer_id"`
	Reason           string    `gorm:"type:text" json:"reason"`
	Status           string    `gorm:"type:varchar(255);default:'Pending'" json:"status"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	OrderItem        OrderItem `gorm:"foreignKey:OrderItemID"`
	Customer           User              `gorm:"foreignKey:CustomerID"`

}

// Settings model (matching TS settings)
type Settings struct {
	gorm.Model
	ID              string                 `gorm:"type:text;primaryKey;default:'global'" json:"id"`
	Fees            float64                `gorm:"type:decimal(10,2);not null;default:5.00" json:"fees"`
	TaxRate         float64                `gorm:"type:decimal(10,2);not null;default:0.00" json:"tax_rate"`
	ShippingOptions map[string]interface{} `gorm:"type:jsonb;not null" json:"shipping_options"`
}

```

---
### db/models/inventory.go
- Size: 2.92 KB
- Lines: 64
- Last Modified: 2025-09-30 12:22:22

```go
package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// type Inventory struct {
// 	gorm.Model
// 	ProductID         string   `gorm:"not null" json:"product_id"`
// 	StockQuantity     int    `gorm:"not null" json:"stock_quantity"`
// 	LowStockThreshold int    `gorm:"not null;default:10" json:"low_stock_threshold"`
// 	Product           Product `gorm:"foreignKey:ProductID"`
// }

// type VendorInventory struct {
// 	ID                string    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
// 	VariantID         string    `gorm:"type:uuid;not null;unique;index" json:"variant_id"`
// 	MerchantID        string    `gorm:"type:uuid;not null;index" json:"merchant_id"`
// 	ProductID         *string   `gorm:"type:uuid;index"` // Nullable: For simple products
// 	Quantity          int       `gorm:"default:0;not null;check:quantity >= 0" json:"quantity"`
// 	ReservedQuantity  int       `gorm:"default:0;check:reserved_quantity >= 0" json:"reserved_quantity"`
// 	LowStockThreshold int       `gorm:"default:10" json:"low_stock_threshold"`
// 	BackorderAllowed  bool      `gorm:"default:false" json:"backorder_allowed"`
// 	CreatedAt         time.Time `json:"created_at"`
// 	UpdatedAt         time.Time `json:"updated_at"`

// 	Variant  *Variant `gorm:"foreignKey:VariantID"`
// 	Product  *Product `gorm:"foreignKey:ProductID"`
// 	Merchant Merchant `gorm:"foreignKey:MerchantID"`
// }



type Inventory struct {
	ID                string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ProductID         *string    `gorm:"type:uuid;index" json:"product_id,omitempty"` // Optional: For simple products
	VariantID         *string    `gorm:"type:uuid;index" json:"variant_id,omitempty"` // Optional: For variants
	MerchantID        string    `gorm:"type:uuid;not null;index" json:"merchant_id"` // Required: Vendor-specific
	Quantity          int       `gorm:"default:0;not null;check:quantity >= 0" json:"quantity"`
	ReservedQuantity  int       `gorm:"default:0;not null;check:reserved_quantity >= 0" json:"reserved_quantity"`
	LowStockThreshold int       `gorm:"default:5" json:"low_stock_threshold"`
	BackorderAllowed  bool      `gorm:"default:false" json:"backorder_allowed"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	Product  *Product  `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"` // Optional belongs to Product
	Variant  *Variant  `gorm:"foreignKey:VariantID;constraint:OnDelete:CASCADE"` // Optional belongs to Variant
	Merchant Merchant  `gorm:"foreignKey:MerchantID;constraint:OnDelete:RESTRICT"` // Belongs to Merchant, no cascade
}


func (vi *Inventory) BeforeCreate(tx *gorm.DB) error {
	if vi.ID == "" {
		vi.ID = uuid.New().String()
	}
	if (vi.VariantID != nil && vi.ProductID != nil) || (vi.VariantID == nil && vi.ProductID == nil) {
		return errors.New("exactly one of VariantID or ProductID must be set")
	}
	return nil
}

```

---
### db/models/merchant.go
- Size: 10.44 KB
- Lines: 179
- Last Modified: 2025-09-30 12:28:15

```go
package models

import (
	"gorm.io/datatypes"
	"time"
)

/*
// MerchantStatus defines the possible status values for a merchant
type MerchantBasicInfo struct {
	StoreName     string `gorm:"column:store_name;size:255;not null" json:"store_name" validate:"required"`
	Name          string `gorm:"column:name;size:255;not null" json:"name" validate:"required"`
	PersonalEmail string `gorm:"column:personal_email;size:255;not null;unique" json:"personal_email" validate:"required,email"`
	WorkEmail     string `gorm:"column:work_email;size:255;not null;unique" json:"work_email" validate:"required,email"`
	PhoneNumber   string `gorm:"column:phone_number;size:50" json:"phone_number"`
}

// MerchantAddress holds address information as JSONB
type MerchantAddress struct {
	PersonalAddress datatypes.JSON `gorm:"column:personal_address;type:jsonb;not null" json:"personal_address" validate:"required"`
	WorkAddress     datatypes.JSON `gorm:"column:work_address;type:jsonb;not null" json:"work_address" validate:"required"`
}

// MerchantBusinessInfo holds business-related information
type MerchantBusinessInfo struct {
	BusinessType               string `gorm:"column:business_type;size:100" json:"business_type"`
	Website                    string `gorm:"column:website;size:255" json:"website"`
	BusinessDescription        string `gorm:"column:business_description;type:text" json:"business_description"`
	BusinessRegistrationNumber string `gorm:"column:business_registration_number;size:255;not null;unique" json:"business_registration_number" validate:"required"`
}

// MerchantDocuments holds document-related information
type MerchantDocuments struct {
	StoreLogoURL                   string `gorm:"column:store_logo_url;size:255" json:"store_logo_url"`
	BusinessRegistrationCertificate string `gorm:"column:business_registration_certificate;size:255" json:"business_registration_certificate"`
}

// MerchantApplication holds the information required for a merchant onboarding application
type MerchantApplication struct {
	ID                string            `gorm:"primaryKey;column:id;type:uuid;default:uuid_generate_v4()" json:"id,omitempty"`
	MerchantBasicInfo                   `gorm:"embedded"`
	MerchantAddress                     `gorm:"embedded"`
	MerchantBusinessInfo                `gorm:"embedded"`
	MerchantDocuments                   `gorm:"embedded"`
	Status            string            `gorm:"column:status;type:varchar(20);default:pending;not null" json:"status"`
	CreatedAt         time.Time         `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time         `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (MerchantApplication) TableName() string {
	return "merchant_application"
}

// MerchantStatus defines the possible statuses for a merchant
type MerchantStatus string

const (
	MerchantStatusActive    MerchantStatus = "active"
	MerchantStatusSuspended MerchantStatus = "suspended"
)

// Merchant holds the active merchant account details after approval
type Merchant struct {
	ID                string            `gorm:"primaryKey;column:id;type:uuid;default:uuid_generate_v4()" json:"id,omitempty"`
	ApplicationID     string            `gorm:"column:application_id;type:uuid;not null;unique" json:"application_id"`
	MerchantID            string            `gorm:"column:merchant_id;type:uuid;not null;unique" json:"user_id"`
	MerchantBasicInfo                   `gorm:"embedded"`
	MerchantAddress                     `gorm:"embedded"`
	MerchantBusinessInfo                `gorm:"embedded"`
	MerchantDocuments                   `gorm:"embedded"`
	Password          string            `gorm:"column:password;size:255;not null" json:"password" validate:"required"`
	Status            MerchantStatus    `gorm:"column:status;type:varchar(20);default:active;index" json:"status"`
	CommissionTier    string            `gorm:"column:commission_tier;default:standard" json:"commission_tier"`
	CommissionRate    float64           `gorm:"column:commission_rate;default:5.00" json:"commission_rate"`
	AccountBalance    float64           `gorm:"column:account_balance;default:0.00" json:"account_balance"`
	TotalSales        float64           `gorm:"column:total_sales;default:0.00" json:"total_sales"`
	TotalPayouts      float64           `gorm:"column:total_payouts;default:0.00" json:"total_payouts"`
	StripeAccountID   string            `gorm:"column:stripe_account_id" json:"stripe_account_id"`
	PayoutSchedule    string            `gorm:"column:payout_schedule;default:weekly" json:"payout_schedule"`
	LastPayoutDate    *time.Time        `gorm:"column:last_payout_date" json:"last_payout_date"`
	Banner            string            `gorm:"column:banner;size:255" json:"banner"`
	Policies          datatypes.JSON    `gorm:"column:policies;type:jsonb" json:"policies"`
	CreatedAt         time.Time         `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time         `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	//Products          []Product         `gorm:"foreignKey:MerchantID" json:"products,omitempty"`
	//CartItems         []CartItem        `gorm:"foreignKey:MerchantID" json:"cart_items,omitempty"`
	//OrderItems        []OrderItem       `gorm:"foreignKey:MerchantID" json:"order_items,omitempty"`
	//Payouts           []Payout          `gorm:"foreignKey:MerchantID" json:"payouts,omitempty"`
}

func (Merchant) TableName() string {
	return "merchant"
}
*/

type MerchantBasicInfo struct {
	StoreName     string `gorm:"column:store_name;size:255;not null" json:"store_name" validate:"required"`
	Name          string `gorm:"column:name;size:255;not null" json:"name" validate:"required"`
	PersonalEmail string `gorm:"column:personal_email;size:255;not null;unique" json:"personal_email" validate:"required,email"`
	WorkEmail     string `gorm:"column:work_email;size:255;not null;unique" json:"work_email" validate:"required,email"`
	PhoneNumber   string `gorm:"column:phone_number;size:50" json:"phone_number"`
}

// MerchantAddress holds address information as JSONB
type MerchantAddress struct {
	PersonalAddress datatypes.JSON `gorm:"column:personal_address;type:jsonb;not null" json:"personal_address" validate:"required"`
	WorkAddress     datatypes.JSON `gorm:"column:work_address;type:jsonb;not null" json:"work_address" validate:"required"`
}

// MerchantBusinessInfo holds business-related information
type MerchantBusinessInfo struct {
	BusinessType               string `gorm:"column:business_type;size:100" json:"business_type"`
	Website                    string `gorm:"column:website;size:255" json:"website"`
	BusinessDescription        string `gorm:"column:business_description;type:text" json:"business_description"`
	BusinessRegistrationNumber string `gorm:"column:business_registration_number;size:255;not null;unique" json:"business_registration_number" validate:"required"`
}

// MerchantDocuments holds document-related information
type MerchantDocuments struct {
	StoreLogoURL                    string `gorm:"column:store_logo_url;size:255" json:"store_logo_url"`
	BusinessRegistrationCertificate string `gorm:"column:business_registration_certificate;size:255" json:"business_registration_certificate"`
}

// MerchantApplication holds the information required for a merchant onboarding application
type MerchantApplication struct {
	ID                   string `gorm:"primaryKey;column:id;type:uuid;default:uuid_generate_v4()" json:"id,omitempty"`
	MerchantBasicInfo    `gorm:"embedded"`
	MerchantAddress      `gorm:"embedded"`
	MerchantBusinessInfo `gorm:"embedded"`
	MerchantDocuments    `gorm:"embedded"`
	Status               string    `gorm:"column:status;type:varchar(20);default:pending;not null" json:"status"`
	CreatedAt            time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (MerchantApplication) TableName() string {
	return "merchant_application"
}

// MerchantStatus defines the possible statuses for a merchant
type MerchantStatus string

const (
	MerchantStatusActive    MerchantStatus = "active"
	MerchantStatusSuspended MerchantStatus = "suspended"
)

// Merchant holds the active merchant account details after approval
type Merchant struct {
	ID                   string `gorm:"primaryKey;column:id;type:uuid;default:uuid_generate_v4()" json:"id,omitempty"`
	ApplicationID        string `gorm:"column:application_id;type:uuid;not null;unique" json:"application_id"`
	MerchantID           string `gorm:"column:merchant_id;type:uuid;not null;unique" json:"user_id"`
	MerchantBasicInfo    `gorm:"embedded"`
	MerchantAddress      `gorm:"embedded"`
	MerchantBusinessInfo `gorm:"embedded"`
	MerchantDocuments    `gorm:"embedded"`
	Password             string         `gorm:"column:password;size:255;not null" json:"password" validate:"required"`
	Status               MerchantStatus `gorm:"column:status;type:varchar(20);default:active;index" json:"status"`
	CommissionTier       string         `gorm:"column:commission_tier;default:standard" json:"commission_tier"`
	CommissionRate       float64        `gorm:"column:commission_rate;default:5.00" json:"commission_rate"`
	AccountBalance       float64        `gorm:"column:account_balance;default:0.00" json:"account_balance"`
	TotalSales           float64        `gorm:"column:total_sales;default:0.00" json:"total_sales"`
	TotalPayouts         float64        `gorm:"column:total_payouts;default:0.00" json:"total_payouts"`
	StripeAccountID      string         `gorm:"column:stripe_account_id" json:"stripe_account_id"`
	PayoutSchedule       string         `gorm:"column:payout_schedule;default:weekly" json:"payout_schedule"`
	LastPayoutDate       *time.Time     `gorm:"column:last_payout_date" json:"last_payout_date"`
	Banner               string         `gorm:"column:banner;size:255" json:"banner"`
	Policies             datatypes.JSON `gorm:"column:policies;type:jsonb" json:"policies"`
	CreatedAt            time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	//Products          []Product         `gorm:"foreignKey:MerchantID" json:"products,omitempty"`
	//CartItems         []CartItem        `gorm:"foreignKey:MerchantID" json:"cart_items,omitempty"`
	//OrderItems        []OrderItem       `gorm:"foreignKey:MerchantID" json:"order_items,omitempty"`
	//Payouts           []Payout          `gorm:"foreignKey:MerchantID" json:"payouts,omitempty"`
}

func (Merchant) TableName() string {
	return "merchant"
}

```

---
### db/models/order.go
- Size: 2.08 KB
- Lines: 70
- Last Modified: 2025-09-30 12:28:15

```go
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

```

---
### db/models/order_item.go
- Size: 2.44 KB
- Lines: 67
- Last Modified: 2025-09-30 12:22:22

```go
package models

import (
	"fmt"
	"gorm.io/gorm"
)

// FulfillmentStatus defines possible fulfillment status values
type FulfillmentStatus string

const (
	FulfillmentStatusNew     FulfillmentStatus = "New"
	FulfillmentStatusShipped FulfillmentStatus = "Shipped"
)

// Valid checks if the status is one of the allowed values
func (s FulfillmentStatus) Valid() error {
	switch s {
	case FulfillmentStatusNew, FulfillmentStatusShipped:
		return nil
	default:
		return fmt.Errorf("invalid fulfillment status: %s", s)
	}
}

// type OrderItem struct {
// 	gorm.Model
// 	OrderID           uint              `gorm:"not null" json:"order_id"`
// 	ProductID         string              `gorm:"not null" json:"product_id"`
// 	MerchantID        uint              `gorm:"not null" json:"merchant_id"`
// 	Quantity          int               `gorm:"not null" json:"quantity"`
// 	Price             float64           `gorm:"type:decimal(10,2);not null" json:"price"`
// 	FulfillmentStatus FulfillmentStatus `gorm:"type:varchar(20);not null;default:'New'" json:"fulfillment_status"`
// 	Order             Order             `gorm:"foreignKey:OrderID"`
// 	Product           Product           `gorm:"foreignKey:ProductID"`
// 	Merchant          Merchant          `gorm:"foreignKey:MerchantID"`
// }

type OrderItem struct {
	gorm.Model
	OrderID   uint   `gorm:"not null;index" json:"order_id"`
	ProductID string `gorm:"not null;index" json:"product_id"`
	//ProductID         uint              `gorm:"not null;index" json:"product_id"`
	MerchantID        string            `gorm:"not null;index" json:"merchant_id"`
	Quantity          int               `gorm:"not null" json:"quantity"`
	Price             float64           `gorm:"type:decimal(10,2);not null" json:"price"`
	FulfillmentStatus FulfillmentStatus `gorm:"type:varchar(20);not null;default:'New'" json:"fulfillment_status"`
	Order             Order             `gorm:"foreignKey:OrderID"`
	Product           Product           `gorm:"foreignKey:ProductID;references:ID"`
	Merchant          Merchant         `gorm:"foreignKey:MerchantID;references:MerchantID"`
}

// BeforeCreate validates the FulfillmentStatus field
func (oi *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if err := oi.FulfillmentStatus.Valid(); err != nil {
		return err
	}
	return nil
}

// BeforeUpdate validates the FulfillmentStatus field
func (oi *OrderItem) BeforeUpdate(tx *gorm.DB) error {
	if err := oi.FulfillmentStatus.Valid(); err != nil {
		return err
	}
	return nil
}

```

---
### db/models/order_merchant_split.go
- Size: 0.58 KB
- Lines: 19
- Last Modified: 2025-09-30 12:28:15

```go
package models

import (
    "time"
    "github.com/shopspring/decimal"
    "gorm.io/gorm"
)

type OrderMerchantSplit struct {
    gorm.Model
    OrderID    uint    `gorm:"index"`
    MerchantID string  `gorm:"type:uuid;index"`  // Match Merchant.ID type
    AmountDue  decimal.Decimal
    Fee        decimal.Decimal  // Platform cut
    Status     string  `gorm:"default:'pending'"`  // pending, payout_requested, paid, reversed
    HoldUntil  time.Time
	Merchant   Merchant `gorm:"foreignKey:MerchantID;references:MerchantID"`
	Order             Order             `gorm:"foreignKey:OrderID"`
}
```

---
### db/models/payment.go
- Size: 1.87 KB
- Lines: 67
- Last Modified: 2025-09-30 12:28:15

```go
package models

import (
	"fmt"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// PaymentStatus defines possible payment status values
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "Pending"
	PaymentStatusCompleted PaymentStatus = "Completed"
	PaymentStatusFailed    PaymentStatus = "Failed"
	PaymentStatusRefunded    PaymentStatus = "Refunded"
)

// Valid checks if the status is one of the allowed values
func (s PaymentStatus) Valid() error {
	switch s {
	case PaymentStatusPending, PaymentStatusCompleted, PaymentStatusFailed:
		return nil
	default:
		return fmt.Errorf("invalid payment status: %s", s)
	}
}

// type Payment struct {
// 	gorm.Model
// 	OrderID uint          `gorm:"not null" json:"order_id"`
// 	Amount  float64       `gorm:"not null" json:"amount"`
// 	Status  PaymentStatus `gorm:"type:varchar(20);not null;default:'Pending'" json:"status"`
// 	Order   Order         `gorm:"foreignKey:OrderID"`
// }



type Payment struct {
     gorm.Model
     OrderID       uint              `gorm:"not null"`
     Amount        decimal.Decimal   `gorm:"type:decimal(10,2)" json:"amount"`
    Currency      string            `gorm:"type:varchar(3);default:'NGN'" json:"currency"`
     Status        PaymentStatus     `gorm:"type:varchar(20);not null;default:'Pending'" json:"status"`
    TransactionID string            `gorm:"type:varchar(100);unique" json:"transaction_id"`
     AuthorizationURL *string        `gorm:"type:varchar(500)" json:"authorization_url"`
     Order         Order             `gorm:"foreignKey:OrderID"`
 }



// BeforeCreate validates the Status field
func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if err := p.Status.Valid(); err != nil {
		return err
	}
	return nil
}

// BeforeUpdate validates the Status field
func (p *Payment) BeforeUpdate(tx *gorm.DB) error {
	if err := p.Status.Valid(); err != nil {
		return err
	}
	return nil
}

```

---
### db/models/payout.go
- Size: 1.21 KB
- Lines: 49
- Last Modified: 2025-09-30 12:28:15

```go
package models

import (
	"fmt"
	"gorm.io/gorm"
)

// PayoutStatus defines possible payout status values
type PayoutStatus string

const (
	PayoutStatusPending   PayoutStatus = "Pending"
	PayoutStatusCompleted PayoutStatus = "Completed"
)

// Valid checks if the status is one of the allowed values
func (s PayoutStatus) Valid() error {
	switch s {
	case PayoutStatusPending, PayoutStatusCompleted:
		return nil
	default:
		return fmt.Errorf("invalid payout status: %s", s)
	}
}

type Payout struct {
	gorm.Model
	MerchantID      uint         `gorm:"not null" json:"merchant_id"`
	Amount          float64      `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status          PayoutStatus `gorm:"type:varchar(20);not null;default:'Pending'" json:"status"`
	PayoutAccountID string       `gorm:"size:255;not null" json:"payout_account_id"`
	Merchant        Merchant     `gorm:"foreignKey:MerchantID"`
}

// BeforeCreate validates the Status field
func (p *Payout) BeforeCreate(tx *gorm.DB) error {
	if err := p.Status.Valid(); err != nil {
		return err
	}
	return nil
}

// BeforeUpdate validates the Status field
func (p *Payout) BeforeUpdate(tx *gorm.DB) error {
	if err := p.Status.Valid(); err != nil {
		return err
	}
	return nil
}

```

---
### db/models/product.go
- Size: 3.77 KB
- Lines: 124
- Last Modified: 2025-09-30 12:28:15

```go
package models

/*

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

*/

```

---
### db/models/product2.go
- Size: 9.61 KB
- Lines: 243
- Last Modified: 2025-09-30 12:28:15

```go
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
	PublicID  string     `gorm:"index" json:"public_id"`
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
```

---
### db/models/user.go
- Size: 0.54 KB
- Lines: 15
- Last Modified: 2025-09-30 12:28:15

```go
package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"unique;not null"`
	Name     string `gorm:"type:varchar(100);not null"`
	Password string // Empty for OAuth users
	//Role     string `gorm:"not null"` // "customer" (default) or "merchant" (upgraded by admin)
	GoogleID string // Google ID for OAuth
	Country  string `gorm:"type:varchar(100)"` // Optional country field
	//Carts    []Cart  `gorm:"foreignKey:UserID" json:"carts,omitempty"`
	//Orders   []Order `gorm:"foreignKey:UserID" json:"orders,omitempty"`
}

```

---
### db/repositories/cart_item_repository.go
- Size: 9.73 KB
- Lines: 268
- Last Modified: 2025-09-30 12:28:15

```go
package repositories

import (
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	"fmt"

	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause" // Added for Locking
)

var (
	ErrCartItemNotFound  = errors.New("cart item not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrReservationFailed = errors.New("failed to reserve stock")
)

type CartItemRepository struct {
	db *gorm.DB
}

func NewCartItemRepository() *CartItemRepository {
	return &CartItemRepository{db: db.DB}
}

/*
func (r *CartItemRepository) Create(ctx context.Context, cartItem *models.CartItem) error {
	return r.db.WithContext(ctx).Create(cartItem).Error
}

func (r *CartItemRepository) FindByID(ctx context.Context, id uint) (*models.CartItem, error) {
	var cartItem models.CartItem
	err := r.db.WithContext(ctx).
		Preload("Cart.User").
		Preload("Product.Merchant").
		Preload("Merchant").
		First(&cartItem, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrCartItemNotFound
	}
	return &cartItem, err
}

func (r *CartItemRepository) FindByCartID(ctx context.Context, cartID uint) ([]models.CartItem, error) {
	var cartItems []models.CartItem
	err := r.db.WithContext(ctx).
		Preload("Product.Merchant").
		Preload("Merchant").
		Where("cart_id = ?", cartID).Find(&cartItems).Error
	return cartItems, err
}

func (r *CartItemRepository) FindByProductIDAndCartID(ctx context.Context, productID string, cartID uint) (*models.CartItem, error) {
	var cartItem models.CartItem
	err := r.db.WithContext(ctx).
		Preload("Product.Merchant").
		Preload("Merchant").
		Where("product_id = ? AND cart_id = ?", productID, cartID).First(&cartItem).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrCartItemNotFound
	}
	return &cartItem, err
}

func (r *CartItemRepository) UpdateQuantityWithReservation(ctx context.Context, itemID uint, newQuantity int, inventoryID uint) error { // Changed to uint
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var item models.CartItem
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&item, itemID).Error; err != nil {
			return fmt.Errorf("failed to lock item: %w", err)
		}

		delta := newQuantity - item.Quantity
		if delta > 0 {
			// Reserve (check first via repo if needed)
			if err := tx.Model(&models.Inventory{}).Where("id = ?", inventoryID).
				Update("reserved_quantity", gorm.Expr("reserved_quantity + ?", delta)).Error; err != nil { // Fixed: Assigned err
				return fmt.Errorf("stock reservation failed: %w", ErrReservationFailed)
			}
		} else if delta < 0 {
			err := tx.Model(&models.Inventory{}).Where("id = ?", inventoryID).
				Update("reserved_quantity", gorm.Expr("reserved_quantity - ?", -delta)).Error // Fixed: Assigned err (unused but for consistency)
			if err != nil {
				return fmt.Errorf("stock unreservation failed: %w", err)
			}
		}

		return tx.Model(&models.CartItem{}).Where("id = ?", itemID).Update("quantity", newQuantity).Error
	})
}

func (r *CartItemRepository) Update(ctx context.Context, cartItem *models.CartItem) error {
	return r.db.WithContext(ctx).Save(cartItem).Error
}

func (r *CartItemRepository) DeleteWithUnreserve(ctx context.Context, id uint, inventoryID uint) error { // uint
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var item models.CartItem
		if err := tx.First(&item, id).Error; err != nil {
			return ErrCartItemNotFound
		}
		// Release
		err := tx.Model(&models.Inventory{}).Where("id = ?", inventoryID).
			Update("reserved_quantity", gorm.Expr("reserved_quantity - ?", item.Quantity)).Error // Fixed: Assigned err
		if err != nil {
			return fmt.Errorf("stock unreservation failed: %w", err)
		}
		return tx.Delete(&models.CartItem{}, id).Error
	})
}

func (r *CartItemRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.CartItem{}, id).Error
}

// Stub for inventory_repo.go if missing:
// func (r *InventoryRepository) UpdateInventoryQuantity(ctx context.Context, inventoryID uint, delta int) error { // Add this method
// 	return r.db.WithContext(ctx).Model(&models.Inventory{}).Where("id = ?", inventoryID).
// 		Update("quantity", gorm.Expr("quantity + ?", delta)).Error
// }


*/

func (r *CartItemRepository) Create(ctx context.Context, cartItem *models.CartItem) error {
	return r.db.WithContext(ctx).Create(cartItem).Error
}

func (r *CartItemRepository) FindByID(ctx context.Context, id uint) (*models.CartItem, error) {
	var cartItem models.CartItem
	err := r.db.WithContext(ctx).
		Preload("Cart.User").
		Preload("Product.Merchant").
		Preload("Merchant").
		First(&cartItem, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrCartItemNotFound
	}
	return &cartItem, err
}

func (r *CartItemRepository) FindByCartID(ctx context.Context, cartID uint) ([]models.CartItem, error) {
	var cartItems []models.CartItem
	err := r.db.WithContext(ctx).
		Preload("Product.Merchant").
		Preload("Merchant").
		Where("cart_id = ?", cartID).Find(&cartItems).Error
	return cartItems, err
}

// func (r *CartItemRepository) FindByProductIDAndCartID(ctx context.Context, productID string, cartID uint) (*models.CartItem, error) {
// 	var cartItem models.CartItem
// 	err := r.db.WithContext(ctx).
// 		Preload("Product.Merchant").
// 		Preload("Merchant").
// 		Where("product_id = ? AND cart_id = ?", productID, cartID).
// 		First(&cartItem).Error
// 	if errors.Is(err, gorm.ErrRecordNotFound) {
// 		return nil, ErrCartItemNotFound
// 	}
// 	return &cartItem, err
// }

func (r *CartItemRepository) FindByProductIDAndCartID(ctx context.Context, productID string, variantID *string, cartID uint) (*models.CartItem, error) {
	var item models.CartItem
	query := r.db.WithContext(ctx).Where("product_id = ? AND cart_id = ?", productID, cartID)
	if variantID != nil {
		query = query.Where("variant_id = ?", *variantID)
	}
	err := query.First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &item, err
}

// UpdateQuantityWithReservation updates the quantity in cart + reserved stock in VendorInventory
// func (r *CartItemRepository) UpdateQuantityWithReservation(ctx context.Context, itemID uint, newQuantity int, vendorInvID uint) error {
// 	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
// 		var item models.CartItem
// 		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&item, itemID).Error; err != nil {
// 			return fmt.Errorf("failed to lock item: %w", err)
// 		}

// 		delta := newQuantity - item.Quantity
// 		if delta > 0 {
// 			// Reserve extra stock
// 			if err := tx.Model(&models.VendorInventory{}).Where("id = ?", vendorInvID).
// 				Update("reserved_quantity", gorm.Expr("reserved_quantity + ?", delta)).Error; err != nil {
// 				return fmt.Errorf("stock reservation failed: %w", err)
// 			}
// 		} else if delta < 0 {
// 			// Unreserve stock
// 			if err := tx.Model(&models.VendorInventory{}).Where("id = ?", vendorInvID).
// 				Update("reserved_quantity", gorm.Expr("reserved_quantity - ?", -delta)).Error; err != nil {
// 				return fmt.Errorf("stock unreservation failed: %w", err)
// 			}
// 		}

// 		return tx.Model(&models.CartItem{}).Where("id = ?", itemID).Update("quantity", newQuantity).Error
// 	})
// }

func (r *CartItemRepository) UpdateQuantityWithReservation(ctx context.Context, itemID uint, newQuantity int, vendorInvID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var item models.CartItem
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&item, itemID).Error; err != nil {
			return fmt.Errorf("failed to lock item: %w", err)
		}

		delta := newQuantity - item.Quantity
		if delta > 0 {
			// Reserve extra stock
			if err := tx.Model(&models.Inventory{}).Where("id = ?", vendorInvID).
				Update("reserved_quantity", gorm.Expr("reserved_quantity + ?", delta)).Error; err != nil {
				return fmt.Errorf("stock reservation failed: %w", err)
			}
		} else if delta < 0 {
			// Unreserve stock
			if err := tx.Model(&models.Inventory{}).Where("id = ?", vendorInvID).
				Update("reserved_quantity", gorm.Expr("reserved_quantity - ?", -delta)).Error; err != nil {
				return fmt.Errorf("stock unreservation failed: %w", err)
			}
		}

		return tx.Model(&models.CartItem{}).Where("id = ?", itemID).Update("quantity", newQuantity).Error
	})
}

func (r *CartItemRepository) Update(ctx context.Context, cartItem *models.CartItem) error {
	return r.db.WithContext(ctx).Save(cartItem).Error
}

// func (r *CartItemRepository) DeleteWithUnreserve(ctx context.Context, id uint, vendorInvID uint) error {
// 	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
// 		var item models.CartItem
// 		if err := tx.First(&item, id).Error; err != nil {
// 			return ErrCartItemNotFound
// 		}
// 		// Release reserved stock
// 		if err := tx.Model(&models.VendorInventory{}).Where("id = ?", vendorInvID).
// 			Update("reserved_quantity", gorm.Expr("reserved_quantity - ?", item.Quantity)).Error; err != nil {
// 			return fmt.Errorf("stock unreservation failed: %w", err)
// 		}
// 		return tx.Delete(&models.CartItem{}, id).Error
// 	})
// }

func (r *CartItemRepository) DeleteWithUnreserve(ctx context.Context, id uint, vendorInvID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var item models.CartItem
		if err := tx.First(&item, id).Error; err != nil {
			return ErrCartItemNotFound
		}
		// Release reserved stock
		if err := tx.Model(&models.Inventory{}).Where("id = ?", vendorInvID).
			Update("reserved_quantity", gorm.Expr("reserved_quantity - ?", item.Quantity)).Error; err != nil {
			return fmt.Errorf("stock unreservation failed: %w", err)
		}
		return tx.Delete(&models.CartItem{}, id).Error
	})
}

func (r *CartItemRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.CartItem{}, id).Error
}

```

---
### db/repositories/cart_repository.go
- Size: 2.95 KB
- Lines: 89
- Last Modified: 2025-09-30 12:28:15

```go
package repositories

import (
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

var ErrCartNotFound = errors.New("cart not found")

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository() *CartRepository {
	return &CartRepository{db: db.DB}
}

// Create adds a new cart
func (r *CartRepository) Create(ctx context.Context, cart *models.Cart) error {
	return r.db.WithContext(ctx).Create(cart).Error
}

// FindByID retrieves a cart by ID with associated User and CartItems
// func (r *CartRepository) FindByID(id uint) (*models.Cart, error) {
// 	var cart models.Cart
// 	err := r.db.Preload("User").Preload("CartItems.Product.Merchant").First(&cart, id).Error
// 	return &cart, err
// }

func (r *CartRepository) FindByID(ctx context.Context, id uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("CartItems.Product.Media"). // Efficient: preload media for UI
		Preload("CartItems.Product.Variants.Inventory").
		First(&cart, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrCartNotFound
	}
	return &cart, err
}

// FindActiveCart retrieves the user's most recent active cart
// func (r *CartRepository) FindActiveCart(userID uint) (*models.Cart, error) {
// 	var cart models.Cart
// 	err := r.db.Preload("CartItems.Product.Merchant").Where("user_id = ? AND status = ?", userID, models.CartStatusActive).Order("created_at DESC").First(&cart).Error
// 	return &cart, err
// }

func (r *CartRepository) FindActiveCart(ctx context.Context, userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.WithContext(ctx).
		Preload("CartItems.Product.Merchant"). // As before
		Where("user_id = ? AND status = ?", userID, models.CartStatusActive).
		Order("created_at DESC").First(&cart).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound //ErrCartNotFound
	}
	return &cart, err
}

// FindByUserIDAndStatus retrieves carts for a user by status
func (r *CartRepository) FindByUserIDAndStatus(ctx context.Context, userID uint, status models.CartStatus) ([]models.Cart, error) {
	var carts []models.Cart
	err := r.db.WithContext(ctx).
		Preload("CartItems.Product.Merchant").Where("user_id = ? AND status = ?", userID, status).Find(&carts).Error
	return carts, err
}

// FindByUserID retrieves all carts for a user
func (r *CartRepository) FindByUserID(ctx context.Context, userID uint) ([]models.Cart, error) {
	var carts []models.Cart
	err := r.db.WithContext(ctx).Preload("CartItems.Product.Merchant").Where("user_id = ?", userID).Find(&carts).Error
	return carts, err
}

// Update modifies an existing cart
func (r *CartRepository) Update(ctx context.Context, cart *models.Cart) error {
	return r.db.WithContext(ctx).Save(cart).Error
}

// Delete removes a cart by ID
func (r *CartRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Cart{}, id).Error
}

```

---
### db/repositories/category_repositry.go
- Size: 1.14 KB
- Lines: 45
- Last Modified: 2025-09-30 12:28:15

```go
package repositories

import (
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{db: db.DB}
}

// Create adds a new category
func (r *CategoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

// FindByID retrieves a category by ID with parent category
func (r *CategoryRepository) FindByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.Preload("Parent").First(&category, id).Error
	return &category, err
}

// FindAll retrieves all categories
func (r *CategoryRepository) FindAll() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Preload("Parent").Find(&categories).Error
	return categories, err
}

// Update modifies an existing category
func (r *CategoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

// Delete removes a category by ID
func (r *CategoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}

```

---
### db/repositories/disputes_repository.go
- Size: 3.77 KB
- Lines: 111
- Last Modified: 2025-10-02 14:57:10

```go
package repositories

import (
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	"context"
	"time"

	"gorm.io/gorm"
)

type DisputeRepository struct {
	db *gorm.DB
}

type ReturnRequestRepository struct{
	db *gorm.DB
}


func NewDisputeRepository() *DisputeRepository {
	return &DisputeRepository{db: db.DB}
}


func NewReturnRequestRepository() *ReturnRequestRepository {
	return &ReturnRequestRepository{db: db.DB}
}



// Create adds a new order item
func (r *DisputeRepository) Create(ctx context.Context,dispute *models.Dispute) error {
	return r.db.WithContext(ctx).Create(dispute).Error
}


// FindMediaByID fetches media
func (r *DisputeRepository) FindDisputeByID(ctx context.Context, id string) (*models.Dispute, error) {
	var dispute models.Dispute
	err := r.db.WithContext(ctx).Scopes(r.activeScope()).First(&dispute, "id = ?", id).Error
	return &dispute, err
}

// UpdateMedia updates fields
func (r *DisputeRepository) Update(ctx context.Context, dispute *models.Dispute) error {
	return r.db.WithContext(ctx).Save(dispute).Error
}

// DeleteMedia soft-deletes
func (r *DisputeRepository) DeleteDispute(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&models.Dispute{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}


// activeScope (if soft delete)
func (r *DisputeRepository) activeScope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB { return db.Where("deleted_at IS NULL") }
}


func (r *ReturnRequestRepository) Create(ctx context.Context,returnrequests *models.ReturnRequest) error {
	return r.db.WithContext(ctx).Create(returnrequests).Error
}


// FindReturnRequestByID fetches a return request by ID
func (r *ReturnRequestRepository) FindReturnRequestByID(ctx context.Context, id string) (*models.ReturnRequest, error) {
    var returnRequest models.ReturnRequest
    err := r.db.WithContext(ctx).Scopes(r.activeScope()).First(&returnRequest, "id = ?", id).Error
    return &returnRequest, err
}

// Update updates a return request
// func (r *ReturnRequestRepository) Update(ctx context.Context, returnRequest *models.ReturnRequest) error {
//     return r.db.WithContext(ctx).Save(returnRequest).Error
// }

// DeleteReturnRequest soft-deletes a return request
func (r *ReturnRequestRepository) DeleteReturnRequest(ctx context.Context, id string) error {
    return r.db.WithContext(ctx).Model(&models.ReturnRequest{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

// activeScope filters out soft-deleted records
func (r *ReturnRequestRepository) activeScope() func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB { return db.Where("deleted_at IS NULL") }
}

func (r *ReturnRequestRepository) FindByID(ctx context.Context, id string) (*models.ReturnRequest, error) {
    var returnRequest models.ReturnRequest
    err := r.db.WithContext(ctx).Scopes(r.activeScope()).First(&returnRequest, "id = ?", id).Error
    return &returnRequest, err
}

func (r *ReturnRequestRepository) Update(ctx context.Context, returnRequest *models.ReturnRequest) error {
    return r.db.WithContext(ctx).Save(returnRequest).Error
}

func (r *ReturnRequestRepository) Delete(ctx context.Context, id string) error {
    return r.db.WithContext(ctx).Model(&models.ReturnRequest{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

// func (r *ReturnRequestRepository) activeScope() func(db *gorm.DB) *gorm.DB {
//     return func(db *gorm.DB) *gorm.DB { return db.Where("deleted_at IS NULL") }
// }

func (r *ReturnRequestRepository) FindByCustomerID(ctx context.Context, customerID uint) ([]models.ReturnRequest, error) {
    var returnRequests []models.ReturnRequest
    err := r.db.WithContext(ctx).Scopes(r.activeScope()).Where("customer_id = ?", customerID).Find(&returnRequests).Error
    return returnRequests, err
}
```

---
### db/repositories/inventory_repository.go
- Size: 5.13 KB
- Lines: 156
- Last Modified: 2025-09-30 12:28:15

```go
package repositories

import (
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	"errors"

	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
type InventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository() *InventoryRepository {
	return &InventoryRepository{db: db.DB}
}

// Create adds a new inventory record
func (r *InventoryRepository) Create(inventory *models.Inventory) error {
	return r.db.Create(inventory).Error
}

// FindByProductID retrieves inventory by product ID
func (r *InventoryRepository) FindByProductID(productID string) (*models.Inventory, error) {
	var inventory models.Inventory
	err := r.db.Where("product_id = ?", productID).First(&inventory).Error
	return &inventory, err
}

// UpdateStock updates the stock quantity for a product
func (r *InventoryRepository) UpdateStock(productID string, quantityChange int) error {
	return r.db.Model(&models.Inventory{}).Where("product_id = ?", productID).
		Update("stock_quantity", gorm.Expr("stock_quantity + ?", quantityChange)).Error
}

// Update modifies an existing inventory record
func (r *InventoryRepository) Update(inventory *models.Inventory) error {
	return r.db.Save(inventory).Error
}

// Delete removes an inventory record by ID
func (r *InventoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Inventory{}, id).Error
}

 func (r *InventoryRepository) UpdateInventoryQuantity(ctx context.Context, inventoryID uint, delta int) error { // Add this method
 	return r.db.WithContext(ctx).Model(&models.Inventory{}).Where("id = ?", inventoryID).
 		Update("quantity", gorm.Expr("quantity + ?", delta)).Error
 }
*/
type InventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository() *InventoryRepository {
	return &InventoryRepository{db: db.DB}
}

// Create adds a new vendor inventory record
func (r *InventoryRepository) Create(ctx context.Context, inv *models.Inventory) error {
	return r.db.WithContext(ctx).Create(inv).Error
}

// FindByVariantID retrieves vendor inventory by variant ID
func (r *InventoryRepository) FindByVariantID(ctx context.Context, variantID, merchantID string) (*models.Inventory, error) {
	var inv models.Inventory
	return &inv, r.db.WithContext(ctx).
		Where("variant_id = ? AND merchant_id = ?", variantID, merchantID).First(&inv).Error
}

// FindByProductID (for simple products without variants)
func (r *InventoryRepository) FindByProductID(ctx context.Context, productID string, merchantID string) (*models.Inventory, error) {
	var inv models.Inventory
	err := r.db.WithContext(ctx).
		Where("product_id = ? AND merchant_id = ?", productID, merchantID).
		First(&inv).Error
	return &inv, err
}

// UpdateStock adjusts quantity (can be negative for reservations)
func (r *InventoryRepository) UpdateStock(ctx context.Context, invID uint, delta int) error {
	return r.db.WithContext(ctx).
		Model(&models.Inventory{}).
		Where("id = ?", invID).
		Update("quantity", gorm.Expr("quantity + ?", delta)).
		Error
}

// ReserveStock increments reserved quantity
func (r *InventoryRepository) ReserveStock(ctx context.Context, invID uint, qty int) error {
	return r.db.WithContext(ctx).
		Model(&models.Inventory{}).
		Where("id = ?", invID).
		Update("reserved_quantity", gorm.Expr("reserved_quantity + ?", qty)).
		Error
}

// ReleaseStock decrements reserved quantity
func (r *InventoryRepository) ReleaseStock(ctx context.Context, invID uint, qty int) error {
	return r.db.WithContext(ctx).
		Model(&models.Inventory{}).
		Where("id = ?", invID).
		Update("reserved_quantity", gorm.Expr("reserved_quantity - ?", qty)).
		Error
}

// Delete removes a vendor inventory record by ID
func (r *InventoryRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Inventory{}, id).Error
}

// UpdateInventoryQuantity updates Quantity (can be negative)

func (r *InventoryRepository) UpdateInventoryQuantity(ctx context.Context, inventoryID string, delta int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var inv models.Inventory
		if err := tx.First(&inv, "id = ?", inventoryID).Error; err != nil {
			return err
		}
		newQ := inv.Quantity + delta
		if newQ < 0 && !inv.BackorderAllowed {
			return errors.New("insufficient stock and backorders not allowed")
		}
		inv.Quantity = newQ
		return tx.Save(&inv).Error
	})
}




// Add method for lookup by product and merchant (no VariantID)
func (r *InventoryRepository) FindByProductAndMerchant(ctx context.Context, productID, merchantID string) (*models.Inventory, error) {
	var inv models.Inventory
	err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}). // Lock for update
		Where("product_id = ? AND merchant_id = ?", productID, merchantID).
		First(&inv).Error
	return &inv, err
}

// UpdateInventory updates quantity/reserved (delta positive for unreserve)
func (r *InventoryRepository) UpdateInventory(ctx context.Context, id string, delta int) error {
	return r.db.WithContext(ctx).
		Model(&models.Inventory{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"quantity":        gorm.Expr("quantity + ?", delta),
			"reserved_quantity": gorm.Expr("GREATEST(reserved_quantity - ?, 0)", delta),
		}).Error
}
```

---
### db/repositories/merchant_repository.go
- Size: 3.07 KB
- Lines: 95
- Last Modified: 2025-09-30 12:28:15

```go
package repositories

import (
	"context"
	"log"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	//"gorm.io/gorm"
)

// MerchantApplicationRepository handles CRUD for merchant applications
// Note: Admin service (in Express) will be responsible for updating status/approval.
type MerchantApplicationRepository struct{}

func NewMerchantApplicationRepository() *MerchantApplicationRepository {
	return &MerchantApplicationRepository{}
}

func (r *MerchantApplicationRepository) Create(ctx context.Context, m *models.MerchantApplication) error {
	err := db.DB.WithContext(ctx).Create(m).Error
	if err != nil {
		log.Printf("Failed to create merchant application: %v", err)
		return err
	}
	return nil
}

func (r *MerchantApplicationRepository) GetByID(ctx context.Context, id string) (*models.MerchantApplication, error) {
	var m models.MerchantApplication
	if err := db.DB.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		log.Printf("Failed to get merchant application by ID %s: %v", id, err)
		return nil, err
	}
	return &m, nil
}

func (r *MerchantApplicationRepository) GetByUserEmail(ctx context.Context, email string) (*models.MerchantApplication, error) {
	var m models.MerchantApplication
	if err := db.DB.WithContext(ctx).Where("personal_email = ? OR work_email = ?", email, email).First(&m).Error; err != nil {
		log.Printf("Failed to get merchant application by email %s: %v", email, err)
		return nil, err
	}
	return &m, nil
}

// MerchantRepository handles active merchants
type MerchantRepository struct{}

func NewMerchantRepository() *MerchantRepository {
	return &MerchantRepository{}
}

func (r *MerchantRepository) GetByID(ctx context.Context, id string) (*models.Merchant, error) {
	var m models.Merchant
	if err := db.DB.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		log.Printf("Failed to get merchant by ID %s: %v", id, err)
		return nil, err
	}
	return &m, nil
}

func (r *MerchantRepository) GetByMerchantID(ctx context.Context, uid string) (*models.Merchant, error) {
	var m models.Merchant
	if err := db.DB.WithContext(ctx).Where("merchant_id = ?", uid).First(&m).Error; err != nil {
		log.Printf("Failed to get merchant by user ID %s: %v", uid, err)
		return nil, err
	}
	return &m, nil
}

func (r *MerchantRepository) GetByWorkEmail(ctx context.Context, email string) (*models.Merchant, error) {
	var m models.Merchant
	if err := db.DB.WithContext(ctx).Where("personal_email = ? OR work_email = ?", email, email).First(&m).Error; err != nil {
		log.Printf("Failed to get merchant  by email %s: %v", email, err)
		return nil, err
	}
	return &m, nil
}



func (r *MerchantRepository) UpdateBankDetails(ctx context.Context, merchantID string, details dto.BankDetailsRequest) error {
	// Use WithContext so DB operations respect request lifecycle
	if err := db.DB.WithContext(ctx).
		Model(&models.MerchantBankDetails{}).
		Where("merchant_id = ?", merchantID).
		Save(details).Error; err != nil {
		
		log.Printf("Failed to update bank details for merchant %s: %v", merchantID, err)
		return err
	}
	return nil
}

```

---
### db/repositories/order_item_repository.go
- Size: 1.63 KB
- Lines: 55
- Last Modified: 2025-09-30 12:28:15

```go
package repositories

import (
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	"context"

	"gorm.io/gorm"
)

type OrderItemRepository struct {
	db *gorm.DB
}

func NewOrderItemRepository() *OrderItemRepository {
	return &OrderItemRepository{db: db.DB}
}

// Create adds a new order item
func (r *OrderItemRepository) Create(orderItem *models.OrderItem) error {
	return r.db.Create(orderItem).Error
}

// FindByID retrieves an order item by ID with associated Order, Product, and Merchant
func (r *OrderItemRepository) FindByID(id uint) (*models.OrderItem, error) {
	var orderItem models.OrderItem
	err := r.db.Preload("Order.User").Preload("Product.Merchant").Preload("Merchant").First(&orderItem, id).Error
	return &orderItem, err
}

// FindByOrderID retrieves all order items for an order
func (r *OrderItemRepository) FindByOrderID(orderID uint) ([]models.OrderItem, error) {
	var orderItems []models.OrderItem
	err := r.db.Preload("Product.Merchant").Preload("Merchant").Where("order_id = ?", orderID).Find(&orderItems).Error
	return orderItems, err
}

// Update modifies an existing order item
func (r *OrderItemRepository) Update(orderItem *models.OrderItem) error {
	return r.db.Save(orderItem).Error
}

// Delete removes an order item by ID
func (r *OrderItemRepository) Delete(id uint) error {
	return r.db.Delete(&models.OrderItem{}, id).Error
}



// In orderItemRepository
func (r *OrderItemRepository) FindOrderItemsByOrderID(ctx context.Context, orderID uint) ([]models.OrderItem, error) {
	var items []models.OrderItem
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).Find(&items).Error
	return items, err
}
```

---
### db/repositories/order_repository.go
- Size: 3.12 KB
- Lines: 99
- Last Modified: 2025-10-02 13:31:15

```go
package repositories

import (
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	"context"

	//"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{db: db.DB}
}

// Create adds a new order
// func (r *OrderRepository) Create(order *models.Order) error {
// 	return r.db.Create(order).Error
// }

type OrderInterface interface {
	FindByID(ctx context.Context, id uint) (*models.Order, error)
	// Add other methods as needed
}

func (r *OrderRepository) Create(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

// FindByID retrieves an order by ID with associated User and OrderItems
func (r *OrderRepository) FindByID(ctx context.Context, id uint) (*models.Order, error) {
	var order models.Order
	//err := r.db.Preload("User").Preload("OrderItems.Product.Merchant").First(&order, id).Error
	err := r.db.WithContext(ctx).Preload("User").Preload("OrderItems").Preload("OrderItems.Product").Preload("OrderItems.Merchant").First(&order, id).Error
	return &order, err
}

// FindByUserID retrieves all orders for a user
func (r *OrderRepository) FindByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("OrderItems.Product.Merchant").Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}

// FindByMerchantID retrieves all orders containing items from a merchant
func (r *OrderRepository) FindByMerchantID(merchantID uint) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("OrderItems.Product").Joins("JOIN order_items oi ON oi.order_id = orders.id").
		Where("oi.merchant_id = ?", merchantID).Find(&orders).Error
	return orders, err
}

// Update modifies an existing order
func (r *OrderRepository) Update(ctx context.Context,order *models.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

// Delete removes an order by ID
func (r *OrderRepository) Delete(id uint) error {
	return r.db.Delete(&models.Order{}, id).Error
}


// FindByIDWithPreloads fetches with ownership check and preloads (avoids N+1)
func (r *OrderRepository) FindByIDWithPreloads(ctx context.Context, id uint) (*models.Order, error) {
	var order models.Order
	// Preload OrderItems (no deeper Inventory preload to avoid N+1; fetch separately if needed)
	err := r.db.WithContext(ctx).
		Scopes(r.activeScope()). // Soft delete filter
		Preload("OrderItems").
		Preload("Payment").
		First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// UpdateStatus updates order status (with locking for concurrency)
func (r *OrderRepository) UpdateStatus(ctx context.Context, id uint, status models.OrderStatus) error {
	return r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Model(&models.Order{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// activeScope for soft deletes (if Order has DeletedAt)
func (r *OrderRepository) activeScope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Unscoped().Where("deleted_at IS NULL")
	}
}
```

---
### db/repositories/payment_repository.go
- Size: 2.23 KB
- Lines: 71
- Last Modified: 2025-09-30 12:28:15

```go
package repositories

import (
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

var ErrPaymentNotFound     = errors.New("payment not found")

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository() *PaymentRepository {
	return &PaymentRepository{db: db.DB}
}

// Create adds a new payment
func (r *PaymentRepository) Create(ctx context.Context,payment *models.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}


 func (r *PaymentRepository) FindByTransactionID(ctx context.Context, txID string) (*models.Payment, error) {
     var payment models.Payment
     err := r.db.WithContext(ctx).Where("transaction_id = ?", txID).First(&payment).Error
     if err != nil {
         if errors.Is(err, gorm.ErrRecordNotFound) {
             return nil, ErrPaymentNotFound
         }
         return nil, err
     }
     return &payment, nil
 }



// FindByID retrieves a payment by ID with associated Order and User
func (r *PaymentRepository) FindByID(ctx context.Context, id uint) (*models.Payment, error) {
	var payment models.Payment
	err :=  r.db.WithContext(ctx).Preload("Order.User").First(&payment, id).Error
	return &payment, err
}

// FindByOrderID retrieves a payment by order ID
func (r *PaymentRepository) FindByOrderID(ctx context.Context ,orderID uint) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.WithContext(ctx).Preload("Order.User").Where("order_id = ?", orderID).First(&payment).Error
	return &payment, err
}

// FindByUserID retrieves all payments for a user
func (r *PaymentRepository) FindByUserID(ctx context.Context ,userID uint) ([]models.Payment, error) {
	var payments []models.Payment
	err := r.db.WithContext(ctx).Preload("Order.User").Joins("JOIN orders ON orders.id = payments.order_id").Where("orders.user_id = ?", userID).Find(&payments).Error
	return payments, err
}

// Update modifies an existing payment
func (r *PaymentRepository) Update(ctx context.Context,payment *models.Payment) error {
	return r.db.WithContext(ctx).Save(payment).Error
}

// Delete removes a payment by ID
func (r *PaymentRepository) Delete(ctx context.Context,id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Payment{}, id).Error
}

```

---
### db/repositories/payout_repository.go
- Size: 1.18 KB
- Lines: 45
- Last Modified: 2025-09-30 12:28:15

```go
package repositories

import (
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"

	"gorm.io/gorm"
)

type PayoutRepository struct {
	db *gorm.DB
}

func NewPayoutRepository() *PayoutRepository {
	return &PayoutRepository{db: db.DB}
}

// Create adds a new payout record
func (r *PayoutRepository) Create(payout *models.Payout) error {
	return r.db.Create(payout).Error
}

// FindByID retrieves a payout by ID with associated Merchant
func (r *PayoutRepository) FindByID(id uint) (*models.Payout, error) {
	var payout models.Payout
	err := r.db.Preload("Merchant").First(&payout, id).Error
	return &payout, err
}

// FindByMerchantID retrieves all payouts for a merchant
func (r *PayoutRepository) FindByMerchantID(merchantID uint) ([]models.Payout, error) {
	var payouts []models.Payout
	err := r.db.Preload("Merchant").Where("merchant_id = ?", merchantID).Find(&payouts).Error
	return payouts, err
}

// Update modifies an existing payout
func (r *PayoutRepository) Update(payout *models.Payout) error {
	return r.db.Save(payout).Error
}

// Delete removes a payout by ID
func (r *PayoutRepository) Delete(id uint) error {
	return r.db.Delete(&models.Payout{}, id).Error
}

```

---
### db/repositories/product_repo.go
- Size: 15.25 KB
- Lines: 503
- Last Modified: 2025-09-30 12:28:15

```go
package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

var (
	ErrProductNotFound  = errors.New("product not found")
	ErrDuplicateSKU     = errors.New("duplicate SKU")
	ErrInvalidInventory = errors.New("invalid inventory setup")
	ErrMerchantNotFound = errors.New("merchant not found")
)




type ProductFilter struct {
    CategoryName   *string
    CategoryID     *uint
    MinPrice       *decimal.Decimal
    MaxPrice       *decimal.Decimal
    InStock        *bool
    VariantAttrs   map[string]interface{}
    MerchantName   *string
}


type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{db: db.DB}
}

// func (r *ProductRepository) FindBySKU(sku string) (*models.Product, error) {
// 	var product models.Product
// 	err := r.db.Where("sku = ?", sku).First(&product).Error
// 	if errors.Is(err, gorm.ErrRecordNotFound) {
// 		return nil, ErrProductNotFound
// 	} else if err != nil {
// 		return nil, fmt.Errorf("failed to find product by SKU: %w", err)
// 	}
// 	return &product, nil
// }



func (r *ProductRepository) FindBySKU(ctx context.Context, sku string) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).Where("sku = ? AND deleted_at IS NULL", sku).First(&product).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrProductNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to find product by SKU: %w", err)
	}
	return &product, nil
}



// func (r *ProductRepository) FindByID(id string, preloads ...string) (*models.Product, error) {
// 	var product models.Product
// 	query := r.db.Where("id = ?", id)
// 	for _, preload := range preloads {
// 		query = query.Preload(preload)
// 	}
// 	err := query.First(&product).Error
// 	if errors.Is(err, gorm.ErrRecordNotFound) {
// 		return nil, ErrProductNotFound
// 	} else if err != nil {
// 		return nil, fmt.Errorf("failed to find product by ID: %w", err)
// 	}
// 	return &product, nil
// }






func (r *ProductRepository) FindByID(ctx context.Context, id string, preloads ...string) (*models.Product, error) {
	var product models.Product
	query := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id)
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	err := query.First(&product).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrProductNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to find product by ID: %w", err)
	}
	return &product, nil
}





// func (r *ProductRepository) ListByMerchant(merchantID string, limit, offset int, filterActive bool) ([]models.Product, error) {
// 	var products []models.Product
// 	query := r.db.Where("merchant_id = ?", merchantID).Limit(limit).Offset(offset)
// 	if filterActive {
// 		query = query.Where("deleted_at IS NULL")
// 	}
// 	err := query.Find(&products).Error
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to list products: %w", err)
// 	}
// 	return products, nil
// }






func (r *ProductRepository) ListByMerchant(ctx context.Context, merchantID string, limit, offset int, filterActive bool) ([]models.Product, error) {
	var products []models.Product
	query := r.db.WithContext(ctx).Where("merchant_id = ?", merchantID).Limit(limit).Offset(offset)
	if filterActive {
		query = query.Where("deleted_at IS NULL")
	}
	err := query.Find(&products).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	return products, nil
}







// func (r *ProductRepository) GetAllProducts(limit, offset int, categoryID *uint, preloads ...string) ([]models.Product, int64, error) {
// 	var products []models.Product
// 	query := r.db.Model(&models.Product{}).Where("deleted_at IS NULL")
// 	if categoryID != nil {
// 		query = query.Where("category_id = ?", *categoryID)
// 	}
// 	var total int64
// 	if err := query.Count(&total).Error; err != nil {
// 		return nil, 0, fmt.Errorf("failed to count products: %w", err)
// 	}
// 	for _, preload := range preloads {
// 		query = query.Preload(preload)
// 	}
// 	query = query.Limit(limit).Offset(offset).Order("created_at DESC")
// 	err := query.Find(&products).Error
// 	if err != nil {
// 		return nil, 0, fmt.Errorf("failed to fetch products: %w", err)
// 	}
// 	return products, total, nil
// }






func (r *ProductRepository) GetAllProducts(ctx context.Context, limit, offset int, categoryID *uint, preloads ...string) ([]models.Product, int64, error) {
	var products []models.Product
	query := r.db.WithContext(ctx).Model(&models.Product{}).Where("deleted_at IS NULL")
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	query = query.Limit(limit).Offset(offset).Order("created_at DESC")
	err := query.Find(&products).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch products: %w", err)
	}
	return products, total, nil
}






func (r *ProductRepository) ProductsFilter(
    ctx context.Context,
    filter ProductFilter,
    limit, offset int,
    preloads ...string,
) ([]models.Product, int64, error) {
    var products []models.Product

    query := r.db.WithContext(ctx).
        Model(&models.Product{}).
        Joins("LEFT JOIN categories ON categories.id = products.category_id").
        Joins("LEFT JOIN merchants ON merchants.id = products.merchant_id").
        Joins("LEFT JOIN variants ON variants.product_id = products.id").
        Joins("LEFT JOIN inventories ON inventories.product_id = products.id OR inventories.variant_id = variants.id").
        Where("products.deleted_at IS NULL")

    // --- Apply filters ---
    if filter.CategoryID != nil {
        query = query.Where("products.category_id = ?", *filter.CategoryID)
    }
    if filter.CategoryName != nil {
        query = query.Where("categories.name ILIKE ?", "%"+*filter.CategoryName+"%")
    }
    if filter.MinPrice != nil {
        query = query.Where("products.base_price >= ?", *filter.MinPrice)
    }
    if filter.MaxPrice != nil {
        query = query.Where("products.base_price <= ?", *filter.MaxPrice)
    }
    if filter.InStock != nil {
        if *filter.InStock {
            query = query.Where("(inventories.quantity - inventories.reserved_quantity) > 0")
        } else {
            query = query.Where("(inventories.quantity - inventories.reserved_quantity) <= 0")
        }
    }
    if filter.MerchantName != nil {
        query = query.Where("merchant.store_name ILIKE ?", "%"+*filter.MerchantName+"%")
    }
    if len(filter.VariantAttrs) > 0 {
        for key, val := range filter.VariantAttrs {
            // Postgres JSONB query on variant.attributes
            query = query.Where("variants.attributes ->> ? = ?", key, fmt.Sprintf("%v", val))
        }
    }

    // --- Count total ---
    var total int64
    if err := query.Distinct("products.id").Count(&total).Error; err != nil {
        return nil, 0, fmt.Errorf("failed to count products: %w", err)
    }

    // --- Preloads ---
    for _, preload := range preloads {
        query = query.Preload(preload)
    }

    // --- Fetch results ---
    err := query.Distinct("products.id").
        Limit(limit).
        Offset(offset).
        Order("products.created_at DESC").
        Find(&products).Error

    if err != nil {
        return nil, 0, fmt.Errorf("failed to fetch products: %w", err)
    }

    return products, total, nil
}





// func (r *ProductRepository) CreateProductWithVariantsAndInventory(ctx context.Context, product *models.Product, variants []models.Variant, variantInputs []dto.VariantInput, media []models.Media, simpleInitialStock *int, isSimple bool) error {
// 	if isSimple && len(variants) > 0 {
// 		return ErrInvalidInventory // Cannot have variants for simple products
// 	}
// 	if !isSimple && (len(variants) == 0 || len(variants) != len(variantInputs)) {
// 		return ErrInvalidInventory // Must provide matching variants and inputs
// 	}

// 	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
// 		// Create product
// 		if err := tx.Create(product).Error; err != nil {
// 			if errors.Is(err, gorm.ErrDuplicatedKey) {
// 				return ErrDuplicateSKU
// 			}
// 			return fmt.Errorf("failed to create product: %w", err)
// 		}

// 		if !isSimple {
// 			// Variant-based product
// 			for i := range variants {
// 				variants[i].ProductID = product.ID
// 				if err := tx.Create(&variants[i]).Error; err != nil {
// 					return fmt.Errorf("failed to create variant: %w", err)
// 				}
// 				variantIDPtr := variants[i].ID
// 				inventory := models.Inventory{
// 					VariantID:         variantIDPtr,
// 					ProductID:         nil,
// 					MerchantID:        product.MerchantID,
// 					Quantity:          variantInputs[i].InitialStock,
// 					ReservedQuantity:  0,
// 					LowStockThreshold: 10,
// 					BackorderAllowed:  false,
// 				}
// 				if err := tx.Create(&inventory).Error; err != nil {
// 					return fmt.Errorf("failed to create variant inventory: %w", err)
// 				}
// 				variants[i].Inventory = inventory
// 			}
// 		}
// 		// Note: Skip VendorInventory creation for simple products

// 		// Create media
// 		for i := range media {
// 			media[i].ProductID = product.ID
// 			if err := tx.Create(&media[i]).Error; err != nil {
// 				return fmt.Errorf("failed to create media: %w", err)
// 			}
// 		}

// 		// Reload with preloads
// 		preloadQuery := tx.Where("id = ?", product.ID)
// 		if !isSimple {
// 			preloadQuery = preloadQuery.Preload("Variants.Inventory")
// 		}
// 		preloadQuery = preloadQuery.Preload("Media")
// 		if err := preloadQuery.First(product).Error; err != nil {
// 			return fmt.Errorf("failed to preload associations: %w", err)
// 		}

// 		return nil
// 	})
// }




func (r *ProductRepository) CreateProductWithVariantsAndInventory(ctx context.Context, product *models.Product, variants []models.Variant, variantInputs []dto.VariantInput, media []models.Media, simpleInitialStock *int, isSimple bool) error {
	// Validate Merchant exists
	var merchant models.Merchant
	if err := r.db.WithContext(ctx).Where("merchant_id = ?", product.MerchantID).First(&merchant).Error; err != nil {
		return ErrMerchantNotFound
	}

	// Validate SKU uniqueness
	if p, _ := r.FindBySKU(ctx, product.SKU); p != nil {
		return ErrDuplicateSKU
	}
	// for _, v := range variants {
	// 	if v2, _ := r.db.WithContext(ctx).Where("sku = ? AND deleted_at IS NULL", v.SKU).First(&models.Variant{}).Error; v2 == nil {
	// 		return ErrDuplicateSKU
	// 	}
	// }

	for _, v := range variants {
var temp models.Variant
if err := r.db.WithContext(ctx).Where("sku = ? AND deleted_at IS NULL", v.SKU).First(&temp).Error; err == nil {
return ErrDuplicateSKU
} else if !errors.Is(err, gorm.ErrRecordNotFound) {
return err  // Propagate unexpected errors
}
}

	// Validate inputs
	if isSimple && len(variants) > 0 {
		return ErrInvalidInventory // Cannot have variants for simple products
	}
	if !isSimple && (len(variants) == 0 || len(variants) != len(variantInputs)) {
		return ErrInvalidInventory // Must provide matching variants and inputs
	}
	for _, vi := range variantInputs {
		if vi.InitialStock < 0 {
			return errors.New("initial stock cannot be negative")
		}
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create product
		if err := tx.Create(product).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return ErrDuplicateSKU
			}
			return fmt.Errorf("failed to create product: %w", err)
		}

		if isSimple {
			// Create inventory for simple product
			if simpleInitialStock == nil {
				return errors.New("simpleInitialStock required for simple products")
			}
			inventory := models.Inventory{
				ProductID:         &product.ID,
				MerchantID:        product.MerchantID,
				Quantity:          *simpleInitialStock,
				ReservedQuantity:  0,
				LowStockThreshold: 5, // From merged model
				BackorderAllowed:  false,
			}
			if err := tx.Create(&inventory).Error; err != nil {
				return fmt.Errorf("failed to create simple inventory: %w", err)
			}
			product.SimpleInventory = &inventory
		} else {
			// Create variants and their inventories
			for i := range variants {
				variants[i].ProductID = product.ID
				if err := tx.Create(&variants[i]).Error; err != nil {
					return fmt.Errorf("failed to create variant: %w", err)
				}
				inventory := models.Inventory{
					//ProductID:         &product.ID, // Explicit link to product (per requirement)
					VariantID:         &variants[i].ID,
					MerchantID:        product.MerchantID,
					Quantity:          variantInputs[i].InitialStock,
					ReservedQuantity:  0,
					LowStockThreshold: 5,
					BackorderAllowed:  false,
				}
				if err := tx.Create(&inventory).Error; err != nil {
					return fmt.Errorf("failed to create variant inventory: %w", err)
				}
				variants[i].Inventory = inventory
			}
		}

		// Create media
		for i := range media {
			media[i].ProductID = product.ID
			if err := tx.Create(&media[i]).Error; err != nil {
				return fmt.Errorf("failed to create media: %w", err)
			}
		}

		// Reload with preloads
		preloadQuery := tx.Where("id = ? AND deleted_at IS NULL", product.ID)
		if !isSimple {
			preloadQuery = preloadQuery.Preload("Variants.Inventory")
		}
		preloadQuery = preloadQuery.Preload("SimpleInventory").Preload("Media")
		if err := preloadQuery.First(product).Error; err != nil {
			return fmt.Errorf("failed to preload associations: %w", err)
		}

		return nil
	})
}




func (r *ProductRepository) UpdateInventoryQuantity(inventoryID string, delta int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var inventory models.Inventory
		if err := tx.First(&inventory, "id = ?", inventoryID).Error; err != nil {
			return fmt.Errorf("failed to find inventory: %w", err)
		}
		newQuantity := inventory.Quantity + delta
		if newQuantity < 0 && !inventory.BackorderAllowed {
			return errors.New("insufficient stock and backorders not allowed")
		}
		inventory.Quantity = newQuantity
		return tx.Save(&inventory).Error
	})
}

func (r *ProductRepository) SoftDeleteProduct(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Product{}).Error
}





//For media uploads
func (r *ProductRepository) CreateMedia(ctx context.Context, media *models.Media) error {
	return r.db.WithContext(ctx).Create(media).Error
}

// FindMediaByID fetches media
func (r *ProductRepository) FindMediaByID(ctx context.Context, id string) (*models.Media, error) {
	var media models.Media
	err := r.db.WithContext(ctx).Scopes(r.activeScope()).First(&media, "id = ?", id).Error
	return &media, err
}

// UpdateMedia updates fields
func (r *ProductRepository) UpdateMedia(ctx context.Context, id string, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&models.Media{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteMedia soft-deletes
func (r *ProductRepository) DeleteMedia(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&models.Media{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

// activeScope (if soft delete)
func (r *ProductRepository) activeScope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB { return db.Where("deleted_at IS NULL") }
}
```

---
### db/repositories/product_repositry.go
- Size: 5.06 KB
- Lines: 148
- Last Modified: 2025-09-30 12:28:15

```go
package repositories

// import (
// 	"api-customer-merchant/internal/db"
// 	"api-customer-merchant/internal/db/models"
// 	"errors"

// 	"gorm.io/gorm"
// )

/*
type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{db: db.DB}
}

// Create adds a new product to the database
func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

// FindByID retrieves a product by ID with associated Merchant and Category
func (r *ProductRepository) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Merchant").Preload("Category").First(&product, id).Error
	return &product, err
}

// FindBySKU retrieves a product by SKU with associated Merchant and Category
func (r *ProductRepository) FindBySKU(sku string) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Merchant").Preload("Category").Where("sku = ?", sku).First(&product).Error
	return &product, err
}

// SearchByName retrieves products matching a name (partial match)
func (r *ProductRepository) SearchByName(name string) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Merchant").Preload("Category").Where("name ILIKE ?", "%"+name+"%").Find(&products).Error
	return products, err
}

// FindByMerchantID retrieves all products for a merchant
func (r *ProductRepository) FindByMerchantID(merchantID uint) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Category").Where("merchant_id = ?", merchantID).Find(&products).Error
	return products, err
}

// FindByCategoryID retrieves all products in a category
func (r *ProductRepository) FindByCategoryID(categoryID uint) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Merchant").Where("category_id = ?", categoryID).Find(&products).Error
	return products, err
}

// Update modifies an existing product
func (r *ProductRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

// Delete removes a product by ID
func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}
// In ProductRepository
func (r *ProductRepository) FindByCategoryWithPagination(categoryID uint, limit, offset int) ([]models.Product, error) {
    var products []models.Product
    err := r.db.Preload("Merchant").Preload("Category").Where("category_id = ?", categoryID).Limit(limit).Offset(offset).Find(&products).Error
    return products, err
}
*/

/*
type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{db: db.DB}
}

// Create adds a new product to the database
func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

// FindByID retrieves a product by ID with associated Merchant and Category
func (r *ProductRepository) FindByID(id string) (*models.Product, error) {
    var product models.Product
    err := r.db.Preload("Merchant").Preload("Category").First(&product, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("product not found")
        }
        return nil, err
    }
    return &product, nil
}

// FindBySKU retrieves a product by SKU with associated Merchant and Category
func (r *ProductRepository) FindBySKU(sku string) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Merchant").Preload("Category").Where("sku = ?", sku).First(&product).Error
	return &product, err
}

// SearchByName retrieves products matching a name (partial match)
func (r *ProductRepository) SearchByName(name string) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Merchant").Preload("Category").Where("name ILIKE ?", "%"+name+"%").Find(&products).Error
	return products, err
}

// FindByMerchantID retrieves all products for a merchant
func (r *ProductRepository) FindByMerchantID(merchantID string) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Category").Where("merchant_id = ?", merchantID).Find(&products).Error
	return products, err
}

// FindByCategoryID retrieves all products in a category
func (r *ProductRepository) FindByCategoryID(categoryID string) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Merchant").Where("category_id = ?", categoryID).Find(&products).Error
	return products, err
}

// Update modifies an existing product
func (r *ProductRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

// Delete removes a product by ID
func (r *ProductRepository) Delete(id string) error {
	return r.db.Delete(&models.Product{}, id).Error
}
// In ProductRepository
func (r *ProductRepository) FindByCategoryWithPagination(categoryID string, limit, offset int) ([]models.Product, error) {
    var products []models.Product
    err := r.db.Preload("Merchant").Preload("Category").Where("category_id = ?", categoryID).Limit(limit).Offset(offset).Find(&products).Error
    return products, err
}

*/

```

---
### db/repositories/user_repository.go
- Size: 1.12 KB
- Lines: 53
- Last Modified: 2025-09-30 12:28:15

```go
package repositories

import (
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	"log"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{db: db.DB}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		log.Printf("Failed to find user by email %s: %v", email, err)
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByGoogleID(googleID string) (*models.User, error) {
	var user models.User
	err := r.db.Where("google_id = ?", googleID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

```

---
### db/repositories/variant_repository.go
- Size: 0.50 KB
- Lines: 18
- Last Modified: 2025-09-30 12:22:22

```go
package repositories

import (
	//"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"

	"context"

	"gorm.io/gorm"
)

type VariantRepository struct{ db *gorm.DB }

func NewVariantRepository(db *gorm.DB) *VariantRepository { return &VariantRepository{db} }
func (r *VariantRepository) FindByID(ctx context.Context, id string) (*models.Variant, error) {
	var variant models.Variant
	return &variant, r.db.WithContext(ctx).Preload("Inventory").First(&variant, "id = ?", id).Error
}

```

---
### services/cart/cart_service.go
- Size: 21.08 KB
- Lines: 727
- Last Modified: 2025-09-30 12:28:15

```go
package cart

import (
	//"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/api/dto" // Assuming dto.BulkUpdateRequest is defined here with ProductID string, Quantity int
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"
	"context"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
	//"gorm.io/gorm/logger"
)

var (
	ErrInvalidUserID     = errors.New("invalid user ID")
	ErrInvalidQuantity   = errors.New("quantity must be positive")
	ErrProductNotFound   = errors.New("product not found")
	ErrInventoryNotFound = errors.New("inventory not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrTransactionFailed   = errors.New("transaction failed")
)

type CartService struct {
	cartRepo      *repositories.CartRepository
	cartItemRepo  *repositories.CartItemRepository
	productRepo   *repositories.ProductRepository
	inventoryRepo *repositories.InventoryRepository
	logger        *zap.Logger
	validator     *validator.Validate
}

func NewCartService(cartRepo *repositories.CartRepository, cartItemRepo *repositories.CartItemRepository, productRepo *repositories.ProductRepository, inventoryRepo *repositories.InventoryRepository, logger *zap.Logger) *CartService {
	return &CartService{
		cartRepo:      cartRepo,
		cartItemRepo:  cartItemRepo,
		productRepo:   productRepo,
		inventoryRepo: inventoryRepo,
		logger:        logger,
		validator:     validator.New(),
	}
}

// GetActiveCart retrieves or creates an active cart for a user
func (s *CartService) GetActiveCart(ctx context.Context, userID uint) (*dto.CartResponse, error) {
	if userID == 0 {
		return nil, ErrInvalidUserID
	}
	cart, err := s.cartRepo.FindActiveCart(ctx, userID)
	// if err == nil {
	// 	return cart, nil
	// }
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("Failed to find active cart", zap.Uint("user_id", userID), zap.Error(err))
		return nil, fmt.Errorf("db error: %w", err)
	}

	cart = &models.Cart{UserID: userID, Status: models.CartStatusActive}
	if err := s.cartRepo.Create(ctx, cart); err != nil {
		s.logger.Error("Failed to create cart", zap.Error(err))
		return nil, fmt.Errorf("create failed: %w", err)
	}
	//return s.cartRepo.FindByID(ctx, cart.ID)
	cart, err = s.cartRepo.FindByID(ctx, cart.ID)
	if err != nil {
		s.logger.Error("Failed to get active cart", zap.Error(err))
		return nil, fmt.Errorf("failed to get active cart: %w", err)
	}
	response := &dto.CartResponse{
	ID:        cart.ID,
	UserID:    cart.UserID,
	Status:    cart.Status,
	Items:     make([]dto.CartItemResponse, len(cart.CartItems)),
	Total:     cart.GrandTotal, // Assuming decimal.Decimal
	CreatedAt: cart.CreatedAt,
	UpdatedAt: cart.UpdatedAt,
}
for i, item := range cart.CartItems {
	response.Items[i] = dto.CartItemResponse{
			ID:        item.ID,
		ProductID: item.ProductID,
		VariantID: item.VariantID, // Fixed from m.URL
		Quantity:  item.Quantity,
		Subtotal:  item.Cart.SubTotal,
	}
}
return response, nil

}

// func (s *CartService) GetCart(ctx context.Context, userID uint) (*models.Cart, error) {
// 	cart, err := s.GetActiveCart(ctx, userID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// Fixed: Use model method
// 	cart.ComputeTotals()
// 	s.logger.Info("Cart fetched", zap.Uint("user_id", userID), zap.Float64("total", cart.GrandTotal))
// 	return cart, nil
// }

// AddItemToCart adds a product to the user's active cart
/*
func (s *CartService) AddItemToCart(ctx context.Context, userID uint, quantity uint, productID string) (*models.Cart, error) {
	if userID == 0 {
		return nil, ErrInvalidUserID
	}
	if productID == "" {
		return nil, errors.New("invalid product ID")
	}
	if quantity == 0 {
		return nil, ErrInvalidQuantity
	}

	cart, err := s.GetActiveCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return nil, ErrProductNotFound
	}

	inventory, err := s.inventoryRepo.FindByProductID(ctx,productID,product.MerchantID)
	if err != nil {
		return nil, ErrInventoryNotFound
	}
	if inventory.StockQuantity < int(quantity) {
		return nil, ErrInsufficientStock
	}

	err = db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		cartItem, err := s.cartItemRepo.FindByProductIDAndCartID(ctx, productID, cart.ID)
		newQty := quantity
		if err == nil {
			newQty += uint(cartItem.Quantity)
			if inventory.StockQuantity < int(newQty) {
				return ErrInsufficientStock
			}
			if err := s.cartItemRepo.UpdateQuantityWithReservation(ctx, cartItem.ID, int(newQty), inventory.ID); err != nil {
				return err
			}
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			cartItem = &models.CartItem{
				CartID:     cart.ID,
				ProductID:  productID,
				Quantity:   int(quantity),
				MerchantID: product.MerchantID,
			}
			if err := s.cartItemRepo.Create(ctx, cartItem); err != nil {
				return err
			}
			if err := s.inventoryRepo.UpdateInventoryQuantity(ctx, inventory.ID, -int(quantity)); err != nil {
				return err
			}
		} else {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	// Reload
	return s.cartRepo.FindByID(ctx, cart.ID)
}
*/


/*
func (s *CartService) AddItemToCart(ctx context.Context, userID uint, quantity uint, productID string) (*models.Cart, error) {
	if userID == 0 {
		return nil, ErrInvalidUserID
	}
	if quantity == 0 {
		return nil, ErrInvalidQuantity
	}

	cart, err := s.GetActiveCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	product, err := s.productRepo.FindByID(ctx,productID)
	if err != nil {
		return nil, err
	}

	inventory, err := s.inventoryRepo.FindByProductID(ctx, productID, product.MerchantID)
	if err != nil {
		return nil, ErrInventoryNotFound
	}

	if inventory.Quantity < int(quantity) {
		return nil, ErrInsufficientStock
	}

	// Assuming no variant for simplicity; pass nil for variantID
	existing, err := s.cartItemRepo.FindByProductIDAndCartID(ctx, productID, nil, cart.ID)
	if err == nil {
		// Existing item: increment quantity
		newQty := existing.Quantity + int(quantity)
		if newQty > inventory.Quantity {
			return nil, ErrInsufficientStock
		}
		err = s.cartItemRepo.UpdateQuantityWithReservation(ctx, existing.ID, newQty, inventory.ID)
		if err != nil {
			return nil, err
		}
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// New item
		cartItem := &models.CartItem{
			CartID:    cart.ID,
			ProductID: productID,
			VariantID: nil,
			Quantity:  int(quantity),
			//PriceSnapshot: product.BasePrice.InexactFloat64(),
			MerchantID: product.MerchantID,
		}
		if err := s.cartItemRepo.Create(ctx, cartItem); err != nil {
			return nil, err
		}
		if err := s.inventoryRepo.UpdateInventoryQuantity(ctx, inventory.ID, -int(quantity)); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	return s.cartRepo.FindByID(ctx, cart.ID)
}
*/




/*
func (s *CartService) AddItemToCart(ctx context.Context, userID uint, quantity uint, productID string, variantID *string) (*models.Cart, error) {
    if userID == 0 {
        return nil, ErrInvalidUserID
    }
    if productID == "" {
        return nil, errors.New("invalid product ID")
    }
    if quantity == 0 {
        return nil, ErrInvalidQuantity
    }

    cart, err := s.GetActiveCart(ctx, userID)
    if err != nil {
        return nil, err
    }

    product, err := s.productRepo.FindByID(ctx, productID)
    if err != nil {
        return nil, ErrProductNotFound
    }

    var inventory *models.VendorInventory
    var priceSnapshot decimal.Decimal = product.BasePrice  // Default to base

    if variantID != nil {
        // Variant product
        variant, err := s.variantRepo.FindByID(ctx, *variantID)
        if err != nil {
            return nil, ErrInvalidVariant
        }
        if variant.ProductID != product.ID {
            return nil, errors.New("variant does not belong to product")
        }
        inventory, err = s.inventoryRepo.FindByVariantID(ctx, *variantID, product.MerchantID)
        if err != nil {
            return nil, ErrInventoryNotFound
        }
        priceSnapshot = product.BasePrice.Add(variant.PriceAdjustment)
    } else {
        // Simple product
        if len(product.Variants) > 0 {
            return nil, errors.New("variant required for this product")
        }
        inventory, err = s.inventoryRepo.FindByProductID(ctx, productID, product.MerchantID)
        if err != nil {
            return nil, ErrInventoryNotFound
        }
    }

    // Quick pre-check
    available := inventory.Quantity - inventory.ReservedQuantity
    if available < int(quantity) {
        return nil, ErrInsufficientStock
    }

    // Transaction
    err = db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        var freshInv *models.VendorInventory
        var freshErr error

        if variantID != nil {
            // Re-fetch variant inventory inside tx for freshness
            freshInv, freshErr = s.inventoryRepo.FindByVariantID(tx, *variantID, product.MerchantID)
        } else {
            // Re-fetch product inventory inside tx for freshness
            freshInv, freshErr = s.inventoryRepo.FindByProductID(tx, productID, product.MerchantID)
        }

        if freshErr != nil {
            return freshErr
        }
        available = freshInv.Quantity - freshInv.ReservedQuantity
        if available < int(quantity) {
            return ErrInsufficientStock
        }

        // Existing item check (use tx)
        cartItem, err := s.cartItemRepo.FindByProductAndVariant(tx, productID, variantID, cart.ID)
        newQty := int(quantity)

        if err == nil {
            newQty += cartItem.Quantity
            if available < newQty {
                return ErrInsufficientStock
            }
            // Update with price if needed (for MVP, assume snapshot set on create)
            return s.cartItemRepo.UpdateQuantityWithReservation(tx, cartItem.ID, newQty, freshInv.ID)
        } else if errors.Is(err, gorm.ErrRecordNotFound) {
            cartItem = &models.CartItem{
                CartID:        cart.ID,
                ProductID:     productID,
                VariantID:     variantID,  // Nil-safe
                Quantity:      newQty,
                PriceSnapshot: priceSnapshot,  // Set once
                MerchantID:    product.MerchantID,
            }
            if err := s.cartItemRepo.Create(tx, cartItem); err != nil {
                return err
            }
            // Reserve (for create, delta = quantity)
            return s.inventoryRepo.ReserveStock(tx, freshInv.ID, newQty)
        }
        return err
    })
    if err != nil {
        s.logger.Error("Add to cart failed", zap.Error(err), zap.Uint("user_id", userID))
        return nil, err
    }

    return s.cartRepo.FindByID(ctx, cart.ID)
}


*/




















func (s *CartService) AddItemToCart(ctx context.Context, userID uint, quantity int, productID string, variantID *string) (*dto.CartResponse, error) {
	if userID == 0 {
		return nil, ErrInvalidUserID
	}
	if quantity <= 0 {
		return nil, ErrInvalidQuantity
	}

	cart, err := s.GetActiveCart(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get active cart", zap.Uint("user_id", userID), zap.Error(err))
		return nil, err
	}

	// Fetch product with preloaded Variants.Inventory and SimpleInventory
	product, err := s.productRepo.FindByID(ctx, productID, "Variants.Inventory", "SimpleInventory")
	if err != nil {
		s.logger.Error("Product not found", zap.String("product_id", productID), zap.Error(err))
		return nil, ErrProductNotFound
	}
	if product.DeletedAt.Valid {
		s.logger.Error("Product is soft-deleted", zap.String("product_id", productID))
		return nil, ErrProductNotFound
	}

	// Determine inventory: focus on variants if they exist, else simple
	var inventory *models.Inventory
	var price decimal.Decimal = product.BasePrice
	//var varID string
	if variantID != nil && len(product.Variants) > 0 {
		//varID = *variantID
		for _, v := range product.Variants {
			if v.ID == *variantID && v.IsActive {
				inventory = &v.Inventory
				price = price.Add(v.PriceAdjustment)
				break
			}
		}
	} else if variantID == nil && product.SimpleInventory != nil {
		inventory = product.SimpleInventory
	} else {
		s.logger.Error("Inventory not found", zap.String("product_id", productID), zap.Stringp("variant_id", variantID))
		return nil, ErrInventoryNotFound
	}
	if inventory == nil {
		s.logger.Error("No valid inventory", zap.String("product_id", productID), zap.Stringp("variant_id", variantID))
		return nil, ErrInventoryNotFound
	}

	// Check available stock
	available := inventory.Quantity - inventory.ReservedQuantity
	if available < quantity {
		s.logger.Warn("Insufficient stock", zap.Int("available", available), zap.Int("requested", quantity))
		return nil, ErrInsufficientStock
	}

	// Transaction: Update cart item and reserve inventory
	err = db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Find existing cart item
		existing, err := s.cartItemRepo.FindByProductIDAndCartID(ctx, productID, nil,cart.ID)
		if err == nil {
			// Update existing item
			newQty := existing.Quantity + quantity
			if newQty > available {
				return ErrInsufficientStock
			}
			if err := s.cartItemRepo.UpdateQuantityWithReservation(ctx, existing.ID, newQty, inventory.ID); err != nil {
				return fmt.Errorf("failed to update cart item: %w", err)
			}
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("failed to check existing cart item: %w", err)
		}

		// Create new cart item
		cartItem := &models.CartItem{
			CartID:    cart.ID,
			ProductID: productID,
			VariantID: variantID, // Assume VariantID is *string in model
			Quantity:  quantity,
			MerchantID: product.MerchantID,
		}
		if err := s.cartItemRepo.Create(ctx, cartItem); err != nil {
			return fmt.Errorf("failed to create cart item: %w", err)
		}
		inventory.ReservedQuantity += quantity // Manual update since method undefined
		if err := tx.Save(inventory).Error; err != nil {
			return fmt.Errorf("failed to reserve inventory: %w", err)
		}
		return nil
	})
	if err != nil {
		s.logger.Error("Transaction failed", zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrTransactionFailed, err)
	}

	// Return updated cart
	updatedCart, err := s.cartRepo.FindByID(ctx, cart.ID)
    if err != nil {
        s.logger.Error("Failed to fetch updated cart", zap.Uint("cart_id", cart.ID), zap.Error(err))
        return nil, err
    }
    // Fix: Preload CartItems with related data
    if err := db.DB.WithContext(ctx).
        Preload("CartItems.Product.Media").
        Preload("CartItems.Product.Variants.Inventory").
        Preload("CartItems.Variant").
        Find(updatedCart).Error; err != nil {
        s.logger.Error("Failed to preload cart items", zap.Error(err))
        return nil, err
    }
   // return updatedCart, nil
	response := &dto.CartResponse{
		ID: updatedCart.ID,
		UserID: updatedCart.UserID,
		Status: updatedCart.Status,
		Items: make([]dto.CartItemResponse, len(updatedCart.CartItems)),
		Total: updatedCart.GrandTotal,
		CreatedAt: updatedCart.CreatedAt,
		UpdatedAt: updatedCart.UpdatedAt,
	}
	for i, item := range updatedCart.CartItems {
	response.Items[i] = dto.CartItemResponse{
		ID:        item.ID,
		ProductID: item.ProductID,
		VariantID: item.VariantID, // Fixed from m.URL
		Quantity:  item.Quantity,
		Subtotal:  item.Cart.SubTotal,

	}
}
return response, nil
}











// UpdateCartItemQuantity updates the quantity of a cart item
func (s *CartService) UpdateCartItemQuantity(ctx context.Context, cartItemID uint, quantity int) (*models.Cart, error) {
	if cartItemID == 0 {
		return nil, errors.New("invalid cart item ID")
	}
	if quantity <= 0 {
		return nil, ErrInvalidQuantity
	}

	// load cart item (contains MerchantID and ProductID)
	cartItem, err := s.cartItemRepo.FindByID(ctx, cartItemID)
	if err != nil {
		return nil, repositories.ErrCartItemNotFound
	}

	// ensure we have a merchantID to scope the inventory lookup
	merchantID := cartItem.MerchantID
	if merchantID == "" {
		// fallback: fetch product to get merchant (shouldn't usually happen if cart items store merchant)
		prod, perr := s.productRepo.FindByID(ctx,cartItem.ProductID)
		if perr != nil {
			return nil, ErrInventoryNotFound
		}
		merchantID = prod.MerchantID
	}

	// NOTE: FindByProductID signature is (ctx, productID, merchantID)
	inventory, err := s.inventoryRepo.FindByProductID(ctx, cartItem.ProductID, merchantID)
	if err != nil {
		return nil, ErrInventoryNotFound
	}

	// model field is Quantity (not StockQuantity)
	if inventory.Quantity < quantity {
		return nil, ErrInsufficientStock
	}

	// UpdateQuantityWithReservation now expects vendor inventory ID as string
	if err := s.cartItemRepo.UpdateQuantityWithReservation(ctx, cartItemID, quantity, inventory.ID); err != nil {
		return nil, err
	}

	return s.cartRepo.FindByID(ctx, cartItem.CartID)
}

func (s *CartService) RemoveCartItem(ctx context.Context, cartItemID uint) (*models.Cart, error) {
	if cartItemID == 0 {
		return nil, errors.New("invalid cart item ID")
	}

	cartItem, err := s.cartItemRepo.FindByID(ctx, cartItemID)
	if err != nil {
		return nil, repositories.ErrCartItemNotFound
	}

	// use the merchant stored on the cart item to find the correct vendor inventory
	merchantID := cartItem.MerchantID
	if merchantID == "" {
		// fallback: fetch product to get merchant
		prod, perr := s.productRepo.FindByID(ctx,cartItem.ProductID)
		if perr != nil {
			return nil, ErrInventoryNotFound
		}
		merchantID = prod.MerchantID
	}

	// pass ctx and merchantID as required by repo
	inventory, err := s.inventoryRepo.FindByProductID(ctx, cartItem.ProductID, merchantID)
	if err != nil {
		return nil, ErrInventoryNotFound
	}

	// DeleteWithUnreserve expects vendor inventory ID as string
	if err := s.cartItemRepo.DeleteWithUnreserve(ctx, cartItemID, inventory.ID); err != nil {
		return nil, err
	}

	return s.cartRepo.FindByID(ctx, cartItem.CartID)
}





func (s *CartService) GetCartItemByID(ctx context.Context, cartItemID uint) (*models.CartItem, error) {
	if cartItemID == 0 {
		return nil, errors.New("invalid cart item ID")
	}
	return s.cartItemRepo.FindByID(ctx, cartItemID)
}





// ClearCart, BulkAddItems ... (add ctx to all calls; stub Bulk if not used)
func (s *CartService) ClearCart(ctx context.Context, userID uint) error {
	cart, err := s.cartRepo.FindActiveCart(ctx, userID)
	if err != nil {
		return err
	}
	items, err := s.cartItemRepo.FindByCartID(ctx, cart.ID)
	if err != nil {
		return err
	}
	for _, item := range items {
		s.RemoveCartItem(ctx, item.ID)
	}
	cart.Status = models.CartStatusAbandoned
	return s.cartRepo.Update(ctx, cart)
}

// BulkAddItems stub (implement as needed; fixed DTO)
// func (s *CartService) BulkAddItems(ctx context.Context, userID uint, items []dto.BulkUpdateRequest) (*models.Cart, error) {
// 	// Validation loop...
// 	for _, item := range items {
// 		if err := s.validator.Struct(&item); err != nil {
// 			return nil, err
// 		}
// 		// Add each (loop AddItemToCart)
// 		_, err := s.AddItemToCart(ctx, userID, uint(item.Quantity), item.ProductID)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	return s.GetCart(ctx, userID)
// }












// func (s *CartService) BulkAddItems(ctx context.Context, userID uint, items dto.BulkUpdateRequest) (*models.Cart, error) {
// 	if userID == 0 {
// 		return nil, ErrInvalidUserID
// 	}
// 	if err := s.validator.Struct(&items); err != nil {
// 		return nil, fmt.Errorf("validation failed: %w", err)
// 	}

// 	// cart, err := s.GetActiveCart(ctx, userID)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	for _, item := range items.Items {
// 		// Convert uint ProductID to string for consistency
// 		productID := fmt.Sprint(item.ProductID)
// 		if _, err := s.AddItemToCart(ctx, userID, uint(item.Quantity), productID); err != nil {
// 			return nil, fmt.Errorf("failed to add item %s: %w", productID, err)
// 		}
// 	}
// 	return s.GetCart(ctx, userID)
// }






func (s *CartService) BulkAddItems(ctx context.Context, userID uint, items dto.BulkUpdateRequest) (*models.Cart, error) {
	if userID == 0 {
		return nil, ErrInvalidUserID
	}
	if len(items.Items) == 0 {
		return nil, errors.New("no items provided")
	}
	if err := s.validator.Struct(&items); err != nil {
		s.logger.Error("Validation failed", zap.Uint("user_id", userID), zap.Error(err))
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	cart, err := s.GetActiveCart(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get active cart", zap.Uint("user_id", userID), zap.Error(err))
		return nil, err
	}

	err = db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i, item := range items.Items {
			if _, err := s.AddItemToCart(ctx, userID, item.Quantity, item.ProductID, item.VariantID); err != nil {
				s.logger.Error("Failed to add item", zap.String("product_id", item.ProductID), zap.Stringp("variant_id", item.VariantID), zap.Error(err))
				return fmt.Errorf("failed to add item %d (product %s): %w", i+1, item.ProductID, err)
			}
		}
		return nil
	})
	if err != nil {
		s.logger.Error("Transaction failed", zap.Uint("user_id", userID), zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrTransactionFailed, err)
	}

	updatedCart, err := s.cartRepo.FindByID(ctx, cart.ID)
	if err != nil {
		s.logger.Error("Failed to fetch updated cart", zap.Uint("cart_id", cart.ID), zap.Error(err))
		return nil, err
	}
	return updatedCart, nil
}
```

---
### services/dispute/dispute_service.go
- Size: 2.63 KB
- Lines: 106
- Last Modified: 2025-10-02 13:34:34

```go
package dispute
import (
	"context"
	"errors"
	"fmt"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidInput = errors.New("invalid input")
)

type DisputeService struct {
	disputeRepo *repositories.DisputeRepository
	orderRepo   *repositories.OrderRepository
	logger      *zap.Logger
}

func NewDisputeService(
	disputeRepo *repositories.DisputeRepository,
	orderRepo *repositories.OrderRepository,
	logger *zap.Logger,
) *DisputeService {
	return &DisputeService{
		disputeRepo: disputeRepo,
		orderRepo:   orderRepo,
		logger:      logger,
	}
}

// CreateDispute creates a new dispute
func (s *DisputeService) CreateDispute(ctx context.Context, userID uint, req dto.CreateDisputeDTO) (*dto.DisputeResponseDTO, error) {
	logger := s.logger.With(zap.String("operation", "CreateDispute"), zap.Uint("user_id", userID))

	// Verify order exists and belongs to user
	order, err := s.orderRepo.FindByID(ctx, parseUint(req.OrderID))
	if err != nil {
		logger.Error("Order not found", zap.Error(err))
		return nil, fmt.Errorf("order not found: %w", err)
	}

	if order.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Create dispute model
	dispute := &models.Dispute{
		ID:          uuid.NewString(),
		OrderID:     req.OrderID,
		CustomerID:  userID,
		MerchantID:  order.OrderItems[0].MerchantID, // Assume first item's merchant
		Reason:      req.Reason,
		Description: req.Description,
		Status:      "open",
	}

	if err := s.disputeRepo.Create(ctx, dispute); err != nil {
		logger.Error("Failed to create dispute", zap.Error(err))
		return nil, err
	}

	return mapDisputeToDTO(dispute), nil
}

// GetDispute retrieves a dispute by ID
func (s *DisputeService) GetDispute(ctx context.Context, disputeID string, userID uint) (*dto.DisputeResponseDTO, error) {
	dispute, err := s.disputeRepo.FindDisputeByID(ctx, disputeID)
	if err != nil {
		return nil, err
	}

	if dispute.CustomerID != userID {
		return nil, ErrUnauthorized
	}

	return mapDisputeToDTO(dispute), nil
}

func mapDisputeToDTO(d *models.Dispute) *dto.DisputeResponseDTO {
	return &dto.DisputeResponseDTO{
		ID:          d.ID,
		OrderID:     d.OrderID,
		CustomerID:  d.CustomerID,
		MerchantID:  d.MerchantID,
		Reason:      d.Reason,
		Description: d.Description,
		Status:      d.Status,
		Resolution:  d.Resolution,
		CreatedAt:   d.CreatedAt,
		ResolvedAt:  d.ResolvedAt,
	}
}

func parseUint(s string) uint {
	// Simplified parser (add robust error handling in production)
	var id uint
	fmt.Sscanf(s, "%d", &id)
	return id
}
```

---
### services/merchant/merchant_service.go
- Size: 8.31 KB
- Lines: 306
- Last Modified: 2025-10-02 12:36:21

```go
package merchant

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/bank"

	//"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"

	//"github.com/gray-adeyi/paystack"

	"golang.org/x/crypto/bcrypt"
)

/*
	type MerchantService struct {
		appRepo  *repositories.MerchantApplicationRepository
		repo     *repositories.MerchantRepository
		validate *validator.Validate
	}

	func NewMerchantService(appRepo *repositories.MerchantApplicationRepository, repo *repositories.MerchantRepository) *MerchantService {
		return &MerchantService{
			appRepo:  appRepo,
			repo:     repo,
			validate: validator.New(),
		}
	}

// SubmitApplication allows a prospective merchant to submit an application.

	func (s *MerchantService) SubmitApplication(ctx context.Context, app *models.MerchantApplication) (*models.MerchantApplication, error) {
		// Validate required fields
		if err := s.validate.Struct(app); err != nil {
			return nil, errors.New("validation failed: " + err.Error())
		}

		// Validate JSONB fields
		if len(app.PersonalAddress) == 0 {
			return nil, errors.New("personal_address cannot be empty")
		}
		if len(app.WorkAddress) == 0 {
			return nil, errors.New("work_address cannot be empty")
		}
		var temp map[string]interface{}
		if err := json.Unmarshal(app.PersonalAddress, &temp); err != nil {
			return nil, errors.New("invalid personal_address JSON")
		}
		if err := json.Unmarshal(app.WorkAddress, &temp); err != nil {
			return nil, errors.New("invalid work_address JSON")
		}




		// Set Status to pending and ensure ID is not set
		app.Status = "pending"
		app.ID = ""

		if err := s.appRepo.Create(ctx, app); err != nil {
			return nil, err
		}
		return app, nil
	}

// GetApplication returns an application by ID (for applicant to check their own status).

	func (s *MerchantService) GetApplication(ctx context.Context, id string) (*models.MerchantApplication, error) {
		if id == "" {
			return nil, errors.New("application ID cannot be empty")
		}
		return s.appRepo.GetByID(ctx, id)
	}

// GetMerchantByUserID returns an active merchant account for the authenticated user.

	func (s *MerchantService) GetMerchantByUserID(ctx context.Context, uid string) (*models.Merchant, error) {
		if uid == "" {
			return nil, errors.New("user ID cannot be empty")
		}
		return s.repo.GetByUserID(ctx, uid)
	}

// GetMerchantByID returns an active merchant by ID.

	func (s *MerchantService) GetMerchantByID(ctx context.Context, id string) (*models.Merchant, error) {
		if id == "" {
			return nil, errors.New("merchant ID cannot be empty")
		}
		return s.repo.GetByID(ctx, id)
	}
*/
type MerchantService struct {
	appRepo  *repositories.MerchantApplicationRepository
	repo     *repositories.MerchantRepository
	validate *validator.Validate
}

func NewMerchantService(appRepo *repositories.MerchantApplicationRepository, repo *repositories.MerchantRepository) *MerchantService {
	return &MerchantService{
		appRepo:  appRepo,
		repo:     repo,
		validate: validator.New(),
	}
}

// SubmitApplication allows a prospective merchant to submit an application.
func (s *MerchantService) SubmitApplication(ctx context.Context, app *models.MerchantApplication) (*models.MerchantApplication, error) {
	// Validate required fields
	if err := s.validate.Struct(app); err != nil {
		return nil, errors.New("validation failed: " + err.Error())
	}

	// Validate JSONB fields
	if len(app.PersonalAddress) == 0 {
		return nil, errors.New("personal_address cannot be empty")
	}
	if len(app.WorkAddress) == 0 {
		return nil, errors.New("work_address cannot be empty")
	}
	var temp map[string]interface{}
	if err := json.Unmarshal(app.PersonalAddress, &temp); err != nil {
		return nil, errors.New("invalid personal_address JSON")
	}
	if err := json.Unmarshal(app.WorkAddress, &temp); err != nil {
		return nil, errors.New("invalid work_address JSON")
	}

	// Set Status to pending and ensure ID is not set
	app.Status = "pending"
	app.ID = ""

	if err := s.appRepo.Create(ctx, app); err != nil {
		return nil, err
	}
	return app, nil
}

// GetApplication returns an application by ID (for applicant to check their own status).
func (s *MerchantService) GetApplication(ctx context.Context, id string) (*models.MerchantApplication, error) {
	if id == "" {
		return nil, errors.New("application ID cannot be empty")
	}
	return s.appRepo.GetByID(ctx, id)
}

// GetMerchantByUserID returns an active merchant account for the authenticated user.
func (s *MerchantService) GetMerchantByUserID(ctx context.Context, uid string) (*models.Merchant, error) {
	if uid == "" {
		return nil, errors.New("user ID cannot be empty")
	}
	return s.repo.GetByMerchantID(ctx, uid)
}

// GetMerchantByID returns an active merchant by ID.
func (s *MerchantService) GetMerchantByID(ctx context.Context, id string) (*models.Merchant, error) {
	if id == "" {
		return nil, errors.New("merchant ID cannot be empty")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *MerchantService) LoginMerchant(ctx context.Context, work_email, password string) (*models.Merchant, error) {
	merchant, err := s.repo.GetByWorkEmail(ctx, work_email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(merchant.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return merchant, nil
}

func (s *MerchantService) GenerateJWT(entity interface{}) (string, error) {
	var id string
	var entityType string

	switch e := entity.(type) {
	case *models.Merchant:
		id = e.ID
		entityType = "merchant"

	default:
		return "", errors.New("invalid entity type")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         id,
		"entityType": entityType,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET not set")
	}

	return token.SignedString([]byte(secret))
}













 //func (s *MerchantService) AddBankDetails(merchantID string, details MerchantBankDetails) error {
     // Validate inputs (e.g., bank code format)
   //  details.MerchantID = merchantID
//     details.Status = "pending"

//     // Create Paystack recipient
//     verify_client := paystack.VerificationClient(config.PaystackSecretKey)  // Assume injected
// 	var response models.Response[models.BankAccountInfo]
// 	if err := client.Verification.ResolveAccount(context.TODO(), &response,p.WithQuery("account_number","0022728151"),p.WithQuery("bank_code","063")); err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(response)
// }
//     recipientReq := &verifyclient. {
//         Type:          "nuban",
//         Name:          details.AccountName,
//         AccountNumber: details.AccountNumber,
//         BankCode:      details.BankCode,
//         Currency:      details.Currency,
//     }
//     resp, err := paystack.Recipient.Create(recipientReq)
//     if err != nil {
//         return err
//     }
//     details.RecipientCode = resp.Data.RecipientCode
//     details.Status = "verified"  // If Paystack verifies

//     return db.DB.Create(&details).Error
// }


func (s *MerchantService) AddBankDetails(ctx context.Context, merchantID string, details dto.BankDetailsRequest) error {
	// Validate bank name
	bankSvc := bank.GetBankService()
	if err := bankSvc.LoadBanks(); err != nil {
		return fmt.Errorf("failed to load banks: %w", err)
	}

	bankCode, err := bankSvc.GetBankCode(details.BankName)
	if err != nil {
		return fmt.Errorf("invalid bank name: %w", err)
	}

	// Override with validated code
	details.BankCode = bankCode

	// Persist via repository
	if err := s.repo.UpdateBankDetails(ctx, merchantID, details); err != nil {
		return fmt.Errorf("failed to save bank details: %w", err)
	}

	return nil
}


func (s *MerchantService) UpdateBankDetails(ctx context.Context, merchantID string ,details  dto.BankDetailsRequest) error {
     // Similar, but use Save or Update
	 if details.BankName == "" {
		return errors.New("empty bank name")
	}

	if details.AccountNumber == "" {
		return errors.New("empty bank name")
	}
	

	

	err := s.repo.UpdateBankDetails(ctx ,merchantID, details)
	if err != nil {
		return  err
	}

	//payment.Status = models.PaymentStatus(status)
	

	return nil

  
 }
```

---
### services/notifications/notifcation_service.go
- Size: 0.71 KB
- Lines: 31
- Last Modified: 2025-09-30 12:22:22

```go
package notifications

import (
	"context"
	"fmt"
	// Assume SMTP or Twilio lib
)

type NotificationService struct {
	// Email/SMS clients
}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

func (s *NotificationService) SendOrderConfirmation(ctx context.Context, orderID uint, email string) error {
	// Use template: "Your order {orderID} is confirmed"
	fmt.Println("Sending email to", email) // Replace with real send
	return nil
}

func (s *NotificationService) SendMerchantNewOrder(merchantID uint, orderID uint) error {
	// Fetch merchant email, send
	return nil
}

func (s *NotificationService) SendStockAlert(merchantID uint, productID uint) error {
	// On low stock
	return nil
}

```

---
### services/order/order_service.go
- Size: 14.00 KB
- Lines: 431
- Last Modified: 2025-09-30 12:28:15

```go
package order

import (
	"context"
	"errors"
	"fmt"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"

	//"go.uber.org/zap"
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// OrderService provides business logic for handling orders.
type OrderService struct {
	orderRepo     *repositories.OrderRepository
	orderItemRepo *repositories.OrderItemRepository
	cartRepo      *repositories.CartRepository
	cartItemRepo  *repositories.CartItemRepository
	productRepo   *repositories.ProductRepository
	inventoryRepo *repositories.InventoryRepository
	logger      *zap.Logger
	validator   *validator.Validate
	db            *gorm.DB
}

// NewOrderService creates a new instance of OrderService.
func NewOrderService(
	orderRepo *repositories.OrderRepository,
	orderItemRepo *repositories.OrderItemRepository,
	cartRepo *repositories.CartRepository,
	cartItemRepo *repositories.CartItemRepository,
	productRepo *repositories.ProductRepository,
	inventoryRepo *repositories.InventoryRepository,
) *OrderService {
	return &OrderService{
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
		cartRepo:      cartRepo,
		cartItemRepo:  cartItemRepo,
		productRepo:   productRepo,
		inventoryRepo: inventoryRepo,
		db:            db.DB,
	}
}

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrInvalidOrderStatus = errors.New("order cannot be cancelled")
	ErrUnauthorizedOrder  = errors.New("unauthorized to cancel this order")
	ErrRefundFailed       = errors.New("failed to initiate refund")
	ErrNotificationFailed = errors.New("failed to send notification")
)

// CreateOrder converts a user's active cart into an order.
// It performs several actions within a single database transaction:
// 1. Finds the user's active cart.
// 2. Validates that the cart is not empty.
// 3. For each item in the cart, it moves the reserved stock to committed stock.
// 4. Creates an Order record.
// 5. Creates OrderItem records corresponding to the CartItems.
// 6. Deletes the cart items.
// 7. Updates the cart status to 'Converted'.
// 8. Returns a DTO representing the newly created order.
/*
func (s *OrderService) CreateOrder(ctx context.Context, userID uint) (*dto.OrderResponse, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}

	cart, err := s.cartRepo.FindActiveCart(ctx, userID)
	if err != nil {
		if errors.Is(err, repositories.ErrCartNotFound) {
			return nil, errors.New("no active cart found")
		}
		return nil, err
	}

	if len(cart.CartItems) == 0 {
		return nil, errors.New("cart is empty")
	}

	var newOrder *models.Order
	var totalAmount float64

	// Use a transaction to ensure atomicity
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Calculate total and create order items
		var orderItems []models.OrderItem
		for _, item := range cart.CartItems {
			price := item.Product.BasePrice.InexactFloat64()
			totalAmount += float64(item.Quantity) * price
			orderItems = append(orderItems, models.OrderItem{
				ProductID:  item.ProductID,
				MerchantID: item.Product.MerchantID,
				Quantity:   item.Quantity,
				Price:      price,
				FulfillmentStatus: models.FulfillmentStatusNew,
			})

			// Here you would typically move reserved stock to committed stock.
			// For now, we assume cart reservation handled this.
			// We'll just update the main inventory.
			// This logic might need to be more robust depending on inventory strategy.
			if err := tx.Model(&models.Inventory{}).
				Where("product_id = ?", item.ProductID).
				Updates(map[string]interface{}{
					"quantity":          gorm.Expr("quantity - ?", item.Quantity),
					"reserved_quantity": gorm.Expr("reserved_quantity - ?", item.Quantity),
				}).Error; err != nil {
				return fmt.Errorf("failed to commit stock for product %s: %w", item.ProductID, err)
			}
		}

		// Create the order
		newOrder = &models.Order{
			UserID:      userID,
			TotalAmount: totalAmount,
			Status:      models.OrderStatusPending,
			OrderItems:  orderItems,
		}
		if err := tx.Create(newOrder).Error; err != nil {
			return err
		}

		// Associate order items with the new order ID and create them
		for i := range orderItems {
			orderItems[i].OrderID = newOrder.ID
		}
		if err := tx.Create(&orderItems).Error; err != nil {
			return err
		}

		// Manually associate for the response DTO
		newOrder.OrderItems = orderItems

		// Clear cart items
		if err := tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
			return err
		}

		// Mark cart as converted
		cart.Status = models.CartStatusConverted
		if err := tx.Save(cart).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Convert to DTO for response
	orderResponse := &dto.OrderResponse{
		ID:         newOrder.ID,
		UserID:     newOrder.UserID,
		Status:     string(newOrder.Status),
		OrderItems: make([]dto.OrderItemResponse, len(newOrder.OrderItems)),
	}
	for i, item := range newOrder.OrderItems {
		orderResponse.OrderItems[i] = dto.OrderItemResponse{
			ProductID: fmt.Sprint(item.ProductID),
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	return orderResponse, nil
}
*/


func (s *OrderService) CreateOrder(ctx context.Context, userID uint) (*dto.OrderResponse, error) {
    if userID == 0 {
        return nil, errors.New("invalid user ID")
    }

    cart, err := s.cartRepo.FindActiveCart(ctx, userID)
    if err != nil {
        if errors.Is(err, repositories.ErrCartNotFound) {
            return nil, errors.New("no active cart found")
        }
        return nil, err
    }

    if len(cart.CartItems) == 0 {
        return nil, errors.New("cart is empty")
    }

    var newOrder *models.Order
    var totalAmount float64

    // Use a transaction to ensure atomicity
    err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // Calculate total and create order items
        var orderItems []models.OrderItem
        for _, item := range cart.CartItems {
            price := item.Product.BasePrice.InexactFloat64()
            if item.VariantID != nil && item.Variant != nil {
                price = item.Variant.TotalPrice.InexactFloat64() // Use variant price if available
            }
            totalAmount += float64(item.Quantity) * price
            orderItems = append(orderItems, models.OrderItem{
                ProductID:         item.ProductID,
                MerchantID:        item.Product.MerchantID,
                Quantity:          item.Quantity,
                Price:             price,
                FulfillmentStatus: models.FulfillmentStatusNew,
                // ID is omitted to let GORM/database auto-generate
            })

            // Update inventory: move reserved stock to committed stock
            inventoryQuery := tx.Model(&models.Inventory{}).Where("product_id = ?", item.ProductID)
            if item.VariantID != nil {
                inventoryQuery = inventoryQuery.Where("variant_id = ?", *item.VariantID)
            }
            if err := inventoryQuery.Updates(map[string]interface{}{
                "quantity":          gorm.Expr("quantity - ?", item.Quantity),
                "reserved_quantity": gorm.Expr("reserved_quantity - ?", item.Quantity),
            }).Error; err != nil {
                return fmt.Errorf("failed to commit stock for product %s: %w", item.ProductID, err)
            }
        }

        // Create the order
        newOrder = &models.Order{
            UserID:      userID,
            TotalAmount: decimal.NewFromFloat(totalAmount),
            Status:      models.OrderStatusPending,
            OrderItems:  orderItems,
        }
        if err := tx.Create(newOrder).Error; err != nil {
            return fmt.Errorf("failed to create order: %w", err)
        }

        // Associate order items with the new order ID
        for i := range orderItems {
            orderItems[i].OrderID = newOrder.ID
            orderItems[i].ID = 0 // Explicitly reset ID to ensure auto-generation
        }
        if err := tx.Create(&orderItems).Error; err != nil {
            return fmt.Errorf("failed to create order items: %w", err)
        }

        // Manually associate for the response DTO
        newOrder.OrderItems = orderItems

        // Clear cart items
        if err := tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
            return fmt.Errorf("failed to clear cart items: %w", err)
        }

        // Mark cart as converted
        cart.Status = models.CartStatusConverted
        if err := tx.Save(cart).Error; err != nil {
            return fmt.Errorf("failed to update cart status: %w", err)
        }

        return nil
    })

    if err != nil {
        //s.logger.Error("Transaction failed", zap.Error(err))
        return nil, fmt.Errorf("transaction failed: %w", err)
    }

    // Convert to DTO for response
    orderResponse := &dto.OrderResponse{
        ID:         newOrder.ID,
        UserID:     newOrder.UserID,
        Status:     string(newOrder.Status),
        OrderItems: make([]dto.OrderItemResponse, len(newOrder.OrderItems)),
    }
    for i, item := range newOrder.OrderItems {
        orderResponse.OrderItems[i] = dto.OrderItemResponse{
            ProductID: fmt.Sprint(item.ProductID),
            Quantity:  item.Quantity,
            Price:     item.Price,
        }
    }
	//  user, _ := userRepo.FindByID(userID)  // Assume userRepo injected
    // authURL, ref, err := paymentService.InitiateTransaction(ctx, order, user.Email)
    // if err != nil {
    //     return nil, err
    // }
    // Return order with authURL in response

    return orderResponse, nil
}



// GetOrder retrieves a single order by its ID.
func (s *OrderService) GetOrder(ctx context.Context, id uint) (*models.Order, error) {
	if id == 0 {
		return nil, errors.New("invalid order ID")
	}
	// The repository already preloads necessary associations.
	return s.orderRepo.FindByID(ctx, id)
	//return s.orderRepo.FindByID(id)
}




// GetOrdersByUserID retrieves all orders for a user
func (s *OrderService) GetOrdersByUserID(userID uint) ([]models.Order, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}
	return s.orderRepo.FindByUserID(userID)
}

// GetOrdersByMerchantID retrieves orders containing a merchant's products
func (s *OrderService) GetOrdersByMerchantID(merchantID uint) ([]models.Order, error) {
	if merchantID == 0 {
		return nil, errors.New("invalid merchant ID")
	}
	return s.orderRepo.FindByMerchantID(merchantID)
}

// UpdateOrderStatus updates the status of an order
func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID uint, status string) (*models.Order, error) {
	if orderID == 0 {
		return nil, errors.New("invalid order ID")
	}
	if err := models.OrderStatus(status).Valid(); err != nil {
		return nil, err
	}

	order, err := s.orderRepo.FindByID(ctx,orderID)
	if err != nil {
		return nil, err
	}

	order.Status = models.OrderStatus(status)
	if err := s.orderRepo.Update(ctx ,order); err != nil {
		return nil, err
	}

	return s.orderRepo.FindByID(ctx ,orderID)
}



// CancelOrder orchestrates cancellation (business logic here)
func (s *OrderService) CancelOrder(ctx context.Context, orderID uint, userID uint, reason string) error {
	logger := s.logger.With(zap.String("operation", "CancelOrder"), zap.Uint("order_id", orderID), zap.Uint("user_id", userID))

	// Fetch order (ownership checked in repo for efficiency)
	order, err := s.orderRepo.FindByIDWithPreloads(ctx, orderID)
	if err != nil {
		logger.Error("Failed to fetch order", zap.Error(err))
		return err
	}
	if order.UserID != userID {
		logger.Warn("Unauthorized cancellation attempt")
		return ErrUnauthorizedOrder
	}
	if order.Status != models.OrderStatusPending { // Adjust enum as per model
		logger.Warn("Invalid status for cancellation", zap.String("status", string(order.Status)))
		return ErrInvalidOrderStatus
	}

	// Transaction for atomicity
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update order status
		if err := s.orderRepo.UpdateStatus(ctx, orderID, models.OrderStatusCancelled); err != nil {
			return err
		}

		// Unreserve inventory for items (no VariantID, so use ProductID + MerchantID)
		items, err := s.orderItemRepo.FindOrderItemsByOrderID(ctx, orderID)
		if err != nil {
			return err
		}
		for _, item := range items {
			inventory, err := s.inventoryRepo.FindByProductAndMerchant(ctx, item.ProductID, item.MerchantID)
			if err != nil {
				return fmt.Errorf("inventory lookup failed for product %s: %w", item.ProductID, err)
			}
			// Unreserve: Add back to Quantity, subtract from ReservedQuantity
			inventory.Quantity += item.Quantity
			if inventory.ReservedQuantity >= item.Quantity {
				inventory.ReservedQuantity -= item.Quantity
			} else {
				inventory.ReservedQuantity = 0
			}
			if err := s.inventoryRepo.UpdateInventory(ctx, inventory.ID, item.Quantity); err != nil { // Assume repo method for update
				return fmt.Errorf("failed to update inventory %s: %w", inventory.ID, err)
			}
		}

		// Initiate refund if paid
		// if order.Payment != nil && order.Payment.Status == "success" {
		// 	if err := s.paymentService.InitiateRefund(ctx, orderID); err != nil {
		// 		logger.Error("Refund initiation failed", zap.Error(err))
		// 		return ErrRefundFailed
		// 	}
		// }

		return nil
	})
	if err != nil {
		logger.Error("Transaction failed", zap.Error(err))
		return err
	}

	// Notifications (outside tx, fire-and-forget)
	// if err := s.notificationSvc.NotifyUser(ctx, userID, "Order Cancelled", fmt.Sprintf("Order %d cancelled: %s", orderID, reason)); err != nil {
	// 	logger.Warn("User notification failed", zap.Error(err)) // Soft fail
	// }
	// // For multi-vendor: Notify per merchant (stub; loop over items if needed)
	// for _, item := range items { // From earlier fetch
	// 	if err := s.notificationSvc.NotifyMerchant(ctx, item.MerchantID, "Order Item Cancelled", fmt.Sprintf("Item for order %d cancelled", orderID)); err != nil {
	// 		logger.Warn("Merchant notification failed", zap.Error(err))
	// 	}
	// }

	logger.Info("Order cancelled successfully")
	return nil
}


```

---
### services/payment/payment_service.go
- Size: 11.02 KB
- Lines: 349
- Last Modified: 2025-09-30 12:28:15

```go
package payment


import (
	"context"
	"errors"
	"fmt"
	"time"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/config"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"

	"github.com/gray-adeyi/paystack"
	m "github.com/gray-adeyi/paystack/models"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

var (
	ErrPaymentFailed     = errors.New("payment initialization failed")
	ErrVerificationFailed = errors.New("payment verification failed")
	ErrRefundFailed      = errors.New("refund failed")
)


type PaymentService struct {
	paymentRepo repositories.PaymentRepository
	orderRepo   repositories.OrderRepository
	//client      *paystack.Client
	config      *config.Config
	logger      *zap.Logger
}

func NewPaymentService(
	paymentRepo repositories.PaymentRepository,
	orderRepo repositories.OrderRepository,
	conf *config.Config,
	logger *zap.Logger,
) *PaymentService {
	//client := paystack.NewClient(paystack.WithSecretKey(conf.PaystackSecretKey))
	return &PaymentService{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
		//client:      client,
		config:      conf,
		logger:      logger,
	}
}

/*
func (s *PaymentService) InitializeCheckout(ctx context.Context, req dto.InitializePaymentRequest) (*dto.PaymentResponse, error) {
	logger := s.logger.With(zap.Uint("order_id", req.OrderID))

	// Validate and fetch order
	order, err := s.orderRepo.FindByID(ctx,req.OrderID)
	if err != nil {
		logger.Error("Order not found", zap.Error(err))
		return nil, fmt.Errorf("order not found: %w", err)
	}
	amountInKobo := order.SubTotal.Mul(decimal.NewFromInt(100)).String()
	psClient := paystack.NewClient(paystack.WithSecretKey(s.config.PaystackSecretKey))
	// Verify amount matches order total
	expectedTotal := order.SubTotal.InexactFloat64() // Assuming decimal.Decimal
	if req.Amount != expectedTotal {
		logger.Error("Amount mismatch", zap.Float64("expected", expectedTotal), zap.Float64("got", req.Amount))
		return nil, fmt.Errorf("amount mismatch: expected %v, got %v", expectedTotal, req.Amount)
	}

	// Convert amount to kobo for Paystack
	//amountKobo := int(req.Amount * 100)

	// Initialize Paystack transaction
	initReq := &paystack.InitializeTransactionRequest{
		Amount:    amountInKobo,
		Email:     req.Email,
		Currency:  req.Currency,
		Reference: fmt.Sprintf("order_%d", order.ID),
	}
	resp, err := psClient.Transactions.Initialize(context.TODO(),initReq)
	if err != nil {
		logger.Error("Paystack init failed", zap.Error(err))
		return nil, fmt.Errorf("paystack init failed: %w", err)
	}

	// Save payment
	payment := &models.Payment{
		OrderID:       order.ID,
		Amount:        decimal.NewFromFloat(req.Amount),
		Currency:      req.Currency,
		Status:        "pending",
		TransactionID: resp.Data.Reference,
	}
	if err := s.paymentRepo.Create(payment); err != nil {
		logger.Error("Failed to save payment", zap.Error(err))
		return nil, fmt.Errorf("failed to save payment: %w", err)
	}

	// Manual mapping
	response := &dto.PaymentResponse{
		ID:             payment.ID,
		OrderID:        payment.OrderID,
		Amount:         payment.Amount.InexactFloat64(),
		Currency:       payment.Currency,
		Status:         payment.Status,
		TransactionID:  payment.TransactionID,
		AuthorizationURL: resp.Data.AuthorizationURL, // For frontend redirect
		CreatedAt:      payment.CreatedAt,
		UpdatedAt:      payment.UpdatedAt,
	}
	return response, nil
}
*/









func (s *PaymentService) InitializeCheckout(ctx context.Context, req dto.InitializePaymentRequest) (*dto.PaymentResponse, error) {
	logger := s.logger.With(zap.Uint("order_id", req.OrderID))

	// Validate and fetch order
	order, err := s.orderRepo.FindByID(ctx, req.OrderID)
	if err != nil {
		logger.Error("Order not found", zap.Error(err))
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// Compute expected amount in kobo (integer) from order subtotal (decimal.Decimal)
	amountKoboFromOrderInt64 := order.SubTotal.Mul(decimal.NewFromInt(100)).IntPart()
	amountKobo := int(amountKoboFromOrderInt64) // paystack client expects int (kobo)
	// Compute requested amount in kobo and compare as integers to avoid float equality issues
	reqAmountKoboInt64 := decimal.NewFromFloat(req.Amount).Mul(decimal.NewFromInt(100)).IntPart()
	reqAmountKobo := int(reqAmountKoboInt64)

	if reqAmountKobo != amountKobo {
		logger.Error("Amount mismatch", zap.Int("expected_kobo", amountKobo), zap.Int("got_kobo", reqAmountKobo))
		return nil, fmt.Errorf("amount mismatch: expected %d kobo, got %d kobo", amountKobo, reqAmountKobo)
	}

	// Create Paystack client
	psClient := paystack.NewClient(paystack.WithSecretKey(s.config.PaystackSecretKey))

	// Call Transactions.Initialize(amount int, email string, response any, optionalPayloads ...)
	var psResp m.Response[m.InitTransaction]
	// pass currency as optional payload if provided
	var initErr error
	if req.Currency != "" {
		initErr = psClient.Transactions.Initialize(ctx, amountKobo, req.Email, &psResp, paystack.WithOptionalPayload("currency", req.Currency))
	} else {
		initErr = psClient.Transactions.Initialize(ctx, amountKobo, req.Email, &psResp)
	}
	if initErr != nil {
		logger.Error("Paystack initialize transaction failed", zap.Error(initErr))
		return nil, fmt.Errorf("paystack initialize failed: %w", initErr)
	}

	// Ensure response data exists
	if psResp.Data.Reference == "" {
		logger.Error("Paystack initialize returned empty reference", zap.Any("response", psResp))
		return nil, fmt.Errorf("paystack initialize returned empty reference")
	}

	// Save payment. Use order.SubTotal (the canonical amount) to avoid any float round issues.
	payment := &models.Payment{
		OrderID:       order.ID,
		Amount:        order.SubTotal, // decimal.Decimal
		Currency:      req.Currency,
		Status:        "pending",
		TransactionID: psResp.Data.Reference,
	}
	if err := s.paymentRepo.Create(ctx,payment); err != nil {
		logger.Error("Failed to save payment", zap.Error(err))
		return nil, fmt.Errorf("failed to save payment: %w", err)
	}

	// Map to DTO response
	response := &dto.PaymentResponse{
		ID:               payment.ID,
		OrderID:          payment.OrderID,
		Amount:           payment.Amount.InexactFloat64(),
		Currency:         payment.Currency,
		Status:          string(payment.Status),
		TransactionID:    payment.TransactionID,
		AuthorizationURL: psResp.Data.AuthorizationUrl, // used by frontend for redirect
		CreatedAt:        payment.CreatedAt,
		UpdatedAt:        payment.UpdatedAt,
	}

	return response, nil
}








/*
func (s *PaymentService) VerifyPayment(ctx context.Context, reference string) (*dto.PaymentResponse, error) {
	logger := s.logger.With(zap.String("reference", reference))
	


	// Verify with Paystack
	psClient := paystack.NewClient(paystack.WithSecretKey(s.config.PaystackSecretKey))
	var psResp m.Response[m.Transaction]
	 err := psClient.Transactions.Verify(ctx, reference, &psResp)
	if err != nil {
		logger.Error("Paystack verify failed", zap.Error(err))
		return nil, fmt.Errorf("paystack verify failed: %w", err)
	}
	if psResp.Data.Status != "success" {
		logger.Warn("Payment not successful", zap.String("status", string(psResp.Data.Status)))
		return nil, fmt.Errorf("payment not successful: status %s",string(psResp.Data.Status))
	}

	// Fetch and update payment
	payment, err := s.paymentRepo.FindByTransactionID(ctx, reference)
	if err != nil {
		logger.Error("Payment not found", zap.Error(err))
		return nil, fmt.Errorf("payment not found: %w", err)
	}
	payment.Status = "success"
	if err := s.paymentRepo.Update(payment); err != nil {
		logger.Error("Failed to update payment", zap.Error(err))
		return nil, fmt.Errorf("failed to update payment: %w", err)
	}

	// Update order status
	order, err := s.orderRepo.FindByID(ctx ,payment.OrderID)
	if err != nil {
		logger.Error("Order not found", zap.Error(err))
		return nil, fmt.Errorf("order not found: %w", err)
	}
	order.Status = "paid"
	if err := s.orderRepo.Update(order); err != nil {
		logger.Error("Failed to update order", zap.Error(err))
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	// Manual mapping
	response := &dto.PaymentResponse{
		ID:            payment.ID,
		OrderID:       payment.OrderID,
		Amount:        payment.Amount.InexactFloat64(),
		Currency:      payment.Currency,
		Status:       string(payment.Status),
		TransactionID: payment.TransactionID,
		CreatedAt:     payment.CreatedAt,
		UpdatedAt:     payment.UpdatedAt,
	}
	return response, nil
}
*/


func (s *PaymentService) VerifyPayment(ctx context.Context, reference string) (*dto.PaymentResponse, error) {
	logger := s.logger.With(zap.String("operation", "VerifyPayment"), zap.String("reference", reference))

	payment, perr := s.paymentRepo.FindByTransactionID(ctx, reference)
	if perr != nil {
		return nil, fmt.Errorf("payment not found: %w", perr)
	}

	psClient := paystack.NewClient(paystack.WithSecretKey(s.config.PaystackSecretKey))
	var resp m.Response[m.Transaction]
	 err := psClient.Transactions.Verify(ctx, reference, &resp)
	if err != nil || !resp.Status  || resp.Data.Status != "success" {
		logger.Error("Paystack verification failed", zap.Error(err))
		// Update to failed
		payment.Status = models.PaymentStatusFailed
		s.paymentRepo.Update(ctx, payment)
		return nil, ErrVerificationFailed
	}

	// Update success
	payment.Status = models.PaymentStatusCompleted
	payment.UpdatedAt = time.Now()
	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return nil, err
	}

	// Update order status
	order, err := s.orderRepo.FindByID(ctx, payment.OrderID)
	if err == nil {
		order.Status = models.OrderStatusCompleted
		s.orderRepo.Update(ctx, order)
	}

	logger.Info("Payment verified", zap.Uint("payment_id", payment.ID))
	response := &dto.PaymentResponse{
		ID:            payment.ID,
		OrderID:       payment.OrderID,
		Amount:        payment.Amount.InexactFloat64(),
		Currency:      payment.Currency,
		Status:       string(payment.Status),
		TransactionID: payment.TransactionID,
		CreatedAt:     payment.CreatedAt,
		UpdatedAt:     payment.UpdatedAt,
	}
	return response, nil
}



// GetPaymentByOrderID retrieves a payment by order ID
func (s *PaymentService) GetPaymentByOrderID(ctx context.Context, orderID uint) (*models.Payment, error) {
	if orderID == 0 {
		return nil, errors.New("invalid order ID")
	}
	return s.paymentRepo.FindByOrderID(ctx,orderID)
}

// GetPaymentsByUserID retrieves all payments for a user
func (s *PaymentService) GetPaymentsByUserID(ctx context.Context,userID uint) ([]models.Payment, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}
	return s.paymentRepo.FindByUserID(ctx,userID)
}

// UpdatePaymentStatus updates the status of a payment
func (s *PaymentService) UpdatePaymentStatus(ctx context.Context,paymentID uint, status string) (*models.Payment, error) {
	if paymentID == 0 {
		return nil, errors.New("invalid payment ID")
	}
	if err := models.PaymentStatus(status).Valid(); err != nil {
		return nil, err
	}

	payment, err := s.paymentRepo.FindByID(ctx ,paymentID)
	if err != nil {
		return nil, err
	}

	payment.Status = models.PaymentStatus(status)
	if err := s.paymentRepo.Update(ctx ,payment); err != nil {
		return nil, err
	}

	return s.paymentRepo.FindByID(ctx ,paymentID)
}
```

---
### services/payout/payout_service.go
- Size: 2.58 KB
- Lines: 91
- Last Modified: 2025-09-30 12:28:15

```go
package payout
/*
import (
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"
	"context"
	"errors"

	"github.com/shopspring/decimal"
)

type PayoutService struct {
	payoutRepo *repositories.PayoutRepository
}

func NewPayoutService(payoutRepo *repositories.PayoutRepository) *PayoutService {
	return &PayoutService{
		payoutRepo: payoutRepo,
	}
}

// CreatePayout creates a payout for a merchant
func (s *PayoutService) CreatePayout(merchantID uint, amount float64) (*models.Payout, error) {
	if merchantID == 0 {
		return nil, errors.New("invalid merchant ID")
	}
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	// Simulate payout processing (placeholder for Stripe)
	payout := &models.Payout{
		MerchantID: merchantID,
		Amount:     amount,
		Status:     models.PayoutStatusPending,
	}
	if err := s.payoutRepo.Create(payout); err != nil {
		return nil, err
	}

	// Simulate successful payout
	payout.Status = models.PayoutStatusCompleted
	if err := s.payoutRepo.Update(payout); err != nil {
		return nil, err
	}

	return s.payoutRepo.FindByID(payout.ID)
}

// GetPayoutByID retrieves a payout by ID
func (s *PayoutService) GetPayoutByID(id uint) (*models.Payout, error) {
	if id == 0 {
		return nil, errors.New("invalid payout ID")
	}
	return s.payoutRepo.FindByID(id)
}

// GetPayoutsByMerchantID retrieves all payouts for a merchant
func (s *PayoutService) GetPayoutsByMerchantID(merchantID uint) ([]models.Payout, error) {
	if merchantID == 0 {
		return nil, errors.New("invalid merchant ID")
	}
	return s.payoutRepo.FindByMerchantID(merchantID)
}


func (s *PayoutService) RequestPayout(ctx context.Context, merchantID string) (*models.Payout, error) {
    // Calc eligible: sum splits where status=pending AND hold_until < now
    var totalDue decimal.Decimal
    db.DB.Model(&models.OrderMerchantSplit{}).
        Where("merchant_id = ? AND status = 'pending' AND hold_until < ?", merchantID, time.Now()).
        Select("SUM(amount_due)").Scan(&totalDue)
    if totalDue.LessThanOrEqual(decimal.Zero) {
        return nil, errors.New("no eligible balance")
    }

    payout := &models.Payout{
        MerchantID: merchantID,
        Amount:     totalDue,
        Status:     "pending",  // Admin approves/sends
    }
    if err := db.DB.Create(payout).Error; err != nil {
        return nil, err
    }
    // Update splits to 'payout_requested'
    db.DB.Model(&models.OrderMerchantSplit{}).
        Where("merchant_id = ? AND status = 'pending' AND hold_until < ?", merchantID, time.Now()).
        Update("status", "payout_requested")
    return payout, nil
}
*/
```

---
### services/pricing/pricing_service.go
- Size: 1.04 KB
- Lines: 42
- Last Modified: 2025-09-30 12:22:22

```go
package pricing

/*
import (
    "errors"
    "api-customer-merchant/internal/db/models"
)

type PricingService struct {
    // Repos if needed
}

func NewPricingService() *PricingService {
    return &PricingService{}
}

func (s *PricingService) CalculateShipping(cart *models.Cart, address string) (float64, error) {
    // Integrate with shipping API (e.g., UPS); placeholder
    if address == "" {
        return 0, errors.New("address required")
    }
    return 10.00, nil // Flat rate per vendor count
}

func (s *PricingService) CalculateTax(cart *models.Cart, country string) (float64, error) {
    // Tax API or rules; placeholder
    var subtotal float64
    for _, item := range cart.CartItems {
        subtotal += float64(item.Quantity) * item.Product.Price
    }
    rate := 0.1 // 10% VAT for luxury
    return subtotal * rate, nil
}

func (s *PricingService) ApplyPromotion(cart *models.Cart, code string) (float64, error) {
    // Validate coupon; placeholder
    if code == "" {
        return 0, nil
    }
    return 5.00, nil // Discount
}
*/

```

---
### services/product/proudct_service.go
- Size: 22.13 KB
- Lines: 686
- Last Modified: 2025-09-30 12:28:15

```go
package product

import (
	"context"
	"errors"
	"fmt"
	//"os"
	"path/filepath"

	//"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/config"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var (
	ErrInvalidProduct    = errors.New("invalid product data")
	//ErrInvalidSKU        = errors.New("invalid SKU format")
	ErrInvalidMediaURL   = errors.New("invalid media URL")
	ErrInvalidAttributes = errors.New("invalid variant attributes")
	ErrUnauthorized      = errors.New("unauthorized operation")
	ErrUploadFailed     = errors.New("upload to Cloudinary failed")
	ErrUpdateFailed     = errors.New("update failed")
	ErrDeleteFailed     = errors.New("delete failed")
	ErrUnauthorizedMedia = errors.New("unauthorized for this media")
)

// SKU validation regex: alphanumeric, hyphens, underscores, max 100 chars
//var skuRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,100}$`)

type ProductService struct {
	productRepo *repositories.ProductRepository
	logger      *zap.Logger
	validator   *validator.Validate
	cld         *cloudinary.Cloudinary
	config  *config.Config
	
}

func NewProductService(productRepo *repositories.ProductRepository,  cfg *config.Config,logger *zap.Logger) *ProductService {
	cld, err := cloudinary.NewFromParams(cfg.CloudinaryCloudName, cfg.CloudinaryAPIKey, cfg.CloudinaryAPISecret)
	if err != nil {
		logger.Fatal("Cloudinary init failed", zap.Error(err))
	}

	return &ProductService{
		productRepo: productRepo,
		logger:      logger,
		validator:   validator.New(),
		cld:         cld,
	}
}

// CreateProductWithVariants creates a product from input DTO
func (s *ProductService) CreateProductWithVariants(ctx context.Context, input *dto.ProductInput) (*dto.ProductResponse, error) {
	logger := s.logger.With(zap.String("operation", "CreateProductWithVariants"))

	// Validate input
	if err := s.validator.Struct(input); err != nil {
		logger.Error("Input validation failed", zap.Error(err))
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	// Additional validation
	// if !skuRegex.MatchString(input.SKU) {
	// 	logger.Error("Invalid SKU format", zap.String("sku", input.SKU))
	// 	return nil, ErrInvalidSKU
	// }

	isSimple := len(input.Variants) == 0
	if isSimple && input.InitialStock == nil {
		logger.Error("Initial stock required for simple product")
		return nil, ErrInvalidProduct
	}
	

	// Map DTO to models
	product := &models.Product{
		Name:        strings.TrimSpace(input.Name),
		MerchantID:  strings.TrimSpace(input.MerchantID),
		Description: strings.TrimSpace(input.Description),
		//SKU:         strings.TrimSpace(input.SKU),
		BasePrice:   decimal.NewFromFloat(input.BasePrice),
		CategoryID:  input.CategoryID,
	}
	variants := make([]models.Variant, len(input.Variants))
	for i, v := range input.Variants {
		variants[i] = models.Variant{
			//SKU:             strings.TrimSpace(v.SKU),
			PriceAdjustment: decimal.NewFromFloat(v.PriceAdjustment),
			Attributes:      v.Attributes,
			IsActive:        true,
		}
	}
	media := make([]models.Media, len(input.Media))
	for i, m := range input.Media {
		media[i] = models.Media{
			URL:  strings.TrimSpace(m.URL),
			Type: models.MediaType(m.Type),
		}
	}


	product.GenerateSKU(input.MerchantID)
	for i := range variants {
		variants[i].GenerateSKU(product.SKU)
	}

	// Delegate to repo
	var simpleStock *int
if isSimple {
    simpleStock = input.InitialStock
}
	err := s.productRepo.CreateProductWithVariantsAndInventory(ctx, product, variants, input.Variants, media, simpleStock, isSimple)
	if err != nil {
		if errors.Is(err, repositories.ErrDuplicateSKU) {
			return nil, fmt.Errorf("duplicate SKU: %w", err)
		}
		if errors.Is(err, repositories.ErrInvalidInventory) {
			return nil, fmt.Errorf("invalid inventory setup: %w", err)
		}
		logger.Error("Failed to create product", zap.Error(err))
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// Map to response DTO
	response := &dto.ProductResponse{
		ID:          product.ID,
		MerchantID:  product.MerchantID,
		Name:        product.Name,
		Description: product.Description,
		//SKU:         product.SKU,
		BasePrice:   (product.BasePrice).InexactFloat64(),
		CategoryID:  product.CategoryID,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
		Media:       make([]dto.MediaResponse, len(product.Media)),
		Variants:    make([]dto.VariantResponse, len(product.Variants)),
	}
	for i, v := range product.Variants {
		response.Variants[i] = dto.VariantResponse{
			ID:              v.ID,
			ProductID:       v.ProductID,
			//SKU:             v.SKU,
			PriceAdjustment: v.PriceAdjustment.InexactFloat64(),
			TotalPrice:      v.TotalPrice.InexactFloat64(),
			Attributes:      v.Attributes,
			IsActive:        v.IsActive,
			CreatedAt:       v.CreatedAt,
			UpdatedAt:       v.UpdatedAt,
			Inventory: dto.InventoryResponse{
				ID:                v.Inventory.ID,
				Quantity:          v.Inventory.Quantity,
				ReservedQuantity:  v.Inventory.ReservedQuantity,
				LowStockThreshold: v.Inventory.LowStockThreshold,
				BackorderAllowed:  v.Inventory.BackorderAllowed,
			},
		}
	}

	// Map media
	for i, m := range product.Media {
		response.Media[i] = dto.MediaResponse{
			ID:        m.ID,
			ProductID: m.ProductID,
			URL:       m.URL,
			Type:      string(m.Type),
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		}
	}

	// SimpleInventory is always nil for simple products
	//response.SimpleInventory = nil
	if product.SimpleInventory != nil {
    response.SimpleInventory = &dto.InventoryResponse{
        ID:                product.SimpleInventory.ID,
        Quantity:          product.SimpleInventory.Quantity,
        ReservedQuantity:  product.SimpleInventory.ReservedQuantity,
        LowStockThreshold: product.SimpleInventory.LowStockThreshold,
        BackorderAllowed:  product.SimpleInventory.BackorderAllowed,
    }
}

	logger.Info("Product created successfully", zap.String("product_id", product.ID))
	return response, nil
}

// GetProductByID fetches a product with optional preloads
func (s *ProductService) GetProductByID(ctx context.Context, id string, preloads ...string) (*dto.ProductResponse, error) {
	logger := s.logger.With(zap.String("operation", "GetProductByID"), zap.String("product_id", id))
	product, err := s.productRepo.FindByID(ctx, id, preloads...)  // Fixed: Added ctx
	if err != nil {
		if errors.Is(err, repositories.ErrProductNotFound) {
			return nil, err
		}
		logger.Error("Failed to fetch product", zap.Error(err))
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}

	response := &dto.ProductResponse{
		ID:          product.ID,
		MerchantID:  product.MerchantID,
		Name:        product.Name,
		Description: product.Description,
		SKU:         product.SKU,
		BasePrice:   (product.BasePrice).InexactFloat64(),
		CategoryID:  product.CategoryID,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
		Variants:    make([]dto.VariantResponse, len(product.Variants)),
		Media:       make([]dto.MediaResponse, len(product.Media)),
	}
	for i, v := range product.Variants {
		response.Variants[i] = dto.VariantResponse{
			ID:              v.ID,
			ProductID:       v.ProductID,
			SKU:             v.SKU,
			PriceAdjustment: (v.PriceAdjustment).InexactFloat64(),
			TotalPrice:      (v.TotalPrice).InexactFloat64(),
			Attributes:      v.Attributes,
			IsActive:        v.IsActive,
			CreatedAt:       v.CreatedAt,
			UpdatedAt:       v.UpdatedAt,
			Inventory: dto.InventoryResponse{
				ID:                v.Inventory.ID,
				Quantity:          v.Inventory.Quantity,
				ReservedQuantity:  v.Inventory.ReservedQuantity,
				LowStockThreshold: v.Inventory.LowStockThreshold,
				BackorderAllowed:  v.Inventory.BackorderAllowed,
			},
		}
	}
	for i, m := range product.Media {
		response.Media[i] = dto.MediaResponse{
			ID:        m.ID,
			ProductID: m.ProductID,
			URL:       m.URL,
			Type:      string(m.Type),
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		}
	}
	if product.SimpleInventory != nil {
		response.SimpleInventory = &dto.InventoryResponse{
			ID:                product.SimpleInventory.ID,
			Quantity:          product.SimpleInventory.Quantity,
			ReservedQuantity:  product.SimpleInventory.ReservedQuantity,
			LowStockThreshold: product.SimpleInventory.LowStockThreshold,
			BackorderAllowed:  product.SimpleInventory.BackorderAllowed,
		}
	}

	logger.Info("Product fetched successfully")
	return response, nil
}

// ListProductsByMerchant lists products for a merchant
func (s *ProductService) ListProductsByMerchant(ctx context.Context, merchantID string, limit, offset int, activeOnly bool) ([]dto.ProductResponse, error) {
	logger := s.logger.With(zap.String("operation", "ListProductsByMerchant"), zap.String("merchant_id", merchantID))
	products, err := s.productRepo.ListByMerchant(ctx, merchantID, limit, offset, activeOnly)  // Fixed: Added ctx
	if err != nil {
		logger.Error("Failed to list products", zap.Error(err))
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	responses := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		responses[i] = dto.ProductResponse{
			ID:          p.ID,
			MerchantID:  p.MerchantID,
			Name:        p.Name,
			Description: p.Description,
			SKU:         p.SKU,
			BasePrice:   (p.BasePrice).InexactFloat64(),
			CategoryID:  p.CategoryID,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
			Variants:    make([]dto.VariantResponse, len(p.Variants)),
			Media:       make([]dto.MediaResponse, len(p.Media)),
		}
		for j, v := range p.Variants {
			responses[i].Variants[j] = dto.VariantResponse{
				ID:              v.ID,
				ProductID:       v.ProductID,
				SKU:             v.SKU,
				PriceAdjustment: (v.PriceAdjustment).InexactFloat64(),
				TotalPrice:      (v.TotalPrice).InexactFloat64(),
				Attributes:      v.Attributes,
				IsActive:        v.IsActive,
				CreatedAt:       v.CreatedAt,
				UpdatedAt:       v.UpdatedAt,
				Inventory: dto.InventoryResponse{
					ID:                v.Inventory.ID,
					Quantity:          v.Inventory.Quantity,
					ReservedQuantity:  v.Inventory.ReservedQuantity,
					LowStockThreshold: v.Inventory.LowStockThreshold,
					BackorderAllowed:  v.Inventory.BackorderAllowed,
				},
			}
		}
		for j, m := range p.Media {
			responses[i].Media[j] = dto.MediaResponse{
				ID:        m.ID,
				ProductID: m.ProductID,
				URL:       m.URL,
				Type:      string(m.Type),
				CreatedAt: m.CreatedAt,
				UpdatedAt: m.UpdatedAt,
			}
		}
		if p.SimpleInventory != nil {
			responses[i].SimpleInventory = &dto.InventoryResponse{
				ID:                p.SimpleInventory.ID,
				Quantity:          p.SimpleInventory.Quantity,
				ReservedQuantity:  p.SimpleInventory.ReservedQuantity,
				LowStockThreshold: p.SimpleInventory.LowStockThreshold,
				BackorderAllowed:  p.SimpleInventory.BackorderAllowed,
			}
		}
	}

	logger.Info("Products listed successfully", zap.Int("count", len(responses)))
	return responses, nil
}

// GetAllProducts fetches all active products for the landing page
func (s *ProductService) GetAllProducts(ctx context.Context, limit, offset int, categoryID *uint) ([]dto.ProductResponse, int64, error) {
	logger := s.logger.With(zap.String("operation", "GetAllProducts"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	products, total, err := s.productRepo.GetAllProducts(ctx, limit, offset, categoryID, "Media", "Variants", "Variants.Inventory", "SimpleInventory")  // Fixed: Added ctx (resolves type shifts)
	if err != nil {
		logger.Error("Failed to fetch all products", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to fetch products: %w", err)
	}

	responses := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		responses[i] = dto.ProductResponse{
			ID:          p.ID,
			MerchantID:  "", // Exclude for customer-facing API
			Name:        p.Name,
			Description: p.Description,
			SKU:         p.SKU,
			BasePrice:   (p.BasePrice).InexactFloat64(),
			CategoryID:  p.CategoryID,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
			Variants:    make([]dto.VariantResponse, len(p.Variants)),
			Media:       make([]dto.MediaResponse, len(p.Media)),
		}
		for j, v := range p.Variants {
			responses[i].Variants[j] = dto.VariantResponse{
				ID:              v.ID,
				ProductID:       v.ProductID,
				SKU:             v.SKU,
				PriceAdjustment: (v.PriceAdjustment).InexactFloat64(),
				TotalPrice:      (v.TotalPrice).InexactFloat64(),
				Attributes:      v.Attributes,
				IsActive:        v.IsActive,
				CreatedAt:       v.CreatedAt,
				UpdatedAt:       v.UpdatedAt,
				Inventory: dto.InventoryResponse{
					ID:                v.Inventory.ID,
					Quantity:          v.Inventory.Quantity,
					ReservedQuantity:  v.Inventory.ReservedQuantity,
					LowStockThreshold: v.Inventory.LowStockThreshold,
					BackorderAllowed:  v.Inventory.BackorderAllowed,
				},
			}
		}
		for j, m := range p.Media {
			responses[i].Media[j] = dto.MediaResponse{
				ID:        m.ID,
				ProductID: m.ProductID,
				URL:       m.URL,
				Type:      string(m.Type),
				CreatedAt: m.CreatedAt,
				UpdatedAt: m.UpdatedAt,
			}
		}
		if p.SimpleInventory != nil {
			responses[i].SimpleInventory = &dto.InventoryResponse{
				ID:                p.SimpleInventory.ID,
				Quantity:          p.SimpleInventory.Quantity,
				ReservedQuantity:  p.SimpleInventory.ReservedQuantity,
				LowStockThreshold: p.SimpleInventory.LowStockThreshold,
				BackorderAllowed:  p.SimpleInventory.BackorderAllowed,
			}
		}
	}

	logger.Info("Products fetched for landing page", zap.Int("count", len(responses)), zap.Int64("total", total))
	return responses, total, nil
}





// GetAllProducts fetches all active products for the landing page
// Assumes ProductFilter is defined in the same package or imported.
 type ProductFilter struct {
     CategoryName   *string
     CategoryID     *uint
     MinPrice       *decimal.Decimal
     MaxPrice       *decimal.Decimal
     InStock        *bool
     VariantAttrs   map[string]interface{}
     MerchantName   *string
 }

func (s *ProductService) FilterProducts(ctx context.Context, filter ProductFilter, limit, offset int) ([]dto.ProductResponse, int64, error) {
	logger := s.logger.With(zap.String("operation", "FilterProducts"))

	// --- pagination sanitization ---
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}




	// --- fetch products from repository using the provided filter ---
	repoFilter := repositories.ProductFilter{
    CategoryName: filter.CategoryName,
    CategoryID:   filter.CategoryID,
    MinPrice:     filter.MinPrice,
    MaxPrice:     filter.MaxPrice,
    InStock:      filter.InStock,
    VariantAttrs: filter.VariantAttrs,
    MerchantName: filter.MerchantName,
}

products, total, err := s.productRepo.ProductsFilter(ctx, repoFilter, limit, offset, "Media", "Variants", "Variants.Inventory", "SimpleInventory")

	if err != nil {
		logger.Error("Failed to fetch products", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to fetch products: %w", err)
	}

	// --- map DB models -> DTOs ---
	responses := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		// convert base price once
		basePriceFloat := p.BasePrice.InexactFloat64()

		responses[i] = dto.ProductResponse{
			ID:          p.ID,
			MerchantID:  "", // hide merchant id from customer-facing API
			Name:        p.Name,
			Description: p.Description,
			SKU:         p.SKU,
			BasePrice:   basePriceFloat,
			CategoryID:  p.CategoryID,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
			Variants:    make([]dto.VariantResponse, len(p.Variants)),
			Media:       make([]dto.MediaResponse, len(p.Media)),
		}

		// Variants: compute TotalPrice = BasePrice + PriceAdjustment
		for j, v := range p.Variants {
			totalPriceDecimal := p.BasePrice.Add(v.PriceAdjustment) // decimal + decimal
			responses[i].Variants[j] = dto.VariantResponse{
				ID:              v.ID,
				ProductID:       v.ProductID,
				SKU:             v.SKU,
				PriceAdjustment: v.PriceAdjustment.InexactFloat64(),
				TotalPrice:      totalPriceDecimal.InexactFloat64(),
				Attributes:      v.Attributes,
				IsActive:        v.IsActive,
				CreatedAt:       v.CreatedAt,
				UpdatedAt:       v.UpdatedAt,
				Inventory: dto.InventoryResponse{
					ID:                v.Inventory.ID,
					Quantity:          v.Inventory.Quantity,
					ReservedQuantity:  v.Inventory.ReservedQuantity,
					LowStockThreshold: v.Inventory.LowStockThreshold,
					BackorderAllowed:  v.Inventory.BackorderAllowed,
				},
			}
		}

		// Media
		for j, m := range p.Media {
			responses[i].Media[j] = dto.MediaResponse{
				ID:        m.ID,
				ProductID: m.ProductID,
				URL:       m.URL,
				Type:      string(m.Type),
				CreatedAt: m.CreatedAt,
				UpdatedAt: m.UpdatedAt,
			}
		}

		// SimpleInventory (for non-variant products)
		if p.SimpleInventory != nil {
			responses[i].SimpleInventory = &dto.InventoryResponse{
				ID:                p.SimpleInventory.ID,
				Quantity:          p.SimpleInventory.Quantity,
				ReservedQuantity:  p.SimpleInventory.ReservedQuantity,
				LowStockThreshold: p.SimpleInventory.LowStockThreshold,
				BackorderAllowed:  p.SimpleInventory.BackorderAllowed,
			}
		}
	}

	logger.Info("Products fetched for filter", zap.Int("count", len(responses)), zap.Int64("total", total))
	return responses, total, nil
}







// UpdateInventory adjusts stock for a given inventory ID
func (s *ProductService) UpdateInventory(ctx context.Context, inventoryID string, delta int) error {
	logger := s.logger.With(zap.String("operation", "UpdateInventory"), zap.String("inventory_id", inventoryID))
	err := s.productRepo.UpdateInventoryQuantity(inventoryID, delta)
	if err != nil {
		logger.Error("Failed to update inventory", zap.Error(err))
		return fmt.Errorf("failed to update inventory: %w", err)
	}
	logger.Info("Inventory updated successfully", zap.Int("delta", delta))
	return nil
}

// DeleteProduct soft-deletes a product
func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
	logger := s.logger.With(zap.String("operation", "DeleteProduct"), zap.String("product_id", id))
	err := s.productRepo.SoftDeleteProduct(id)
	if err != nil {
		logger.Error("Failed to delete product", zap.Error(err))
		return fmt.Errorf("failed to delete product: %w", err)
	}
	logger.Info("Product deleted successfully")
	return nil
}







//Media service


// UploadMedia uploads file to Cloudinary, saves to DB
func (s *ProductService) UploadMedia(ctx context.Context, productID, merchantID, filePath, mediaType string) (*models.Media, error) {
	logger := s.logger.With(zap.String("operation", "UploadMedia"), zap.String("product_id", productID))

	// Validate merchant owns product
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil || product.MerchantID != merchantID {
		return nil, ErrUnauthorizedMedia
	}

	// Upload to Cloudinary
	params := uploader.UploadParams{
		Folder:     "merchant_media", // Organized folder
		ResourceType: mediaType, // image/video
		PublicID:    fmt.Sprintf("%s_%s", productID, filepath.Base(filePath)), // Unique ID
	}
	resp, err := s.cld.Upload.Upload(ctx, filePath, params)
	if err != nil {
		logger.Error("Cloudinary upload failed", zap.Error(err))
		return nil, ErrUploadFailed
	}

	// Save to DB
	//mediaType models.Media
	media := &models.Media{
		ProductID: productID,
		URL:       resp.SecureURL,
		Type:      models.MediaType(mediaType),
		PublicID:  resp.PublicID, // Store for delete/update (add to model if missing)
	}
	if err := s.productRepo.CreateMedia(ctx, media); err != nil {
		// Cleanup on failure
		s.cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: media.PublicID})
		return nil, err
	}

	logger.Info("Media uploaded", zap.String("public_id", resp.PublicID))
	return media, nil
}

// UpdateMedia re-uploads or updates URL
func (s *ProductService) UpdateMedia(ctx context.Context, mediaID, productID, merchantID string, req *dto.MediaUpdateRequest) (*models.Media, error) {
	logger := s.logger.With(zap.String("operation", "UpdateMedia"), zap.String("media_id", mediaID))

	media, err := s.productRepo.FindMediaByID(ctx, mediaID)
	if err != nil || media.ProductID != productID || !s.merchantOwnsProduct(ctx, productID, merchantID) {
		return nil, ErrUnauthorizedMedia
	}

	var newURL string
	var newPublicID string
	if req.File != nil {
		// Re-upload
		resp, err := s.cld.Upload.Upload(ctx, *req.File, uploader.UploadParams{
			PublicID:    media.PublicID, // Overwrite existing
			ResourceType: string(media.Type),
		})
		if err != nil {
			logger.Error("Cloudinary re-upload failed", zap.Error(err))
			return nil, ErrUpdateFailed
		}
		newURL = resp.SecureURL
		newPublicID = resp.PublicID
	} else if req.URL != nil {
		newURL = *req.URL
	}

	// Update DB
	updates := map[string]interface{}{"url": newURL}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if newPublicID != "" {
		updates["public_id"] = newPublicID
	}
	if err := s.productRepo.UpdateMedia(ctx, mediaID, updates); err != nil {
		return nil, err
	}

	media.URL = newURL
	if req.Type != nil {
		media.Type = models.MediaType(*req.Type)
	}
	return media, nil
}

// DeleteMedia destroys on Cloudinary, deletes from DB
func (s *ProductService) DeleteMedia(ctx context.Context, mediaID, productID, merchantID, reason string) error {
	logger := s.logger.With(zap.String("operation", "DeleteMedia"), zap.String("media_id", mediaID))

	media, err := s.productRepo.FindMediaByID(ctx, mediaID)
	if err != nil || media.ProductID != productID || !s.merchantOwnsProduct(ctx, productID, merchantID) {
		return ErrUnauthorizedMedia
	}

	// Destroy on Cloudinary
	_, err = s.cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: media.PublicID})
	if err != nil {
		logger.Error("Cloudinary destroy failed", zap.Error(err))
		return ErrDeleteFailed
	}

	// Soft delete from DB
	if err := s.productRepo.DeleteMedia(ctx, mediaID); err != nil {
		return err
	}

	logger.Info("Media deleted", zap.String("public_id", media.PublicID), zap.String("reason", reason))
	return nil
}

// Helper: Check merchant owns product
func (s *ProductService) merchantOwnsProduct(ctx context.Context, productID, merchantID string) bool {
	product, err := s.productRepo.FindByID(ctx, productID)
	return err == nil && product.MerchantID == merchantID
}

```

---
### services/return_request/return_request_service.go
- Size: 2.37 KB
- Lines: 89
- Last Modified: 2025-10-02 15:19:40

```go
package return_request

import (
	"context"
	"errors"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"

	"github.com/google/uuid"
)


var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidInput = errors.New("invalid input")
)


type ReturnRequestService struct {
	repo *repositories.ReturnRequestRepository
}

func NewReturnRequestService(repo *repositories.ReturnRequestRepository) *ReturnRequestService {
	return &ReturnRequestService{repo: repo}
}

func (s *ReturnRequestService) CreateReturnRequest(ctx context.Context, userID uint, req dto.CreateReturnRequestDTO) (*dto.ReturnRequestResponseDTO, error) {
	returnReq := &models.ReturnRequest{
		ID:          uuid.NewString(),
		OrderItemID: req.OrderItemID,
		CustomerID:  userID,
		Reason:      req.Reason,
		Status:      "Pending",
	}

	if err := s.repo.Create(ctx, returnReq); err != nil {
		return nil, err
	}

	return &dto.ReturnRequestResponseDTO{
		ID:          returnReq.ID,
		OrderItemID: returnReq.OrderItemID,
		CustomerID:  returnReq.CustomerID,
		Reason:      returnReq.Reason,
		Status:      returnReq.Status,
		CreatedAt:   returnReq.CreatedAt,
		UpdatedAt:   returnReq.UpdatedAt,
	}, nil
}

func (s *ReturnRequestService) GetReturnRequest(ctx context.Context, id string, userID uint) (*dto.ReturnRequestResponseDTO, error) {
    returnReq, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    if returnReq.CustomerID != userID {
        return nil, ErrUnauthorized
    }

    return mapReturnRequestToDTO(returnReq), nil
}

func (s *ReturnRequestService) GetCustomerReturnRequests(ctx context.Context, userID uint) ([]dto.ReturnRequestResponseDTO, error) {
    returnRequests, err := s.repo.FindByCustomerID(ctx, userID)
    if err != nil {
        return nil, err
    }

    dtos := make([]dto.ReturnRequestResponseDTO, len(returnRequests))
    for i, req := range returnRequests {
        dtos[i] = *mapReturnRequestToDTO(&req)
    }
    return dtos, nil
}

func mapReturnRequestToDTO(r *models.ReturnRequest) *dto.ReturnRequestResponseDTO {
    return &dto.ReturnRequestResponseDTO{
        ID:          r.ID,
        OrderItemID: r.OrderItemID,
        CustomerID:  r.CustomerID,
        Reason:      r.Reason,
        Status:      r.Status,
        CreatedAt:   r.CreatedAt,
        UpdatedAt:   r.UpdatedAt,
    }
}
```

---
### services/user/user_service.go
- Size: 5.99 KB
- Lines: 234
- Last Modified: 2025-09-30 12:22:22

```go
package user

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	//"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"

	//"google.golang.org/api/oauth2/v2"

	//"api-customer-merchant/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

type googleUserInfo struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (s *AuthService) RegisterUser(email, name, password, country string) (*models.User, error) {
	_, err := s.userRepo.FindByEmail(email)
	if err == nil {
		return nil, errors.New("email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:    email,
		Name:     name,
		Password: string(hashedPassword),
		Country:  country,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) LoginUser(email, password string) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *AuthService) GenerateJWT(entity interface{}) (string, error) {
	var id uint
	var entityType string

	switch e := entity.(type) {
	case *models.User:
		id = e.ID
		entityType = "user"

	default:
		return "", errors.New("invalid entity type")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         float64(id),
		"entityType": entityType,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET not set")
	}

	return token.SignedString([]byte(secret))
}

func (s *AuthService) GetOAuthConfig(entityType string) *oauth2.Config {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		log.Println("BASE_URL environment variable is not set")
		return nil
	}
	redirectURL := baseURL + "/customer/auth/google/callback"
	log.Printf("OAuth redirect URL: %s", redirectURL)
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid"},
		Endpoint: google.Endpoint,
	}
}

func (s *AuthService) GoogleLogin(code, baseURL, entityType string) (*models.User, string, error) {
	if entityType != "customer" {
		log.Printf("Invalid entityType for OAuth: %s", entityType)
		return nil, "", errors.New("OAuth only supported for customers")
	}

	// Get OAuth config
	oauthConfig := s.GetOAuthConfig(entityType)
	if oauthConfig == nil || oauthConfig.ClientID == "" || oauthConfig.ClientSecret == "" {
		log.Println("Google OAuth credentials not set")
		return nil, "", errors.New("OAuth configuration error")
	}

	// Exchange code for access token
	ctx := context.Background()
	token, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		log.Printf("Failed to exchange code: %v", err)
		return nil, "", errors.New("failed to exchange code")
	}

	// Fetch user info
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		log.Printf("Failed to create userinfo request: %v", err)
		return nil, "", errors.New("failed to create userinfo request")
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to get user info: %v", err)
		return nil, "", errors.New("failed to get user info")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Userinfo endpoint returned status: %d", resp.StatusCode)
		return nil, "", errors.New("failed to get user info")
	}

	var userInfo googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		log.Printf("Failed to decode user info: %v", err)
		return nil, "", errors.New("failed to decode user info")
	}

	// Validate email
	if userInfo.Email == "" {
		log.Println("No email provided by Google")
		return nil, "", errors.New("no email provided")
	}

	// Check if user exists
	user, err := s.userRepo.FindByEmail(userInfo.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Failed to find user: %v", err)
		return nil, "", err
	}
	if user == nil {
		// Register new user
		user = &models.User{
			Email:   userInfo.Email,
			Name:    userInfo.Name,
			Country: "", // Set default or prompt later
		}
		if err := s.userRepo.Create(user); err != nil {
			log.Printf("Failed to create user: %v", err)
			return nil, "", err
		}
	}

	// Generate JWT
	jwtToken, err := s.GenerateJWT(user)
	if err != nil {
		log.Printf("Failed to generate JWT: %v", err)
		return nil, "", errors.New("failed to generate JWT")
	}

	return user, jwtToken, nil
}

func (s *AuthService) UpdateProfile(userID uint, name, country string, addresses []string) error {

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	user.Name = name
	user.Country = country
	// Addresses as JSON; add field to User model if needed: Addresses jsonb
	return s.userRepo.Update(user)
}

func (s *AuthService) ResetPassword(email, newPassword string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)
	return s.userRepo.Update(user)
}

```

---

---
## 📊 Summary
- Total files: 41
- Total size: 194.73 KB
- File types: .go
