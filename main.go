package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/SoraDaibu/go-clean-starter/cmd"
	"github.com/SoraDaibu/go-clean-starter/internal/logger"
)

// Version is the value of release tag embed on build.
var Version = "edge"

// Revision is the value of commit hash embed on build.
var Revision = "latest"

func main() {
	// Bootstrap default logger before config is loaded; builder.Resolve later
	// reconfigures it from the loaded config (level/env).
	logger.SetDefault(logger.New("", ""))

	if err := cmd.Run(context.Background(), os.Args[1:], fmt.Sprintf("%s - %s", Version, Revision), os.Stdout); err != nil {
		slog.Error("command failed", "err", err)
		os.Exit(1)
	}
}
