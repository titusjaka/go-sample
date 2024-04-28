package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver
)

// Flags represents PostgreSQL connection flags
// and provides methods to open a connection to the database.
type Flags struct {
	Host     string `kong:"optional,group='Postgres',name=postgres-host,default=localhost,env=POSTGRES_HOST,help='PostgreSQL address host.'"`
	Port     uint32 `kong:"optional,group='Postgres',name=postgres-port,default=5432,env=POSTGRES_PORT,help='PostgreSQL address port.'"`
	Username string `kong:"optional,group='Postgres',name=postgres-username,default=postgres,env=POSTGRES_USERNAME,help='PostgreSQL username.'"`
	Password string `kong:"optional,group='Postgres',name=postgres-password,env=POSTGRES_PASSWORD,help='PostgreSQL password.'" json:"-"`
	Database string `kong:"optional,group='Postgres',name=postgres-database,default=postgres,env=POSTGRES_DATABASE,help='PostgreSQL database.'"`

	TLSMode string `kong:"optional,group='Postgres',name=postgres-tls-mode,default=disable,env=POSTGRES_TLS_MODE,help='PostgreSQL TLS mode.'"`
}

// OpenStdSQLDB opens a new connection to the PostgreSQL database
// using the standard library's sql package.
func (p Flags) OpenStdSQLDB() (*sql.DB, error) {
	db, err := sql.Open("pgx/v5", p.BuildConnectionString())
	if err != nil {
		return nil, fmt.Errorf("open PostgreSQL connection: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ping PostgreSQL: %w", err)
	}
	return db, nil
}

func (p Flags) BuildConnectionString() string {
	connString := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s sslmode=%s ",
		p.Host,
		p.Port,
		p.Username,
		p.Database,
		p.TLSMode,
	)

	if p.Password != "" {
		connString += fmt.Sprintf("password=%s ", p.Password)
	}

	return connString
}
