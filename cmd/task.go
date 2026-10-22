package cmd

import (
	"context"
	"log/slog"

	"github.com/SoraDaibu/go-clean-starter/builder"
	"github.com/SoraDaibu/go-clean-starter/config"
	"github.com/SoraDaibu/go-clean-starter/migration"
)

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
		if err := migration.Up(getMigrationDatabaseURL()); err != nil {
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
