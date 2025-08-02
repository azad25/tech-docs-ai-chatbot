package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"tech-docs-ai/internal/emb"
	"tech-docs-ai/internal/types"
	"tech-docs-ai/internal/vec"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockService is a mock implementation of the Service interface
type MockService struct {
	mock.Mock
}

func (m *MockService) Chat(message string) (string, error) {
	args := m.Called(message)
	return args.String(0), args.Error(1)
}

func (m *MockService) ChatWithHistory(sessionID, message string) (string, error) {
	args := m.Called(sessionID, message)
	return args.String(0), args.Error(1)
}

func (m *MockService) GetChatHistory(sessionID string, limit int) ([]*types.ChatMessage, error) {
	args := m.Called(sessionID, limit)
	return args.Get(0).([]*types.ChatMessage), args.Error(1)
}

func (m *MockService) AddDocument(doc *types.Document) error {
	args := m.Called(doc)
	return args.Error(0)
}

func (m *MockService) SearchDocuments(query string, limit int) ([]*types.Document, error) {
	args := m.Called(query, limit)
	return args.Get(0).([]*types.Document), args.Error(1)
}

func TestHandleChat(t *testing.T) {
	tests := []struct {
		name           string
		inputMessage   string
		expectedCode   int
		serviceError   error
		serviceResponse string
	}{
		{
			name:           "successful chat",
			inputMessage:   "Hello",
			expectedCode:   http.StatusOK,
			serviceError:   nil,
			serviceResponse: "Hello, how can I help you?",
		},
		{
			name:           "empty message",
			inputMessage:   "",
			expectedCode:   http.StatusBadRequest,
			serviceError:   nil,
			serviceResponse: "",
		},
		{
			name:           "service error",
			inputMessage:   "Hello",
			expectedCode:   http.StatusInternalServerError,
			serviceError:   fmt.Errorf("service error"),
			serviceResponse: "",
		},
	}

	for _, tt := range tests {
		t := tt
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockService)
			if tt.inputMessage != "" {
				mockService.On("Chat", tt.inputMessage).Return(tt.serviceResponse, tt.serviceError)
			}

			handler := NewHandler(mockService)

			reqBody := chatRequest{Message: tt.inputMessage}
			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/chat", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			handler.HandleChat(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode == http.StatusOK {
				var response chatResponse
				err := json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.serviceResponse, response.Response)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestHandleSearchDocuments(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		limit          string
		expectedCode   int
		expectedDocs   []*types.Document
		serviceError   error
	}{
		{
			name:         "successful search",
			query:        "test",
			limit:        "10",
			expectedCode: http.StatusOK,
			expectedDocs: []*types.Document{
				{ID: "1", Title: "Test Doc", Content: "Test content"},
			},
			serviceError: nil,
		},
		{
			name:         "empty query",
			query:        "",
			limit:        "10",
			expectedCode: http.StatusBadRequest,
			expectedDocs: nil,
			serviceError: nil,
		},
		{
			name:         "invalid limit",
			query:        "test",
			limit:        "invalid",
			expectedCode: http.StatusBadRequest,
			expectedDocs: nil,
			serviceError: nil,
		},
	}

	for _, tt := range tests {
		t := tt
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockService)
			if tt.query != "" && tt.expectedCode != http.StatusBadRequest {
				limit := 10
				mockService.On("SearchDocuments", tt.query, limit).Return(tt.expectedDocs, tt.serviceError)
			}

			handler := NewHandler(mockService)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/documents/search?q="+tt.query+"&limit="+tt.limit, nil)
			w := httptest.NewRecorder()
			handler.HandleSearchDocuments(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode == http.StatusOK {
				var docs []*types.Document
				err := json.NewDecoder(w.Body).Decode(&docs)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedDocs, docs)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestHandleAddDocument(t *testing.T) {
	tests := []struct {
		name         string
		doc          *types.Document
		expectedCode int
		serviceError error
	}{
		{
			name: "successful add",
			doc: &types.Document{
				Title:   "Test Doc",
				Content: "Test content",
			},
			expectedCode: http.StatusCreated,
			serviceError: nil,
		},
		{
			name: "empty title",
			doc: &types.Document{
				Title:   "",
				Content: "Test content",
			},
			expectedCode: http.StatusBadRequest,
			serviceError: nil,
		},
		{
			name: "service error",
			doc: &types.Document{
				Title:   "Test Doc",
				Content: "Test content",
			},
			expectedCode: http.StatusInternalServerError,
			serviceError: fmt.Errorf("service error"),
		},
	}

	for _, tt := range tests {
		t := tt
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockService)
			if tt.doc.Title != "" {
				mockService.On("AddDocument", tt.doc).Return(tt.serviceError)
			}

			handler := NewHandler(mockService)

			reqBody := documentRequest{
				Title:   tt.doc.Title,
				Content: tt.doc.Content,
			}
			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/documents", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			handler.HandleAddDocument(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}
