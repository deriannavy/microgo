package store

import (
	"context"
	"database/sql"
)

type Account struct {
	Id        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}

type AccountStore struct {
	db *sql.DB
}

func (s *AccountStore) Create(ctx context.Context, account *Account) error {
	query := `
		INSERT INTO account (username, password, email) VALUES ($1, $2, $3) RETURNING id;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		account.Username,
		account.Password,
		account.Email,
	).Scan(
		&account.Id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *AccountStore) GetById(ctx context.Context, accountId int64) (*Account, error) {
	query := `
		SELECT id, username, password, email, created_at FROM account WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	account := &Account{}
	err := s.db.QueryRowContext(
		ctx,
		query,
		accountId,
	).Scan(
		&account.Id,
		&account.Username,
		&account.Password,
		&account.Email,
		&account.CreatedAt,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return account, nil

}
