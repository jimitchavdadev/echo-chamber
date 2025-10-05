package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zoro/echo-chamber/backend/internal/database"
	"github.com/zoro/echo-chamber/backend/internal/models"
)

type PostInput struct {
	Content string `json:"content" binding:"required"`
}

func CreatePost(c *gin.Context) {
	currentUserID, _ := c.Get("userID")

	var input PostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post := models.Post{
		Content: input.Content,
		UserID:  currentUserID.(uint),
	}

	result := database.DB.Create(&post)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	// Preload user to return the author's details
	database.DB.Preload("User").First(&post, post.ID)

	c.JSON(http.StatusCreated, post)
}

func GetPost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var post models.Post
	if err := database.DB.Preload("User").First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Calculate like count
	database.DB.Model(&models.Like{}).Where("post_id = ?", post.ID).Count(&post.LikeCount)

	// Check if the current user (if any) has liked this post
	currentUserID, exists := c.Get("userID")
	if exists {
		var like models.Like
		if err := database.DB.Where("post_id = ? AND user_id = ?", post.ID, currentUserID).First(&like).Error; err == nil {
			post.IsLiked = true
		}
	}

	c.JSON(http.StatusOK, post)
}

func DeletePost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	currentUserID, _ := c.Get("userID")

	var post models.Post
	if err := database.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	if post.UserID != currentUserID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this post"})
		return
	}

	// Using Unscoped().Delete() for a permanent delete
	database.DB.Unscoped().Delete(&post)

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
