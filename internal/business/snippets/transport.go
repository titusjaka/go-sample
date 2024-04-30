package snippets

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/gorilla/schema"

	"github.com/titusjaka/go-sample/internal/infrastructure/api"
	"github.com/titusjaka/go-sample/internal/infrastructure/service"
)

//go:generate go run go.uber.org/mock/mockgen -typed -source=transport.go -destination ./transport_mock_test.go -package snippets_test -mock_names Service=MockService

// Service is used to manipulate data over snippets
type Service interface {
	Get(ctx context.Context, id uint) (Snippet, *service.Error)
	Create(ctx context.Context, snippet Snippet) (Snippet, *service.Error)
	List(ctx context.Context, limit uint, offset uint) ([]Snippet, service.Pagination, *service.Error)
	SoftDelete(ctx context.Context, id uint) *service.Error
}

// Transport is a struct that holds all endpoints for snippets
type Transport struct {
	logger  *slog.Logger
	service Service
}

// NewTransport creates a new Transport instance
func NewTransport(s Service, l *slog.Logger) *Transport {
	return &Transport{
		logger:  l,
		service: s,
	}
}

// Routes initialize all endpoints for route /snippets
func (t *Transport) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", t.listSnippets)
	r.Post("/", t.createSnippet)
	r.Route("/{snippet_id}", func(r chi.Router) {
		r.Get("/", t.getSnippet)
		r.Delete("/", t.deleteSnippet)
	})

	return r
}

// listSnippets in an endpoint for GET /snippets method
func (t *Transport) listSnippets(w http.ResponseWriter, r *http.Request) {
	var listSnippetsRequest ListSnippetsRequest
	if err := schema.NewDecoder().Decode(&listSnippetsRequest, r.URL.Query()); err != nil {
		t.logger.Error("failed to decode request params", slog.Any("err", err))
		_ = render.Render(w, r, api.ErrBadRequest(err))
		return
	}

	snippets, pagination, svcErr := t.service.List(r.Context(), listSnippetsRequest.Limit, listSnippetsRequest.Offset)
	if svcErr != nil {
		t.logger.Error("failed to list snippets", slog.Any("svc_err", svcErr))
		_ = render.Render(w, r, api.NewErrResponse(svcErr))
		return
	}

	render.JSON(w, r, &ListSnippetsResponse{
		Snippets:   convertToListSnippetsResponse(snippets),
		Pagination: pagination,
	})
}

// getSnippet is an endpoint for GET /snippets/{snippet_id} method
func (t *Transport) getSnippet(w http.ResponseWriter, r *http.Request) {
	snippetID, svcErr := parseSnippetID(r)
	if svcErr != nil {
		t.logger.Error("failed to parse snippet id", slog.Any("svc_err", svcErr))
		_ = render.Render(w, r, api.NewErrResponse(svcErr))
		return
	}

	snippet, svcErr := t.service.Get(r.Context(), snippetID)
	if svcErr != nil {
		t.logger.Error("failed to get snippet", slog.Any("svc_err", svcErr))
		_ = render.Render(w, r, api.NewErrResponse(svcErr))
		return
	}

	render.JSON(w, r, convertToSnippetResponse(snippet))
}

// createSnippet in an endpoint for POST /snippets method
func (t *Transport) createSnippet(w http.ResponseWriter, r *http.Request) {
	var createSnippetReq CreateSnippetRequest
	if err := render.Decode(r, &createSnippetReq); err != nil {
		t.logger.Error("failed to decode request params", slog.Any("err", err))
		_ = render.Render(w, r, api.ErrBadRequest(err))
		return
	}

	if validationErr := createSnippetReq.Validate(); validationErr != nil {
		t.logger.Info("request is not valid", slog.Any("validation_err", validationErr))
		_ = render.Render(w, r, api.ErrBadRequest(validationErr))
		return
	}

	newSnippet := Snippet{
		Title:     createSnippetReq.Title,
		Content:   createSnippetReq.Content,
		ExpiresAt: createSnippetReq.ExpiresAt,
	}

	snippet, svcErr := t.service.Create(r.Context(), newSnippet)
	if svcErr != nil {
		t.logger.Error("failed to create snippet", slog.Any("svc_err", svcErr))
		_ = render.Render(w, r, api.NewErrResponse(svcErr))
		return
	}

	render.JSON(w, r, convertToSnippetResponse(snippet))
}

// deleteSnippet in an endpoint for DELETE /snippets/{snippet_id} method
func (t *Transport) deleteSnippet(w http.ResponseWriter, r *http.Request) {
	snippetID, svcErr := parseSnippetID(r)
	if svcErr != nil {
		t.logger.Error("failed to parse snippet id", slog.Any("svc_err", svcErr))
		_ = render.Render(w, r, api.NewErrResponse(svcErr))
		return
	}

	svcErr = t.service.SoftDelete(r.Context(), snippetID)
	if svcErr != nil {
		t.logger.Error("failed to delete snippet", slog.Any("svc_err", svcErr))
		_ = render.Render(w, r, api.NewErrResponse(svcErr))
		return
	}

	render.NoContent(w, r)
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
