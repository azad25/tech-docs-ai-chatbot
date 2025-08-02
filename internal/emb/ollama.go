package emb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// OllamaClient is a client for interacting with Ollama API
type OllamaClient struct {
	apiURL     string
	model      string
	chatModel  string
	httpClient *http.Client
}

// NewOllamaClient creates a new Ollama client
func NewOllamaClient() *OllamaClient {
	apiURL := os.Getenv("OLLAMA_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:11434"
	}

	model := os.Getenv("OLLAMA_MODEL")
	if model == "" {
		model = "nomic-embed-text"
	}

	chatModel := os.Getenv("OLLAMA_CHAT_MODEL")
	if chatModel == "" {
		chatModel = "llama3.2:1b"
	}

	return &OllamaClient{
		apiURL:    apiURL,
		model:     model,
		chatModel: chatModel,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// EmbedRequest represents a request to generate embeddings
type EmbedRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// EmbedResponse represents a response from the embedding API
type EmbedResponse struct {
	Embedding []float32 `json:"embedding"`
}

// ChatRequest represents a request to the chat API
type ChatRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// ChatResponse represents a response from the chat API
type ChatResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// Embed generates embeddings for the given text
func (c *OllamaClient) Embed(text string) ([]float32, error) {
	request := EmbedRequest{
		Model:  c.model,
		Prompt: text,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal embed request: %w", err)
	}

	resp, err := c.httpClient.Post(c.apiURL+"/api/embeddings", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to make embed request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("embed request failed with status: %d", resp.StatusCode)
	}

	var embedResponse EmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&embedResponse); err != nil {
		return nil, fmt.Errorf("failed to decode embed response: %w", err)
	}

	return embedResponse.Embedding, nil
}

// Chat generates a chat response for the given message
func (c *OllamaClient) Chat(message string) (string, error) {
	request := ChatRequest{
		Model:  c.chatModel,
		Prompt: message,
		Stream: false,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal chat request: %w", err)
	}

	resp, err := c.httpClient.Post(c.apiURL+"/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to make chat request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("chat request failed with status: %d", resp.StatusCode)
	}

	var chatResponse ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResponse); err != nil {
		return "", fmt.Errorf("failed to decode chat response: %w", err)
	}

	return chatResponse.Response, nil
}
