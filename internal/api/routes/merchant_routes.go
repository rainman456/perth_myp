package routes

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
// "api-customer-merchant/internal/api/handlers"
// "api-customer-merchant/internal/middleware"
//"api-customer-merchant/internal/db"
//"api-customer-merchant/internal/db/models"
//"api-customer-merchant/internal/db/repositories"

//"api-customer-merchant/internal/middleware"
// "api-customer-merchant/internal/services/merchant"
// "api-customer-merchant/internal/services/product"

// "api-customer-merchant/internal/domain/order"
// "api-customer-merchant/internal/domain/payout"
//"api-customer-merchant/internal/domain/product"
//"api-customer-merchant/internal/domain/invent"

//"github.com/gin-gonic/gin"
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

/*
func RegisterMerchantRoutes(r *gin.Engine) {
    appRepo := repositories.NewMerchantApplicationRepository()
    repo := repositories.NewMerchantRepository()
    service := merchant.NewMerchantService(appRepo, repo)
    productRepo := repositories.NewProductRepository()
    inventoryRepo:= repositories.NewInventoryRepository()
	productService := product.NewProductService(productRepo,inventoryRepo)
    //productRepo := repositories.NewProductRepository()
    //inventoryRepo:= repositories.NewInventoryRepository()
	//productService := product.NewProductService(productRepo,inventoryRepo)

	// Other services (stubs; instantiate as needed)
	// orderService := order.NewOrderService(nil) // Adjust
	// payoutService := payout.NewPayoutService(nil)
	// promotionService := promotions.NewPromotionService(nil)
    //h := handlers.NewMerchantAuthHandler(service)

    authHandler := handlers.NewMerchantAuthHandler(service)
    merchhandler:=handlers.NewMerchantHandlers(productService)
    merchant := r.Group("/merchant")
    {
    // Application submission and status
    // merchant.POST("/apply",  authHandler.Apply)
    // merchant.GET("/application/:id",  authHandler.GetApplication)
    // // Merchant account access (once approved by admin via Express API)
    // merchant.GET("/me",  authHandler.GetMyMerchant)

    //merchant.POST("/create/product", merchhandler.CreateProduct)
	//merchant.GET("/products", merchhandler.GetMyProducts)
	//merchant.PUT("/:id", merchhandler.UpdateProduct)
	//merchant.DELETE("/:id", merchhandler.DeleteProduct)
	//merchant.POST("/bulk-upload", merchhandler.BulkUploadProducts)
    merchant.POST("/apply",  authHandler.Apply)
    merchant.POST("/login",  authHandler.Login)
    merchant.GET("/application/:id",  authHandler.GetApplication)
    // Merchant account access (once approved by admin via Express API)
    //merchant.GET("/me",  authHandler.GetMyMerchant)
    protected := merchant.Group("/")
    protected.Use(middleware.AuthMiddleware("merchant"))
    protected.GET("/me",  authHandler.GetMyMerchant)

    protected.POST("/products/create", merchhandler.CreateProduct)
	protected.GET("/products", merchhandler.GetMyProducts)
	protected.PUT("/products/:id", merchhandler.UpdateProduct)
	protected.DELETE("/products/:id", merchhandler.DeleteProduct)
	protected.POST("/products/bulk-upload", merchhandler.BulkUploadProducts)

    }


}
*/
import (
	"api-customer-merchant/internal/api/handlers"
	"api-customer-merchant/internal/config"
	"api-customer-merchant/internal/db/repositories"
	"api-customer-merchant/internal/middleware"
	"api-customer-merchant/internal/services/merchant"
	"api-customer-merchant/internal/services/product"

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
	//inventoryRepo := repositories.NewInventoryRepository()
	productService := product.NewProductService(productRepo, cfg, logger)

	merchantAuthHandler := handlers.NewMerchantAuthHandler(merchantService)
	//merchantHandler := handlers.NewMerchantHandlers(productService)
	mediaHandler := handlers.NewProductMediaHandler(productService, logger)
	merchantproductHandler := handlers.NewProductHandlers(productService, logger)

	merchantGroup := r.Group("/merchant")
	{
		merchantGroup.POST("/apply", merchantAuthHandler.Apply)
		merchantGroup.GET("/application/:id", merchantAuthHandler.GetApplication)
		merchantGroup.POST("/login", merchantAuthHandler.Login)

		protected := merchantGroup.Group("")
		protected.Use(middleware.AuthMiddleware("merchant"))
		protected.GET("/me", merchantAuthHandler.GetMyMerchant)
		protected.POST("/products", merchantproductHandler.CreateProduct) // Use productHandler for consistency
		protected.GET("/products", func(c *gin.Context) {
			// Override to use merchantID from context
			merchantID, _ := c.Get("merchantID")
			c.Set("id", merchantID.(string)) // Set param for handler
			merchantproductHandler.ListProductsByMerchant(c)
		})
		//protected.PUT("/products/:id", merchantproductHandler.UpdateProduct)
		protected.DELETE("/products/:id", merchantproductHandler.DeleteProduct)
		// protected.POST("/products/bulk-upload", merchantHandler.BulkUploadProducts) // If implemented

		protected.POST("/products/:product_id/media", mediaHandler.UploadMedia)
		protected.PUT("/products/:product_id/media/:media_id", mediaHandler.UpdateMedia)
		protected.DELETE("/products/:product_id/media/:media_id", mediaHandler.DeleteMedia)
		protected.PUT("/products/inventory/:id", merchantproductHandler.UpdateInventory)
	}
}