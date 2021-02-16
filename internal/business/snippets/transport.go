package snippets

import (
	"github.com/go-chi/chi"
	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/titusjaka/go-sample/internal/infrastructure/api"
	"github.com/titusjaka/go-sample/internal/infrastructure/log"
)

// MakeSnippetsHandler initialize all endpoints for route /snippets
func MakeSnippetsHandler(s Service, l log.Logger) chi.Router {
	options := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(api.EncodeError),
		kithttp.ServerErrorHandler(api.NewLogErrorHandler(l)),
	}

	listSnippetsHandler := kithttp.NewServer(
		MakeListSnippetsEndpoint(s),
		decodeListSnippetsRequest,
		api.EncodeResponse,
		options...,
	)

	createSnippetHandler := kithttp.NewServer(
		MakeCreateSnippetEndpoint(s),
		decodeCreateSnippetRequest,
		api.EncodeResponse,
		options...,
	)

	getSnippetHandler := kithttp.NewServer(
		MakeGetSnippetEndpoint(s),
		decodeGetSnippetRequest,
		api.EncodeResponse,
		options...,
	)

	deleteSnippetHandler := kithttp.NewServer(
		MakeDeleteSnippetEndpoint(s),
		decodeDeleteSnippetRequest,
		api.EncodeResponse,
		options...,
	)

	r := chi.NewRouter()
	r.Get("/", listSnippetsHandler.ServeHTTP)
	r.Post("/", createSnippetHandler.ServeHTTP)
	r.Route("/{snippet_id}", func(r chi.Router) {
		r.Get("/", getSnippetHandler.ServeHTTP)
		r.Delete("/", deleteSnippetHandler.ServeHTTP)
	})

	return r
}
