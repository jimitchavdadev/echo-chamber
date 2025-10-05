package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zoro/echo-chamber/backend/internal/database"
	"github.com/zoro/echo-chamber/backend/internal/models"
)

func LikePost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	currentUserID, _ := c.Get("userID")

	like := models.Like{
		PostID: uint(postID),
		UserID: currentUserID.(uint),
	}

	// Use FirstOrCreate. This will find an existing like or create a new one.
	// We will ignore the error here because the most common error is a "duplicate key"
	// violation, which is acceptable in this case (the post is already liked).
	database.DB.Where(like).FirstOrCreate(&like)

	c.JSON(http.StatusOK, gin.H{"message": "Post liked successfully"})
}

func UnlikePost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	currentUserID, _ := c.Get("userID")

	like := models.Like{
		PostID: uint(postID),
		UserID: currentUserID.(uint),
	}

	result := database.DB.Where(like).Delete(&models.Like{})
	if result.Error != nil || result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "You have not liked this post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post unliked successfully"})
}
