package pgtest_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/titusjaka/go-sample/v2/internal/infrastructure/postgres/pgtest"
	"github.com/titusjaka/go-sample/v2/internal/infrastructure/postgres/pgtest/testdata"
)

const (
	envFile = "../../../../.env"
)

func TestInitTestDatabase(t *testing.T) {
	t.Parallel()

	// ====================================================
	// Init DB flags
	flags := pgtest.FlagsFromEnv(t, envFile)
	flags.TestDatabaseTemplate = ""

	// ====================================================
	// Init test database
	conn := pgtest.InitTestDatabase(
		t,
		pgtest.WithFlags(flags),
		pgtest.WithConfigFiles(envFile),
		pgtest.WithApplyMigrations(testdata.Migrations),
	)

	// ====================================================
	// Run test
	var id int32
	err := conn.QueryRow("SELECT id FROM test").Scan(&id)
	assert.NoError(t, err)

	assert.Equal(t, int32(1), id)
}
