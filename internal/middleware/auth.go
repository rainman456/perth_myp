package middleware

import (
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
		id := claims["id"].(string)
		switch entityType {
		case "user":
			c.Set("userID", id)
		case "merchant":
			c.Set("merchantID", id)
		}
		c.Next()
	}
}
