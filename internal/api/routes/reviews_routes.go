package routes

import (
	"api-customer-merchant/internal/api/handlers"
	"api-customer-merchant/internal/db/repositories"
	"api-customer-merchant/internal/middleware"
	"api-customer-merchant/internal/services/review"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupReviewRoutes(r *gin.Engine) {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	revrepo := repositories.NewReviewRepository()
	ordrepo:=repositories.NewOrderRepository()
	service := review.NewReviewService(revrepo,ordrepo, logger)
	//customer := r.Group("/products")
	{
		reviewHandler := handlers.NewReviewHandler(service)
		r.GET("/:productID/reviews",reviewHandler.GetReviewsByProduct)

		protected := r.Group("")
		protected.Use(middleware.AuthMiddleware("customer"))
		protected.POST("/review",reviewHandler.CreateReview)
		protected.GET("/reviews/:id", reviewHandler.GetReview)
		//protected.GET("/reviews/:id", reviewHandler.GetReview)
		//protected.GET("/:productID/reviews",reviewHandler.GetReviewsByProduct)
		protected.GET("/reviews",reviewHandler.GetAllUserReviews)
		protected.PUT("/reviews/:id", reviewHandler.UpdateReview)
		protected.DELETE("/reviews/:id", reviewHandler.DeleteReview)
	}
}
