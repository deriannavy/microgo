package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Account interface {
		Register(context.Context, *Account) error
	}
	Transaction interface {
		Create(context.Context, *Transaction) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Account:     &AccountStore{db},
		Transaction: &TransactionStore{db},
	}
}
