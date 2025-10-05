package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zoro/echo-chamber/backend/internal/database"
	"github.com/zoro/echo-chamber/backend/internal/models"
)

// GetConversations retrieves a list of users the current user has chatted with
func GetConversations(c *gin.Context) {
	currentUserID, _ := c.Get("userID")

	// This new query is simpler and avoids GORM's soft-delete confusion.
	// 1. Get the distinct IDs of conversation partners.
	var partnerIDs []uint
	database.DB.Model(&models.ChatMessage{}).
		Select("CASE WHEN sender_id = ? THEN receiver_id ELSE sender_id END", currentUserID).
		Where("sender_id = ? OR receiver_id = ?", currentUserID, currentUserID).
		Distinct().
		Pluck("CASE WHEN sender_id = ? THEN receiver_id ELSE sender_id END", &partnerIDs)

	if len(partnerIDs) == 0 {
		c.JSON(http.StatusOK, []models.User{})
		return
	}

	// 2. Fetch the user details for those partner IDs.
	var users []models.User
	database.DB.Where("id IN ?", partnerIDs).Find(&users)

	c.JSON(http.StatusOK, users)
}


// GetChatHistory retrieves the message history between the current user and another user
func GetChatHistory(c *gin.Context) {
	currentUserID, _ := c.Get("userID")
	otherUserIDStr := c.Param("userId")
	otherUserID, _ := strconv.ParseUint(otherUserIDStr, 10, 32)
	
	var messages []models.ChatMessage
	database.DB.Preload("Sender").Preload("Receiver").
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			currentUserID, otherUserID, otherUserID, currentUserID).
		Order("created_at asc").
		Find(&messages)
		
	c.JSON(http.StatusOK, messages)
}
