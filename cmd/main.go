package main

import (
	"fmt"
	"log"
	"os"

	"api-customer-merchant/internal/api/handlers"
//"api-customer-merchant/internal/api/handlers"
	"api-customer-merchant/internal/api/routes"
//"api-customer-merchant/internal/bank"
	"api-customer-merchant/internal/config"

	//"api-customer-merchant/internal/services/cart"
	//"api-customer-merchant/internal/services/order"

	//"api-customer-merchant/internal/middleware"
	"api-customer-merchant/internal/db"
	//"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	//_ "api-customer-merchant/docs" // Import generated docs
	"github.com/gin-contrib/cors"
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
	//models.BackfillCategorySlugs(db.DB)
	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS","PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	// Create single router
	//r := gin.Default()

	// Customer routes under /customer
	routes.RegisterCustomerRoutes(r)
	routes.SetupMerchantRoutes(r)
	routes.SetupProductRoutes(r)
	routes.SetupOrderRoutes(r)
	routes.SetupCartRoutes(r)
	routes.SetupDisputeRoutes(r)
	routes.SetupReturnRequestRoutes(r)
	routes.SetupReviewRoutes(r)
	routes.SetupWishlistRoutes(r)
	routes.RegisterPaymentRoutes(r)


	//svc := bank.NewFetchBankService()
	// Optionally set a different cache file:
	// svc.CacheFile = "data/banks.json" // if you expose that field in the constructor

	handler := handlers.NewBankHandler()

	r.GET("/banks", handler.GetBanks)
	

	// Swagger endpoint
	// Serve the OpenAPI spec file
	//r.StaticFile("/swagger.yaml", "./swagger.yaml") // Or "/swagger.json" if using JSON
	r.StaticFile("/swagger.json", "./swagger.json")
	// Configure Swagger UI to load the spec
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger.json"))) // Or "/swagger.json"

	// Get port from environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Run on 0.0.0.0:port for Railway compatibility
	addr := fmt.Sprintf("0.0.0.0:%s", port)
	log.Printf("Example app listening on port %s", port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("API failed: %v", err)
	}
}
