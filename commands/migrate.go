package commands

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/titusjaka/go-sample/v2/commands/flags"
	"github.com/titusjaka/go-sample/v2/internal/infrastructure/postgres"
	"github.com/titusjaka/go-sample/v2/internal/infrastructure/postgres/pgmigrator"
	"github.com/titusjaka/go-sample/v2/internal/infrastructure/postgres/pgtest"
	"github.com/titusjaka/go-sample/v2/migrations"
)

// MigrateCmd implements kong.Command for migrations. To use this command you need to add migrate.Command
// to the application structure and bind a migration source.
//
// Usage:
//
//	type App struct {
//		// Add other application-related flags here
//		// ...
//		Migrate commands.MigrateCmd `kong:"cmd,name=migrate,help='Apply database migrations.'"`
//	}
//
//	func main() {
//		var app App
//		kCtx := kong.Parse(
//			&app,
//			// Add other params, if needed
//			// …
//		)
//		kCtx.FatalIfErrorf(kCtx.Run())
//	}
//
// CLI usage:
//
//	$ go run main.go migrate up
type MigrateCmd struct {
	Create CreateCmd `kong:"cmd,name=create,help='Create a new blank migration file. Pass a [name] as the first argument.'"`
	Up     UpCmd     `kong:"cmd,name=up,default=1,help='Apply all database migrations.'"`
	Down   DownCmd   `kong:"cmd,name=down,help='Rollback a [number] of migrations. Pass a [number] as the flag.'"`

	TestDB InitTestDBCmd `kong:"cmd,name=init-test-db,help='Init test database template.'"`
}

// ============================================================================
// Sub-commands

// CreateCmd represents a CLI sub-command to create a new migration file
type CreateCmd struct {
	Directory string `kong:"default='./migrations',help='Directory to store migration files'"`
	Name      string `kong:"arg,required,help='Migration name'"`

	Logger flags.Logger `kong:"embed"`
}

// UpCmd represents a CLI sub-command to apply all migrations to DB
type UpCmd struct {
	Postgres postgres.Flags `kong:"embed"`
	Logger   flags.Logger   `kong:"embed"`
}

// DownCmd represents a CLI sub-command to roll back a specified number of migrations
type DownCmd struct {
	Postgres postgres.Flags `kong:"embed"`
	Logger   flags.Logger   `kong:"embed"`

	Steps int `kong:"required,default='1',name=steps,help='Number of migrations to revert'"`
}

// InitTestDBCmd represents a CLI sub-command to create a new template database for testing
type InitTestDBCmd struct {
	Postgres pgtest.Flags `kong:"embed"`
	Logger   flags.Logger `kong:"embed"`
}

// ============================================================================
// Actions

// Run (CreateCmd) creates a new migration file
func (c CreateCmd) Run() error {
	logger := c.Logger.Init()

	filename, err := pgmigrator.Create(c.Directory, c.Name)
	if err != nil {
		return err
	}

	logger.Info("🧶 ➡ created new migration", slog.String("path", filename))

	return nil
}

// Run (UpCmd) applies all migrations to DB
func (c UpCmd) Run() error {
	logger := c.Logger.Init()

	db, err := c.Postgres.OpenStdSQLDB()
	if err != nil {
		return fmt.Errorf("init DB: %w", err)
	}

	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			logger.Error("close pg connection", slog.Any("err", closeErr))
		}
	}()

	applied, err := pgmigrator.NewMigrator(db, migrations.Dir).Up(context.Background())
	if err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}

	logger.Info("👟 ➡ migration(s) applied successfully", slog.Int("applied", applied))

	return nil
}

// Run (DownCmd) reverts a specified number of migrations
func (c DownCmd) Run() error {
	logger := c.Logger.Init()

	db, err := c.Postgres.OpenStdSQLDB()
	if err != nil {
		return fmt.Errorf("init DB: %w", err)
	}

	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			logger.Error("close pg connection", slog.Any("err", closeErr))
		}
	}()

	reverted, err := pgmigrator.NewMigrator(db, migrations.Dir).Down(context.Background(), c.Steps)
	if err != nil {
		return fmt.Errorf("revert migrations: %w", err)
	}

	logger.Info("🤖 ➡ migration(s) reverted successfully", slog.Int("reverted", reverted))

	return nil
}

// Run (InitTestDBCmd) creates a new template database for testing
func (c InitTestDBCmd) Run() error {
	logger := c.Logger.Init()

	logger.Info("🧱 ➡ creating template database...")

	if err := pgtest.CreateTemplateDatabase(c.Postgres, migrations.Dir); err != nil {
		return fmt.Errorf("create template DB: %w", err)
	}

	logger.Info("🏠 ➡ template database created successfully", slog.String("name", c.Postgres.TestDatabaseTemplate))

	return nil
}
