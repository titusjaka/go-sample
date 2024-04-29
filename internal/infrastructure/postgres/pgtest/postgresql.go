package pgtest

import (
	"cmp"
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"strings"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/stretchr/testify/require"

	kongdotenv "github.com/titusjaka/kong-dotenv-go"

	"github.com/titusjaka/go-sample/internal/infrastructure/postgres/pgmigrator"
)

// InitTestDatabase creates initializes a new dummy database in s PostgreSQL server
func InitTestDatabase(t testing.TB, opts ...Option) *sql.DB {
	t.Helper()

	conf := defaultConfig()
	for _, opt := range opts {
		opt(conf)
	}

	// Init PostgreSQL flags that will be used to create test DB
	initialFlags := cmp.Or(conf.flags, FlagsFromEnv(t, conf.configFiles...))

	// Init PostgreSQL connection, this is needed to create test DB
	initialDB, err := initialFlags.Postgres.OpenStdSQLDB()
	require.NoError(t, err)

	defer func() {
		// Close initial connection
		closeErr := initialDB.Close()
		require.NoError(t, closeErr)
	}()

	// Create test DB
	databaseName := DatabaseName(t)
	CreateDatabase(t, initialDB, databaseName, initialFlags.TestDatabaseTemplate)

	// Init connection with test DB
	testFlags := initialFlags
	testFlags.Postgres.Database = databaseName
	testDB, err := testFlags.Postgres.OpenStdSQLDB()
	require.NoError(t, err)

	// apply DB actions
	for _, pgxFunc := range conf.pgxFuncs {
		pgxFunc(t, testDB)
	}

	t.Cleanup(func() {
		tearDownTestDatabase(t, testDB, initialFlags, databaseName)
	})

	return testDB
}

// FlagsFromEnv returns an instance of flags.Postgres populated from environment variables or .env files.
func FlagsFromEnv(t testing.TB, paths ...string) Flags {
	t.Helper()

	var pg Flags

	parser, err := kong.New(&pg, kong.Configuration(kongdotenv.ENVFileReader, paths...))
	require.NoError(t, err)

	_, err = parser.Parse(nil)
	require.NoError(t, err)

	return pg
}

// CreateDatabase creates a new database within a PostgreSQL connection
func CreateDatabase(t testing.TB, db *sql.DB, databaseName, testDatabaseTemplateName string) {
	t.Helper()

	query := fmt.Sprintf("CREATE DATABASE %s", databaseName)
	if testDatabaseTemplateName != "" {
		query += fmt.Sprintf(" TEMPLATE = %s", testDatabaseTemplateName)
	}

	_, err := db.Exec(query)
	require.NoError(t, err)
}

// DropDatabase drops database
func DropDatabase(t testing.TB, db *sql.DB, databaseName string, force bool) {
	t.Helper()

	query := fmt.Sprintf("DROP DATABASE %s", databaseName)
	if force {
		query += " WITH (FORCE)"
	}

	_, err := db.Exec(query)
	require.NoError(t, err)
}

// ApplyMigrations applies migration to a provided database
func ApplyMigrations(t testing.TB, db *sql.DB, source fs.FS) {
	t.Helper()

	migrator := pgmigrator.NewMigrator(db, source)

	_, err := migrator.Up(context.Background())
	require.NoError(t, err)
}

// tearDownTestDatabase cleans DB after test
func tearDownTestDatabase(t testing.TB, testConn *sql.DB, initFlags Flags, dbNameToDrop string) {
	require.Truef(t, strings.HasPrefix(dbNameToDrop, "test"), "DB name must start with 'test' prefix")

	// close connection to test DB
	closeErr := testConn.Close()
	require.NoError(t, closeErr)

	// Init Common PostgreSQL connection to drop test DB
	commonConn, err := initFlags.Postgres.OpenStdSQLDB()
	require.NoError(t, err)

	// Drop test DB
	DropDatabase(t, commonConn, dbNameToDrop, true)
}
