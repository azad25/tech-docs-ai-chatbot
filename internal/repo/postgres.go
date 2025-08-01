package repo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"tech-docs-ai/internal/types"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

// PostgresStore is a PostgreSQL implementation of document storage.
type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore creates a new PostgreSQL document store.
func NewPostgresStore() (*PostgresStore, error) {
	// Get database connection details from environment
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		port = "5432"
	}

	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		user = "user"
	}

	password := os.Getenv("POSTGRES_PASSWORD")
	if password == "" {
		password = "password"
	}

	dbname := os.Getenv("POSTGRES_DB")
	if dbname == "" {
		dbname = "tech-docs-ai"
	}

	// Create connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	store := &PostgresStore{db: db}

	// Initialize database schema
	if err := store.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return store, nil
}

// initSchema creates the necessary database tables.
func (p *PostgresStore) initSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS documents (
		id VARCHAR(255) PRIMARY KEY,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		category VARCHAR(100) NOT NULL,
		tags TEXT[] NOT NULL,
		author VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL,
		metadata JSONB
	);

	CREATE INDEX IF NOT EXISTS idx_documents_category ON documents(category);
	CREATE INDEX IF NOT EXISTS idx_documents_tags ON documents USING GIN(tags);
	CREATE INDEX IF NOT EXISTS idx_documents_created_at ON documents(created_at);
	`

	_, err := p.db.Exec(query)
	return err
}

// StoreDocument stores a document in the database.
func (p *PostgresStore) StoreDocument(doc *types.Document) error {
	// Generate ID if not provided
	if doc.ID == "" {
		doc.ID = fmt.Sprintf("doc_%d", time.Now().UnixNano())
	}

	// Set timestamps
	now := time.Now()
	if doc.CreatedAt.IsZero() {
		doc.CreatedAt = now
	}
	doc.UpdatedAt = now

	// Convert metadata to JSON
	metadataJSON, err := json.Marshal(doc.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `
		INSERT INTO documents (id, title, content, category, tags, author, created_at, updated_at, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			content = EXCLUDED.content,
			category = EXCLUDED.category,
			tags = EXCLUDED.tags,
			author = EXCLUDED.author,
			updated_at = EXCLUDED.updated_at,
			metadata = EXCLUDED.metadata
	`

	_, err = p.db.Exec(query,
		doc.ID,
		doc.Title,
		doc.Content,
		doc.Category,
		pq.Array(doc.Tags),
		doc.Author,
		doc.CreatedAt,
		doc.UpdatedAt,
		metadataJSON,
	)

	if err != nil {
		return fmt.Errorf("failed to store document: %w", err)
	}

	return nil
}

// GetDocument retrieves a document by ID.
func (p *PostgresStore) GetDocument(id string) (*types.Document, error) {
	query := `
	SELECT id, title, content, category, tags, author, created_at, updated_at, metadata
	FROM documents WHERE id = $1
	`

	var doc types.Document
	var metadataJSON []byte

	err := p.db.QueryRow(query, id).Scan(
		&doc.ID,
		&doc.Title,
		&doc.Content,
		&doc.Category,
		pq.Array(&doc.Tags),
		&doc.Author,
		&doc.CreatedAt,
		&doc.UpdatedAt,
		&metadataJSON,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	// Parse metadata
	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &doc.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	return &doc, nil
}

// SearchDocuments searches for documents by text query.
func (p *PostgresStore) SearchDocuments(query string, limit int) ([]*types.Document, error) {
	// Simple text search using ILIKE
	searchQuery := `
	SELECT id, title, content, category, tags, author, created_at, updated_at, metadata
	FROM documents 
	WHERE title ILIKE $1 OR content ILIKE $1 OR category ILIKE $1
	ORDER BY created_at DESC
	LIMIT $2
	`

	searchTerm := "%" + query + "%"
	rows, err := p.db.Query(searchQuery, searchTerm, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search documents: %w", err)
	}
	defer rows.Close()

	var documents []*types.Document
	for rows.Next() {
		var doc types.Document
		var metadataJSON []byte

		err := rows.Scan(
			&doc.ID,
			&doc.Title,
			&doc.Content,
			&doc.Category,
			pq.Array(&doc.Tags),
			&doc.Author,
			&doc.CreatedAt,
			&doc.UpdatedAt,
			&metadataJSON,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan document: %w", err)
		}

		// Parse metadata
		if len(metadataJSON) > 0 {
			if err := json.Unmarshal(metadataJSON, &doc.Metadata); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
		}

		documents = append(documents, &doc)
	}

	return documents, nil
}

// Close closes the database connection.
func (p *PostgresStore) Close() error {
	return p.db.Close()
}
