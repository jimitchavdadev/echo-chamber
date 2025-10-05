package api

import (
	"net/http"

	"github.com/zoro/echo-chamber/backend/internal/auth"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware requires a valid token and aborts if it's not present.
// Used for actions that REQUIRE a logged-in user (e.g., posting, liking).
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := auth.GetUserIDFromToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}

// MaybeAuthMiddleware checks for a token, but does not fail if it's missing.
// It simply sets the userID in the context if the user is logged in.
// Used for public endpoints that show extra info for logged-in users (e.g., isLiked status).
func MaybeAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := auth.GetUserIDFromToken(c)
		if err == nil {
			// User is logged in, set their ID in the context
			c.Set("userID", userID)
		}
		// Continue with the request regardless of auth status
		c.Next()
	}
}
