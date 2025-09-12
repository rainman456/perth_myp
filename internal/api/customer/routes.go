package customer

   import (
       "api-customer-merchant/internal/api/customer/handlers"
       "api-customer-merchant/internal/middleware"
       "api-customer-merchant/internal/db/repositories"
       "api-customer-merchant/internal/domain/user"


       "github.com/gin-gonic/gin"
   )

   func RegisterRoutes(r *gin.Engine) {
    repo := repositories.NewUserRepository()
    service := user.NewAuthService( repo)
       customer := r.Group("/customer")
       {
           authHandler := handlers.NewAuthHandler(service)
           customer.POST("/register", authHandler.Register)
           customer.POST("/login", authHandler.Login)
           customer.GET("/auth/google", authHandler.GoogleAuth)
           customer.GET("/auth/google/callback", authHandler.GoogleCallback)

           protected := customer.Group("/")
           protected.Use(middleware.AuthMiddleware("customer"))
           protected.POST("/logout", authHandler.Logout)
       }
   }