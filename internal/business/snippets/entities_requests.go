package snippets

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/schema"

	"github.com/titusjaka/go-sample/internal/infrastructure/service"
)

// getSnippetRequest represents a request struct for GET /snippets/{snippet_id} method
type getSnippetRequest struct {
	SnippetID uint
}

// decodeGetSnippetRequest implements go-kit DecodeRequestFunc for GET /snippets/{snippet_id} method
func decodeGetSnippetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id, svcErr := parseSnippetID(r)
	if svcErr != nil {
		return nil, svcErr
	}
	return &getSnippetRequest{
		SnippetID: id,
	}, nil
}

// ListSnippetsRequest represents a request struct for GET /snippets?limit=<x>&offset=<y> method
type ListSnippetsRequest struct {
	Limit  uint
	Offset uint
}

// decodeListSnippetsRequest implements go-kit DecodeRequestFunc for GET /snippets?limit=<x>&offset=<y> method
func decodeListSnippetsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req ListSnippetsRequest
	err := schema.NewDecoder().Decode(&req, r.URL.Query())
	if err != nil {
		return nil, &service.Error{
			Type: service.BadRequest,
			Base: fmt.Errorf("failed to decode request params: %w", err),
		}
	}
	return &req, err
}

// createSnippetRequest represents a request struct for POST /snippets method
type createSnippetRequest struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	ExpiresAt string `json:"expires_at"`
}

// Validate implements ozzo-validation.Validatable interface and used to check user request
func (r *createSnippetRequest) Validate() error {
	now := time.Now().UTC()
	yearAfter := now.Add(366 * 24 * time.Hour)
	rules := []*validation.FieldRules{
		validation.Field(&r.Title, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.Content, validation.Required, validation.Length(1, 10000)),
		validation.Field(&r.ExpiresAt, validation.Date(time.RFC3339).Min(now).Max(yearAfter).RangeError(
			fmt.Sprintf(
				"expires_at (%s) must be a valid RFC3339 date within range from now (%s) to now + 1 year (%s)",
				r.ExpiresAt,
				now.Format(time.RFC3339),
				yearAfter.Format(time.RFC3339),
			))),
	}

	return validation.ValidateStruct(r, rules...)
}

// decodeCreateSnippetRequest implements go-kit DecodeRequestFunc for POST /snippets method
func decodeCreateSnippetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req createSnippetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, &service.Error{
			Type: service.BadRequest,
			Base: err,
		}
	}
	return &req, nil
}

// deleteSnippetRequest represents a request struct for GET /snippets/{snippet_id} method
type deleteSnippetRequest struct {
	SnippetID uint
}

// decodeDeleteSnippetRequest implements go-kit DecodeRequestFunc for DELETE /snippets/{snippet_id} method
func decodeDeleteSnippetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id, svcErr := parseSnippetID(r)
	if svcErr != nil {
		return nil, svcErr
	}
	return &deleteSnippetRequest{SnippetID: id}, nil
}

// parseSnippetID fetches URLParam from go-chi request Context and check it. In case of error service.Error is returned
func parseSnippetID(r *http.Request) (uint, *service.Error) {
	id, err := strconv.Atoi(chi.URLParam(r, "snippet_id"))
	switch {
	case err != nil:
		return 0, &service.Error{
			Type: service.BadRequest,
			Base: err,
		}
	case id <= 0:
		return 0, &service.Error{
			Type: service.BadRequest,
			Base: fmt.Errorf("invalid id param: %d", id),
		}
	default:
		return uint(id), nil
	}
}
