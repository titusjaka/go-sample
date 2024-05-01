package snippets

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/titusjaka/go-sample/internal/infrastructure/service"
)

// ErrNotFound error used to signal higher level about sql.ErrNoRows error
var ErrNotFound = errors.New("not found")

// PGStorage implements storage interface and provides methods to manipulate data in PostgreSQL storage
type PGStorage struct {
	conn *sql.DB
}

// NewPGStorage returns a new instance of PGStorage
func NewPGStorage(conn *sql.DB) *PGStorage {
	return &PGStorage{conn: conn}
}

// Get returns a single snippet from storage
func (pg *PGStorage) Get(ctx context.Context, id uint) (Snippet, error) {
	query := `
		SELECT 
			id, 
			title,
			content,
			created_at,
			updated_at,
			expires_at
		FROM 
			snippets
		WHERE id = $1
	`

	var snippet Snippet
	switch err := pg.conn.QueryRowContext(ctx, query, id).Scan(
		&snippet.ID,
		&snippet.Title,
		&snippet.Content,
		&snippet.CreatedAt,
		&snippet.UpdatedAt,
		&snippet.ExpiresAt,
	); {
	case err == nil:
		return snippet, nil
	case errors.Is(err, sql.ErrNoRows):
		return Snippet{}, ErrNotFound
	default:
		return Snippet{}, fmt.Errorf("failed to scan snippet: %w", err)
	}
}

// Create saves a single snippet to storage
func (pg *PGStorage) Create(ctx context.Context, snippet Snippet) (uint, error) {
	query := `
		INSERT INTO snippets
		(
			title,
			content,
			created_at,
			updated_at,
			expires_at
		)
		VALUES
		(
			$1,
			$2,
			$3,
			$4,
			$5
		)
		RETURNING id
	`

	var id uint
	switch err := pg.conn.QueryRowContext(
		ctx,
		query,
		snippet.Title,
		snippet.Content,
		snippet.CreatedAt,
		snippet.UpdatedAt,
		snippet.ExpiresAt,
	).Scan(&id); err {
	case nil:
		return id, nil
	default:
		return 0, fmt.Errorf("failed to add snippet: %w", err)
	}
}

// List returns a list of snippets from storage
func (pg *PGStorage) List(ctx context.Context, pagination service.Pagination) ([]Snippet, error) {
	query := `
		SELECT
			id,
			title,
			content,
			created_at,
			updated_at,
			expires_at
		FROM snippets
		WHERE
			expires_at > NOW()
		ORDER BY created_at DESC
		%s
	`

	paginationExpression := ConvertPaginationToSQLExpression(pagination)

	rows, err := pg.conn.QueryContext(
		ctx,
		fmt.Sprintf(query, paginationExpression),
	)
	switch {
	case err == nil:
		break
	default:
		return nil, fmt.Errorf("failed to list snippets: %w", err)
	}

	defer func() {
		_ = rows.Close()
	}()

	var results []Snippet
	for rows.Next() {
		var snippet Snippet
		err := rows.Scan(
			&snippet.ID,
			&snippet.Title,
			&snippet.Content,
			&snippet.CreatedAt,
			&snippet.UpdatedAt,
			&snippet.ExpiresAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan snippet row: %w", err)
		}

		results = append(results, snippet)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error from iterating snippets rows: %w", err)
	}

	return results, nil
}

// SoftDelete set `expires_at` to `now()`, so snippet is considered deleted
func (pg *PGStorage) SoftDelete(ctx context.Context, id uint) error {
	wrapErr := func(err error) error {
		return fmt.Errorf("failed to soft delete snippet from DB (ID: %d): %w", id, err)
	}

	query := `
		UPDATE snippets
		SET
			updated_at = NOW(),
			expires_at = NOW()
		WHERE
			id = $1
	`

	result, err := pg.conn.ExecContext(ctx, query, id)
	if err != nil {
		return wrapErr(err)
	}

	affectedRows, err := result.RowsAffected()
	switch {
	case err != nil:
		return wrapErr(err)
	case affectedRows == 0:
		return ErrNotFound
	default:
		return nil
	}
}

// Total counts a total number of snippets
func (pg *PGStorage) Total(ctx context.Context) (uint, error) {
	query := `
		SELECT COUNT(*) 
		FROM snippets
	`

	row := pg.conn.QueryRowContext(ctx, query)

	var count uint
	err := row.Scan(&count)
	return count, err
}
