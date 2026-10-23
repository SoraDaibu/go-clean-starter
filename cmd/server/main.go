package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/SoraDaibu/go-clean-starter/builder"
	"github.com/SoraDaibu/go-clean-starter/config"
	"github.com/SoraDaibu/go-clean-starter/internal/http"
	"github.com/SoraDaibu/go-clean-starter/internal/logger"
	"github.com/SoraDaibu/go-clean-starter/migration"
)

// Version and Revision are set with -ldflags at build time.
var Version = "edge"
var Revision = "latest"

func main() {
	logger.SetDefault(logger.New("", ""))
	if err := run(os.Args[1:], os.Stdout); err != nil {
		slog.Error("command failed", "err", err)
		os.Exit(1)
	}
}

func run(args []string, output io.Writer) error {
	flags := flag.NewFlagSet("server", flag.ContinueOnError)
	flags.SetOutput(output)
	version := flags.Bool("version", false, "Print version and revision")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("server: unexpected arguments %q", flags.Args())
	}
	if *version {
		_, err := fmt.Fprintf(output, "%s - %s\n", Version, Revision)
		return err
	}
	return serve()
}

func serve() error {
	// run server
	slog.Info("starting server...")

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
		if err := migration.Up(migration.DatabaseURLFromEnv()); err != nil {
			return err
		}
	}

	server := http.NewServer(d)
	defer func() { err = server.Close() }()
	server.Run()

	return nil
}
