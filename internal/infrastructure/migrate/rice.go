package migrate

import (
	"os"

	rice "github.com/GeertJohan/go.rice"
	sqlmigrate "github.com/rubenv/sql-migrate"
)

// RiceSource holds rice Box with migrations
type RiceSource struct {
	box *rice.Box
}

// NewRiceSource returns a new RiceSource
func NewRiceSource(box *rice.Box) *RiceSource {
	return &RiceSource{box}
}

// FindMigrations implements rice.MigrationSource interface
func (m *RiceSource) FindMigrations() ([]*sqlmigrate.Migration, error) {
	migrations := make([]*sqlmigrate.Migration, 0)

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, openErr := m.box.Open(path)
		if openErr != nil {
			return openErr
		}

		defer func() {
			_ = file.Close()
		}()

		migration, parseErr := sqlmigrate.ParseMigration(path, file)
		if parseErr != nil {
			return parseErr
		}

		migrations = append(migrations, migration)
		return nil
	}

	if err := m.box.Walk("", walkFunc); err != nil {
		return nil, err
	}

	return migrations, nil
}
