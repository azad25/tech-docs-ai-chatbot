// internal/emb/fake.go
package emb

import (
	"errors"
)

// FakeClient is a fake implementation of an embedding client for testing.
type FakeClient struct{}

// NewFakeClient creates and returns a new FakeClient.
func NewFakeClient() *FakeClient {
	return &FakeClient{}
}

// Embed takes a text and returns a fake vector embedding.
func (c *FakeClient) Embed(text string) ([]float32, error) {
	// Return a static, fake embedding for demonstration.
	return []float32{0.1, 0.2, 0.3}, nil
}

// Chat returns a hardcoded chat response.
func (c *FakeClient) Chat(message string) (string, error) {
	if message == "" {
		return "", errors.New("empty message")
	}
	return "Fake LLM response for: " + message, nil
}
