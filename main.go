package main

import (
	"log"
	"os"
	"sort"

	"github.com/urfave/cli/v2"

	"github.com/titusjaka/go-sample/commands"
)

func main() {
	app := &cli.App{
		Name:  "go-sample",
		Usage: "Use it as a starting point for your Go backend application.",
		Description: "Go backend application using a modular project layout.\n" +
			"Ready to go REST-API service with a specimen route “snippets.”\n" +
			"It also contains some infrastructure code to simplify routine operations.",
		Commands: []*cli.Command{
			commands.NewServerCmd(),
			commands.NewMigrateCmd(),
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
