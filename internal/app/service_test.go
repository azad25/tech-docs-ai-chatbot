package app

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"tech-docs-ai/internal/emb" // Import emb package
	"tech-docs-ai/internal/vec" // Import vec package
)

func TestServiceChat(t *testing.T) {
	embClient := emb.NewFakeClient()
	vecClient := vec.NewQdrantClient()      // Add vecClient
	svc := NewService(embClient, vecClient) // Pass both arguments

	response, err := svc.Chat("Hello")

	assert.NoError(t, err)
	assert.Equal(t, "Fake LLM response for: Hello", response)
}
