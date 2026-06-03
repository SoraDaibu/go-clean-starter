package migration

import (
	"embed"
	"errors"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed sql/*.sql
var sqlFiles embed.FS

func prepare(source string) (*migrate.Migrate, error) {
	driver, err := iofs.New(sqlFiles, "sql")
	if err != nil {
		return nil, err
	}

	return migrate.NewWithSourceInstance("iofs", driver, source)
}

func Up(source string) error {
	m, err := prepare(source)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("migration failed", "err", err)
		return err
	}

	slog.Info("Migration Up completed")
	return nil
}

func Down(source string) error {
	m, err := prepare(source)
	if err != nil {
		return err
	}

	if err := m.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("migration failed", "err", err)
		return err
	}

	slog.Info("Migration Down completed")
	return nil
}
