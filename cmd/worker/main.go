// cmd/worker/main.go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"tech-docs-ai/internal/emb"
	"tech-docs-ai/internal/repo"
	"tech-docs-ai/internal/vec"
)

func main() {
	log.Println("Starting Tech Docs AI Worker Service...")

	// Initialize components
	ollamaClient := emb.NewOllamaClient()
	qdrantClient := vec.NewQdrantClient()

	postgresStore, err := repo.NewPostgresStore()
	if err != nil {
		log.Fatalf("Failed to initialize PostgreSQL store: %v", err)
	}
	defer postgresStore.Close()

	// Create a simple worker that processes scraping jobs
	// Note: We'll use concrete types instead of interfaces for now
	worker := NewScrapingWorker(ollamaClient, qdrantClient, postgresStore)

	// Set up graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal %v, shutting down...", sig)
		cancel()
	}()

	// Start the worker
	if err := worker.Start(ctx); err != nil {
		log.Fatalf("Worker failed: %v", err)
	}
}

// ScrapingWorker is a simplified worker that processes scraping jobs
type ScrapingWorker struct {
	embClient *emb.OllamaClient
	vecClient *vec.QdrantClient
	docStore  *repo.PostgresStore
}

// NewScrapingWorker creates a new scraping worker
func NewScrapingWorker(embClient *emb.OllamaClient, vecClient *vec.QdrantClient, docStore *repo.PostgresStore) *ScrapingWorker {
	return &ScrapingWorker{
		embClient: embClient,
		vecClient: vecClient,
		docStore:  docStore,
	}
}

// Start starts the worker (placeholder for now)
func (w *ScrapingWorker) Start(ctx context.Context) error {
	log.Println("Worker started. Press Ctrl+C to stop.")

	// Keep the worker running until context is cancelled
	<-ctx.Done()

	log.Println("Worker stopped.")
	return nil
}
