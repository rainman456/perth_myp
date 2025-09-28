package routes

/*

import (
	"api-customer-merchant/internal/api/handlers"
	"api-customer-merchant/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterPaymentRoutes(r *gin.Engine, paymentHandler *handlers.PaymentHandler) {
	payment := r.Group("/payments")
	payment.Use(middleware.RateLimitMiddleware())
	{
		payment.POST("/initialize", paymentHandler.Initialize)
		payment.GET("/verify/:reference", paymentHandler.Verify)
		payment.POST("/webhook", paymentHandler.Webhook) // No auth for webhook
	}
}

*/