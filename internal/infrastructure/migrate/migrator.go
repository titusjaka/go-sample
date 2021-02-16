package migrate

import (
	"context"
	"database/sql"
	"fmt"

	migrate "github.com/rubenv/sql-migrate"
)

const migrationsTable = "migrations"

// Migrator represents a migration service
type Migrator struct {
	db     *sql.DB
	source migrate.MigrationSource
}

// NewMigrator returns a new Migrator
func NewMigrator(db *sql.DB, source migrate.MigrationSource) *Migrator {
	return &Migrator{
		db:     db,
		source: source,
	}
}

// Up applies all migrations
func (m *Migrator) Up(ctx context.Context) (int, error) {
	return m.apply(ctx, migrate.Up, 0)
}

// Down rollback a number of migrations
func (m *Migrator) Down(ctx context.Context, max int) (int, error) {
	return m.apply(ctx, migrate.Down, max)
}

func (m *Migrator) apply(ctx context.Context, direction migrate.MigrationDirection, max int) (int, error) {
	tx, err := m.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("can't begin db transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	migrate.SetTable(migrationsTable)

	// Ensure we got single migrator running at a time
	// Other concurrent sessions should wait
	if _, err = tx.ExecContext(ctx, `SELECT pg_advisory_xact_lock(1)`); err != nil {
		return 0, fmt.Errorf("can't acquire advisory lock: %w", err)
	}

	applied, err := migrate.ExecMax(m.db, "postgres", m.source, direction, max)
	if err != nil {
		return 0, fmt.Errorf("can't apply database migrations: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("can't commit db transaction: %w", err)
	}

	return applied, nil
}
