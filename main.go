package main

import (
	"context"
	"fmt"
	"os"

	"github.com/SoraDaibu/go-clean-starter/cmd"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

// Version is the value of release tag embed on build.
var Version = "edge"

// Revision is the value of commit hash embed on build.
var Revision = "latest"

func main() {
	app := &cli.Command{
		Name:    "go-clean-starter",
		Version: fmt.Sprintf("%s - %s", Version, Revision),
		Commands: []*cli.Command{
			cmd.ServeCommand,
			cmd.TaskCommand,
			cmd.MigrationCommand,
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Error().Err(err).Msg("")
	}
}
