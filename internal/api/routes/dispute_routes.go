package routes

import (
	"api-customer-merchant/internal/api/handlers"
	"api-customer-merchant/internal/db/repositories"
	"api-customer-merchant/internal/middleware"
	"api-customer-merchant/internal/services/dispute"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupDisputeRoutes(r *gin.Engine) {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logger.Sync() // Ensure logger flushes logs
	disputeRepo := repositories.NewDisputeRepository()
	orderRepo := repositories.NewOrderRepository()
	disputeService := dispute.NewDisputeService(disputeRepo,orderRepo,logger)
	disputeHandler := handlers.NewDisputeHandler(disputeService)

	disputeGroup := r.Group("/disputes")
	protected := disputeGroup.Use(middleware.AuthMiddleware("customer"))
	protected.POST("", disputeHandler.CreateDispute)
	protected.GET("/:id", disputeHandler.GetDispute)
}