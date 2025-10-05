package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/zoro/echo-chamber/backend/internal/api"
	"github.com/zoro/echo-chamber/backend/internal/database"
)

func main() {
	database.ConnectDatabase()
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// --- Route Groups ---
	
	// Group for all API routes
	apiGroup := r.Group("/api")

	// Public routes that DO NOT need to know about the user
	apiGroup.POST("/register", api.Register)
	apiGroup.POST("/login", api.Login)
	apiGroup.POST("/logout", api.Logout)
	apiGroup.GET("/posts/:id/comments", api.GetCommentsForPost)
	
	// Public routes that CAN be enhanced by knowing the user (optional auth)
	apiGroup.Use(api.MaybeAuthMiddleware())
	{
		apiGroup.GET("/users/:username", api.GetUserByUsername)
		apiGroup.GET("/posts/:id", api.GetPost)
	}

	// Protected routes that REQUIRE a user to be logged in
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
	}

	log.Println("Starting server on port 8080...")
	r.Run(":8080")
}
