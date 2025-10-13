package routes

import (
	"api-customer-merchant/internal/api/handlers"
	"api-customer-merchant/internal/db/repositories"
	"api-customer-merchant/internal/middleware"
	"api-customer-merchant/internal/services/order"

	"github.com/gin-gonic/gin"
)

func SetupOrderRoutes(r *gin.Engine) {
	orderRepo := repositories.NewOrderRepository()
	orderitemRepo := repositories.NewOrderItemRepository()
	cartRepo := repositories.NewCartRepository()
	cartitemRepo := repositories.NewCartItemRepository()
	productRepo := repositories.NewProductRepository()
inventoryRepo := repositories.NewInventoryRepository() // Fixed: Added "New"	
 orderService := order.NewOrderService(orderRepo, orderitemRepo, cartRepo, cartitemRepo, productRepo,inventoryRepo)
	orderHandler := handlers.NewOrderHandler(orderService)
	protected := middleware.AuthMiddleware("customer") // Consider adding auth middleware
	r.POST("/orders",protected, orderHandler.CreateOrder) //create order
	r.GET("/orders/:id", protected,orderHandler.GetOrder) //get order
	r.POST("orders/:id/cancel", protected,orderHandler.CancelOrder)
	r.GET("/orders",protected,orderHandler.GetUserOrders)

	//r.GET("/orders/:id", protected, orderHandler.GetOrder)
}
