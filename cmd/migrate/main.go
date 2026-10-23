package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"syscall"

	"github.com/SoraDaibu/go-clean-starter/internal/logger"
	"github.com/SoraDaibu/go-clean-starter/migration"

	"golang.org/x/term"
)

func main() {
	logger.SetDefault(logger.New("", ""))
	if err := run(os.Args[1:], os.Stdout); err != nil {
		slog.Error("command failed", "err", err)
		os.Exit(1)
	}
}

type options struct {
	direction   string
	interactive bool
}

func parseFlags(args []string, output io.Writer) (options, error) {
	var opts options
	flags := flag.NewFlagSet("migrate", flag.ContinueOnError)
	flags.SetOutput(output)
	flags.StringVar(&opts.direction, "direction", "up", "Migration direction: up or down (one step)")
	flags.BoolVar(&opts.interactive, "interactive", false, "Prompt for database connection details")
	if err := flags.Parse(args); err != nil {
		return opts, err
	}
	if flags.NArg() != 0 {
		return opts, fmt.Errorf("migrate: unexpected arguments %q", flags.Args())
	}
	if opts.direction != "up" && opts.direction != "down" {
		return opts, fmt.Errorf("invalid migration direction %q", opts.direction)
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
	databaseURL := migration.DatabaseURLFromEnv()
	if opts.interactive {
		logger.SetDefault(logger.New("info", "local"))
		databaseURL, err = scanDatasource()
		if err != nil {
			return err
		}
	}
	if opts.direction == "down" {
		return migration.Down(databaseURL)
	}
	return migration.Up(databaseURL)
}

//nolint:forbidigo
func scanDatasource() (string, error) {
	var host,
		port,
		database,
		username string

	fmt.Println("Enter database connection info.")

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("Host (default: %s)> ", os.Getenv("DB_HOST"))
	scanner.Scan()

	if host = scanner.Text(); host == "" {
		host = os.Getenv("DB_HOST")
	}

	fmt.Printf("Port (default: %s)> ", os.Getenv("DB_PORT"))
	scanner.Scan()

	if port = scanner.Text(); port == "" {
		port = os.Getenv("DB_PORT")
	}

	fmt.Printf("DB Name (default: %s)> ", os.Getenv("DB_NAME"))
	scanner.Scan()

	if database = scanner.Text(); database == "" {
		database = os.Getenv("DB_NAME")
	}

	fmt.Printf("Username (default: %s)> ", os.Getenv("DB_USER"))
	scanner.Scan()

	if username = scanner.Text(); username == "" {
		username = os.Getenv("DB_USER")
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("read connection details: %w", err)
	}

	fmt.Print("Password > ")

	bytePassword, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return "", fmt.Errorf("read password: %w", err)
	}

	fmt.Println("")

	return migration.DatabaseURL(host, port, database, username, string(bytePassword)), nil
}
