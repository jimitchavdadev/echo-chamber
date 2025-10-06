package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	ws "github.com/zoro/echo-chamber/backend/internal/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ServeWs handles websocket requests from the peer.
func ServeWs(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	userID, _ := c.Get("userID")

	client := &ws.Client{
		Hub:    ws.HubInstance,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		UserID: userID.(uint),
	}
	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
