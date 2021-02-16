package snippets

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"

	"github.com/titusjaka/go-sample/internal/infrastructure/service"
)

//go:generate mockgen -source=endpoints.go -destination ./mocks_service_test.go -package snippets_test -mock_names Service=MockService

// Service is used to manipulate data over snippets
type Service interface {
	Get(ctx context.Context, id uint) (Snippet, *service.Error)
	Create(ctx context.Context, snippet Snippet) (Snippet, *service.Error)
	List(ctx context.Context, limit uint, offset uint) ([]Snippet, service.Pagination, *service.Error)
	SoftDelete(ctx context.Context, id uint) *service.Error
}

// MakeListSnippetsEndpoint creates a new go-kit Endpoint for GET /snippets method
func MakeListSnippetsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// nolint:errcheck
		req := request.(*ListSnippetsRequest)
		snippets, pagination, svcErr := s.List(ctx, req.Limit, req.Offset)

		return &listSnippetsResponse{
			Snippets:   convertToListSnippetsResponse(snippets),
			Pagination: pagination,
			err:        svcErr,
		}, nil
	}
}

// MakeGetSnippetEndpoint creates a new go-kit Endpoint for GET /snippets/{snippet_id} method
func MakeGetSnippetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// nolint:errcheck
		req := request.(*getSnippetRequest)
		snippet, svcErr := s.Get(ctx, req.SnippetID)
		return &getSnippetResponse{
			snippetResponse: convertToSnippetResponse(snippet),
			err:             svcErr,
		}, nil
	}
}

// MakeCreateSnippetEndpoint creates a new go-kit Endpoint for POST /snippets method
func MakeCreateSnippetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		wrapErr := func(err error) (interface{}, error) {
			return &createSnippetResponse{
				err: &service.Error{
					Type: service.BadRequest,
					Base: err,
				},
			}, nil
		}

		// nolint:errcheck
		req := request.(*createSnippetRequest)
		if err := req.Validate(); err != nil {
			return wrapErr(err)
		}

		expiresAt, err := time.Parse(time.RFC3339, req.ExpiresAt)
		if err != nil {
			return wrapErr(err)
		}
		newSnippet := Snippet{
			Title:     req.Title,
			Content:   req.Content,
			ExpiresAt: expiresAt,
		}

		snippet, svcErr := s.Create(ctx, newSnippet)

		return &createSnippetResponse{
			snippetResponse: convertToSnippetResponse(snippet),
			err:             svcErr,
		}, nil
	}
}

// MakeDeleteSnippetEndpoint creates a new go-kit Endpoint for DELETE /snippets/{snippet_id} method
func MakeDeleteSnippetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// nolint:errcheck
		req := request.(*deleteSnippetRequest)
		svcErr := s.SoftDelete(ctx, req.SnippetID)
		return &deleteSnippetResponse{
			err: svcErr,
		}, nil
	}
}
