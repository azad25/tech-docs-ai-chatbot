package vec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"tech-docs-ai/internal/types"
)

// QdrantClient is a real implementation for vector storage using Qdrant.
type QdrantClient struct {
	apiURL     string
	collection string
	httpClient *http.Client
}

// NewQdrantClient creates and returns a new QdrantClient.
func NewQdrantClient() *QdrantClient {
	apiURL := os.Getenv("QDRANT_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:6333"
	}

	collection := os.Getenv("QDRANT_COLLECTION")
	if collection == "" {
		collection = "tech_docs_knowledge"
	}

	client := &QdrantClient{
		apiURL:     apiURL,
		collection: collection,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	// Initialize the collection if it doesn't exist
	client.ensureCollection()

	return client
}

// ensureCollection creates the collection if it doesn't exist.
func (c *QdrantClient) ensureCollection() error {
	// Check if collection exists
	checkURL := fmt.Sprintf("%s/collections/%s", c.apiURL, c.collection)
	resp, err := c.httpClient.Get(checkURL)
	if err != nil {
		return fmt.Errorf("failed to check collection: %w", err)
	}
	defer resp.Body.Close()

	// If collection doesn't exist, create it
	if resp.StatusCode == http.StatusNotFound {
		createURL := fmt.Sprintf("%s/collections/%s", c.apiURL, c.collection)
		createBody := map[string]interface{}{
			"vectors": map[string]interface{}{
				"size":     768, // Default size for nomic-embed-text
				"distance": "Cosine",
			},
		}

		jsonData, err := json.Marshal(createBody)
		if err != nil {
			return fmt.Errorf("failed to marshal create collection request: %w", err)
		}

		req, err := http.NewRequest("PUT", createURL, bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")

		createResp, err := c.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("failed to create collection: %w", err)
		}
		defer createResp.Body.Close()

		if createResp.StatusCode != http.StatusOK && createResp.StatusCode != http.StatusCreated {
			return fmt.Errorf("failed to create collection with status: %d", createResp.StatusCode)
		}
	}

	return nil
}

// StoreVector stores a vector in Qdrant with metadata.
func (c *QdrantClient) StoreVector(vector []float32, metadata map[string]interface{}) error {
	if len(vector) == 0 {
		return fmt.Errorf("empty vector")
	}

	// Generate a unique numeric ID for the vector
	vectorID := time.Now().UnixNano()

	// Prepare the upsert request
	upsertBody := map[string]interface{}{
		"points": []map[string]interface{}{
			{
				"id":      vectorID,
				"vector":  vector,
				"payload": metadata,
			},
		},
	}

	jsonData, err := json.Marshal(upsertBody)
	if err != nil {
		return fmt.Errorf("failed to marshal upsert request: %w", err)
	}

	upsertURL := fmt.Sprintf("%s/collections/%s/points", c.apiURL, c.collection)
	req, err := http.NewRequest("PUT", upsertURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to store vector: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to store vector with status: %d", resp.StatusCode)
	}

	return nil
}

// SearchVector searches for similar vectors in Qdrant and returns results with metadata.
func (c *QdrantClient) SearchVector(vector []float32, limit int) ([]types.SearchResult, error) {
	if len(vector) == 0 {
		return nil, fmt.Errorf("empty vector")
	}

	// Prepare the search request
	searchBody := map[string]interface{}{
		"vector":       vector,
		"limit":        limit,
		"with_payload": true,
		"with_vectors": false,
	}

	jsonData, err := json.Marshal(searchBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search request: %w", err)
	}

	searchURL := fmt.Sprintf("%s/collections/%s/points/search", c.apiURL, c.collection)
	resp, err := c.httpClient.Post(searchURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to search vectors: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to search vectors with status: %d", resp.StatusCode)
	}

	var searchResult map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	var results []types.SearchResult
	if result, ok := searchResult["result"].([]interface{}); ok {
		for _, item := range result {
			if point, ok := item.(map[string]interface{}); ok {
				result := types.SearchResult{}

				// Extract ID
				if id, ok := point["id"].(string); ok {
					result.ID = id
				}

				// Extract score
				if score, ok := point["score"].(float64); ok {
					result.Score = float32(score)
				}

				// Extract payload/metadata
				if payload, ok := point["payload"].(map[string]interface{}); ok {
					result.Metadata = payload
				}

				results = append(results, result)
			}
		}
	}

	return results, nil
}
