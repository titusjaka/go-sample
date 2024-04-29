package pgtest

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/titusjaka/go-sample/internal/infrastructure/postgres/pgmigrator"
)

func CreateTemplateDatabase(flags Flags, migrations fs.FS) error {
	initDB, err := flags.Postgres.OpenStdSQLDB()
	if err != nil {
		return fmt.Errorf("open initial DB connection: %w", err)
	}

	var exists bool

	query := "SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1);"

	if err = initDB.QueryRow(query, flags.TestDatabaseTemplate).Scan(&exists); err != nil {
		return fmt.Errorf("check DB existance: %w", err)
	}

	if !exists {
		_, err = initDB.Exec(fmt.Sprintf("CREATE DATABASE %s", flags.TestDatabaseTemplate))
		if err != nil {
			return err
		}
	}

	if err = initDB.Close(); err != nil {
		return fmt.Errorf("close inital DB connection: %w", err)
	}

	flags.Postgres.Database = flags.TestDatabaseTemplate

	db, err := flags.Postgres.OpenStdSQLDB()
	if err != nil {
		return err
	}
	defer func() {
		_ = db.Close()
	}()

	if _, err = pgmigrator.NewMigrator(db, migrations).Up(context.Background()); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}

	return nil
}
