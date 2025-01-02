package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotFound = errors.New("Not found")
)

type Storage struct {
	Account interface {
		Register(context.Context, *Account) error
	}
	Transaction interface {
		Create(context.Context, *Transaction) error
		GetById(context.Context, int64) (*Transaction, error)
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
