package builder

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/SoraDaibu/go-clean-starter/config"
	"github.com/SoraDaibu/go-clean-starter/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Dependency struct {
	Config *config.Config
	Logger logger.Logger
	DB     *pgxpool.Pool
	HTTP   *http.Client
}

type (
	NeedsDB bool
)

// DependencyNeeds specifies which dependencies are required
// Add more dependency as needed such as S3, Redis, Stripe, etc.
type DependencyNeeds struct {
	needsDB NeedsDB
}

func NewDependencyNeeds(needsDB NeedsDB) *DependencyNeeds {
	return &DependencyNeeds{
		needsDB: needsDB,
	}
}

func NewDependencyNeedsAllTrue() *DependencyNeeds {
	return &DependencyNeeds{
		needsDB: true,
	}
}

func Resolve(c *config.Config, dn *DependencyNeeds) (*Dependency, error) {
	l := logger.New(c.App.LogLevel, c.App.Env)
	// Make slog's package-level functions (used by bootstrap/CLI code) share the
	// same configured backend.
	logger.SetDefault(l)

	d := &Dependency{
		Config: c,
		Logger: l,
		HTTP: &http.Client{
			Timeout: time.Duration(c.HTTP.TimeoutSeconds) * time.Second,
		},
	}

	if dn.needsDB {
		if err := connectDB(d); err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}
	}

	return d, nil
}

func connectDB(d *Dependency) error {
	// Parse config with pool settings
	config, err := pgxpool.ParseConfig(dbURL(d.Config, false))
	if err != nil {
		return fmt.Errorf("%s: %w", dbURL(d.Config, true), err)
	}

	// Set connection pool settings
	config.MaxConnLifetime = time.Duration(d.Config.DB.Connection.LifetimeSeconds) * time.Second
	config.MaxConnIdleTime = time.Duration(d.Config.DB.Connection.LifetimeSeconds) * time.Second
	config.MinConns = int32(d.Config.DB.Connection.MinIdleConns)
	config.MaxConns = int32(d.Config.DB.Connection.MaxOpen)

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	d.DB = pool
	return nil
}

func dbURL(c *config.Config, mask bool) string {
	password := c.DB.Password
	if mask {
		password = "********"
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.DB.User,
		password,
		c.DB.Host,
		c.DB.Port,
		c.DB.Name,
		c.DB.SSLMode,
	)
}
