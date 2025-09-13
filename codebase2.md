# Codebase Analysis: perth_myp
Generated: 2025-09-13 03:21:02
---

## 📂 Project Structure
```tree
📁 perth_myp
├── cmd/
│   └── main.go
├── internal/
│   ├── api/
│   │   ├── customer/
│   │   │   ├── handlers/
│   │   │   │   ├── auth_handler.go
│   │   │   │   └── customer_handlers.go
│   │   │   └── routes.go
│   │   └── merchant/
│   │       ├── handlers/
│   │       │   ├── auth_handler.go
│   │       │   └── merchant_handlers.go
│   │       └── routes.go
│   ├── config/
│   │   └── config.go
│   ├── db/
│   │   ├── models/
│   │   │   ├── cart.go
│   │   │   ├── cart_item.go
│   │   │   ├── category.go
│   │   │   ├── inventory.go
│   │   │   ├── merchant.go
│   │   │   ├── order.go
│   │   │   ├── order_item.go
│   │   │   ├── payment.go
│   │   │   ├── payout.go
│   │   │   ├── product.go
│   │   │   └── user.go
│   │   ├── repositories/
│   │   │   ├── cart_item_repositry.go
│   │   │   ├── cart_repositry.go
│   │   │   ├── category_repositry.go
│   │   │   ├── inventory_repository.go
│   │   │   ├── merchant_repositry.go
│   │   │   ├── order_item_repository.go
│   │   │   ├── order_repository.go
│   │   │   ├── payment_repository.go
│   │   │   ├── payout_repository.go
│   │   │   ├── product_repositry.go
│   │   │   └── user_repository.go
│   │   └── db.go
│   ├── domain/
│   │   ├── cart/
│   │   │   └── cart_service.go
│   │   ├── merchant/
│   │   │   └── merchant_service.go
│   │   ├── notifications/
│   │   │   └── notifcation_service.go
│   │   ├── order/
│   │   │   └── order_service.go
│   │   ├── payment/
│   │   │   └── payment_service.go
│   │   ├── payout/
│   │   │   └── payout_service.go
│   │   ├── pricing/
│   │   │   └── pricing_service.go
│   │   ├── product/
│   │   │   └── product_service.go
│   │   ├── test/
│   │   │   └── test_service.go
│   │   └── user/
│   │       └── user_service.go
│   ├── events/
│   ├── jobs/
│   ├── middleware/
│   │   ├── auth.go
│   │   └── rate_limit.go
│   ├── shared/
│   ├── tests/
│   │   ├── integration/
│   │   ├── mocks/
│   │   └── unit/
│   │       ├── test_handlers.go
│   │       └── test_service.go
│   └── utils/
│       ├── blacklist.go
│       └── redis.go
├── pkg/
│   ├── payments/
│   └── tests/
│       ├── mocks/
│       └── unit/
└── .env
```
---

## 📄 File Contents
### .env
- Size: 0.72 KB
- Lines: 10
- Last Modified: 2025-09-13 02:27:47

<xaiArtifact artifact_id="6ede2fda-a768-4f6f-9dfd-a22dfa1fecf2" artifact_version_id="0c59c513-4b34-4db9-8c04-43ed06136fa4" title=".env" contentType="text/plain">
```plain
DB_DSN=postgresql://neondb_owner:npg_CcwoeLb6V1XH@ep-wild-haze-adu0bdvq-pooler.c-2.us-east-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require
JWT_SECRET=9072a74677e95918103b4993b96ef0455995408610c82c3cb3433f718d4838e0
GOOGLE_CLIENT_ID=269870327937-9qlv0sl9lt374slkcqicpus76tnk9cle.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=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
REDIS_ADDR=testing-zrth-cbdi-356640.leapcell.cloud:6379
REDIS_PASS=Ae00000HR66rpYGCp6yeR5xEA0hCzBidoSQqoRfbSr78m/kXbVTTMxpwptseLynFm2zhNKt
REDIS_DB=0
REDIS_URL=rediss://default:Ae00000HR66rpYGCp6yeR5xEA0hCzBidoSQqoRfbSr78m/kXbVTTMxpwptseLynFm2zhNKt@testing-zrth-cbdi-356640.leapcell.cloud:6379
BASE_URL=https://perthmyp-production.up.railway.app
PORT=8080

```
</xaiArtifact>

---
### cmd/main.go
- Size: 2.10 KB
- Lines: 86
- Last Modified: 2025-09-12 21:54:43

<xaiArtifact artifact_id="dee36a75-ba7e-4d23-81fa-e7d66eca4183" artifact_version_id="f6505d80-e0ec-4c13-9bd3-41b72735c602" title="cmd/main.go" contentType="text/go">
```go
package main

import (
	"fmt"
	"log"
	"os"

	customer "api-customer-merchant/internal/api/customer"
	merchant "api-customer-merchant/internal/api/merchant"
	"api-customer-merchant/internal/config"

	//"api-customer-merchant/internal/middleware"
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "api-customer-merchant/docs" // Import generated docs
)

// @title Multivendor API
// @version 1.0
// @description API for customer and merchant authentication in a multivendor platform
// @termsOfService http://example.com/terms/
// @contact.name API Support
// @contact.email support@example.com
// @license.name MIT
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	// if err := godotenv.Load(); err != nil {
    //          log.Println("No .env file found, relying on environment variables")
    //      }
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	conf := config.Load()
	utils.InitRedis(conf)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET not set")
	}

		//  if err := godotenv.Load(); err != nil {
		// 	log.Fatal("Error loading .env file")
		// }
		// secret := os.Getenv("JWT_SECRET")
		// if secret == "" {
		// 	log.Fatal("JWT_SECRET not set")
		// }
	// Connect to database and migrate
	db.Connect()
	db.AutoMigrate()
	r := gin.Default()
	r.Use(gin.Recovery())


	// Create single router
	//r := gin.Default()

	// Customer routes under /customer
	customer.RegisterRoutes(r)
    merchant.RegisterRoutes(r)

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Get port from environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Run on 0.0.0.0:port for Railway compatibility
	addr := fmt.Sprintf("0.0.0.0:%s", port)
	log.Printf("Example app listening on port %s", port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("API failed: %v", err)
	}
}
```
</xaiArtifact>

---
### internal/utils/blacklist.go
- Size: 0.93 KB
- Lines: 31
- Last Modified: 2025-09-12 19:23:37

<xaiArtifact artifact_id="e79a0392-c5e7-44ea-99a9-922e6a63da46" artifact_version_id="8c7ea3e4-da19-4427-b331-dae81f799c4a" title="internal/utils/blacklist.go" contentType="text/go">
```go
package utils

import (
    "context"
    "errors"
    "log"
    "time"
    //"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// Add adds a token to the Redis blacklist with an expiration
func Add(token string) error {
    if RedisClient == nil {
        log.Println("RedisClient is nil, cannot add token to blacklist")
        return errors.New("redis client not initialized")
    }
    // Set token in Redis with a 24-hour expiration
    return RedisClient.Set(ctx, "blacklist:"+token, "true", 24*time.Hour).Err()
}

// IsBlacklisted checks if a token is in the Redis blacklist
func IsBlacklisted(token string) bool {
    if RedisClient == nil {
        log.Println("RedisClient is nil, skipping blacklist check")
        return false // Fallback to allow operation if Redis is unavailable
    }
    _, err := RedisClient.Get(ctx, "blacklist:"+token).Result() // Line 22
    return err == nil // Token exists in Redis if no error
}
```
</xaiArtifact>

---
### internal/utils/redis.go
- Size: 1.41 KB
- Lines: 56
- Last Modified: 2025-09-12 19:24:01

<xaiArtifact artifact_id="a3fe1764-87b3-4c37-b17e-dc53433a7fcd" artifact_version_id="af034a85-9564-4d3c-8c69-c78d6a37a041" title="internal/utils/redis.go" contentType="text/go">
```go
package utils

import (
	"api-customer-merchant/internal/config"
	"context"
	"crypto/tls"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis(conf *config.Config) {
    
    RedisClient = redis.NewClient(&redis.Options{
        Addr:     conf.RedisAddr,
        Password: conf.RedisPass,
        DB:       conf.RedisDB,
        TLSConfig: &tls.Config{
            MinVersion: tls.VersionTLS12, 
            InsecureSkipVerify: true, // Use with caution, only for testing
        },
    })

    ctx := context.Background()
    if err := RedisClient.Ping(ctx).Err(); err != nil {
       log.Printf("Failed to connect to Redis: %v, continuing without caching", err)
        RedisClient = nil // Fallback to avoid crashes
    } else {
        log.Println("Connected to Redis successfully")
    }
}

// Helper to get cached value or fetch and cache
func GetOrSetCache(ctx context.Context, key string, ttl time.Duration, fetch func() (any, error)) (any, error) {
    val, err := RedisClient.Get(ctx, key).Result()
    if err == nil {
        return val, nil // Deserialize if needed (e.g., JSON)
    }
    if err != redis.Nil {
        return nil, err
    }

    data, err := fetch()
    if err != nil {
        return nil, err
    }

    // Serialize if complex (e.g., JSON marshal)
    if err := RedisClient.Set(ctx, key, data, ttl).Err(); err != nil {
        return nil, err
    }
    return data, nil
}
```
</xaiArtifact>

---
### internal/config/config.go
- Size: 0.72 KB
- Lines: 31
- Last Modified: 2025-09-12 19:24:45

<xaiArtifact artifact_id="6e8e0c2c-7804-4579-b816-7600f2ca4d01" artifact_version_id="a431c43b-e34b-413c-9913-c9c9b44bde07" title="internal/config/config.go" contentType="text/go">
```go
package config


import (
	//"log"
	//"net/url"
	"os"
	"strconv"
	//"time"
)

type Config struct {
	RedisAddr string
	RedisPass string
	RedisDB   int
	// Other fields...
}

 func Load() *Config {
 	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
 	// AccessTokenExp, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXP_MINUTES")) // e.g., 15
 	// RefreshTokenExp, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXP_DAYS"))  // e.g., 7
 	 // AccessTokenExp = time.Duration(AccessTokenExp) * time.Minute
 	 // RefreshTokenExp = time.Duration(RefreshTokenExp) * 24 * time.Hour
 	return &Config{
 		RedisAddr: os.Getenv("REDIS_ADDR"), // e.g., "localhost:6379"
 		RedisPass: os.Getenv("REDIS_PASS"),
 		RedisDB:   redisDB, // Default 0
 		// ...
 	}
 }

```
</xaiArtifact>

---
### internal/middleware/rate_limit.go
- Size: 0.45 KB
- Lines: 20
- Last Modified: 2025-09-12 19:22:20

<xaiArtifact artifact_id="ac915ca4-319b-460b-b56f-fc84cff9d12b" artifact_version_id="46727b84-3e25-4cb3-b068-b88b73fe1441" title="internal/middleware/rate_limit.go" contentType="text/go">
```go
package middleware

import (
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
    "golang.org/x/time/rate"
)

func RateLimitMiddleware() gin.HandlerFunc {
    limiter := rate.NewLimiter(rate.Every(time.Minute), 100) // 100/min
    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```
</xaiArtifact>

---
### internal/middleware/auth.go
- Size: 1.27 KB
- Lines: 55
- Last Modified: 2025-09-05 10:40:08

<xaiArtifact artifact_id="c8d3dfd6-02a9-4806-89b0-a659172deed6" artifact_version_id="90e1f3dd-1b8b-4572-8f76-cf000486adcd" title="internal/middleware/auth.go" contentType="text/go">
```go
package middleware

import (
	"net/http"
	"os"
	"strings"

	"api-customer-merchant/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(entityType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if utils.IsBlacklisted(tokenString) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is blacklisted"})
			c.Abort()
			return
		}
		key := os.Getenv("JWT_SECRET")

		secret := []byte(key) // Load from env
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || claims["entityType"] != entityType {
				c.JSON(http.StatusForbidden, gin.H{"error": "Invalid entity type"})
				c.Abort()
				return
			}

			c.Set("entityId", claims["id"])
			c.Next()
	}
}
```
</xaiArtifact>

---
### internal/db/db.go
- Size: 6.09 KB
- Lines: 215
- Last Modified: 2025-09-13 03:14:32

<xaiArtifact artifact_id="ae82f87e-21a6-494f-87d7-2ea8f07136a6" artifact_version_id="cafc6d43-ed9f-4d97-a98d-ab548490e3f0" title="internal/db/db.go" contentType="text/go">
```go
package db

import (
	"api-customer-merchant/internal/db/models"
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
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
    // log.Println("Migrating Product ecosystem (Product, Variant, Media)...")
    // if err := DB.AutoMigrate(&models.Product{}, &models.Variant{}, &models.Media{}); err != nil {
    //     log.Printf("Failed to migrate Product/Variant/Media: %v", err)
    //     return
    // }

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

    // log.Println("Migrating Cart...")
    // if err := DB.AutoMigrate(&models.Cart{}); err != nil {
    //     log.Printf("Failed to migrate Cart: %v", err)
    //     return
    // }

    // log.Println("Migrating CartItem...")
    // if err := DB.AutoMigrate(&models.CartItem{}); err != nil {
    //     log.Printf("Failed to migrate CartItem: %v", err)
    //     return
    // }

    // log.Println("Migrating Order...")
    // if err := DB.AutoMigrate(&models.Order{}); err != nil {
    //     log.Printf("Failed to migrate Order: %v", err)
    //     return
    // }

    // log.Println("Migrating OrderItem...")
    // if err := DB.AutoMigrate(&models.OrderItem{}); err != nil {
    //     log.Printf("Failed to migrate OrderItem: %v", err)
    //     return
    // }

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
        &models.User{},
        //&models.MerchantApplication{},
         //&models.Product{},
         //&models.Variant{},
         //&models.Media{},
        // &models.Cart{},
        // &models.Order{},
        // &models.OrderItem{},
        // &models.CartItem{},
         //&models.Category{},
        // &models.Inventory{},
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
</xaiArtifact>

---
### internal/domain/pricing/pricing_service.go
- Size: 1.03 KB
- Lines: 40
- Last Modified: 2025-09-05 12:12:34

<xaiArtifact artifact_id="b33ec3a7-b291-40df-804e-f90ee8642a62" artifact_version_id="b803f000-8732-4fcc-86e3-ffce9fd5f701" title="internal/domain/pricing/pricing_service.go" contentType="text/go">
```go
package pricing

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
```
</xaiArtifact>

---
### internal/domain/product/product_service.go
- Size: 4.00 KB
- Lines: 144
- Last Modified: 2025-09-05 12:08:48

<xaiArtifact artifact_id="2cee6e37-e4f7-49ee-b1e4-10462a9cc99d" artifact_version_id="07e30a31-1d1f-4929-8b95-9781f88ef48f" title="internal/domain/product/product_service.go" contentType="text/go">
```go
package product

import (
	"errors"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"
	"strings"
)

type ProductService struct {
	productRepo  *repositories.ProductRepository
	inventoryRepo *repositories.InventoryRepository
}

func NewProductService(productRepo *repositories.ProductRepository, inventoryRepo *repositories.InventoryRepository) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		inventoryRepo: inventoryRepo,
	}
}

// GetProductByID retrieves a product by its ID
func (s *ProductService) GetProductByID(id uint) (*models.Product, error) {
	if id == 0 {
		return nil, errors.New("invalid product ID")
	}
	return s.productRepo.FindByID(id)
}

// GetProductBySKU retrieves a product by its SKU
func (s *ProductService) GetProductBySKU(sku string) (*models.Product, error) {
	if strings.TrimSpace(sku) == "" {
		return nil, errors.New("SKU cannot be empty")
	}
	return s.productRepo.FindBySKU(sku)
}

// SearchProductsByName searches for products by name (case-insensitive)
func (s *ProductService) SearchProductsByName(name string) ([]models.Product, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("search name cannot be empty")
	}
	return s.productRepo.SearchByName(name)
}

// GetProductsByCategory retrieves products in a category
// In ProductService
func (s *ProductService) GetProductsByCategory(categoryID uint, limit, offset int) ([]models.Product, error) {
    if categoryID == 0 {
        return nil, errors.New("invalid category ID")
    }
    return s.productRepo.FindByCategoryWithPagination(categoryID, limit, offset)
}

// CreateProduct creates a new product for a merchant
func (s *ProductService) CreateProduct(product *models.Product, merchantID uint) error {
	if product == nil {
		return errors.New("product cannot be nil")
	}
	if merchantID == 0 {
		return errors.New("invalid merchant ID")
	}
	if strings.TrimSpace(product.Name) == "" {
		return errors.New("product name cannot be empty")
	}
	if strings.TrimSpace(product.SKU) == "" {
		return errors.New("SKU cannot be empty")
	}
	if product.Price <= 0 {
		return errors.New("price must be positive")
	}
	if product.CategoryID == 0 {
		return errors.New("category ID must be set")
	}

	// Check if SKU is unique
	if _, err := s.productRepo.FindBySKU(product.SKU); err == nil {
		return errors.New("SKU already exists")
	}

	product.MerchantID = merchantID
	return s.productRepo.Create(product)
}

// UpdateProduct updates an existing product (merchant only)
func (s *ProductService) UpdateProduct(product *models.Product, merchantID uint) error {
	if product == nil || product.ID == 0 {
		return errors.New("invalid product or product ID")
	}
	if merchantID == 0 {
		return errors.New("invalid merchant ID")
	}
	if strings.TrimSpace(product.Name) == "" {
		return errors.New("product name cannot be empty")
	}
	if strings.TrimSpace(product.SKU) == "" {
		return errors.New("SKU cannot be empty")
	}
	if product.Price <= 0 {
		return errors.New("price must be positive")
	}
	if product.CategoryID == 0 {
		return errors.New("category ID must be set")
	}

	// Verify product belongs to merchant
	existing, err := s.productRepo.FindByID(product.ID)
	if err != nil {
		return err
	}
	if existing.MerchantID != merchantID {
		return errors.New("product does not belong to merchant")
	}

	// Check if SKU is unique (excluding current product)
	if p, err := s.productRepo.FindBySKU(product.SKU); err == nil && p.ID != product.ID {
		return errors.New("SKU already exists")
	}

	return s.productRepo.Update(product)
}

// DeleteProduct deletes a product (merchant only)
func (s *ProductService) DeleteProduct(id uint, merchantID uint) error {
	if id == 0 {
		return errors.New("invalid product ID")
	}
	if merchantID == 0 {
		return errors.New("invalid merchant ID")
	}

	// Verify product belongs to merchant
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return err
	}
	if product.MerchantID != merchantID {
		return errors.New("product does not belong to merchant")
	}

	return s.productRepo.Delete(id)
}



```
</xaiArtifact>

---
### internal/domain/cart/cart_service.go
- Size: 4.51 KB
- Lines: 172
- Last Modified: 2025-09-04 14:30:14

<xaiArtifact artifact_id="c6ae89ce-1b30-4894-81ec-bbfd3df1ed59" artifact_version_id="fef22c54-2690-4f8c-b7ac-22c7cd8346fa" title="internal/domain/cart/cart_service.go" contentType="text/go">
```go
package cart

import (
	"errors"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"

	"gorm.io/gorm"
)

type CartService struct {
	cartRepo      *repositories.CartRepository
	cartItemRepo  *repositories.CartItemRepository
	productRepo   *repositories.ProductRepository
	inventoryRepo *repositories.InventoryRepository
}

func NewCartService(cartRepo *repositories.CartRepository, cartItemRepo *repositories.CartItemRepository, productRepo *repositories.ProductRepository, inventoryRepo *repositories.InventoryRepository) *CartService {
	return &CartService{
		cartRepo:      cartRepo,
		cartItemRepo:  cartItemRepo,
		productRepo:   productRepo,
		inventoryRepo: inventoryRepo,
	}
}

// GetActiveCart retrieves or creates an active cart for a user
func (s *CartService) GetActiveCart(userID uint) (*models.Cart, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}

	// Try to find an active cart
	cart, err := s.cartRepo.FindActiveCart(userID)
	if err == nil {
		return cart, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Create a new active cart if none exists
	cart = &models.Cart{
		UserID: userID,
		Status: models.CartStatusActive,
	}
	if err := s.cartRepo.Create(cart); err != nil {
		return nil, err
	}
	return s.cartRepo.FindByID(cart.ID)
}

// AddItemToCart adds a product to the user's active cart
func (s *CartService) AddItemToCart(userID, productID uint, quantity int) (*models.Cart, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}
	if productID == 0 {
		return nil, errors.New("invalid product ID")
	}
	if quantity <= 0 {
		return nil, errors.New("quantity must be positive")
	}

	// Get active cart
	cart, err := s.GetActiveCart(userID)
	if err != nil {
		return nil, err
	}

	// Check if product exists
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Check stock availability
	inventory, err := s.inventoryRepo.FindByProductID(productID)
	if err != nil {
		return nil, errors.New("inventory not found")
	}
	if inventory.StockQuantity < quantity {
		return nil, errors.New("insufficient stock")
	}

	// Check if product is already in cart
	cartItem, err := s.cartItemRepo.FindByProductIDAndCartID(productID, cart.ID)
	if err == nil {
		// Update quantity if item exists
		newQuantity := cartItem.Quantity + quantity
		if inventory.StockQuantity < newQuantity {
			return nil, errors.New("insufficient stock for updated quantity")
		}
		if err := s.cartItemRepo.UpdateQuantity(cartItem.ID, newQuantity); err != nil {
			return nil, err
		}
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new cart item
		cartItem = &models.CartItem{
			CartID:     cart.ID,
			ProductID:  productID,
			Quantity:   quantity,
			MerchantID: product.MerchantID,
		}
		if err := s.cartItemRepo.Create(cartItem); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	return s.cartRepo.FindByID(cart.ID)
}

// UpdateCartItemQuantity updates the quantity of a cart item
func (s *CartService) UpdateCartItemQuantity(cartItemID uint, quantity int) (*models.Cart, error) {
	if cartItemID == 0 {
		return nil, errors.New("invalid cart item ID")
	}
	if quantity <= 0 {
		return nil, errors.New("quantity must be positive")
	}

	// Get cart item
	cartItem, err := s.cartItemRepo.FindByID(cartItemID)
	if err != nil {
		return nil, errors.New("cart item not found")
	}

	// Check stock availability
	inventory, err := s.inventoryRepo.FindByProductID(cartItem.ProductID)
	if err != nil {
		return nil, errors.New("inventory not found")
	}
	if inventory.StockQuantity < quantity {
		return nil, errors.New("insufficient stock")
	}

	// Update quantity
	if err := s.cartItemRepo.UpdateQuantity(cartItemID, quantity); err != nil {
		return nil, err
	}

	return s.cartRepo.FindByID(cartItem.CartID)
}

// RemoveCartItem removes an item from the cart
func (s *CartService) RemoveCartItem(cartItemID uint) (*models.Cart, error) {
	if cartItemID == 0 {
		return nil, errors.New("invalid cart item ID")
	}

	// Get cart item to find cart ID
	cartItem, err := s.cartItemRepo.FindByID(cartItemID)
	if err != nil {
		return nil, errors.New("cart item not found")
	}

	// Delete cart item
	if err := s.cartItemRepo.Delete(cartItemID); err != nil {
		return nil, err
	}

	return s.cartRepo.FindByID(cartItem.CartID)
}

func (s *CartService) GetCartItemByID(cartItemID uint) (*models.CartItem, error) {
	if cartItemID == 0 {
		return nil, errors.New("invalid cart item ID")
	}
	return s.cartItemRepo.FindByID(cartItemID)
}
```
</xaiArtifact>

---
### internal/domain/payout/payout_service.go
- Size: 1.54 KB
- Lines: 62
- Last Modified: 2025-09-04 13:02:24

<xaiArtifact artifact_id="bfab3c8e-6d4c-4792-b562-47ca19423337" artifact_version_id="56fa9355-1e58-44d4-8ba9-e39639395db1" title="internal/domain/payout/payout_service.go" contentType="text/go">
```go
package payout


import (
	"errors"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"
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
```
</xaiArtifact>

---
### internal/domain/payment/payment_service.go
- Size: 2.50 KB
- Lines: 96
- Last Modified: 2025-09-04 12:55:06

<xaiArtifact artifact_id="84708222-cece-4e67-908f-ddee1c8ba540" artifact_version_id="479d2f3e-eb1b-4de8-a7df-d33c23aa402c" title="internal/domain/payment/payment_service.go" contentType="text/go">
```go
package payment

import (
	"errors"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"
)

type PaymentService struct {
	paymentRepo *repositories.PaymentRepository
	orderRepo   *repositories.OrderRepository
}

func NewPaymentService(paymentRepo *repositories.PaymentRepository, orderRepo *repositories.OrderRepository) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
	}
}

// ProcessPayment creates a payment for an order (placeholder for Stripe)
func (s *PaymentService) ProcessPayment(orderID uint, amount float64) (*models.Payment, error) {
	if orderID == 0 {
		return nil, errors.New("invalid order ID")
	}
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	// Verify order exists
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	// Verify order amount matches
	if order.TotalAmount != amount {
		return nil, errors.New("payment amount does not match order total")
	}

	// Simulate Stripe payment processing
	payment := &models.Payment{
		OrderID: orderID,
		Amount:  amount,
		Status:  models.PaymentStatusPending,
	}
	if err := s.paymentRepo.Create(payment); err != nil {
		return nil, err
	}

	// Simulate successful payment
	payment.Status = models.PaymentStatusCompleted
	if err := s.paymentRepo.Update(payment); err != nil {
		return nil, err
	}

	return s.paymentRepo.FindByID(payment.ID)
}

// GetPaymentByOrderID retrieves a payment by order ID
func (s *PaymentService) GetPaymentByOrderID(orderID uint) (*models.Payment, error) {
	if orderID == 0 {
		return nil, errors.New("invalid order ID")
	}
	return s.paymentRepo.FindByOrderID(orderID)
}

// GetPaymentsByUserID retrieves all payments for a user
func (s *PaymentService) GetPaymentsByUserID(userID uint) ([]models.Payment, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}
	return s.paymentRepo.FindByUserID(userID)
}

// UpdatePaymentStatus updates the status of a payment
func (s *PaymentService) UpdatePaymentStatus(paymentID uint, status string) (*models.Payment, error) {
	if paymentID == 0 {
		return nil, errors.New("invalid payment ID")
	}
	if err := models.PaymentStatus(status).Valid(); err != nil {
		return nil, err
	}

	payment, err := s.paymentRepo.FindByID(paymentID)
	if err != nil {
		return nil, err
	}

	payment.Status = models.PaymentStatus(status)
	if err := s.paymentRepo.Update(payment); err != nil {
		return nil, err
	}

	return s.paymentRepo.FindByID(paymentID)
}
```
</xaiArtifact>

---
### internal/domain/merchant/merchant_service.go
- Size: 2.47 KB
- Lines: 85
- Last Modified: 2025-09-12 18:42:47

<xaiArtifact artifact_id="399c6d03-03d6-40a4-8364-788dcb9034b1" artifact_version_id="aae6f3d3-cd80-4d26-ba87-b8c4e433b159" title="internal/domain/merchant/merchant_service.go" contentType="text/go">
```go
package merchant


import (
	"context"
	"encoding/json"
	"errors"

	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"
	"github.com/go-playground/validator/v10"
	
)
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
```
</xaiArtifact>

---
### internal/domain/user/user_service.go
- Size: 6.56 KB
- Lines: 242
- Last Modified: 2025-09-12 19:00:40

<xaiArtifact artifact_id="f4ac0f47-482d-4757-9d4b-ee23ce152b14" artifact_version_id="8da19779-f33a-4e75-8acb-37576c3a57ee" title="internal/domain/user/user_service.go" contentType="text/go">
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
	userRepo     *repositories.UserRepository
	
}


func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		
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
		"id":        float64(id),
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
		 Scopes:       []string{
            "https://www.googleapis.com/auth/userinfo.email",
            "https://www.googleapis.com/auth/userinfo.profile",
            "openid"},
		Endpoint:     google.Endpoint,
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
</xaiArtifact>

---
### internal/domain/notifications/notifcation_service.go
- Size: 0.74 KB
- Lines: 31
- Last Modified: 2025-09-05 12:15:35

<xaiArtifact artifact_id="9d672238-7105-407e-992a-07bbedc6ce2b" artifact_version_id="f489d88c-2712-4856-9541-7a2c766a55e6" title="internal/domain/notifications/notifcation_service.go" contentType="text/go">
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
</xaiArtifact>

---
### internal/domain/order/order_service.go
- Size: 4.01 KB
- Lines: 151
- Last Modified: 2025-09-04 12:36:21

<xaiArtifact artifact_id="100f2bc9-1029-4b4e-985e-6d55a07d1365" artifact_version_id="42368950-e091-41c5-99f4-a846b8b72ef5" title="internal/domain/order/order_service.go" contentType="text/go">
```go
package order

import (
	"errors"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"

	"gorm.io/gorm"
)

type OrderService struct {
	orderRepo     *repositories.OrderRepository
	orderItemRepo *repositories.OrderItemRepository
	cartRepo      *repositories.CartRepository
	cartItemRepo  *repositories.CartItemRepository
	inventoryRepo *repositories.InventoryRepository
}

func NewOrderService(orderRepo *repositories.OrderRepository, orderItemRepo *repositories.OrderItemRepository, cartRepo *repositories.CartRepository, cartItemRepo *repositories.CartItemRepository, inventoryRepo *repositories.InventoryRepository) *OrderService {
	return &OrderService{
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
		cartRepo:      cartRepo,
		cartItemRepo:  cartItemRepo,
		inventoryRepo: inventoryRepo,
	}
}

// CreateOrderFromCart creates an order from the user's active cart
func (s *OrderService) CreateOrderFromCart(userID uint) (*models.Order, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}

	// Get active cart
	cart, err := s.cartRepo.FindActiveCart(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no active cart found")
		}
		return nil, err
	}

	// Check if cart has items
	cartItems, err := s.cartItemRepo.FindByCartID(cart.ID)
	if err != nil {
		return nil, err
	}
	if len(cartItems) == 0 {
		return nil, errors.New("cart is empty")
	}

	// Validate stock for all items
	for _, item := range cartItems {
		inventory, err := s.inventoryRepo.FindByProductID(item.ProductID)
		if err != nil {
			return nil, errors.New("inventory not found for product")
		}
		if inventory.StockQuantity < item.Quantity {
			return nil, errors.New("insufficient stock for product")
		}
	}

	// Calculate total amount
	var totalAmount float64
	for _, item := range cartItems {
		totalAmount += float64(item.Quantity) * item.Product.Price
	}

	// Create order
	order := &models.Order{
		UserID:      userID,
		TotalAmount: totalAmount,
		Status:      models.OrderStatusPending,
	}
	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}

	// Create order items and update inventory
	for _, item := range cartItems {
		orderItem := &models.OrderItem{
			OrderID:    order.ID,
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			Price:      item.Product.Price,
			MerchantID: item.MerchantID,
		}
		if err := s.orderItemRepo.Create(orderItem); err != nil {
			return nil, err
		}
		// Update inventory
		if err := s.inventoryRepo.UpdateStock(item.ProductID, -item.Quantity); err != nil {
			return nil, err
		}
	}

	// Mark cart as converted
	cart.Status = models.CartStatusConverted
	if err := s.cartRepo.Update(cart); err != nil {
		return nil, err
	}

	return s.orderRepo.FindByID(order.ID)
}

// GetOrderByID retrieves an order by ID
func (s *OrderService) GetOrderByID(id uint) (*models.Order, error) {
	if id == 0 {
		return nil, errors.New("invalid order ID")
	}
	return s.orderRepo.FindByID(id)
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
func (s *OrderService) UpdateOrderStatus(orderID uint, status string) (*models.Order, error) {
	if orderID == 0 {
		return nil, errors.New("invalid order ID")
	}
	if err := models.OrderStatus(status).Valid(); err != nil {
		return nil, err
	}

	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}

	order.Status = models.OrderStatus(status)
	if err := s.orderRepo.Update(order); err != nil {
		return nil, err
	}

	return s.orderRepo.FindByID(orderID)
}
```
</xaiArtifact>

---
### internal/domain/test/test_service.go
- Size: 5.93 KB
- Lines: 171
- Last Modified: 2025-09-04 12:40:28

<xaiArtifact artifact_id="a5cb05ea-8e37-487e-a188-8b5c65b2f1d5" artifact_version_id="dd4bff31-7409-499b-8049-9be37c196443" title="internal/domain/test/test_service.go" contentType="text/go">
```go
package test

// import (
// 	"fmt"
// 	"api-customer-merchant/internal/db"
// 	"api-customer-merchant/internal/db/models"
// 	"api-customer-merchant/internal/db/repositories"
// 	"api-customer-merchant/internal/domain/cart"
// 	"api-customer-merchant/internal/domain/order"
// 	"api-customer-merchant/internal/domain/payment"
// 	"api-customer-merchant/internal/domain/payout"
// 	"api-customer-merchant/internal/domain/product"
// )

// func TestServices() {
// 	db.Connect()
// 	db.AutoMigrate()

// 	// Setup repositories
// 	productRepo := repositories.NewProductRepository()
// 	inventoryRepo := repositories.NewInventoryRepository()
// 	cartRepo := repositories.NewCartRepository()
// 	cartItemRepo := repositories.NewCartItemRepository()
// 	orderRepo := repositories.NewOrderRepository()
// 	orderItemRepo := repositories.NewOrderItemRepository()
// 	paymentRepo := repositories.NewPaymentRepository()
// 	payoutRepo := repositories.NewPayoutRepository()

// 	// Setup services
// 	productService := product.NewProductService(productRepo, inventoryRepo)
// 	cartService := cart.NewCartService(cartRepo, cartItemRepo, productRepo, inventoryRepo)
// 	orderService := order.NewOrderService(orderRepo, orderItemRepo, cartRepo, cartItemRepo, inventoryRepo)
// 	paymentService := payment.NewPaymentService(paymentRepo, orderRepo)
// 	payoutService := payout.NewPayoutService(payoutRepo)

// 	// Insert test data
// 	user := &models.User{Email: "test@example.com", Name: "Test User", Password: "$2a$10$examplehashedpassword", Country: "Nigeria"}
// 	if err := db.DB.Create(user).Error; err != nil {
// 		fmt.Println("Error creating user:", err)
// 	}

// 	merchant := &models.Merchant{MerchantBasicInfo: models.MerchantBasicInfo{Name: "Test Merchant", StoreName: "Test Store", PersonalEmail: "personal@example.com", WorkEmail: "work@example.com", Password: "$2a$10$examplehashedpassword"}, Status: models.MerchantStatusApproved}
// 	if err := db.DB.Create(merchant).Error; err != nil {
// 		fmt.Println("Error creating merchant:", err)
// 	}

// 	category := &models.Category{Name: "Electronics", Attributes: map[string]interface{}{"color": []string{"black"}}}
// 	if err := db.DB.Create(category).Error; err != nil {
// 		fmt.Println("Error creating category:", err)
// 	}

// 	// Test ProductService
// 	product := &models.Product{Name: "Smartphone", SKU: "SM123", Price: 599.99, CategoryID: category.ID, MerchantID: merchant.ID}
// 	if err := productService.CreateProduct(product, merchant.ID); err != nil {
// 		fmt.Println("Error creating product:", err)
// 	}
// 	p, err := productService.GetProductBySKU("SM123")
// 	if err != nil {
// 		fmt.Println("Error finding product by SKU:", err)
// 	} else {
// 		fmt.Printf("Found product by SKU: %+v\n", p)
// 	}
// 	products, err := productService.SearchProductsByName("Smart")
// 	if err != nil {
// 		fmt.Println("Error searching products:", err)
// 	} else {
// 		fmt.Printf("Found products by name: %+v\n", products)
// 	}

// 	// Add inventory
// 	inventory := &models.Inventory{ProductID: p.ID, StockQuantity: 100, LowStockThreshold: 10}
// 	if err := inventoryRepo.Create(inventory); err != nil {
// 		fmt.Println("Error creating inventory:", err)
// 	}

// 	// Test inventory stock update
// 	if err := inventoryRepo.UpdateStock(p.ID, -10); err != nil {
// 		fmt.Println("Error updating stock:", err)
// 	} else {
// 		inv, _ := inventoryRepo.FindByProductID(p.ID)
// 		fmt.Printf("Updated stock: %+v\n", inv)
// 	}

// 	// Test CartService
// 	cart, err := cartService.GetActiveCart(user.ID)
// 	if err != nil {
// 		fmt.Println("Error getting active cart:", err)
// 	} else {
// 		fmt.Printf("Active cart: %+v\n", cart)
// 	}

// 	// Add item to cart
// 	cart, err = cartService.AddItemToCart(user.ID, p.ID, 2)
// 	if err != nil {
// 		fmt.Println("Error adding item to cart:", err)
// 	} else {
// 		fmt.Printf("Cart after adding item: %+v\n", cart)
// 	}

// 	// Test OrderService
// 	order, err := orderService.CreateOrderFromCart(user.ID)
// 	if err != nil {
// 		fmt.Println("Error creating order:", err)
// 	} else {
// 		fmt.Printf("Created order: %+v\n", order)
// 	}

// 	// Verify stock after order
// 	inv, err := inventoryRepo.FindByProductID(p.ID)
// 	if err != nil {
// 		fmt.Println("Error checking stock after order:", err)
// 	} else {
// 		fmt.Printf("Stock after order: %+v\n", inv)
// 	}

// 	// Retrieve order
// 	o, err := orderService.GetOrderByID(order.ID)
// 	if err != nil {
// 		fmt.Println("Error finding order by ID:", err)
// 	} else {
// 		fmt.Printf("Found order by ID: %+v\n", o)
// 	}

// 	// Update order status
// 	o, err = orderService.UpdateOrderStatus(order.ID, string(models.OrderStatusShipped))
// 	if err != nil {
// 		fmt.Println("Error updating order status:", err)
// 	} else {
// 		fmt.Printf("Updated order status: %+v\n", o)
// 	}

// 	// Test PaymentService
// 	payment, err := paymentService.ProcessPayment(order.ID, order.TotalAmount)
// 	if err != nil {
// 		fmt.Println("Error processing payment:", err)
// 	} else {
// 		fmt.Printf("Processed payment: %+v\n", payment)
// 	}

// 	// Retrieve payment
// 	pay, err := paymentService.GetPaymentByOrderID(order.ID)
// 	if err != nil {
// 		fmt.Println("Error finding payment by order ID:", err)
// 	} else {
// 		fmt.Printf("Found payment by order ID: %+v\n", pay)
// 	}

// 	// Update payment status
// 	pay, err = paymentService.UpdatePaymentStatus(payment.ID, string(models.PaymentStatusCompleted))
// 	if err != nil {
// 		fmt.Println("Error updating payment status:", err)
// 	} else {
// 		fmt.Printf("Updated payment status: %+v\n", pay)
// 	}

// 	// Test PayoutService
// 	payout, err := payoutService.CreatePayout(merchant.ID, 500.00)
// 	if err != nil {
// 		fmt.Println("Error creating payout:", err)
// 	} else {
// 		fmt.Printf("Created payout: %+v\n", payout)
// 	}

// 	// Retrieve payout
// 	po, err := payoutService.GetPayoutByID(payout.ID)
// 	if err != nil {
// 		fmt.Println("Error finding payout by ID:", err)
// 	} else {
// 		fmt.Printf("Found payout by ID: %+v\n", po)
// 	}
// }
```
</xaiArtifact>

---
### internal/db/repositories/payment_repository.go
- Size: 1.51 KB
- Lines: 52
- Last Modified: 2025-09-04 13:06:49

<xaiArtifact artifact_id="c2ec416d-dcf9-4680-b2bb-d55ce04a282f" artifact_version_id="175dd546-b28b-4d81-bc2a-14558f7d755a" title="internal/db/repositories/payment_repository.go" contentType="text/go">
```go
package repositories

import (
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"

	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository() *PaymentRepository {
	return &PaymentRepository{db: db.DB}
}

// Create adds a new payment
func (r *PaymentRepository) Create(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

// FindByID retrieves a payment by ID with associated Order and User
func (r *PaymentRepository) FindByID(id uint) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Preload("Order.User").First(&payment, id).Error
	return &payment, err
}

// FindByOrderID retrieves a payment by order ID
func (r *PaymentRepository) FindByOrderID(orderID uint) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Preload("Order.User").Where("order_id = ?", orderID).First(&payment).Error
	return &payment, err
}

// FindByUserID retrieves all payments for a user
func (r *PaymentRepository) FindByUserID(userID uint) ([]models.Payment, error) {
	var payments []models.Payment
	err := r.db.Preload("Order.User").Joins("JOIN orders ON orders.id = payments.order_id").Where("orders.user_id = ?", userID).Find(&payments).Error
	return payments, err
}

// Update modifies an existing payment
func (r *PaymentRepository) Update(payment *models.Payment) error {
	return r.db.Save(payment).Error
}

// Delete removes a payment by ID
func (r *PaymentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Payment{}, id).Error
}
```
</xaiArtifact>

---
### internal/db/repositories/user_repository.go
- Size: 1.11 KB
- Lines: 53
- Last Modified: 2025-09-12 18:34:45

<xaiArtifact artifact_id="9855490e-b3c5-4a73-9725-e48eb866b1f2" artifact_version_id="270021b3-4d64-491f-94a7-f1f8ac035b6b" title="internal/db/repositories/user_repository.go" contentType="text/go">
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
</xaiArtifact>

---
### internal/db/repositories/cart_repositry.go
- Size: 1.82 KB
- Lines: 59
- Last Modified: 2025-09-03 22:35:07

<xaiArtifact artifact_id="4d553849-3e45-49e3-abca-2697eff9d439" artifact_version_id="fceb9b82-9d14-41c0-ab86-f9e860bb9799" title="internal/db/repositories/cart_repositry.go" contentType="text/go">
```go
package repositories

import (
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"

	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository() *CartRepository {
	return &CartRepository{db: db.DB}
}

// Create adds a new cart
func (r *CartRepository) Create(cart *models.Cart) error {
	return r.db.Create(cart).Error
}

// FindByID retrieves a cart by ID with associated User and CartItems
func (r *CartRepository) FindByID(id uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.Preload("User").Preload("CartItems.Product.Merchant").First(&cart, id).Error
	return &cart, err
}

// FindActiveCart retrieves the user's most recent active cart
func (r *CartRepository) FindActiveCart(userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.Preload("CartItems.Product.Merchant").Where("user_id = ? AND status = ?", userID, models.CartStatusActive).Order("created_at DESC").First(&cart).Error
	return &cart, err
}

// FindByUserIDAndStatus retrieves carts for a user by status
func (r *CartRepository) FindByUserIDAndStatus(userID uint, status models.CartStatus) ([]models.Cart, error) {
	var carts []models.Cart
	err := r.db.Preload("CartItems.Product.Merchant").Where("user_id = ? AND status = ?", userID, status).Find(&carts).Error
	return carts, err
}

// FindByUserID retrieves all carts for a user
func (r *CartRepository) FindByUserID(userID uint) ([]models.Cart, error) {
	var carts []models.Cart
	err := r.db.Preload("CartItems.Product.Merchant").Where("user_id = ?", userID).Find(&carts).Error
	return carts, err
}

// Update modifies an existing cart
func (r *CartRepository) Update(cart *models.Cart) error {
	return r.db.Save(cart).Error
}

// Delete removes a cart by ID
func (r *CartRepository) Delete(id uint) error {
	return r.db.Delete(&models.Cart{}, id).Error
}
```
</xaiArtifact>

---
### internal/db/repositories/payout_repository.go
- Size: 1.18 KB
- Lines: 45
- Last Modified: 2025-09-03 20:05:20

<xaiArtifact artifact_id="4ab130ac-b977-44a2-83fc-c0aab821f3ad" artifact_version_id="20d9b5f3-540a-4864-be3b-1893a3b61894" title="internal/db/repositories/payout_repository.go" contentType="text/go">
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
</xaiArtifact>

---
### internal/db/repositories/cart_item_repositry.go
- Size: 1.90 KB
- Lines: 57
- Last Modified: 2025-09-03 22:55:06

<xaiArtifact artifact_id="026d1ba3-e7c1-4b31-8026-0eb7a6b1a51a" artifact_version_id="42cbd909-ea37-4115-92d3-374914c3f829" title="internal/db/repositories/cart_item_repositry.go" contentType="text/go">
```go
package repositories

import (
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"

	"gorm.io/gorm"
)

type CartItemRepository struct {
	db *gorm.DB
}

func NewCartItemRepository() *CartItemRepository {
	return &CartItemRepository{db: db.DB}
}

// Create adds a new cart item
func (r *CartItemRepository) Create(cartItem *models.CartItem) error {
	return r.db.Create(cartItem).Error
}

// FindByID retrieves a cart item by ID with associated Cart, Product, and Merchant
func (r *CartItemRepository) FindByID(id uint) (*models.CartItem, error) {
	var cartItem models.CartItem
	err := r.db.Preload("Cart.User").Preload("Product.Merchant").Preload("Merchant").First(&cartItem, id).Error
	return &cartItem, err
}

// FindByCartID retrieves all cart items for a cart
func (r *CartItemRepository) FindByCartID(cartID uint) ([]models.CartItem, error) {
	var cartItems []models.CartItem
	err := r.db.Preload("Product.Merchant").Preload("Merchant").Where("cart_id = ?", cartID).Find(&cartItems).Error
	return cartItems, err
}

// FindByProductIDAndCartID retrieves a cart item by product ID and cart ID
func (r *CartItemRepository) FindByProductIDAndCartID(productID, cartID uint) (*models.CartItem, error) {
	var cartItem models.CartItem
	err := r.db.Preload("Product.Merchant").Preload("Merchant").Where("product_id = ? AND cart_id = ?", productID, cartID).First(&cartItem).Error
	return &cartItem, err
}

// UpdateQuantity updates the quantity of a cart item
func (r *CartItemRepository) UpdateQuantity(id uint, quantity int) error {
	return r.db.Model(&models.CartItem{}).Where("id = ?", id).Update("quantity", quantity).Error
}

// Update modifies an existing cart item
func (r *CartItemRepository) Update(cartItem *models.CartItem) error {
	return r.db.Save(cartItem).Error
}

// Delete removes a cart item by ID
func (r *CartItemRepository) Delete(id uint) error {
	return r.db.Delete(&models.CartItem{}, id).Error
}
```
</xaiArtifact>

---
### internal/db/repositories/category_repositry.go
- Size: 1.14 KB
- Lines: 45
- Last Modified: 2025-09-03 19:34:54

<xaiArtifact artifact_id="6cff4bf8-6fdb-474d-bc3f-3066f50942b5" artifact_version_id="a17661fd-6b1d-4035-93e8-442643b647d6" title="internal/db/repositories/category_repositry.go" contentType="text/go">
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
</xaiArtifact>

---
### internal/db/repositories/inventory_repository.go
- Size: 1.29 KB
- Lines: 44
- Last Modified: 2025-09-04 12:34:54

<xaiArtifact artifact_id="ee124e01-eff6-49c8-b313-909f6a98c6d0" artifact_version_id="e8c2ec57-7cf6-4bcd-8249-ac2fcf83d6af" title="internal/db/repositories/inventory_repository.go" contentType="text/go">
```go
package repositories

import (
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"

	"gorm.io/gorm"
)

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
func (r *InventoryRepository) FindByProductID(productID uint) (*models.Inventory, error) {
	var inventory models.Inventory
	err := r.db.Where("product_id = ?", productID).First(&inventory).Error
	return &inventory, err
}

// UpdateStock updates the stock quantity for a product
func (r *InventoryRepository) UpdateStock(productID uint, quantityChange int) error {
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
```
</xaiArtifact>

---
### internal/db/repositories/order_repository.go
- Size: 1.54 KB
- Lines: 53
- Last Modified: 2025-09-03 19:55:29

<xaiArtifact artifact_id="dfe0f436-0930-4755-bbcf-c24132517cdb" artifact_version_id="d36753bb-f5a2-40ba-8775-4ea5f2ecab2e" title="internal/db/repositories/order_repository.go" contentType="text/go">
```go
package repositories

import (
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{db: db.DB}
}

// Create adds a new order
func (r *OrderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

// FindByID retrieves an order by ID with associated User and OrderItems
func (r *OrderRepository) FindByID(id uint) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("User").Preload("OrderItems.Product.Merchant").First(&order, id).Error
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
func (r *OrderRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

// Delete removes an order by ID
func (r *OrderRepository) Delete(id uint) error {
	return r.db.Delete(&models.Order{}, id).Error
}
```
</xaiArtifact>

---
### internal/db/repositories/product_repositry.go
- Size: 2.48 KB
- Lines: 72
- Last Modified: 2025-09-05 12:09:22

<xaiArtifact artifact_id="774c99e5-6ef3-4a11-b168-d08dbacba1d7" artifact_version_id="d1e7cf4e-3ed8-4cfe-8225-6944f166b6ee" title="internal/db/repositories/product_repositry.go" contentType="text/go">
```go
package repositories

import (
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"

	"gorm.io/gorm"
)

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
```
</xaiArtifact>

---
### internal/db/repositories/merchant_repositry.go
- Size: 2.22 KB
- Lines: 70
- Last Modified: 2025-09-12 18:43:22

<xaiArtifact artifact_id="53647ade-8517-4ea2-8da4-6ace53b85cbd" artifact_version_id="288f8acb-00f8-49e0-9168-a4fef30ab029" title="internal/db/repositories/merchant_repositry.go" contentType="text/go">
```go
package repositories

import (
	"context"
	"log"

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

func (r *MerchantRepository) GetByUserID(ctx context.Context, uid string) (*models.Merchant, error) {
	var m models.Merchant
	if err := db.DB.WithContext(ctx).Where("user_id = ?", uid).First(&m).Error; err != nil {
		log.Printf("Failed to get merchant by user ID %s: %v", uid, err)
		return nil, err
	}
	return &m, nil
}
```
</xaiArtifact>

---
### internal/db/repositories/order_item_repository.go
- Size: 1.35 KB
- Lines: 45
- Last Modified: 2025-09-03 19:58:35

<xaiArtifact artifact_id="56c601a5-869c-40f3-87f1-5a235eb84b85" artifact_version_id="c60435d8-9e33-419b-aa7d-2ff50af01c99" title="internal/db/repositories/order_item_repository.go" contentType="text/go">
```go
package repositories

import (
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"

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
```
</xaiArtifact>

---
### internal/db/models/payout.go
- Size: 1.21 KB
- Lines: 49
- Last Modified: 2025-09-03 17:35:09

<xaiArtifact artifact_id="efc5e534-20c0-4d93-83ae-e1d15bbf5f4b" artifact_version_id="135dc518-5c57-4a74-8a4c-fe8f0c0a79c8" title="internal/db/models/payout.go" contentType="text/go">
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
	MerchantID       uint         `gorm:"not null" json:"merchant_id"`
	Amount           float64      `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status           PayoutStatus `gorm:"type:varchar(20);not null;default:'Pending'" json:"status"`
	PayoutAccountID  string       `gorm:"size:255;not null" json:"payout_account_id"`
	Merchant         Merchant     `gorm:"foreignKey:MerchantID"`
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
</xaiArtifact>

---
### internal/db/models/cart_item.go
- Size: 0.44 KB
- Lines: 16
- Last Modified: 2025-09-03 17:40:23

<xaiArtifact artifact_id="38f5e457-2598-4a32-b23e-291155b436a8" artifact_version_id="3399b1bd-0839-49ad-9d47-138adf3680f3" title="internal/db/models/cart_item.go" contentType="text/go">
```go
package models

import (
	"gorm.io/gorm"
)

type CartItem struct {
	gorm.Model
	CartID     uint    `gorm:"not null" json:"cart_id"`
	ProductID  uint    `gorm:"not null" json:"product_id"`
	Quantity   int     `gorm:"not null" json:"quantity"`
	MerchantID uint    `gorm:"not null" json:"merchant_id"`
	Cart       Cart    `gorm:"foreignKey:CartID"`
	Product    Product `gorm:"foreignKey:ProductID"`
	Merchant   Merchant `gorm:"foreignKey:MerchantID"`
}
```
</xaiArtifact>

---
### internal/db/models/payment.go
- Size: 1.16 KB
- Lines: 49
- Last Modified: 2025-09-04 13:05:31

<xaiArtifact artifact_id="736e8cd4-687f-401f-b86d-4f2d579d8bcc" artifact_version_id="a93703c6-12ae-458e-96e6-27947962ed8c" title="internal/db/models/payment.go" contentType="text/go">
```go
package models

import (
	"fmt"
	"gorm.io/gorm"
)

// PaymentStatus defines possible payment status values
type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "Pending"
	PaymentStatusCompleted PaymentStatus = "Completed"
	PaymentStatusFailed   PaymentStatus = "Failed"
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

type Payment struct {
	gorm.Model
	OrderID uint          `gorm:"not null" json:"order_id"`
	Amount  float64       `gorm:"not null" json:"amount"`
	Status  PaymentStatus `gorm:"type:varchar(20);not null;default:'Pending'" json:"status"`
	Order   Order         `gorm:"foreignKey:OrderID"`
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
</xaiArtifact>

---
### internal/db/models/cart.go
- Size: 1.09 KB
- Lines: 49
- Last Modified: 2025-09-03 22:38:41

<xaiArtifact artifact_id="007f921e-2155-4a2f-8320-1802f49ced0d" artifact_version_id="c5d4ebdc-143d-43c9-bf1f-e81463b83cad" title="internal/db/models/cart.go" contentType="text/go">
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
	UserID uint       `gorm:"not null" json:"user_id"`
	Status CartStatus `gorm:"type:varchar(20);not null;default:'Active'" json:"status"`
	User   User       `gorm:"foreignKey:UserID"`
	CartItems []CartItem `gorm:"foreignKey:CartID"`
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
```
</xaiArtifact>

---
### internal/db/models/inventory.go
- Size: 0.34 KB
- Lines: 13
- Last Modified: 2025-09-03 17:39:11

<xaiArtifact artifact_id="afac9ce7-2d52-4af8-857b-c23a62626df1" artifact_version_id="d657953e-32fe-48f0-8847-a5a8f05dd179" title="internal/db/models/inventory.go" contentType="text/go">
```go
package models

import (
	"gorm.io/gorm"
)

type Inventory struct {
	gorm.Model
	ProductID         uint   `gorm:"not null" json:"product_id"`
	StockQuantity     int    `gorm:"not null" json:"stock_quantity"`
	LowStockThreshold int    `gorm:"not null;default:10" json:"low_stock_threshold"`
	Product           Product `gorm:"foreignKey:ProductID"`
}
```
</xaiArtifact>

---
### internal/db/models/merchant.go
- Size: 5.36 KB
- Lines: 93
- Last Modified: 2025-09-12 18:44:39

<xaiArtifact artifact_id="ad3ee8b4-d628-4ab0-bf26-1ca7bceaec8d" artifact_version_id="ef77cfb9-f133-43d4-8fe4-49a956d74981" title="internal/db/models/merchant.go" contentType="text/go">
```go
package models


import (
	"time"
	"gorm.io/datatypes"
)

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
```
</xaiArtifact>

---
### internal/db/models/product.go
- Size: 0.54 KB
- Lines: 17
- Last Modified: 2025-09-03 17:38:23

<xaiArtifact artifact_id="ef6cc1a6-4106-4051-9723-843a95fd7070" artifact_version_id="8f71b991-233e-444c-af04-8301d0e4b171" title="internal/db/models/product.go" contentType="text/go">
```go
package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	MerchantID  uint    `gorm:"not null" json:"merchant_id"`
	Name        string  `gorm:"size:255;not null" json:"name"`
	Description string  `gorm:"type:text" json:"description"`
	SKU         string  `gorm:"size:100;unique;not null" json:"sku"`
	Price       float64 `gorm:"type:decimal(10,2);not null" json:"price"`
	CategoryID  uint    `gorm:"not null" json:"category_id"`
	Merchant    Merchant `gorm:"foreignKey:MerchantID"`
	Category    Category `gorm:"foreignKey:CategoryID"`
}
```
</xaiArtifact>

---
### internal/db/models/order_item.go
- Size: 1.59 KB
- Lines: 53
- Last Modified: 2025-09-03 17:42:59

<xaiArtifact artifact_id="2d20a88a-525c-40d4-a32f-46eec6f3e9d7" artifact_version_id="93a86bbc-b17d-481c-ade9-db5135aa7793" title="internal/db/models/order_item.go" contentType="text/go">
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

type OrderItem struct {
	gorm.Model
	OrderID           uint              `gorm:"not null" json:"order_id"`
	ProductID         uint              `gorm:"not null" json:"product_id"`
	MerchantID        uint              `gorm:"not null" json:"merchant_id"`
	Quantity          int               `gorm:"not null" json:"quantity"`
	Price             float64           `gorm:"type:decimal(10,2);not null" json:"price"`
	FulfillmentStatus FulfillmentStatus `gorm:"type:varchar(20);not null;default:'New'" json:"fulfillment_status"`
	Order             Order             `gorm:"foreignKey:OrderID"`
	Product           Product           `gorm:"foreignKey:ProductID"`
	Merchant          Merchant          `gorm:"foreignKey:MerchantID"`
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
</xaiArtifact>

---
### internal/db/models/user.go
- Size: 0.54 KB
- Lines: 15
- Last Modified: 2025-09-12 21:58:48

<xaiArtifact artifact_id="54ba02ac-29a5-465b-9d87-4fede6dc75d5" artifact_version_id="b04445a3-aaea-4436-a288-da23069090f0" title="internal/db/models/user.go" contentType="text/go">
```go
package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"unique;not null"`
	Name     string `gorm:"type:varchar(100);not null"`
	Password string // Empty for OAuth users
	//Role     string `gorm:"not null"` // "customer" (default) or "merchant" (upgraded by admin)
	GoogleID string    // Google ID for OAuth
	Country  string `gorm:"type:varchar(100)"` // Optional country field
	//Carts    []Cart  `gorm:"foreignKey:UserID" json:"carts,omitempty"`
	//Orders   []Order `gorm:"foreignKey:UserID" json:"orders,omitempty"`
}
```
</xaiArtifact>

---
### internal/db/models/category.go
- Size: 0.34 KB
- Lines: 13
- Last Modified: 2025-09-03 17:38:43

<xaiArtifact artifact_id="d2572282-e623-4bd5-a02b-90fe087b6e7a" artifact_version_id="e9c3bed5-50ac-4ccf-9556-df72d61bde6b" title="internal/db/models/category.go" contentType="text/go">
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
</xaiArtifact>

---
### internal/db/models/order.go
- Size: 1.16 KB
- Lines: 49
- Last Modified: 2025-09-03 17:42:12

<xaiArtifact artifact_id="df07a79d-1dc5-4735-a68f-d2de58486764" artifact_version_id="e8549a66-2de8-4829-8f78-95c0d5e6f1de" title="internal/db/models/order.go" contentType="text/go">
```go
package models

import (
	"fmt"
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

type Order struct {
	gorm.Model
	UserID      uint        `gorm:"not null" json:"user_id"`
	TotalAmount float64     `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	Status      OrderStatus `gorm:"type:varchar(20);not null;default:'Pending'" json:"status"`
	User        User        `gorm:"foreignKey:UserID"`
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
</xaiArtifact>

---
### internal/tests/unit/test_handlers.go
- Size: 0.01 KB
- Lines: 1
- Last Modified: 2025-09-12 19:27:54

<xaiArtifact artifact_id="37dfdf57-4b1f-42db-8280-9a9c1b627eb4" artifact_version_id="46d0792f-47b1-48d7-9199-0874f754f32e" title="internal/tests/unit/test_handlers.go" contentType="text/go">
```go
package unit
```
</xaiArtifact>

---
### internal/tests/unit/test_service.go
- Size: 0.01 KB
- Lines: 1
- Last Modified: 2025-09-12 19:27:13

<xaiArtifact artifact_id="70523bd1-6747-4e90-9505-8d0c9fabb136" artifact_version_id="54c2b65d-dacc-4936-92bd-2ffcfbaaa211" title="internal/tests/unit/test_service.go" contentType="text/go">
```go
package unit

```
</xaiArtifact>

---
### internal/api/merchant/routes.go
- Size: 4.00 KB
- Lines: 112
- Last Modified: 2025-09-12 18:56:50

<xaiArtifact artifact_id="027f44bb-02a8-41a0-9bdc-a2ee4a6b0e64" artifact_version_id="2227639d-8fc2-468b-bfaa-daeef41ba4cf" title="internal/api/merchant/routes.go" contentType="text/go">
```go
package merchant
/*
    import (
       "api-customer-merchant/internal/api/merchant/handlers"
       "api-customer-merchant/internal/middleware"
       "api-customer-merchant/internal/db/repositories"

        "github.com/gin-gonic/gin"
    )

    func RegisterRoutes(r *gin.Engine) {
        merchant := r.Group("/merchant")
        {
            authHandler := handlers.NewAuthHandler()
            merchant.POST("/submitApplication", authHandler.Register)
            merchant.POST("/login", authHandler.Login)
            //merchant.GET("/auth/google", authHandler.GoogleAuth)
            //merchant.GET("/auth/google/callback", authHandler.GoogleCallback)

            protected := merchant.Group("/")
            protected.Use(middleware.AuthMiddleware("merchant"))
            protected.POST("/logout", authHandler.Logout)
        }
    }
*/


import (
	"api-customer-merchant/internal/api/merchant/handlers"
	//"api-customer-merchant/internal/db"
    //"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"

	//"api-customer-merchant/internal/middleware"
	"api-customer-merchant/internal/domain/merchant"
	// "api-customer-merchant/internal/domain/order"
	// "api-customer-merchant/internal/domain/payout"
	//"api-customer-merchant/internal/domain/product"
	//"api-customer-merchant/internal/domain/invent"

	"github.com/gin-gonic/gin"
)

/*
   func RegisterRoutes(r *gin.Engine) {
       merchant := r.Group("/merchant")
       {
          //userRepo := repositories.NewUserRepository()
           // merchantRepo := repositories.NewMerchantRepository()
           appRepo := repositories.NewMerchantApplicationRepository()
           repo := repositories.NewMerchantRepository()
           service := merchant.NewMerchantService(appRepo, repo)
           h := handlers.NewMerchantHandler(service)

           // authHandler := handlers.NewAuthHandler(merchantRepo)
           // merchant.POST("/submitApplication", authHandler.Register)
           // merchant.POST("/login", authHandler.Login)
           //merchant.GET("/auth/google", authHandler.GoogleAuth)
           //merchant.GET("/auth/google/callback", authHandler.GoogleCallback)
           merchant.POST("/apply",  h.Apply)
           merchant.GET("/application/:id",  h.GetApplication)


           // Merchant account access (once approved by admin via Express API)
           //merchant.POST("/me",  h.GetMyMerchant)



           merchant := merchant.Group("/")
           merchant.Use(middleware.AuthMiddleware("merchant"))
           merchant.POST("/me",  h.GetMyMerchant)
           //merchant.POST("/logout", authHandler.Logout)

       }
   }
*/

func RegisterRoutes(r *gin.Engine) {
    appRepo := repositories.NewMerchantApplicationRepository()
    repo := repositories.NewMerchantRepository()
    service := merchant.NewMerchantService(appRepo, repo)
    //productRepo := repositories.NewProductRepository()
    //inventoryRepo:= repositories.NewInventoryRepository()
	//productService := product.NewProductService(productRepo,inventoryRepo)

	// Other services (stubs; instantiate as needed)
	// orderService := order.NewOrderService(nil) // Adjust
	// payoutService := payout.NewPayoutService(nil)
	// promotionService := promotions.NewPromotionService(nil)
    //h := handlers.NewMerchantAuthHandler(service)
    
    authHandler := handlers.NewMerchantAuthHandler(service)
    //merchhandler:=handlers.NewMerchantHandlers(productService)
    merchant := r.Group("/merchant")
    {
    // Application submission and status
    merchant.POST("/apply",  authHandler.Apply)
    merchant.GET("/application/:id",  authHandler.GetApplication)
    // Merchant account access (once approved by admin via Express API)
    merchant.GET("/me",  authHandler.GetMyMerchant)

    //merchant.POST("/create/product", merchhandler.CreateProduct)
	//merchant.GET("/products", merchhandler.GetMyProducts)
	//merchant.PUT("/:id", merchhandler.UpdateProduct)
	//merchant.DELETE("/:id", merchhandler.DeleteProduct)
	//merchant.POST("/bulk-upload", merchhandler.BulkUploadProducts)


    }

   
}
```
</xaiArtifact>

---
### internal/api/customer/routes.go
- Size: 0.94 KB
- Lines: 28
- Last Modified: 2025-09-12 19:01:11

<xaiArtifact artifact_id="254b02a8-2214-4461-8328-8e8b41a00cc5" artifact_version_id="53a1a303-ee1f-4737-9965-2876f2cda0fc" title="internal/api/customer/routes.go" contentType="text/go">
```go
package customer

   import (
       "api-customer-merchant/internal/api/customer/handlers"
       "api-customer-merchant/internal/middleware"
       "api-customer-merchant/internal/db/repositories"
       "api-customer-merchant/internal/domain/user"


       "github.com/gin-gonic/gin"
   )

   func RegisterRoutes(r *gin.Engine) {
    repo := repositories.NewUserRepository()
    service := user.NewAuthService( repo)
       customer := r.Group("/customer")
       {
           authHandler := handlers.NewAuthHandler(service)
           customer.POST("/register", authHandler.Register)
           customer.POST("/login", authHandler.Login)
           customer.GET("/auth/google", authHandler.GoogleAuth)
           customer.GET("/auth/google/callback", authHandler.GoogleCallback)

           protected := customer.Group("/")
           protected.Use(middleware.AuthMiddleware("customer"))
           protected.POST("/logout", authHandler.Logout)
       }
   }
```
</xaiArtifact>

---
### internal/api/merchant/handlers/merchant_handlers.go
- Size: 3.52 KB
- Lines: 135
- Last Modified: 2025-09-04 13:43:00

<xaiArtifact artifact_id="4f60856d-cd8f-4115-a0ad-6ac5741ce492" artifact_version_id="d87961cc-6945-4571-a284-a5148c63a968" title="internal/api/merchant/handlers/merchant_handlers.go" contentType="text/go">
```go
package handlers

import (
	"net/http"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/domain/order"
	"api-customer-merchant/internal/domain/payout"
	"api-customer-merchant/internal/domain/product"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MerchantHandlers struct {
	productService *product.ProductService
	orderService   *order.OrderService
	payoutService  *payout.PayoutService
}

func NewMerchantHandlers(productService *product.ProductService, orderService *order.OrderService, payoutService *payout.PayoutService) *MerchantHandlers {
	return &MerchantHandlers{
		productService: productService,
		orderService:   orderService,
		payoutService:  payoutService,
	}
}

// CreateProduct handles POST /merchant/products
func (h *MerchantHandlers) CreateProduct(c *gin.Context) {
	merchantID, exists := c.Get("merchantID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.productService.CreateProduct(&product, merchantID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// UpdateProduct handles PUT /merchant/products/:id
func (h *MerchantHandlers) UpdateProduct(c *gin.Context) {
	merchantID, exists := c.Get("merchantID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product.ID = uint(id)

	if err := h.productService.UpdateProduct(&product, merchantID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct handles DELETE /merchant/products/:id
func (h *MerchantHandlers) DeleteProduct(c *gin.Context) {
	merchantID, exists := c.Get("merchantID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	if err := h.productService.DeleteProduct(uint(id), merchantID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted"})
}

// GetOrders handles GET /merchant/orders
func (h *MerchantHandlers) GetOrders(c *gin.Context) {
	merchantID, exists := c.Get("merchantID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	orders, err := h.orderService.GetOrdersByMerchantID(merchantID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetPayouts handles GET /merchant/payouts
func (h *MerchantHandlers) GetPayouts(c *gin.Context) {
	merchantID, exists := c.Get("merchantID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	payouts, err := h.payoutService.GetPayoutsByMerchantID(merchantID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payouts)
}
```
</xaiArtifact>

---
### internal/api/merchant/handlers/auth_handler.go
- Size: 5.45 KB
- Lines: 114
- Last Modified: 2025-09-13 02:56:54

<xaiArtifact artifact_id="59637335-eba0-437a-a2f3-d45b07f21406" artifact_version_id="94ec4f3e-01a3-4697-a681-0a984b8e2640" title="internal/api/merchant/handlers/auth_handler.go" contentType="text/go">
```go
package handlers

import (
	"encoding/json"
	"net/http"

	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/domain/merchant"

	"github.com/gin-gonic/gin"
)

type MerchantHandler struct {
	service *merchant.MerchantService
}

func NewMerchantAuthHandler(s *merchant.MerchantService) *MerchantHandler {
	return &MerchantHandler{service: s}
}

// Apply godoc
// @Summary Submit a new merchant application
// @Description Allows a prospective merchant to submit an application with personal, business, and address information
// @Tags Merchant
// @Accept json
// @Produce json
// @Param body body object{first_name=string,last_name=string,email=string,phone=string,personal_address=object{street=string,city=string,state=string,postal_code=string,country=string},work_address=object{street=string,city=string,state=string,postal_code=string,country=string},business_name=string,business_type=string,tax_id=string,documents=object{business_license=string,identification=string}} true "Merchant application details"
// @Success 201 {object} object{id=string,first_name=string,last_name=string,email=string,phone=string,personal_address=object{street=string,city=string,state=string,postal_code=string,country=string},work_address=object{street=string,city=string,state=string,postal_code=string,country=string},business_name=string,business_type=string,tax_id=string,documents=object{business_license=string,identification=string},status=string,created_at=string} "Created application"
// @Failure 400 {object} object{error=string} "Invalid request body or malformed JSON"
// @Failure 500 {object} object{error=string} "Failed to submit application"
// @Router /merchant/apply [post]
func (h *MerchantHandler) Apply(c *gin.Context) {
	var req struct {
		models.MerchantBasicInfo
		PersonalAddress            map[string]any `json:"personal_address" validate:"required"`
		WorkAddress                map[string]any `json:"work_address" validate:"required"`
		models.MerchantBusinessInfo
		models.MerchantDocuments
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	// Convert personal_address and work_address to JSONB
	personalAddressJSON, err := json.Marshal(req.PersonalAddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid personal_address JSON: " + err.Error()})
		return
	}
	workAddressJSON, err := json.Marshal(req.WorkAddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid work_address JSON: " + err.Error()})
		return
	}

	app := &models.MerchantApplication{
		MerchantBasicInfo:    req.MerchantBasicInfo,
		MerchantAddress:      models.MerchantAddress{PersonalAddress: personalAddressJSON, WorkAddress: workAddressJSON},
		MerchantBusinessInfo: req.MerchantBusinessInfo,
		MerchantDocuments:    req.MerchantDocuments,
	}

	app, err = h.service.SubmitApplication(c.Request.Context(), app)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to submit application: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, app)
}

// GetApplication godoc
// @Summary Retrieve a merchant application by ID
// @Description Allows an applicant to view the status and details of their merchant application
// @Tags Merchant
// @Produce json
// @Param id path string true "Application ID" format(uuid)
// @Success 200 {object} object{id=string,first_name=string,last_name=string,email=string,phone=string,personal_address=object{street=string,city=string,state=string,postal_code=string,country=string},work_address=object{street=string,city=string,state=string,postal_code=string,country=string},business_name=string,business_type=string,tax_id=string,documents=object{business_license=string,identification=string},status=string,created_at=string} "Application details"
// @Failure 400 {object} object{error=string} "Invalid application ID format"
// @Failure 404 {object} object{error=string} "Application not found"
// @Router /merchant/application/{id} [get]
func (h *MerchantHandler) GetApplication(c *gin.Context) {
	id := c.Param("id")
	app, err := h.service.GetApplication(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, app)
}

// GetMyMerchant godoc
// @Summary Retrieve current merchant account
// @Description Fetches the merchant account details for the authenticated user, if their application has been approved
// @Tags Merchant
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{id=string,user_id=string,business_name=string,business_type=string,tax_id=string,personal_address=object{street=string,city=string,state=string,postal_code=string,country=string},work_address=object{street=string,city=string,state=string,postal_code=string,country=string},status=string,created_at=string,updated_at=string} "Merchant account details"
// @Failure 401 {object} object{error=string} "Unauthorized: Missing or invalid authentication"
// @Failure 404 {object} object{error=string} "Merchant account not found"
// @Router /merchant/me [get]
func (h *MerchantHandler) GetMyMerchant(c *gin.Context) {
	userID, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	m, err := h.service.GetMerchantByUserID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "merchant not found: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, m)
}
```
</xaiArtifact>

---
### internal/api/customer/handlers/customer_handlers.go
- Size: 9.13 KB
- Lines: 334
- Last Modified: 2025-09-05 12:10:28

<xaiArtifact artifact_id="5ee7069d-4731-4c43-8996-558d7175eb74" artifact_version_id="7bf09ba4-1e17-4c61-9e23-c4bca509a7a7" title="internal/api/customer/handlers/customer_handlers.go" contentType="text/go">
```go
package handlers

import (
	"net/http"
	//"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/domain/cart"
	"api-customer-merchant/internal/domain/order"
	"api-customer-merchant/internal/domain/payment"
	"api-customer-merchant/internal/domain/product"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CustomerHandlers struct {
	productService *product.ProductService
	cartService    *cart.CartService
	orderService   *order.OrderService
	paymentService *payment.PaymentService
}

func NewCustomerHandlers(productService *product.ProductService, cartService *cart.CartService, orderService *order.OrderService, paymentService *payment.PaymentService) *CustomerHandlers {
	return &CustomerHandlers{
		productService: productService,
		cartService:    cartService,
		orderService:   orderService,
		paymentService: paymentService,
	}
}

// GetProductByID handles GET /customer/products/:id
func (h *CustomerHandlers) GetProductByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	product, err := h.productService.GetProductByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// GetProductBySKU handles GET /customer/products/sku/:sku
func (h *CustomerHandlers) GetProductBySKU(c *gin.Context) {
	sku := c.Param("sku")
	if sku == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SKU cannot be empty"})
		return
	}

	product, err := h.productService.GetProductBySKU(sku)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// SearchProductsByName handles GET /customer/products/search?name={name}
func (h *CustomerHandlers) SearchProductsByName(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "search name cannot be empty"})
		return
	}

	products, err := h.productService.SearchProductsByName(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

// AddItemToCart handles POST /customer/cart/add
func (h *CustomerHandlers) AddItemToCart(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.cartService.AddItemToCart(userID.(uint), req.ProductID, req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// UpdateCartItemQuantity handles PUT /customer/cart/update/:cartItemID
func (h *CustomerHandlers) UpdateCartItemQuantity(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	cartItemIDStr := c.Param("cartItemID")
	cartItemID, err := strconv.ParseUint(cartItemIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cart item ID"})
		return
	}
	// Verify cart item belongs to user's active cart
	cartItem, err := h.cartService.GetCartItemByID(uint(cartItemID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cart item not found"})
		return
	}
	cart, err := h.cartService.GetActiveCart(userID.(uint))
	if err != nil || cart.ID != cartItem.CartID {
		c.JSON(http.StatusForbidden, gin.H{"error": "cart item does not belong to user"})
		return
	}

	var req struct {
		Quantity int `json:"quantity" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedCart, err := h.cartService.UpdateCartItemQuantity(uint(cartItemID), req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedCart)
}

// RemoveCartItem handles DELETE /customer/cart/remove/:cartItemID
func (h *CustomerHandlers) RemoveCartItem(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	cartItemIDStr := c.Param("cartItemID")
	cartItemID, err := strconv.ParseUint(cartItemIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cart item ID"})
		return
	}
	// Verify cart item belongs to user's active cart
	cartItem, err := h.cartService.GetCartItemByID(uint(cartItemID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cart item not found"})
		return
	}
	cart, err := h.cartService.GetActiveCart(userID.(uint))
	if err != nil || cart.ID != cartItem.CartID {
		c.JSON(http.StatusForbidden, gin.H{"error": "cart item does not belong to user"})
		return
	}
	

	updatedCart, err := h.cartService.RemoveCartItem(uint(cartItemID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedCart)
}

// CreateOrder handles POST /customer/orders/create
func (h *CustomerHandlers) CreateOrder(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	order, err := h.orderService.CreateOrderFromCart(userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetOrderByID handles GET /customer/orders/:id
func (h *CustomerHandlers) GetOrderByID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	order, err := h.orderService.GetOrderByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	if order.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "order does not belong to user"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetOrders handles GET /customer/orders
func (h *CustomerHandlers) GetOrders(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	orders, err := h.orderService.GetOrdersByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// ProcessPayment handles POST /customer/payments/process
func (h *CustomerHandlers) ProcessPayment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		OrderID uint    `json:"order_id" binding:"required"`
		Amount  float64 `json:"amount" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify order belongs to user
	order, err := h.orderService.GetOrderByID(req.OrderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	if order.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "order does not belong to user"})
		return
	}

	payment, err := h.paymentService.ProcessPayment(req.OrderID, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// GetPaymentByOrderID handles GET /customer/payments/:orderID
func (h *CustomerHandlers) GetPaymentByOrderID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	orderIDStr := c.Param("orderID")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	// Verify order belongs to user
	order, err := h.orderService.GetOrderByID(uint(orderID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	if order.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "order does not belong to user"})
		return
	}

	payment, err := h.paymentService.GetPaymentByOrderID(uint(orderID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		return
	}

	c.JSON(http.StatusOK, payment)
}

func (h *CustomerHandlers) GetProductsByCategory(c *gin.Context) {
    categoryIDStr := c.Param("categoryID")
    categoryID, _ := strconv.ParseUint(categoryIDStr, 10, 32)
    limit, _ := strconv.Atoi(c.Query("limit"))
    offset, _ := strconv.Atoi(c.Query("offset"))
    products, err := h.productService.GetProductsByCategory(uint(categoryID), limit, offset)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, products)
}
```
</xaiArtifact>

---
### internal/api/customer/handlers/auth_handler.go
- Size: 6.12 KB
- Lines: 198
- Last Modified: 2025-09-12 19:33:09

<xaiArtifact artifact_id="0f84d8ec-c40f-400f-9e84-f224e3b61e0f" artifact_version_id="4308a009-c3d6-4b21-a655-889e088b6103" title="internal/api/customer/handlers/auth_handler.go" contentType="text/go">
```go
package handlers

import (
	"log"
	"net/http"
	"os"

	//"os"
	"strings"

	//"api-customer-merchant/internal/db/models"
	//"api-customer-merchant/internal/db/repositories"
	services "api-customer-merchant/internal/domain/user"
	"api-customer-merchant/internal/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	service *services.AuthService
}

// In customer/handlers/auth_handler.go AND merchant/handlers/auth_handler.go
func NewAuthHandler(s *services.AuthService) *AuthHandler {
    return &AuthHandler{
        service: s,
    }
}
// Register godoc
// @Summary Register a new customer
// @Description Creates a new customer account with email, name, password, and optional country
// @Tags Customer
// @Accept json
// @Produce json
// @Param body body object{email=string,name=string,password=string,country=string} true "Customer registration details"
// @Success 200 {object} object{token=string} "JWT token"
// @Failure 400 {object} object{error=string} "Invalid request"
// @Failure 500 {object} object{error=string} "Server error"
// @Router /customer/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
		Country  string `json:"country"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.RegisterUser(req.Email, req.Name, req.Password, req.Country)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token})
}


// Login godoc
// @Summary Customer login
// @Description Authenticates a customer using email and password
// @Tags Customer
// @Accept json
// @Produce json
// @Param body body object{email=string,password=string} true "Customer login credentials"
// @Success 200 {object} object{token=string} "JWT token"
// @Failure 400 {object} object{error=string} "Invalid request"
// @Failure 401 {object} object{error=string} "Unauthorized"
// @Failure 403 {object} object{error=string} "Invalid role"
// @Failure 500 {object} object{error=string} "Server error"
// @Router /customer/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.LoginUser(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// if user.Role != "customer" {
	// 	c.JSON(http.StatusForbidden, gin.H{"error": "Invalid role for this API"})
	// 	return
	// }

	token, err := h.service.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}


// GoogleAuth godoc
// @Summary Initiate Google OAuth for customer
// @Description Redirects to Google OAuth login page
// @Tags Customer
// @Produce json
// @Success 307 {object} object{} "Redirect to Google OAuth"
// @Router /customer/auth/google [get]
func (h *AuthHandler) GoogleAuth(c *gin.Context) {
	url := h.service.GetOAuthConfig("customer").AuthCodeURL("state-customer", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
	 //c.JSON(http.StatusOK, gin.H{"url": url})
}



// GoogleCallback godoc
// @Summary Handle Google OAuth callback for customer
// @Description Processes Google OAuth callback and returns JWT token
// @Tags Customer
// @Produce json
// @Param code query string true "OAuth code"
// @Success 200 {object} object{token=string} "JWT token"
// @Failure 400 {object} object{error=string} "Code not provided"
// @Failure 500 {object} object{error=string} "Server error"
// @Router /customer/auth/google/callback [get]
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
    code := c.Query("code")
    state := c.Query("state")
    if code == "" || state == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code or state"})
        return
    }
    // Verify state (in production, check against stored value)
    if state != "state-customer" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state parameter"})
        return
    }

    user, token, err := h.service.GoogleLogin(code, os.Getenv("BASE_URL"), "customer")
    if err != nil {
        log.Printf("Google login failed: %v", err)
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"token": token, "user": user})
}



// Logout godoc
// @Summary Customer logout
// @Description Invalidates the customer's JWT token
// @Tags Customer
// @Security BearerAuth
// @Produce json
// @Success 200 {object} object{message=string} "Logout successful"
// @Failure 400 {object} object{error=string} "Authorization header required"
// @Router /customer/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header required"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	utils.Add(tokenString)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
    userID, _ := c.Get("userID")
    var req struct { Name string; Country string; Addresses []string }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := h.service.UpdateProfile(userID.(uint), req.Name, req.Country, req.Addresses); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "updated"})
}
```
</xaiArtifact>

---

---
## 📊 Summary
- Total files: 48
- Total size: 107.47 KB
- File types: .go, unknown
