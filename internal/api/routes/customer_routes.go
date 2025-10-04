package routes

import (
	"api-customer-merchant/internal/api/handlers"
	"api-customer-merchant/internal/db/repositories"
	"api-customer-merchant/internal/middleware"
	"api-customer-merchant/internal/services/user"

	"github.com/gin-gonic/gin"
)

func RegisterCustomerRoutes(r *gin.Engine) {
	repo := repositories.NewUserRepository()
	service := user.NewAuthService(repo)
	customer := r.Group("/customer")
	{
		authHandler := handlers.NewAuthHandler(service)
		customer.POST("/register", authHandler.Register)
		customer.POST("/login", authHandler.Login)
		customer.GET("/auth/google", authHandler.GoogleAuth)
		customer.GET("/auth/google/callback", authHandler.GoogleCallback)

		protected := customer.Group("/")
		protected.Use(middleware.AuthMiddleware("customer"))
		protected.PATCH("/update",authHandler.UpdateProfile)
		protected.POST("/logout", authHandler.Logout)
	}
}
