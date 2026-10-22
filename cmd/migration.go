package cmd

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"os"
	"strings"
	"syscall"

	"github.com/SoraDaibu/go-clean-starter/internal/logger"
	"github.com/SoraDaibu/go-clean-starter/migration"

	"golang.org/x/term"
)

func migrate(direction string) error {
	if direction == "" {
		return migration.Up(getMigrationDatabaseURL())
	}
	logger.SetDefault(logger.New("info", "local"))
	databaseURL, err := scanDatasource()
	if err != nil {
		return err
	}
	if direction == "up" {
		err = migration.Up(databaseURL)
	} else {
		err = migration.Down(databaseURL)
	}
	if err != nil {
		return err
	}
	slog.Info("migration complete", "direction", direction)
	return nil
}

func getMigrationDatabaseURL() string {
	return migrationDatabaseURL(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"))
}

func migrationDatabaseURL(host, port, database, username, password string) string {
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(strings.TrimSpace(username), password),
		Host:     net.JoinHostPort(strings.TrimSpace(host), strings.TrimSpace(port)),
		Path:     "/" + strings.TrimSpace(database),
		RawQuery: url.Values{"sslmode": []string{migrationSSLMode()}}.Encode(),
	}
	return u.String()
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

	return migrationDatabaseURL(host, port, database, username, string(bytePassword)), nil
}

func migrationSSLMode() string {
	if mode := os.Getenv("PGSSLMODE"); mode != "" {
		return mode
	}
	return "disable"
}
