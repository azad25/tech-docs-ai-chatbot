package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"tech-docs-ai/internal/emb"
	"tech-docs-ai/internal/repo"
	"tech-docs-ai/internal/scraper"
	"tech-docs-ai/internal/types"
	"tech-docs-ai/internal/vec"

	"github.com/segmentio/kafka-go"
)

// Consumer is a Kafka consumer for processing scraping jobs.
type Consumer struct {
	reader          *kafka.Reader
	w3schoolsScraper *scraper.W3SchoolsScraper
	universalScraper *scraper.UniversalScraper
	docStore        *repo.PostgresStore
	embClient       *emb.OllamaClient
	vecClient       *vec.QdrantClient
	workerPool      *WorkerPool
}

// NewConsumer creates a new Kafka consumer.
func NewConsumer(docStore *repo.PostgresStore, embClient *emb.OllamaClient, vecClient *vec.QdrantClient) *Consumer {
	kafkaURL := os.Getenv("KAFKA_URL")
	if kafkaURL == "" {
		kafkaURL = "localhost:9092"
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaURL},
		Topic:    "scrape-jobs",
		GroupID:  "scraper-workers",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &Consumer{
		reader:           reader,
		w3schoolsScraper: scraper.NewW3SchoolsScraper(),
		universalScraper: scraper.NewUniversalScraper(),
		docStore:         docStore,
		embClient:        embClient,
		vecClient:        vecClient,
		workerPool:       NewWorkerPool(5), // 5 workers
	}
}

// Start starts consuming messages from Kafka.
func (c *Consumer) Start(ctx context.Context) error {
	log.Println("Starting Kafka consumer...")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			m, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}

			// Process message in worker pool
			c.workerPool.Submit(func() {
				c.processMessage(m)
			})
		}
	}
}

// processMessage processes a single Kafka message.
func (c *Consumer) processMessage(m kafka.Message) {
	log.Printf("Processing message: %s", string(m.Value))

	var job types.ScrapeJob
	if err := json.Unmarshal(m.Value, &job); err != nil {
		log.Printf("Failed to unmarshal job: %v", err)
		return
	}

	// Check if content already exists for this URL/topic
	existingDocs, err := c.docStore.SearchDocuments(job.URL, 1)
	if err == nil && len(existingDocs) > 0 {
		log.Printf("Content already exists for URL: %s, skipping scrape", job.URL)
		return
	}

	// Choose appropriate scraper based on URL
	var content *scraper.ScrapedContent
	var doc *types.Document

	if strings.Contains(job.URL, "w3schools.com") {
		// Use W3Schools-specific scraper
		var err error
		content, err = c.w3schoolsScraper.ScrapePage(job.URL)
		if err != nil {
			log.Printf("Failed to scrape %s with W3Schools scraper: %v", job.URL, err)
			return
		}
		doc = c.w3schoolsScraper.ConvertToDocument(content)
	} else {
		// Use universal scraper for all other URLs
		var err error
		content, err = c.universalScraper.ScrapePage(job.URL)
		if err != nil {
			log.Printf("Failed to scrape %s with universal scraper: %v", job.URL, err)
			return
		}
		doc = c.universalScraper.ConvertToDocument(content)
	}

	// Override category and tags if provided in job
	if job.Category != "" {
		doc.Category = job.Category
	}
	if len(job.Tags) > 0 {
		doc.Tags = append(doc.Tags, job.Tags...)
	}

	// Store document and vector
	if err := c.storeDocumentWithVector(doc, job.URL); err != nil {
		log.Printf("Failed to store document: %v", err)
		return
	}

	log.Printf("Successfully processed job %s for URL: %s", job.JobID, job.URL)
}

// storeDocumentWithVector stores a document in both database and vector store.
func (c *Consumer) storeDocumentWithVector(doc *types.Document, url string) error {
	// Generate embedding for the document content
	vector, err := c.embClient.Embed(doc.Content)
	if err != nil {
		return fmt.Errorf("failed to embed document: %w", err)
	}

	// Store document in database
	if err := c.docStore.StoreDocument(doc); err != nil {
		return fmt.Errorf("failed to store document: %w", err)
	}

	// Store vector with document metadata
	source := "universal"
	if strings.Contains(url, "w3schools.com") {
		source = "w3schools"
	}
	
	metadata := map[string]interface{}{
		"document_id": doc.ID,
		"title":       doc.Title,
		"category":    doc.Category,
		"tags":        doc.Tags,
		"author":      doc.Author,
		"source":      source,
	}

	if err := c.vecClient.StoreVector(vector, metadata); err != nil {
		return fmt.Errorf("failed to store vector: %w", err)
	}

	return nil
}

// Close closes the Kafka consumer.
func (c *Consumer) Close() error {
	return c.reader.Close()
}

// WorkerPool manages a pool of workers for processing jobs.
type WorkerPool struct {
	workers  int
	jobQueue chan func()
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewWorkerPool creates a new worker pool.
func NewWorkerPool(workers int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	pool := &WorkerPool{
		workers:  workers,
		jobQueue: make(chan func(), workers*2),
		ctx:      ctx,
		cancel:   cancel,
	}

	// Start workers
	for i := 0; i < workers; i++ {
		pool.wg.Add(1)
		go pool.worker(i)
	}

	return pool
}

// worker is a single worker goroutine.
func (p *WorkerPool) worker(id int) {
	defer p.wg.Done()

	log.Printf("Worker %d started", id)

	for {
		select {
		case job := <-p.jobQueue:
			job()
		case <-p.ctx.Done():
			log.Printf("Worker %d stopping", id)
			return
		}
	}
}

// Submit submits a job to the worker pool.
func (p *WorkerPool) Submit(job func()) {
	select {
	case p.jobQueue <- job:
		// Job submitted successfully
	case <-p.ctx.Done():
		log.Println("Worker pool is shutting down, job rejected")
	}
}

// Shutdown gracefully shuts down the worker pool.
func (p *WorkerPool) Shutdown() {
	log.Println("Shutting down worker pool...")
	p.cancel()
	p.wg.Wait()
	log.Println("Worker pool shutdown complete")
}
