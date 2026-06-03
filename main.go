package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/SoraDaibu/go-clean-starter/cmd"
	"github.com/SoraDaibu/go-clean-starter/internal/logger"

	"github.com/urfave/cli/v3"
)

// Version is the value of release tag embed on build.
var Version = "edge"

// Revision is the value of commit hash embed on build.
var Revision = "latest"

func main() {
	// Bootstrap default logger before config is loaded; builder.Resolve later
	// reconfigures it from the loaded config (level/env).
	logger.SetDefault(logger.New("", ""))

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
		slog.Error("command failed", "err", err)
	}
}
