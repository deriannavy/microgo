package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("Not found")
	ErrDuplicateEmail    = errors.New("Duplicate email")
	ErrDuplicateUsername = errors.New("Duplicate username")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Account interface {
		Create(ctx context.Context, tx *sql.Tx, account *Account) error
		CreateAndConfirm(ctx context.Context, account *Account, token string, expiry time.Duration) error
		Activate(context.Context, string) error
		GetById(context.Context, int64) (*Account, error)
		GetByEmail(ctx context.Context, email string) (*Account, error)
		Delete(ctx context.Context, tx *sql.Tx, accountId int64) error
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
