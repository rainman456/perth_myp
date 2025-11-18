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
    if RedisClient == nil {
        log.Printf("redis not initialized, skipping cache for key %s", key)
        return fetch()
    }

    val, err := RedisClient.Get(ctx, key).Result()
    if err == nil {
        return val, nil
    }
    if err != redis.Nil {
        log.Printf("redis GET error for key %s: %v", key, err)
    }

    data, err := fetch()
    if err != nil {
        return nil, err
    }

    go func(k string, d any) {
        if RedisClient == nil {
            return
        }
        _ = RedisClient.Set(context.Background(), k, d, ttl).Err()
    }(key, data)

    return data, nil
}


func GetOrSetCacheJSON[T any](ctx context.Context, key string, ttl time.Duration, fetch func() (T, error)) (T, error) {
    var result T

    // If Redis isn't initialized, just fetch and return (no cache)
    if RedisClient == nil {
        // Optional: log once for visibility (don't spam)
        log.Printf("redis not initialized, skipping cache for key %s", key)
        return fetch()
    }

    // Try to get from cache
    val, err := RedisClient.Get(ctx, key).Result()
    if err == nil {
        if err := json.Unmarshal([]byte(val), &result); err == nil {
            return result, nil
        }
        // If unmarshal fails, fallthrough to fetch fresh
        log.Printf("failed to unmarshal cached value for key %s: %v", key, err)
    } else if err != redis.Nil {
        // Real Redis error — log and continue to fetch fresh data
        log.Printf("redis GET error for key %s: %v", key, err)
    }

    // Cache miss or error — fetch fresh data
    result, err = fetch()
    if err != nil {
        return result, err
    }

    // Store in cache asynchronously (best-effort)
    go func(data T) {
        if RedisClient == nil {
            return
        }
        jsonData, merr := json.Marshal(data)
        if merr != nil {
            return
        }
        _ = RedisClient.Set(context.Background(), key, jsonData, ttl).Err()
    }(result)

    return result, nil
}


// InvalidateCache - Delete single key
// FIXED: Handle nil RedisClient
func InvalidateCache(ctx context.Context, key string) error {
	if RedisClient == nil {
		log.Printf("redis not initialized, skipping cache invalidation for key %s", key)
		return nil // Not an error - just no-op when Redis is unavailable
	}
	return RedisClient.Del(ctx, key).Err()
}

// InvalidateCachePattern - Delete keys matching pattern
// FIXED: Handle nil RedisClient
func InvalidateCachePattern(ctx context.Context, pattern string) error {
	if RedisClient == nil {
		log.Printf("redis not initialized, skipping cache pattern invalidation for pattern %s", pattern)
		return nil // Not an error - just no-op when Redis is unavailable
	}

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