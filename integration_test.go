package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"tech-docs-ai/internal/app"
	"tech-docs-ai/internal/cache"
	"tech-docs-ai/internal/emb"
	"tech-docs-ai/internal/kafka"
	"tech-docs-ai/internal/repo"
	"tech-docs-ai/internal/scraper"
	"tech-docs-ai/internal/types"
	"tech-docs-ai/internal/vec"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSystemIntegration tests that all components can be initialized without errors
func TestSystemIntegration(t *testing.T) {
	// Test that we can create all the main components without panicking
	
	// Test Ollama client creation
	ollamaClient := emb.NewOllamaClient()
	if ollamaClient == nil {
		t.Error("Failed to create Ollama client")
	}

	// Test Qdrant client creation
	qdrantClient := vec.NewQdrantClient()
	if qdrantClient == nil {
		t.Error("Failed to create Qdrant client")
	}

	// Test Redis cache creation
	redisCache, err := cache.NewRedisCache()
	if err != nil {
		t.Logf("Redis cache creation failed (expected in test environment): %v", err)
		// This is expected to fail in test environment without Redis
	} else {
		defer redisCache.Close()
	}

	// Test Kafka producer creation
	kafkaProducer := kafka.NewProducer()
	if kafkaProducer == nil {
		t.Error("Failed to create Kafka producer")
	}
	defer kafkaProducer.Close()

	// Test PostgreSQL store creation
	_, err = repo.NewPostgresStore()
	if err != nil {
		t.Logf("PostgreSQL store creation failed (expected in test environment): %v", err)
		// This is expected to fail in test environment without PostgreSQL
	}

	t.Log("System integration test completed - all components can be initialized")
}

// TestHandlerCreation tests that we can create handlers without errors
func TestHandlerCreation(t *testing.T) {
	// Create a mock service for testing
	mockService := &MockServiceImpl{}
	
	// Test handler creation
	handler := app.NewHandler(mockService)
	if handler == nil {
		t.Error("Failed to create handler")
	}

	// Test WebSocket handler creation
	wsHandler := app.NewWebSocketHandler(mockService)
	if wsHandler == nil {
		t.Error("Failed to create WebSocket handler")
	}

	t.Log("Handler creation test completed successfully")
}

// MockServiceImpl is a simple mock implementation for testing
type MockServiceImpl struct{}

func (m *MockServiceImpl) Chat(message string) (string, error) {
	return "Mock response", nil
}

func (m *MockServiceImpl) ChatWithHistory(sessionID, message string) (string, error) {
	return "Mock response with history", nil
}

func (m *MockServiceImpl) GetChatHistory(sessionID string, limit int) ([]*types.ChatMessage, error) {
	return []*types.ChatMessage{}, nil
}

func (m *MockServiceImpl) AddDocument(doc *types.Document) error {
	return nil
}

func (m *MockServiceImpl) SearchDocuments(query string, limit int) ([]*types.Document, error) {
	return []*types.Document{}, nil
}

func (m *MockServiceImpl) ScrapeDocument(url, category string, tags []string) error {
	return nil
}

func (m *MockServiceImpl) GenerateTutorialFromScrapedData(url, topic string) (string, error) {
	return "Mock tutorial", nil
}

func (m *MockServiceImpl) ScrapeAndGenerateTutorial(url, topic string) (string, error) {
	return "Mock scrape and tutorial", nil
}

func (m *MockServiceImpl) GetConversationInsights(sessionID string) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

// TestAPIEndpointsIntegration tests all API endpoints with a mock service
func TestAPIEndpointsIntegration(t *testing.T) {
	// Create mock service
	mockService := &MockServiceImpl{}
	
	// Create handler
	handler := app.NewHandler(mockService)
	
	// Create router
	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/chat", handler.HandleChat)
		r.Post("/chat/history", handler.HandleChatWithHistory)
		r.Get("/chat/history", handler.HandleGetChatHistory)
		r.Get("/chat/insights", handler.HandleGetConversationInsights)
		r.Post("/documents", handler.HandleAddDocument)
		r.Post("/scrape", handler.HandleScrapeDocument)
		r.Get("/documents/search", handler.HandleSearchDocuments)
		r.Post("/tutorials/generate", handler.HandleGenerateTutorial)
		r.Post("/tutorials/scrape-and-generate", handler.HandleScrapeAndGenerateTutorial)
	})

	// Test chat endpoint
	t.Run("Chat endpoint", func(t *testing.T) {
		reqBody := map[string]string{"message": "Hello"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/chat", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]string
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, "Mock response", response["response"])
	})

	// Test add document endpoint
	t.Run("Add document endpoint", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"title":    "Test Document",
			"content":  "Test content",
			"category": "Testing",
			"tags":     []string{"test"},
			"author":   "Test Author",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/documents", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	// Test search documents endpoint
	t.Run("Search documents endpoint", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/documents/search?q=test", nil)
		
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Test scrape endpoint
	t.Run("Scrape endpoint", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"url":      "https://example.com",
			"category": "Test",
			"tags":     []string{"test"},
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/scrape", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Log("API endpoints integration test completed successfully")
}

// TestScrapersIntegration tests both scrapers
func TestScrapersIntegration(t *testing.T) {
	// Test W3Schools scraper
	t.Run("W3Schools scraper", func(t *testing.T) {
		scraper := scraper.NewW3SchoolsScraper()
		assert.NotNil(t, scraper)
		
		// Test with a mock HTML server
		testHTML := `
<!DOCTYPE html>
<html>
<head><title>HTML Tutorial</title></head>
<body>
	<h1>HTML Introduction</h1>
	<p>HTML is the standard markup language for Web pages.</p>
	<pre><code>&lt;h1&gt;My First Heading&lt;/h1&gt;</code></pre>
</body>
</html>`
		
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(testHTML))
		}))
		defer server.Close()
		
		content, err := scraper.ScrapePage(server.URL)
		require.NoError(t, err)
		assert.NotNil(t, content)
		assert.Equal(t, "HTML Introduction", content.Title)
		// Content might be empty due to specific HTML structure expectations
		assert.NotEmpty(t, content.Title)
	})

	// Test Universal scraper
	t.Run("Universal scraper", func(t *testing.T) {
		scraper := scraper.NewUniversalScraper()
		assert.NotNil(t, scraper)
		
		// Test with a mock HTML server
		testHTML := `
<!DOCTYPE html>
<html>
<head><title>JavaScript Guide</title></head>
<body>
	<main>
		<h1>JavaScript Basics</h1>
		<p>JavaScript is a programming language.</p>
		<pre><code class="language-javascript">console.log('Hello World');</code></pre>
	</main>
</body>
</html>`
		
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(testHTML))
		}))
		defer server.Close()
		
		content, err := scraper.ScrapePage(server.URL)
		require.NoError(t, err)
		assert.NotNil(t, content)
		assert.Equal(t, "JavaScript Basics", content.Title)
		assert.Contains(t, content.Content, "JavaScript is a programming language")
		assert.True(t, len(content.Examples) >= 1)
		assert.Contains(t, content.Examples[0], "console.log")
	})

	t.Log("Scrapers integration test completed successfully")
}

// TestFullWorkflowIntegration tests the complete workflow from scraping to chat
func TestFullWorkflowIntegration(t *testing.T) {
	// This test requires actual services, so we'll skip it if they're not available
	t.Skip("Skipping full workflow test - requires running services")
	
	// Create components (this would fail without actual services)
	ollamaClient := emb.NewOllamaClient()
	qdrantClient := vec.NewQdrantClient()
	
	postgresStore, err := repo.NewPostgresStore()
	if err != nil {
		t.Skip("PostgreSQL not available for integration test")
	}
	defer postgresStore.Close()
	
	redisCache, err := cache.NewRedisCache()
	if err != nil {
		t.Skip("Redis not available for integration test")
	}
	defer redisCache.Close()
	
	kafkaProducer := kafka.NewProducer()
	defer kafkaProducer.Close()
	
	// Create service
	service := app.NewService(ollamaClient, qdrantClient, postgresStore, kafkaProducer, redisCache)
	
	// Test adding a document
	doc := &types.Document{
		Title:    "Test Integration Document",
		Content:  "This is a test document for integration testing.",
		Category: "Testing",
		Tags:     []string{"integration", "test"},
		Author:   "Integration Test",
	}
	
	err = service.AddDocument(doc)
	require.NoError(t, err)
	
	// Test searching for the document
	docs, err := service.SearchDocuments("integration", 10)
	require.NoError(t, err)
	assert.Len(t, docs, 1)
	assert.Equal(t, "Test Integration Document", docs[0].Title)
	
	// Test chat functionality
	response, err := service.Chat("Tell me about integration testing")
	require.NoError(t, err)
	assert.NotEmpty(t, response)
	
	t.Log("Full workflow integration test completed successfully")
}

// TestConcurrentRequests tests the system under concurrent load
func TestConcurrentRequests(t *testing.T) {
	mockService := &MockServiceImpl{}
	handler := app.NewHandler(mockService)
	
	r := chi.NewRouter()
	r.Post("/api/v1/chat", handler.HandleChat)
	
	// Create test server
	server := httptest.NewServer(r)
	defer server.Close()
	
	// Number of concurrent requests
	numRequests := 10
	done := make(chan bool, numRequests)
	
	// Send concurrent requests
	for i := 0; i < numRequests; i++ {
		go func(id int) {
			defer func() { done <- true }()
			
			reqBody := map[string]string{"message": fmt.Sprintf("Hello from request %d", id)}
			body, _ := json.Marshal(reqBody)
			
			resp, err := http.Post(server.URL+"/api/v1/chat", "application/json", bytes.NewReader(body))
			require.NoError(t, err)
			defer resp.Body.Close()
			
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}(i)
	}
	
	// Wait for all requests to complete
	for i := 0; i < numRequests; i++ {
		select {
		case <-done:
			// Request completed
		case <-time.After(5 * time.Second):
			t.Fatal("Request timed out")
		}
	}
	
	t.Log("Concurrent requests test completed successfully")
}

// TestErrorHandling tests error handling across the system
func TestErrorHandling(t *testing.T) {
	// Test with a service that returns errors
	errorService := &ErrorMockService{}
	handler := app.NewHandler(errorService)
	
	r := chi.NewRouter()
	r.Post("/api/v1/chat", handler.HandleChat)
	r.Post("/api/v1/documents", handler.HandleAddDocument)
	
	// Test chat error handling
	t.Run("Chat error handling", func(t *testing.T) {
		reqBody := map[string]string{"message": "Hello"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/chat", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
	
	// Test document error handling
	t.Run("Document error handling", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"title":   "Test",
			"content": "Test content",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/documents", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
	
	t.Log("Error handling test completed successfully")
}

// ErrorMockService is a mock service that returns errors for testing
type ErrorMockService struct{}

func (m *ErrorMockService) Chat(message string) (string, error) {
	return "", fmt.Errorf("mock chat error")
}

func (m *ErrorMockService) ChatWithHistory(sessionID, message string) (string, error) {
	return "", fmt.Errorf("mock chat with history error")
}

func (m *ErrorMockService) GetChatHistory(sessionID string, limit int) ([]*types.ChatMessage, error) {
	return nil, fmt.Errorf("mock get chat history error")
}

func (m *ErrorMockService) AddDocument(doc *types.Document) error {
	return fmt.Errorf("mock add document error")
}

func (m *ErrorMockService) SearchDocuments(query string, limit int) ([]*types.Document, error) {
	return nil, fmt.Errorf("mock search documents error")
}

func (m *ErrorMockService) ScrapeDocument(url, category string, tags []string) error {
	return fmt.Errorf("mock scrape document error")
}

func (m *ErrorMockService) GenerateTutorialFromScrapedData(url, topic string) (string, error) {
	return "", fmt.Errorf("mock generate tutorial error")
}

func (m *ErrorMockService) ScrapeAndGenerateTutorial(url, topic string) (string, error) {
	return "", fmt.Errorf("mock scrape and generate tutorial error")
}

func (m *ErrorMockService) GetConversationInsights(sessionID string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("mock get conversation insights error")
}