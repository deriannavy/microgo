package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("Not found")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Account interface {
		Create(context.Context, *Account) error
		CreateAndConfirm(context.Context, *Account, string) error
		GetById(context.Context, int64) (*Account, error)
	}
	Transaction interface {
		Create(context.Context, *Transaction) error
		GetById(context.Context, int64) (*Transaction, error)
		GetByAccountId(context.Context, int64, PaginatedTransactionQuery) ([]TransactionWithMetadata, error)
		Update(context.Context, *Transaction) error
		Delete(context.Context, int64) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Account:     &AccountStore{db},
		Transaction: &TransactionStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return nil
	}

	return tx.Commit()
}
