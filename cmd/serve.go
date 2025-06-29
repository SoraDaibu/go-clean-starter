package cmd

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"

	"github.com/SoraDaibu/go-clean-starter/builder"
	"github.com/SoraDaibu/go-clean-starter/config"
	"github.com/SoraDaibu/go-clean-starter/internal/http"
	"github.com/SoraDaibu/go-clean-starter/migration"
)

var ServeCommand = &cli.Command{
	Name:  "serve",
	Usage: "To run a backend server",
	Action: cli.ActionFunc(func(ctx context.Context, c *cli.Command) error {
		// run server
		log.Info().Msg("starting server by `serve` command...")

		cnf, err := config.Load()
		if err != nil {
			return err
		}

		dn := builder.NewDependencyNeedsAllTrue()
		d, err := builder.Resolve(cnf, dn)
		if err != nil {
			log.Error().Err(err).Msg("failed to resolve dependencies")
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
	}),
}
