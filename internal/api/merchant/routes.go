package merchant

    import (
       "api-customer-merchant/internal/api/merchant/handlers"
       "api-customer-merchant/internal/middleware"

        "github.com/gin-gonic/gin"
    )

    func RegisterRoutes(r *gin.Engine) {
        merchant := r.Group("/merchant")
        {
            authHandler := handlers.NewAuthHandler()
            merchant.POST("/submitApplcation", authHandler.Register)
            merchant.POST("/login", authHandler.Login)
            //merchant.GET("/auth/google", authHandler.GoogleAuth)
            //merchant.GET("/auth/google/callback", authHandler.GoogleCallback)

            protected := merchant.Group("/")
            protected.Use(middleware.AuthMiddleware("merchant"))
            protected.POST("/logout", authHandler.Logout)
        }
    }