package routes

import (
	"api-customer-merchant/internal/api/handlers"
	"api-customer-merchant/internal/config"
	"api-customer-merchant/internal/db/repositories"
	"api-customer-merchant/internal/middleware"
	"api-customer-merchant/internal/services/order"
	"api-customer-merchant/internal/services/payment"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupOrderRoutes(r *gin.Engine) {
// 	orderRepo := repositories.NewOrderRepository()
// 	orderitemRepo := repositories.NewOrderItemRepository()
// 	cartRepo := repositories.NewCartRepository()
// 	cartitemRepo := repositories.NewCartItemRepository()
// 	productRepo := repositories.NewProductRepository()
// inventoryRepo := repositories.NewInventoryRepository() // Fixed: Added "New"
// userRepo := repositories.NewUserRepository()                    // ADD THIS

//  orderService := order.NewOrderService(orderRepo, orderitemRepo, cartRepo, cartitemRepo, productRepo,inventoryRepo)
// 	orderHandler := handlers.NewOrderHandler(orderService)
// 	protected := middleware.AuthMiddleware("customer") // Consider adding auth middleware
// 	r.POST("/orders",protected, orderHandler.CreateOrder) //create order
// 	r.GET("/orders/:id", protected,orderHandler.GetOrder) //get order
// 	r.POST("orders/:id/cancel", protected,orderHandler.CancelOrder)
// 	r.GET("/orders",protected,orderHandler.GetUserOrders)

// 	//r.GET("/orders/:id", protected, orderHandler.GetOrder)




logger, _ := zap.NewProduction()
    
orderRepo := repositories.NewOrderRepository()
orderitemRepo := repositories.NewOrderItemRepository()
cartRepo := repositories.NewCartRepository()
cartitemRepo := repositories.NewCartItemRepository()
productRepo := repositories.NewProductRepository()
inventoryRepo := repositories.NewInventoryRepository()
userRepo := repositories.NewUserRepository()                    // ADD THIS

// Payment service initialization
conf := config.Load()
paymentRepo := repositories.NewPaymentRepository()
payoutRepo := repositories.NewPayoutRepository()
merchantRepo := repositories.NewMerchantRepository()
paymentService := payment.NewPaymentService(                   // ADD THIS
	paymentRepo,
	orderRepo,
	payoutRepo,
	merchantRepo,
	conf,
	logger,
)

orderService := order.NewOrderService(
	orderRepo,
	orderitemRepo,
	cartRepo,
	cartitemRepo,
	productRepo,
	inventoryRepo,
	userRepo,
	paymentService,
	logger,  
)

orderHandler := handlers.NewOrderHandler(orderService)
protected := middleware.AuthMiddleware("customer")

r.POST("/orders", protected, orderHandler.CreateOrder)
r.GET("/orders/:id", protected, orderHandler.GetOrder)
r.POST("/orders/:id/cancel", protected, orderHandler.CancelOrder)
r.GET("/orders", protected, orderHandler.GetUserOrders)
}
