package snippets

import (
	"time"

	"github.com/titusjaka/go-sample/internal/infrastructure/service"
)

// SnippetResponse represents a common snippet-response struct
type SnippetResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// ListSnippetsResponse represents a response struct for GET /snippets?limit=<x>&offset=<y> method
type ListSnippetsResponse struct {
	Snippets   []SnippetResponse  `json:"snippets,omitempty"`
	Pagination service.Pagination `json:"pagination"`
}

// convertToListSnippetsResponse is used to map []Snippet -> []SnippetResponse
func convertToListSnippetsResponse(snippets []Snippet) []SnippetResponse {
	response := make([]SnippetResponse, len(snippets))
	for i := range snippets {
		response[i] = convertToSnippetResponse(snippets[i])
	}
	return response
}

// convertToSnippetResponse is used to map Snippet -> SnippetResponse
func convertToSnippetResponse(snippet Snippet) SnippetResponse {
	// nolint:gocritic
	return SnippetResponse{
		ID:        snippet.ID,
		Title:     snippet.Title,
		Content:   snippet.Content,
		CreatedAt: snippet.CreatedAt,
		ExpiresAt: snippet.ExpiresAt,
	}
}
