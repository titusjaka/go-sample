package pgtest

import (
	"database/sql"
	"io/fs"
	"testing"
)

type config struct {
	configFiles []string
	pgxFuncs    []func(t testing.TB, db *sql.DB)
	flags       Flags
}

func defaultConfig() *config {
	return &config{}
}

type Option func(*config)

func WithFlags(flags Flags) Option {
	return func(c *config) {
		c.flags = flags
	}
}

func WithConfigFiles(files ...string) Option {
	return func(c *config) {
		c.configFiles = files
	}
}

func WithApplyMigrations(source fs.FS) Option {
	return func(c *config) {
		c.pgxFuncs = append(c.pgxFuncs, func(t testing.TB, db *sql.DB) {
			ApplyMigrations(t, db, source)
		})
	}
}
