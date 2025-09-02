package main

import (
	"fmt"
	"log"
	"os"

	customerHandlers "api-customer-merchant/internal/customer/handlers"
	merchantHandlers "api-customer-merchant/internal/merchant/handlers"
	"api-customer-merchant/internal/middleware"
	"api-customer-merchant/internal/shared/auth/models"
	"api-customer-merchant/internal/shared/db"

	"github.com/gin-gonic/gin"

	"github.com/swaggo/gin-swagger"
    swaggerFiles "github.com/swaggo/files"
	//"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "api-customer-merchant/docs"

)

func main() {
	// Connect to database and migrate
	db.Connect()
	db.DB.AutoMigrate(&models.User{})

	// Create single router
	r := gin.Default()
	

	// Customer routes under /customer
	customer := r.Group("/customer")
	{
		// Public routes
		customerAuth := customerHandlers.NewAuthHandler()
		customer.POST("/register", customerAuth.Register)
		customer.POST("/login", customerAuth.Login)
		customer.GET("/auth/google", customerAuth.GoogleAuth)
		customer.GET("/auth/google/callback", customerAuth.GoogleCallback)

		// Protected routes
		protected := customer.Group("/")
		protected.Use(middleware.AuthMiddleware("customer"))
		protected.POST("/logout", customerAuth.Logout)
	}

	// Merchant routes under /merchant
	merchant := r.Group("/merchant")
	{
		// Public routes
		merchantAuth := merchantHandlers.NewAuthHandler()
		merchant.POST("/register", merchantAuth.Register)
		merchant.POST("/login", merchantAuth.Login)
		merchant.GET("/auth/google", merchantAuth.GoogleAuth)
		merchant.GET("/auth/google/callback", merchantAuth.GoogleCallback)

		// Protected routes
		protected := merchant.Group("/")
		protected.Use(middleware.AuthMiddleware("merchant"))
		protected.POST("/logout", merchantAuth.Logout)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/swagger.json")))
	// Run on :8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Run on 0.0.0.0:port for Render compatibility
	addr := fmt.Sprintf("0.0.0.0:%s", port)
	log.Printf("Example app listening on port %s", port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("API failed: %v", err)
	}
}