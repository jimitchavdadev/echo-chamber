package api

import (
	"net/http"
	"strconv"

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

// FollowUser creates a follow relationship
func FollowUser(c *gin.Context) {
	currentUserID, _ := c.Get("userID")
	
	targetUserIDStr := c.Param("id")
	targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if currentUserID.(uint) == uint(targetUserID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You cannot follow yourself"})
		return
	}

	// Check if user exists
	var targetUser models.User
	if err := database.DB.First(&targetUser, targetUserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User to follow not found"})
		return
	}

	follow := models.Follower{
		FollowerID:  currentUserID.(uint),
		FollowingID: uint(targetUserID),
	}

	// Use FirstOrCreate to prevent duplicate follow entries
	result := database.DB.Where(follow).FirstOrCreate(&follow)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not follow user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully followed user"})
}

// UnfollowUser removes a follow relationship
func UnfollowUser(c *gin.Context) {
	currentUserID, _ := c.Get("userID")

	targetUserIDStr := c.Param("id")
	targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	follow := models.Follower{
		FollowerID:  currentUserID.(uint),
		FollowingID: uint(targetUserID),
	}

	result := database.DB.Where(follow).Delete(&models.Follower{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not unfollow user"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "You are not following this user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully unfollowed user"})
}

type UserProfileResponse struct {
	ID          uint   `json:"id"`
	Username    string `json:"username"`
	Bio         string `json:"bio"`
	IsFollowing bool   `json:"isFollowing"`
}

// GetUserByUsername finds a user by their username for public profiles
func GetUserByUsername(c *gin.Context) {
	username := c.Param("username")
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	response := UserProfileResponse{
		ID:       user.ID,
		Username: user.Username,
		Bio:      user.Bio,
	}

	// Check if the request is from an authenticated user
	currentUserID, exists := c.Get("userID")
	if exists {
		// If authenticated, check the follow status
		var follow models.Follower
		err := database.DB.Where("follower_id = ? AND following_id = ?", currentUserID, user.ID).First(&follow).Error
		response.IsFollowing = err == nil // isFollowing is true if a record was found
	}

	c.JSON(http.StatusOK, response)
}
