package snippets

import (
	"net/http"
	"time"

	"github.com/titusjaka/go-sample/internal/infrastructure/service"
)

// snippetResponse represents a common snippet-response struct
type snippetResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// listSnippetsResponse represents a response struct for GET /snippets?limit=<x>&offset=<y> method
type listSnippetsResponse struct {
	Snippets   []snippetResponse  `json:"snippets,omitempty"`
	Pagination service.Pagination `json:"pagination"`
	err        *service.Error
}

// ServiceError is an implementation of api.ServiceErrorer interface.
// It's used to handle business-logic errors as HTTP responses
func (r *listSnippetsResponse) ServiceError() *service.Error {
	return r.err
}

// getSnippetResponse represents a response struct for GET /snippets/{snippet_id} method
type getSnippetResponse struct {
	snippetResponse
	err *service.Error
}

// ServiceError is an implementation of api.ServiceErrorer interface.
// It's used to handle business-logic errors as HTTP responses
func (r *getSnippetResponse) ServiceError() *service.Error {
	return r.err
}

// createSnippetResponse represents a response struct for POST /snippets method
type createSnippetResponse struct {
	snippetResponse
	err *service.Error
}

// StatusCode is an implementation of go-kit StatusCoder interface.
// It's used to override default 200 OK response
func (r *createSnippetResponse) StatusCode() int {
	return http.StatusCreated
}

// ServiceError is an implementation of api.ServiceErrorer interface.
// It's used to handle business-logic errors as HTTP responses
func (r *createSnippetResponse) ServiceError() *service.Error {
	return r.err
}

// deleteSnippetResponse represents a response struct for DELETE /snippets/{snippet_id} method
type deleteSnippetResponse struct {
	err *service.Error
}

// StatusCode is an implementation of go-kit StatusCoder interface.
// It's used to override default 200 OK response
func (r *deleteSnippetResponse) StatusCode() int {
	return http.StatusNoContent
}

// ServiceError is an implementation of api.ServiceErrorer interface.
// It's used to handle business-logic errors as HTTP responses
func (r *deleteSnippetResponse) ServiceError() *service.Error {
	return r.err
}

// convertToListSnippetsResponse is used to map []Snippet -> []snippetResponse
func convertToListSnippetsResponse(snippets []Snippet) []snippetResponse {
	response := make([]snippetResponse, len(snippets))
	for i := range snippets {
		response[i] = convertToSnippetResponse(snippets[i])
	}
	return response
}

// convertToSnippetResponse is used to map Snippet -> snippetResponse
func convertToSnippetResponse(snippet Snippet) snippetResponse {
	// nolint:gocritic
	return snippetResponse{
		ID:        snippet.ID,
		Title:     snippet.Title,
		Content:   snippet.Content,
		CreatedAt: snippet.CreatedAt,
		ExpiresAt: snippet.ExpiresAt,
	}
}
