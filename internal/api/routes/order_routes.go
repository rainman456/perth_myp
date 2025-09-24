package routes

import (
	"api-customer-merchant/internal/api/handlers"
	"api-customer-merchant/internal/db/repositories"
	"api-customer-merchant/internal/services/order"

	"github.com/gin-gonic/gin"
)

func SetupOrderRoutes(r *gin.Engine) {
	orderRepo := repositories.NewOrderRepository()
	orderitemRepo := repositories.NewOrderItemRepository()
	cartRepo := repositories.NewCartRepository()
	cartitemRepo := repositories.NewCartItemRepository()
	productRepo := repositories.NewProductRepository()
	orderService := order.NewOrderService(orderRepo, orderitemRepo, cartRepo, cartitemRepo, productRepo)
	orderHandler := handlers.NewOrderHandler(orderService)
	//protected := middleware.AuthMiddleware("user") // Consider adding auth middleware
	r.POST("/orders", orderHandler.CreateOrder) //create order
	r.GET("/orders/:id", orderHandler.GetOrder) //get order

	//r.GET("/orders/:id", protected, orderHandler.GetOrder)
}
