package pgmigrator

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"text/template"
	"time"

	migrate "github.com/rubenv/sql-migrate"
)

const migrationsTableName = "migrations"

const templateContent = `-- +migrate Up

-- +migrate Down
`

// Migrator represents a migration service
type Migrator struct {
	db     *sql.DB
	source migrate.MigrationSource
}

// NewMigrator returns a new Migrator
func NewMigrator(db *sql.DB, source fs.FS) *Migrator {
	return &Migrator{
		db: db,
		source: &migrate.HttpFileSystemMigrationSource{
			FileSystem: http.FS(source),
		},
	}
}

// Create a new migration file with a given name
func Create(migrationsDir, name string) (fullPath string, err error) {
	dir := filepath.Clean(migrationsDir)

	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", fmt.Errorf("create directory for migrations: %w", err)
	}

	fileName := fmt.Sprintf("%d_%s.sql", time.Now().Unix(), name)
	pathName := path.Join(dir, fileName)

	file, err := os.Create(path.Clean(pathName))
	if err != nil {
		return "", fmt.Errorf("create migration file (%q): %w", pathName, err)
	}

	tpl, err := template.New("new_migration").Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("parse migration template: %w", err)
	}

	if err = tpl.Execute(file, nil); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	return pathName, nil
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
		return 0, fmt.Errorf("begin db transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	migrate.SetTable(migrationsTableName)

	// Ensure we got single migrator running at a time
	// Other concurrent sessions should wait
	if _, err = tx.ExecContext(ctx, `SELECT PG_ADVISORY_XACT_LOCK(1)`); err != nil {
		return 0, fmt.Errorf("acquire advisory lock: %w", err)
	}

	applied, err := migrate.ExecMax(m.db, "postgres", m.source, direction, max)
	if err != nil {
		return 0, fmt.Errorf("apply database migrations: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit db transaction: %w", err)
	}

	return applied, nil
}
