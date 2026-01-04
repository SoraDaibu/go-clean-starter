package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	App struct {
		Env        string
		LogLevel   string
		ListenPort uint16
	}
	DB struct {
		Host       string
		Port       int
		User       string
		Password   string
		Name       string
		SSLMode    string
		Connection struct {
			MinIdleConns    int
			MaxOpen         int
			LifetimeSeconds int
		}
	}
	HTTP struct {
		TimeoutSeconds int
	}
}

func Load() (*Config, error) {
	cnf := &Config{}
	var err error

	// app
	cnf.App.Env = os.Getenv("APP_ENV")
	cnf.App.LogLevel = os.Getenv("APP_LOG_LEVEL")
	listenPort, err := strconv.Atoi(os.Getenv("APP_LISTEN_PORT"))
	if err != nil {
		return nil, fmt.Errorf("failed to get LISTEN_PORT: %w", err)
	}
	cnf.App.ListenPort = uint16(listenPort)

	// database
	cnf.DB.Host = os.Getenv("DB_HOST")
	cnf.DB.Port, err = strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, fmt.Errorf("failed to get DB_PORT: %w", err)
	}
	cnf.DB.User = os.Getenv("DB_USER")
	cnf.DB.Password = os.Getenv("DB_PASSWORD")
	cnf.DB.Name = os.Getenv("DB_NAME")
	cnf.DB.SSLMode = os.Getenv("PGSSLMODE")

	// Add missing database connection pool configuration
	cnf.DB.Connection.MinIdleConns, err = strconv.Atoi(os.Getenv("DB_MIN_IDLE_CONNS"))
	if err != nil {
		return nil, fmt.Errorf("failed to get DB_MIN_IDLE_CONNS: %w", err)
	}
	cnf.DB.Connection.MaxOpen, err = strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNS"))
	if err != nil {
		return nil, fmt.Errorf("failed to get DB_MAX_OPEN_CONNS: %w", err)
	}
	cnf.DB.Connection.LifetimeSeconds, err = strconv.Atoi(os.Getenv("DB_CONN_LIFETIME_SECONDS"))
	if err != nil {
		return nil, fmt.Errorf("failed to get DB_CONN_LIFETIME_SECONDS: %w", err)
	}

	// http
	// timeout seconds for http request
	cnf.HTTP.TimeoutSeconds, err = strconv.Atoi(os.Getenv("HTTP_TIMEOUT_SECONDS"))
	if err != nil {
		return nil, fmt.Errorf("failed to get HTTP_TIMEOUT_SECONDS: %w", err)
	}

	return cnf, nil
}
