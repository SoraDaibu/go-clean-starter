package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
)

// Run dispatches application commands after parsing their flags.
// args excludes the executable name. Help and version do not initialize services.
func Run(ctx context.Context, args []string, version string, output io.Writer) error {
	command, err := parseCommand(args, output)
	if errors.Is(err, flag.ErrHelp) {
		return nil
	}
	if err != nil {
		return err
	}
	switch command.name {
	case "version":
		_, err = fmt.Fprintln(output, version)
		return err
	case "serve":
		return serve()
	case "migrate", "migrate up", "migrate down":
		return migrate(command.direction)
	case "task import":
		return importItems(ctx, command.sourceDir, command.dryRun)
	default:
		return fmt.Errorf("unknown command %q", command.name)
	}
}

type commandOptions struct {
	name      string
	direction string
	sourceDir string
	dryRun    bool
}

func parseCommand(args []string, output io.Writer) (commandOptions, error) {
	var command commandOptions
	root := flag.NewFlagSet("go-clean-starter", flag.ContinueOnError)
	root.SetOutput(output)
	version := root.Bool("version", false, "Print version and revision")
	root.BoolVar(version, "v", false, "Print version and revision")
	root.Usage = func() {
		fmt.Fprintln(output, "Usage: go-clean-starter <command> [flags]")
		fmt.Fprintln(output, "Commands: serve, migrate [up|down], task import")
		root.PrintDefaults()
	}
	if err := root.Parse(args); err != nil {
		return command, err
	}
	args = root.Args()
	if *version {
		if len(args) != 0 {
			return command, fmt.Errorf("--version does not accept arguments")
		}
		command.name = "version"
		return command, nil
	}
	if len(args) == 0 {
		root.Usage()
		return command, flag.ErrHelp
	}
	if args[0] == "help" {
		if len(args) == 1 {
			root.Usage()
			return command, flag.ErrHelp
		}
		return parseCommand(append(append([]string{}, args[1:]...), "--help"), output)
	}
	command.name = args[0]
	args = args[1:]
	switch command.name {
	case "serve":
	case "migrate":
		if len(args) > 0 && (args[0] == "up" || args[0] == "down") {
			command.direction = args[0]
			command.name += " " + args[0]
			args = args[1:]
		}
	case "task":
		if len(args) == 0 || (len(args) == 1 && (args[0] == "-h" || args[0] == "--help")) {
			fmt.Fprintln(output, "Usage: go-clean-starter task import [flags]")
			return command, flag.ErrHelp
		}
		if args[0] != "import" {
			return command, fmt.Errorf("unknown task %q", args[0])
		}
		command.name = "task import"
		args = args[1:]
	default:
		return command, fmt.Errorf("unknown command %q", command.name)
	}

	flags := flag.NewFlagSet(command.name, flag.ContinueOnError)
	flags.SetOutput(output)
	flags.Usage = func() {
		fmt.Fprintf(output, "Usage: go-clean-starter %s [flags]\n", command.name)
		if command.name == "migrate" {
			fmt.Fprintln(output, "Subcommands: up, down (interactive); omit to migrate using environment variables")
		}
		flags.PrintDefaults()
	}
	if command.name == "task import" {
		flags.StringVar(&command.sourceDir, "source-dir", "./internal/task/item/data", "Directory containing CSV files")
		flags.BoolVar(&command.dryRun, "dry-run", false, "Validate files without importing")
	}
	if err := flags.Parse(args); err != nil {
		return command, err
	}
	if flags.NArg() != 0 {
		return command, fmt.Errorf("%s: unexpected arguments %q", command.name, flags.Args())
	}
	return command, nil
}
