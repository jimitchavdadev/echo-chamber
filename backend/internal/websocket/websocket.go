package websocket

import "encoding/json"

// HubInstance is the global hub that manages all websocket connections.
var HubInstance = newHub()

// WsMessage defines the structure for messages sent over WebSocket.
type WsMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// ChatMessagePayload is the specific payload for chat messages.
type ChatMessagePayload struct {
	ReceiverID uint   `json:"receiverId"`
	Content    string `json:"content"`
}
