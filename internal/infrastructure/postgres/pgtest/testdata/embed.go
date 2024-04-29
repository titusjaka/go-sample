package testdata

import (
	"embed"
)

// Migrations is used to hold embedded migration files.
//
//go:embed *.sql
var Migrations embed.FS
