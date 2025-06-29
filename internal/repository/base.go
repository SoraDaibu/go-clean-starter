package repository

import (
	"context"
	"database/sql"

	"github.com/SoraDaibu/go-clean-starter/internal/sqlc"
)

// BaseRepository provides common database operations
// Following OCP: provides extensible base functionality without requiring modification
type BaseRepository struct {
	db *sql.DB
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(db *sql.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

// GetQueries returns sqlc queries with proper session handling
// Following DRY: centralizes query creation logic
func (r *BaseRepository) GetQueries(ctx context.Context) *sqlc.Queries {
	return GetQueriesWithSession(ctx, r.db)
}

// GetDB returns the database connection
// Following encapsulation: provides controlled access to database
func (r *BaseRepository) GetDB() *sql.DB {
	return r.db
}
