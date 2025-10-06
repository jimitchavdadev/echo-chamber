package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zoro/echo-chamber/backend/internal/database"
	"github.com/zoro/echo-chamber/backend/internal/models"
	"github.com/zoro/echo-chamber/backend/internal/websocket"
)

func LikePost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, _ := strconv.ParseUint(postIDStr, 10, 32)
	currentUserID, _ := c.Get("userID")

	// Find the post to get the author's ID
	var post models.Post
	if err := database.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	like := models.Like{
		PostID: uint(postID),
		UserID: currentUserID.(uint),
	}
	
	result := database.DB.Where(like).FirstOrCreate(&like)

	// If the like was newly created (not found before)
	if result.RowsAffected > 0 && post.UserID != currentUserID.(uint) {
		notification := models.Notification{
			UserID:   post.UserID, // Notify the author of the post
			ActorID:  currentUserID.(uint),
			Type:     models.NotificationTypeLike,
			EntityID: uint(postID),
		}
		database.DB.Create(&notification)
		websocket.HubInstance.SendNotification(&notification)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post liked successfully"})
}

func UnlikePost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, _ := strconv.ParseUint(postIDStr, 10, 32)
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
	// Also delete the corresponding notification
	database.DB.Where("type = ? AND actor_id = ? AND entity_id = ?", models.NotificationTypeLike, currentUserID, postID).Delete(&models.Notification{})


	c.JSON(http.StatusOK, gin.H{"message": "Post unliked successfully"})
}
