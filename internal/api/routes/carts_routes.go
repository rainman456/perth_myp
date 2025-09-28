package routes

import (
	"api-customer-merchant/internal/api/handlers"
	"api-customer-merchant/internal/db/repositories"

	//"api-customer-merchant/internal/middleware"
	"api-customer-merchant/internal/services/cart"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupCartRoutes(r *gin.Engine) {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logger.Sync() // Ensure logger flushes logs
	inventoryRepo := repositories.NewInventoryRepository()
	cartitemRepo := repositories.NewCartItemRepository()
	cartRepo := repositories.NewCartRepository()
	productRepo := repositories.NewProductRepository()
	cartService := cart.NewCartService(cartRepo, cartitemRepo, productRepo, inventoryRepo, logger)
	cartHandlers := handlers.NewCartHandler(cartService,logger)
	//protected := middleware.AuthMiddleware("user")
	r.GET("/cart", cartHandlers.GetCart)
	r.POST("/cart/items", cartHandlers.AddToCart)
	r.GET("/cart/items/:id", cartHandlers.GetCartItem)
	r.PUT("/cart/items/:id", cartHandlers.UpdateCartItemQuantity)
	r.DELETE("/cart/items/:id", cartHandlers.RemoveCartItem)
	//r.POST("/cart/clear", protected, customerHandlers.ClearCart)
}
