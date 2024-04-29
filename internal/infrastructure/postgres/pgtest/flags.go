package pgtest

import (
	"github.com/titusjaka/go-sample/internal/infrastructure/postgres"
)

type Flags struct {
	Postgres             postgres.Flags `kong:"embed"`
	TestDatabaseTemplate string         `kong:"optional,group='Postgres',name=postgres-test-database,default=test,env=POSTGRES_TEST_DATABASE_TEMPLATE,help='PostgreSQL database template for tests.'"`
}
