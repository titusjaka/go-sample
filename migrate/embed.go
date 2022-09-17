package migrate

import (
	"embed"
	"io/fs"
	"net/http"

	sqlmigrate "github.com/rubenv/sql-migrate"
)

const dirName = "migrations"

// migrationsDir is used to hold embedded migration files.
// It must be global var due to embed package convention.
//
//go:embed migrations/*.sql
var migrationsDir embed.FS //nolint:gochecknoglobals

// NewEmbeddedSource returns a new sqlmigrate.MigrationSource
func NewEmbeddedSource() (*sqlmigrate.HttpFileSystemMigrationSource, error) {
	subDir, err := fs.Sub(migrationsDir, dirName)
	if err != nil {
		return nil, err
	}

	return &sqlmigrate.HttpFileSystemMigrationSource{FileSystem: http.FS(subDir)}, nil
}
