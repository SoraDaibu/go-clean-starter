package repository

import (
	"context"

	"github.com/SoraDaibu/go-clean-starter/internal/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

// BaseRepository provides common database operations
// Following OCP: provides extensible base functionality without requiring modification
type BaseRepository struct {
	// pgxpool.Pool is safe to use from multiple goroutines simultaneously
	pool *pgxpool.Pool
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(pool *pgxpool.Pool) *BaseRepository {
	return &BaseRepository{pool: pool}
}

// GetQueries returns sqlc queries with proper session handling
// Following DRY: centralizes query creation logic
func (r *BaseRepository) GetQueries(ctx context.Context) *sqlc.Queries {
	return GetQueriesWithSession(ctx, r.pool)
}

// GetPool returns the connection pool
// Following encapsulation: provides controlled access to connection pool
func (r *BaseRepository) GetPool() *pgxpool.Pool {
	return r.pool
}
