package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/zoro/echo-chamber/backend/internal/api"
	"github.com/zoro/echo-chamber/backend/internal/database"
	"github.com/zoro/echo-chamber/backend/internal/websocket"
)

func main() {
	database.ConnectDatabase()
	
	// Run the WebSocket hub in a separate goroutine
	go websocket.HubInstance.Run()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	apiGroup := r.Group("/api")
	apiGroup.POST("/register", api.Register)
	apiGroup.POST("/login", api.Login)
	apiGroup.POST("/logout", api.Logout)
	apiGroup.GET("/posts/:id/comments", api.GetCommentsForPost)
	
	apiGroup.Use(api.MaybeAuthMiddleware())
	{
		apiGroup.GET("/users/:username", api.GetUserByUsername)
		apiGroup.GET("/posts/:id", api.GetPost)
	}

	protected := apiGroup.Group("/")
	protected.Use(api.AuthMiddleware())
	{
		protected.GET("/me", api.GetCurrentUser)
		protected.PUT("/profile", api.UpdateCurrentUserProfile)
		protected.POST("/users/:id/follow", api.FollowUser)
		protected.DELETE("/users/:id/unfollow", api.UnfollowUser)
		protected.GET("/feed", api.GetFeed)
		protected.POST("/posts", api.CreatePost)
		protected.DELETE("/posts/:id", api.DeletePost)
		protected.POST("/posts/:id/comments", api.CreateComment)
		protected.POST("/posts/:id/like", api.LikePost)
		protected.DELETE("/posts/:id/like", api.UnlikePost)
		// Add the WebSocket endpoint
		protected.GET("/ws", api.ServeWs)
		protected.GET("/chat/conversations", api.GetConversations)
		protected.GET("/chat/history/:userId", api.GetChatHistory)
		protected.GET("/users/search", api.SearchUsers)
	}

	log.Println("Starting server on port 8080...")
	r.Run(":8080")
}
