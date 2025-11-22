package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/repositories"
)

// TransactionManager handles database transactions
type TransactionManager struct {
	pool repositories.ConnectionPool
}

func NewTransactionManager(pool repositories.ConnectionPool) *TransactionManager {
	return &TransactionManager{pool: pool}
}

// ExecuteInTransaction executes a function within a transaction
func (tm *TransactionManager) ExecuteInTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := tm.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p) // re-throw panic after rollback
		}
	}()

	if err := fn(tx); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return fmt.Errorf("rollback error: %v (original error: %w)", rollbackErr, err)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetTx returns a transaction that can be passed to repositories
func (tm *TransactionManager) GetTx(ctx context.Context) (pgx.Tx, error) {
	return tm.pool.Begin(ctx)
}
