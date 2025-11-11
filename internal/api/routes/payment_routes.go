package routes

import (
	"os"

	"api-customer-merchant/internal/api/handlers"
	"api-customer-merchant/internal/config"
	"api-customer-merchant/internal/db/repositories"
	"api-customer-merchant/internal/middleware"
	"api-customer-merchant/internal/services/payment"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterPaymentRoutes(r *gin.Engine) {
	// Initialize logger (consider injecting instead of creating here for production)
	logger, err := zap.NewProduction()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	// Do NOT defer logger.Sync() here — sync on app shutdown instead

	// Load config (assume config has a way to load from env or file; adjust if needed)
	conf := &config.Config{
		PaystackSecretKey: os.Getenv("PAYSTACK_SECRET_KEY"),
		// Add other config fields as needed
	}

	// Initialize repositories (assuming they use a global db.DB or inject if needed)
	orderRepo := repositories.NewOrderRepository()
	merchantRepo := repositories.NewMerchantRepository()
	payoutRepo := repositories.NewPayoutRepository()
	
	paymentRepo := repositories.NewPaymentRepository()

	// Initialize service with all dependencies
	paymentService := payment.NewPaymentService(
		paymentRepo,
		orderRepo,
		payoutRepo,
		merchantRepo,
		conf,
		logger,
	)

	// Initialize handler
	paymentHandler := handlers.NewPaymentHandler(paymentService, conf, logger)

	payment := r.Group("/payments")

	// Webhook: No auth, no rate limit (Paystack retries on failure; avoid limiting)
	payment.POST("/webhook", paymentHandler.Webhook)

	// Protected customer routes with rate limiting
	protectedGroup := payment.Group("")
	protectedGroup.Use(middleware.AuthMiddleware("customer"), middleware.RateLimitMiddleware())
	{
		protectedGroup.POST("/initialize", paymentHandler.Initialize)
		protectedGroup.GET("/verify/:reference", paymentHandler.Verify)
	}
}