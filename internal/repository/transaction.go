package repository

import (
	"context"

	"github.com/SoraDaibu/go-clean-starter/internal/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Transaction represents an interface for grouping business procedures.
//
// Usage is similar to RDBMS transactions.
// However, transactions at the Usecase layer are business knowledge,
// and are not necessarily used only with RDBMS repositories.
type Transaction interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type dbTransaction struct {
	pool *pgxpool.Pool
}

func NewTransaction(pool *pgxpool.Pool) Transaction {
	return &dbTransaction{pool: pool}
}

func (tx *dbTransaction) Do(
	ctx context.Context,
	fn func(context.Context) error,
) error {
	t, err := tx.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer t.Rollback(ctx)

	if err := fn(SetSession(ctx, t)); err != nil {
		return err
	}

	return t.Commit(ctx)
}

type key struct{ value string }

var _contextKeyTx = &key{"_contextKeyTx"}

// GetSessionOr returns the current transaction or the fallback pool.
// It'll use a transaction if it exists, otherwise it'll use the fallback pool.
func GetSessionOr(ctx context.Context, fallback *pgxpool.Pool) sqlc.DBTX {
	if tx := ctx.Value(_contextKeyTx); tx != nil {
		if t, ok := tx.(pgx.Tx); ok {
			return t
		}
	}

	return fallback
}

// SetSession sets the current transaction to the context.
func SetSession(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, _contextKeyTx, tx)
}

// GetQueriesWithSession returns a new sqlc.Queries instance with the current transaction or the fallback pool.
func GetQueriesWithSession(ctx context.Context, fallback *pgxpool.Pool) *sqlc.Queries {
	return sqlc.New(GetSessionOr(ctx, fallback))
}
