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
	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)
	
	// 1. Find all users the current user is following
	var followingIDs []uint
	database.DB.Model(&models.Follower{}).Where("follower_id = ?", currentUserID).Pluck("following_id", &followingIDs)
	followingIDs = append(followingIDs, currentUserID.(uint)) // Include user's own posts

	// 2. Fetch a page of posts from those users
	var posts []models.Post
	database.DB.Preload("User").
		Where("user_id IN ?", followingIDs).
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&posts)

	if len(posts) == 0 {
		c.JSON(http.StatusOK, []models.Post{})
		return
	}

	// 3. Efficiently fetch like counts and liked status for the retrieved posts
	var postIDs []uint
	for _, post := range posts {
		postIDs = append(postIDs, post.ID)
	}

	// Get all like counts for the posts in one query
	type LikeCountResult struct {
		PostID uint
		Count  int64
	}
	var likeCounts []LikeCountResult
	database.DB.Model(&models.Like{}).Select("post_id, count(*) as count").Where("post_id IN ?", postIDs).Group("post_id").Scan(&likeCounts)
	
	// Create a map for easy lookup
	likeCountMap := make(map[uint]int64)
	for _, result := range likeCounts {
		likeCountMap[result.PostID] = result.Count
	}

	// Get all likes by the current user for the posts in one query
	var userLikes []models.Like
	database.DB.Where("post_id IN ? AND user_id = ?", postIDs, currentUserID).Find(&userLikes)
	
	// Create a set for easy lookup
	userLikedPostIDs := make(map[uint]bool)
	for _, like := range userLikes {
		userLikedPostIDs[like.PostID] = true
	}

	// 4. Attach the calculated counts and statuses to the post objects
	for i := range posts {
		posts[i].LikeCount = likeCountMap[posts[i].ID]
		posts[i].IsLiked = userLikedPostIDs[posts[i].ID]
	}

	c.JSON(http.StatusOK, posts)
}
