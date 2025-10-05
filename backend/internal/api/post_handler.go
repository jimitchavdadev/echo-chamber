package api

import (
	"net/http"

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
