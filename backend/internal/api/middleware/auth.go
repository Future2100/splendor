package middleware

import (
	"net/http"
	"strings"

	"splendor-backend/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Check Bearer format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format. Use 'Bearer <token>'"})
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		claims, err := jwt.ValidateToken(token, jwtSecret)
		if err != nil {
			if err == jwt.ErrExpiredToken {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			}
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}

// OptionalAuthMiddleware validates JWT tokens but doesn't require them
func OptionalAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			token := parts[1]
			claims, err := jwt.ValidateToken(token, jwtSecret)
			if err == nil {
				c.Set("userID", claims.UserID)
				c.Set("username", claims.Username)
			}
		}

		c.Next()
	}
}
