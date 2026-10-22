package cmd

import (
	"log/slog"

	"github.com/SoraDaibu/go-clean-starter/builder"
	"github.com/SoraDaibu/go-clean-starter/config"
	"github.com/SoraDaibu/go-clean-starter/internal/http"
	"github.com/SoraDaibu/go-clean-starter/migration"
)

func serve() error {
	// run server
	slog.Info("starting server by `serve` command...")

	cnf, err := config.Load()
	if err != nil {
		return err
	}

	dn := builder.NewDependencyNeedsAllTrue()
	d, err := builder.Resolve(cnf, dn)
	if err != nil {
		slog.Error("failed to resolve dependencies", "err", err)
		return err
	}

	// migrate if local
	if cnf.App.Env == "local" {
		if err := migration.Up(getMigrationDatabaseURL()); err != nil {
			return err
		}
	}

	server := http.NewServer(d)
	defer func() { err = server.Close() }()
	server.Run()

	return nil
}
