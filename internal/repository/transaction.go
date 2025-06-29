package repository

import (
	"context"
	"database/sql"

	"github.com/SoraDaibu/go-clean-starter/internal/sqlc"
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
	db *sql.DB
}

func NewTransaction(db *sql.DB) Transaction {
	return &dbTransaction{db: db}
}

func (tx *dbTransaction) Do(
	ctx context.Context,
	fn func(context.Context) error,
) error {
	t, err := tx.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer t.Rollback()

	if err := fn(SetSession(ctx, t)); err != nil {
		return err
	}

	return t.Commit()
}

type key struct{ value string }

var _contextKeyTx = &key{"_contextKeyTx"}

// GetSessionOr returns the current transaction or the fallback database.
// It'll use a transaction if it exists, otherwise it'll use the fallback database.
func GetSessionOr(ctx context.Context, fallback *sql.DB) sqlc.DBTX {
	if tx := ctx.Value(_contextKeyTx); tx != nil {
		if t, ok := tx.(*sql.Tx); ok {
			return t
		}
	}

	return fallback
}

// SetSession sets the current transaction to the context.
func SetSession(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, _contextKeyTx, tx)
}

// GetQueriesWithSession returns a new sqlc.Queries instance with the current transaction or the fallback database.
func GetQueriesWithSession(ctx context.Context, fallback *sql.DB) *sqlc.Queries {
	return sqlc.New(GetSessionOr(ctx, fallback))
}
