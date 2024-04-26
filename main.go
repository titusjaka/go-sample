package main

import (
	"github.com/alecthomas/kong"

	kongdotenv "github.com/titusjaka/kong-dotenv-go"

	"github.com/titusjaka/go-sample/commands"
	"github.com/titusjaka/go-sample/internal/infrastructure/kongflag"
)

type App struct {
	EnvFile kongdotenv.ENVFileConfig `kong:"optional,name=env-file,help='Path to .env file'"`

	Migrate commands.MigrateCmd `kong:"cmd,name=migrate,help='Create a new migration, apply (or rollback) migrations to DB.'"`
	Server  commands.ServerCmd  `kong:"cmd,name=server,default=1,help='Start the HTTP server.'"`
}

var (
	serviceName  = "go-sample"
	Version      = "v0.0.1"
	GitCommitSHA = "unknown"
	GitBranch    = "unknown"
)

func main() {
	var app App
	kCtx := kong.Parse(
		&app,
		kong.Name(serviceName),
		kong.Description(
			"Go backend application using a modular project layout.\n"+
				"Ready to go REST-API service with a specimen route “snippets”.\n"+
				"It also contains some infrastructure code to simplify routine operations.",
		),
		kong.Vars{
			kongflag.ServiceName:    serviceName,
			kongflag.ServiceVersion: Version,
			kongflag.GitCommitSHA:   GitCommitSHA,
			kongflag.GitBranch:      GitBranch,
		},
		kongflag.WithBuildInfo(),
		kongflag.WithVersion(Version),
		kongflag.WithDumpEnvs(),
	)
	kCtx.FatalIfErrorf(kCtx.Run())
}
