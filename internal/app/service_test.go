package app

import (
	"testing"

	"tech-docs-ai/internal/emb"
	"tech-docs-ai/internal/types"
	"tech-docs-ai/internal/vec"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEmbeddingClient mocks the embedding client
type MockEmbeddingClient struct {
	mock.Mock
}

func (m *MockEmbeddingClient) Embed(text string) ([]float32, error) {
	args := m.Called(text)
	return args.Get(0).([]float32), args.Error(1)
}

// MockVectorClient mocks the vector client
type MockVectorClient struct {
	mock.Mock
}

func (m *MockVectorClient) AddDocument(doc *types.Document, embedding []float32) error {
	args := m.Called(doc, embedding)
	return args.Error(0)
}

func (m *MockVectorClient) SearchDocuments(embedding []float32, limit int) ([]*types.Document, error) {
	args := m.Called(embedding, limit)
	return args.Get(0).([]*types.Document), args.Error(1)
}

func TestService_Chat(t *testing.T) {
	tests := []struct {
		name           string
		message        string
		expectedResp   string
		embeddingError error
		searchError    error
		searchResults  []*types.Document
		expectedError  bool
	}{
		{
			name:           "successful chat",
			message:        "Hello",
			expectedResp:   "Response based on: Test Doc",
			embeddingError: nil,
			searchError:    nil,
			searchResults: []*types.Document{
				{Title: "Test Doc", Content: "Test content"},
			},
			expectedError: false,
		},
		{
			name:           "embedding error",
			message:        "Hello",
			expectedResp:   "",
			embeddingError: fmt.Errorf("embedding error"),
			searchError:    nil,
			searchResults:  nil,
			expectedError:  true,
		},
		{
			name:           "search error",
			message:        "Hello",
			expectedResp:   "",
			embeddingError: nil,
			searchError:    fmt.Errorf("search error"),
			searchResults:  nil,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t := tt
		t.Run(tt.name, func(t *testing.T) {
			mockEmb := new(MockEmbeddingClient)
			mockVec := new(MockVectorClient)

			embedding := []float32{0.1, 0.2, 0.3}
			mockEmb.On("Embed", tt.message).Return(embedding, tt.embeddingError)
			if tt.embeddingError == nil {
				mockVec.On("SearchDocuments", embedding, 5).Return(tt.searchResults, tt.searchError)
			}

			svc := NewService(mockEmb, mockVec)
			resp, err := svc.Chat(tt.message)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Empty(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResp, resp)
			}

			mockEmb.AssertExpectations(t)
			mockVec.AssertExpectations(t)
		})
	}
}

func TestService_AddDocument(t *testing.T) {
	tests := []struct {
		name           string
		doc            *types.Document
		embeddingError error
		addError       error
		expectedError  bool
	}{
		{
			name: "successful add",
			doc: &types.Document{
				Title:   "Test Doc",
				Content: "Test content",
			},
			embeddingError: nil,
			addError:       nil,
			expectedError:  false,
		},
		{
			name: "embedding error",
			doc: &types.Document{
				Title:   "Test Doc",
				Content: "Test content",
			},
			embeddingError: fmt.Errorf("embedding error"),
			addError:       nil,
			expectedError:  true,
		},
		{
			name: "add error",
			doc: &types.Document{
				Title:   "Test Doc",
				Content: "Test content",
			},
			embeddingError: nil,
			addError:       fmt.Errorf("add error"),
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t := tt
		t.Run(tt.name, func(t *testing.T) {
			mockEmb := new(MockEmbeddingClient)
			mockVec := new(MockVectorClient)

			embedding := []float32{0.1, 0.2, 0.3}
			mockEmb.On("Embed", tt.doc.Content).Return(embedding, tt.embeddingError)
			if tt.embeddingError == nil {
				mockVec.On("AddDocument", tt.doc, embedding).Return(tt.addError)
			}

			svc := NewService(mockEmb, mockVec)
			err := svc.AddDocument(tt.doc)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockEmb.AssertExpectations(t)
			mockVec.AssertExpectations(t)
		})
	}
}

func TestService_SearchDocuments(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		limit          int
		embeddingError error
		searchError    error
		searchResults  []*types.Document
		expectedError  bool
	}{
		{
			name:           "successful search",
			query:          "test",
			limit:          5,
			embeddingError: nil,
			searchError:    nil,
			searchResults: []*types.Document{
				{Title: "Test Doc", Content: "Test content"},
			},
			expectedError: false,
		},
		{
			name:           "embedding error",
			query:          "test",
			limit:          5,
			embeddingError: fmt.Errorf("embedding error"),
			searchError:    nil,
			searchResults:  nil,
			expectedError:  true,
		},
		{
			name:           "search error",
			query:          "test",
			limit:          5,
			embeddingError: nil,
			searchError:    fmt.Errorf("search error"),
			searchResults:  nil,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t := tt
		t.Run(tt.name, func(t *testing.T) {
			mockEmb := new(MockEmbeddingClient)
			mockVec := new(MockVectorClient)

			embedding := []float32{0.1, 0.2, 0.3}
			mockEmb.On("Embed", tt.query).Return(embedding, tt.embeddingError)
			if tt.embeddingError == nil {
				mockVec.On("SearchDocuments", embedding, tt.limit).Return(tt.searchResults, tt.searchError)
			}

			svc := NewService(mockEmb, mockVec)
			docs, err := svc.SearchDocuments(tt.query, tt.limit)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, docs)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.searchResults, docs)
			}

			mockEmb.AssertExpectations(t)
			mockVec.AssertExpectations(t)
		})
	}
}
