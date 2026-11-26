// Package migrations embeds the SQL migration files into the backend binary so
// that migrations can be applied without shipping the .sql files separately.
package migrations

import "embed"

// FS contains all migration files in this directory.
//
//go:embed *.sql
var FS embed.FS
