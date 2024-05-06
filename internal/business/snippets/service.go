package snippets

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/titusjaka/go-sample/v2/internal/infrastructure/service"
)

//go:generate go run go.uber.org/mock/mockgen -typed -source=service.go -destination ./service_mock_test.go -package snippets_test -mock_names Storage=MockStorage

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
	logger  *slog.Logger

	now func() time.Time
}

// NewService returns new instance of SnippetService
func NewService(
	storage Storage,
	logger *slog.Logger,
	nowFunc func() time.Time,
) *SnippetService {
	return &SnippetService{
		storage: storage,
		logger:  logger,

		now: nowFunc,
	}
}

// Get returns a single snippet
func (s *SnippetService) Get(ctx context.Context, id uint) (Snippet, *service.Error) {
	snippet, err := s.storage.Get(ctx, id)
	switch {
	case err == nil:
		return snippet, nil
	case errors.Is(err, ErrNotFound):
		return Snippet{}, &service.Error{
			Type: service.NotFound,
			Base: ErrNotFound,
		}
	default:
		s.logger.Error("failed to get a snippet", slog.Any("err", err))
		return Snippet{}, &service.Error{
			Type: service.InternalError,
			Base: fmt.Errorf("failed to list snippets: %w", err),
		}
	}
}

// Create creates a single snippet
func (s *SnippetService) Create(ctx context.Context, snippet Snippet) (Snippet, *service.Error) {
	createdAt := s.now()
	snippet.CreatedAt = createdAt
	snippet.UpdatedAt = createdAt
	snippet.ExpiresAt = snippet.ExpiresAt.UTC()

	id, err := s.storage.Create(ctx, snippet)
	if err != nil {
		s.logger.Error("failed to create snippet", slog.Any("err", err))
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
		s.logger.Error("failed to query total amount of snippets", slog.Any("err", err))
		return nil, service.Pagination{}, &service.Error{
			Type: service.InternalError,
			Base: fmt.Errorf("failed to query total amount of snippets: %w", err),
		}
	}

	pagination := NewPagination(limit, offset, snippetsCount)

	snippets, err := s.storage.List(ctx, pagination)
	if err != nil {
		s.logger.Error("failed to list snippets", slog.Any("err", err))
		return nil, pagination, &service.Error{
			Type: service.InternalError,
			Base: fmt.Errorf("failed to list snippets: %w", err),
		}
	}

	return snippets, pagination, nil
}

// SoftDelete mark a single snippet as deleted
func (s *SnippetService) SoftDelete(ctx context.Context, id uint) *service.Error {
	switch err := s.storage.SoftDelete(ctx, id); {
	case err == nil:
		return nil
	case errors.Is(err, ErrNotFound):
		return &service.Error{
			Type: service.NotFound,
			Base: ErrNotFound,
		}
	default:
		s.logger.Error(
			"failed to soft delete snippet",
			slog.Uint64("id", uint64(id)),
			slog.Any("err", err),
		)
		return &service.Error{
			Type: service.InternalError,
			Base: fmt.Errorf("failed to delete snippet: %w", err),
		}
	}
}
