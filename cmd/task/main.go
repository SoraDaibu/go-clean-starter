package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/SoraDaibu/go-clean-starter/builder"
	"github.com/SoraDaibu/go-clean-starter/config"
	"github.com/SoraDaibu/go-clean-starter/internal/logger"
	"github.com/SoraDaibu/go-clean-starter/migration"
)

func main() {
	logger.SetDefault(logger.New("", ""))
	if err := run(os.Args[1:], os.Stdout); err != nil {
		slog.Error("command failed", "err", err)
		os.Exit(1)
	}
}

type options struct {
	sourceDir string
	dryRun    bool
}

func parseFlags(args []string, output io.Writer) (options, error) {
	var opts options
	flags := flag.NewFlagSet("task", flag.ContinueOnError)
	flags.SetOutput(output)
	flags.StringVar(&opts.sourceDir, "source-dir", "./internal/task/item/data", "Directory containing CSV files")
	flags.BoolVar(&opts.dryRun, "dry-run", false, "Validate files without importing")
	if err := flags.Parse(args); err != nil {
		return opts, err
	}
	if flags.NArg() != 0 {
		return opts, fmt.Errorf("task: unexpected arguments %q", flags.Args())
	}
	return opts, nil
}

func run(args []string, output io.Writer) error {
	opts, err := parseFlags(args, output)
	if errors.Is(err, flag.ErrHelp) {
		return nil
	}
	if err != nil {
		return err
	}
	return importItems(context.Background(), opts.sourceDir, opts.dryRun)
}

func importItems(ctx context.Context, sourceDir string, dryRun bool) error {
	cnf, err := config.Load()
	if err != nil {
		return err
	}

	dependencies, err := builder.InitializeDependency(cnf)
	if err != nil {
		return err
	}

	defer dependencies.DB.Close()

	slog.Info("importing items", "source-dir", sourceDir, "dry-run", dryRun)

	// migrate if local
	if cnf.App.Env == "local" {
		if err := migration.Up(migration.DatabaseURLFromEnv()); err != nil {
			return err
		}
	}

	task := builder.InitializeItemTaskUsecase(dependencies)
	err = task.ImportItems(ctx, sourceDir, dryRun)
	if err != nil {
		return err
	}

	slog.Info("item import success 🎉")

	return nil
}
