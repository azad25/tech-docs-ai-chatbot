package types

import (
	"time"
)

// Document represents a documentation or tutorial piece.
type Document struct {
	ID        string            `json:"id"`
	Title     string            `json:"title"`
	Content   string            `json:"content"`
	Category  string            `json:"category"`
	Tags      []string          `json:"tags"`
	Author    string            `json:"author"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Metadata  map[string]string `json:"metadata"`
}

// SearchResult represents a search result from vector database.
type SearchResult struct {
	ID       string                 `json:"id"`
	Score    float32                `json:"score"`
	Metadata map[string]interface{} `json:"metadata"`
}

// ScrapeJob represents a scraping job message.
type ScrapeJob struct {
	URL      string   `json:"url"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
	JobID    string   `json:"job_id"`
}

// ChatSession represents a chat session
type ChatSession struct {
	ID        string         `json:"id"`
	UserID    string         `json:"user_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Messages  []*ChatMessage `json:"messages"`
}

// ChatMessage represents a chat message
type ChatMessage struct {
	ID        string    `json:"id"`
	Role      string    `json:"role"` // "user" or "assistant"
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}
