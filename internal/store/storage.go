package store

import (
	"context"
	"database/sql"
)

var (
	ErrNotFound = error.New("Not found")
)

type Storage struct {
	Account interface {
		Register(context.Context, *Account) error
	}
	Transaction interface {
		Create(context.Context, *Transaction) error
		GetById(context.Context, int64) (*Transaction, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Account:     &AccountStore{db},
		Transaction: &TransactionStore{db},
	}
}
