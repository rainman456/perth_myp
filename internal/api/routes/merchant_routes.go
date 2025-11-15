package routes

import (
	"api-customer-merchant/internal/api/handlers" // "api-customer-merchant/internal/api/handlers"
	"api-customer-merchant/internal/config"
	"api-customer-merchant/internal/db/repositories"
	"api-customer-merchant/internal/middleware"
	"api-customer-merchant/internal/services/dispute"
	"api-customer-merchant/internal/services/merchant"
	"api-customer-merchant/internal/services/order"
	"api-customer-merchant/internal/services/payment"
	"api-customer-merchant/internal/services/payout"
	"api-customer-merchant/internal/services/product"
	"api-customer-merchant/internal/services/email"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

)

func SetupMerchantRoutes(r *gin.Engine) {
	cfg := config.Load()
	logger, _ := zap.NewProduction()

	appRepo := repositories.NewMerchantApplicationRepository()
	merchantRepo := repositories.NewMerchantRepository()

	merchantService := merchant.NewMerchantService(appRepo, merchantRepo)

	productRepo := repositories.NewProductRepository()
	// reviewRepo:= repositories.NewReviewRepository()
	productService := product.NewProductService(productRepo, cfg, logger)

	// Order service for merchant orders
	orderRepo := repositories.NewOrderRepository()
	orderitemRepo := repositories.NewOrderItemRepository()
	cartRepo := repositories.NewCartRepository()
	cartitemRepo := repositories.NewCartItemRepository()
	inventoryRepo := repositories.NewInventoryRepository()
	userRepo := repositories.NewUserRepository()

	// Payment service initialization
	paymentRepo := repositories.NewPaymentRepository()
	payoutRepo := repositories.NewPayoutRepository()
	paymentService := payment.NewPaymentService(
		paymentRepo,
		orderRepo,
		payoutRepo,
		merchantRepo,
		cfg,
		logger,
	)

	// Email service initialization
	emailService := email.NewEmailService()

	orderService := order.NewOrderService(
		orderRepo,
		orderitemRepo,
		cartRepo,
		cartitemRepo,
		productRepo,
		inventoryRepo,
		userRepo,
		paymentService,
		emailService,
		cfg,
		logger,
	)

	// Dispute service
	disputeRepo := repositories.NewDisputeRepository()
	disputeService := dispute.NewDisputeService(disputeRepo, orderRepo, logger)

	// Payout service
	payoutService := payout.NewPayoutService(payoutRepo)

	merchantOrderHandler := handlers.NewMerchantOrderHandler(orderService, logger)
	merchantPayoutHandler := handlers.NewPayoutHandler(payoutService, logger)
	merchantDisputeHandler := handlers.NewMerchantDisputeHandler(disputeService)

	merchantAuthHandler := handlers.NewMerchantAuthHandler(merchantService)
	mediaHandler := handlers.NewProductMediaHandler(productService, logger)
	merchantproductHandler := handlers.NewProductHandlers(productService, logger)

	merchantGroup := r.Group("/merchant")
	{
		merchantGroup.POST("/apply", merchantAuthHandler.Apply)
		merchantGroup.GET("/application/:id", merchantAuthHandler.GetApplication)
		merchantGroup.POST("/login", merchantAuthHandler.Login)

		protected := merchantGroup.Group("")
		protected.Use(middleware.AuthMiddleware("merchant"))
		{
			protected.GET("/me", merchantAuthHandler.GetMyMerchant)
			protected.PUT("/profile", merchantAuthHandler.UpdateProfile)
			protected.POST("/logout", merchantAuthHandler.Logout)

			// Merchant orders
			ordersGroup := protected.Group("/orders")
			{
				ordersGroup.GET("", merchantOrderHandler.GetMerchantOrders)
				ordersGroup.GET("/:id", merchantOrderHandler.GetMerchantOrder)

				// Merchant order item actions
				orderItemsGroup := ordersGroup.Group("/items")
				{
					orderItemsGroup.POST("/:id/accept", merchantOrderHandler.AcceptOrderItem)
					orderItemsGroup.POST("/:id/decline", merchantOrderHandler.DeclineOrderItem)
					orderItemsGroup.POST("/:id/sent-to-aronova-hub", merchantOrderHandler.UpdateOrderItemToSentToAronovaHub)
				}
			}

			// Merchant disputes
			disputesGroup := protected.Group("/disputes")
			{
				disputesGroup.GET("", merchantDisputeHandler.ListMerchantDisputes)
				disputesGroup.PUT("/:id", merchantDisputeHandler.UpdateDispute)
			}

			// Merchant payouts
			payoutsGroup := protected.Group("/payouts")
			{
				payoutsGroup.GET("", merchantPayoutHandler.GetMerchantPayouts)
				payoutsGroup.POST("/request", merchantPayoutHandler.RequestPayout)
			}

			productsGroup := protected.Group("/products")
			{
				productsGroup.POST("", merchantproductHandler.CreateProduct)
				productsGroup.POST("/bulk-upload", merchantproductHandler.BulkUploadProducts)           // Add bulk upload route
				productsGroup.PUT("/bulk-update", merchantproductHandler.BulkUpdateProducts)            // Add bulk update route
				productsGroup.PUT("/bulk-inventory-update", merchantproductHandler.BulkUpdateInventory) // Add bulk inventory update route
				productsGroup.PUT("/:id", merchantproductHandler.UpdateProduct)                         // Add update product route
				productsGroup.GET("", func(c *gin.Context) {
					// Override to use merchantID from context for list
					merchantID, _ := c.Get("merchantID")
					c.Set("id", merchantID.(string)) // Temporary set for handler compatibility
					merchantproductHandler.ListProductsByMerchant(c)
				})

				// All :id actions (delete, media) under /:id subgroup
				singleProductGroup := productsGroup.Group("/:id")
				{
					singleProductGroup.DELETE("", merchantproductHandler.DeleteProduct)

					// Media nested under /:id/media
					mediaSubGroup := singleProductGroup.Group("/media")
					{
						mediaSubGroup.POST("", mediaHandler.UploadMedia)
						mediaSubGroup.PUT("/:media_id", mediaHandler.UpdateMedia)
						mediaSubGroup.DELETE("/:media_id", mediaHandler.DeleteMedia)
					}
				}

				// Inventory separate (not under /:id to avoid nesting issues)
				productsGroup.PUT("/inventory/:id", merchantproductHandler.UpdateInventory)
				// Add update variant route
				productsGroup.PUT("/variants/:id", merchantproductHandler.UpdateVariant)
			}
		}
	}
}
