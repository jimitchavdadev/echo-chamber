package websocket

import (
	"encoding/json"
	"log"

	"github.com/zoro/echo-chamber/backend/internal/database"
	"github.com/zoro/echo-chamber/backend/internal/models"
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	clients    map[uint]*Client
	Register   chan *Client
	Unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[uint]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.clients[client.UserID] = client
			log.Printf("Client registered: UserID %d", client.UserID)
		case client := <-h.Unregister:
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
				log.Printf("Client unregistered: UserID %d", client.UserID)
			}
		}
	}
}

// SendNotification sends a notification to a specific user if they are online.
func (h *Hub) SendNotification(notification *models.Notification) {
	client, ok := h.clients[notification.UserID]
	if !ok {
		log.Printf("User %d is not online, notification not sent in real-time.", notification.UserID)
		return
	}

	database.DB.Preload("Actor").First(&notification, notification.ID)
	
	payloadBytes, err := json.Marshal(notification)
	if err != nil {
		log.Printf("error marshalling notification: %v", err)
		return
	}

	msg := WsMessage{
		Type:    "new_notification",
		Payload: json.RawMessage(payloadBytes),
	}
	
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("error marshalling ws message: %v", err)
		return
	}

	client.Send <- msgBytes
}

func (h *Hub) handleChatMessage(senderID uint, message []byte) {
	var payload ChatMessagePayload
	if err := json.Unmarshal(message, &payload); err != nil {
		log.Printf("error unmarshalling chat message: %v", err)
		return
	}

	chatMsg := models.ChatMessage{
		SenderID:   senderID,
		ReceiverID: payload.ReceiverID,
		Content:    payload.Content,
	}

	if err := database.DB.Create(&chatMsg).Error; err != nil {
		log.Printf("error saving chat message to DB: %v", err)
		return
	}

	receiverClient, ok := h.clients[payload.ReceiverID]
	if ok {
		// This is the fully corrected line
		database.DB.Preload("Sender").Preload("Receiver").First(&chatMsg, chatMsg.ID)
		
		payloadBytes, err := json.Marshal(chatMsg)
		if err != nil {
			log.Printf("error marshalling chat message for delivery: %v", err)
			return
		}
		
		wsMsg := WsMessage{
			Type: "new_chat_message",
			Payload: json.RawMessage(payloadBytes),
		}

		msgBytes, err := json.Marshal(wsMsg)
		if err != nil {
			log.Printf("error marshalling ws message for chat: %v", err)
			return
		}
		receiverClient.Send <- msgBytes
	}
}
