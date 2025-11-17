package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"api-customer-merchant/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(entityType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if utils.IsBlacklisted(tokenString) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is blacklisted"})
			c.Abort()
			return
		}
		key := os.Getenv("JWT_SECRET")

		secret := []byte(key) // Load from env
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["entityType"] != entityType {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid entity type"})
			c.Abort()
			return
		}

		//c.Set("entityId", claims["id"])
		idInterface := claims["id"]
		id := fmt.Sprintf("%v", idInterface) // Convert to string regardless of type (handles float64, string, etc.)
		switch entityType {
		case "user", "customer":
			c.Set("userID", id)
		case "merchant":
			c.Set("merchantID", id)
		}
		c.Next()
	}
}



// OptionalAuthMiddleware is similar to AuthMiddleware but doesn't require authentication
// It sets user info if available but doesn't block if not authenticated
func OptionalAuthMiddleware(entityType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// Try to get token from Authorization header first
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// If no Authorization header, try to get token from cookie
			cookie, err := c.Cookie("auth_token")
			if err != nil {
				// No authentication found, continue without setting user
				c.Next()
				return
			}
			tokenString = cookie
		}

		// Check if token is blacklisted
		if utils.IsBlacklisted(tokenString) {
			// Blacklisted token, continue without authentication
			c.Next()
			return
		}

		// Validate JWT token
		key := os.Getenv("JWT_SECRET")
		secret := []byte(key)
		
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
			// Invalid token, continue without authentication
			c.Next()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["entityType"] != entityType {
			// Invalid entity type, continue without authentication
			c.Next()
			return
		}

		// Extract and set user/merchant ID
		idInterface := claims["id"]
		id := fmt.Sprintf("%v", idInterface)
		
		switch entityType {
		case "user", "customer":
			c.Set("userID", id)
		case "merchant":
			c.Set("merchantID", id)
		}
		
		c.Next()
	}
}