package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"tech-docs-ai/internal/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockServiceForTesting is a proper mock implementation for testing
type MockServiceForTesting struct {
	mock.Mock
}

func (m *MockServiceForTesting) Chat(message string) (string, error) {
	args := m.Called(message)
	return args.String(0), args.Error(1)
}

func (m *MockServiceForTesting) ChatWithHistory(sessionID, message string) (string, error) {
	args := m.Called(sessionID, message)
	return args.String(0), args.Error(1)
}

func (m *MockServiceForTesting) GetChatHistory(sessionID string, limit int) ([]*types.ChatMessage, error) {
	args := m.Called(sessionID, limit)
	return args.Get(0).([]*types.ChatMessage), args.Error(1)
}

func (m *MockServiceForTesting) AddDocument(doc *types.Document) error {
	args := m.Called(doc)
	return args.Error(0)
}

func (m *MockServiceForTesting) SearchDocuments(query string, limit int) ([]*types.Document, error) {
	args := m.Called(query, limit)
	return args.Get(0).([]*types.Document), args.Error(1)
}

func (m *MockServiceForTesting) ScrapeDocument(url, category string, tags []string) error {
	args := m.Called(url, category, tags)
	return args.Error(0)
}

func (m *MockServiceForTesting) GenerateTutorialFromScrapedData(url, topic string) (string, error) {
	args := m.Called(url, topic)
	return args.String(0), args.Error(1)
}

func (m *MockServiceForTesting) ScrapeAndGenerateTutorial(url, topic string) (string, error) {
	args := m.Called(url, topic)
	return args.String(0), args.Error(1)
}

func (m *MockServiceForTesting) GetConversationInsights(sessionID string) (map[string]interface{}, error) {
	args := m.Called(sessionID)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func TestHandler_HandleChat_Success(t *testing.T) {
	mockService := new(MockServiceForTesting)
	mockService.On("Chat", "Hello").Return("Hi there!", nil)

	handler := NewHandler(mockService)

	reqBody := chatRequest{Message: "Hello"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/chat", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.HandleChat(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response chatResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "Hi there!", response.Response)

	mockService.AssertExpectations(t)
}

func TestHandler_HandleChat_EmptyMessage(t *testing.T) {
	mockService := new(MockServiceForTesting)
	handler := NewHandler(mockService)

	reqBody := chatRequest{Message: ""}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/chat", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.HandleChat(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errorResp ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&errorResp)
	require.NoError(t, err)
	assert.Equal(t, ErrValidation, errorResp.Code)
	assert.Contains(t, errorResp.Error, "message cannot be empty")
}

func TestHandler_HandleChat_MessageTooLong(t *testing.T) {
	mockService := new(MockServiceForTesting)
	handler := NewHandler(mockService)

	longMessage := string(make([]byte, 2001)) // Exceeds 2000 character limit
	reqBody := chatRequest{Message: longMessage}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/chat", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.HandleChat(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errorResp ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&errorResp)
	require.NoError(t, err)
	assert.Equal(t, ErrValidation, errorResp.Code)
	assert.Contains(t, errorResp.Error, "message too long")
}

func TestHandler_HandleChat_ServiceError(t *testing.T) {
	mockService := new(MockServiceForTesting)
	mockService.On("Chat", "Hello").Return("", fmt.Errorf("service error"))

	handler := NewHandler(mockService)

	reqBody := chatRequest{Message: "Hello"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/chat", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.HandleChat(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var errorResp ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&errorResp)
	require.NoError(t, err)
	assert.Equal(t, ErrInternalServer, errorResp.Code)

	mockService.AssertExpectations(t)
}

func TestHandler_HandleAddDocument_Success(t *testing.T) {
	mockService := new(MockServiceForTesting)
	mockService.On("AddDocument", mock.AnythingOfType("*types.Document")).Return(nil)

	handler := NewHandler(mockService)

	reqBody := documentRequest{
		Title:    "Test Document",
		Content:  "Test content",
		Category: "Testing",
		Tags:     []string{"test"},
		Author:   "Test Author",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/documents", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.HandleAddDocument(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]string
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "Document added successfully", response["status"])

	mockService.AssertExpectations(t)
}

func TestHandler_HandleAddDocument_EmptyTitle(t *testing.T) {
	mockService := new(MockServiceForTesting)
	handler := NewHandler(mockService)

	reqBody := documentRequest{
		Title:   "",
		Content: "Test content",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/documents", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.HandleAddDocument(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errorResp ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&errorResp)
	require.NoError(t, err)
	assert.Equal(t, ErrValidation, errorResp.Code)
	assert.Contains(t, errorResp.Error, "title cannot be empty")
}

func TestHandler_HandleSearchDocuments_Success(t *testing.T) {
	mockService := new(MockServiceForTesting)
	expectedDocs := []*types.Document{
		{ID: "1", Title: "Test Doc", Content: "Test content"},
	}
	mockService.On("SearchDocuments", "test", 10).Return(expectedDocs, nil)

	handler := NewHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/documents/search?q=test", nil)
	w := httptest.NewRecorder()
	handler.HandleSearchDocuments(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var docs []*types.Document
	err := json.NewDecoder(w.Body).Decode(&docs)
	require.NoError(t, err)
	assert.Equal(t, expectedDocs, docs)

	mockService.AssertExpectations(t)
}

func TestHandler_HandleSearchDocuments_EmptyQuery(t *testing.T) {
	mockService := new(MockServiceForTesting)
	handler := NewHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/documents/search?q=", nil)
	w := httptest.NewRecorder()
	handler.HandleSearchDocuments(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errorResp ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&errorResp)
	require.NoError(t, err)
	assert.Equal(t, ErrValidation, errorResp.Code)
	assert.Contains(t, errorResp.Error, "Query parameter 'q' is required")
}

func TestHandler_HandleSearchDocuments_InvalidLimit(t *testing.T) {
	mockService := new(MockServiceForTesting)
	handler := NewHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/documents/search?q=test&limit=invalid", nil)
	w := httptest.NewRecorder()
	handler.HandleSearchDocuments(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errorResp ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&errorResp)
	require.NoError(t, err)
	assert.Equal(t, ErrValidation, errorResp.Code)
	assert.Contains(t, errorResp.Error, "Invalid limit parameter")
}

func TestHandler_HandleScrapeDocument_Success(t *testing.T) {
	mockService := new(MockServiceForTesting)
	mockService.On("ScrapeDocument", "https://example.com", "Test", []string{"test"}).Return(nil)

	handler := NewHandler(mockService)

	reqBody := scrapeRequest{
		URL:      "https://example.com",
		Category: "Test",
		Tags:     []string{"test"},
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/scrape", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.HandleScrapeDocument(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "Scraping job queued successfully", response["status"])

	mockService.AssertExpectations(t)
}

func TestHandler_HandleChatWithHistory_Success(t *testing.T) {
	mockService := new(MockServiceForTesting)
	mockService.On("ChatWithHistory", "session123", "Hello").Return("Hi there!", nil)

	handler := NewHandler(mockService)

	reqBody := chatWithHistoryRequest{
		SessionID: "session123",
		Message:   "Hello",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/chat/history", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.HandleChatWithHistory(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response chatResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "Hi there!", response.Response)

	mockService.AssertExpectations(t)
}

func TestHandler_HandleGetChatHistory_Success(t *testing.T) {
	mockService := new(MockServiceForTesting)
	expectedHistory := []*types.ChatMessage{
		{ID: "1", Role: "user", Content: "Hello"},
		{ID: "2", Role: "assistant", Content: "Hi there!"},
	}
	mockService.On("GetChatHistory", "session123", 20).Return(expectedHistory, nil)

	handler := NewHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/chat/history?session_id=session123", nil)
	w := httptest.NewRecorder()
	handler.HandleGetChatHistory(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "session123", response["session_id"])
	assert.Equal(t, float64(2), response["count"]) // JSON numbers are float64

	mockService.AssertExpectations(t)
}

func TestHandler_HandleGetConversationInsights_Success(t *testing.T) {
	mockService := new(MockServiceForTesting)
	expectedInsights := map[string]interface{}{
		"total_messages": 5,
		"topics":         []string{"javascript", "react"},
	}
	mockService.On("GetConversationInsights", "session123").Return(expectedInsights, nil)

	handler := NewHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/chat/insights?session_id=session123", nil)
	w := httptest.NewRecorder()
	handler.HandleGetConversationInsights(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "session123", response["session_id"])
	assert.NotNil(t, response["insights"])

	mockService.AssertExpectations(t)
}

func TestValidateRequest(t *testing.T) {
	// Test valid chat request
	t.Run("Valid chat request", func(t *testing.T) {
		reqBody := chatRequest{Message: "Hello"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		var target chatRequest
		err := validateRequest(req, &target)
		assert.NoError(t, err)
		assert.Equal(t, "Hello", target.Message)
	})

	// Test empty chat message
	t.Run("Empty chat message", func(t *testing.T) {
		reqBody := chatRequest{Message: ""}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		var target chatRequest
		err := validateRequest(req, &target)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "message cannot be empty")
	})

	// Test long chat message
	t.Run("Long chat message", func(t *testing.T) {
		reqBody := chatRequest{Message: string(make([]byte, 2001))}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		var target chatRequest
		err := validateRequest(req, &target)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "message too long")
	})

	// Test valid document request
	t.Run("Valid document request", func(t *testing.T) {
		reqBody := documentRequest{Title: "Test", Content: "Content"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		var target documentRequest
		err := validateRequest(req, &target)
		assert.NoError(t, err)
		assert.Equal(t, "Test", target.Title)
	})

	// Test empty document title
	t.Run("Empty document title", func(t *testing.T) {
		reqBody := documentRequest{Title: "", Content: "Content"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		var target documentRequest
		err := validateRequest(req, &target)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "title cannot be empty")
	})

	// Test valid scrape request
	t.Run("Valid scrape request", func(t *testing.T) {
		reqBody := scrapeRequest{URL: "https://example.com"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		var target scrapeRequest
		err := validateRequest(req, &target)
		assert.NoError(t, err)
		assert.Equal(t, "https://example.com", target.URL)
	})

	// Test invalid scrape URL
	t.Run("Invalid scrape URL", func(t *testing.T) {
		reqBody := scrapeRequest{URL: "ht tp://invalid url with spaces"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		var target scrapeRequest
		err := validateRequest(req, &target)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid URL format")
	})
}