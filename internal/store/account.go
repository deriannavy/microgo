package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Account struct {
	Id        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	CreatedAt string   `json:"created_at"`
	IsActive  bool     `json:"is_active"`
	RoleId    int64    `json:"role_id"`
	Role      Role     `json:"role"`
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

func (p *password) Compare(text string) error {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(text))
}

type AccountStore struct {
	db *sql.DB
}

func (s *AccountStore) Activate(ctx context.Context, token string) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		account, err := s.GetAccountByToken(ctx, tx, token)
		if err != nil {
			return err
		}

		account.IsActive = true
		if err := s.Update(ctx, tx, account); err != nil {
			return err
		}

		if err := s.DeleteAccountConfirmation(ctx, tx, account.Id); err != nil {
			return err
		}

		return nil
	})
}
func (s *AccountStore) Create(ctx context.Context, tx *sql.Tx, account *Account) error {
	query := `
		INSERT INTO account (username, password, email, role_id) VALUES ($1, $2, $3,(SELECT id FROM roles WHERE name = $4)) RETURNING id;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	role := account.Role.Name
	if role == "" {
		role = "user"
	}

	err := tx.QueryRowContext(
		ctx,
		query,
		account.Username,
		account.Password.hash,
		account.Email,
		role,
	).Scan(
		&account.Id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *AccountStore) CreateAndConfirm(ctx context.Context, account *Account, token string, expiry time.Duration) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.Create(ctx, tx, account); err != nil {
			return err
		}

		if err := s.CreateAccountConfirmation(ctx, tx, token, expiry, account.Id); err != nil {
			return err
		}
		return nil
	})
}

func (s *AccountStore) CreateAccountConfirmation(ctx context.Context, tx *sql.Tx, token string, expiry time.Duration, accountId int64) error {
	query := `INSERT INTO account_confirmation (token, account_id, expiry) VALUES ($1, $2, $3);`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, accountId, time.Now().Add(expiry))
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "account_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "account_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}
	return nil
}

func (s *AccountStore) GetById(ctx context.Context, accountId int64) (*Account, error) {
	query := `
		SELECT 
		    a.id, a.username, a.password, a.email, a.created_at, 
		    r.id, r.name, r.description, r.level
		FROM
		    account a
		INNER JOIN role r ON 
			r.id = a.role_id
		WHERE 
			a.id = $1 and a.is_active = true;
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
		&account.Password.hash,
		&account.Email,
		&account.CreatedAt,
		&account.Role.Id,
		&account.Role.Name,
		&account.Role.Description,
		&account.Role.Level,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return account, nil

}

func (s *AccountStore) GetByEmail(ctx context.Context, email string) (*Account, error) {
	query := `
		SELECT id, username, password, email, created_at FROM account WHERE email = $1 and is_active = true;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	account := &Account{}
	err := s.db.QueryRowContext(
		ctx,
		query,
		email,
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

func (s *AccountStore) DeleteAccountAndConfirmations(ctx context.Context, accountId int64) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.Delete(ctx, tx, accountId); err != nil {
			return err
		}

		if err := s.DeleteAccountConfirmation(ctx, tx, accountId); err != nil {
			return err
		}

		return nil

	})
}

func (s *AccountStore) GetAccountByToken(ctx context.Context, tx *sql.Tx, token string) (*Account, error) {
	query := `
		SELECT 
		    a.id, a.username, a.email, a.created_at, a.is_active 
		FROM 
		    account a 
		INNER JOIN account_confirmation ac ON 
			a.id = ac.account_id
		WHERE 
		    ac.token = $1 AND ac.expiry > $2;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	account := &Account{}

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	err := tx.QueryRowContext(
		ctx,
		query,
		hashToken,
		time.Now(),
	).Scan(
		&account.Id,
		&account.Username,
		&account.Email,
		&account.CreatedAt,
		&account.IsActive,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return account, nil
}

func (s *AccountStore) Update(ctx context.Context, tx *sql.Tx, account *Account) error {
	query := `UPDATE account SET username = $1, email = $2, is_active = $3 WHERE id = $4;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, account.Username, account.Email, account.IsActive, account.Id)
	if err != nil {
		return err
	}

	return nil
}

func (s *AccountStore) Delete(ctx context.Context, tx *sql.Tx, id int64) error {
	query := `DELETE FROM account WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *AccountStore) DeleteAccountConfirmation(ctx context.Context, tx *sql.Tx, accountId int64) error {
	query := `DELETE FROM account_confirmation WHERE account_id = $1;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, accountId)

	if err != nil {
		return err
	}

	return nil
}
