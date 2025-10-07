package main

import (
	"fmt"
	"log"

	//"os"

	//"time"

	"api-customer-merchant/internal/db/models" // Adjust to your models' package path

	//"github.com/google/uuid"
	//"github.com/shopspring/decimal"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	//"github.com/google/uuid"
	//"golang.org/x/crypto/bcrypt"
)

func main() {
	// Connect to DB
	// dsn := os.Getenv("DB_DSN")
	// if dsn == "" {
	//     log.Fatal("DB_DSN environment variable not set")
	// }
	dsn := "postgresql://neondb_owner:npg_CcwoeLb6V1XH@ep-wild-haze-adu0bdvq-pooler.c-2.us-east-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require"

	db, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn}), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// Run migrations first (optional, for dev)
	//db.AutoMigrate(  &models.Product{},&models.Variant{},models.Media{},models.Category{})

	// Seed Users
	// hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	// user1 := models.User{
	// 	Email:    "user1@example.com",
	// 	Name:     "John Doe",
	// 	Password: string(hashedPassword),
	// 	Country:  "US",
	// }
	// user2 := models.User{
	// 	Email:    "merchant1@example.com",
	// 	Name:     "Jane Merchant",
	// 	Password: string(hashedPassword),
	// 	Country:  "UK",
	// }
	// if err := db.Create([]*models.User{&user1, &user2}).Error; err != nil {
	// 	log.Fatalf("Failed to seed users: %v", err)
	// }

	// Seed Merchant

	// Seed Categories
	//  category := models.Category{
	//  	Name:       "Electronics",
	//  	Attributes: map[string]interface{}{"type": "gadget"},
	//  }
	//  if err := db.Create(&category).Error; err != nil {
	//  	log.Fatalf("Failed to seed category: %v", err)
	//  }

	// Seed Products
	/*
	   	 product1 := models.Product{
	   	 	ID:         uuid.New().String(),
	   	 	MerchantID: "68a63ffc-f988-47a3-bc74-989b498b1e01",
	   	 	CategoryID: 1,
	   	 	Name:       "Smartphone",
	   	 	SKU:        "SM-001",
	   	 	BasePrice:      699.99,
	   	 	Media:      []models.Media{{URL: "image1.png"},{URL: "image2.png"}},
	   		Variants: []models.Variant{
	           {
	               SKU:   "SNKR-001-BLK-42",
	               Price: 79.99,
	               Attributes: map[string]string{
	                   "color": "black",
	                   "size":  "42",
	               },
	           },
	           {
	               SKU:   "SNKR-001-WHT-43",
	               Price: 84.99,
	               Attributes: map[string]string{
	                   "color": "white",
	                   "size":  "43",
	               },
	           },
	       },
	   	 }
	*/
	//  product2 := models.Product{
	//  	ID:         uuid.New().String(),
	//  	MerchantID: merchantID,
	//  	CategoryID: category.ID,
	//  	Name:       "Laptop",
	//  	SKU:        "LP-001",
	//  	Price:      1299.99,
	//  	Currency:   "USD",
	//  	IsActive:   true,
	//  }

	// product2 := models.Product{
	// 	ID:          uuid.New().String(),
	// 	MerchantID:  "984d6da6-29c4-4506-abaf-608b3498cc04",
	// 	CategoryID:  1,
	// 	Name:        "Earpod",
	// 	Description: "A high-end earpod with advanced features",
	// 	SKU:         "EAR-001",
	// 	BasePrice:   decimal.NewFromFloat(700.00),
	// 	Media: []models.Media{
	// 		{
	// 			URL:  "https://www.example.com/image4.png",
	// 			Type: "",
	// 		},
	// 		{
	// 			URL:  "https://www.example.com/image5.png",
	// 			Type: "",
	// 		},
	// 	},
	// 	Variants: []models.Variant{
	// 		{
	// 			ProductID:       "", // Will be set automatically after product creation
	// 			SKU:             "EAR-01-BLK-64",
	// 			PriceAdjustment: decimal.NewFromFloat(50.00),
	// 			TotalPrice:      decimal.NewFromFloat(0.00), // Will be computed in BeforeCreate
	// 			Attributes: models.AttributesMap{
	// 				"color": "black",
	// 				"size":  "large",
	// 			},
	// 			IsActive: true,
	// 			Inventory: models.Inventory{
	// 				Quantity:          100,
	// 				ReservedQuantity:  0,
	// 				LowStockThreshold: 10,
	// 				BackorderAllowed:  false,
	// 				MerchantID:        "984d6da6-29c4-4506-abaf-608b3498cc04",
	// 			},
	// 		},
	// 		{
	// 			ProductID:       "", // Will be set automatically after product creation
	// 			SKU:             "EAR-001-WHT-128",
	// 			PriceAdjustment: decimal.NewFromFloat(50.00),
	// 			TotalPrice:      decimal.NewFromFloat(0.00), // Will be computed in BeforeCreate
	// 			Attributes: models.AttributesMap{
	// 				"color": "white",
	// 				"size":  "medium",
	// 			},
	// 			IsActive: true,
	// 			Inventory: models.Inventory{
	// 				Quantity:          50,
	// 				ReservedQuantity:  0,
	// 				LowStockThreshold: 10,
	// 				BackorderAllowed:  false,
	// 				MerchantID:        "984d6da6-29c4-4506-abaf-608b3498cc04",
	// 			},
	// 		},
	// 	},
	// }
	// if err := db.Create([]*models.Product{&product2}).Error; err != nil {
	// 	log.Fatalf("Failed to seed products: %v", err)
	// }

	// // Seed Inventory
	// inventory1 := models.Inventory{
	// 	ProductID:     product1.ID,
	// 	StockQuantity: 100,
	// 	LowThreshold:  10,
	// }
	// inventory2 := models.Inventory{
	// 	ProductID:     product2.ID,
	// 	StockQuantity: 50,
	// 	LowThreshold:  5,
	// }
	// if err := db.Create([]*models.Inventory{&inventory1, &inventory2}).Error; err != nil {
	// 	log.Fatalf("Failed to seed inventory: %v", err)
	// }

	// // Seed Cart and CartItem
	// cart := models.Cart{
	// 	UserID: user1.ID,
	// 	Status: models.CartStatusActive,
	// }
	// if err := db.Create(&cart).Error; err != nil {
	// 	log.Fatalf("Failed to seed cart: %v", err)
	// }
	// cartItem := models.CartItem{
	// 	CartID:        cart.ID,
	// 	ProductID:     product1.ID,
	// 	Quantity:      2,
	// 	PriceSnapshot: product1.Price,
	// 	MerchantID:    merchantID,
	// }
	// if err := db.Create(&cartItem).Error; err != nil {
	// 	log.Fatalf("Failed to seed cart item: %v", err)
	// }


	//Seed Categories
	// category := models.Category{
	// 	Name:       "Electronics",
	// 	Attributes: map[string]interface{}{"type": "gadget"},
	// }
	// if err := db.Create(&category).Error; err != nil {
	// 	log.Fatalf("Failed to seed category Electronics: %v", err)
	// }

	// Luxury Clothing Category 1: Haute Couture
	category := models.Category{
		Name:       "Haute Couture",
		Attributes: map[string]interface{}{
			"material":  "silk",
			"style":     "evening wear",
			"exclusivity": "limited edition",
		},
	}
	if err := db.Create(&category).Error; err != nil {
		log.Fatalf("Failed to seed category Haute Couture: %v", err)
	}

	// Luxury Clothing Category 2: Designer Footwear
	category = models.Category{
		Name:       "Designer Footwear",
		Attributes: map[string]interface{}{
			"material":  "leather",
			"heel_type": "stiletto",
			"brand_tier": "premium",
		},
	}
	if err := db.Create(&category).Error; err != nil {
		log.Fatalf("Failed to seed category Designer Footwear: %v", err)
	}

	// Luxury Clothing Category 3: Luxury Accessories
	category = models.Category{
		Name:       "Luxury Accessories",
		Attributes: map[string]interface{}{
			"type":      "handbags",
			"material":  "exotic leather",
			"hardware":  "gold-plated",
		},
	}
	if err := db.Create(&category).Error; err != nil {
		log.Fatalf("Failed to seed category Luxury Accessories: %v", err)
	}

	// Luxury Clothing Category 4: Bespoke Tailoring
	category = models.Category{
		Name:       "Bespoke Tailoring",
		Attributes: map[string]interface{}{
			"fit":       "custom",
			"fabric":    "wool",
			"occasion":  "formal",
		},
	}
	if err := db.Create(&category).Error; err != nil {
		log.Fatalf("Failed to seed category Bespoke Tailoring: %v", err)
	}

	// Luxury Clothing Category 5: Premium Outerwear
	category = models.Category{
		Name:       "Premium Outerwear",
		Attributes: map[string]interface{}{
			"material":  "cashmere",
			"insulation": "down",
			"style":     "trench",
		},
	}
	if err := db.Create(&category).Error; err != nil {
		log.Fatalf("Failed to seed category Premium Outerwear: %v", err)
	}

	log.Println("Successfully seeded all categories")




	fmt.Println("Database seeded successfully!")
}
