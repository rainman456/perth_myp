package db

import (
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

	//  log.Println("Migrating Merchant...")
	//  if err := DB.AutoMigrate(&models.Merchant{}); err != nil {
	//      log.Printf("Failed to migrate Merchant: %v", err)
	//      return
	//  }

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
	//&models.Dispute{},
	//&models.OrderMerchantSplit{},
	//&models.ReturnRequest{},
//&models.MerchantBankDetails{},
	//&models.User{},
    //&models.UserAddress{},
	//&models.MerchantApplication{},
	// &models.Product{},
	// &models.Variant{},
	// &models.Media{},
	 //&models.Cart{},
	// &models.Order{},
	// &models.OrderItem{},
	 //&models.CartItem{},
	//&models.Category{},
	 //&models.Inventory{},
	// &models.Promotion{},
	// &models.ReturnRequest{},
	 //&models.Payout{},
	//&models.Payment{},
	//&models.Review{},
	//&models.UserWishlist{},
	
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
