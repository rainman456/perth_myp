package utils

import (
	"api-customer-merchant/internal/config"
	"context"
	"crypto/tls"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis(conf *config.Config) {
    
    RedisClient = redis.NewClient(&redis.Options{
        Addr:     conf.RedisAddr,
        Password: conf.RedisPass,
        DB:       conf.RedisDB,
        TLSConfig: &tls.Config{
            MinVersion: tls.VersionTLS12, 
            InsecureSkipVerify: true, // Use with caution, only for testing
        },
    })

    ctx := context.Background()
    if err := RedisClient.Ping(ctx).Err(); err != nil {
       log.Printf("Failed to connect to Redis: %v, continuing without caching", err)
        RedisClient = nil // Fallback to avoid crashes
    } else {
        log.Println("Connected to Redis successfully")
    }
}

// Helper to get cached value or fetch and cache
func GetOrSetCache(ctx context.Context, key string, ttl time.Duration, fetch func() (any, error)) (any, error) {
    val, err := RedisClient.Get(ctx, key).Result()
    if err == nil {
        return val, nil // Deserialize if needed (e.g., JSON)
    }
    if err != redis.Nil {
        return nil, err
    }

    data, err := fetch()
    if err != nil {
        return nil, err
    }

    // Serialize if complex (e.g., JSON marshal)
    if err := RedisClient.Set(ctx, key, data, ttl).Err(); err != nil {
        return nil, err
    }
    return data, nil
}