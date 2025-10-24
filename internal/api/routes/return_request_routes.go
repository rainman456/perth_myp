package routes

import (
	"api-customer-merchant/internal/api/handlers"
	"api-customer-merchant/internal/db/repositories"
	"api-customer-merchant/internal/middleware"
	"api-customer-merchant/internal/services/return_request"

	"github.com/gin-gonic/gin"
)

func SetupReturnRequestRoutes(r *gin.Engine) {
	returnReqRepo := repositories.NewReturnRequestRepository()
	returnReqService := return_request.NewReturnRequestService(returnReqRepo)
	returnReqHandler := handlers.NewReturnRequestHandler(returnReqService)

	returnReqGroup := r.Group("/return-requests")
	protected := returnReqGroup.Use(middleware.AuthMiddleware("customer"))
	protected.POST("", returnReqHandler.CreateReturnRequest)
	protected.GET("/:id", returnReqHandler.GetReturnRequest)
	protected.GET("/:orderId", returnReqHandler.GetReturnRequestsByOrderID)
	protected.GET("", returnReqHandler.ListCustomerReturnRequests)
}