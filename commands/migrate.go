package commands

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"text/template"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/urfave/cli/v2"

	"github.com/titusjaka/go-sample/internal/infrastructure/log"
	"github.com/titusjaka/go-sample/internal/infrastructure/migrate"
)

const templateContent = `-- +migrate Up

-- +migrate Down
`

// NewMigrateCmd creates a new migrate CLI sub-command
func NewMigrateCmd() *cli.Command {
	return &cli.Command{
		Name:        "migrate",
		Usage:       "create a new migration, apply (or rollback) migrations to DB.",
		Description: "Utility for easy migrations handling.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dsn",
				Usage:    "Data Source Name for PostgreSQL database server",
				EnvVars:  []string{"DATABASE_DSN", "POSTGRES_DSN"},
				Required: true,
			},
		},
		Subcommands: []*cli.Command{
			{
				Name:   "up",
				Usage:  "apply all migrations to DB",
				Action: migrateUp,
			},
			{
				Name:   "down",
				Usage:  "Rollback a [number] of migrations. Pass a [number] as the first argument.",
				Action: migrateDown,
			},
			{
				Name:   "create",
				Usage:  "Create a new blank migration file. Pass a [name] as the first argument.",
				Action: migrateCreate,
			},
		},
	}
}

func migrateUp(c *cli.Context) error {
	logger := log.New()

	migrator, err := initMigrator(c.String("dsn"))
	if err != nil {
		return fmt.Errorf("can't init migrator: %w", err)
	}

	applied, err := migrator.Up(context.Background())
	if err != nil {
		return fmt.Errorf("can't apply migrations: %w", err)
	}

	logger.Info("👟 ➡ migration(s) applied successfully", log.Field("applied", applied))
	return nil
}

func migrateDown(c *cli.Context) error {
	logger := log.New()

	steps, err := strconv.Atoi(c.Args().First())
	if err != nil {
		return fmt.Errorf("can't parse number of steps %w", err)
	}

	migrator, err := initMigrator(c.String("dsn"))
	if err != nil {
		return fmt.Errorf("can't init migrator: %w", err)
	}

	reverted, err := migrator.Down(context.Background(), steps)
	if err != nil {
		return fmt.Errorf("can't revert migrations: %w", err)
	}

	logger.Info("🤖 ➡ migration(s) reverted successfully", log.Field("reverted", reverted))
	return nil
}

func migrateCreate(c *cli.Context) error {
	logger := log.New()

	dir := filepath.Clean("migrations")
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("can't create directory for migrations: %w", err)
	}

	fileName := fmt.Sprintf("%d_%s.sql", time.Now().Unix(), c.Args().First())
	pathName := path.Join(dir, fileName)
	file, err := os.Create(pathName)
	if err != nil {
		return fmt.Errorf("can't create migration file (%q): %w", pathName, err)
	}

	tpl, err := template.New("new_migration").Parse(templateContent)
	if err != nil {
		return fmt.Errorf("can't parse migration template: %w", err)
	}

	if err := tpl.Execute(file, nil); err != nil {
		return fmt.Errorf("can't execute template: %w", err)
	}

	logger.Info("🧶 ➡ created new migration", log.Field("path", pathName))
	return nil
}

func initMigrator(dsn string) (*migrate.Migrator, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("can't open db connection: %w", err)
	}

	source := migrate.NewRiceSource(rice.MustFindBox("../migrations"))
	migrator := migrate.NewMigrator(db, source)

	return migrator, nil
}
