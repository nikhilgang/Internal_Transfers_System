package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/InternalTransfer/internal/apperror"
	"github.com/InternalTransfer/internal/model"
	"github.com/shopspring/decimal"
)

type AccountRepository struct {
	pool *pgxpool.Pool
}

func NewAccountRepository(pool *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{pool: pool}
}

func (r *AccountRepository) Create(ctx context.Context, accountID int64, initialBalance decimal.Decimal) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO accounts (account_id, balance) VALUES ($1, $2)`, accountID, initialBalance)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return &apperror.ErrConflict{Entity: "account", ID: accountID}
		}
		return fmt.Errorf("inserting account: %w", err)
	}
	return nil
}

func (r *AccountRepository) GetByID(ctx context.Context, accountID int64) (*model.Account, error) {
	var a model.Account
	err := r.pool.QueryRow(ctx,
		`SELECT account_id, balance, created_at, updated_at FROM accounts WHERE account_id = $1`,
		accountID,
	).Scan(&a.AccountID, &a.Balance, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &apperror.ErrNotFound{Entity: "account", ID: accountID}
		}
		return nil, fmt.Errorf("querying account: %w", err)
	}
	return &a, nil
}

func (r *AccountRepository) GetByIDForUpdate(ctx context.Context, tx pgx.Tx, accountID int64) (*model.Account, error) {
	var a model.Account
	err := tx.QueryRow(ctx,
		`SELECT account_id, balance, created_at, updated_at FROM accounts WHERE account_id = $1 FOR UPDATE`,
		accountID,
	).Scan(&a.AccountID, &a.Balance, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &apperror.ErrNotFound{Entity: "account", ID: accountID}
		}
		return nil, fmt.Errorf("locking account: %w", err)
	}
	return &a, nil
}

func (r *AccountRepository) UpdateBalance(ctx context.Context, tx pgx.Tx, accountID int64, newBalance decimal.Decimal) error {
	tag, err := tx.Exec(ctx, `UPDATE accounts SET balance = $1, updated_at = NOW() WHERE account_id = $2`, newBalance, accountID)
	if err != nil {
		return fmt.Errorf("updating balance: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return &apperror.ErrNotFound{Entity: "account", ID: accountID}
	}
	return nil
}
