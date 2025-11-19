package main

import (
	"fmt"
	"log"

	 "api-customer-merchant/internal/db/models"

	// "github.com/google/uuid"
	// "github.com/shopspring/decimal"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Connect to DB
	dsn := "postgresql://neondb_owner:npg_CcwoeLb6V1XH@ep-wild-haze-adu0bdvq-pooler.c-2.us-east-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require"

	db, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn}), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// Run migrations first (optional, for dev)
	 //db.AutoMigrate(&models.Category{}, &models.Product{}, &models.Variant{}, &models.Media{}, &models.Inventory{})
	 db.AutoMigrate(&models.Settings{})
	// Merchant IDs to use
	// merchantID1 := "68a63ffc-f988-47a3-bc74-989b498b1e01"
	// merchantID2 := "984d6da6-29c4-4506-abaf-608b3498cc04"

	// // Seed Categories
	// categories := []models.Category{
	// 	{
	// 		Name:         "Clothes",
	// 		CategorySlug: "clothes",
	// 		Attributes:   map[string]interface{}{"type": "apparel", "gender": "unisex"},
	// 	},
	// 	{
	// 		Name:         "Watches",
	// 		CategorySlug: "watches",
	// 		Attributes:   map[string]interface{}{"type": "accessory", "material": "metal"},
	// 	},
	// 	{
	// 		Name:         "Footwears",
	// 		CategorySlug: "footwears",
	// 		Attributes:   map[string]interface{}{"type": "footwear", "material": "leather"},
	// 	},
	// 	{
	// 		Name:         "Ankaras",
	// 		CategorySlug: "ankaras",
	// 		Attributes:   map[string]interface{}{"type": "fabric", "origin": "african"},
	// 	},
	// 	{
	// 		Name:         "Neckwears",
	// 		CategorySlug: "neckwears",
	// 		Attributes:   map[string]interface{}{"type": "accessory", "material": "fabric"},
	// 	},
	// 	{
	// 		Name:         "Bags",
	// 		CategorySlug: "bags",
	// 		Attributes:   map[string]interface{}{"type": "accessory", "material": "leather"},
	// 	},
	// }

	// // Create categories
	// for i := range categories {
	// 	if err := db.Create(&categories[i]).Error; err != nil {
	// 		log.Printf("Warning: Failed to seed category %s: %v", categories[i].Name, err)
	// 	}
	// }

	// // Get category IDs after creation
	// var clothesCategory, watchesCategory, footwearsCategory, ankarasCategory, neckwearsCategory, bagsCategory models.Category
	// db.Where("category_slug = ?", "clothes").First(&clothesCategory)
	// db.Where("category_slug = ?", "watches").First(&watchesCategory)
	// db.Where("category_slug = ?", "footwears").First(&footwearsCategory)
	// db.Where("category_slug = ?", "ankaras").First(&ankarasCategory)
	// db.Where("category_slug = ?", "neckwears").First(&neckwearsCategory)
	// db.Where("category_slug = ?", "bags").First(&bagsCategory)

	// // Seed Products with Variants and Media
	// products := []models.Product{
	// 	// Clothes Products
	// 	{
	// 		ID:          uuid.New().String(),
	// 		MerchantID:  merchantID1,
	// 		CategoryID:  clothesCategory.ID,
	// 		Name:        "Children's Handmade Clothes",
	// 		Description: "Beautiful children's clothes sewn with natural fabrics",
	// 		SKU:         "CLOTH-CHILD-001",
	// 		BasePrice:   decimal.NewFromFloat(45.99),
	// 		Media: []models.Media{
	// 			{URL: "https://media.istockphoto.com/id/1370454967/photo/clothes-for-children-are-sewn-with-their-own-hands-sale-of-clothes-made-of-natural-fabrics.jpg?s=612x612&w=0&k=20&c=Ghwz3LHzujIR3PiNWbpwSZLb_IQ_Mm06QOqfzi3lyQw=", Type: models.MediaTypeImage},
	// 		},
	// 		Variants: []models.Variant{
	// 			{
	// 				SKU:             "CLOTH-CHILD-001-S",
	// 				PriceAdjustment: decimal.NewFromFloat(0.00),
	// 				Attributes:      models.AttributesMap{"size": "S", "color": "multi"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          30,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 5,
	// 					BackorderAllowed:  true,
	// 					MerchantID:        merchantID1,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "CLOTH-CHILD-001-M",
	// 				PriceAdjustment: decimal.NewFromFloat(5.00),
	// 				Attributes:      models.AttributesMap{"size": "M", "color": "multi"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          25,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 5,
	// 					BackorderAllowed:  true,
	// 					MerchantID:        merchantID1,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "CLOTH-CHILD-001-L",
	// 				PriceAdjustment: decimal.NewFromFloat(10.00),
	// 				Attributes:      models.AttributesMap{"size": "L", "color": "multi"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          20,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 3,
	// 					BackorderAllowed:  false,
	// 					MerchantID:        merchantID1,
	// 				},
	// 			},
	// 		},
	// 	},
	// 	{
	// 		ID:          uuid.New().String(),
	// 		MerchantID:  merchantID2,
	// 		CategoryID:  clothesCategory.ID,
	// 		Name:        "Wedding Dress",
	// 		Description: "Beautiful wedding dress for the perfect bride",
	// 		SKU:         "CLOTH-WED-001",
	// 		BasePrice:   decimal.NewFromFloat(350.00),
	// 		Media: []models.Media{
	// 			{URL: "https://www.shutterstock.com/image-photo/beautiful-wedding-dresses-bridal-dress-260nw-2673306709.jpg", Type: models.MediaTypeImage},
	// 		},
	// 		Variants: []models.Variant{
	// 			{
	// 				SKU:             "CLOTH-WED-001-S",
	// 				PriceAdjustment: decimal.NewFromFloat(-20.00),
	// 				Attributes:      models.AttributesMap{"size": "S", "color": "white"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          5,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 1,
	// 					BackorderAllowed:  false,
	// 					MerchantID:        merchantID2,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "CLOTH-WED-001-M",
	// 				PriceAdjustment: decimal.NewFromFloat(0.00),
	// 				Attributes:      models.AttributesMap{"size": "M", "color": "white"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          3,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 1,
	// 					BackorderAllowed:  false,
	// 					MerchantID:        merchantID2,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "CLOTH-WED-001-L",
	// 				PriceAdjustment: decimal.NewFromFloat(20.00),
	// 				Attributes:      models.AttributesMap{"size": "L", "color": "white"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          2,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 1,
	// 					BackorderAllowed:  false,
	// 					MerchantID:        merchantID2,
	// 				},
	// 			},
	// 		},
	// 	},
	// 	{
	// 		ID:          uuid.New().String(),
	// 		MerchantID:  merchantID1,
	// 		CategoryID:  clothesCategory.ID,
	// 		Name:        "Traditional Qatari Dress",
	// 		Description: "Authentic traditional dress from Qatar",
	// 		SKU:         "CLOTH-QAT-001",
	// 		BasePrice:   decimal.NewFromFloat(120.00),
	// 		Media: []models.Media{
	// 			{URL: "https://www.shutterstock.com/image-photo/doha-qatar-january-05-2022-260nw-2247009105.jpg", Type: models.MediaTypeImage},
	// 		},
	// 		Variants: []models.Variant{
	// 			{
	// 				SKU:             "CLOTH-QAT-001-S",
	// 				PriceAdjustment: decimal.NewFromFloat(0.00),
	// 				Attributes:      models.AttributesMap{"size": "S", "color": "blue"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          15,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 3,
	// 					BackorderAllowed:  true,
	// 					MerchantID:        merchantID1,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "CLOTH-QAT-001-M",
	// 				PriceAdjustment: decimal.NewFromFloat(0.00),
	// 				Attributes:      models.AttributesMap{"size": "M", "color": "blue"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          12,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 2,
	// 					BackorderAllowed:  true,
	// 					MerchantID:        merchantID1,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "CLOTH-QAT-001-L",
	// 				PriceAdjustment: decimal.NewFromFloat(0.00),
	// 				Attributes:      models.AttributesMap{"size": "L", "color": "blue"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          10,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 2,
	// 					BackorderAllowed:  false,
	// 					MerchantID:        merchantID1,
	// 				},
	// 			},
	// 		},
	// 	},

	// 	// Watches Products
	// 	{
	// 		ID:          uuid.New().String(),
	// 		MerchantID:  merchantID2,
	// 		CategoryID:  watchesCategory.ID,
	// 		Name:        "Luxury Gold Women's Watch",
	// 		Description: "Elegant luxury watch for women with gold finish",
	// 		SKU:         "WATCH-GOLD-001",
	// 		BasePrice:   decimal.NewFromFloat(450.00),
	// 		Media: []models.Media{
	// 			{URL: "https://media.istockphoto.com/id/1180453576/photo/luxury-watch-isolated-on-white-background-with-clipping-path-gold-watch-women-watch-female.jpg?s=612x612&w=0&k=20&c=7156SpeDaeLHq7506ULnp6ZQrzbuoaHvOfnK6RT4L2A=", Type: models.MediaTypeImage},
	// 		},
	// 		Variants: []models.Variant{
	// 			{
	// 				SKU:             "WATCH-GOLD-001-S",
	// 				PriceAdjustment: decimal.NewFromFloat(-50.00),
	// 				Attributes:      models.AttributesMap{"size": "small", "material": "gold"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          8,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 2,
	// 					BackorderAllowed:  false,
	// 					MerchantID:        merchantID2,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "WATCH-GOLD-001-M",
	// 				PriceAdjustment: decimal.NewFromFloat(0.00),
	// 				Attributes:      models.AttributesMap{"size": "medium", "material": "gold"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          6,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 1,
	// 					BackorderAllowed:  false,
	// 					MerchantID:        merchantID2,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "WATCH-GOLD-001-L",
	// 				PriceAdjustment: decimal.NewFromFloat(50.00),
	// 				Attributes:      models.AttributesMap{"size": "large", "material": "gold"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          4,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 1,
	// 					BackorderAllowed:  false,
	// 					MerchantID:        merchantID2,
	// 				},
	// 			},
	// 		},
	// 	},
	// 	{
	// 		ID:          uuid.New().String(),
	// 		MerchantID:  merchantID1,
	// 		CategoryID:  watchesCategory.ID,
	// 		Name:        "Luxury White Dial Watch",
	// 		Description: "Stylish luxury watch with white dial",
	// 		SKU:         "WATCH-WHITE-001",
	// 		BasePrice:   decimal.NewFromFloat(380.00),
	// 		Media: []models.Media{
	// 			{URL: "https://media.istockphoto.com/id/1193931855/photo/luxury-watch-isolated-on-white-background-with-clipping-path-for-artwork-or-design-white.jpg?s=612x612&w=0&k=20&c=vw1ceQ7rq04cCkOvzqaywVwP34fLs0QvdI0pp8-elkM=", Type: models.MediaTypeImage},
	// 		},
	// 		Variants: []models.Variant{
	// 			{
	// 				SKU:             "WATCH-WHITE-001-S",
	// 				PriceAdjustment: decimal.NewFromFloat(-30.00),
	// 				Attributes:      models.AttributesMap{"size": "small", "material": "silver"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          10,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 2,
	// 					BackorderAllowed:  true,
	// 					MerchantID:        merchantID1,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "WATCH-WHITE-001-M",
	// 				PriceAdjustment: decimal.NewFromFloat(0.00),
	// 				Attributes:      models.AttributesMap{"size": "medium", "material": "silver"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          7,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 1,
	// 					BackorderAllowed:  false,
	// 					MerchantID:        merchantID1,
	// 				},
	// 			},
	// 		},
	// 	},
	// 	{
	// 		ID:          uuid.New().String(),
	// 		MerchantID:  merchantID2,
	// 		CategoryID:  watchesCategory.ID,
	// 		Name:        "Luxury Silver Watch",
	// 		Description: "Premium luxury silver watch",
	// 		SKU:         "WATCH-SILVER-001",
	// 		BasePrice:   decimal.NewFromFloat(520.00),
	// 		Media: []models.Media{
	// 			{URL: "https://www.shutterstock.com/image-photo/luxury-watch-isolated-on-white-260nw-2198958671.jpg", Type: models.MediaTypeImage},
	// 		},
	// 		Variants: []models.Variant{
	// 			{
	// 				SKU:             "WATCH-SILVER-001-S",
	// 				PriceAdjustment: decimal.NewFromFloat(-50.00),
	// 				Attributes:      models.AttributesMap{"size": "small", "material": "silver"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          6,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 1,
	// 					BackorderAllowed:  false,
	// 					MerchantID:        merchantID2,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "WATCH-SILVER-001-M",
	// 				PriceAdjustment: decimal.NewFromFloat(0.00),
	// 				Attributes:      models.AttributesMap{"size": "medium", "material": "silver"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          5,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 1,
	// 					BackorderAllowed:  false,
	// 					MerchantID:        merchantID2,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "WATCH-SILVER-001-L",
	// 				PriceAdjustment: decimal.NewFromFloat(50.00),
	// 				Attributes:      models.AttributesMap{"size": "large", "material": "silver"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          3,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 1,
	// 					BackorderAllowed:  false,
	// 					MerchantID:        merchantID2,
	// 				},
	// 			},
	// 		},
	// 	},

	// 	// Footwears Products
	// 	{
	// 		ID:          uuid.New().String(),
	// 		MerchantID:  merchantID1,
	// 		CategoryID:  footwearsCategory.ID,
	// 		Name:        "Elegant Beige High Heel Shoes",
	// 		Description: "Stylish beige high heel shoes for special occasions",
	// 		SKU:         "SHOE-BEIGE-001",
	// 		BasePrice:   decimal.NewFromFloat(120.00),
	// 		Media: []models.Media{
	// 			{URL: "https://www.shutterstock.com/image-photo/elegant-beige-high-heel-shoes-260nw-2635369899.jpg", Type: models.MediaTypeImage},
	// 		},
	// 		Variants: []models.Variant{
	// 			{
	// 				SKU:             "SHOE-BEIGE-001-37",
	// 				PriceAdjustment: decimal.NewFromFloat(0.00),
	// 				Attributes:      models.AttributesMap{"size": "37", "color": "beige"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          12,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 2,
	// 					BackorderAllowed:  true,
	// 					MerchantID:        merchantID1,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "SHOE-BEIGE-001-38",
	// 				PriceAdjustment: decimal.NewFromFloat(0.00),
	// 				Attributes:      models.AttributesMap{"size": "38", "color": "beige"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          10,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 2,
	// 					BackorderAllowed:  true,
	// 					MerchantID:        merchantID1,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "SHOE-BEIGE-001-39",
	// 				PriceAdjustment: decimal.NewFromFloat(0.00),
	// 				Attributes:      models.AttributesMap{"size": "39", "color": "beige"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          8,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 1,
	// 					BackorderAllowed:  false,
	// 					MerchantID:        merchantID1,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "SHOE-BEIGE-001-40",
	// 				PriceAdjustment: decimal.NewFromFloat(0.00),
	// 				Attributes:      models.AttributesMap{"size": "40", "color": "beige"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          6,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 1,
	// 					BackorderAllowed:  false,
	// 					MerchantID:        merchantID1,
	// 				},
	// 			},
	// 		},
	// 	},
	// 	{
	// 		ID:          uuid.New().String(),
	// 		MerchantID:  merchantID2,
	// 		CategoryID:  footwearsCategory.ID,
	// 		Name:        "Colorful Women's Sneakers",
	// 		Description: "Flying colorful women's sneakers for sports and casual wear",
	// 		SKU:         "SHOE-SNEAK-001",
	// 		BasePrice:   decimal.NewFromFloat(85.00),
	// 		Media: []models.Media{
	// 			{URL: "https://media.istockphoto.com/id/1436061606/photo/flying-colorful-womens-sneaker-isolated-on-white-background-fashionable-stylish-sports-shoe.jpg?s=612x612&w=0&k=20&c=2KKjX9tXo0ibmBaPlflnJNdtZ-J77wrprVStaPL2Gj4=", Type: models.MediaTypeImage},
	// 		},
	// 		Variants: []models.Variant{
	// 			{
	// 				SKU:             "SHOE-SNEAK-001-36",
	// 				PriceAdjustment: decimal.NewFromFloat(0.00),
	// 				Attributes:      models.AttributesMap{"size": "36", "color": "multi"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          20,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 3,
	// 					BackorderAllowed:  true,
	// 					MerchantID:        merchantID2,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "SHOE-SNEAK-001-37",
	// 				PriceAdjustment: decimal.NewFromFloat(0.00),
	// 				Attributes:      models.AttributesMap{"size": "37", "color": "multi"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          18,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 3,
	// 					BackorderAllowed:  true,
	// 					MerchantID:        merchantID2,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "SHOE-SNEAK-001-38",
	// 				PriceAdjustment: decimal.NewFromFloat(0.00),
	// 				Attributes:      models.AttributesMap{"size": "38", "color": "multi"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          15,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 2,
	// 					BackorderAllowed:  true,
	// 					MerchantID:        merchantID2,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "SHOE-SNEAK-001-39",
	// 				PriceAdjustment: decimal.NewFromFloat(0.00),
	// 				Attributes:      models.AttributesMap{"size": "39", "color": "multi"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          12,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 2,
	// 					BackorderAllowed:  false,
	// 					MerchantID:        merchantID2,
	// 				},
	// 			},
	// 		},
	// 	},
	// 	{
	// 		ID:          uuid.New().String(),
	// 		MerchantID:  merchantID1,
	// 		CategoryID:  footwearsCategory.ID,
	// 		Name:        "Purple Sports Sneakers",
	// 		Description: "Comfortable purple sports sneakers for active lifestyle",
	// 		SKU:         "SHOE-PURPLE-001",
	// 		BasePrice:   decimal.NewFromFloat(95.00),
	// 		Media: []models.Media{
	// 			{URL: "https://media.istockphoto.com/id/1411635454/photo/colorful-purple-sneakers-isolated-over-white-studio-background-comfortable-shoes-sport.jpg?s=612x612&w=0&k=20&c=AETXTH7nNFzE2eLrPY8Ke4ZbklXM9xSs_y3e6SzQ4x8=", Type: models.MediaTypeImage},
	// 		},
	// 		Variants: []models.Variant{
	// 			{
	// 				SKU:             "SHOE-PURPLE-001-36",
	// 				PriceAdjustment: decimal.NewFromFloat(-5.00),
	// 				Attributes:      models.AttributesMap{"size": "36", "color": "purple"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          15,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 2,
	// 					BackorderAllowed:  true,
	// 					MerchantID:        merchantID1,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "SHOE-PURPLE-001-37",
	// 				PriceAdjustment: decimal.NewFromFloat(0.00),
	// 				Attributes:      models.AttributesMap{"size": "37", "color": "purple"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          14,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 2,
	// 					BackorderAllowed:  true,
	// 					MerchantID:        merchantID1,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "SHOE-PURPLE-001-38",
	// 				PriceAdjustment: decimal.NewFromFloat(0.00),
	// 				Attributes:      models.AttributesMap{"size": "38", "color": "purple"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          12,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 2,
	// 					BackorderAllowed:  true,
	// 					MerchantID:        merchantID1,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "SHOE-PURPLE-001-39",
	// 				PriceAdjustment: decimal.NewFromFloat(5.00),
	// 				Attributes:      models.AttributesMap{"size": "39", "color": "purple"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          10,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 1,
	// 					BackorderAllowed:  false,
	// 					MerchantID:        merchantID1,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "SHOE-PURPLE-001-40",
	// 				PriceAdjustment: decimal.NewFromFloat(5.00),
	// 				Attributes:      models.AttributesMap{"size": "40", "color": "purple"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          8,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 1,
	// 					BackorderAllowed:  false,
	// 					MerchantID:        merchantID1,
	// 				},
	// 			},
	// 		},
	// 	},

	// 	// Ankaras Products
	// 	{
	// 		ID:          uuid.New().String(),
	// 		MerchantID:  merchantID2,
	// 		CategoryID:  ankarasCategory.ID,
	// 		Name:        "African Ethnic Native Pattern",
	// 		Description: "Traditional African ethnic native pattern fabric",
	// 		SKU:         "ANKARA-PAT-001",
	// 		BasePrice:   decimal.NewFromFloat(25.00),
	// 		Media: []models.Media{
	// 			{URL: "https://www.shutterstock.com/image-vector/african-ethnic-native-patterntraditional-kenteankarakitengechitengecapulana-260nw-2660177915.jpg", Type: models.MediaTypeImage},
	// 		},
	// 		Variants: []models.Variant{
	// 			{
	// 				SKU:             "ANKARA-PAT-001-S",
	// 				PriceAdjustment: decimal.NewFromFloat(0.00),
	// 				Attributes:      models.AttributesMap{"size": "small", "pattern": "kente"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          50,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 10,
	// 					BackorderAllowed:  true,
	// 					MerchantID:        merchantID2,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "ANKARA-PAT-001-M",
	// 				PriceAdjustment: decimal.NewFromFloat(5.00),
	// 				Attributes:      models.AttributesMap{"size": "medium", "pattern": "kente"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          40,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 8,
	// 					BackorderAllowed:  true,
	// 					MerchantID:        merchantID2,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "ANKARA-PAT-001-L",
	// 				PriceAdjustment: decimal.NewFromFloat(10.00),
	// 				Attributes:      models.AttributesMap{"size": "large", "pattern": "kente"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          30,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 5,
	// 					BackorderAllowed:  true,
	// 					MerchantID:        merchantID2,
	// 				},
	// 			},
	// 		},
	// 	},
	// 	{
	// 		ID:          uuid.New().String(),
	// 		MerchantID:  merchantID1,
	// 		CategoryID:  ankarasCategory.ID,
	// 		Name:        "Red Abstract Floral Ankara",
	// 		Description: "Beautiful red abstract floral traditional African fabric",
	// 		SKU:         "ANKARA-RED-001",
	// 		BasePrice:   decimal.NewFromFloat(30.00),
	// 		Media: []models.Media{
	// 			{URL: "https://www.shutterstock.com/image-vector/red-abstract-floral-traditional-african-260nw-2660124201.jpg", Type: models.MediaTypeImage},
	// 		},
	// 		Variants: []models.Variant{
	// 			{
	// 				SKU:             "ANKARA-RED-001-S",
	// 				PriceAdjustment: decimal.NewFromFloat(0.00),
	// 				Attributes:      models.AttributesMap{"size": "small", "pattern": "floral"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          45,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 8,
	// 					BackorderAllowed:  true,
	// 					MerchantID:        merchantID1,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "ANKARA-RED-001-M",
	// 				PriceAdjustment: decimal.NewFromFloat(5.00),
	// 				Attributes:      models.AttributesMap{"size": "medium", "pattern": "floral"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          35,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 6,
	// 					BackorderAllowed:  true,
	// 					MerchantID:        merchantID1,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "ANKARA-RED-001-L",
	// 				PriceAdjustment: decimal.NewFromFloat(10.00),
	// 				Attributes:      models.AttributesMap{"size": "large", "pattern": "floral"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          25,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 4,
	// 					BackorderAllowed:  false,
	// 					MerchantID:        merchantID1,
	// 				},
	// 			},
	// 		},
	// 	},
	// 	{
	// 		ID:          uuid.New().String(),
	// 		MerchantID:  merchantID2,
	// 		CategoryID:  ankarasCategory.ID,
	// 		Name:        "African Tribal Clash Ornament",
	// 		Description: "Traditional African ethnic tribal clash ornament fabric",
	// 		SKU:         "ANKARA-TRIB-001",
	// 		BasePrice:   decimal.NewFromFloat(35.00),
	// 		Media: []models.Media{
	// 			{URL: "https://www.shutterstock.com/image-vector/african-ethnic-tribal-clash-ornament-260nw-2674095379.jpg", Type: models.MediaTypeImage},
	// 		},
	// 		Variants: []models.Variant{
	// 			{
	// 				SKU:             "ANKARA-TRIB-001-S",
	// 				PriceAdjustment: decimal.NewFromFloat(0.00),
	// 				Attributes:      models.AttributesMap{"size": "small", "pattern": "tribal"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          40,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 7,
	// 					BackorderAllowed:  true,
	// 					MerchantID:        merchantID2,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "ANKARA-TRIB-001-M",
	// 				PriceAdjustment: decimal.NewFromFloat(5.00),
	// 				Attributes:      models.AttributesMap{"size": "medium", "pattern": "tribal"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          30,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 5,
	// 					BackorderAllowed:  true,
	// 					MerchantID:        merchantID2,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "ANKARA-TRIB-001-L",
	// 				PriceAdjustment: decimal.NewFromFloat(10.00),
	// 				Attributes:      models.AttributesMap{"size": "large", "pattern": "tribal"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          20,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 3,
	// 					BackorderAllowed:  false,
	// 					MerchantID:        merchantID2,
	// 				},
	// 			},
	// 		},
	// 	},

	// 	// Neckwears Products
	// 	{
	// 		ID:          uuid.New().String(),
	// 		MerchantID:  merchantID1,
	// 		CategoryID:  neckwearsCategory.ID,
	// 		Name:        "Black Handkerchief",
	// 		Description: "Classic black handkerchief for formal and casual wear",
	// 		SKU:         "NECK-BLK-001",
	// 		BasePrice:   decimal.NewFromFloat(15.00),
	// 		Media: []models.Media{
	// 			{URL: "https://www.shutterstock.com/image-photo/one-black-handkerchief-isolated-on-600nw-2593330633.jpg", Type: models.MediaTypeImage},
	// 		},
	// 		Variants: []models.Variant{
	// 			{
	// 				SKU:             "NECK-BLK-001-S",
	// 				PriceAdjustment: decimal.NewFromFloat(0.00),
	// 				Attributes:      models.AttributesMap{"size": "small", "color": "black"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          100,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 20,
	// 					BackorderAllowed:  true,
	// 					MerchantID:        merchantID1,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "NECK-BLK-001-M",
	// 				PriceAdjustment: decimal.NewFromFloat(2.00),
	// 				Attributes:      models.AttributesMap{"size": "medium", "color": "black"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          80,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 15,
	// 					BackorderAllowed:  true,
	// 					MerchantID:        merchantID1,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "NECK-BLK-001-L",
	// 				PriceAdjustment: decimal.NewFromFloat(4.00),
	// 				Attributes:      models.AttributesMap{"size": "large", "color": "black"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          60,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 10,
	// 					BackorderAllowed:  true,
	// 					MerchantID:        merchantID1,
	// 				},
	// 			},
	// 		},
	// 	},
	// 	{
	// 		ID:          uuid.New().String(),
	// 		MerchantID:  merchantID2,
	// 		CategoryID:  neckwearsCategory.ID,
	// 		Name:        "Stylish Neck Tie",
	// 		Description: "Modern stylish neck tie for professional look",
	// 		SKU:         "NECK-TIE-001",
	// 		BasePrice:   decimal.NewFromFloat(25.00),
	// 		Media: []models.Media{
	// 			{URL: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQSyNHU99LxG-qQUOuUUH9QwBLC5o3-I0ORFMc4zL3D58ymeFI&s", Type: models.MediaTypeImage},
	// 		},
	// 		Variants: []models.Variant{
	// 			{
	// 				SKU:             "NECK-TIE-001-S",
	// 				PriceAdjustment: decimal.NewFromFloat(0.00),
	// 				Attributes:      models.AttributesMap{"size": "small", "color": "blue"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          50,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 10,
	// 					BackorderAllowed:  true,
	// 					MerchantID:        merchantID2,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "NECK-TIE-001-M",
	// 				PriceAdjustment: decimal.NewFromFloat(0.00),
	// 				Attributes:      models.AttributesMap{"size": "medium", "color": "blue"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          40,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 8,
	// 					BackorderAllowed:  true,
	// 					MerchantID:        merchantID2,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "NECK-TIE-001-L",
	// 				PriceAdjustment: decimal.NewFromFloat(5.00),
	// 				Attributes:      models.AttributesMap{"size": "large", "color": "blue"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          30,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 5,
	// 					BackorderAllowed:  false,
	// 					MerchantID:        merchantID2,
	// 				},
	// 			},
	// 		},
	// 	},
	// 	{
	// 		ID:          uuid.New().String(),
	// 		MerchantID:  merchantID1,
	// 		CategoryID:  neckwearsCategory.ID,
	// 		Name:        "Blue Bow Tie",
	// 		Description: "Elegant blue bow tie for formal occasions",
	// 		SKU:         "NECK-BOW-001",
	// 		BasePrice:   decimal.NewFromFloat(20.00),
	// 		Media: []models.Media{
	// 			{URL: "https://media.istockphoto.com/id/1298105203/photo/neckties-men-accessories-mens-fashion-bow-blue-tie-isolated-on-white.jpg?s=612x612&w=0&k=20&c=AGkxwafREqTb84EbzxJwNKXmdCxMrejEIGQLSxru79g=", Type: models.MediaTypeImage},
	// 		},
	// 		Variants: []models.Variant{
	// 			{
	// 				SKU:             "NECK-BOW-001-S",
	// 				PriceAdjustment: decimal.NewFromFloat(0.00),
	// 				Attributes:      models.AttributesMap{"size": "small", "color": "blue"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          60,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 12,
	// 					BackorderAllowed:  true,
	// 					MerchantID:        merchantID1,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "NECK-BOW-001-M",
	// 				PriceAdjustment: decimal.NewFromFloat(2.00),
	// 				Attributes:      models.AttributesMap{"size": "medium", "color": "blue"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          50,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 10,
	// 					BackorderAllowed:  true,
	// 					MerchantID:        merchantID1,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "NECK-BOW-001-L",
	// 				PriceAdjustment: decimal.NewFromFloat(4.00),
	// 				Attributes:      models.AttributesMap{"size": "large", "color": "blue"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          40,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 8,
	// 					BackorderAllowed:  true,
	// 					MerchantID:        merchantID1,
	// 				},
	// 			},
	// 		},
	// 	},

	// 	// Bags Products
	// 	{
	// 		ID:          uuid.New().String(),
	// 		MerchantID:  merchantID2,
	// 		CategoryID:  bagsCategory.ID,
	// 		Name:        "White Female Handbags Collection",
	// 		Description: "Elegant collection of white female handbags",
	// 		SKU:         "BAG-WHITE-001",
	// 		BasePrice:   decimal.NewFromFloat(75.00),
	// 		Media: []models.Media{
	// 			{URL: "https://www.shutterstock.com/image-photo/white-female-handbags-collection-on-260nw-739041304.jpg", Type: models.MediaTypeImage},
	// 		},
	// 		Variants: []models.Variant{
	// 			{
	// 				SKU:             "BAG-WHITE-001-S",
	// 				PriceAdjustment: decimal.NewFromFloat(-10.00),
	// 				Attributes:      models.AttributesMap{"size": "small", "color": "white"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          25,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 5,
	// 					BackorderAllowed:  true,
	// 					MerchantID:        merchantID2,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "BAG-WHITE-001-M",
	// 				PriceAdjustment: decimal.NewFromFloat(0.00),
	// 				Attributes:      models.AttributesMap{"size": "medium", "color": "white"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          20,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 4,
	// 					BackorderAllowed:  true,
	// 					MerchantID:        merchantID2,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "BAG-WHITE-001-L",
	// 				PriceAdjustment: decimal.NewFromFloat(10.00),
	// 				Attributes:      models.AttributesMap{"size": "large", "color": "white"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          15,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 3,
	// 					BackorderAllowed:  false,
	// 					MerchantID:        merchantID2,
	// 				},
	// 			},
	// 		},
	// 	},
	// 	{
	// 		ID:          uuid.New().String(),
	// 		MerchantID:  merchantID1,
	// 		CategoryID:  bagsCategory.ID,
	// 		Name:        "Blue Fashion Purse",
	// 		Description: "Stylish blue fashion purse for everyday use",
	// 		SKU:         "BAG-BLUE-001",
	// 		BasePrice:   decimal.NewFromFloat(65.00),
	// 		Media: []models.Media{
	// 			{URL: "https://media.istockphoto.com/id/1365118618/photo/blue-fashion-purse-handbag-on-white-background-isolated.jpg?s=612x612&w=0&k=20&c=VNszfC0cxenqZGhjlr3gqqvzHWREuhdY_H3CKF1B38g=", Type: models.MediaTypeImage},
	// 		},
	// 		Variants: []models.Variant{
	// 			{
	// 				SKU:             "BAG-BLUE-001-S",
	// 				PriceAdjustment: decimal.NewFromFloat(-5.00),
	// 				Attributes:      models.AttributesMap{"size": "small", "color": "blue"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          30,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 6,
	// 					BackorderAllowed:  true,
	// 					MerchantID:        merchantID1,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "BAG-BLUE-001-M",
	// 				PriceAdjustment: decimal.NewFromFloat(0.00),
	// 				Attributes:      models.AttributesMap{"size": "medium", "color": "blue"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          25,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 5,
	// 					BackorderAllowed:  true,
	// 					MerchantID:        merchantID1,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "BAG-BLUE-001-L",
	// 				PriceAdjustment: decimal.NewFromFloat(5.00),
	// 				Attributes:      models.AttributesMap{"size": "large", "color": "blue"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          20,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 4,
	// 					BackorderAllowed:  false,
	// 					MerchantID:        merchantID1,
	// 				},
	// 			},
	// 		},
	// 	},
	// 	{
	// 		ID:          uuid.New().String(),
	// 		MerchantID:  merchantID2,
	// 		CategoryID:  bagsCategory.ID,
	// 		Name:        "Premium Leather Handbag",
	// 		Description: "Luxury premium leather handbag for fashion enthusiasts",
	// 		SKU:         "BAG-LEATH-001",
	// 		BasePrice:   decimal.NewFromFloat(120.00),
	// 		Media: []models.Media{
	// 			{URL: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTAUFyEw7i2H_NI_Ipf0Pe0uCks3vOqsc-KQLk00b-0MYL1j0lO&s", Type: models.MediaTypeImage},
	// 		},
	// 		Variants: []models.Variant{
	// 			{
	// 				SKU:             "BAG-LEATH-001-S",
	// 				PriceAdjustment: decimal.NewFromFloat(-15.00),
	// 				Attributes:      models.AttributesMap{"size": "small", "color": "brown"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          15,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 3,
	// 					BackorderAllowed:  false,
	// 					MerchantID:        merchantID2,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "BAG-LEATH-001-M",
	// 				PriceAdjustment: decimal.NewFromFloat(0.00),
	// 				Attributes:      models.AttributesMap{"size": "medium", "color": "brown"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          12,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 2,
	// 					BackorderAllowed:  false,
	// 					MerchantID:        merchantID2,
	// 				},
	// 			},
	// 			{
	// 				SKU:             "BAG-LEATH-001-L",
	// 				PriceAdjustment: decimal.NewFromFloat(15.00),
	// 				Attributes:      models.AttributesMap{"size": "large", "color": "brown"},
	// 				IsActive:        true,
	// 				Inventory: models.Inventory{
	// 					Quantity:          8,
	// 					ReservedQuantity:  0,
	// 					LowStockThreshold: 1,
	// 					BackorderAllowed:  false,
	// 					MerchantID:        merchantID2,
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	// // Create products with variants and media
	// for i := range products {
	// 	if err := db.Create(&products[i]).Error; err != nil {
	// 		log.Printf("Warning: Failed to seed product %s: %v", products[i].Name, err)
	// 	}
	// }



	//fmt.Println("Database seeded successfully with 30+ products across 6 categories!")
	fmt.Println("Database migrated ")
}
