package app

import (
	"log"
	"net/http"
	"sync"

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

	log.Printf("WebSocket connection established from %s", r.RemoteAddr)

	// Create a mutex to protect WebSocket writes
	var writeMutex sync.Mutex

	// Send a welcome message to confirm connection
	welcomeMsg := WebSocketMessage{
		Type:     "connection",
		Response: "WebSocket connection established successfully",
	}
	writeMutex.Lock()
	if err := conn.WriteJSON(welcomeMsg); err != nil {
		writeMutex.Unlock()
		log.Printf("Failed to send welcome message: %v", err)
		return
	}
	writeMutex.Unlock()

	for {
		var msg WebSocketMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("WebSocket read error from %s: %v", r.RemoteAddr, err)
			break
		}

		log.Printf("Received WebSocket message: type=%s, message=%s", msg.Type, msg.Message)

		switch msg.Type {
		case "chat":
			// Send typing indicator
			h.sendTypingIndicatorSafe(conn, &writeMutex)
			
			// Process chat in a goroutine to avoid blocking the WebSocket connection
			go func(connection *websocket.Conn, message string, mutex *sync.Mutex) {
				log.Printf("Starting chat processing for message: %s", message)
				response, err := h.service.Chat(message)
				if err != nil {
					log.Printf("Chat processing failed: %v", err)
					h.sendErrorSafe(connection, mutex, "Failed to process chat message", err)
					return
				}
				log.Printf("Chat processing completed, sending response")
				h.sendResponseSafe(connection, mutex, "chat_response", response)
				log.Printf("Response sent successfully")
			}(conn, msg.Message, &writeMutex)

		case "chat_with_history":
			if msg.SessionID == "" {
				h.sendErrorSafe(conn, &writeMutex, "Session ID is required for chat with history", nil)
				continue
			}
			
			// Send typing indicator
			h.sendTypingIndicatorSafe(conn, &writeMutex)
			
			// Process chat with history in a goroutine to avoid blocking
			go func(connection *websocket.Conn, sessionID, message string, mutex *sync.Mutex) {
				response, err := h.service.ChatWithHistory(sessionID, message)
				if err != nil {
					h.sendErrorSafe(connection, mutex, "Failed to process chat with history", err)
					return
				}
				h.sendResponseSafe(connection, mutex, "chat_response", response)
			}(conn, msg.SessionID, msg.Message, &writeMutex)

		default:
			h.sendErrorSafe(conn, &writeMutex, "Unknown message type", nil)
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

// Thread-safe versions of the helper functions

// sendResponseSafe sends a successful response over WebSocket with mutex protection
func (h *WebSocketHandler) sendResponseSafe(conn *websocket.Conn, mutex *sync.Mutex, msgType, response string) {
	msg := WebSocketMessage{
		Type:     msgType,
		Response: response,
	}
	mutex.Lock()
	defer mutex.Unlock()
	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("WebSocket write error: %v", err)
	}
}

// sendErrorSafe sends an error response over WebSocket with mutex protection
func (h *WebSocketHandler) sendErrorSafe(conn *websocket.Conn, mutex *sync.Mutex, message string, err error) {
	errorMsg := message
	if err != nil {
		errorMsg += ": " + err.Error()
	}

	msg := WebSocketMessage{
		Type:  "error",
		Error: errorMsg,
	}
	mutex.Lock()
	defer mutex.Unlock()
	if writeErr := conn.WriteJSON(msg); writeErr != nil {
		log.Printf("WebSocket write error: %v", writeErr)
	}
}

// sendTypingIndicatorSafe sends a typing indicator over WebSocket with mutex protection
func (h *WebSocketHandler) sendTypingIndicatorSafe(conn *websocket.Conn, mutex *sync.Mutex) {
	msg := WebSocketMessage{
		Type: "typing",
	}
	mutex.Lock()
	defer mutex.Unlock()
	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("WebSocket write error: %v", err)
	}
}