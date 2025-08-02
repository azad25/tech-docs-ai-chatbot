// internal/app/handler.go
package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"tech-docs-ai/internal/types"
)

// ServiceInterface defines the interface that Service implements
type ServiceInterface interface {
	Chat(message string) (string, error)
	ChatWithHistory(sessionID, message string) (string, error)
	GetChatHistory(sessionID string, limit int) ([]*types.ChatMessage, error)
	AddDocument(doc *types.Document) error
	SearchDocuments(query string, limit int) ([]*types.Document, error)
	ScrapeDocument(url, category string, tags []string) error
	GenerateTutorialFromScrapedData(url, topic string) (string, error)
	ScrapeAndGenerateTutorial(url, topic string) (string, error)
	GetConversationInsights(sessionID string) (map[string]interface{}, error)
}

// Handler handles HTTP requests for the application.
type Handler struct {
	service ServiceInterface
}

// NewHandler creates a new Handler with the given service.
func NewHandler(svc ServiceInterface) *Handler {
	return &Handler{service: svc}
}

// chatRequest defines the structure for an incoming chat request.
type chatRequest struct {
	Message string `json:"message"`
}

// chatResponse defines the structure for a chat response.
type chatResponse struct {
	Response string `json:"response"`
}

// documentRequest defines the structure for adding a document.
type documentRequest struct {
	Title    string            `json:"title"`
	Content  string            `json:"content"`
	Category string            `json:"category"`
	Tags     []string          `json:"tags"`
	Author   string            `json:"author"`
	Metadata map[string]string `json:"metadata"`
}

// scrapeRequest defines the structure for triggering a scrape.
type scrapeRequest struct {
	URL      string   `json:"url"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
}

// tutorialRequest defines the structure for generating a tutorial.
type tutorialRequest struct {
	URL   string `json:"url"`
	Topic string `json:"topic"`
}

// scrapeTutorialRequest defines the structure for scraping and generating a tutorial.
type scrapeTutorialRequest struct {
	URL   string `json:"url"`
	Topic string `json:"topic"`
}

// chatWithHistoryRequest defines the structure for chat with history.
type chatWithHistoryRequest struct {
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
}

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

// Error codes for API responses
const (
	ErrInvalidRequest    = "INVALID_REQUEST"
	ErrInternalServer    = "INTERNAL_SERVER_ERROR"
	ErrValidation        = "VALIDATION_ERROR"
	ErrResourceNotFound  = "RESOURCE_NOT_FOUND"
	ErrUnauthorized      = "UNAUTHORIZED"
)

// sendError sends a standardized error response
func sendError(w http.ResponseWriter, code int, errCode, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: message,
		Code:  errCode,
	})
}

// validateRequest validates the request body against a schema
func validateRequest(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return fmt.Errorf("invalid request body: %w", err)
	}

	// Add validation for specific request types
	switch req := v.(type) {
	case *chatRequest:
		if strings.TrimSpace(req.Message) == "" {
			return fmt.Errorf("message cannot be empty")
		}
		if len(req.Message) > 2000 {
			return fmt.Errorf("message too long (max 2000 characters)")
		}
	case *documentRequest:
		if strings.TrimSpace(req.Title) == "" {
			return fmt.Errorf("title cannot be empty")
		}
		if strings.TrimSpace(req.Content) == "" {
			return fmt.Errorf("content cannot be empty")
		}
	case *scrapeRequest:
		if strings.TrimSpace(req.URL) == "" {
			return fmt.Errorf("URL cannot be empty")
		}
		if _, err := url.Parse(req.URL); err != nil {
			return fmt.Errorf("invalid URL format")
		}
	}

	return nil
}

// HandleChat handles requests to the chat endpoint with improved error handling
func (h *Handler) HandleChat(w http.ResponseWriter, r *http.Request) {
	var req chatRequest
	if err := validateRequest(r, &req); err != nil {
		sendError(w, http.StatusBadRequest, ErrValidation, err.Error())
		return
	}

	response, err := h.service.Chat(req.Message)
	if err != nil {
		sendError(w, http.StatusInternalServerError, ErrInternalServer, "Failed to process chat request")
		log.Printf("Chat error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chatResponse{Response: response})
}

// HandleAddDocument handles requests to add a new document with improved validation
func (h *Handler) HandleAddDocument(w http.ResponseWriter, r *http.Request) {
	var req documentRequest
	if err := validateRequest(r, &req); err != nil {
		sendError(w, http.StatusBadRequest, ErrValidation, err.Error())
		return
	}

	doc := &types.Document{
		Title:    req.Title,
		Content:  req.Content,
		Category: req.Category,
		Tags:     req.Tags,
		Author:   req.Author,
		Metadata: req.Metadata,
	}

	if err := h.service.AddDocument(doc); err != nil {
		sendError(w, http.StatusInternalServerError, ErrInternalServer, "Failed to add document")
		log.Printf("Add document error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "Document added successfully"})
}

// HandleScrapeDocument handles requests to scrape a document from a URL.
func (h *Handler) HandleScrapeDocument(w http.ResponseWriter, r *http.Request) {
	var req scrapeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.ScrapeDocument(req.URL, req.Category, req.Tags); err != nil {
		http.Error(w, "Failed to scrape document", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "Scraping job queued successfully"})
}

// HandleSearchDocuments handles requests to search documents with improved validation
func (h *Handler) HandleSearchDocuments(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if strings.TrimSpace(query) == "" {
		sendError(w, http.StatusBadRequest, ErrValidation, "Query parameter 'q' is required")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil || l < 1 || l > 100 {
			sendError(w, http.StatusBadRequest, ErrValidation, "Invalid limit parameter (must be between 1 and 100)")
			return
		}
		limit = l
	}

	documents, err := h.service.SearchDocuments(query, limit)
	if err != nil {
		sendError(w, http.StatusInternalServerError, ErrInternalServer, "Failed to search documents")
		log.Printf("Search documents error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(documents)
}

// HandleGenerateTutorial handles requests to generate a tutorial from scraped data.
func (h *Handler) HandleGenerateTutorial(w http.ResponseWriter, r *http.Request) {
	var req tutorialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Topic == "" {
		http.Error(w, "Topic is required", http.StatusBadRequest)
		return
	}

	tutorial, err := h.service.GenerateTutorialFromScrapedData(req.URL, req.Topic)
	if err != nil {
		http.Error(w, "Failed to generate tutorial", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"tutorial": tutorial,
		"topic":    req.Topic,
	})
}

// HandleScrapeAndGenerateTutorial handles requests to scrape content and generate a tutorial.
func (h *Handler) HandleScrapeAndGenerateTutorial(w http.ResponseWriter, r *http.Request) {
	var req scrapeTutorialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.URL == "" || req.Topic == "" {
		http.Error(w, "URL and topic are required", http.StatusBadRequest)
		return
	}

	tutorial, err := h.service.ScrapeAndGenerateTutorial(req.URL, req.Topic)
	if err != nil {
		http.Error(w, "Failed to scrape and generate tutorial", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"tutorial": tutorial,
		"topic":    req.Topic,
		"url":      req.URL,
	})
}

// HandleChatWithHistory handles requests to chat with conversation history.
func (h *Handler) HandleChatWithHistory(w http.ResponseWriter, r *http.Request) {
	var req chatWithHistoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.SessionID == "" || req.Message == "" {
		http.Error(w, "Session ID and message are required", http.StatusBadRequest)
		return
	}

	response, err := h.service.ChatWithHistory(req.SessionID, req.Message)
	if err != nil {
		http.Error(w, "Failed to get chat response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chatResponse{Response: response})
}

// HandleGetChatHistory handles requests to get chat history.
func (h *Handler) HandleGetChatHistory(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	history, err := h.service.GetChatHistory(sessionID, limit)
	if err != nil {
		http.Error(w, "Failed to get chat history", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"session_id": sessionID,
		"history":    history,
		"count":      len(history),
	})
}

// HandleGetConversationInsights handles requests to get conversation insights.
func (h *Handler) HandleGetConversationInsights(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	insights, err := h.service.GetConversationInsights(sessionID)
	if err != nil {
		http.Error(w, "Failed to get conversation insights", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"session_id": sessionID,
		"insights":   insights,
	})
}
