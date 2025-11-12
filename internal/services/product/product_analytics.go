package product


/*
import (
	"context"

	"github.com/redis/go-redis/v9"
)


type ProductAnalytics struct {
    redis *redis.Client
}

func (s *ProductService) TrackProductView(ctx context.Context, productID string, userID uint) {
    // Increment view count (sorted set for trending)
    s.redis.ZIncrBy(ctx, "product:views:daily", 1, productID)
    
    // Store in time-series for analytics
    day := time.Now().Format("2006-01-02")
    s.redis.HIncrBy(ctx, fmt.Sprintf("product:views:%s", day), productID, 1)
    
    // Set expiry on daily keys (30 days)
    s.redis.Expire(ctx, fmt.Sprintf("product:views:%s", day), 30*24*time.Hour)
}

func (s *ProductService) GetTrendingProducts(ctx context.Context, limit int) ([]string, error) {
    // Get top viewed products from last 24 hours
    return s.redis.ZRevRange(ctx, "product:views:daily", 0, int64(limit-1)).Result()
}

// Use in GetProductByID
func (s *ProductService) GetProductByID(ctx context.Context, id string, userID uint, preloads ...string) (*dto.ProductResponse, error) {
    // Track view asynchronously
    go s.TrackProductView(context.Background(), id, userID)
    
    // ... rest of code
}

*/