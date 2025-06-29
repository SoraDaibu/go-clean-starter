package cmd

import (
	"bufio"
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"syscall"

	"github.com/SoraDaibu/go-clean-starter/migration"

	"golang.org/x/term"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

var cliOutput = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

var MigrationCommand = &cli.Command{
	Name: "migrate",
	Action: cli.ActionFunc(func(ctx context.Context, c *cli.Command) error {
		return migration.Up(getMigrationDatabaseURL())
	}),
	Commands: []*cli.Command{{
		Name: "up",
		Action: cli.ActionFunc(func(ctx context.Context, c *cli.Command) error {
			log.Logger = cliOutput
			if err := migration.Up("postgres://" + scanDatasource()); err != nil {
				return err
			}

			log.Info().Msg("Successfully upped ⤴️")

			return nil
		}),
	}, {
		Name: "down",
		Action: cli.ActionFunc(func(ctx context.Context, c *cli.Command) error {
			log.Logger = cliOutput
			if err := migration.Down("postgres://" + scanDatasource()); err != nil {
				return err
			}

			log.Info().Msg("Successfully downed ⤵️")

			return nil
		}),
	}},
}

func getMigrationDatabaseURL() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_NAME")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	sslMode := os.Getenv("PGSSLMODE")
	if sslMode == "" {
		sslMode = "disable"
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		strings.TrimSpace(username),
		strings.TrimSpace(password),
		strings.TrimSpace(host),
		strings.TrimSpace(port),
		strings.TrimSpace(database),
	) + fmt.Sprintf("?%s", url.Values{
		"sslmode": []string{sslMode},
	}.Encode())
}

//nolint:forbidigo
func scanDatasource() string {
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

	fmt.Print("Password > ")

	bytePassword, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		log.Panic().Err(err).Msg("")
	}

	fmt.Println("")

	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		strings.TrimSpace(username),
		strings.TrimSpace(string(bytePassword)),
		strings.TrimSpace(host),
		strings.TrimSpace(port),
		strings.TrimSpace(database),
	) + fmt.Sprintf("?%s", url.Values{
		"charset": []string{"utf8mb4"},
		"loc":     []string{"UTC"},
	}.Encode())
}
