// backend/cmd/server/main.go
package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run() // listens and serves on 0.0.0.0:8080 by default
}
