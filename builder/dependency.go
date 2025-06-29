package builder

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/SoraDaibu/go-clean-starter/config"

	_ "github.com/lib/pq"
)

type Dependency struct {
	Config *config.Config
	DB     *sql.DB
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
	d := &Dependency{
		Config: c,
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
	db, err := sql.Open("postgres", dbURL(d.Config, false))
	if err != nil {
		return fmt.Errorf("%s: %w", dbURL(d.Config, true), err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetConnMaxLifetime(time.Duration(d.Config.DB.Connection.LifetimeSeconds) * time.Second)
	db.SetMaxIdleConns(d.Config.DB.Connection.MaxIdle)
	db.SetMaxOpenConns(d.Config.DB.Connection.MaxOpen)

	d.DB = db
	return nil
}

func dbURL(c *config.Config, mask bool) string {
	password := c.DB.Password
	if mask {
		password = "********"
	}

	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DB.Host,
		c.DB.Port,
		c.DB.User,
		password,
		c.DB.Name,
		c.DB.SSLMode,
	)
}
