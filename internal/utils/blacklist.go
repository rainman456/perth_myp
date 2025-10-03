package utils

import (
	"context"
	"errors"
	"log"
	"time"
	//"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// Add adds a token to the Redis blacklist with an expiration
func Add(token string) error {
	if RedisClient == nil {
		log.Println("RedisClient is nil, cannot add token to blacklist")
		return errors.New("redis client not initialized")
	}
	// Set token in Redis with a 24-hour expiration
	return RedisClient.Set(ctx, "blacklist:"+token, "true", 24*time.Hour).Err()
}

// IsBlacklisted checks if a token is in the Redis blacklist
func IsBlacklisted(token string) bool {
	if RedisClient == nil {
		log.Println("RedisClient is nil, skipping blacklist check")
		return false // Fallback to allow operation if Redis is unavailable
	}
	if RedisClient == nil {
		return false
	}
	_, err := RedisClient.Get(ctx, "blacklist:"+token).Result() // Line 22
	return err == nil                                           // Token exists in Redis if no error
}
