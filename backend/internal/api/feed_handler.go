package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zoro/echo-chamber/backend/internal/database"
	"github.com/zoro/echo-chamber/backend/internal/models"
)

// GetFeed retrieves posts from users that the current user follows
func GetFeed(c *gin.Context) {
	currentUserID, _ := c.Get("userID")

	// --- Pagination ---
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}
	// ------------------

	// 1. Find all users the current user is following
	var following []models.Follower
	database.DB.Where("follower_id = ?", currentUserID).Find(&following)

	var followingIDs []uint
	for _, f := range following {
		followingIDs = append(followingIDs, f.FollowingID)
	}
	// Also include the user's own posts in their feed
	followingIDs = append(followingIDs, currentUserID.(uint))

	// 2. Fetch posts from those users
	var posts []models.Post
	result := database.DB.Preload("User").
		Where("user_id IN ?", followingIDs).
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&posts)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve feed"})
		return
	}

	c.JSON(http.StatusOK, posts)
}
