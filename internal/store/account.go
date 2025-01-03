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

func (s *AccountStore) Register(ctx context.Context, account *Account) error {
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
