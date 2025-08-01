package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"tech-docs-ai/internal/emb" // Import emb package
	"tech-docs-ai/internal/vec" // Import vec package

	"github.com/stretchr/testify/assert"
)

func TestHandleChat(t *testing.T) {
	embClient := emb.NewFakeClient()
	vecClient := vec.NewQdrantClient()      // Add vecClient
	svc := NewService(embClient, vecClient) // Pass both arguments
	handler := NewHandler(svc)

	reqBody := chatRequest{Message: "Hello"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/chat", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.HandleChat(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var chatResp chatResponse
	json.NewDecoder(resp.Body).Decode(&chatResp)
	assert.Equal(t, "Fake LLM response for: Hello", chatResp.Response)
}
