package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// DBTX defines the database operations interface
type DBTX interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// ConnectionPool defines the complete database pool interface
type ConnectionPool interface {
	DBTX
	Begin(ctx context.Context) (pgx.Tx, error)
	Close()
	Ping(ctx context.Context) error
}

// TransactionManager interface (as before)
type TransactionManager interface {
	ExecuteInTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error
	GetTx(ctx context.Context) (pgx.Tx, error)
}

// TxProvider provides transaction capability
type TxProvider[T any] interface {
	WithTx(tx pgx.Tx) T // Set transaction for the repository
}
