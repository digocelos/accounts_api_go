package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/digocelo/account-api/internal/account"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepo struct {
	pool *pgxpool.Pool
}

func NewAccountRepo(pool *pgxpool.Pool) *AccountRepo {
	return &AccountRepo{pool: pool}
}

func (r *AccountRepo) InsertIdempotentByDocument(ctx context.Context, a account.Account) (bool, error) {
	const q = `
	INSERT INTO accounts (id, document, name, email, version, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT (document) DO NOTHING;
	`
	ct, err := r.pool.Exec(ctx, q, a.ID, a.Document, a.Name, a.Email, a.Version, a.CreatedAt, a.UpdatedAt)

	if err != nil {
		return false, fmt.Errorf("insert account: %w", err)
	}
	return ct.RowsAffected() == 1, nil
}

func (r *AccountRepo) GetById(ctx context.Context, id string) (account.Account, error) {
	const q = `
	SELECT id, document, name, email, version, created_at, updated_at
	  FROM accounts
	 WHERE  id = $1
	`

	var a account.Account

	err := r.pool.QueryRow(ctx, q, id).Scan(
		&a.ID, &a.Document, &a.Name, &a.Email, &a.Version, &a.CreatedAt, &a.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return account.Account{}, account.ErrNotFound
		}
		return account.Account{}, fmt.Errorf("get by id: %w", err)
	}
	return a, nil
}

func (r *AccountRepo) GetByDocument(ctx context.Context, document string) (account.Account, error) {
	const q = `
	SELECT id, document, version, name, email, created_at, updated_at
	  FROM accounts
	 WHERE document = $1`

	var a account.Account
	err := r.pool.QueryRow(ctx, q, document).Scan(
		&a.ID, &a.Document, &a.Version, &a.Name, &a.Email, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return account.Account{}, account.ErrNotFound
		}
		return account.Account{}, fmt.Errorf("get by document: %w", err)
	}
	return a, nil
}

func (r *AccountRepo) UpdateWithOptimisticLock(ctx context.Context, id string, expectedVersion int, name *string, email *string) (account.Account, error) {
	const q = `
	UPDATE account SET
	  name = &COALESCE($1, name),
	  email = COALESCE($2, email),
	  version = version + 1,
	  updated_at = now()
	WHERE id = $3 and version = $4
	RETURNING id, document, name, email, version, created_at, updated_at;
	`

	var a account.Account

	err := r.pool.QueryRow(ctx, q, name, email, id, expectedVersion).Scan(
		&a.ID, &a.Document, &a.Name, &a.Email, &a.Version, &a.CreatedAt, &a.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Could be not found or conflict; we'll differentiate:
			_, getErr := r.GetById(ctx, id)
			if getErr != nil {
				return account.Account{}, account.ErrConfict
			}
			if errors.Is(err, account.ErrNotFound) {
				return account.Account{}, account.ErrNotFound
			}
			return account.Account{}, fmt.Errorf("update conflict check: %w", getErr)
		}
		return account.Account{}, fmt.Errorf("update account: %w", err)
	}
	return a, nil
}
