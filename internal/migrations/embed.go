// Package migrations exposes goose migration files as an embed.FS so they can
// be compiled into the binary and applied without external files.
package migrations

import "embed"

// Files holds all SQL migration files under ./sql.
//
//go:embed sql
var Files embed.FS
