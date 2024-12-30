package store

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
)

type Transaction struct {
	ID        int64  `json:"id"`
	AccountID int64  `json:"account_id"`
	Date      string `json:"date"`
	Amount    int32  `json:"amount"`
	// accountOut
	// accountIN
	Place       string   `json:"place"`
	Description string   `json:"description"`
	Tag         []string `json:"tag"`
}

type TransactionStore struct {
	db *sql.DB
}

func (s *TransactionStore) Create(ctx context.Context, transaction *Transaction) error {
	query := `
		INSERT INTO transaction (id, account_id, date, amount, place, description, tag)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;
	`
	err := s.db.QueryRowContext(
		ctx,
		query,
		transaction.ID,
		transaction.AccountID,
		transaction.Date,
		transaction.Amount,
		transaction.Place,
		transaction.Description,
		pq.Array(transaction.Tag),
	).Scan(
		&transaction.ID,
	)

	if err != nil {
		return err
	}

	return nil
}
