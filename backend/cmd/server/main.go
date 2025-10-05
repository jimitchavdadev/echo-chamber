package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/zoro/echo-chamber/backend/internal/api"
	"github.com/zoro/echo-chamber/backend/internal/database"
)

func main() {
	// Initialize Database
	database.ConnectDatabase()

	r := gin.Default()

	// CORS Middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Public routes
	public := r.Group("/api")
	{
		public.POST("/register", api.Register)
		public.POST("/login", api.Login)
		public.POST("/logout", api.Logout)
		public.GET("/users/:username", api.GetUserByUsername)
	}

	// Protected routes
	protected := r.Group("/api")
	protected.Use(api.AuthMiddleware())
	{
		protected.GET("/me", api.GetCurrentUser)
		protected.PUT("/profile", api.UpdateCurrentUserProfile)
		protected.POST("/users/:id/follow", api.FollowUser)
		protected.DELETE("/users/:id/unfollow", api.UnfollowUser)
		protected.GET("/feed", api.GetFeed)
		protected.POST("/posts", api.CreatePost)
	}

	log.Println("Starting server on port 8080...")
	r.Run(":8080")
}
