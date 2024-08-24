package postgres

import (
	"context"

	"gorm.io/gorm"
)

// TransactionManager defines an interface for managing database transactions.
// It includes methods to begin a transaction, commit it, or roll it back.
type TransactionManager interface {
	// Begins a transaction and returns a new context with the transaction.
	Begin(ctx context.Context) (context.Context, error)
	// Commits the current transaction associated with the context.
	Commit(ctx context.Context) error
	// Rolls back the current transaction associated with the context.
	Rollback(ctx context.Context) error
}

// transactionManagerImpl is a concrete implementation of the TransactionManager interface.
// It holds a reference to the gorm.DB instance, which is used to interact with the database.
type transactionManagerImpl struct {
	db *gorm.DB
}

// NewTransactionManager creates a new instance of transactionManagerImpl.
// It accepts a gorm.DB instance and returns it as a TransactionManager interface.
func NewTransactionManager(db *gorm.DB) TransactionManager {
	return &transactionManagerImpl{db}
}

// Begin starts a new transaction and stores the transaction in the context.
// It returns a new context with the transaction or an error if the transaction fails to start.
func (tm *transactionManagerImpl) Begin(ctx context.Context) (context.Context, error) {
	tx := tm.db.Begin()
	if tx.Error != nil {
		return ctx, tx.Error
	}

	return context.WithValue(ctx, dbKey, tx), nil
}

// Commit commits the transaction associated with the context.
// It retrieves the transaction from the context and commits it to the database.
func (tm *transactionManagerImpl) Commit(ctx context.Context) error {
	tx, ok := ctx.Value(dbKey).(*gorm.DB)
	if !ok {
		return nil
	}
	return tx.Commit().Error
}

// Rollback rolls back the transaction associated with the context.
// It retrieves the transaction from the context and rolls it back.
func (tm *transactionManagerImpl) Rollback(ctx context.Context) error {
	tx, ok := ctx.Value(dbKey).(*gorm.DB)
	if !ok {
		return nil
	}
	return tx.Rollback().Error
}
