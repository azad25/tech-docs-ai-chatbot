// cmd/worker/main.go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"tech-docs-ai/internal/emb"
	"tech-docs-ai/internal/kafka"
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

	// Create Kafka consumer for processing scraping jobs
	consumer := kafka.NewConsumer(postgresStore, ollamaClient, qdrantClient)
	defer consumer.Close()

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

	// Start the consumer
	if err := consumer.Start(ctx); err != nil {
		log.Fatalf("Consumer failed: %v", err)
	}
}
