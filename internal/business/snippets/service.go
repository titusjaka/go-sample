package snippets

import (
	"context"
	"fmt"
	"time"

	"github.com/titusjaka/go-sample/internal/infrastructure/log"
	"github.com/titusjaka/go-sample/internal/infrastructure/service"
)

//go:generate mockgen -source=service.go -destination ./mocks_storage_test.go -package snippets_test -mock_names Storage=MockStorage

// Storage is used to manipulate data in DB
type Storage interface {
	Get(ctx context.Context, id uint) (Snippet, error)
	Create(ctx context.Context, snippet Snippet) (uint, error)
	List(ctx context.Context, pagination service.Pagination) ([]Snippet, error)
	SoftDelete(ctx context.Context, id uint) error
	Total(ctx context.Context) (uint, error)
}

// SnippetService represents service struct. It holds storage and logger.
type SnippetService struct {
	storage Storage
	logger  log.Logger
}

// NewService returns new instance of SnippetService
func NewService(storage Storage, logger log.Logger) *SnippetService {
	return &SnippetService{
		storage: storage,
		logger:  logger,
	}
}

// Get returns a single snippet
func (s *SnippetService) Get(ctx context.Context, id uint) (Snippet, *service.Error) {
	snippet, err := s.storage.Get(ctx, id)
	switch err {
	case nil:
		return snippet, nil
	case ErrNotFound:
		return Snippet{}, &service.Error{
			Type: service.NotFound,
			Base: ErrNotFound,
		}
	default:
		s.logger.Error("failed to get a snippet", log.Field("err", err))
		return Snippet{}, &service.Error{
			Type: service.InternalError,
			Base: fmt.Errorf("failed to list snippets: %w", err),
		}
	}
}

// Create creates a single snippet
func (s *SnippetService) Create(ctx context.Context, snippet Snippet) (Snippet, *service.Error) {
	createdAt := time.Now().UTC()
	snippet.CreatedAt = createdAt
	snippet.UpdatedAt = createdAt
	snippet.ExpiresAt = snippet.ExpiresAt.UTC()

	id, err := s.storage.Create(ctx, snippet)
	if err != nil {
		s.logger.Error("failed to create snippet", log.Field("err", err))
		return Snippet{}, &service.Error{
			Type: service.InternalError,
			Base: fmt.Errorf("failed to create snippet: %w", err),
		}
	}

	snippet.ID = id
	return snippet, nil
}

// List returns a list of snippets and a pagination struct
func (s *SnippetService) List(ctx context.Context, limit uint, offset uint) ([]Snippet, service.Pagination, *service.Error) {
	snippetsCount, err := s.storage.Total(ctx)
	if err != nil {
		s.logger.Error("failed to query total amount of snippets", log.Field("err", err))
		return nil, service.Pagination{}, &service.Error{
			Type: service.InternalError,
			Base: fmt.Errorf("failed to query total amount of snippets: %w", err),
		}
	}

	pagination := NewPagination(limit, offset, snippetsCount)

	snippets, err := s.storage.List(ctx, pagination)
	if err != nil {
		s.logger.Error("failed to list snippets", log.Field("err", err))
		return nil, pagination, &service.Error{
			Type: service.InternalError,
			Base: fmt.Errorf("failed to list snippets: %w", err),
		}
	}

	return snippets, pagination, nil
}

// SoftDelete mark a single snippet as deleted
func (s *SnippetService) SoftDelete(ctx context.Context, id uint) *service.Error {
	switch err := s.storage.SoftDelete(ctx, id); err {
	case nil:
		return nil
	case ErrNotFound:
		return &service.Error{
			Type: service.NotFound,
			Base: ErrNotFound,
		}
	default:
		s.logger.Error(
			"failed to soft delete snippet",
			log.Field("id", id),
			log.Field("err", err),
		)
		return &service.Error{
			Type: service.InternalError,
			Base: fmt.Errorf("failed to delete snippet: %w", err),
		}
	}
}
