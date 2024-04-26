package flags

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // import pg driver
)

// PostgreSQL represents PostgreSQL connection flags
// and provides methods to open a connection to the database.
type PostgreSQL struct {
	DSN string `kong:"required,group='Postgres',name=postgres-dsn,default=localhost,env='POSTGRES_DSN,DATABASE_DSN',help='Data Source Name for PostgreSQL database server.'"`
}

// OpenStdSQLDB opens a new connection to the PostgreSQL database
// using the standard library's sql package.
func (p PostgreSQL) OpenStdSQLDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", p.DSN)
	if err != nil {
		return nil, fmt.Errorf("open PostgreSQL connection: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ping PostgreSQL: %w", err)
	}
	return db, nil
}
