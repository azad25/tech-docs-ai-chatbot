package app

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// WebSocketHandler handles WebSocket connections for real-time chat
type WebSocketHandler struct {
	service  ServiceInterface
	upgrader websocket.Upgrader
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(service ServiceInterface) *WebSocketHandler {
	return &WebSocketHandler{
		service: service,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Allow connections from any origin in development
				// In production, you should validate the origin
				return true
			},
		},
	}
}

// WebSocketMessage represents a message sent over WebSocket
type WebSocketMessage struct {
	Type      string `json:"type"`
	SessionID string `json:"session_id,omitempty"`
	Message   string `json:"message,omitempty"`
	Response  string `json:"response,omitempty"`
	Error     string `json:"error,omitempty"`
}

// HandleWebSocket handles WebSocket connections
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	log.Println("WebSocket connection established")

	for {
		var msg WebSocketMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		switch msg.Type {
		case "chat":
			// Send typing indicator
			h.sendTypingIndicator(conn)
			
			response, err := h.service.Chat(msg.Message)
			if err != nil {
				h.sendError(conn, "Failed to process chat message", err)
				continue
			}
			h.sendResponse(conn, "chat_response", response)

		case "chat_with_history":
			if msg.SessionID == "" {
				h.sendError(conn, "Session ID is required for chat with history", nil)
				continue
			}
			
			// Send typing indicator
			h.sendTypingIndicator(conn)
			
			response, err := h.service.ChatWithHistory(msg.SessionID, msg.Message)
			if err != nil {
				h.sendError(conn, "Failed to process chat with history", err)
				continue
			}
			h.sendResponse(conn, "chat_response", response)

		default:
			h.sendError(conn, "Unknown message type", nil)
		}
	}
}

// sendResponse sends a successful response over WebSocket
func (h *WebSocketHandler) sendResponse(conn *websocket.Conn, msgType, response string) {
	msg := WebSocketMessage{
		Type:     msgType,
		Response: response,
	}
	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("WebSocket write error: %v", err)
	}
}

// sendError sends an error response over WebSocket
func (h *WebSocketHandler) sendError(conn *websocket.Conn, message string, err error) {
	errorMsg := message
	if err != nil {
		errorMsg += ": " + err.Error()
	}

	msg := WebSocketMessage{
		Type:  "error",
		Error: errorMsg,
	}
	if writeErr := conn.WriteJSON(msg); writeErr != nil {
		log.Printf("WebSocket write error: %v", writeErr)
	}
}

// sendTypingIndicator sends a typing indicator over WebSocket
func (h *WebSocketHandler) sendTypingIndicator(conn *websocket.Conn) {
	msg := WebSocketMessage{
		Type: "typing",
	}
	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("WebSocket write error: %v", err)
	}
}