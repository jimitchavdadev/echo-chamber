package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zoro/echo-chamber/backend/internal/database"
	"github.com/zoro/echo-chamber/backend/internal/models"
	"github.com/zoro/echo-chamber/backend/internal/websocket"
)

type CommentInput struct {
	Content string `json:"content" binding:"required"`
}

func CreateComment(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, _ := strconv.ParseUint(postIDStr, 10, 32)
	currentUserID, _ := c.Get("userID")

	var post models.Post
	if err := database.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	
	var input CommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment := models.Comment{
		Content: input.Content,
		PostID:  uint(postID),
		UserID:  currentUserID.(uint),
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	// Send notification if not commenting on own post
	if post.UserID != currentUserID.(uint) {
		notification := models.Notification{
			UserID: post.UserID,
			ActorID: currentUserID.(uint),
			Type: models.NotificationTypeComment,
			EntityID: uint(postID),
		}
		database.DB.Create(&notification)
		websocket.HubInstance.SendNotification(&notification)
	}


	database.DB.Preload("User").First(&comment, comment.ID)
	c.JSON(http.StatusCreated, comment)
}

func GetCommentsForPost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, _ := strconv.ParseUint(postIDStr, 10, 32)

	var comments []models.Comment
	database.DB.Preload("User").Where("post_id = ?", postID).Order("created_at asc").Find(&comments)

	c.JSON(http.StatusOK, comments)
}
