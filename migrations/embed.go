package migrations

import (
	"embed"
)

// Dir is used to hold embedded migration files.
//
//go:embed *.sql
var Dir embed.FS
