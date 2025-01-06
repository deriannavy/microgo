package store

import (
	"context"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	Id        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	CreatedAt string   `json:"created_at"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.text = &text
	p.hash = hash

	return nil
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

func (s *AccountStore) CreateAndConfirm(ctx context.Context, account *Account, token string) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.Create(ctx, account); err != nil {
			return err
		}
		return nil
	})
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
