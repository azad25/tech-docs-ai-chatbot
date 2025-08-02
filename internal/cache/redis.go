package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"tech-docs-ai/internal/types"

	"github.com/redis/go-redis/v9"
)

// Cache configuration constants
const (
	MaxMemoryPolicy    = "allkeys-lru"
	MaxMemoryBytes     = 1024 * 1024 * 1024 // 1GB
	MaxKeySize         = 1024               // 1KB
	MaxValueSize       = 5 * 1024 * 1024    // 5MB
	MaxConnectionRetry = 3
	RetryDelay         = time.Second * 2
)

// RedisCache is a Redis-based cache implementation with connection pooling.
type RedisCache struct {
	client *redis.Client
	mutex  sync.RWMutex
}

// NewRedisCache creates a new Redis cache client with improved configuration.
func NewRedisCache() (*RedisCache, error) {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		opts = &redis.Options{
			Addr:         redisURL,
			Password:     "",
			DB:           0,
			PoolSize:     10,
			MinIdleConns: 5,
			MaxRetries:   MaxConnectionRetry,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
			PoolTimeout:  4 * time.Second,
		}
	}

	client := redis.NewClient(opts)

	// Configure Redis with maxmemory and eviction policy
	for i := 0; i < MaxConnectionRetry; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Test connection
		if err := client.Ping(ctx).Err(); err != nil {
			if i == MaxConnectionRetry-1 {
				return nil, fmt.Errorf("failed to connect to Redis after %d retries: %w", MaxConnectionRetry, err)
			}
			time.Sleep(RetryDelay)
			continue
		}

		// Set maxmemory and policy
		if err := client.ConfigSet(ctx, "maxmemory", fmt.Sprintf("%d", MaxMemoryBytes)).Err(); err != nil {
			return nil, fmt.Errorf("failed to set Redis maxmemory: %w", err)
		}
		if err := client.ConfigSet(ctx, "maxmemory-policy", MaxMemoryPolicy).Err(); err != nil {
			return nil, fmt.Errorf("failed to set Redis maxmemory-policy: %w", err)
		}

		break
	}

	return &RedisCache{client: client}, nil
}

// validateKey checks if the key size is within limits
func (c *RedisCache) validateKey(key string) error {
	if len(key) > MaxKeySize {
		return fmt.Errorf("key size exceeds maximum allowed size of %d bytes", MaxKeySize)
	}
	return nil
}

// validateValue checks if the value size is within limits

// setWithValidation sets a value in Redis with size validation

// getWithValidation gets a value from Redis with error handling

// Cache keys
const (
	DocumentKeyPrefix     = "doc:"
	EmbeddingKeyPrefix    = "emb:"
	SearchResultKeyPrefix = "search:"
	ChatSessionKeyPrefix  = "chat:"
	DefaultTTL            = 24 * time.Hour
)

// DocumentCache caches document data
func (c *RedisCache) DocumentCache() DocumentCache {
	return &redisDocumentCache{cache: c}
}

// EmbeddingCache caches embedding vectors
func (c *RedisCache) EmbeddingCache() EmbeddingCache {
	return &redisEmbeddingCache{cache: c}
}

// SearchCache caches search results
func (c *RedisCache) SearchCache() SearchCache {
	return &redisSearchCache{cache: c}
}

// ChatCache caches chat sessions
func (c *RedisCache) ChatCache() ChatCache {
	return &redisChatCache{cache: c}
}

type RedisConfig struct {
	Addr         string
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
}

// Removed duplicate methods validateKey, validateValue, setWithValidation, and getWithValidation to keep only the first implementations
// These methods are already defined earlier in the file, so this block is removed to fix redeclaration errors.
// Removed the second validateKey method to fix duplicate declaration error
// This method is already defined earlier in the file
func (c *RedisCache) validateValue(value []byte) error {
	if len(value) > MaxValueSize {
		return fmt.Errorf("value size exceeds maximum allowed size of %d bytes", MaxValueSize)
	}
	return nil
}

func (c *RedisCache) setWithValidation(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if err := c.validateKey(key); err != nil {
		return err
	}

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	if err := c.validateValue(data); err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *RedisCache) getWithValidation(ctx context.Context, key string, value interface{}) error {
	if err := c.validateKey(key); err != nil {
		return err
	}

	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil // Cache miss
		}
		return fmt.Errorf("failed to get value from cache: %w", err)
	}

	if err := json.Unmarshal(data, value); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return nil
}

// DocumentCache interface for caching documents
type DocumentCache interface {
	Get(ctx context.Context, id string) (*types.Document, error)
	Set(ctx context.Context, doc *types.Document, ttl time.Duration) error
	Delete(ctx context.Context, id string) error
	InvalidateByCategory(ctx context.Context, category string) error
}

// EmbeddingCache interface for caching embeddings
type EmbeddingCache interface {
	Get(ctx context.Context, text string) ([]float32, error)
	Set(ctx context.Context, text string, embedding []float32, ttl time.Duration) error
	Delete(ctx context.Context, text string) error
}

// SearchCache interface for caching search results
type SearchCache interface {
	Get(ctx context.Context, query string, limit int) ([]*types.Document, error)
	Set(ctx context.Context, query string, limit int, docs []*types.Document, ttl time.Duration) error
	Delete(ctx context.Context, query string) error
}

// ChatCache interface for caching chat sessions
type ChatCache interface {
	GetSession(ctx context.Context, sessionID string) (*types.ChatSession, error)
	SetSession(ctx context.Context, session *types.ChatSession, ttl time.Duration) error
	AddMessage(ctx context.Context, sessionID string, message *types.ChatMessage) error
	GetHistory(ctx context.Context, sessionID string, limit int) ([]*types.ChatMessage, error)
}

// redisDocumentCache implements DocumentCache
type redisDocumentCache struct {
	cache *RedisCache
}

func (c *redisDocumentCache) Get(ctx context.Context, id string) (*types.Document, error) {
	key := DocumentKeyPrefix + id
	data, err := c.cache.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get document from cache: %w", err)
	}

	var doc types.Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("failed to unmarshal document: %w", err)
	}

	return &doc, nil
}

func (c *redisDocumentCache) Set(ctx context.Context, doc *types.Document, ttl time.Duration) error {
	key := DocumentKeyPrefix + doc.ID
	return c.cache.setWithValidation(ctx, key, doc, ttl)
}

func (c *redisDocumentCache) Delete(ctx context.Context, id string) error {
	key := DocumentKeyPrefix + id
	return c.cache.client.Del(ctx, key).Err()
}

func (c *redisDocumentCache) InvalidateByCategory(ctx context.Context, category string) error {
	// This is a simplified implementation
	// In a real scenario, you might want to maintain a separate index of documents by category
	pattern := DocumentKeyPrefix + "*"
	keys, err := c.cache.client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get keys: %w", err)
	}

	for _, key := range keys {
		// Get the document and check its category
		data, err := c.cache.client.Get(ctx, key).Bytes()
		if err != nil {
			continue
		}

		var doc types.Document
		if err := json.Unmarshal(data, &doc); err != nil {
			continue
		}

		if doc.Category == category {
			c.cache.client.Del(ctx, key)
		}
	}

	return nil
}

// redisEmbeddingCache implements EmbeddingCache
type redisEmbeddingCache struct {
	cache *RedisCache
}

func (c *redisEmbeddingCache) Get(ctx context.Context, text string) ([]float32, error) {
	key := EmbeddingKeyPrefix + hashText(text)
	data, err := c.cache.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("cache miss") // Return error for cache miss
		}
		return nil, fmt.Errorf("failed to get embedding from cache: %w", err)
	}

	var embedding []float32
	if err := json.Unmarshal(data, &embedding); err != nil {
		return nil, fmt.Errorf("failed to unmarshal embedding: %w", err)
	}

	return embedding, nil
}

func (c *redisEmbeddingCache) Set(ctx context.Context, text string, embedding []float32, ttl time.Duration) error {
	key := EmbeddingKeyPrefix + hashText(text)
	data, err := json.Marshal(embedding)
	if err != nil {
		return fmt.Errorf("failed to marshal embedding: %w", err)
	}

	if err := c.cache.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set embedding in cache: %w", err)
	}

	return nil
}

func (c *redisEmbeddingCache) Delete(ctx context.Context, text string) error {
	key := EmbeddingKeyPrefix + hashText(text)
	return c.cache.client.Del(ctx, key).Err()
}

// redisSearchCache implements SearchCache
type redisSearchCache struct {
	cache *RedisCache
}

func (c *redisSearchCache) Get(ctx context.Context, query string, limit int) ([]*types.Document, error) {
	// This is a simplified implementation for demonstration.
	// In a real search, you'd hash the query and limit.
	// For now, we'll just use a placeholder key.
	key := SearchResultKeyPrefix + "query:" + query + ":limit:" + fmt.Sprintf("%d", limit)
	data, err := c.cache.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get search results from cache: %w", err)
	}

	var docs []*types.Document
	if err := json.Unmarshal(data, &docs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search results: %w", err)
	}

	return docs, nil
}

func (c *redisSearchCache) Set(ctx context.Context, query string, limit int, docs []*types.Document, ttl time.Duration) error {
	// This is a simplified implementation for demonstration.
	// In a real search, you'd hash the query and limit.
	// For now, we'll just use a placeholder key.
	key := SearchResultKeyPrefix + "query:" + query + ":limit:" + fmt.Sprintf("%d", limit)
	data, err := json.Marshal(docs)
	if err != nil {
		return fmt.Errorf("failed to marshal search results: %w", err)
	}

	if err := c.cache.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set search results in cache: %w", err)
	}

	return nil
}

func (c *redisSearchCache) Delete(ctx context.Context, query string) error {
	// This is a simplified implementation for demonstration.
	// In a real search, you'd hash the query.
	// For now, we'll just use a placeholder key.
	key := SearchResultKeyPrefix + "query:" + query
	return c.cache.client.Del(ctx, key).Err()
}

// redisChatCache implements ChatCache
type redisChatCache struct {
	cache *RedisCache
}

func (c *redisChatCache) GetSession(ctx context.Context, sessionID string) (*types.ChatSession, error) {
	key := ChatSessionKeyPrefix + sessionID
	data, err := c.cache.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get chat session from cache: %w", err)
	}

	var session types.ChatSession
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal chat session: %w", err)
	}

	return &session, nil
}

func (c *redisChatCache) SetSession(ctx context.Context, session *types.ChatSession, ttl time.Duration) error {
	key := ChatSessionKeyPrefix + session.ID
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal chat session: %w", err)
	}

	if err := c.cache.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set chat session in cache: %w", err)
	}

	return nil
}

func (c *redisChatCache) AddMessage(ctx context.Context, sessionID string, message *types.ChatMessage) error {
	// Get existing session
	session, err := c.GetSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	if session == nil {
		session = &types.ChatSession{
			ID:        sessionID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Messages:  []*types.ChatMessage{},
		}
	}

	// Add message
	session.Messages = append(session.Messages, message)
	session.UpdatedAt = time.Now()

	// Save back to cache
	return c.SetSession(ctx, session, DefaultTTL)
}

func (c *redisChatCache) GetHistory(ctx context.Context, sessionID string, limit int) ([]*types.ChatMessage, error) {
	session, err := c.GetSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if session == nil {
		return []*types.ChatMessage{}, nil
	}

	// Return last N messages
	if len(session.Messages) <= limit {
		return session.Messages, nil
	}

	return session.Messages[len(session.Messages)-limit:], nil
}

// hashText creates a simple hash for text to use as cache keys
func hashText(text string) string {
	hash := 0
	for _, char := range text {
		hash = ((hash << 5) - hash) + int(char)
		hash = hash & hash // Convert to 32-bit integer
	}
	return fmt.Sprintf("%d", hash)
}

// Close closes the Redis client and its connections
func (rc *RedisCache) Close() error {
	if rc.client != nil {
		return rc.client.Close()
	}
	return nil
}

// Health checks if Redis is healthy
func (rc *RedisCache) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := rc.client.Ping(ctx).Result()
	return err
}
