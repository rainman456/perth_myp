package routes



import (
	"api-customer-merchant/internal/api/handlers"
	"api-customer-merchant/internal/db/repositories"
	"api-customer-merchant/internal/middleware"
	"api-customer-merchant/internal/services/wishlist"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupWishlistRoutes(r *gin.Engine) {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	wishrepo := repositories.NewWishlistRepository()
	service := wishlist.NewWishlistService(wishrepo, logger)
	{
		wishlistHandler := handlers.NewWishlistHandler(service)

		protected := middleware.AuthMiddleware("customer") // Consider adding auth middleware
		r.POST("/wishlist/:productID", protected, wishlistHandler.AddToWishlist)
		r.DELETE("/wishlist/:productID", protected,   wishlistHandler.RemoveFromWishlist)
		r.GET("/wishlist", protected, wishlistHandler.GetWishlist)
		r.GET("/wishlist/:productID/check", protected, wishlistHandler.IsInWishlist)
		r.DELETE("/wishlist/clear", protected ,wishlistHandler.ClearWishlist)
	}
}