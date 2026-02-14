package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type TransactionRepository struct {
	pool *pgxpool.Pool
}

func NewTransactionRepository(pool *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{pool: pool}
}

func (r *TransactionRepository) Create(ctx context.Context, tx pgx.Tx, sourceID, destID int64, amount decimal.Decimal) error {
	_, err := tx.Exec(ctx, `INSERT INTO transactions (source_account_id, destination_account_id, amount) VALUES ($1, $2, $3)`, sourceID, destID, amount)
	if err != nil {
		return fmt.Errorf("inserting transaction: %w", err)
	}
	return nil
}
