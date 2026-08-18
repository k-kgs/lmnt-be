// Package migrations embeds the SQL migration files directly into the
// compiled binary, so a deployed container can self-migrate on boot without
// needing the /migrations directory mounted or a separate deploy step —
// see Run() in this package, called from cmd/server/main.go before the HTTP
// server starts. This is what makes a Render deploy on push actually apply
// schema changes, not just start a container against a stale schema.
package migrations

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed *.sql
var files embed.FS

// Run applies all pending "up" migrations against databaseURL. Safe to call
// on every boot — migrate.ErrNoChange (nothing pending) is not an error.
func Run(databaseURL string) error {
	src, err := iofs.New(files, ".")
	if err != nil {
		return fmt.Errorf("load embedded migrations: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", src, databaseURL)
	if err != nil {
		return fmt.Errorf("init migrator: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}
