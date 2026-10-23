package migration

import (
	"net"
	"net/url"
	"os"
	"strings"
)

// DatabaseURLFromEnv builds a PostgreSQL URL from the database environment variables.
func DatabaseURLFromEnv() string {
	return DatabaseURL(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"))
}

// DatabaseURL builds a PostgreSQL URL, preserving and escaping credentials.
func DatabaseURL(host, port, database, username, password string) string {
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(strings.TrimSpace(username), password),
		Host:     net.JoinHostPort(strings.TrimSpace(host), strings.TrimSpace(port)),
		Path:     "/" + strings.TrimSpace(database),
		RawQuery: url.Values{"sslmode": []string{migrationSSLMode()}}.Encode(),
	}
	return u.String()
}

func migrationSSLMode() string {
	if mode := os.Getenv("PGSSLMODE"); mode != "" {
		return mode
	}
	return "disable"
}
