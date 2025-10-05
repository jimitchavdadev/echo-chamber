package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zoro/echo-chamber/backend/internal/database"
	"github.com/zoro/echo-chamber/backend/internal/models"
)

type ProfileInput struct {
	Bio string `json:"bio"`
}

// UpdateCurrentUserProfile handles updating the logged-in user's profile
func UpdateCurrentUserProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var input ProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Bio = input.Bio
	database.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"user": gin.H{"id": user.ID, "username": user.Username, "email": user.Email, "bio": user.Bio}})
}
