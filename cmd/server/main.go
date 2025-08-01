// cmd/server/main.go
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"tech-docs-ai/internal/app"
	"tech-docs-ai/internal/cache"
	"tech-docs-ai/internal/emb"
	"tech-docs-ai/internal/kafka"
	"tech-docs-ai/internal/repo"
	"tech-docs-ai/internal/vec"
)

func main() {
	// Initialize a new router
	r := chi.NewRouter()

	// Apply middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Create a new Ollama embedding client
	ollamaClient := emb.NewOllamaClient()

	// Initialize Qdrant client
	qdrantClient := vec.NewQdrantClient()

	// Initialize PostgreSQL care resource store
	postgresStore, err := repo.NewPostgresStore()
	if err != nil {
		log.Fatalf("Failed to initialize PostgreSQL store: %v", err)
	}
	defer postgresStore.Close()

	// Initialize Redis cache
	redisCache, err := cache.NewRedisCache()
	if err != nil {
		log.Fatalf("Failed to initialize Redis cache: %v", err)
	}
	defer redisCache.Close()

	// Initialize Kafka producer
	kafkaProducer := kafka.NewProducer()
	defer kafkaProducer.Close()

	// Create the main application service and handler
	svc := app.NewService(ollamaClient, qdrantClient, postgresStore, kafkaProducer, redisCache)
	handler := app.NewHandler(svc)
	wsHandler := app.NewWebSocketHandler(svc)

	// Set up API routes
	r.Post("/api/v1/chat", handler.HandleChat)
	r.Post("/api/v1/chat/history", handler.HandleChatWithHistory)
	r.Get("/api/v1/chat/history", handler.HandleGetChatHistory)
	r.Get("/api/v1/chat/insights", handler.HandleGetConversationInsights)
	r.Post("/api/v1/documents", handler.HandleAddDocument)
	r.Post("/api/v1/scrape", handler.HandleScrapeDocument)
	r.Get("/api/v1/documents/search", handler.HandleSearchDocuments)
	r.Post("/api/v1/tutorials/generate", handler.HandleGenerateTutorial)
	r.Post("/api/v1/tutorials/scrape-and-generate", handler.HandleScrapeAndGenerateTutorial)

	// WebSocket endpoint for real-time chat
	r.HandleFunc("/ws", wsHandler.HandleWebSocket)

	// Add health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy", "service": "tech-docs-ai"}`))
	})

	// Run the seeder
	app.Seeder()

	// Start the server
	port := "8080"
	log.Printf("Elderly Care AI Server starting on port %s", port)
	log.Printf("Ollama API URL: %s", getEnv("OLLAMA_API_URL", "http://localhost:11434"))
	log.Printf("Qdrant API URL: %s", getEnv("QDRANT_API_URL", "http://localhost:6333"))
	log.Printf("Kafka URL: %s", getEnv("KAFKA_URL", "localhost:9092"))
	log.Printf("Redis URL: %s", getEnv("REDIS_URL", "localhost:6379"))

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
