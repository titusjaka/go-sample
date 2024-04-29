package main_test

import (
	"flag"
	"os"
	"testing"

	"github.com/alecthomas/kong"

	kongdotenv "github.com/titusjaka/kong-dotenv-go"

	"github.com/titusjaka/go-sample/internal/infrastructure/postgres/pgtest"
	"github.com/titusjaka/go-sample/migrations"
)

func TestMain(m *testing.M) {
	flag.Parse()

	if testing.Short() {
		return
	}

	var pg pgtest.Flags

	parser, err := kong.New(&pg, kong.Configuration(kongdotenv.ENVFileReader, ".env"))
	if err != nil {
		panic(err)
	}

	if _, err = parser.Parse(nil); err != nil {
		panic(err)
	}

	if err = pgtest.CreateTemplateDatabase(pg, migrations.Dir); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}
