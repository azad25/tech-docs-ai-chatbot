package app

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocketHandler handles WebSocket connections for real-time chat
type WebSocketHandler struct {
	service  *Service
	upgrader websocket.Upgrader
	clients  map[*websocket.Conn]bool
	mutex    sync.RWMutex
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(svc *Service) *WebSocketHandler {
	return &WebSocketHandler{
		service: svc,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for development
			},
		},
		clients: make(map[*websocket.Conn]bool),
	}
}

// WebSocketMessage represents a WebSocket message
type WebSocketMessage struct {
	Type      string `json:"type"` // "chat", "typing", "status"
	Message   string `json:"message"`
	SessionID string `json:"session_id,omitempty"`
	UserType  string `json:"user_type,omitempty"`
}

// WebSocketResponse represents a WebSocket response
type WebSocketResponse struct {
	Type      string `json:"type"` // "chat_response", "status", "error"
	Message   string `json:"message"`
	SessionID string `json:"session_id,omitempty"`
	Error     string `json:"error,omitempty"`
}

// HandleWebSocket handles WebSocket connections
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to WebSocket
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	// Add client to the list
	h.mutex.Lock()
	h.clients[conn] = true
	h.mutex.Unlock()

	// Send welcome message
	welcomeMsg := WebSocketResponse{
		Type:    "status",
		Message: "Connected to Tech Docs AI WebSocket",
	}
	conn.WriteJSON(welcomeMsg)

	log.Printf("WebSocket client connected: %s", conn.RemoteAddr())

	// Handle incoming messages
	for {
		var msg WebSocketMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		// Handle different message types
		switch msg.Type {
		case "chat":
			h.handleChatMessage(conn, msg)
		case "typing":
			h.handleTypingMessage(conn, msg)
		default:
			// Send error response for unknown message type
			response := WebSocketResponse{
				Type:  "error",
				Error: "Unknown message type: " + msg.Type,
			}
			conn.WriteJSON(response)
		}
	}

	// Remove client from the list
	h.mutex.Lock()
	delete(h.clients, conn)
	h.mutex.Unlock()

	log.Printf("WebSocket client disconnected: %s", conn.RemoteAddr())
}

// handleChatMessage processes chat messages and sends responses
func (h *WebSocketHandler) handleChatMessage(conn *websocket.Conn, msg WebSocketMessage) {
	if msg.Message == "" {
		response := WebSocketResponse{
			Type:  "error",
			Error: "Empty message",
		}
		conn.WriteJSON(response)
		return
	}

	// Send typing indicator
	typingResponse := WebSocketResponse{
		Type:    "typing",
		Message: "AI is thinking...",
	}
	conn.WriteJSON(typingResponse)

	// Get response from service
	var response string
	var err error

	if msg.SessionID != "" {
		// Use chat with history
		response, err = h.service.ChatWithHistory(msg.SessionID, msg.Message)
	} else {
		// Use regular chat
		response, err = h.service.Chat(msg.Message)
	}

	if err != nil {
		errorResponse := WebSocketResponse{
			Type:  "error",
			Error: "Failed to get response: " + err.Error(),
		}
		conn.WriteJSON(errorResponse)
		return
	}

	// Send the response
	chatResponse := WebSocketResponse{
		Type:      "chat_response",
		Message:   response,
		SessionID: msg.SessionID,
	}
	conn.WriteJSON(chatResponse)
}

// handleTypingMessage handles typing indicators
func (h *WebSocketHandler) handleTypingMessage(conn *websocket.Conn, msg WebSocketMessage) {
	// Broadcast typing indicator to other clients (if needed)
	// For now, just acknowledge
	response := WebSocketResponse{
		Type: "typing_ack",
	}
	conn.WriteJSON(response)
}

// Broadcast sends a message to all connected clients
func (h *WebSocketHandler) Broadcast(msg WebSocketResponse) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for client := range h.clients {
		err := client.WriteJSON(msg)
		if err != nil {
			log.Printf("Failed to send message to client: %v", err)
			client.Close()
			delete(h.clients, client)
		}
	}
}

// GetConnectedClientsCount returns the number of connected clients
func (h *WebSocketHandler) GetConnectedClientsCount() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return len(h.clients)
}
