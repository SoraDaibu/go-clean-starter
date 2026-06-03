package cmd

import (
	"context"
	"log/slog"

	"github.com/SoraDaibu/go-clean-starter/builder"
	"github.com/SoraDaibu/go-clean-starter/config"
	"github.com/SoraDaibu/go-clean-starter/migration"

	"github.com/urfave/cli/v3"
)

var TaskCommand = &cli.Command{
	Name:  "task",
	Usage: "Parent Task Command. Exact subcommand is required to run a task.",
	Commands: []*cli.Command{
		{
			Name:  "import",
			Usage: "Import item from files",
			Action: cli.ActionFunc(func(ctx context.Context, c *cli.Command) error {
				cnf, err := config.Load()
				if err != nil {
					return err
				}

				dependencies, err := builder.InitializeDependency(cnf)
				if err != nil {
					return err
				}

				// args
				sourceDir := c.String("source-dir")
				dryRun := c.Bool("dry-run")
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
			}),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "source-dir",
					Usage: "Directory containing CSV files",
					Value: "./internal/task/item/data",
				},
				&cli.BoolFlag{
					Name:  "dry-run",
					Usage: "Validate files without importing",
				},
			},
		},
		// NOTE: Add more subcommands here for new tasks
	},
}
