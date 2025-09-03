package customer

   import (
       "api-customer-merchant/internal/api/customer/handlers"
       "api-customer-merchant/internal/middleware"

       "github.com/gin-gonic/gin"
   )

   func RegisterRoutes(r *gin.Engine) {
       customer := r.Group("/customer")
       {
           authHandler := handlers.NewAuthHandler()
           customer.POST("/register", authHandler.Register)
           customer.POST("/login", authHandler.Login)
           customer.GET("/auth/google", authHandler.GoogleAuth)
           customer.GET("/auth/google/callback", authHandler.GoogleCallback)

           protected := customer.Group("/")
           protected.Use(middleware.AuthMiddleware("customer"))
           protected.POST("/logout", authHandler.Logout)
       }
   }