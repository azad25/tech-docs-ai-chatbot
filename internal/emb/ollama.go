package emb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// OllamaClient is a real implementation of an embedding client that communicates with Ollama.
type OllamaClient struct {
	apiURL     string
	embedModel string
	chatModel  string
	httpClient *http.Client
}

// NewOllamaClient creates and returns a new OllamaClient.
func NewOllamaClient() *OllamaClient {
	apiURL := os.Getenv("OLLAMA_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:11434"
	}

	embedModel := os.Getenv("OLLAMA_MODEL")
	if embedModel == "" {
		embedModel = "nomic-embed-text"
	}

	chatModel := os.Getenv("OLLAMA_CHAT_MODEL")
	if chatModel == "" {
		chatModel = "tinyllama"
	}

	return &OllamaClient{
		apiURL:     apiURL,
		embedModel: embedModel,
		chatModel:  chatModel,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// EmbedRequest represents the request structure for Ollama embeddings.
type EmbedRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// EmbedResponse represents the response structure for Ollama embeddings.
type EmbedResponse struct {
	Embedding []float32 `json:"embedding"`
}

// ChatRequest represents the request structure for Ollama chat.
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

// Message represents a chat message.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse represents the response structure for Ollama chat.
type ChatResponse struct {
	Model     string  `json:"model"`
	Message   Message `json:"message"`
	Done      bool    `json:"done"`
	CreatedAt string  `json:"created_at"`
}

// Embed takes a text and returns its vector representation using Ollama.
func (c *OllamaClient) Embed(text string) ([]float32, error) {
	reqBody := EmbedRequest{
		Model:  c.embedModel,
		Prompt: text,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal embed request: %w", err)
	}

	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/api/embeddings", c.apiURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to make embed request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("embed request failed with status: %d", resp.StatusCode)
	}

	var embedResp EmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&embedResp); err != nil {
		return nil, fmt.Errorf("failed to decode embed response: %w", err)
	}

	return embedResp.Embedding, nil
}

// Chat takes a message and returns a response from the LLM using Ollama.
func (c *OllamaClient) Chat(message string) (string, error) {
	reqBody := ChatRequest{
		Model: c.chatModel,
		Messages: []Message{
			{
				Role:    "user",
				Content: message,
			},
		},
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal chat request: %w", err)
	}

	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/api/chat", c.apiURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("failed to make chat request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("chat request failed with status: %d", resp.StatusCode)
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("failed to decode chat response: %w", err)
	}

	return chatResp.Message.Content, nil
}
