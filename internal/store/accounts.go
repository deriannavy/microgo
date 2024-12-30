package store

import (
	"context"
	"database/sql"
)

type Account struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}

type AccountStore struct {
	db *sql.DB
}

func (s *AccountStore) Register(ctx context.Context, account *Account) error {
	query := `
		INSERT INTO account (username, password, email) VALUES ($1, $2, $3) RETURNING id;
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		account.Username,
		account.Password,
		account.Email,
	).Scan(
		&account.ID,
	)
	if err != nil {
		return err
	}
	return nil
}
