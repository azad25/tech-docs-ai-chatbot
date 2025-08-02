// internal/app/service.go
package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"tech-docs-ai/internal/cache"
	"tech-docs-ai/internal/types"
)

// Service holds the core business logic of the application.
// It uses interfaces to communicate with other layers.
type Service struct {
	embClient embClient
	vecClient vecClient
	docStore  docStore
	kafkaProd kafkaProducer
	cache     *cache.RedisCache
}

// NewService creates a new Service instance.
func NewService(embClient embClient, vecClient vecClient, docStore docStore, kafkaProd kafkaProducer, cache *cache.RedisCache) *Service {
	return &Service{
		embClient: embClient,
		vecClient: vecClient,
		docStore:  docStore,
		kafkaProd: kafkaProd,
		cache:     cache,
	}
}

// embClient is an interface that defines the contract for an embedding client.
type embClient interface {
	// Embed will take a text and return its vector representation.
	Embed(text string) ([]float32, error)
	// Chat will take a message and return a response from the LLM.
	Chat(message string) (string, error)
}

// vecClient is an interface for vector storage.
type vecClient interface {
	StoreVector(vector []float32, metadata map[string]interface{}) error
	SearchVector(vector []float32, limit int) ([]types.SearchResult, error)
}

// docStore is an interface for document storage.
type docStore interface {
	StoreDocument(doc *types.Document) error
	GetDocument(id string) (*types.Document, error)
	SearchDocuments(query string, limit int) ([]*types.Document, error)
}

// kafkaProducer is an interface for Kafka messaging.
type kafkaProducer interface {
	SendMessage(topic string, message []byte) error
}

// Chat handles the core logic for a chat interaction with RAG and caching.
func (s *Service) Chat(message string) (string, error) {
	ctx := context.Background()

	// Step 1: Check cache for embedding
	embeddingCache := s.cache.EmbeddingCache()
	queryVector, err := embeddingCache.Get(ctx, message)
	if err != nil {
		// Cache miss, generate embedding
		queryVector, err = s.embClient.Embed(message)
		if err != nil {
			return "", fmt.Errorf("failed to embed query: %w", err)
		}
		// Cache the embedding for future use
		embeddingCache.Set(ctx, message, queryVector, cache.DefaultTTL)
	}

	// Step 2: Search for relevant documents in vector database
	searchResults, err := s.vecClient.SearchVector(queryVector, 5)
	if err != nil {
		return "", fmt.Errorf("failed to search vectors: %w", err)
	}

	// Step 3: Retrieve relevant documents from cache/database
	var contextDocs []string
	var tutorialData []string
	var hasRelevantContent bool

	for _, result := range searchResults {
		if docID, ok := result.Metadata["document_id"].(string); ok {
			doc, err := s.getDocumentWithCache(ctx, docID)
			if err == nil && doc != nil {
				// Check if the content is relevant (score threshold)
				if result.Score > 0.7 { // Adjust threshold as needed
					contextDocs = append(contextDocs, fmt.Sprintf("Title: %s\nContent: %s", doc.Title, doc.Content))
					tutorialData = append(tutorialData, doc.Content)
					hasRelevantContent = true
				}
			}
		}
	}

	// Step 4: Build context-aware prompt for tutorial generation
	context := ""
	if hasRelevantContent && len(contextDocs) > 0 {
		context = "Based on the following relevant documentation:\n\n" + strings.Join(contextDocs, "\n\n") + "\n\n"
	}

	// Step 5: Generate tutorial response using LLM with context
	var tutorialPrompt string
	if hasRelevantContent {
		tutorialPrompt = fmt.Sprintf("%sUser Question: %s\n\nPlease generate a comprehensive tutorial in Markdown format based on the documentation above. Structure your response as follows:\n\n# [Topic Name] Tutorial\n\n## Overview\n[Brief 2-3 sentence explanation of the topic]\n\n## What You'll Learn\n- [Learning objective 1]\n- [Learning objective 2]\n- [Learning objective 3]\n\n## Prerequisites\n- [Prerequisite 1]\n- [Prerequisite 2]\n\n## Step-by-Step Guide\n\n### Step 1: [First Step]\n[Detailed explanation with examples]\n\n### Step 2: [Second Step]\n[Detailed explanation with examples]\n\n### Step 3: [Third Step]\n[Detailed explanation with examples]\n\n## Code Examples\n\n### Basic Example\n```[language]\n[Code example here]\n```\n\n### Advanced Example\n```[language]\n[More complex code example]\n```\n\n## Best Practices\n- [Best practice 1]\n- [Best practice 2]\n- [Best practice 3]\n\n## Common Pitfalls to Avoid\n- [Pitfall 1]\n- [Pitfall 2]\n\n## Summary\n[Brief summary of what was covered]\n\n## Next Steps\n- [What to learn next 1]\n- [What to learn next 2]\n\nKeep the tutorial comprehensive, well-structured, and beginner-friendly. Use proper Markdown formatting with headers, code blocks, and bullet points.", context, message)
	} else {
		// No relevant content found, provide general response
		tutorialPrompt = fmt.Sprintf(`User Question: %s

Please provide a helpful and informative tutorial in Markdown format about this topic. If you don't have specific information, provide general guidance and suggest where they might find more detailed information.

Structure your response as a well-formatted Markdown tutorial with:

# [Topic Name]

## Overview
[Brief explanation]

## Key Concepts
- [Concept 1]
- [Concept 2]

## Getting Started
[Basic steps]

## Examples
[Practical examples]

## Resources
[Where to learn more]

Use proper Markdown formatting with headers, code blocks, and bullet points.`, message)
	}

	response, err := s.embClient.Chat(tutorialPrompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate tutorial response: %w", err)
	}

	// Step 6: Store the response in vector database for future learning
	go s.storeResponseForLearning(message, response, hasRelevantContent)

	return response, nil
}

// storeResponseForLearning stores the LLM response in the vector database for future learning
func (s *Service) storeResponseForLearning(userQuery, llmResponse string, wasBasedOnScrapedData bool) {
	ctx := context.Background()

	// Generate embedding for the response
	responseVector, err := s.embClient.Embed(llmResponse)
	if err != nil {
		log.Printf("Failed to embed response for learning: %v", err)
		return
	}

	// Create a document for the response
	responseDoc := &types.Document{
		ID:        fmt.Sprintf("response_%d", time.Now().UnixNano()),
		Title:     fmt.Sprintf("AI Response: %s", truncateString(userQuery, 50)),
		Content:   llmResponse,
		Category:  "AI_Response",
		Tags:      []string{"ai-response", "user-generated", "learning"},
		Author:    "AI_Assistant",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Metadata: map[string]string{
			"user_query":           userQuery,
			"was_based_on_scraped": fmt.Sprintf("%t", wasBasedOnScrapedData),
			"response_type":        "tutorial",
		},
	}

	// Store in database
	if err := s.docStore.StoreDocument(responseDoc); err != nil {
		log.Printf("Failed to store response document: %v", err)
		return
	}

	// Cache the document
	docCache := s.cache.DocumentCache()
	docCache.Set(ctx, responseDoc, cache.DefaultTTL)

	// Store vector with metadata
	metadata := map[string]interface{}{
		"document_id":   responseDoc.ID,
		"title":         responseDoc.Title,
		"category":      responseDoc.Category,
		"tags":          responseDoc.Tags,
		"author":        responseDoc.Author,
		"source":        "ai-response",
		"user_query":    userQuery,
		"response_type": "tutorial",
		"learning_data": true,
	}

	if err := s.vecClient.StoreVector(responseVector, metadata); err != nil {
		log.Printf("Failed to store response vector: %v", err)
		return
	}

	log.Printf("Stored AI response for learning: %s", responseDoc.ID)
}

// truncateString truncates a string to the specified length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// getDocumentWithCache retrieves a document with caching
func (s *Service) getDocumentWithCache(ctx context.Context, id string) (*types.Document, error) {
	docCache := s.cache.DocumentCache()

	// Try cache first
	doc, err := docCache.Get(ctx, id)
	if err == nil && doc != nil {
		return doc, nil
	}

	// Cache miss, get from database
	doc, err = s.docStore.GetDocument(id)
	if err != nil {
		return nil, err
	}

	// Cache the document for future use
	if doc != nil {
		docCache.Set(ctx, doc, cache.DefaultTTL)
	}

	return doc, nil
}

// AddDocument adds a new documentation piece to the system with caching.
func (s *Service) AddDocument(doc *types.Document) error {
	ctx := context.Background()

	// Generate embedding for the document content
	embeddingCache := s.cache.EmbeddingCache()
	vector, err := embeddingCache.Get(ctx, doc.Content)
	if err != nil {
		// Cache miss, generate embedding
		vector, err = s.embClient.Embed(doc.Content)
		if err != nil {
			return fmt.Errorf("failed to embed document: %w", err)
		}
		// Cache the embedding
		embeddingCache.Set(ctx, doc.Content, vector, cache.DefaultTTL)
	}

	// Store document in database
	if err := s.docStore.StoreDocument(doc); err != nil {
		return fmt.Errorf("failed to store document: %w", err)
	}

	// Cache the document
	docCache := s.cache.DocumentCache()
	docCache.Set(ctx, doc, cache.DefaultTTL)

	// Store vector with document metadata
	metadata := map[string]interface{}{
		"document_id": doc.ID,
		"title":       doc.Title,
		"category":    doc.Category,
		"tags":        doc.Tags,
		"author":      doc.Author,
	}

	if err := s.vecClient.StoreVector(vector, metadata); err != nil {
		return fmt.Errorf("failed to store vector: %w", err)
	}

	// Invalidate search cache for this category
	searchCache := s.cache.SearchCache()
	searchCache.Delete(ctx, doc.Category)

	return nil
}

// ScrapeDocument queues a scraping job via Kafka.
func (s *Service) ScrapeDocument(url, category string, tags []string) error {
	job := types.ScrapeJob{
		URL:      url,
		Category: category,
		Tags:     tags,
		JobID:    fmt.Sprintf("job_%d", time.Now().UnixNano()),
	}

	// Serialize the job
	jobData, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal scrape job: %w", err)
	}

	// Send to Kafka topic
	if err := s.kafkaProd.SendMessage("scrape-jobs", jobData); err != nil {
		return fmt.Errorf("failed to send scrape job to Kafka: %w", err)
	}

	return nil
}

// SearchDocuments searches for documents by text query with caching.
func (s *Service) SearchDocuments(query string, limit int) ([]*types.Document, error) {
	ctx := context.Background()
	searchCache := s.cache.SearchCache()

	// Try cache first
	docs, err := searchCache.Get(ctx, query, limit)
	if err == nil && docs != nil {
		return docs, nil
	}

	// Cache miss, search database
	docs, err = s.docStore.SearchDocuments(query, limit)
	if err != nil {
		return nil, err
	}

	// Cache the search results
	if docs != nil {
		searchCache.Set(ctx, query, limit, docs, cache.DefaultTTL)
	}

	return docs, nil
}

// GetChatSession retrieves or creates a chat session
func (s *Service) GetChatSession(sessionID string) (*types.ChatSession, error) {
	ctx := context.Background()
	chatCache := s.cache.ChatCache()

	session, err := chatCache.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if session == nil {
		// Create new session
		session = &types.ChatSession{
			ID:        sessionID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Messages:  []*types.ChatMessage{},
		}
		chatCache.SetSession(ctx, session, cache.DefaultTTL)
	}

	return session, nil
}

// AddChatMessage adds a message to a chat session
func (s *Service) AddChatMessage(sessionID string, role, content string) error {
	ctx := context.Background()
	chatCache := s.cache.ChatCache()

	message := &types.ChatMessage{
		ID:        fmt.Sprintf("msg_%d", time.Now().UnixNano()),
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
	}

	return chatCache.AddMessage(ctx, sessionID, message)
}

// GetChatHistory retrieves chat history for a session
func (s *Service) GetChatHistory(sessionID string, limit int) ([]*types.ChatMessage, error) {
	ctx := context.Background()
	chatCache := s.cache.ChatCache()

	return chatCache.GetHistory(ctx, sessionID, limit)
}

// GenerateTutorialFromScrapedData generates a tutorial from scraped content
func (s *Service) GenerateTutorialFromScrapedData(url, topic string) (string, error) {
	// Search for existing documents related to this topic
	docs, err := s.docStore.SearchDocuments(topic, 10)
	if err != nil {
		return "", fmt.Errorf("failed to search for existing documents: %w", err)
	}

	var tutorialContent []string
	var contextDocs []string

	// Collect content from existing documents
	for _, doc := range docs {
		if strings.Contains(strings.ToLower(doc.Title), strings.ToLower(topic)) ||
			strings.Contains(strings.ToLower(doc.Category), strings.ToLower(topic)) {
			tutorialContent = append(tutorialContent, doc.Content)
			contextDocs = append(contextDocs, fmt.Sprintf("Title: %s\nContent: %s", doc.Title, doc.Content))
		}
	}

	// If no existing content, trigger scraping
	if len(tutorialContent) == 0 {
		// Queue scraping job
		if err := s.ScrapeDocument(url, topic, []string{"tutorial", "documentation"}); err != nil {
			return "", fmt.Errorf("failed to queue scraping job: %w", err)
		}

		return fmt.Sprintf("I've queued a scraping job for %s. The tutorial will be generated once the content is scraped and processed. Please try again in a few minutes.", topic), nil
	}

	// Generate tutorial from collected content
	context := ""
	if len(contextDocs) > 0 {
		context = "Based on the following scraped documentation:\n\n" + strings.Join(contextDocs, "\n\n") + "\n\n"
	}

	tutorialPrompt := fmt.Sprintf("%sGenerate a comprehensive tutorial for: %s\n\nPlease create a well-structured tutorial in Markdown format based on the scraped documentation above. Structure your response as follows:\n\n# Complete Tutorial: %s\n\n## Overview\n[Provide a clear, concise overview of the topic]\n\n## Prerequisites\n[List any prerequisites or basic knowledge needed]\n\n## Step-by-Step Guide\n\n### Step 1: [First Step]\n[Detailed explanation with examples]\n\n### Step 2: [Second Step]\n[Detailed explanation with examples]\n\n### Step 3: [Third Step]\n[Detailed explanation with examples]\n\n## Code Examples\n\n### Basic Example\n```[language]\n[Provide practical code examples]\n```\n\n### Advanced Example\n```[language]\n[More complex examples]\n```\n\n## Best Practices\n- [Best practice 1]\n- [Best practice 2]\n- [Best practice 3]\n\n## Common Pitfalls to Avoid\n- [Pitfall 1]\n- [Pitfall 2]\n\n## Summary\n[Brief summary of what was covered]\n\n## Next Steps\n[Suggest what to learn next]\n\nMake the tutorial comprehensive yet easy to follow, with practical examples and clear explanations. Use proper Markdown formatting.", context, topic, topic)

	response, err := s.embClient.Chat(tutorialPrompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate tutorial: %w", err)
	}

	return response, nil
}

// ScrapeAndGenerateTutorial scrapes content and immediately generates a tutorial
func (s *Service) ScrapeAndGenerateTutorial(url, topic string) (string, error) {
	// First, try to get existing content
	docs, err := s.docStore.SearchDocuments(topic, 5)
	if err != nil {
		return "", fmt.Errorf("failed to search documents: %w", err)
	}

	var tutorialContent []string
	var contextDocs []string

	// Collect relevant content
	for _, doc := range docs {
		if strings.Contains(strings.ToLower(doc.Title), strings.ToLower(topic)) ||
			strings.Contains(strings.ToLower(doc.Category), strings.ToLower(topic)) {
			tutorialContent = append(tutorialContent, doc.Content)
			contextDocs = append(contextDocs, fmt.Sprintf("Title: %s\nContent: %s", doc.Title, doc.Content))
		}
	}

	// If we have content, generate tutorial immediately
	if len(tutorialContent) > 0 {
		context := "Based on the following scraped documentation:\n\n" + strings.Join(contextDocs, "\n\n") + "\n\n"

		tutorialPrompt := fmt.Sprintf("%sGenerate a quick tutorial for: %s\n\nPlease create a concise tutorial in Markdown format based on the scraped documentation above. Structure your response as follows:\n\n# Quick Tutorial: %s\n\n## What is %s?\n[Brief explanation]\n\n## Key Concepts:\n- [Concept 1]\n- [Concept 2]\n- [Concept 3]\n\n## Basic Example:\n```[language]\n[Simple, practical example]\n```\n\n## Common Use Cases:\n- [Use case 1]\n- [Use case 2]\n\n## Tips:\n- [Tip 1]\n- [Tip 2]\n\nKeep it concise and practical for beginners. Use proper Markdown formatting.", context, topic, topic, topic)

		response, err := s.embClient.Chat(tutorialPrompt)
		if err != nil {
			return "", fmt.Errorf("failed to generate tutorial: %w", err)
		}

		return response, nil
	}

	// If no content exists, queue scraping and return a message
	if err := s.ScrapeDocument(url, topic, []string{"tutorial", "documentation"}); err != nil {
		return "", fmt.Errorf("failed to queue scraping job: %w", err)
	}

	return fmt.Sprintf("I'm scraping content for %s. The tutorial will be available shortly. Please try again in a few minutes.", topic), nil
}

// ChatWithHistory handles chat with conversation history
func (s *Service) ChatWithHistory(sessionID, message string) (string, error) {
	ctx := context.Background()

	// Get or create chat session
	_, err := s.GetChatSession(sessionID)
	if err != nil {
		return "", fmt.Errorf("failed to get chat session: %w", err)
	}

	// Get conversation history
	history, err := s.GetChatHistory(sessionID, 10) // Last 10 messages
	if err != nil {
		return "", fmt.Errorf("failed to get chat history: %w", err)
	}

	// Build conversation context
	var conversationContext strings.Builder
	if len(history) > 0 {
		conversationContext.WriteString("Previous conversation:\n")
		for _, msg := range history {
			conversationContext.WriteString(fmt.Sprintf("%s: %s\n", msg.Role, msg.Content))
		}
		conversationContext.WriteString("\n")
	}

	// Search for relevant content in vector database
	queryVector, err := s.embClient.Embed(message)
	if err != nil {
		return "", fmt.Errorf("failed to embed query: %w", err)
	}

	searchResults, err := s.vecClient.SearchVector(queryVector, 5)
	if err != nil {
		return "", fmt.Errorf("failed to search vectors: %w", err)
	}

	// Retrieve relevant documents
	var contextDocs []string
	var hasRelevantContent bool

	for _, result := range searchResults {
		if docID, ok := result.Metadata["document_id"].(string); ok {
			doc, err := s.getDocumentWithCache(ctx, docID)
			if err == nil && doc != nil {
				if result.Score > 0.7 {
					contextDocs = append(contextDocs, fmt.Sprintf("Title: %s\nContent: %s", doc.Title, doc.Content))
					hasRelevantContent = true
				}
			}
		}
	}

	// Build comprehensive prompt with history and context
	var prompt strings.Builder
	prompt.WriteString(conversationContext.String())

	if hasRelevantContent && len(contextDocs) > 0 {
		prompt.WriteString("Based on the following relevant documentation:\n\n")
		prompt.WriteString(strings.Join(contextDocs, "\n\n"))
		prompt.WriteString("\n\n")
	}

	prompt.WriteString(fmt.Sprintf("Current user question: %s\n\n", message))

	if hasRelevantContent {
		prompt.WriteString("Please generate a short, focused tutorial in Markdown format based on the documentation above. Consider the conversation history for context. Structure your response as follows:\n\n# Quick Tutorial: [Topic Name]\n\n## What is [Topic]?\n[Brief 1-2 sentence explanation]\n\n## Key Concepts:\n- [Concept 1]\n- [Concept 2]\n- [Concept 3]\n\n## Basic Example:\n```[language]\n[Provide a simple, practical example]\n```\n\n## Common Use Cases:\n- [Use case 1]\n- [Use case 2]\n\n## Tips:\n- [Tip 1]\n- [Tip 2]\n\nKeep the tutorial concise, practical, and beginner-friendly. Use proper Markdown formatting.")
	} else {
		prompt.WriteString("Please provide a helpful and informative tutorial in Markdown format about this topic. Consider the conversation history for context. If you don't have specific information, provide general guidance and suggest where they might find more detailed information.\n\nStructure your response as a well-formatted Markdown tutorial with proper headers, code blocks, and bullet points.")
	}

	response, err := s.embClient.Chat(prompt.String())
	if err != nil {
		return "", fmt.Errorf("failed to generate response: %w", err)
	}

	// Store user message in history
	if err := s.AddChatMessage(sessionID, "user", message); err != nil {
		log.Printf("Failed to store user message: %v", err)
	}

	// Store AI response in history
	if err := s.AddChatMessage(sessionID, "assistant", response); err != nil {
		log.Printf("Failed to store AI response: %v", err)
	}

	// Store the response in vector database for learning
	go s.storeResponseForLearning(message, response, hasRelevantContent)

	return response, nil
}

// GetConversationInsights analyzes conversation history for insights
func (s *Service) GetConversationInsights(sessionID string) (map[string]interface{}, error) {
	history, err := s.GetChatHistory(sessionID, 50) // Get more history for analysis
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation history: %w", err)
	}

	insights := map[string]interface{}{
		"total_messages": len(history),
		"topics":         []string{},
		"user_questions": []string{},
		"ai_responses":   []string{},
	}

	var userMessages, aiMessages []string
	for _, msg := range history {
		if msg.Role == "user" {
			userMessages = append(userMessages, msg.Content)
		} else if msg.Role == "assistant" {
			aiMessages = append(aiMessages, msg.Content)
		}
	}

	insights["user_questions"] = userMessages
	insights["ai_responses"] = aiMessages

	// Analyze topics (simple keyword extraction)
	allText := strings.Join(userMessages, " ")
	words := strings.Fields(strings.ToLower(allText))
	wordCount := make(map[string]int)

	for _, word := range words {
		if len(word) > 3 { // Only count words longer than 3 characters
			wordCount[word]++
		}
	}

	// Get top 5 most common words as topics
	var topics []string
	for word, count := range wordCount {
		if count > 1 && len(topics) < 5 {
			topics = append(topics, word)
		}
	}
	insights["topics"] = topics

	return insights, nil
}
