package utils

import (
	"api-customer-merchant/internal/config"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
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
			MinVersion:         tls.VersionTLS12,
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


func GetOrSetCacheJSON[T any](ctx context.Context, key string, ttl time.Duration, fetch func() (T, error)) (T, error) {
	var result T

	// Try to get from cache
	val, err := RedisClient.Get(ctx, key).Result()
	if err == nil {
		// Cache hit - deserialize
		if err := json.Unmarshal([]byte(val), &result); err == nil {
			return result, nil
		}
		// If unmarshal fails, continue to fetch fresh data
	}

	if err != nil && err != redis.Nil {
		// Redis error (not just cache miss)
		// Continue to fetch from DB but log the error
	}

	// Cache miss or error - fetch fresh data
	result, err = fetch()
	if err != nil {
		return result, err
	}

	// Store in cache (async to not block response)
	go func() {
		jsonData, err := json.Marshal(result)
		if err != nil {
			return
		}
		RedisClient.Set(context.Background(), key, jsonData, ttl).Err()
	}()

	return result, nil
}

// InvalidateCache - Delete single key
func InvalidateCache(ctx context.Context, key string) error {
	return RedisClient.Del(ctx, key).Err()
}

// InvalidateCachePattern - Delete keys matching pattern
func InvalidateCachePattern(ctx context.Context, pattern string) error {
	iter := RedisClient.Scan(ctx, 0, pattern, 0).Iterator()
	keys := []string{}

	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return err
	}

	if len(keys) > 0 {
		return RedisClient.Del(ctx, keys...).Err()
	}

	return nil
}

// Cache key generators
func ProductCacheKey(productID string) string {
	return "product:details:" + productID
}

func ProductListCacheKey(page, limit int, filters string) string {
	return fmt.Sprintf("product:list:p%d:l%d:%s", page, limit, filters)
}

func ProductSearchCacheKey(query string, limit int) string {
	return fmt.Sprintf("product:search:%s:l%d", query, limit)
}