package api

import (
	"net/http"
	"github.com/zoro/echo-chamber/backend/internal/auth"
	"github.com/gin-gonic/gin"
)

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
