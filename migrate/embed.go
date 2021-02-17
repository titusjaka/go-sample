package migrate

import (
	"bytes"
	"embed"
	"io/fs"
	"path"

	sqlmigrate "github.com/rubenv/sql-migrate"
)

const dirName = "migrations"

// migrationsDir is used to hold embedded migration files.
// It must be global var due to embed package convention.
//go:embed migrations
var migrationsDir embed.FS //nolint:gochecknoglobals

// EmbeddedSource holds embedded FS introduced in go 1.16
type EmbeddedSource struct {
	fs embed.FS
}

// NewEmbeddedSource returns a new EmbeddedSource
func NewEmbeddedSource() *EmbeddedSource {
	return &EmbeddedSource{
		fs: migrationsDir,
	}
}

// FindMigrations implements sql-migrate.MigrationSource interface
func (m *EmbeddedSource) FindMigrations() ([]*sqlmigrate.Migration, error) {
	migrationFiles, err := m.fs.ReadDir(dirName)
	if err != nil {
		return nil, err
	}

	migrations := make([]*sqlmigrate.Migration, 0, len(migrationFiles))

	parseMigrationFunc := func(entry fs.DirEntry) error {
		if entry.IsDir() {
			return nil
		}

		file, err := m.fs.Open(path.Join(dirName, entry.Name()))
		if err != nil {
			return err
		}

		defer func() {
			_ = file.Close()
		}()

		buf := bufSeeker{}
		if _, err = buf.ReadFrom(file); err != nil {
			return err
		}

		migration, err := sqlmigrate.ParseMigration(entry.Name(), &buf)
		if err != nil {
			return err
		}

		migrations = append(migrations, migration)
		return nil
	}

	for i := range migrationFiles {
		if err := parseMigrationFunc(migrationFiles[i]); err != nil {
			return nil, err
		}
	}

	return migrations, nil
}

// bufSeeker is a dumb struct implementing io.Seeker interface
// to satisfy sqlmigrate.ParseMigration method
type bufSeeker struct {
	bytes.Buffer
}

// Seek return empty response because sqlmigrate.ParseMigration use it just to check error
func (b bufSeeker) Seek(_ int64, _ int) (int64, error) { return 0, nil }
