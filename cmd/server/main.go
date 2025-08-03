// cmd/server/main.go
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"tech-docs-ai/internal/app"
	"tech-docs-ai/internal/cache"
	"tech-docs-ai/internal/emb"
	"tech-docs-ai/internal/kafka"
	"tech-docs-ai/internal/repo"
	"tech-docs-ai/internal/vec"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error       string `json:"error"`
	Code        string `json:"code"`
	Description string `json:"description,omitempty"`
}

// Logger represents a structured logger
type Logger struct {
	*log.Logger
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp string      `json:"timestamp"`
	Level     string      `json:"level"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
}

// NewLogger creates a new structured logger
func NewLogger() *Logger {
	return &Logger{log.New(os.Stdout, "", 0)}
}

// Info logs an info message
func (l *Logger) Info(msg string, data interface{}) {
	l.log("INFO", msg, data)
}

// Error logs an error message
func (l *Logger) Error(msg string, err error, data interface{}) {
	l.log("ERROR", msg, map[string]interface{}{
		"error": err.Error(),
		"data":  data,
	})
}

// log writes a structured log entry
func (l *Logger) log(level, msg string, data interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level,
		Message:   msg,
		Data:      data,
	}

	jsonEntry, _ := json.Marshal(entry)
	l.Logger.Println(string(jsonEntry))
}

func main() {
	// Initialize structured logger
	logger := NewLogger()

	// Initialize a new router
	r := chi.NewRouter()

	// Add rate limiter configuration
	rateLimiter := app.NewRateLimiter(10, 30) // 10 requests per second, burst of 30

	// Update router middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(app.RateLimitMiddleware(rateLimiter))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:8080", "http://localhost"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Sec-WebSocket-Key", "Sec-WebSocket-Version", "Sec-WebSocket-Extensions", "Sec-WebSocket-Protocol"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Create a new Ollama embedding client
	ollamaClient := emb.NewOllamaClient()

	// Initialize Qdrant client
	qdrantClient := vec.NewQdrantClient()

	// Initialize PostgreSQL store
	postgresStore, err := repo.NewPostgresStore()
	if err != nil {
		logger.Error("Failed to initialize PostgreSQL store", err, nil)
		os.Exit(1)
	}
	defer postgresStore.Close()

	// Initialize Redis cache
	redisCache, err := cache.NewRedisCache()
	if err != nil {
		logger.Error("Failed to initialize Redis cache", err, nil)
		os.Exit(1)
	}
	defer redisCache.Close()
	
	// Add health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := redisCache.Health(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Redis connection failed"})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Initialize Kafka producer
	kafkaProducer := kafka.NewProducer()
	defer kafkaProducer.Close()

	// Create the main application service and handlers
	svc := app.NewService(ollamaClient, qdrantClient, postgresStore, kafkaProducer, redisCache)
	handler := app.NewHandler(svc)
	wsHandler := app.NewWebSocketHandler(svc)

	// Set up API routes
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

	// WebSocket endpoint for real-time chat
	r.HandleFunc("/ws", wsHandler.HandleWebSocket)

	// Add health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "healthy",
			"service": "tech-docs-ai",
			"version": os.Getenv("APP_VERSION"),
		})
	})

	// Run the seeder
	app.Seeder()

	// Start the server
	port := getEnv("PORT", "8080")
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       600 * time.Second,  // Increased for LLM responses
		WriteTimeout:      600 * time.Second,  // Increased for LLM responses
		IdleTimeout:       120 * time.Second,
	}

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start server in a goroutine
	go func() {
		logger.Info("Server starting", map[string]string{"port": port})
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start", err, nil)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()

	// Shutdown gracefully
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server shutdown failed", err, nil)
	}

	logger.Info("Server stopped gracefully", nil)
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	// Mask sensitive values in logs
	if strings.Contains(strings.ToLower(key), "password") ||
		strings.Contains(strings.ToLower(key), "secret") ||
		strings.Contains(strings.ToLower(key), "key") {
		return "[REDACTED]"
	}
	return value
}
