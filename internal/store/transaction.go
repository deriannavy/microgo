package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Transaction struct {
	Id        int64  `json:"id"`
	AccountId int64  `json:"account_id"`
	Date      string `json:"date"`
	Amount    int32  `json:"amount"`
	// accountOut
	// accountIN
	Place       string   `json:"place"`
	Description string   `json:"description"`
	Tag         []string `json:"tag"`
	Version     int32    `json:"version"`
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
		transaction.Id,
		transaction.AccountId,
		transaction.Date,
		transaction.Amount,
		transaction.Place,
		transaction.Description,
		pq.Array(transaction.Tag),
	).Scan(
		&transaction.Id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *TransactionStore) GetById(ctx context.Context, id int64) (*Transaction, error) {
	query := `SELECT id, account_id, date, amount, place, description, tag, version FROM transaction WHERE id = $1;`

	var transaction Transaction
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&transaction.Id,
		&transaction.AccountId,
		&transaction.Date,
		&transaction.Amount,
		&transaction.Place,
		&transaction.Description,
		pq.Array(&transaction.Tag),
		&transaction.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &transaction, nil
}

func (s *TransactionStore) Update(ctx context.Context, transaction *Transaction) error {
	query := `
		UPDATE 
			transaction
		SET
			date = $1
			amount = $2
			place = $3
			description = $4
			tag = $5
		WHERE
			id = $6 AND 
			version = $7
		RETURNING 
			version
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		transaction.Date,
		transaction.Amount,
		transaction.Place,
		transaction.Description,
		pq.Array(transaction.Tag),
		transaction.Id,
		transaction.Version,
	).Scan(
		&transaction.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}

func (s *TransactionStore) Delete(ctx context.Context, transactionId int64) error {
	query := `DELETE FROM transaction WHERE id = $1`

	res, err := s.db.ExecContext(ctx, query, transactionId)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}
